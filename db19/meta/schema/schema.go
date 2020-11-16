// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package schema is a separate package so it can be used by query parsing.
package schema

import (
	"strings"

	"github.com/apmckinlay/gsuneido/db19/index/ixspec"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Schema struct {
	Table string
	// Columns are the physical fields in the records, in order
	Columns []string
	// Derived are the rules (capitalized) and _lower!
	Derived []string
	Indexes []Index
}

type Index struct {
	Columns []string
	Ixspec  ixspec.T
	// Mode is 'k' for key, 'i' for index, 'u' for unique index
	Mode      int
	Fktable   string
	Fkmode    int
	Fkcolumns []string
}

// fkmode bits
const (
	Block          = 0
	CascadeUpdates = 1
	CascadeDeletes = 2
	Cascade        = CascadeUpdates | CascadeDeletes
)

func (sc *Schema) String() string {
	var sb strings.Builder
	if sc.Columns != nil || sc.Derived != nil {
		var cb str.CommaBuilder
		for _, col := range sc.Columns {
			cb.Add(col)
		}
		for _, col := range sc.Derived {
			cb.Add(col)
		}
		sb.WriteString("(")
		sb.WriteString(cb.String())
		sb.WriteString(") ")
	}
	sep := ""
	for i := range sc.Indexes {
		sb.WriteString(sep)
		sb.WriteString(sc.Indexes[i].String())
		sep = " "
	}
	return sb.String()
}

func (ix *Index) String() string {
	s := map[int]string{'k': "key", 'i': "index", 'u': "index unique"}[ix.Mode]
	s += str.Join("(,)", ix.Columns...)
	if ix.Fktable != "" {
		s += " in " + ix.Fktable
		if len(ix.Fkcolumns) > 0 {
			sep := "("
			for _, f := range ix.Fkcolumns {
				s += sep + f
				sep = ","
			}
			s += ")"
		}
		if ix.Fkmode&Cascade != 0 {
			s += " cascade"
			if ix.Fkmode == CascadeUpdates {
				s += " update"
			}
		}
	}
	return s
}
