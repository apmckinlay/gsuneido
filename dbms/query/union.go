// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/setord"
	"github.com/apmckinlay/gsuneido/util/setset"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/strs"
)

type Union struct {
	Compatible
	fixed    []Fixed // lazy, calculated by Fixed()
	strategy unionStrategy
	rewound  bool
	empty1   Row
	empty2   Row
	src1     bool
	src2     bool
	key1     string
	key2     string
	row1     Row
	row2     Row
	curKey   string
}

type unionApproach struct {
	keyIndex   []string
	strategy   unionStrategy
	idx1, idx2 []string
	reverse    bool
}

type unionStrategy int

const (
	// unionMerge is a merge of source and source2
	unionMerge unionStrategy = iota + 2
	// unionLookup is source not in source2, followed by source2 (unordered)
	unionLookup
)

func (u *Union) String() string {
	op := "UNION"
	switch u.strategy {
	case unionMerge:
		op += "-MERGE"
	case unionLookup:
		op += "-LOOKUP"
	}
	return u.String2(op)
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
		keys[i] = sset.AddUnique(keys[i], u.disjoint)
	}
	// exclude any keys that are super-sets of another key
	var keys2 [][]string
outer:
	for i := 0; i < len(keys); i++ {
		for j := 0; j < len(keys); j++ {
			if i != j && sset.Subset(keys[i], keys[j]) {
				continue outer
			}
		}
		keys2 = append(keys2, keys[i])
	}
	return keys2
}

func (u *Union) Indexes() [][]string {
	// lookup can read via any index
	return setord.Union(u.source.Indexes(), u.source2.Indexes())
}

func (u *Union) nrows() int {
	return u.nrowsCalc(u.source.nrows(), u.source2.nrows())
}

func (u *Union) nrowsCalc(n1, n2 int) int {
	if u.disjoint != "" {
		return n1 + n2
	}
	min := ints.Max(n1, n2) // smaller could be all duplicates
	max := n1 + n2          // could be no duplicates
	return (min + max) / 2  // estimate half way between
}

func (u *Union) Transform() Query {
	u.source = u.source.Transform()
	u.source2 = u.source2.Transform()
	return u
}

func (u *Union) Fixed() []Fixed {
	if u.fixed != nil { // once only
		return u.fixed
	}
	fixed1 := u.source.Fixed()
	fixed2 := u.source2.Fixed()
	for _, f1 := range fixed1 {
		for _, f2 := range fixed2 {
			if f1.col == f2.col {
				u.fixed = append(u.fixed,
					Fixed{f1.col, sset.Union(f1.values, f2.values)})
				break
			}
		}
	}
	cols2 := u.source2.Columns()
	emptyStr := []string{""}
	for _, f1 := range fixed1 {
		if !sset.Contains(cols2, f1.col) {
			u.fixed = append(u.fixed,
				Fixed{f1.col, sset.Union(f1.values, emptyStr)})
		}
	}
	cols1 := u.source.Columns()
	for _, f2 := range fixed2 {
		if !sset.Contains(cols1, f2.col) {
			u.fixed = append(u.fixed,
				Fixed{f2.col, sset.Union(f2.values, emptyStr)})
		}
	}
	return u.fixed
}

func (u *Union) optimize(mode Mode, index []string) (Cost, interface{}) {
	// if there is a required index, use Merge
	if index != nil {
		// if not disjoint then index must also be a key
		if u.disjoint == "" && (!setset.Contains(u.source.Keys(), index) ||
			!setset.Contains(u.source2.Keys(), index)) {
			return impossible, nil
		}
		cost := Optimize(u.source, mode, index) + Optimize(u.source2, mode, index)
		approach := &unionApproach{keyIndex: index, strategy: unionMerge,
			idx1: index, idx2: index}
		return cost, approach
	}
	// else no required index
	if u.disjoint != "" {
		cost := Optimize(u.source, mode, nil) + Optimize(u.source2, mode, nil)
		approach := &unionApproach{strategy: unionLookup}
		return cost, approach
	}
	// else not disjoint
	mergeCost, mergeApp := u.optMerge(u.source, u.source2, mode)
	lookupCost, lookupApp := u.optLookup(u.source, u.source2, mode)
	lookupCostRev, lookupAppRev := u.optLookup(u.source2, u.source, mode)
	cost, approach := min3(mergeCost, mergeApp,
		lookupCost, lookupApp, lookupCostRev, lookupAppRev)
	if cost >= impossible {
		return impossible, nil
	}
	return cost, approach
}

func (*Union) optMerge(source, source2 Query, mode Mode) (Cost, interface{}) {
	// need key (unique) index to eliminate duplicates
	keys := setset.Intersect(source.Keys(), source2.Keys())
	var bestKey, bestIdx1, bestIdx2 []string
	bestCost := impossible
	for _, key := range keys {
		for _, idx1 := range source.Indexes() {
			if !sset.Subset(idx1, key) {
				continue
			}
			ik1 := sset.Intersect(idx1, key)
			for _, idx2 := range source2.Indexes() {
				ik2 := sset.Intersect(idx2, key)
				if strs.Equal(ik1, ik2) {
					cost := Optimize(source, mode, idx1) +
						Optimize(source2, mode, idx2)
					if cost < bestCost {
						bestKey = key
						bestCost = cost
						bestIdx1, bestIdx2 = idx1, idx2
					}
				}
			}
		}
	}
	approach := &unionApproach{keyIndex: bestKey, strategy: unionMerge,
		idx1: bestIdx1, idx2: bestIdx2}
	return bestCost, approach
}

func (u *Union) optLookup(source, source2 Query, mode Mode) (Cost, interface{}) {
	var bestKey []string
	bestCost := impossible
	for _, key := range source2.Keys() {
		cost := Optimize(source, mode, nil) +
			LookupCost(source2, mode, key, source.nrows())
		if cost < bestCost && Optimize(source2, mode, key) < impossible {
			bestKey = key
			bestCost = cost
		}
	}
	approach := &unionApproach{keyIndex: bestKey, strategy: unionLookup,
		idx1: nil, idx2: bestKey}
	if source == u.source2 {
		approach.reverse = true
		bestCost += outOfOrder
	}
	return bestCost, approach
}

func (u *Union) setApproach(_ []string, approach interface{}, tran QueryTran) {
	app := approach.(*unionApproach)
	u.strategy = app.strategy
	u.keyIndex = app.keyIndex
	if app.reverse {
		u.source, u.source2 = u.source2, u.source
	}
	u.source = SetApproach(u.source, app.idx1, tran)
	u.source2 = SetApproach(u.source2, app.idx2, tran)

	u.empty1 = make(Row, len(u.source.Header().Fields))
	u.empty2 = make(Row, len(u.source2.Header().Fields))

	u.rewound = true
}

// execution --------------------------------------------------------

func (u *Union) Rewind() {
	u.source.Rewind()
	u.source2.Rewind()
	u.rewound = true
}

func (u *Union) Get(dir Dir) Row {
	defer func() { u.rewound = false }()
	switch u.strategy {
	case unionLookup:
		return u.getLookup(dir)
	case unionMerge:
		return u.getMerge(dir)
	}
	panic("shouldn't reach here")
}

func (u *Union) getLookup(dir Dir) Row {
	if u.rewound {
		u.src1 = (dir == Next)
	}
	var row Row
	for {
		if u.src1 {
			for {
				row = u.source.Get(dir)
				if row == nil {
					break
				}
				if !u.source2Has(row) {
					return JoinRows(row, u.empty2)
				}
			}
			if dir == Prev {
				return nil
			}
			u.src1 = false
		} else { // source2
			row = u.source2.Get(dir)
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

func (u *Union) getMerge(dir Dir) Row {
	if u.hdr1 == nil {
		u.hdr1 = u.source.Header()
		u.hdr2 = u.source2.Header()
	}

	// read from the appropriate source(s)
	if u.rewound {
		u.fetch1(dir)
		u.fetch2(dir)
	} else {
		// curkey is required for changing direction
		if u.src1 || u.before(dir, u.key1, u.curKey, true) {
			u.fetch1(dir)
		}
		if u.src2 || u.before(dir, u.key2, u.curKey, false) {
			u.fetch2(dir)
		}
	}

	u.src1, u.src2 = false, false
	if u.row1 == nil && u.row2 == nil {
		u.curKey = u.key1
		u.src1 = true
		return nil
	} else if u.row1 != nil && u.row2 != nil && u.equal(u.row1, u.row2) {
		// rows same so return either one
		u.curKey = u.key1
		u.src1, u.src2 = true, true
		return JoinRows(u.row1, u.empty2)
	} else if u.row1 != nil &&
		(u.row2 == nil || u.before(dir, u.key1, u.key2, true)) {
		u.curKey = u.key1
		u.src1 = true
		return JoinRows(u.row1, u.empty2)
	} else {
		u.curKey = u.key2
		u.src2 = true
		return JoinRows(u.empty1, u.row2)
	}
}

func (u *Union) fetch1(dir Dir) {
	u.row1 = u.source.Get(dir)
	if u.row1 == nil {
		u.key1 = endKey(dir)
	} else {
		u.key1 = projectKey(u.row1, u.hdr1, u.keyIndex)
	}
}

func (u *Union) fetch2(dir Dir) {
	u.row2 = u.source2.Get(dir)
	if u.row2 == nil {
		u.key2 = endKey(dir)
	} else {
		u.key2 = projectKey(u.row2, u.hdr2, u.keyIndex)
	}
}

func (*Union) before(dir Dir, key1, key2 string, x bool) bool {
	if key1 == key2 {
		if dir == Next {
			return x
		}
		return !x
	}
	if dir == Next {
		return key1 < key2
	}
	return key1 > key2
}

func endKey(dir Dir) string {
	if dir == Next {
		return ixkey.Max
	} // else Prev
	return ixkey.Min
}

func (u *Union) Select(cols, vals []string) {
	u.source.Select(cols, vals)
	u.source2.Select(cols, vals)
	u.rewound = true
}
