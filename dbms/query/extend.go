// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import "github.com/apmckinlay/gsuneido/util/sset"

type Extend struct {
	Query1
	cols  []string
	exprs []Expr
}

func (e *Extend) String() string {
	s := e.Query1.String() + " extend "
	sep := ""
	for i, c := range e.cols {
		s += sep + c
		sep = ", "
		if e.exprs[i] != nil {
			s += " = " + e.exprs[i].String()
		}
	}
	return s
}

func (e *Extend) Columns() []string {
	return sset.Union(e.source.Columns(), e.cols)
}

func (e *Extend) Transform() Query {
	// remove empty Extends
	if len(e.cols) == 0 {
		return e.source.Transform()
	}
	// combine Renames
	for e2, ok := e.source.(*Extend); ok; e2, ok = e.source.(*Extend) {
		e.cols = append(e2.cols, e.cols...)
		e.exprs = append(e2.exprs, e.exprs...)
		//TODO e.init()
		e.source = e2.source
	}
	e.source = e.source.Transform()
	return e
}
