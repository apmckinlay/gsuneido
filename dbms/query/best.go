// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"slices"

	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/slc"
)

type bestIndex struct {
	index   []string
	fixcost Cost
	varcost Cost
}

func newBestIndex() bestIndex {
	return bestIndex{fixcost: impossible, varcost: impossible}
}

// update returns true if the new values are the lowest cost so far
func (bi *bestIndex) update(index []string, fixcost, varcost Cost) bool {
	if fixcost+varcost < bi.fixcost+bi.varcost {
		*bi = bestIndex{index: index, fixcost: fixcost, varcost: varcost}
		return true
	}
	return false
}

func (bi *bestIndex) cost() int {
	return bi.fixcost + bi.varcost
}

func (bi *bestIndex) String() string {
	if bi.cost() >= impossible {
		return "impossible"
	}
	return fmt.Sprint("{", bi.index, " ",
		trace.Number(bi.fixcost), " + ", trace.Number(bi.varcost),
		" = ", trace.Number(bi.cost()), "}")
}

//-------------------------------------------------------------------

// bestGrouped finds the best index with cols (in any order) as a prefix
// taking fixed into consideration.
// It is used by Project, Summarize, and Join.
func bestGrouped(source Query, mode Mode, index []string, frac float64, cols []string) bestIndex {
	var indexes [][]string
	if index == nil {
		indexes = source.Indexes()
	} else {
		indexes = [][]string{index}
	}
	best := bestGrouped2(source, mode, indexes, frac, cols)
	if index == nil {
		fixcost, varcost := Optimize(source, mode, cols, frac)
		best.update(cols, fixcost, varcost)
	}
	return best
}

func bestGrouped2(source Query, mode Mode, indexes [][]string, frac float64, cols []string) bestIndex {
	fixed := source.Fixed()
	nColsUnfixed := countUnfixed(cols, fixed)
	best := newBestIndex()
	for _, idx := range indexes {
		if grouped(idx, cols, nColsUnfixed, fixed) {
			fixcost, varcost := Optimize(source, mode, idx, frac)
			best.update(idx, fixcost, varcost)
		}
	}
	return best
}

// bestLookupIndex finds the best index for nrows lookup operations.
// cols restricts candidates to those grouped by cols (for Join to-one);
// pass nil to allow any lookup-eligible index (for Intersect, Minus, Union).
// If no physical index qualifies, falls back to logical keys.
func bestLookupIndex(source Query, mode Mode, nrows int, frac float64, cols []string) bestIndex {
	fixed := source.Fixed()
	keys := source.Keys()
	best := newBestIndex()
	var nColsUnfixed int
	if cols != nil {
		nColsUnfixed = countUnfixed(cols, fixed)
	}
	for _, idx := range source.Indexes() {
		if lookupIndexEligible(idx, keys, fixed) &&
			(cols == nil || grouped(idx, cols, nColsUnfixed, fixed)) {
			fixcost, varcost := LookupCost(source, mode, idx, nrows, frac)
			best.update(idx, fixcost, varcost)
		}
	}
	if best.index != nil {
		return best
		// could check the fallback regardless, but in practice it rarely helps
	}
	// fallback: no qualifying physical index (e.g. system tables with nil Indexes)
	fallbackKeys := keys
	if cols != nil {
		fallbackKeys = [][]string{cols}
	}
	for _, k := range fallbackKeys {
		fixcost, varcost := LookupCost(source, mode, k, nrows, frac)
		best.update(k, fixcost, varcost)
	}
	return best
}

func lookupIndexEligible(index []string, keys [][]string, fixed Fixed) bool {
	for _, key := range keys {
		nColsUnfixed := countUnfixed(key, fixed)
		if nColsUnfixed == 0 || grouped(index, key, nColsUnfixed, fixed) {
			return true
		}
	}
	return false
}

func countUnfixed(cols []string, fixed Fixed) int {
	nunfixed := 0
	for _, col := range cols {
		if !fixed.Single(col) {
			nunfixed++
		}
	}
	return nunfixed
}

// grouped returns whether an index has cols (in any order) as a prefix
// taking fixed into consideration
func grouped(index []string, cols []string, nColsUnfixed int, fixed Fixed) bool {
	if len(index) < nColsUnfixed {
		return false
	}
	n := 0
	for _, col := range index {
		if fixed.Single(col) {
			continue
		}
		if !slices.Contains(cols, col) {
			return false
		}
		n++
		if n == nColsUnfixed {
			return true
		}
	}
	return false
}

// ordered returns whether an index supplies an order
// taking fixed into consideration.
// It is used by Where and Sort.
func ordered(index []string, order []string, fixed Fixed) bool {
	return orderedn(index, order, fixed) >= len(order)
}

// orderedn returns the number of fields in order that are satisfied
func orderedn(index []string, order []string, fixed Fixed) int {
	i := 0
	o := 0
	in := len(index)
	on := len(order)
	for i < in && o < on {
		if index[i] == order[o] {
			o++
			i++
		} else if fixed.Single(index[i]) {
			i++
		} else if fixed.Single(order[o]) {
			o++
		} else {
			return o
		}
	}
	for o < on && fixed.Single(order[o]) {
		o++
	}
	return o
}

//-------------------------------------------------------------------

// withoutDupsOrSupersets simplifies a set of keys
// by removing duplicates and supersets
func withoutDupsOrSupersets(keys [][]string) [][]string {
	om := newOptMod(keys)
outer:
	for _, k1 := range keys {
		for _, k2 := range keys {
			if len(k1) > len(k2) && set.Subset(k1, k2) {
				continue outer // skip/exclude k1 - superset
			}
		}
		if !slc.ContainsFn(om.result(), k1, set.Equal[string]) { // exclude duplicates
			om.add(k1)
		}
	}
	return om.result()
}

// optmod is useful when building a new version
// which is likely to be the same as the original.
// It avoids constructing a new version unless there are changes,
// without having to redundantly check in advance.
type optmod struct {
	orig [][]string
	mod  [][]string
	i    int
}

func newOptMod(orig [][]string) *optmod {
	return &optmod{orig: orig}
}

func (b *optmod) add(x []string) {
	if b.mod == nil {
		if b.i < len(b.orig) && set.Equal(x, b.orig[b.i]) {
			b.i++ // same as orig
			return
		}
		b.mod = append(b.mod, b.orig[:b.i]...)
	}
	b.mod = append(b.mod, x)
}

func (b *optmod) result() [][]string {
	if b.mod == nil {
		return b.orig[:b.i:b.i]
	}
	return slices.Clip(b.mod)
}

// isEmptyKey returns true if the indexes have a single empty index
func isEmptyKey(indexes [][]string) bool {
	return len(indexes) == 1 && len(indexes[0]) == 0
}
