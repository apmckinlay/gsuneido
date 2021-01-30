// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Rename struct {
	Query1
	from []string
	to   []string
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
	return paren(r.source) + " RENAME " + sb.String()
}

func (r *Rename) Columns() []string {
	return renameCols(r.source.Columns(), r.from, r.to)
}

func renameCols(cols, from, to []string) []string {
	newCols := sset.Copy(cols)
	for i := 0; i < len(cols); i++ {
		j := str.List(from).Index(cols[i])
		if j != -1 {
			newCols[i] = to[j]
		}
	}
	return newCols
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
