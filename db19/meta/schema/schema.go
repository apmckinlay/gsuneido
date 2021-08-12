// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package schema is a separate package so it can be used by query parsing.
package schema

import (
	"strings"

	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/strs"
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
	Ixspec  ixkey.Spec
	// Mode is 'k' for key, 'i' for index, 'u' for unique index
	Mode int
	Fk   Fkey
	// FkToHere is other foreign keys that reference this one
	FkToHere []Fkey
}

type Fkey struct {
	Table   string
	Columns []string
	Mode    int
}

// Fkey mode bits
const (
	Block          = 0
	CascadeUpdates = 1
	CascadeDeletes = 2
	Cascade        = CascadeUpdates | CascadeDeletes
)

func (sc *Schema) String() string {
	var sb strings.Builder
	sb.WriteString(sc.Table)
	sb.WriteString(" ")
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
	s += strs.Join("(,)", ix.Columns)
	if ix.Fk.Table != "" {
		s += " in " + ix.Fk.Table
		if !strs.Equal(ix.Fk.Columns, ix.Columns) {
			s += strs.Join("(,)", ix.Fk.Columns)
		}
		if ix.Fk.Mode&Cascade != 0 {
			s += " cascade"
			if ix.Fk.Mode == CascadeUpdates {
				s += " update"
			}
		}
	}
	for _, fk := range ix.FkToHere {
		s += " from " + fk.Table + strs.Join("(,)", fk.Columns)
	}
	return s
}
