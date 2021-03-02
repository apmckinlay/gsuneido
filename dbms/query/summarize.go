// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/ssset"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Summarize struct {
	Query1
	by       []string
	cols     []string
	ops      []string
	ons      []string
	wholeRow bool
}

func (su *Summarize) Init() {
	su.Query1.Init()
	if !sset.Subset(su.source.Columns(), su.by) {
		panic("summarize: nonexistent columns: " +
			str.Join(", ", sset.Difference(su.by, su.source.Columns())))
	}
	check(su.by)
	check(su.ons)
	for i := 0; i < len(su.cols); i++ {
		if su.cols[i] == "" {
			if su.ons[i] == "" {
				su.cols[i] = "count"
			} else {
				su.cols[i] = su.ops[i] + "_" + su.ons[i]
			}
		}
	}
	su.wholeRow = su.minmax1() && ssset.Contains(su.source.Keys(), su.ons)
}

func check(cols []string) {
	for _, c := range cols {
		if strings.HasSuffix(c, "_lower!") {
			panic("can't summarize _lower! fields")
		}
	}
}

func (su *Summarize) minmax1() bool {
	if len(su.by) > 0 || len(su.ops) != 1 {
		return false
	}
	fn := str.ToLower(su.ops[0])
	return fn == "min" || fn == "max"
}

func (su *Summarize) String() string {
	s := paren(su.source) + " SUMMARIZE "
	if len(su.by) > 0 {
		s += str.Join(", ", su.by) + ", "
	}
	sep := ""
	for i := range su.cols {
		s += sep
		sep = ", "
		if su.cols[i] != "" {
			s += su.cols[i] + " = "
		}
		s += su.ops[i]
		if su.ops[i] != "count" {
			s += " " + su.ons[i]
		}
	}
	return s
}

func (su *Summarize) Columns() []string {
	if su.wholeRow {
		return sset.Union(su.cols, su.source.Columns())
	}
	return sset.Union(su.by, su.cols)
}

func (su *Summarize) Keys() [][]string {
	return projectKeys(su.source.Keys(), su.by)
}

func (su *Summarize) Indexes() [][]string {
	if len(su.by) == 0 || containsKey(su.by, su.source.Keys()) {
		return su.source.Indexes()
	}
	var idxs [][]string
	for _, src := range su.source.Indexes() {
		if sset.StartsWithSet(src, su.by) {
			idxs = append(idxs, src)
		}
	}
	return idxs
}

// containsKey returns true if a set of columns contain one of the keys
func containsKey(cols []string, keys [][]string) bool {
	for _, key := range keys {
		if sset.Subset(cols, key) {
			return true
		}
	}
	return false
}

func (su *Summarize) Transform() Query {
	su.source = su.source.Transform()
	return su
}
