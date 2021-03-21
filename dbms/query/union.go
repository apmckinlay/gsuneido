// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/ssset"
)

type Union struct {
	Compatible
	fixed    []Fixed // lazy, calculated by Fixed()
	strategy unionStrategy
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
	unionMerge unionStrategy = iota + 1
	// unionLookup is source that aren't in source2, followed by source2
	unionLookup
	// unionFollow is source followed by source2 (disjoint)
	unionFollow
)

func (u *Union) String() string {
	op := "UNION"
	switch u.strategy {
	case unionMerge:
		op += "-MERGE"
	case unionLookup:
		op += "-LOOKUP"
	case unionFollow:
		op += "-FOLLOW"
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
	// NOTE: there are more possible indexes
	return ssset.Intersect(
		ssset.Intersect(u.source.Keys(), u.source.Indexes()),
		ssset.Intersect(u.source2.Keys(), u.source2.Indexes()))
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
					Fixed{f1.col, vUnion(f1.values, f2.values)})
				break
			}
		}
	}
	cols2 := u.source2.Columns()
	emptyStr := []runtime.Value{runtime.EmptyStr}
	for _, f1 := range fixed1 {
		if !sset.Contains(cols2, f1.col) {
			u.fixed = append(u.fixed,
				Fixed{f1.col, vUnion(f1.values, emptyStr)})
		}
	}
	cols1 := u.source.Columns()
	for _, f2 := range fixed2 {
		if !sset.Contains(cols1, f2.col) {
			u.fixed = append(u.fixed,
				Fixed{f2.col, vUnion(f2.values, emptyStr)})
		}
	}
	return u.fixed
}

func (u *Union) optimize(mode Mode, index []string) (Cost, interface{}) {
	// if there is a required index, use Merge
	if index != nil {
		// if not disjoint then index must also be a key
		if u.disjoint == "" && (!ssset.Contains(u.source.Keys(), index) ||
			!ssset.Contains(u.source2.Keys(), index)) {
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
		approach := &unionApproach{strategy: unionFollow}
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
	keyidxs := ssset.Intersect(
		ssset.Intersect(source.Keys(), source.Indexes()),
		ssset.Intersect(source2.Keys(), source2.Indexes()))
	var mergeKey []string
	mergeCost := impossible
	for _, k := range keyidxs {
		cost := Optimize(source, mode, k) + Optimize(source2, mode, k)
		if cost < mergeCost {
			mergeKey = k
			mergeCost = cost
		}
	}
	approach := &unionApproach{keyIndex: mergeKey, strategy: unionMerge,
		idx1: mergeKey, idx2: mergeKey}
	return mergeCost, approach
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
}
