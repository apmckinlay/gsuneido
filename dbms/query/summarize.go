// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/util/ints"
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
	strategy sumStrategy
	via      []string
}

type sumStrategy int

const (
	sumSeq sumStrategy = iota + 1
	sumMap
	sumIdx
)

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
	s := paren(su.source) + " SUMMARIZE"
	switch su.strategy {
	case sumSeq:
		s += "-SEQ"
	case sumMap:
		s += "-MAP"
	case sumIdx:
		s += "-IDX"
	}
	if len(su.by) > 0 {
		s += " " + str.Join(", ", su.by) + ","
	}
	sep := " "
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

func (su *Summarize) Updateable() bool {
	return false // override Query1 source.Updateable
}

func (su *Summarize) Transform() Query {
	su.source = su.source.Transform()
	return su
}

func (su *Summarize) optimize(mode Mode, index []string, act action) Cost {
	seqCost := su.seqCost(mode, index, assess)
	idxCost := su.idxCost(mode, assess)
	mapCost := su.mapCost(mode, index, assess)

	if act == assess {
		return ints.Min(seqCost, ints.Min(idxCost, mapCost))
	}

	if seqCost <= idxCost && seqCost <= mapCost {
		return su.seqCost(mode, index, freeze)
	} else if idxCost <= mapCost {
		return su.idxCost(mode, freeze)
	} else {
		return su.mapCost(mode, index, freeze)
	}
}

func (su *Summarize) seqCost(mode Mode, index []string, act action) Cost {
	if act == freeze {
		su.strategy = sumSeq
	}
	if len(su.by) == 0 || containsKey(su.by, su.source.Keys()) {
		if len(su.by) == 0 {
			su.via = nil
		} else {
			su.via = index
		}
		return Optimize(su.source, mode, su.via, freeze)
	}
	best := su.bestPrefixed(su.sourceIndexes(index), su.by, mode)
	if act == assess || best.cost >= impossible {
		return best.cost
	}
	su.via = best.index
	// optimize1 to bypass temp index
	return optimize1(su.source, mode, best.index, freeze)
}
func (su *Summarize) sourceIndexes(index []string) [][]string {
	if index == nil {
		return su.source.Indexes()
	}
	fixed := su.source.Fixed()
	var indexes [][]string
	for _, idx := range su.source.Indexes() {
		if su.prefixed(idx, index, fixed) {
			indexes = append(indexes, idx)
		}
	}
	return indexes
}

func (su *Summarize) idxCost(mode Mode, act action) Cost {
	if !su.minmax1() {
		return impossible
	}
	// optimize1 to bypass temp index
	// dividing by nrecords since we're only reading one record
	nr := ints.Max(1, su.source.nrows())
	cost := optimize1(su.source, mode, su.ons, freeze) / nr
	if act == freeze {
		su.strategy = sumIdx
		su.via = su.ons
	}
	return cost
}

func (su *Summarize) mapCost(mode Mode, index []string, act action) Cost {
	// can only provide 'by' as index
	if !str.List(su.by).HasPrefix(index) {
		return impossible
	}
	// optimize1 to bypass temp index
	cost := optimize1(su.source, mode, nil, freeze)
	// add 50% for map overhead
	cost += cost / 2
	if act == freeze {
		su.strategy = sumMap
	}
	return cost
}
