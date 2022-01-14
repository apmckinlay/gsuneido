// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package schema is a separate package so it can be used by query parsing.
package schema

import (
	"sort"
	"strings"

	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/hash"
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
	// FkToHere is other foreign keys that reference this index
	FkToHere []Fkey // filled in by meta
}

type Fkey struct {
	Table   string
	Columns []string
	IIndex  int
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
	return sc.string(false)
}

// String2 includes fkToHere information
func (sc *Schema) String2() string {
	return sc.string(true)
}

func (sc *Schema) string(fktohere bool) string {
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
		sb.WriteString(sc.Indexes[i].string(fktohere))
		sep = " "
	}
	return sb.String()
}

func (ix *Index) String() string {
	return ix.string(false)
}

func (ix *Index) string(fktohere bool) string {
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
	if fktohere {
		toHere := make([]string, len(ix.FkToHere))
		for i, fk := range ix.FkToHere {
			toHere[i] = " from " + fk.Table + strs.Join("(,)", fk.Columns)
		}
		// sort for consistency in tests
		sort.Slice(toHere, func(i, j int) bool { return toHere[i] < toHere[j] })
		s += strs.Join("", toHere)
	}
	return s
}

// FindIndex returns a pointer to the Index with the given columns
// or else nil if not found
func (sc *Schema) FindIndex(cols []string) *Index {
	for i := range sc.Indexes {
		idx := &sc.Indexes[i]
		if strs.Equal(cols, idx.Columns) {
			return idx
		}
	}
	return nil
}

// IIndex returns the position of the index with the given columns
// or else it panics
func (sc *Schema) IIndex(cols []string) int {
	for i := range sc.Indexes {
		idx := &sc.Indexes[i]
		if strs.Equal(cols, idx.Columns) {
			return i
		}
	}
	panic("IIndex not found" + sc.Table + strs.Join(",", cols))
}

func (ix *Index) Equal(iy *Index) bool {
	return strs.Equal(ix.Columns, iy.Columns) &&
		ix.Mode == iy.Mode &&
		ix.Fk.Table == iy.Fk.Table &&
		ix.Fk.Mode == iy.Fk.Mode &&
		strs.Equal(ix.Fk.Columns, iy.Fk.Columns)
}

func (sc *Schema) Check() {
	sc.checkLower()
	sc.checkForKey()
	CheckIndexes(sc.Table, sc.Columns, sc.Indexes)
}

func (sc *Schema) checkLower() {
	for _, col := range sc.Derived {
		if strings.HasSuffix(col, "_lower!") &&
			!strs.Contains(sc.Columns, strings.TrimSuffix(col, "_lower!")) {
			panic("_lower! nonexistent column: " +
				strings.TrimSuffix(col, "_lower!"))
		}
	}
}

func (sc *Schema) checkForKey() {
	hasKey := false
	for i := range sc.Indexes {
		ix := &sc.Indexes[i]
		if ix.Mode == 'k' {
			hasKey = true
		}
	}
	if !hasKey {
		panic("key required in " + sc.Table)
	}
}

func CheckIndexes(table string, cols []string, idxs []Index) {
	for i := range idxs {
		ix := &idxs[i]
		if ix.Mode != 'k' && len(ix.Columns) == 0 {
			panic("index columns must not be empty")
		}
		for _, col := range ix.Columns {
			if !strs.Contains(cols, col) &&
				!strs.Contains(cols, str.RemoveSuffix(col, "_lower!")) {
				panic("invalid index column: " +
					col + " in " + table)
			}
		}
	}
}

func (sc *Schema) Cksum() uint32 {
	cksum := hash.HashString(sc.Table)
	for _, col := range sc.Columns {
		cksum += hash.HashString(col)
	}
	for i := range sc.Indexes {
		cksum += sc.Indexes[i].Cksum()
	}
	return cksum
}

func (ix *Index) Cksum() uint32 {
	cksum := uint32(ix.Mode)
	for _, col := range ix.Columns {
		cksum += hash.HashString(col)
	}
	return cksum
}
