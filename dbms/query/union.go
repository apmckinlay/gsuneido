// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"slices"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/slc"
	"github.com/apmckinlay/gsuneido/util/tsc"
)

var (
	unionMergeCount    atomic.Int64
	unionLookupCount   atomic.Int64
	unionDisjointCount atomic.Int64
	unionMergeDisjoint atomic.Int64
)

var _ = AddInfo("query.union.merge", &unionMergeCount)
var _ = AddInfo("query.union.lookup", &unionLookupCount)
var _ = AddInfo("query.union.disjoint", &unionDisjointCount)
var _ = AddInfo("query.union.merge-disjoint", &unionMergeDisjoint)

type Union struct {
	Compatible
	src2get   func(*Thread, Dir) Row
	src1get   func(*Thread, Dir) Row
	mergeCols []string
	row2      Row
	empty1    Row
	empty2    Row
	row1      Row
	strat     unionStrategy
	src2      bool
	prevDir   Dir
	src1      bool
	state
}

type unionApproach struct {
	keyIndex   []string // not necessarily a key if disjoint
	strat      unionStrategy
	reverse    bool
	req1, req2 Require
}

type unionStrategy int

const (
	// unionMerge is an ordered merge of source1 and source2
	unionMerge unionStrategy = iota + 2
	// unionLookup is source1 not in source2, followed by source2 (unordered).
	// Also used for disjoint, but without lookups.
	unionLookup
)

func (us unionStrategy) String() string {
	switch us {
	case unionMerge:
		return "merge"
	case unionLookup:
		return "lookup"
	default:
		return "unknown"
	}
}

func NewUnion(src1, src2 Query) *Union {
	u := &Union{Compatible: *newCompatible(src1, src2)}
	u.header = JoinHeaders(src1.Header(), src2.Header())
	u.indexes = u.getIndexes()
	u.setNrows(u.getNrows())
	u.rowSiz.Set((u.source1.rowSize() + u.source2.rowSize()) / 2)
	return u
}

func (u *Union) String() string {
	s := "union"
	if u.strat == 0 { // not optimized
		if u.disjoint == "" {
			s += " /*NOT DISJOINT*/"
		}
		return s
	}
	if u.disjoint != "" {
		s += "-disjoint(" + u.disjoint + ")"
	}
	switch u.strat {
	case unionMerge:
		s += "-merge" //+ str.Join("(,)", u.keyIndex)
	case unionLookup:
		if u.disjoint == "" {
			s += "-lookup"
		}
	}
	return s
}

func (u *Union) Keys() [][]string {
	if u.keys == nil {
		if u.disjoint == "" {
			u.keys = [][]string{u.allCols}
		} else {
			keys := u.keypairs()
			for i := range keys {
				// keypairs must ensure that appending is valid
				keys[i] = set.AddUnique(keys[i], u.disjoint)
			}
			u.keys = minimizeKeys(keys)
		}
		assert.That(u.keys != nil)
	}
	return u.keys
}

func (*Union) fastSingle() bool {
	return false
}

func (u *Union) getIndexes() [][]string {
	// lookup can read via any index
	// merge cannot but that will be handled by optimize
	idx1 := u.source1.Indexes()
	idx2 := u.source2.Indexes()
	if isEmptyKey(idx1) {
		return idx2
	} else if isEmptyKey(idx2) {
		return idx1
	}
	return set.UnionFn(idx1, idx2, slices.Equal)
}

func (u *Union) getNrows() (int, int) {
	n1, p1 := u.source1.Nrows()
	n2, p2 := u.source2.Nrows()
	return u.nrowsCalc(n1, n2), u.nrowsCalc(p1, p2)
}

func (u *Union) nrowsCalc(n1, n2 int) int {
	if u.disjoint != "" {
		return n1 + n2
	}
	min := max(n1, n2)     // smaller could be all duplicates
	max := n1 + n2         // could be no duplicates
	return (min + max) / 2 // estimate half way between
}

func (u *Union) Transform() Query {
	src1 := u.source1.Transform()
	src2 := u.source2.Transform()
	if _, ok := src1.(*Nothing); ok {
		// remove unnecessary Union
		return keepCols(src2, src1, u.header)
	}
	if _, ok := src2.(*Nothing); ok {
		// remove unnecessary Union
		return keepCols(src1, src2, u.header)
	}
	if src1 != u.source1 || src2 != u.source2 {
		return NewUnion(src1, src2)
	}
	return u
}

func keepCols(src, nothing Query, hdr *Header) Query {
	cols := set.Difference(nothing.Columns(), src.Columns())
	if len(cols) == 0 {
		return src
	}
	var empty ast.Expr = &ast.Constant{Val: EmptyStr}
	exprs := slc.Repeat(empty, len(cols))
	for i, col := range cols {
		if !hdr.HasField(col) {
			exprs[i] = nil
		}
	}
	// need to transform in case e.g. src is another extend
	return NewExtend(src, cols, exprs).Transform()
}

func (u *Union) Fixed() Fixed {
	if u.fixed == nil {
		u.fixed = u.getFixed()
		assert.That(u.fixed != nil)
	}
	return u.fixed
}

func (u *Union) getFixed() Fixed {
	fixed1 := u.source1.Fixed()
	fixed2 := u.source2.Fixed()
	fixed := make(Fixed, 0, len(fixed1)+len(fixed2))
	// add ones that are in both
	for _, f1 := range fixed1 {
		for _, f2 := range fixed2 {
			if f1.col == f2.col {
				fixed = append(fixed,
					Fix{col: f1.col, values: set.Union(f1.values, f2.values)})
				break
			}
		}
	}
	// fixed on columns that are only on one source
	// can treat the other source as fixed = ""
	cols2 := u.source2.Columns()
	emptyStr := []string{""}
	for _, f1 := range fixed1 {
		if !slices.Contains(cols2, f1.col) {
			fixed = append(fixed,
				Fix{col: f1.col, values: set.Union(f1.values, emptyStr)})
		}
	}
	cols1 := u.source1.Columns()
	for _, f2 := range fixed2 {
		if !slices.Contains(cols1, f2.col) {
			fixed = append(fixed,
				Fix{col: f2.col, values: set.Union(f2.values, emptyStr)})
		}
	}
	return fixed
}

// optimize ---------------------------------------------------------

func (u *Union) optimize(mode Mode, req Require) (Cost, Cost, any) {
	// try merge versus lookup
	fcMerge, vcMerge, appMerge := u.optMerge(mode, req)
	// The lookup strategy just concatenates source1 then source2, with no
	// merging. For ReqUnique/ReqGroup that is only valid when disjoint if
	// req.cols includes the disjoint column - otherwise the same req.cols
	// values can occur on both sides (they only differ by the disjoint
	// column) and would end up split into two non-adjacent groups/lookups
	// instead of being combined.
	if req.use == ReqNone ||
		(u.disjoint != "" && (req.use == ReqUnique || req.use == ReqGroup) &&
			slices.Contains(req.cols, u.disjoint)) {
		fcLookup, vcLookup, appLookup := u.optLookup(mode, req)
		if fcLookup+vcLookup < fcMerge+vcMerge {
			return fcLookup, vcLookup, appLookup
		}
	}
	return fcMerge, vcMerge, appMerge
}

func (u *Union) optLookup(mode Mode, req Require) (Cost, Cost, *unionApproach) {
	// try forward versus reverse
	fc, vc, app := u.optLookup2(mode, req)

	u.source1, u.source2 = u.source2, u.source1
	fcRev, vcRev, appRev := u.optLookup2(mode, req)
	u.source1, u.source2 = u.source2, u.source1
	fcRev += outOfOrder

	if fcRev+vcRev < fc+vc {
		appRev.reverse = true
		return fcRev, vcRev, appRev
	}
	return fc, vc, app
}

func (u *Union) optLookup2(mode Mode, req Require) (Cost, Cost, *unionApproach) {
	nrows1, _ := u.source1.Nrows()
	nseeks := req.SeekCount(nrows1)
	req1 := req
	if req.use == ReqUnique || req.use == ReqGroup {
		req1 = GroupReq(req.cols, req.SelectFrac(nrows1), int32(nseeks))
	}
	fc1, vc1 := Optimize(u.source1, mode, req1)
	if fc1+vc1 >= impossible {
		return impossible, impossible, nil
	}
	if u.disjoint != "" {
		fc2, vc2 := Optimize(u.source2, mode, req1)
		if fc2+vc2 >= impossible {
			return impossible, impossible, nil
		}
		return fc1 + fc2, vc1 + vc2,
			&unionApproach{strat: unionLookup, req1: req1, req2: req1}
	}
	// else not disjoint so we need lookups on source2
	req2 := UniqueReq(u.source2.Columns(), int32(nseeks))
	fc2, vc2 := Optimize(u.source2, mode, req2)
	if fc2+vc2 >= impossible {
		return impossible, impossible, nil
	}
	return fc1 + fc2, vc1 + vc2,
		&unionApproach{strat: unionLookup, req1: req1, req2: req2, keyIndex: req2.cols}
}

func (u *Union) optMerge(mode Mode, req Require) (Cost, Cost, *unionApproach) {
	// merge is symmetrical so we don't need to try forward and reverse
	if u.disjoint != "" {
		mr := OrderReq(req.cols, req.frac)
		fc1, vc1 := Optimize(u.source1, mode, mr)
		fc2, vc2 := Optimize(u.source2, mode, mr)
		if fc1+vc1 < impossible && fc2+vc2 < impossible {
			return fc1 + fc2, vc1 + vc2,
				&unionApproach{keyIndex: req.cols, strat: unionMerge, req1: mr, req2: mr}
		}
	}
	// Special case: if both sources have empty keys, allow unordered merge
	// This matches the old optMergeWithOrder behavior
	keys1 := u.source1.Keys()
	keys2 := u.source2.Keys()
	if isEmptyKey(keys1) && isEmptyKey(keys2) {
		// With empty keys, each source has at most 1 row, so union has at most 2 rows.
		// The merge will compare using allCols, which must satisfy the req.
		allCols := u.allCols
		if req.SatisfiedByWithFixed(allCols, u.source1.Fixed()) &&
			req.SatisfiedByWithFixed(allCols, u.source2.Fixed()) {
			mr := NoneReq(req.frac)
			fc1, vc1 := Optimize(u.source1, mode, mr)
			fc2, vc2 := Optimize(u.source2, mode, mr)
			if fc1+vc1 < impossible && fc2+vc2 < impossible {
				return fc1 + fc2, vc1 + vc2,
					&unionApproach{strat: unionMerge, req1: mr, req2: mr}
			}
		}
	}
	type b struct {
		order []string
		req   Require
	}
	best := newBest[b]()
	orders := u.mergeIndexes(req)
	for _, order := range orders {
		srcReq := OrderReq(order, req.frac)
		fc1, vc1 := Optimize(u.source1, mode, srcReq)
		if fc1+vc1 < impossible {
			fc2, vc2 := Optimize(u.source2, mode, srcReq)
			best.update(fc1+fc2, vc1+vc2, b{order: order, req: srcReq})
		}
	}
	if best.none() {
		return impossible, impossible, nil
	}
	return best.fixcost, best.varcost,
		&unionApproach{strat: unionMerge, keyIndex: best.data.order,
			req1: best.data.req, req2: best.data.req}
}

// mergeIndexes finds orders that unionMerge can use to read both sources.
// For each pair of source indexes, their common prefix is a candidate order
// (both sources can physically produce it) if each source has a key covered
// by that prefix (taking fixed into account) and the order satisfies the req.
// It also tries bare common keys, which may require temp indexes.
//
// For ReqUnique (Lookup) the order must also be a key of the union RESULT,
// not just a key of each source. Otherwise the same order values can produce
// distinct rows from the two sources (e.g. when one source has an extra
// extend column), and Lookup via Select+Get would return more than one row.
// When no order qualifies, optMerge returns impossible and the optimizer
// falls back to wrapping the union in a TempIndex.
func (u *Union) mergeIndexes(req Require) [][]string {
	keys := u.Keys()
	fixed := u.Fixed()
	fixed1 := u.source1.Fixed()
	keys1 := u.source1.Keys()
	indexes1 := u.source1.Indexes()
	fixed2 := u.source2.Fixed()
	keys2 := u.source2.Keys()
	indexes2 := u.source2.Indexes()
	needResultKey := req.use == ReqUnique
	var results [][]string
	for _, idx1 := range indexes1 {
		for _, idx2 := range indexes2 {
			order := slc.CommonPrefix(idx1, idx2)
			if len(order) > 0 &&
				hasKey(order, keys1, fixed1) &&
				hasKey(order, keys2, fixed2) &&
				req.SatisfiedByWithFixed(order, fixed1) &&
				req.SatisfiedByWithFixed(order, fixed2) &&
				(!needResultKey || hasKey(order, keys, fixed)) &&
				!slc.ContainsFn(results, order, slices.Equal) {
				results = append(results, order)
			}
		}
	}
	// try bare common keys, probably requiring temp indexes
	commonKeys := set.IntersectFn(keys1, keys2, set.Equal)
	for _, key := range commonKeys {
		if req.SatisfiedByWithFixed(key, fixed1) &&
			req.SatisfiedByWithFixed(key, fixed2) &&
			(!needResultKey || hasKey(key, keys, fixed)) &&
			!slc.ContainsFn(results, key, slices.Equal) {
			results = append(results, key)
		}
	}
	return results
}

func (u *Union) setApproach(req Require, approach any, tran QueryTran) {
	app := approach.(*unionApproach)
	u.strat = app.strat
	if app.strat == 0 {
		u.strat = unionLookup
	}
	if u.strat == unionMerge {
		unionMergeCount.Add(1)
		if u.disjoint != "" {
			unionMergeDisjoint.Add(1)
		}
	} else {
		unionLookupCount.Add(1)
	}
	if u.disjoint != "" {
		unionDisjointCount.Add(1)
	}
	u.keyIndex = app.keyIndex
	if app.reverse {
		u.source1, u.source2 = u.source2, u.source1
	}
	u.source1 = SetApproach(u.source1, app.req1, tran)
	u.source2 = SetApproach(u.source2, app.req2, tran)
	u.header = JoinHeaders(u.source1.Header(), u.source2.Header())
	u.src1Only = set.Difference(u.source1.Columns(), u.source2.Columns())
	u.empty1 = make(Row, len(u.source1.Header().Fields))
	u.empty2 = make(Row, len(u.source2.Header().Fields))
	u.state = rewound
	u.src1get = u.source1.Get
	u.src2get = u.source2.Get
}

// keyPrefixOfIndex returns the prefix of index up to and including
// the last field that belongs to key.
// This is the minimum index prefix that both sources must share
// for the merge to iterate in a compatible order.
func keyPrefixOfIndex(index, key []string) []string {
	for i := len(index) - 1; i >= 0; i-- {
		if slices.Contains(key, index[i]) {
			return index[:i+1]
		}
	}
	return nil
}

// indexContainsKey returns a key from keys if the index contains all fields
// of that key, otherwise nil.
// See also [hasKey]
func indexContainsKey(index []string, keys [][]string) []string {
	for _, key := range keys {
		if set.Subset(index, key) {
			return key
		}
	}
	return nil
}

// keyFieldOrder returns the order of the key fields as they appear in the index.
func keyFieldOrder(index, key []string) []string {
	result := make([]string, 0, len(key))
	for _, col := range index {
		if slices.Contains(key, col) {
			result = append(result, col)
		}
	}
	return result
}

// sameKeyFieldOrder returns true if the key fields appear in the same order
// in the index as in keyOrder.
func sameKeyFieldOrder(index, key []string, keyOrder []string) bool {
	order := keyFieldOrder(index, key)
	return slices.Equal(order, keyOrder)
}

func mergeIndexes(keys, indexes1, indexes2 [][]string,
	callback func(key []string, i1, i2 int)) {
	for _, key := range keys {
		callback(key, -1, -1) // -1 means key
		for i1, idx1 := range indexes1 {
			if keyperm := keyPerm(idx1, key); keyperm != nil {
				for i2, idx2 := range indexes2 {
					if slc.HasPrefix(idx2, keyperm) {
						callback(key, i1, i2)
					}
				}
			}
		}
	}
}

func keyPerm(index, key []string) []string {
	if len(index) >= len(key) {
		index = index[:len(key)]
		if set.Equal(index, key) {
			return index
		}
	}
	return nil
}

// execution --------------------------------------------------------

func (u *Union) Rewind() {
	u.source1.Rewind()
	u.source2.Rewind()
	u.state = rewound
}

func (u *Union) Get(th *Thread, dir Dir) Row {
	defer func(t uint64) { u.tget += tsc.Read() - t }(tsc.Read())
	if u.state == eof {
		return nil
	}
	var row Row
	switch u.strat {
	case unionLookup:
		row = u.getLookup(th, dir)
	case unionMerge:
		if u.disjoint != "" {
			row = u.getMergeDisjoint(th, dir)
		} else {
			row = u.getMerge(th, dir)
		}
	default:
		panic(assert.ShouldNotReachHere())
	}
	if row != nil {
		u.state = within
		u.ngets++
	} else {
		u.state = eof
	}
	return row
}

func (u *Union) getLookup(th *Thread, dir Dir) Row {
	if u.state == rewound {
		u.src1 = (dir == Next)
	}
	var row Row
	for {
		if u.src1 {
			for {
				row = u.src1get(th, dir)
				if row == nil {
					break
				}
				if !u.source2Has(th, row) {
					return JoinRows(row, u.empty2)
				}
			}
			if dir == Prev {
				return nil
			}
			u.source2.Rewind() // source2 may be stuck at eof from a prior Prev
			u.src1 = false
		} else { // source2
			row = u.src2get(th, dir)
			if row != nil {
				return JoinRows(u.empty1, row)
			}
			if dir == Next {
				return nil
			}
			u.source1.Rewind() // source1 may be stuck at eof from a prior Next
			u.src1 = true
		}
		// continue
	}
}

func (u *Union) getMerge(th *Thread, dir Dir) (r Row) {
	if u.mergeCols == nil {
		// compare keyIndex fields first
		u.mergeCols = set.Union(u.keyIndex, u.allCols)
	}

	// refill row1 and row2
	if u.state == rewound || (u.src1 && u.src2) {
		u.get1(th, dir)
		u.get2(th, dir)
	} else if u.src1 {
		u.get1(th, dir)
		if dir != u.prevDir {
			u.get2(th, dir)
		}
	} else if u.src2 {
		u.get2(th, dir)
		if dir != u.prevDir {
			u.get1(th, dir)
		}
	}

	u.prevDir = dir
	u.src1, u.src2 = false, false
	if u.row1 == nil && u.row2 == nil {
		u.src1, u.src2 = true, true
		return nil
	} else if u.row2 == nil {
		u.src1 = true
		return JoinRows(u.row1, u.empty2)
	} else if u.row1 == nil {
		u.src2 = true
		return JoinRows(u.empty1, u.row2)
	}
	cmp := u.compare(th, u.row1, u.row2, u.mergeCols)
	if cmp == 0 {
		// rows identical, arbitrarily return row1
		u.src1, u.src2 = true, true
		return JoinRows(u.row1, u.empty2)
	}
	if dir == Prev {
		cmp *= -1
	}
	if cmp < 0 {
		u.src1 = true
		return JoinRows(u.row1, u.empty2)
	} else {
		u.src2 = true
		return JoinRows(u.empty1, u.row2)
	}
}

func (u *Union) get1(th *Thread, dir Dir) {
	if dir != u.prevDir && u.row1 == nil {
		u.source1.Rewind()
	}
	u.row1 = u.src1get(th, dir)
}

func (u *Union) get2(th *Thread, dir Dir) {
	if dir != u.prevDir && u.row2 == nil {
		u.source2.Rewind()
	}
	u.row2 = u.src2get(th, dir)
}

func (u *Union) compare(th *Thread, row1, row2 Row, cols []string) int {
	for _, col := range cols {
		x1 := row1.GetRawVal(u.source1.Header(), col, th, u.st)
		x2 := row2.GetRawVal(u.source2.Header(), col, th, u.st)
		if c := strings.Compare(x1, x2); c != 0 {
			return c
		}
	}
	return 0
}

func (u *Union) getMergeDisjoint(th *Thread, dir Dir) (r Row) {
	// refill row1 and row2
	if u.state == rewound {
		u.get1(th, dir)
		u.get2(th, dir)
	} else if u.src1 {
		u.get1(th, dir)
		if dir != u.prevDir {
			u.get2(th, dir)
		}
	} else if u.src2 {
		u.get2(th, dir)
		if dir != u.prevDir {
			u.get1(th, dir)
		}
	}

	u.prevDir = dir
	u.src1, u.src2 = false, false
	if u.row1 == nil && u.row2 == nil {
		u.src1, u.src2 = true, true
		return nil
	} else if u.row2 == nil {
		u.src1 = true
		return JoinRows(u.row1, u.empty2)
	} else if u.row1 == nil {
		u.src2 = true
		return JoinRows(u.empty1, u.row2)
	}

	cmp := u.compare(th, u.row1, u.row2, u.keyIndex)
	if dir == Next {
		if cmp <= 0 {
			u.src1 = true
			return JoinRows(u.row1, u.empty2)
		} else {
			u.src2 = true
			return JoinRows(u.empty1, u.row2)
		}
	} else { // Prev
		if cmp <= 0 {
			u.src2 = true
			return JoinRows(u.empty1, u.row2)
		} else {
			u.src1 = true
			return JoinRows(u.row1, u.empty2)
		}
	}
}

func nothing(*Thread, Dir) Row { return nil }

func (u *Union) Select(sels Sels) {
	u.nsels++
	u.state = rewound
	u.src1get = u.source1.Get
	u.src2get = u.source2.Get
	if sels == nil { // clear
		u.source1.Select(nil)
		u.source2.Select(nil)
		return
	}
	if selConflict(u.source1.Columns(), sels) {
		u.src1get = nothing
	} else {
		u.source1.Select(removeNonexistentEmpty(u.source1.Columns(), sels))
	}
	if selConflict(u.source2.Columns(), sels) {
		u.src2get = nothing
	} else {
		u.source2.Select(removeNonexistentEmpty(u.source2.Columns(), sels))
	}
}

func removeNonexistentEmpty(srccols []string, sels Sels) Sels {
	for i, sel := range sels {
		if !slices.Contains(srccols, sel.col) && sel.val == "" {
			newsels := slices.Clip(sels[:i])
			for ; i < len(sels); i++ {
				if slices.Contains(srccols, sels[i].col) || sels[i].val != "" {
					newsels = append(newsels, sels[i])
				}
			}
			if len(newsels) == 0 {
				return nil
			}
			return newsels
		}
	}
	return sels
}

// selConflict is also used by Table
func selConflict(srcCols []string, sels Sels) bool {
	for _, sel := range sels {
		if sel.val != "" && !slices.Contains(srcCols, sel.col) {
			return true
		}
	}
	return false
}

func (u *Union) Lookup(th *Thread, sels Sels) Row {
	u.nlooks++
	return lookupViaSelectGet(u, th, sels)
}

func (u *Union) Simple(th *Thread) []Row {
	// rows1 + rows2 not in rows1
	empty1 := make(Row, len(u.source1.Header().Fields))
	empty2 := make(Row, len(u.source2.Header().Fields))
	rows1 := u.source1.Simple(th)
	rows2 := u.source2.Simple(th)
	rows := rows1
outer:
	for _, row2 := range rows2 {
		for _, row1 := range rows1 {
			if u.equal(th, row1, row2) {
				continue outer
			}
		}
		rows = append(rows, JoinRows(empty1, row2))
	}
	for i := range rows[:len(rows1)] {
		rows[i] = JoinRows(rows[i], empty2)
	}
	return rows
}
