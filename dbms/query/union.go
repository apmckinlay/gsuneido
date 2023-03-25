// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/ord"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/exp/slices"
)

type Union struct {
	src1get   func(*Thread, Dir) Row
	src2get   func(*Thread, Dir) Row
	row1      Row
	mergeCols []string
	empty1    Row
	empty2    Row
	row2      Row
	Compatible
	strategy unionStrategy
	src2     bool
	prevDir  Dir
	src1     bool
	rewound  bool
}

type unionApproach struct {
	keyIndex   []string
	idx1, idx2 []string
	strategy   unionStrategy
	reverse    bool
}

type unionStrategy int

const (
	// unionMerge is a merge of source and source2
	unionMerge unionStrategy = iota + 2
	// unionLookup is source not in source2, followed by source2 (unordered)
	unionLookup
)

func NewUnion(src1, src2 Query) *Union {
	u := &Union{Compatible: Compatible{
		Query2: Query2{source1: src1, source2: src2}}}
	u.init(u.calcFixed)
	return u
}

func (u *Union) String() string {
	return u.String2(u.stringOp())
}

func (u *Union) stringOp() string {
	strategy := ""
	switch u.strategy {
	case unionMerge:
		strategy += "-MERGE"
	case unionLookup:
		if u.disjoint == "" {
			strategy += "-LOOKUP"
		}
	}
	if u.keyIndex != nil {
		strategy += str.Join("(,)", u.keyIndex)
	}
	return u.Compatible.stringOp("UNION", strategy)
}

func (u *Union) Columns() []string {
	return u.allCols
}

func (u *Union) Keys() [][]string {
	if u.disjoint == "" {
		return [][]string{u.allCols}
	}
	keys := u.keypairs()
	for i := range keys {
		// keypairs must ensure that appending is valid
		keys[i] = set.AddUnique(keys[i], u.disjoint)
	}
	return withoutDupsOrSupersets(keys)
}

func (*Union) fastSingle() bool {
	return false
}

func (u *Union) Indexes() [][]string {
	// lookup can read via any index
	return set.UnionFn(u.source1.Indexes(), u.source2.Indexes(),
		slices.Equal[string])
}

func (u *Union) Nrows() (int, int) {
	n1, p1 := u.source1.Nrows()
	n2, p2 := u.source2.Nrows()
	return u.nrowsCalc(n1, n2), u.nrowsCalc(p1, p2)
}

func (u *Union) nrowsCalc(n1, n2 int) int {
	if u.disjoint != "" {
		return n1 + n2
	}
	min := ord.Max(n1, n2) // smaller could be all duplicates
	max := n1 + n2         // could be no duplicates
	return (min + max) / 2 // estimate half way between
}

func (u *Union) rowSize() int {
	return (u.source1.rowSize() + u.source2.rowSize()) / 2
}

func (u *Union) Transform() Query {
	src1 := u.source1.Transform()
	src2 := u.source2.Transform()
	if _, ok := src1.(*Nothing); ok {
		// remove unnecessary Union
		return keepCols(src2, src1, u.Header())
	}
	if _, ok := src2.(*Nothing); ok {
		// remove unnecessary Union
		return keepCols(src1, src2, u.Header())
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

func (u *Union) calcFixed(fixed1, fixed2 []Fixed) []Fixed {
	fixed := make([]Fixed, 0, len(fixed1)+len(fixed2))
	// add ones that are in both
	for _, f1 := range fixed1 {
		for _, f2 := range fixed2 {
			if f1.col == f2.col {
				fixed = append(fixed,
					Fixed{col: f1.col, values: set.Union(f1.values, f2.values)})
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
				Fixed{col: f1.col, values: set.Union(f1.values, emptyStr)})
		}
	}
	cols1 := u.source1.Columns()
	for _, f2 := range fixed2 {
		if !slices.Contains(cols1, f2.col) {
			fixed = append(fixed,
				Fixed{col: f2.col, values: set.Union(f2.values, emptyStr)})
		}
	}
	return fixed
}

func (u *Union) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	// if there is a required index, use Merge
	if index != nil {
		// if not disjoint then index must also be a key
		if u.disjoint == "" &&
			(!handlesIndex(u.source1.Keys(), index) ||
				!handlesIndex(u.source2.Keys(), index)) {
			return impossible, impossible, nil
		}
		fixcost1, varcost1 := Optimize(u.source1, mode, index, frac)
		fixcost2, varcost2 := Optimize(u.source2, mode, index, frac)
		approach := &unionApproach{keyIndex: index, strategy: unionMerge,
			idx1: index, idx2: index}
		return fixcost1 + fixcost2, varcost1 + varcost2, approach
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

func handlesIndex(keys [][]string, index []string) bool {
	if len(keys) == 1 && len(keys[0]) == 0 {
		return true // singleton
	}
	return slc.ContainsFn(keys, index, set.Equal[string])
}

func (*Union) optMerge(src1, src2 Query, mode Mode, frac float64) (Cost, Cost, any) {
	// if we get here, there is no required index, and it's not disjoint
	// we need a common key (unique) index to eliminate duplicates
	fixed1 := src1.Fixed()
	indexes1 := src1.Indexes()
	idxs1 := withoutFixed2(indexes1, fixed1)
	fixed2 := src2.Fixed()
	indexes2 := src2.Indexes()
	idxs2 := withoutFixed2(indexes2, fixed2)

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
			bestKey = key
			bestFixCost = fixcost1 + fixcost2
			bestVarCost = varcost1 + varcost2
			bestIdx1, bestIdx2 = index1, index2
		}
	}
	keys1 := withoutFixed2(src1.Keys(), fixed1)
	keys2 := withoutFixed2(src2.Keys(), fixed2)
	// interesect using set.Equal to ignore order
	keys := set.IntersectFn(keys1, keys2, set.Equal[string])
	mergeIndexes(keys, idxs1, idxs2, opt)
	approach := &unionApproach{keyIndex: bestKey, strategy: unionMerge,
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
	best := newBestIndex()
	fixcost1, varcost1 := Optimize(src1, mode, nil, frac)
	nrows1, _ := src1.Nrows()
	for _, key := range src2.Keys() {
		fixcost2, varcost2 :=
			LookupCost(src2, mode, key, int(float64(nrows1)*frac))
		best.update(key, fixcost2, varcost2)
	}
	approach := &unionApproach{keyIndex: best.index, strategy: unionLookup,
		idx1: nil, idx2: best.index}
	if src1 == u.source2 {
		approach.reverse = true
		best.fixcost += outOfOrder
	}
	return fixcost1 + best.fixcost, varcost1 + best.varcost, approach
}

func (u *Union) setApproach(_ []string, frac float64, approach any, tran QueryTran) {
	app := approach.(*unionApproach)
	u.strategy = app.strategy
	if app.strategy == 0 {
		u.strategy = unionLookup
	}
	u.keyIndex = app.keyIndex
	if app.reverse {
		u.source1, u.source2 = u.source2, u.source1
	}
	u.source1 = SetApproach(u.source1, app.idx1, frac, tran)
	if app.strategy == unionLookup {
		frac = 0
	}
	u.source2 = SetApproach(u.source2, app.idx2, frac, tran)

	u.empty1 = make(Row, len(u.source1.Header().Fields))
	u.empty2 = make(Row, len(u.source2.Header().Fields))

	u.rewound = true
	u.src1get = u.source1.Get
	u.src2get = u.source2.Get
}

// execution --------------------------------------------------------

func (u *Union) Rewind() {
	u.source1.Rewind()
	u.source2.Rewind()
	u.rewound = true
}

func (u *Union) Get(th *Thread, dir Dir) Row {
	defer func() { u.rewound = false }()
	switch u.strategy {
	case unionLookup:
		return u.getLookup(th, dir)
	case unionMerge:
		return u.getMerge(th, dir)
	}
	panic(assert.ShouldNotReachHere())
}

func (u *Union) getLookup(th *Thread, dir Dir) Row {
	if u.rewound {
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
			u.src1 = false
			u.source2.Rewind()
		} else { // source2
			row = u.src2get(th, dir)
			if row != nil {
				return JoinRows(u.empty1, row)
			}
			if dir == Next {
				return nil
			}
			u.src1 = true
			// continue
		}
	}
}

func (u *Union) getMerge(th *Thread, dir Dir) (r Row) {
	if u.hdr1 == nil {
		u.hdr1 = u.source1.Header()
		u.hdr2 = u.source2.Header()
	}
	if u.mergeCols == nil {
		// compare keyIndex fields first
		u.mergeCols = set.Union(u.keyIndex, u.allCols)
	}
	get1 := func() {
		if dir != u.prevDir && u.row1 == nil {
			u.source1.Rewind()
		}
		u.row1 = u.src1get(th, dir)
	}
	get2 := func() {
		if dir != u.prevDir && u.row2 == nil {
			u.source2.Rewind()
		}
		u.row2 = u.src2get(th, dir)
	}

	// refill row1 and row2
	if u.rewound || (u.src1 && u.src2) {
		get1()
		get2()
	} else if u.src1 {
		get1()
		if dir != u.prevDir {
			get2()
		}
	} else if u.src2 {
		get2()
		if dir != u.prevDir {
			get1()
		}
	}
	// fmt.Println("row1:", u.row1)
	// fmt.Println("row2:", u.row2)

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
	cmp := u.compare(th, u.row1, u.row2, u.hdr1, u.hdr2)
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

func (u *Union) compare(th *Thread, row1, row2 Row, hdr1, hdr2 *Header) int {
	for _, col := range u.mergeCols {
		x1 := row1.GetRawVal(hdr1, col, th, u.st)
		x2 := row2.GetRawVal(hdr2, col, th, u.st)
		if c := strings.Compare(x1, x2); c != 0 {
			return c
		}
	}
	return 0
}

func nothing(*Thread, Dir) Row { return nil }

func (u *Union) Select(cols, vals []string) {
	// fmt.Println("Union Select", cols, unpack(vals))
	u.rewound = true
	u.src1get = u.source1.Get
	u.src2get = u.source2.Get
	if cols == nil { // clear
		u.source1.Select(nil, nil)
		u.source2.Select(nil, nil)
		return
	}
	if u.hdr1 == nil {
		u.hdr1 = u.source1.Header()
		u.hdr2 = u.source2.Header()
	}
	if selConflict(u.hdr1.Columns, cols, vals) {
		u.src1get = nothing
	} else {
		u.source1.Select(cols, vals)
	}
	if selConflict(u.hdr2.Columns, cols, vals) {
		u.src2get = nothing
	} else {
		u.source2.Select(cols, vals)
	}
}

func selConflict(srcCols, cols, vals []string) bool {
	for i, col := range cols {
		if vals[i] != "" && !slices.Contains(srcCols, col) {
			return true
		}
	}
	return false
}

func (u *Union) Lookup(th *Thread, cols, vals []string) Row {
	u.Select(cols, vals)
	defer u.Select(nil, nil) // clear select
	return u.Get(th, Next)
}
