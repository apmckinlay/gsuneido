// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"strings"

	"github.com/apmckinlay/gsuneido/database/db19/ixspec"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Schema struct {
	Table   string
	Columns []ColumnSchema
	Indexes []IndexSchema
	//TODO foreign key target stuff
	// mutable is used to know whether to persist
	mutable bool
}

//go:generate genny -in ../../../genny/hamt/hamt.go -out schemahamt.go -pkg meta gen "Item=*Schema KeyType=string"
//go:generate genny -in ../../../genny/hamt/hamt2.go -out schemahamt2.go -pkg meta gen "Item=*Schema KeyType=string"

func SchemaKey(ti *Schema) string {
	return ti.Table
}

func SchemaHash(key string) uint32 {
	return hash.HashString(key)
}

type ColumnSchema struct {
	Name  string
	Field int
}

type IndexSchema struct {
	Fields []int
	Ixspec ixspec.T
	// Mode is 'k' for key, 'i' for index, 'u' for unique index
	Mode     int
	Fktable  string
	Fkmode   int
	Fkfields []int
}

// fkmode bits
const (
	BLOCK           = 0
	CASCADE_UPDATES = 1
	CASCADE_DELETES = 2
	CASCADE         = CASCADE_UPDATES | CASCADE_DELETES
)

func (ts *Schema) storSize() int {
	size := 2 + len(ts.Table)
	size += 2
	for _, col := range ts.Columns {
		size += 2 + 2 + len(col.Name)
	}
	size++
	for i := range ts.Indexes {
		idx := &ts.Indexes[i]
		size += 1 + 2 + 2*len(idx.Fields) +
			2 + len(idx.Fktable) + 1 + 1 + 2*len(idx.Fkfields)
	}
	return size
}

func (ts *Schema) Write(w *stor.Writer) {
	w.PutStr(ts.Table)
	w.Put2(len(ts.Columns))
	for i := range ts.Columns {
		col := &ts.Columns[i]
		w.Put2(col.Field).PutStr(col.Name)
	}
	w.Put1(len(ts.Indexes))
	for i := range ts.Indexes {
		idx := &ts.Indexes[i]
		w.Put1(idx.Mode).Put2Ints(idx.Fields)
		w.PutStr(idx.Fktable).Put1(idx.Fkmode).Put1Ints(idx.Fkfields)
	}
}

func ReadSchema(_ *stor.Stor, r *stor.Reader) *Schema {
	ts := Schema{}
	ts.Table = r.GetStr()
	n := r.Get2()
	ts.Columns = make([]ColumnSchema, n)
	for i := 0; i < n; i++ {
		ts.Columns[i] = ColumnSchema{Field: r.Get2(), Name: r.GetStr()}
	}
	n = r.Get1()
	ts.Indexes = make([]IndexSchema, n)
	for i := 0; i < n; i++ {
		ts.Indexes[i] = IndexSchema{
			Mode:     r.Get1(),
			Fields:   r.Get2Ints(),
			Fktable:  r.GetStr(),
			Fkmode:   r.Get1(),
			Fkfields: r.Get1Ints(),
		}
	}
	ts.Ixspecs()
	return &ts
}

func (ts *Schema) Ixspecs() {
	key := ts.firstShortestKey()
	for i := range ts.Indexes {
		ix := &ts.Indexes[i]
		ix.Ixspec.Cols = ix.Fields
		switch ts.Indexes[i].Mode {
		case 'u':
			ix.Ixspec.Cols2 = key
		case 'i':
			ix.Ixspec.Cols = append(ix.Fields, key...)
		}
	}
}

func (ts *Schema) firstShortestKey() []int {
	var key []int
	for i := range ts.Indexes {
		ix := &ts.Indexes[i]
		if ix.usableKey() &&
			(key == nil || len(ix.Fields) < len(key)) {
			key = ix.Fields
		}
	}
	return key
}

func (ix *IndexSchema) usableKey() bool {
	return ix.Mode == 'k' && len(ix.Fields) > 0 && !hasSpecial(ix.Fields)
}

func hasSpecial(fields []int) bool {
	for _, f := range fields {
		if f < 0 {
			return true
		}
	}
	return false
}

//-------------------------------------------------------------------

func (ts *Schema) String() string {
	var sb strings.Builder
	var cb str.CommaBuilder
	for i := range ts.Columns {
		col := ts.Columns[i]
		if col.Field == -1 {
			cb.Add("-")
		} else {
			cb.Add(ts.Columns[i].Name)
		}
	}
	sb.WriteString("(")
	sb.WriteString(cb.String())
	sb.WriteString(")")
	for i := range ts.Indexes {
		sb.WriteString(" ")
		sb.WriteString(ts.Indexes[i].string(ts.Columns))
	}
	return sb.String()
}

func (ix *IndexSchema) string(cols []ColumnSchema) string {
	var cb str.CommaBuilder
	for _, c := range ix.Fields {
		if c < 0 {
			cb.Add(cols[-c-2].Name + "_lower!")
		} else {
			cb.Add(cols[c].Name)
		}
	}
	s := map[int]string{'k': "key", 'i': "index", 'u': "index unique"}[ix.Mode]
	s += "(" + cb.String() + ")"
	if ix.Fktable != "" {
		s += " in " + ix.Fktable
		if ix.Fkfields != nil {
			sep := "("
			for _, f := range ix.Fkfields {
				s += sep + cols[f].Name
				sep = ","
			}
		}
		if ix.Fkmode&CASCADE != 0 {
			s += " cascade"
			if ix.Fkmode == CASCADE_UPDATES {
				s += " update"
			}
		}
	}
	return s
}
