// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/ssset"
)

type Union struct {
	Compatible
	fixed []Fixed
}

func (u *Union) String() string {
	return u.Query2.String("UNION")
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
