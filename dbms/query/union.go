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
	idx1, idx2 []string
	strat      unionStrategy
	frac2      float64
	reverse    bool
	req1, req2 Require
}

type unionStrategy int

const (
	// unionMerge is a merge of source1 and source2
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
	return u.Compatible.String(s)
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

func (u *Union) optimize2(mode Mode, req Require) (Cost, Cost, any) {
	switch req.Use() {
	case ReqUnordered:
		return u.opt2Unordered(mode, req)
	case ReqOrdered, ReqGrouped:
		return u.opt2Merge(mode, req)
	case ReqLookup:
		if u.disjoint != "" {
			return u.opt2Lookup(mode, req)
		}
		// Non-disjoint: the lookup strategy's source2Has does inner
		// Lookup calls that clobber the parent's Select state, causing
		// duplicates. Use merge instead (optTempIndex2 wraps it for
		// efficient lookups), matching v1 optimize(index != nil).
		return u.opt2Merge(mode, req)
	}
	panic(assert.ShouldNotReachHere())
}

func (u *Union) opt2Unordered(mode Mode, req Require) (Cost, Cost, any) {
	if u.disjoint != "" {
		mr := UnorderedReq(req.frac)
		fc1, vc1 := Optimize2(u.source1, mode, mr)
		fc2, vc2 := Optimize2(u.source2, mode, mr)
		return fc1 + fc2, vc1 + vc2,
			&unionApproach{strat: unionLookup, req1: mr, req2: mr}
	}
	mergeFix, mergeVar, mergeApp := u.opt2Merge(mode, req)
	lookupFix, lookupVar, lookupApp :=
		u.opt2Lookup(mode, req)
	lookupRevFix, lookupRevVar, lookupRevApp :=
		u.opt2LookupRev(mode, req)
	return min3(mergeFix, mergeVar, mergeApp,
		lookupFix, lookupVar, lookupApp,
		lookupRevFix, lookupRevVar, lookupRevApp)
}

func (u *Union) opt2Merge(mode Mode, req Require) (Cost, Cost, any) {
	if u.disjoint != "" {
		mr := OrderedReq(req.cols, req.frac)
		fc1, vc1 := Optimize2(u.source1, mode, mr)
		fc2, vc2 := Optimize2(u.source2, mode, mr)
		if fc1+vc1 >= impossible || fc2+vc2 >= impossible {
			return impossible, impossible, nil
		}
		return fc1 + fc2, vc1 + vc2,
			&unionApproach{keyIndex: req.cols, strat: unionMerge,
				idx1: req.cols, idx2: req.cols, req1: mr, req2: mr}
	}
	if req.Use() == ReqUnordered {
		return u.optMergeNoOrder2(mode, req)
	}
	return u.optMergeWithOrder2(mode, req)
}

func (u *Union) optMergeNoOrder2(mode Mode, req Require) (Cost, Cost, any) {
	fixed1 := u.source1.Fixed()
	indexes1 := u.source1.Indexes()
	idxs1 := fixed1.RemoveFrom2(indexes1)
	fixed2 := u.source2.Fixed()
	indexes2 := u.source2.Indexes()
	idxs2 := fixed2.RemoveFrom2(indexes2)
	keys1 := fixed1.RemoveFrom2(u.source1.Keys())
	keys2 := fixed2.RemoveFrom2(u.source2.Keys())
	commonKeys := set.IntersectFn(keys1, keys2, set.Equal[string])

	bestFixCost := impossible
	bestVarCost := impossible
	var bestKey, bestIdx1, bestIdx2 []string
	var bestReq1, bestReq2 Require
	for _, key := range commonKeys {
		// try key itself
		mr := OrderedReq(key, req.frac)
		fc1, vc1 := Optimize2(u.source1, mode, mr)
		fc2, vc2 := Optimize2(u.source2, mode, mr)
		if fc1+vc1 < impossible && fc2+vc2 < impossible &&
			fc1+vc1+fc2+vc2 < bestFixCost+bestVarCost {
			bestFixCost = fc1 + fc2
			bestVarCost = vc1 + vc2
			bestKey = key
			bestIdx1, bestIdx2 = key, key
			bestReq1, bestReq2 = mr, mr
		}
		// try index pairs
		for i1, idx1 := range idxs1 {
			if kp := keyPerm(idx1, key); kp != nil {
				mr1 := OrderedReq(indexes1[i1], req.frac)
				fc1i, vc1i := Optimize2(u.source1, mode, mr1)
				for i2, idx2 := range idxs2 {
					if slc.HasPrefix(idx2, kp) {
						mr2 := OrderedReq(indexes2[i2], req.frac)
						fc2i, vc2i := Optimize2(u.source2, mode, mr2)
						if fc1i+vc1i < impossible && fc2i+vc2i < impossible &&
							fc1i+vc1i+fc2i+vc2i < bestFixCost+bestVarCost {
							bestFixCost = fc1i + fc2i
							bestVarCost = vc1i + vc2i
							bestKey = indexes1[i1][:len(key)]
							bestIdx1 = indexes1[i1]
							bestIdx2 = indexes2[i2]
							bestReq1, bestReq2 = mr1, mr2
						}
					}
				}
			}
		}
	}
	if bestIdx1 == nil {
		return impossible, impossible, nil
	}
	return bestFixCost, bestVarCost,
		&unionApproach{keyIndex: bestKey, strat: unionMerge,
			idx1: bestIdx1, idx2: bestIdx2, req1: bestReq1, req2: bestReq2}
}

func (u *Union) optMergeWithOrder2(mode Mode, req Require) (Cost, Cost, any) {
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
		fc1, vc1 := Optimize2(u.source1, mode, mr)
		fc2, vc2 := Optimize2(u.source2, mode, mr)
		return fc1 + fc2, vc1 + vc2,
			&unionApproach{strat: unionMerge, req1: mr, req2: mr}
	}

	if emptyKey1 {
		return u.bestMergeIndexOne2(mode, req,
			u.source2, u.source1, indexes2, keys2, order)
	}
	if emptyKey2 {
		return u.bestMergeIndexOne2(mode, req,
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
			fc1, vc1 := Optimize2(u.source1, mode, mr1)
			fc2, vc2 := Optimize2(u.source2, mode, mr2)
			if fc1+vc1 < impossible && fc2+vc2 < impossible &&
				fc1+vc1+fc2+vc2 < bestFixCost+bestVarCost {
				bestFixCost = fc1 + fc2
				bestVarCost = vc1 + vc2
				bestApproach = &unionApproach{keyIndex: keyIdx,
					strat: unionMerge, idx1: index1, idx2: index2,
					req1: mr1, req2: mr2}
			}
		}
	}
	if bestApproach == nil {
		return impossible, impossible, nil
	}
	return bestFixCost, bestVarCost, bestApproach
}

func (u *Union) bestMergeIndexOne2(mode Mode, req Require,
	srcKey, srcEmpty Query, indexes, keys [][]string, order []string) (Cost, Cost, any) {
	bestFixCost := impossible
	bestVarCost := impossible
	var bestApproach *unionApproach
	mr0 := OrderedReq(nil, req.frac)
	fc0, vc0 := Optimize2(srcEmpty, mode, mr0)
	for _, index2 := range indexes {
		if !slc.HasPrefix(index2, order) {
			continue
		}
		key2 := indexContainsKey(index2, keys)
		if key2 == nil {
			continue
		}
		mr2 := OrderedReq(index2, req.frac)
		fc2, vc2 := Optimize2(srcKey, mode, mr2)
		if fc0+vc0+fc2+vc2 < bestFixCost+bestVarCost {
			bestFixCost = fc0 + fc2
			bestVarCost = vc0 + vc2
			var mr1, mr2s Require
			var idx1, idx2 []string
			if srcKey == u.source1 {
				idx1 = index2
				mr1 = mr2
				mr2s = mr0
			} else {
				idx2 = index2
				mr1 = mr0
				mr2s = mr2
			}
			bestApproach = &unionApproach{keyIndex: index2,
				strat: unionMerge, idx1: idx1, idx2: idx2,
				req1: mr1, req2: mr2s}
		}
	}
	if bestApproach == nil {
		return impossible, impossible, nil
	}
	return bestFixCost, bestVarCost, bestApproach
}

func (u *Union) opt2Lookup(mode Mode, req Require) (Cost, Cost, any) {
	return u.opt2LookupDir(mode, req, false)
}

func (u *Union) opt2LookupRev(mode Mode, req Require) (Cost, Cost, any) {
	fixcost, varcost, app := u.opt2LookupDir(mode, req, true)
	if ap, ok := app.(*unionApproach); ok {
		ap.reverse = true
		fixcost += outOfOrder
	}
	return fixcost, varcost, app
}

func (u *Union) opt2LookupDir(mode Mode, req Require, reverse bool) (Cost, Cost, any) {
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
	fc1, vc1 := Optimize2(src1, mode, req1)
	if fc1+vc1 >= impossible {
		return impossible, impossible, nil
	}
	nlookups := req.LookupCount(nrows1)
	if u.disjoint != "" {
		mr2 := req1
		fc2, vc2 := Optimize2(src2, mode, mr2)
		if fc2+vc2 >= impossible {
			return impossible, impossible, nil
		}
		return fc1 + fc2, vc1 + vc2,
			&unionApproach{strat: unionLookup,
				req1: req1, req2: mr2, reverse: reverse}
	}
	req2, fc2, vc2 := anyKeyLookup2(src2, mode, nlookups)
	if fc2+vc2 >= impossible {
		return impossible, impossible, nil
	}
	ki := req2.cols
	if ki == nil {
		ki = []string{}
	}
	return fc1 + fc2, vc1 + vc2,
		&unionApproach{keyIndex: ki, strat: unionLookup,
			req1: req1, req2: req2, reverse: reverse}
}

func (u *Union) setApproach2(_ Require, approach any, tran QueryTran) {
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
	u.source1 = SetApproach2(u.source1, app.req1, tran)
	u.source2 = SetApproach2(u.source2, app.req2, tran)
	u.header = JoinHeaders(u.source1.Header(), u.source2.Header())
	u.empty1 = make(Row, len(u.source1.Header().Fields))
	u.empty2 = make(Row, len(u.source2.Header().Fields))
	u.state = rewound
	u.src1get = u.source1.Get
	u.src2get = u.source2.Get
}

func (u *Union) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	// if there is a required index, use Merge
	if index != nil {
		if u.disjoint != "" {
			fixcost1, varcost1 := Optimize(u.source1, mode, index, frac)
			fixcost2, varcost2 := Optimize(u.source2, mode, index, frac)
			if fixcost1+varcost1 >= impossible || fixcost2+varcost2 >= impossible {
				return impossible, impossible, nil
			}
			approach := &unionApproach{keyIndex: index, strat: unionMerge,
				idx1: index, idx2: index}
			return fixcost1 + fixcost2, varcost1 + varcost2, approach
		}
		idx1, idx2, bestKey, fixcost, varcost := u.bestMergeIndexes(index, mode, frac)
		if fixcost >= impossible {
			return impossible, impossible, nil
		}
		keyIndex := bestKey
		if keyIndex == nil {
			// for singletons, use the required order as the keyIndex
			keyIndex = index
		}
		approach := &unionApproach{keyIndex: keyIndex, strat: unionMerge,
			idx1: idx1, idx2: idx2}
		return fixcost, varcost, approach
	}
	// else no required index
	if u.disjoint != "" {
		fixcost1, varcost1 := Optimize(u.source1, mode, nil, frac)
		fixcost2, varcost2 := Optimize(u.source2, mode, nil, frac)
		approach := &unionApproach{} // will use getLookup, but no lookups
		return fixcost1 + fixcost2, varcost1 + varcost2, approach
	}
	// else not disjoint
	mergeFixCost, mergeVarCost, mergeApp :=
		u.optMerge(u.source1, u.source2, mode, frac)
	lookupFixCost, lookupVarCost, lookupApp :=
		u.optLookup(u.source1, u.source2, mode, frac)
	lookupRevFixCost, lookupRevVarCost, lookupRevApp :=
		u.optLookup(u.source2, u.source1, mode, frac)
	fixcost, varcost, approach := min3(
		mergeFixCost, mergeVarCost, mergeApp,
		lookupFixCost, lookupVarCost, lookupApp,
		lookupRevFixCost, lookupRevVarCost, lookupRevApp)
	// trace.Println("UNION", mode, index, frac)
	// trace.Println("    src1 keys", u.source1.Keys(), "indexes", u.source1.Indexes(),
	// 	"fastSingle", u.source1.fastSingle())
	// trace.Println("    src2 keys", u.source2.Keys(), "indexes", u.source2.Indexes(),
	// 	"fastSingle", u.source2.fastSingle())
	// trace.Println("    merge", mergeFixCost, "+", mergeVarCost,
	// 	"=", mergeFixCost+mergeVarCost)
	// trace.Println("    lookup", lookupFixCost, "+", lookupVarCost,
	// 	"=", lookupFixCost+lookupVarCost)
	// trace.Println("    lookupRev", lookupRevFixCost, "+", lookupRevVarCost,
	// 	"=", lookupRevFixCost+lookupRevVarCost)
	if fixcost >= impossible {
		return impossible, impossible, nil
	}
	return fixcost, varcost, approach
}

// bestMergeIndexes finds the pair of indexes from source1 and source2 that have
// the lowest cost for a merge operation with the required order.
// For each source1 index that is prefixed by the required order and contains a
// source1 key, the keyIndex is the prefix of that index through all key fields.
// A source2 index is only a candidate if it is also prefixed by keyIndex
// (ensuring both sources iterate in exactly the same order for the merge key).
// Returns the pair with the lowest cost.
// Empty keys (singletons) are handled specially:
// - If both sources have empty keys, order is not needed
// - If only one source has an empty key, we need a key index on the other
func (u *Union) bestMergeIndexes(order []string, mode Mode, frac float64) (
	idx1, idx2 []string, bestKey []string, fixcost, varcost Cost) {
	fixcost = impossible
	varcost = impossible
	indexes1 := u.source1.Indexes()
	keys1 := u.source1.Keys()
	indexes2 := u.source2.Indexes()
	keys2 := u.source2.Keys()
	emptyKey1 := isEmptyKey(keys1)
	emptyKey2 := isEmptyKey(keys2)
	if emptyKey1 && emptyKey2 {
		// both sources have empty keys, don't need order for Optimize
		fc1, vc1 := Optimize(u.source1, mode, nil, frac)
		fc2, vc2 := Optimize(u.source2, mode, nil, frac)
		return nil, nil, nil, fc1 + fc2, vc1 + vc2
	}
	if emptyKey1 {
		// only source1 has empty key, need index on source2 that includes a key
		idx2, bestKey, fixcost, varcost =
			bestMergeIndex(u.source1, u.source2, indexes2, keys2, order, mode, frac)
		return
	}
	if emptyKey2 {
		// only source2 has empty key, need index on source1 that includes a key
		idx1, bestKey, fixcost, varcost =
			bestMergeIndex(u.source2, u.source1, indexes1, keys1, order, mode, frac)
		return
	}
	// neither source has empty key
	for _, index1 := range indexes1 {
		if !slc.HasPrefix(index1, order) {
			continue
		}
		key1 := indexContainsKey(index1, keys1)
		if key1 == nil {
			continue
		}
		// keyIndex is the prefix of index1 through all key fields.
		// Both sources must iterate in this order for the merge to be correct.
		keyIndex := keyPrefixOfIndex(index1, key1)
		for _, index2 := range indexes2 {
			if !slc.HasPrefix(index2, keyIndex) {
				continue
			}
			if indexContainsKey(index2, keys2) == nil {
				continue
			}
			// candidate pair found, get cost
			fc1, vc1 := Optimize(u.source1, mode, index1, frac)
			fc2, vc2 := Optimize(u.source2, mode, index2, frac)
			if fc1+vc1+fc2+vc2 < fixcost+varcost {
				idx1 = index1
				idx2 = index2
				bestKey = keyIndex
				fixcost = fc1 + fc2
				varcost = vc1 + vc2
			}
		}
	}
	return
}

// bestMergeIndex finds the best index on src2 that includes a key,
func bestMergeIndex(src1, src2 Query, indexes2, keys2 [][]string,
	order []string, mode Mode, frac float64) (
	idx []string, bestKey []string, fixcost, varcost Cost) {
	fixcost = impossible
	varcost = impossible
	fc1, vc1 := Optimize(src1, mode, nil, frac)
	for _, index2 := range indexes2 {
		if !slc.HasPrefix(index2, order) {
			continue
		}
		key2 := indexContainsKey(index2, keys2)
		if key2 == nil {
			continue
		}
		fc2, vc2 := Optimize(src2, mode, index2, frac)
		if fc1+vc1+fc2+vc2 < fixcost+varcost {
			idx = index2
			bestKey = index2
			fixcost = fc1 + fc2
			varcost = vc1 + vc2
		}
	}
	return
}

// keyPrefixOfIndex returns the prefix of index up to and including
// the last field that belongs to key.
// This is the minimum index prefix that both sources must share
// for the merge to iterate in a compatible order.
func keyPrefixOfIndex(index, key []string) []string {
	last := -1
	for i, col := range index {
		if slices.Contains(key, col) {
			last = i
		}
	}
	return index[:last+1]
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

func (*Union) optMerge(src1, src2 Query, mode Mode, frac float64) (Cost, Cost, any) {
	// if we get here, there is no required index, and it's not disjoint
	// we need a common key (unique) index to eliminate duplicates
	fixed1 := src1.Fixed()
	indexes1 := src1.Indexes()
	idxs1 := fixed1.RemoveFrom2(indexes1)
	fixed2 := src2.Fixed()
	indexes2 := src2.Indexes()
	idxs2 := fixed2.RemoveFrom2(indexes2)

	var bestKey, bestIdx1, bestIdx2 []string
	bestFixCost := impossible
	bestVarCost := impossible
	opt := func(key []string, i1, i2 int) {
		var index1, index2 []string
		if i1 == -1 {
			index1 = key
			index2 = key
		} else {
			index1 = indexes1[i1]
			index2 = indexes2[i2]
		}
		fixcost1, varcost1 := Optimize(src1, mode, index1, frac)
		fixcost2, varcost2 := Optimize(src2, mode, index2, frac)
		if fixcost1+varcost1+fixcost2+varcost2 < bestFixCost+bestVarCost {
			// use the actual index order, not the key order,
			// so the merge comparison matches the iteration order
			bestKey = index1[:len(key)]
			bestFixCost = fixcost1 + fixcost2
			bestVarCost = varcost1 + varcost2
			bestIdx1, bestIdx2 = index1, index2
		}
	}
	keys1 := fixed1.RemoveFrom2(src1.Keys())
	keys2 := fixed2.RemoveFrom2(src2.Keys())
	// intersect using set.Equal to ignore order
	keys := set.IntersectFn(keys1, keys2, set.Equal[string])
	mergeIndexes(keys, idxs1, idxs2, opt)
	approach := &unionApproach{keyIndex: bestKey, strat: unionMerge,
		idx1: bestIdx1, idx2: bestIdx2}
	return bestFixCost, bestVarCost, approach
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

func (u *Union) optLookup(src1, src2 Query, mode Mode, frac float64) (Cost, Cost, any) {
	fixcost1, varcost1 := Optimize(src1, mode, nil, frac)
	nrows1, _ := src1.Nrows()
	nrows2, _ := src2.Nrows()
	lookups := int(float64(nrows1) * frac)
	frac2 := float64(lookups) / float64(max(1, nrows2))
	best := bestLookupIndex(src2, mode, lookups, frac2, nil)
	approach := &unionApproach{keyIndex: best.index, strat: unionLookup,
		idx1: nil, idx2: best.index, frac2: frac2}
	if src1 == u.source2 {
		approach.reverse = true
		best.fixcost += outOfOrder
	}
	return fixcost1 + best.fixcost, varcost1 + best.varcost, approach
}

func (u *Union) setApproach(_ []string, frac float64, approach any, tran QueryTran) {
	if u.disjoint != "" {
		unionDisjointCount.Add(1)
	}
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
	u.keyIndex = app.keyIndex
	if app.reverse {
		u.source1, u.source2 = u.source2, u.source1
	}
	u.source1 = SetApproach(u.source1, app.idx1, frac, tran)
	if app.strat == unionLookup {
		frac = app.frac2
	}
	u.source2 = SetApproach(u.source2, app.idx2, frac, tran)
	u.header = JoinHeaders(u.source1.Header(), u.source2.Header())

	u.empty1 = make(Row, len(u.source1.Header().Fields))
	u.empty2 = make(Row, len(u.source2.Header().Fields))

	u.state = rewound
	u.src1get = u.source1.Get
	u.src2get = u.source2.Get
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
