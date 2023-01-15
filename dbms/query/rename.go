// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/exp/slices"
)

type Rename struct {
	Query1
	from []string
	to   []string
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
	return r
}

func (r *Rename) renameDependencies(src []string) {
	r.from = slices.Clip(r.from)
	r.to = slices.Clip(r.to)
	for i,n := 0, len(r.from); i < n; i++ {
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

func (r *Rename) Columns() []string {
	return slc.Replace(r.source.Columns(), r.from, r.to)
}

func (r *Rename) Keys() [][]string {
	return renameIndexes(r.source.Keys(), r.from, r.to)
}

func (r *Rename) Indexes() [][]string {
	return renameIndexes(r.source.Indexes(), r.from, r.to)
}

func renameIndexes(idxs [][]string, from, to []string) [][]string {
	idxs2 := make([][]string, len(idxs))
	for i, ix := range idxs {
		idxs2[i] = slc.Replace(ix, from, to)
	}
	return idxs2
}

func (r *Rename) Fixed() []Fixed {
	fixed := r.source.Fixed()
	result := make([]Fixed, len(fixed))
	for i, fxd := range fixed {
		j := slices.Index(r.from, fxd.col)
		if j == -1 {
			result[i] = fxd
		} else {
			result[i] = Fixed{col: r.to[j], values: fxd.values}
		}
	}
	return result
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
		from = slc.With(r2.from, from...)
		to = slc.With(r2.to, to...)
		dst := 0
	outer:
		for i := 0; i < len(from); i++ {
			for j := 0; j < i; j++ {
				if to[j] == from[i] {
					to[j] = to[i]
					continue outer
				}
			}
			if i > dst {
				from[dst] = from[i]
				to[dst] = to[i]
			}
			dst++
		}
		from = from[:dst]
		to = to[:dst]
		src = r2.source
	}
	src = src.Transform()
	if _, ok := src.(*Nothing); ok {
		return NewNothing(slc.Replace(src.Columns(), from, to))
	}
	return NewRename(src, from, to)
}

func (r *Rename) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	fixcost, varcost := Optimize(r.source, mode,
		slc.Replace(index, r.to, r.from), frac)
	return fixcost, varcost, nil
}

func (r *Rename) setApproach(index []string, frac float64, _ any, tran QueryTran) {
	r.source = SetApproach(r.source, slc.Replace(index, r.to, r.from), frac, tran)
}

// execution --------------------------------------------------------

func (r *Rename) Header() *runtime.Header {
	hdr := r.source.Header()
	cols := slc.Replace(hdr.Columns, r.from, r.to)
	flds := renameIndexes(hdr.Fields, r.from, r.to)
	return runtime.NewHeader(flds, cols)
}

func (r *Rename) Get(th *runtime.Thread, dir runtime.Dir) runtime.Row {
	return r.source.Get(th, dir)
}

func (r *Rename) Select(cols, vals []string) {
	r.source.Select(slc.Replace(cols, r.to, r.from), vals)
}

func (r *Rename) Lookup(th *runtime.Thread, cols, vals []string) runtime.Row {
	return r.source.Lookup(th, slc.Replace(cols, r.to, r.from), vals)
}
