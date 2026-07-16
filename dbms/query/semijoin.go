// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"sync/atomic"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/set"
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
}

type semiJoinApproach struct {
	req2 Require
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

func (sj *SemiJoin) setApproach(req Require, approach any, tran QueryTran) {
	ap := approach.(*semiJoinApproach)
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
	sj.source1 = SetApproach(sj.source1, req, tran)
	sj.source2 = SetApproach(sj.source2, ap.req2, tran)
	sj.header = sj.source1.Header()
}

func (sj *SemiJoin) Rewind() {
	sj.source1.Rewind()
	sj.source2.Rewind()
}

func (sj *SemiJoin) Get(th *Thread, dir Dir) Row {
	defer func(t uint64) { sj.tget += tsc.Read() - t }(tsc.Read())
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
	sj.source1.Select(sels)
	sj.source2.Select(nil)
	sj.Rewind()
}

func (sj *SemiJoin) Lookup(th *Thread, sels Sels) Row {
	sj.nlooks++
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
