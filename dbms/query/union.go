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
	u.lookCost.Set(src1.lookupCost() + src2.lookupCost())
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
		s += "-merge"
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
	switch req.Use() {
	case ReqUnordered:
		return u.optUnordered(mode, req)
	case ReqOrdered, ReqGrouped:
		return u.optMerge(mode, req)
	case ReqLookup:
		if u.disjoint != "" {
			return u.optLookup(mode, req)
		}
		// Non-disjoint: the lookup strategy's source2Has does inner
		// Lookup calls that clobber the parent's Select state, causing
		// duplicates. Use merge instead (optTempIndex2 wraps it for
		// efficient lookups), matching v1 optimize(index != nil).
		return u.optMerge(mode, req)
	}
	panic(assert.ShouldNotReachHere())
}

func (u *Union) optUnordered(mode Mode, req Require) (Cost, Cost, any) {
	if u.disjoint != "" {
		mr := UnorderedReq(req.frac)
		fc1, vc1 := Optimize(u.source1, mode, mr)
		fc2, vc2 := Optimize(u.source2, mode, mr)
		return fc1 + fc2, vc1 + vc2,
			&unionApproach{strat: unionLookup, req1: mr, req2: mr}
	}
	mergeFix, mergeVar, mergeApp := u.optMerge(mode, req)
	lookupFix, lookupVar, lookupApp := u.optLookup(mode, req)
	lookupRevFix, lookupRevVar, lookupRevApp := u.optLookupRev(mode, req)
	return min3(
		mergeFix, mergeVar, mergeApp,
		lookupFix, lookupVar, lookupApp,
		lookupRevFix, lookupRevVar, lookupRevApp)
}

func (u *Union) optMerge(mode Mode, req Require) (Cost, Cost, any) {
	if u.disjoint != "" {
		mr := OrderedReq(req.cols, req.frac)
		fc1, vc1 := Optimize(u.source1, mode, mr)
		fc2, vc2 := Optimize(u.source2, mode, mr)
		if fc1+vc1 >= impossible || fc2+vc2 >= impossible {
			return impossible, impossible, nil
		}
		return fc1 + fc2, vc1 + vc2,
			&unionApproach{keyIndex: req.cols, strat: unionMerge, req1: mr, req2: mr}
	}
	if req.Use() == ReqUnordered {
		return u.optMergeNoOrder(mode, req)
	}
	return u.optMergeWithOrder(mode, req)
}

func (u *Union) optMergeNoOrder(mode Mode, req Require) (Cost, Cost, any) {
	fixed1 := u.source1.Fixed()
	indexes1 := u.source1.Indexes()
	idxs1 := fixed1.RemoveFromAll(indexes1)
	fixed2 := u.source2.Fixed()
	indexes2 := u.source2.Indexes()
	idxs2 := fixed2.RemoveFromAll(indexes2)
	keys1 := fixed1.RemoveFromAll(u.source1.Keys())
	keys2 := fixed2.RemoveFromAll(u.source2.Keys())
	commonKeys := set.IntersectFn(keys1, keys2, set.Equal[string])

	bestFixCost := impossible
	bestVarCost := impossible
	var bestApproach *unionApproach
	for _, key := range commonKeys {
		// try key itself
		mr := OrderedReq(key, req.frac)
		fc1, vc1 := Optimize(u.source1, mode, mr)
		fc2, vc2 := Optimize(u.source2, mode, mr)
		if fc1+vc1 < impossible && fc2+vc2 < impossible &&
			fc1+vc1+fc2+vc2 < bestFixCost+bestVarCost {
			bestFixCost = fc1 + fc2
			bestVarCost = vc1 + vc2
			bestApproach = &unionApproach{keyIndex: key,
				strat: unionMerge, req1: mr, req2: mr}
		}
		// try index pairs
		for i1, idx1 := range idxs1 {
			if kp := keyPerm(idx1, key); kp != nil {
				mr1 := OrderedReq(indexes1[i1], req.frac)
				fc1i, vc1i := Optimize(u.source1, mode, mr1)
				for i2, idx2 := range idxs2 {
					if slc.HasPrefix(idx2, kp) {
						mr2 := OrderedReq(indexes2[i2], req.frac)
						fc2i, vc2i := Optimize(u.source2, mode, mr2)
						if fc1i+vc1i < impossible && fc2i+vc2i < impossible &&
							fc1i+vc1i+fc2i+vc2i < bestFixCost+bestVarCost {
							bestFixCost = fc1i + fc2i
							bestVarCost = vc1i + vc2i
							bestApproach = &unionApproach{
								keyIndex: indexes1[i1][:len(key)],
								strat:    unionMerge, req1: mr1, req2: mr2}
						}
					}
				}
			}
		}
	}
	if bestApproach == nil {
		return impossible, impossible, nil
	}
	return bestFixCost, bestVarCost, bestApproach
}

func (u *Union) optMergeWithOrder(mode Mode, req Require) (Cost, Cost, any) {
	order := req.cols
	fixed1 := u.source1.Fixed()
	indexes1 := u.source1.Indexes()
	keys1 := u.source1.Keys()
	indexes2 := u.source2.Indexes()
	keys2 := u.source2.Keys()
	emptyKey1 := isEmptyKey(keys1)
	emptyKey2 := isEmptyKey(keys2)

	if emptyKey1 && emptyKey2 {
		mr := OrderedReq(nil, req.frac)
		fc1, vc1 := Optimize(u.source1, mode, mr)
		fc2, vc2 := Optimize(u.source2, mode, mr)
		return fc1 + fc2, vc1 + vc2,
			&unionApproach{strat: unionMerge, req1: mr, req2: mr}
	}

	if emptyKey1 {
		return u.bestMergeIndex(mode, req,
			u.source2, u.source1, indexes2, keys2, order)
	}
	if emptyKey2 {
		return u.bestMergeIndex(mode, req,
			u.source1, u.source2, indexes1, keys1, order)
	}

	bestFixCost := impossible
	bestVarCost := impossible
	var bestApproach *unionApproach

	for _, index1 := range indexes1 {
		if !slc.HasPrefix(index1, order) {
			continue
		}
		if req.Use() == ReqGrouped {
			nColsUnfixedReq := countUnfixed(req.cols, fixed1)
			if !grouped(index1, req.cols, nColsUnfixedReq, fixed1) {
				continue
			}
		}
		key1 := indexContainsKey(index1, keys1)
		if key1 == nil {
			continue
		}
		keyIdx := keyPrefixOfIndex(index1, key1)
		for _, index2 := range indexes2 {
			if !slc.HasPrefix(index2, keyIdx) {
				continue
			}
			if indexContainsKey(index2, keys2) == nil {
				continue
			}
			mr1 := OrderedReq(index1, req.frac)
			mr2 := OrderedReq(index2, req.frac)
			fc1, vc1 := Optimize(u.source1, mode, mr1)
			fc2, vc2 := Optimize(u.source2, mode, mr2)
			if fc1+vc1 < impossible && fc2+vc2 < impossible &&
				fc1+vc1+fc2+vc2 < bestFixCost+bestVarCost {
				bestFixCost = fc1 + fc2
				bestVarCost = vc1 + vc2
				bestApproach = &unionApproach{keyIndex: keyIdx,
					strat: unionMerge, req1: mr1, req2: mr2}
			}
		}
	}
	if bestApproach == nil {
		return impossible, impossible, nil
	}
	return bestFixCost, bestVarCost, bestApproach
}

func (u *Union) bestMergeIndex(mode Mode, req Require,
	srcKey, srcEmpty Query, indexes, keys [][]string, order []string) (Cost, Cost, any) {
	bestFixCost := impossible
	bestVarCost := impossible
	var bestApproach *unionApproach
	mr0 := OrderedReq(nil, req.frac)
	fc0, vc0 := Optimize(srcEmpty, mode, mr0)
	for _, index2 := range indexes {
		if !slc.HasPrefix(index2, order) {
			continue
		}
		key2 := indexContainsKey(index2, keys)
		if key2 == nil {
			continue
		}
		mr2 := OrderedReq(index2, req.frac)
		fc2, vc2 := Optimize(srcKey, mode, mr2)
		if fc0+vc0+fc2+vc2 < bestFixCost+bestVarCost {
			bestFixCost = fc0 + fc2
			bestVarCost = vc0 + vc2
			var mr1, mr2s Require
			if srcKey == u.source1 {
				mr1 = mr2
				mr2s = mr0
			} else {
				mr1 = mr0
				mr2s = mr2
			}
			bestApproach = &unionApproach{keyIndex: index2,
				strat: unionMerge, req1: mr1, req2: mr2s}
		}
	}
	if bestApproach == nil {
		return impossible, impossible, nil
	}
	return bestFixCost, bestVarCost, bestApproach
}

func (u *Union) optLookup(mode Mode, req Require) (Cost, Cost, any) {
	return u.optLookupDir(mode, req, false)
}

func (u *Union) optLookupRev(mode Mode, req Require) (Cost, Cost, any) {
	fixcost, varcost, app := u.optLookupDir(mode, req, true)
	if ap, ok := app.(*unionApproach); ok {
		ap.reverse = true
		fixcost += outOfOrder
	}
	return fixcost, varcost, app
}

func (u *Union) optLookupDir(mode Mode, req Require, reverse bool) (Cost, Cost, any) {
	src1, src2 := u.source1, u.source2
	if reverse {
		src1, src2 = src2, src1
	}
	nrows1, _ := src1.Nrows()
	var req1 Require
	switch req.Use() {
	case ReqUnordered:
		req1 = UnorderedReq(req.frac)
	case ReqLookup:
		req1 = GroupedReq(req.cols, req.SelectFrac(nrows1), req.nlookups)
	case ReqGrouped:
		req1 = req
	default:
		return impossible, impossible, nil
	}
	fc1, vc1 := Optimize(src1, mode, req1)
	if fc1+vc1 >= impossible {
		return impossible, impossible, nil
	}
	nlookups := req.LookupCount(nrows1)
	if u.disjoint != "" {
		mr2 := req1
		fc2, vc2 := Optimize(src2, mode, mr2)
		if fc2+vc2 >= impossible {
			return impossible, impossible, nil
		}
		return fc1 + fc2, vc1 + vc2,
			&unionApproach{strat: unionLookup,
				req1: req1, req2: mr2, reverse: reverse}
	}
	req2 := LookupReq(src2.Columns(), nlookups)
	fc2, vc2 := Optimize(src2, mode, req2)
	if fc2+vc2 >= impossible {
		return impossible, impossible, nil
	}
	ki := req2.cols
	return fc1 + fc2, vc1 + vc2,
		&unionApproach{keyIndex: ki, strat: unionLookup,
			req1: req1, req2: req2, reverse: reverse}
}

func (u *Union) setApproach(_ Require, approach any, tran QueryTran) {
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
	// fmt.Println("Union Select", cols, unpack(vals))
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
	u.Select(sels)
	defer u.Select(nil) // clear select
	return GetNext1(u, th)
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
