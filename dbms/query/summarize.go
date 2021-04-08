// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/runtime"
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
	summarizeApproach
}

type summarizeApproach struct {
	strategy sumStrategy
	index    []string
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
	s := parenQ2(su.source) + " SUMMARIZE"
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

func (su *Summarize) nrows() int {
	nr := su.source.nrows()
	if len(su.by) == 0 {
		nr = 1
	} else if !containsKey(su.by, su.source.Keys()) {
		nr /= 2 // ???
	}
	return nr
}

func (su *Summarize) rowSize() int {
	return su.source.rowSize() + len(su.cols)*8
}

func (su *Summarize) Updateable() bool {
	return false // override Query1 source.Updateable
}

func (*Summarize) Output(runtime.Record) {
	panic("can't output to this query")
}

func (su *Summarize) Transform() Query {
	su.source = su.source.Transform()
	return su
}

func (su *Summarize) optimize(mode Mode, index []string) (Cost, interface{}) {
	seqCost, seqApp := su.seqCost(mode, index)
	idxCost, idxApp := su.idxCost(mode)
	mapCost, mapApp := su.mapCost(mode, index)
	return min3(seqCost, seqApp, idxCost, idxApp, mapCost, mapApp)
}

func (su *Summarize) seqCost(mode Mode, index []string) (Cost, interface{}) {
	approach := &summarizeApproach{strategy: sumSeq}
	if len(su.by) == 0 || containsKey(su.by, su.source.Keys()) {
		if len(su.by) != 0 {
			approach.index = index
		}
		cost := Optimize(su.source, mode, su.index)
		return cost, approach
	}
	best := su.bestPrefixed(su.sourceIndexes(index), su.by, mode)
	if best.cost >= impossible {
		return impossible, nil
	}
	approach.index = best.index
	return best.cost, approach
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

func (su *Summarize) idxCost(mode Mode) (Cost, interface{}) {
	if !su.minmax1() {
		return impossible, nil
	}
	// dividing by nrecords since we're only reading one record
	nr := ints.Max(1, su.source.nrows())
	cost := Optimize(su.source, mode, su.ons) / nr
	approach := &summarizeApproach{strategy: sumIdx, index: su.ons}
	return cost, approach
}

func (su *Summarize) mapCost(mode Mode, index []string) (Cost, interface{}) {
	// can only provide 'by' as index
	if !str.List(su.by).HasPrefix(index) {
		return impossible, nil
	}
	cost := Optimize(su.source, mode, nil)
	cost += cost / 2 // add 50% for map overhead
	approach := &summarizeApproach{strategy: sumMap}
	return cost, approach
}

func (su *Summarize) setApproach(_ []string, approach interface{}, tran QueryTran) {
	su.summarizeApproach = *approach.(*summarizeApproach)
	su.source = SetApproach(su.source, su.index, tran)
}
