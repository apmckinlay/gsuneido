// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"slices"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/tsc"
)

type Rename struct {
	from []string
	to   []string
	Query1
}

func NewRename(src Query, from, to []string) *Rename {
	srcCols := src.Columns()
	checkRename(srcCols, from, to)
	r := &Rename{Query1: Query1{source: src}, from: from, to: to}
	r.renameDependencies(srcCols)
	r.header = r.getHeader()
	r.keys = r.renameIndexes(src.Keys())
	r.indexes = r.renameIndexes(src.Indexes())
	r.setNrows(src.Nrows())
	r.rowSiz.Set(src.rowSize())
	r.fast1.Set(src.fastSingle())
	r.singleTbl.Set(src.SingleTable())
	r.lookCost.Set(src.lookupCost())
	return r
}

func checkRename(srcCols []string, from []string, to []string) {
	cols := slc.Clone(srcCols)
	for i, f := range from {
		j := slices.Index(cols, f)
		if j == -1 {
			panic("rename: nonexistent column: " + f)
		}
		t := to[i]
		if slices.Contains(cols, t) {
			panic("rename: column already exists: " + t)
		}
		cols[j] = t
	}
}

func (r *Rename) getHeader() *Header {
	flds := r.renameIndexes(r.source.Header().Fields)
	cols := r.renameFwd(r.source.Columns())
	return NewHeader(flds, cols)
}

func (r *Rename) renameDependencies(src []string) {
	r.from = slices.Clip(r.from)
	r.to = slices.Clip(r.to)
	for i, f := range r.from {
		deps := f + "_deps"
		if slices.Contains(src, deps) && !slices.Contains(r.from, deps) {
			r.from = append(r.from, deps)
			r.to = append(r.to, r.to[i]+"_deps")
		}
	}
}

func (r *Rename) String() string {
	sep := ""
	var sb strings.Builder
	sb.WriteString("rename ")
	for i, from := range r.from {
		sb.WriteString(sep)
		sb.WriteString(from)
		sb.WriteString(" to ")
		sb.WriteString(r.to[i])
		sep = ", "
	}
	return sb.String()
}

func (r *Rename) renameFwd(list []string) []string {
	cloned := false
	for i := range r.from {
		if j := slices.Index(list, r.from[i]); j != -1 {
			if !cloned {
				list = slc.Clone(list)
				cloned = true
			}
			list[j] = r.to[i]
		}
	}
	return list
}

func (r *Rename) renameRev(list []string) []string {
	cloned := false
	for i := len(r.to) - 1; i >= 0; i-- {
		if j := slices.Index(list, r.to[i]); j != -1 {
			if !cloned {
				list = slc.Clone(list)
				cloned = true
			}
			list[j] = r.from[i]
		}
	}
	return list
}

func (r *Rename) renameIndexes(idxs [][]string) [][]string {
	idxs2 := make([][]string, len(idxs))
	for i, ix := range idxs {
		idxs2[i] = r.renameFwd(ix)
	}
	return idxs2
}

func (r *Rename) Fixed() []Fixed {
	if r.fixed == nil {
		r.fixed = r.source.Fixed()
		cloned := false
		for i, from := range r.from {
			for j, fxd := range r.fixed {
				if fxd.col == from {
					if !cloned {
						r.fixed = slc.Clone(r.fixed)
						cloned = true
					}
					r.fixed[j] = Fixed{col: r.to[i], values: fxd.values}
					break
				}
			}
		}
		if r.fixed == nil {
			r.fixed = []Fixed{}
		}
	}
	return r.fixed
}

func (r *Rename) Transform() Query {
	src := r.source.Transform()

	// remove empty Renames
	if len(r.from) == 0 {
		return src
	}
	// combine Renames
	from := r.from
	to := r.to
	if r2, ok := src.(*Rename); ok {
		from = slc.With(r2.from, from...)
		to = slc.With(r2.to, to...)
		src = r2.source
	}
	if _, ok := src.(*Nothing); ok {
		return NewNothing(r)
	}
	if len(from) == 0 {
		return src
	}
	if src != r.source {
		r = NewRename(src, from, to)
	}
	return r
}

func (r *Rename) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	fixcost, varcost := Optimize(r.source, mode, r.renameRev(index), frac)
	return fixcost, varcost, nil
}

func (r *Rename) setApproach(index []string, frac float64, _ any, tran QueryTran) {
	r.source = SetApproach(r.source, r.renameRev(index), frac, tran)
	r.header = r.getHeader()
}

// execution --------------------------------------------------------

func (r *Rename) Get(th *Thread, dir Dir) Row {
	defer func(t uint64) { r.tget += tsc.Read() - t }(tsc.Read())
	r.ngets++
	return r.source.Get(th, dir)
}

func (r *Rename) Select(cols, vals []string) {
	r.source.Select(r.renameRev(cols), vals)
}

func (r *Rename) Lookup(th *Thread, cols, vals []string) Row {
	return r.source.Lookup(th, r.renameRev(cols), vals)
}

func (r *Rename) Simple(th *Thread) []Row {
	return r.source.Simple(th)
}
