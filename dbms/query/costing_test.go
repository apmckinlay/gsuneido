// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"slices"
	"strconv"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/slc"
)

const ncosting = 10000

func TestCosting_Table(t *testing.T) {
	for range ncosting {
		ft := testFT()
		q := ft.NewFuzzTable()
		q = q.Transform()
		req := NoneReq(1)
		fixcost, varcost := Optimize(q, ReadMode, req)
		q = SetApproach(q, req, ft.rt)
		q.SetTran(ft.rt)
		iterate(q)
		tbl := q.(*Table)
		actualCost := Cost(tbl.ngets) * rowCost(tbl.index)
		assert.This(fixcost + varcost).Is(actualCost)
		ft.db.Close()
	}
}

func TestCosting_Times(t *testing.T) {
	for range ncosting {
		ft := testFT()
		q := fuzzTimes(ft)
		src1 := q.(*Times).source1.(*Table)
		q, fixcost, varcost := costingSetup(q, ft)
		tbl1 := q.(*Times).source1.(*Table)
		tbl2 := q.(*Times).source2.(*Table)
		actualCost := Cost(tbl1.ngets)*rowCost(tbl1.index) +
			Cost(tbl2.ngets)*rowCost(tbl2.index)
		if tbl1 != src1 { // reversed
			actualCost += outOfOrder
		}
		assert.This(fixcost + varcost).Is(actualCost)
		ft.db.Close()
	}
}

func TestCosting_Project(t *testing.T) {
	for range ncosting {
		ft := testFT()
		q := fuzzProjectForCosting(ft)
		q, fixcost, varcost := costingSetup(q, ft)
		proj := q.(*Project)
		tbl := proj.source.(*Table)
		actualCost := Cost(tbl.ngets) * rowCost(tbl.index)
		if proj.strat == projMap {
			actualCost += Cost(tbl.ngets * 20)
		}
		assert.This(fixcost + varcost).Is(actualCost)
		ft.db.Close()
	}
}

// fuzzProjectForCosting creates a randomized Project
// that exercises both projSeq and projMap strategies,
// with data structured to have projDupDiv duplicates per group.
func fuzzProjectForCosting(ft *FT) Query {
	for { // loop is so we can reject unusable and restart
		b := ft.newFT().Sizes(ftMaxRows, 1, ftMaxIndexes).NoEmptyKey().construct()

		// truncate to make len(data) divisible by projDupDiv
		b.data = b.data[:len(b.data)-len(b.data)%projGrpDiv]

		projCols := randomProjectCols(ft.rnd, b.columns, b.indexes)

		if set.Equal(projCols, b.columns) {
			continue // Transform would eliminate
		}

		if indexContainsKey(projCols, b.keys) == nil {
			if slices.ContainsFunc(b.keys,
				func(key []string) bool { return !set.Disjoint(key, projCols) }) {
				// can't alter key fields so we can't make groups duplicate
				continue
			}
			projIdxs := make([]int, len(projCols))
			for i, col := range projCols {
				projIdxs[i] = b.colIndex[col]
			}
			for ri := range b.data {
				g := ri / projGrpDiv
				for i, ci := range projIdxs {
					b.data[ri][ci] = projCols[i] + "_" + strconv.Itoa(g)
				}
			}
			slc.Shuffle(ft.rnd, b.data)
		}

		src := b.finish()
		return NewProject(src, projCols)
	}
}

func TestCosting_Summarize(t *testing.T) {
	for range ncosting {
		ft := testFT()
		q := fuzzSummarizeForCosting(ft)
		q, fixcost, varcost := costingSetup(q, ft)
		su := q.(*Summarize)
		tbl, ok := su.source.(*Table)
		if !ok {
			ft.db.Close()
			continue
		}
		var actualCost Cost
		switch su.strat {
		case sumTbl:
			actualCost = 1
		case sumSeq:
			actualCost = Cost(tbl.ngets) * rowCost(tbl.index)
			nrows, _ := su.Nrows()
			fmt.Println("table:", tbl.ngets, "estimated:", nrows, "actual:", su.ngets)
		case sumIdx:
			actualCost = rowCost(tbl.index)
		case sumMap:
			actualCost = Cost(tbl.ngets)*rowCost(tbl.index) + Cost(tbl.ngets)*20
		}
		diff := fixcost + varcost - actualCost
		if diff < -1 || diff > 1 {
			t.Fatalf("strat: %d, fixcost+varcost: %d, actualCost: %d", su.strat, fixcost+varcost, actualCost)
		}
		ft.db.Close()
	}
}

func fuzzSummarizeForCosting(ft *FT) Query {
	for { // loop is so we can reject unusable and restart
		b := ft.newFT().Sizes(ftMaxRows, 1, ftMaxIndexes).NoEmptyKey().construct()
		if len(b.columns) == 0 {
			continue
		}

		switch ft.rnd.IntN(5) {
		case 0: // sumTbl: no by, single count
			return NewSummarize(b.finish(), "", nil,
				[]string{""}, []string{"count"}, []string{""})

		case 1: // sumIdx: no by, single min/max
			var col string
			if len(b.indexes) > 0 && ft.rnd.IntN(2) == 0 {
				col = random(b.indexes, ft.rnd)[0]
			} else {
				col = random(b.columns, ft.rnd)
			}
			op := "min"
			if ft.rnd.IntN(2) == 0 {
				op = "max"
			}
			return NewSummarize(b.finish(), "", nil,
				[]string{""}, []string{op}, []string{col})

		case 2: // sumSeq with no by: total or average
			col := random(b.columns, ft.rnd)
			op := "total"
			if ft.rnd.IntN(2) == 0 {
				op = "average"
			}
			return NewSummarize(b.finish(), "", nil,
				[]string{""}, []string{op}, []string{col})

		case 3: // sumSeq with by as index prefix
			if len(b.indexes) == 0 {
				continue
			}
			idx := random(b.indexes, ft.rnd)
			if len(idx) == 0 {
				continue
			}
			prefixLen := 1 + ft.rnd.IntN(len(idx))
			by := slices.Clone(idx[:prefixLen])
			if slices.ContainsFunc(b.keys,
				func(key []string) bool { return !set.Disjoint(key, by) }) {
				continue
			}
			const sumDupDiv = 10
			b.data = b.data[:len(b.data)-len(b.data)%sumDupDiv]
			if len(b.data) == 0 {
				continue
			}
			byIdxs := make([]int, len(by))
			for i, col := range by {
				byIdxs[i] = b.colIndex[col]
			}
			for ri := range b.data {
				g := ri / sumDupDiv
				for i, ci := range byIdxs {
					b.data[ri][ci] = by[i] + "_" + strconv.Itoa(g)
				}
			}
			slc.Shuffle(ft.rnd, b.data)
			return NewSummarize(b.finish(), "", by,
				[]string{""}, []string{"count"}, []string{""})

		case 4: // sumMap: by not an index prefix and not a key
			col := random(b.columns, ft.rnd)
			by := []string{col}
			isIdxPrefix := false
			for _, idx := range b.indexes {
				if len(idx) > 0 && idx[0] == col {
					isIdxPrefix = true
					break
				}
			}
			if isIdxPrefix {
				continue
			}
			if hasKey(by, b.keys, nil) {
				continue
			}
			if slices.ContainsFunc(b.keys,
				func(key []string) bool { return slices.Contains(key, col) }) {
				continue
			}
			b.data = b.data[:len(b.data)-len(b.data)%sumGrpDiv]
			if len(b.data) == 0 {
				continue
			}
			ci := b.colIndex[col]
			for ri := range b.data {
				g := ri / sumGrpDiv
				b.data[ri][ci] = col + "_" + strconv.Itoa(g)
			}
			slc.Shuffle(ft.rnd, b.data)
			return NewSummarize(b.finish(), "", by,
				[]string{""}, []string{"count"}, []string{""})
		}
	}
}

func costingSetup(q Query, ft *FT) (Query, Cost, Cost) {
	q = q.Transform()
	req := NoneReq(1)
	fixcost, varcost := Optimize(q, ReadMode, req)
	q = SetApproach(q, req, ft.rt)
	q.SetTran(ft.rt)
	iterate(q)
	return q, fixcost, varcost
}

// rowCost matches the per-get cost used by Table.costFor
func rowCost(index []string) Cost {
	return tableFast + Cost(len(index))*colsBias
}

func iterate(q Query) {
	th := &Thread{}
	for row := q.Get(th, Next); row != nil; row = q.Get(th, Next) {
	}
}
