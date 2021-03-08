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
	fixed    []Fixed
	strategy unionStrategy
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

func (u *Union) optimize(mode Mode, index []string, act action) Cost {
	// if there is a required index, use Merge
	if index != nil {
		// if not disjoint then index must also be a key
		if u.disjoint == "" && (!ssset.Contains(u.source.Keys(), index) ||
			!ssset.Contains(u.source2.Keys(), index)) {
			return impossible
		}
		if act == freeze {
			u.keyIndex = index
			u.strategy = unionMerge
		}
		return Optimize(u.source, mode, index, act) +
			Optimize(u.source2, mode, index, act)
	}
	// else no required index
	if u.disjoint != "" {
		// if disjoint use Follow
		if act == freeze {
			u.strategy = unionFollow
		}
		return Optimize(u.source, mode, nil, act) +
			Optimize(u.source2, mode, nil, act)
	}
	// else not disjoint
	mergeCost, mergeKey := u.optMerge(u.source, u.source2, mode)
	lookupCost, lookupKey := u.optLookup(u.source, u.source2, mode)
	lookupCostRev, lookupKeyRev := u.optLookup(u.source2, u.source, mode)
	lookupCostRev += outOfOrder
	cost := ints.Min(mergeCost, ints.Min(lookupCost, lookupCostRev))
	if cost >= impossible {
		return impossible
	}
	if act == freeze {
		if cost == mergeCost {
			u.strategy = unionMerge
			u.keyIndex = mergeKey
			// optimize1 to bypass temp index
			optimize1(u.source, mode, mergeKey, freeze)
			optimize1(u.source2, mode, mergeKey, freeze)
		} else {
			u.strategy = unionLookup
			if cost == lookupCostRev {
				u.source, u.source2, lookupKey =
					u.source2, u.source, lookupKeyRev // swap
			}
			u.keyIndex = lookupKey
			// optimize1 to bypass temp index
			optimize1(u.source, mode, nil, freeze)
			optimize1(u.source2, mode, u.keyIndex, freeze)
		}
	}
	return cost
}

func (*Union) optMerge(source, source2 Query, mode Mode) (Cost, []string) {
	keyidxs := ssset.Intersect(
		ssset.Intersect(source.Keys(), source.Indexes()),
		ssset.Intersect(source2.Keys(), source2.Indexes()))
	var mergeKey []string
	mergeCost := impossible
	for _, k := range keyidxs {
		// optimize1 to bypass temp index
		cost := optimize1(source, mode, k, assess) +
			optimize1(source2, mode, k, assess)
		if cost < mergeCost {
			mergeKey = k
			mergeCost = cost
		}
	}
	return mergeCost, mergeKey
}

func (*Union) optLookup(source, source2 Query, mode Mode) (Cost, []string) {
	var bestKey []string
	bestCost := impossible
	for _, key := range source2.Keys() {
		cost := optimize1(source, mode, nil, assess) +
			source.nrows()*source2.lookupCost()
		if cost < bestCost {
			bestKey = key
			bestCost = cost
		}
	}
	return bestCost, bestKey
}
