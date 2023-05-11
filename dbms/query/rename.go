// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/exp/slices"
)

type Rename struct {
	from []string
	to   []string
	Query1
}

func NewRename(src Query, from, to []string) *Rename {
	srcCols := src.Columns()
	if !set.Subset(srcCols, from) {
		panic("rename: nonexistent column(s): " +
			str.Join(", ", set.Difference(from, srcCols)))
	}
	if !set.Disjoint(srcCols, to) {
		panic("rename: column(s) already exist: " +
			str.Join(", ", set.Intersect(srcCols, to)))
	}
	r := &Rename{Query1: Query1{source: src}, from: from, to: to}
	r.renameDependencies(srcCols)
	r.header = r.getHeader()
	r.keys = renameIndexes(src.Keys(), r.from, r.to)
	r.indexes = renameIndexes(src.Indexes(), r.from, r.to)
	return r
}

func (r *Rename) getHeader() *Header {
	flds := renameIndexes(r.source.Header().Fields, r.from, r.to)
	cols := slc.Replace(r.source.Columns(), r.from, r.to)
	return NewHeader(flds, cols)
}

func (r *Rename) renameDependencies(src []string) {
	r.from = slices.Clip(r.from)
	r.to = slices.Clip(r.to)
	for i, n := 0, len(r.from); i < n; i++ {
		deps := r.from[i] + "_deps"
		if slices.Contains(src, deps) && !slices.Contains(r.from, deps) {
			r.from = append(r.from, deps)
			r.to = append(r.to, r.to[i]+"_deps")
		}
	}
}

func (r *Rename) String() string {
	return parenQ2(r.source) + " " + r.stringOp()
}

func (r *Rename) stringOp() string {
	sep := ""
	var sb strings.Builder
	sb.WriteString("RENAME ")
	for i, from := range r.from {
		sb.WriteString(sep)
		sb.WriteString(from)
		sb.WriteString(" to ")
		sb.WriteString(r.to[i])
		sep = ", "
	}
	return sb.String()
}

func renameIndexes(idxs [][]string, from, to []string) [][]string {
	idxs2 := make([][]string, len(idxs))
	for i, ix := range idxs {
		idxs2[i] = slc.Replace(ix, from, to)
	}
	return idxs2
}

func (r *Rename) Fixed() []Fixed {
	if r.fixed == nil {
		srcFix := r.source.Fixed()
		result := make([]Fixed, len(srcFix))
		for i, fxd := range srcFix {
			j := slices.Index(r.from, fxd.col)
			if j == -1 {
				result[i] = fxd
			} else {
				result[i] = Fixed{col: r.to[j], values: fxd.values}
			}
		}
		r.fixed = result
		assert.That(r.fixed != nil)
	}
	return r.fixed
}

func (r *Rename) Transform() Query {
	// remove empty Renames
	if len(r.from) == 0 {
		return r.source.Transform()
	}
	// combine Renames
	src := r.source
	from := r.from
	to := r.to
	for r2, ok := src.(*Rename); ok; r2, ok = src.(*Rename) {
		from, to = mergeRename(r2.from, r2.to, from, to)
		src = r2.source
	}
	src = src.Transform()
	if _, ok := src.(*Nothing); ok {
		return NewNothing(slc.Replace(src.Columns(), from, to))
	}
	if len(from) == 0 {
		return src
	}
	return NewRename(src, from, to)
}

func mergeRename(from1, to1, from2, to2 []string) (from, to []string) {
	from = slices.Clone(from1)
	to = slices.Clone(to1)
	for i, f := range from2 {
		t := to2[i]
		if j := slices.Index(to, f); j >= 0 {
			if t == from[j] {
				// rename back to original, so remove
				from = slices.Delete(from, j, j+1)
				to = slices.Delete(to, j, j+1)
			} else {
				// rename again, so update first
				to[j] = t
			}
		} else {
			from = append(from, f)
			to = append(to, t)
		}
	}
	return from, to
}

func (r *Rename) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	fixcost, varcost := Optimize(r.source, mode,
		slc.Replace(index, r.to, r.from), frac)
	return fixcost, varcost, nil
}

func (r *Rename) setApproach(index []string, frac float64, _ any, tran QueryTran) {
	r.source = SetApproach(r.source, slc.Replace(index, r.to, r.from), frac, tran)
	r.header = r.getHeader()
}

// execution --------------------------------------------------------

func (r *Rename) Get(th *Thread, dir Dir) Row {
	return r.source.Get(th, dir)
}

func (r *Rename) Select(cols, vals []string) {
	r.source.Select(slc.Replace(cols, r.to, r.from), vals)
}

func (r *Rename) Lookup(th *Thread, cols, vals []string) Row {
	return r.source.Lookup(th, slc.Replace(cols, r.to, r.from), vals)
}
