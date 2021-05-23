// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/strs"
)

type Rename struct {
	Query1
	from []string
	to   []string
}

func (r *Rename) Init() {
	r.Query1.Init()
	srcCols := r.source.Columns()
	if !sset.Subset(srcCols, r.from) {
		panic("rename: nonexistent column(s): " +
			str.Join(", ", sset.Difference(r.from, srcCols)))
	}
	if !sset.Disjoint(srcCols, r.to) {
		panic("rename: column(s) already exist: " +
			str.Join(", ", sset.Intersect(srcCols, r.to)))
	}
	r.renameDependencies(srcCols)
}

func (r *Rename) renameDependencies(src []string) {
	copy := false
	for i := 0; i < len(r.from); i++ {
		deps := r.from[i] + "_deps"
		if sset.Contains(src, deps) {
			if !copy {
				r.from = sset.Copy(r.from)
				r.to = sset.Copy(r.to)
				copy = true
			}
			r.from = append(r.from, deps)
			r.to = append(r.to, r.to[i]+"_deps")
		}
	}
}

func (r *Rename) String() string {
	sep := ""
	var sb strings.Builder
	for i, from := range r.from {
		sb.WriteString(sep)
		sb.WriteString(from)
		sb.WriteString(" to ")
		sb.WriteString(r.to[i])
		sep = ", "
	}
	return parenQ2(r.source) + " RENAME " + sb.String()
}

func (r *Rename) Columns() []string {
	return strs.Replace(r.source.Columns(), r.from, r.to)
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
		idxs2[i] = strs.Replace(ix, from, to)
	}
	return idxs2
}

func (r *Rename) Fixed() []Fixed {
	fixed := r.source.Fixed()
	result := make([]Fixed, len(fixed))
	for i, fxd := range fixed {
		j := strs.Index(r.from, fxd.col)
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
	for r2, ok := r.source.(*Rename); ok; r2, ok = r.source.(*Rename) {
		from := append(r2.from, r.from...)
		to := append(r2.to, r.to...)
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
		r.from = from[:dst]
		r.to = to[:dst]
		r.source = r2.source
	}
	r.source = r.source.Transform()
	return r
}

func (r *Rename) setApproach(index []string, _ interface{}, tran QueryTran) {
	r.source = SetApproach(r.source, index, tran)
}

// execution --------------------------------------------------------

func (r *Rename) Header() *runtime.Header {
	hdr := r.source.Header()
	cols := strs.Replace(hdr.Columns, r.from, r.to)
	flds := renameIndexes(hdr.Fields, r.from, r.to)
	return runtime.NewHeader(flds, cols)
}

func (r *Rename) Get(dir runtime.Dir) runtime.Row {
	return r.source.Get(dir)
}

func (r *Rename) Select(cols, vals []string) {
	r.source.Select(strs.Replace(cols, r.from, r.to), vals)
}
