// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/hash"
)

//go:generate genny -in ../../../genny/flathash/flathash.go -out schemahtbl.go -pkg meta gen "Key=string Item=Schema"

func (*SchemaHtbl) hash(key string) uint32 {
	return hash.HashString(key)
}

func (*SchemaHtbl) keyOf(ti *Schema) string {
	return ti.Table
}

type Schema struct {
	Table   string
	Columns []ColumnSchema
	Indexes []IndexSchema
	//TODO foreign key target stuff
	// mutable is used to know whether to persist
	mutable bool
}

type ColumnSchema struct {
	Name  string
	Field int
}

type IndexSchema struct {
	Fields []int
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

func (sc *Schema) Write(w *stor.Writer) {
	w.PutStr(sc.Table)
	w.Put2(len(sc.Columns))
	for _, col := range sc.Columns {
		w.Put2(col.Field).PutStr(col.Name)
	}
	w.Put1(len(sc.Indexes))
	for _, idx := range sc.Indexes {
		w.Put1(idx.Mode).PutInts(idx.Fields)
		w.PutStr(idx.Fktable).Put1(idx.Fkmode).PutInts(idx.Fkfields)
	}
}

func ReadSchema(r *stor.Reader) *Schema {
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
			Fields:   r.GetInts(),
			Fktable:  r.GetStr(),
			Fkmode:   r.Get1(),
			Fkfields: r.GetInts(),
		}
	}
	return &ts
}
