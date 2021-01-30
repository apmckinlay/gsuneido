// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/util/sset"
)

type Remove struct {
	Query1
	columns []string
}

func (rm *Remove) String() string {
	return paren(rm.source) + " REMOVE " + strings.Join(rm.columns, ", ")
}

func (rm *Remove) Transform() Query {
	cols := sset.Difference(rm.source.Columns(), rm.columns)
	p := &Project{Query1: rm.Query1, columns: cols}
	return p.Transform()
}
