// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"
	"sync/atomic"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/shmap"
	"github.com/apmckinlay/gsuneido/util/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/tsc"
)

type SemiJoin struct {
	Query2
	qt          QueryTran
	st          *SuTran
	lookupCache lookupCache
	by          []string
	prevFixed1  Fixed
	prevFixed2  Fixed
	joinType    joinType
	optimized   bool
	reverse     bool
	sels1       Sels // from incoming Select, added to source1 probe in reverse mode

	revState      state
	indexed       bool
	row2          Row
	lookupRow     Row
	dedup         *mapType
	warned        bool
	derived       int
	derivedWarned bool
	th            *Thread
}

type semiJoinApproach struct {
	reverse bool
	req1    Require
	req2    Require
}

var (
	semiJoinCacheProbes atomic.Int64
	semiJoinCacheMisses atomic.Int64
	semiJoin11Count     atomic.Int64
	semiJoin1nCount     atomic.Int64
	semiJoinn1Count     atomic.Int64
	semiJoinnnCount     atomic.Int64
)

var _ = AddInfo("query.semijoin.cacheProbes", &semiJoinCacheProbes)
var _ = AddInfo("query.semijoin.cacheMisses", &semiJoinCacheMisses)
var _ = AddInfo("query.semijoin.11", &semiJoin11Count)
var _ = AddInfo("query.semijoin.1n", &semiJoin1nCount)
var _ = AddInfo("query.semijoin.n1", &semiJoinn1Count)
var _ = AddInfo("query.semijoin.nn", &semiJoinnnCount)

func NewSemiJoin(src1, src2 Query, by []string, t QueryTran) *SemiJoin {
	b := set.Intersect(src1.Columns(), src2.Columns())
	if len(b) == 0 {
		panic("semijoin: common columns required")
	}
	if by == nil {
		by = b
	} else if !set.Subset(b, by) {
		panic("semijoinjoin: by must be a subset of the common columns")
	}
	sj := &SemiJoin{qt: t, st: MakeSuTran(t), by: by}
	sj.lookupCache.SetCounters(&semiJoinCacheProbes, &semiJoinCacheMisses)
	sj.joinType = getJoinType(by, src1, src2)
	sj.source1, sj.source2 = src1, src2
	sj.header = src1.Header()
	sj.keys = src1.Keys()
	sj.indexes = src1.Indexes()
	sj.fixed = src1.Fixed()
	sj.setNrows(sj.getNrows())
	sj.rowSiz.Set(src1.rowSize())
	sj.fast1.Set(src1.fastSingle())
	return sj
}

func (sj *SemiJoin) With(src1, src2 Query) *SemiJoin {
	q := NewSemiJoin(src1, src2, sj.by, sj.qt)
	q.prevFixed1 = sj.prevFixed1
	q.prevFixed2 = sj.prevFixed2
	return q
}

func (sj *SemiJoin) String() string {
	op := "semijoin"
	if sj.optimized {
		if sj.reverse {
			op += "-rev"
		}
		op += " " + sj.joinType.String()
	} else if sj.joinType == many_to_many {
		op += " /*MANY TO MANY*/"
	}
	return op + " by" + str.Join("(,)", sj.by)
}

func (sj *SemiJoin) SetTran(t QueryTran) {
	sj.qt = t
	sj.st = MakeSuTran(t)
	sj.lookupCache.Reset()
	sj.Query2.SetTran(t)
	// don't need to clear dedup since it's only used in ReadMode
	// which doesn't use SetTran
}

func (sj *SemiJoin) getNrows() (int, int) {
	n1, p1 := sj.source1.Nrows()
	n2, p2 := sj.source2.Nrows()
	if n1 == 0 || n2 == 0 || p2 == 0 {
		return 0, p1
	}
	est := n1 * n2 / p2
	if est > n1 {
		est = n1
	}
	return est, p1
}

func (sj *SemiJoin) Transform() Query {
	sj.optimized = true
	src1 := sj.source1.Transform()
	if _, ok := src1.(*Nothing); ok {
		return NewNothing(sj)
	}
	src2 := sj.source2.Transform()
	if _, ok := src2.(*Nothing); ok {
		return NewNothing(sj)
	}
	fix1, fix2 := src1.Fixed(), src2.Fixed()
	if !fix1.Equal(sj.prevFixed1) || !fix2.Equal(sj.prevFixed2) {
		src1 = copyFixed(fix2, fix1, src1, sj.by, sj.qt)
		src2 = copyFixed(fix1, fix2, src2, sj.by, sj.qt)
		sj.prevFixed1, sj.prevFixed2 = fix1, fix2
	}
	if src1 != sj.source1 || src2 != sj.source2 {
		return sj.With(src1, src2).Transform()
	}
	return sj
}

func (sj *SemiJoin) optimize(mode Mode, req Require) (Cost, Cost, any) {
	fwdFix, fwdVar, fwdApp := sj.optimizeForward(mode, req)
	revFix, revVar, revApp := sj.optimizeReverse(mode, req)
	if revFix+revVar < fwdFix+fwdVar {
		return revFix, revVar, revApp
	}
	return fwdFix, fwdVar, fwdApp
}

func (sj *SemiJoin) optimizeForward(mode Mode, req Require) (Cost, Cost, any) {
	fixcost1, varcost1 := Optimize(sj.source1, mode, req)
	nrows1, _ := sj.source1.Nrows()
	nrows2, _ := sj.source2.Nrows()
	nseeks := req.SeekCount(nrows1)
	if nseeks <= 0 {
		nseeks = 1
	}
	var req2 Require
	if sj.joinType.toOne() {
		req2 = UniqueReq(sj.by, nseeks)
	} else {
		frac2 := min(float32(1), float32(nseeks)/float32(max(1, nrows2)))
		req2 = GroupReq(sj.by, frac2, nseeks)
	}
	fixcost2, varcost2 := Optimize(sj.source2, mode, req2)
	if fixcost2+varcost2 >= impossible {
		return impossible, impossible, nil
	}
	return fixcost1 + fixcost2, varcost1 + varcost2,
		&semiJoinApproach{req2: req2}
}

// optimizeReverse evaluates the cost of running the semijoin in "reverse" mode,
// where source2 drives the iteration instead of source1. In this strategy:
//   - Source2 is iterated (or seeked) according to the incoming request
//   - For each source2 row, the "by" columns are projected and used to lookup
//     or select matching rows from source1
//   - Source1 rows that match are returned as the semijoin result
//
// This can be more efficient than forward mode when source2 is smaller or has
// better access paths (indexes) for the incoming request pattern. However, it
// requires duplicate elimination for one_to_many and many_to_many join types,
// since multiple source2 rows can share the same "by" values, which would
// otherwise cause the same source1 row to be returned multiple times.
//
// The duplicate elimination uses a map (like Project) to track which source2
// rows have been seen, ensuring each distinct "by" value group is processed
// only once. This adds overhead proportional to the number of source2 rows
// scanned, so reverse mode is only chosen when its total cost (including
// deduplication) is lower than forward mode.
func (sj *SemiJoin) optimizeReverse(mode Mode, req Require) (Cost, Cost, any) {
	nrows2, _ := sj.source2.Nrows()
	fixcost2, varcost2 := Optimize(sj.source2, mode, req)
	if fixcost2+varcost2 >= impossible {
		return impossible, impossible, nil
	}
	nseeks := req.SeekCount(nrows2)
	if nseeks <= 0 {
		nseeks = 1
	}
	var req1 Require
	if sj.joinType == one_to_one || sj.joinType == one_to_many {
		req1 = UniqueReq(sj.by, nseeks)
	} else { // many_to_one, many_to_many
		frac1 := min(float32(1), float32(nseeks)/float32(max(1, nrows2))) // ???
		req1 = GroupReq(sj.by, frac1, nseeks)
	}
	fixcost1, varcost1 := Optimize(sj.source1, mode, req1)
	if fixcost1+varcost1 >= impossible {
		return impossible, impossible, nil
	}
	varcost := varcost1 + varcost2
	if sj.joinType == one_to_many || sj.joinType == many_to_many {
		if mode != ReadMode || nrows2 > mapThreshold {
			return impossible, impossible, nil
		}
		// duplicate elimination cost, proportional to source2 rows scanned
		varcost += Cost(nseeks) * mapCost
	}
	return fixcost1 + fixcost2 + outOfOrder, varcost,
		&semiJoinApproach{reverse: true, req1: req1, req2: req}
}

func (sj *SemiJoin) setApproach(req Require, approach any, tran QueryTran) {
	ap := approach.(*semiJoinApproach)
	sj.reverse = ap.reverse
	switch sj.joinType {
	case one_to_one:
		semiJoin11Count.Add(1)
	case one_to_many:
		semiJoin1nCount.Add(1)
	case many_to_one:
		semiJoinn1Count.Add(1)
	case many_to_many:
		semiJoinnnCount.Add(1)
	}
	if sj.reverse {
		sj.source1 = SetApproach(sj.source1, ap.req1, tran)
		sj.source2 = SetApproach(sj.source2, ap.req2, tran)
	} else {
		sj.source1 = SetApproach(sj.source1, req, tran)
		sj.source2 = SetApproach(sj.source2, ap.req2, tran)
	}
	sj.header = sj.source1.Header()
}

func (sj *SemiJoin) Rewind() {
	sj.source1.Rewind()
	sj.source2.Rewind()
	sj.row2 = nil
	sj.lookupRow = nil
	sj.revState = rewound
	// NOTE: sj.dedup and sj.indexed are NOT reset here - like Project,
	// the dedup map remains valid for the underlying data across Rewind,
	// only the cursor position (state) is reset.
}

func (sj *SemiJoin) Get(th *Thread, dir Dir) Row {
	defer func(t uint64) { sj.tget += tsc.Read() - t }(tsc.Read())
	if sj.reverse {
		if sj.revState == eof {
			return nil
		}
		row := sj.getReverse(th, dir)
		if row != nil {
			sj.revState = within
		} else {
			sj.revState = eof
		}
		return row
	}
	for {
		row := sj.source1.Get(th, dir)
		if row == nil {
			return nil
		}
		if sj.source2Has(th, row, dir) {
			sj.ngets++
			return row
		}
	}
}

// getReverse iterates source2 and looks up/selects source1 by the by
// columns, returning source1 rows. For one_to_many and many_to_many, by is
// not a key of source2, so multiple source2 rows can share the same by
// values; those duplicates must be eliminated (like Project) before doing
// the lookup on source1, otherwise the same source1 row would be returned
// more than once. (one_to_one and many_to_one can't produce duplicates
// since by is a key of source2, so every source2 row is distinct by
// definition.)
//
// The duplicate elimination is keyed on source2's own rows (not the
// resulting source1 row), because - just like Project's source - each
// physical source2 row corresponds to exactly one position in source2's
// own bidirectional sequence, visited at most once per direction traversal.
// That lets us reuse Project's exact technique (canonical-row-by-identity)
// to keep Get(Prev) the precise reverse of Get(Next): duplicate elimination
// always keeps the row encountered first when scanning source2 in Next
// order. If the very first call after a Rewind is a Prev, we must build
// the map (scanning forward) before starting the actual Prev-directed scan.
func (sj *SemiJoin) getReverse(th *Thread, dir Dir) Row {
	sj.th = th
	defer func() { sj.th = nil }()
	needDedup := sj.joinType == one_to_many || sj.joinType == many_to_many
	if needDedup {
		if sj.dedup == nil {
			hdr2 := sj.source2.Header()
			hfn := func(k rowHash) uint64 { return k.hash }
			eqfn := func(x, y rowHash) bool {
				return x.hash == y.hash &&
					equalCols(x.row, y.row, hdr2, sj.by, sj.th, sj.st)
			}
			sj.dedup = shmap.NewMapFuncs[rowHash, struct{}](hfn, eqfn)
		}
		if sj.revState == rewound && dir == Prev && !sj.indexed {
			sj.buildDedup(th)
		}
	}
	for {
		if sj.row2 == nil {
			row2 := sj.source2.Get(th, dir)
			if row2 == nil {
				if dir == Next {
					sj.indexed = true
				}
				return nil
			}
			if needDedup {
				oldRow2, existed := sj.dedupRow(row2)
				if existed && !row2.SameAs(oldRow2) {
					continue // duplicate by values, not the canonical row2
				}
			}
			sj.row2 = row2
			sels := slc.With(sj.sels1, sj.projectRow2(th, sj.row2)...)
			if sj.joinType == one_to_one || sj.joinType == one_to_many {
				sj.lookupRow = sj.source1.Lookup(th, sels)
			} else { // many_to_one, many_to_many
				sj.source1.Select(sels)
			}
		}
		var row1 Row
		if sj.joinType == one_to_one || sj.joinType == one_to_many {
			row1 = sj.lookupRow
			sj.lookupRow = nil
			sj.row2 = nil
		} else { // many_to_one, many_to_many
			row1 = sj.source1.Get(th, dir)
			if row1 == nil {
				sj.row2 = nil
				continue
			}
		}
		if row1 != nil {
			sj.ngets++
			return row1
		}
	}
}

// buildDedup does a full forward (Next) scan of source2, recording the
// canonical (first encountered) row2 for each distinct by-values group.
// This establishes, once and for all, which physical source2 row
// represents each group, so that a subsequent scan in either direction
// (Get Next or Prev) can consistently decide which occurrences to keep
// or skip.
func (sj *SemiJoin) buildDedup(th *Thread) {
	for {
		row2 := sj.source2.Get(th, Next)
		if row2 == nil {
			break
		}
		sj.dedupRow(row2)
	}
	sj.source2.Rewind()
	sj.indexed = true
}

// dedupRow returns the old row2 and true if its by values already existed,
// else the new row2 and false
func (sj *SemiJoin) dedupRow(row2 Row) (Row, bool) {
	hdr2 := sj.source2.Header()
	rh := rowHash{row: row2, hash: hashCols(row2, hdr2, sj.by, sj.th, sj.st)}
	k, existed := sj.dedup.GetInit(rh)
	if existed {
		return k.row, true
	}
	if !sj.warned && sj.dedup.Size() > mapWarn {
		sj.warned = true
		Warning("semijoin-map large >", mapWarn)
	}
	sj.derived += row2.Derived()
	if !sj.derivedWarned && sj.derived > derivedWarn {
		sj.derivedWarned = true
		Warning("semijoin-map derived large >", derivedWarn,
			"average", sj.derived/sj.dedup.Size())
	}
	return row2, false
}

func (sj *SemiJoin) projectRow2(th *Thread, row Row) Sels {
	sels := make(Sels, len(sj.by))
	for i, col := range sj.by {
		sels[i] = Sel{col, row.GetRawVal(sj.source2.Header(), col, th, sj.st)}
	}
	return sels
}

func (sj *SemiJoin) source2Has(th *Thread, row Row, dir Dir) bool {
	sels := make(Sels, len(sj.by))
	for i, col := range sj.by {
		sels[i] = Sel{col, row.GetRawVal(sj.source1.Header(), col, th, sj.st)}
	}
	if sj.joinType == one_to_one {
		return sj.source2.Lookup(th, sels) != nil
	} else if sj.joinType == many_to_one {
		return sj.lookupCache.Lookup(th, sj.source2, sels, sj.st) != nil
	}
	sj.source2.Select(sels)
	return sj.source2.Get(th, dir) != nil
}

func (sj *SemiJoin) Select(sels Sels) {
	sj.nsels++
	if sj.reverse {
		// iterate source2, so apply by columns to it
		// and save the rest to add to source1 probe
		sj.sels1 = nil
		var sel2 Sels
		for _, sel := range sels {
			if slices.Contains(sj.by, sel.col) {
				sel2 = append(sel2, sel)
			} else {
				sj.sels1 = append(sj.sels1, sel)
			}
		}
		sj.source2.Select(sel2)
		if sj.dedup != nil {
			sj.dedup.Clear()
			sj.derived = 0
			sj.derivedWarned = false
			sj.warned = false
		}
		sj.indexed = false
	} else {
		// iterate source1, so apply all sels to source1
		sj.source1.Select(sels)
		sj.source2.Select(nil)
	}
	sj.Rewind()
}

func (sj *SemiJoin) Lookup(th *Thread, sels Sels) Row {
	sj.nlooks++
	if sj.reverse {
		return lookupViaSelectGet(sj, th, sels)
	}
	if sj.joinType == one_to_one || sj.joinType == many_to_one {
		row := sj.source1.Lookup(th, sels)
		if row == nil {
			return nil
		}
		sel2 := make(Sels, len(sj.by))
		for i, col := range sj.by {
			sel2[i] = Sel{col, row.GetRawVal(sj.source1.Header(), col, th, sj.st)}
		}
		var row2 Row
		if sj.joinType == one_to_one {
			row2 = sj.source2.Lookup(th, sel2)
		} else {
			row2 = sj.lookupCache.Lookup(th, sj.source2, sel2, sj.st)
		}
		if row2 == nil {
			return nil
		}
		return row
	}
	return lookupViaSelectGet(sj, th, sels)
}

func (sj *SemiJoin) Simple(th *Thread) []Row {
	rows1 := sj.source1.Simple(th)
	rows2 := sj.source2.Simple(th)
	dst := 0
	for _, row1 := range rows1 {
		for _, row2 := range rows2 {
			if sj.equalBy(th, row1, row2) {
				rows1[dst] = row1
				dst++
				break
			}
		}
	}
	return rows1[:dst]
}

func (sj *SemiJoin) equalBy(th *Thread, row1, row2 Row) bool {
	for _, col := range sj.by {
		if row1.GetRawVal(sj.source1.Header(), col, th, sj.st) !=
			row2.GetRawVal(sj.source2.Header(), col, th, sj.st) {
			return false
		}
	}
	return true
}
