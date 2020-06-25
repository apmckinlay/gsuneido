// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package metadata

import "github.com/apmckinlay/gsuneido/database/db19/stor"

type TableSchema struct {
	Table   string
	Columns []ColumnSchema
	Indexes []IndexSchema
	//TODO foreign key target stuff
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

// WriteSchema saves the schema to external packed format in a stor.
// Tables are sorted by table number.
func (t *TableInfoHtbl) WriteSchema(st *stor.Stor) uint64 {
	return t.Write(st, (*TableInfo).WriteSchema)
}

func (ti *TableInfo) WriteSchema(w *stor.Writer) {
	ts := ti.Schema
	w.PutStr(ti.Table)
	w.Put2(len(ts.Columns))
	for _, col := range ts.Columns {
		w.Put2(col.Field).PutStr(col.Name)
	}
	w.Put1(len(ts.Indexes))
	for _, idx := range ts.Indexes {
		w.Put1(idx.Mode).PutInts(idx.Fields)
		w.PutStr(idx.Fktable).Put1(idx.Fkmode).PutInts(idx.Fkfields)
	}
}

func (t *TableInfoHtbl) ReadSchema(st *stor.Stor, off uint64) *TableInfoHtbl {
	r := st.Reader(off)
	nitems := r.Get2()
	if nitems != t.nitems {
		panic("info - schema size mismatch")
	}
	nfingers := 1 + nitems/itemsPerFinger
	for i := 0; i < nfingers; i++ {
		r.Get3() // skip the fingers
	}
	for i := 0; i < nitems; i++ {
		ts := ReadTableSchema(r)
		t.Get(ts.Table).Schema = ts
	}
	return t
}

func ReadTableSchema(r *stor.Reader) *TableSchema {
	ts := TableSchema{}
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

//-------------------------------------------------------------------

type SchemaPacked struct {
	packed
}

func NewSchemaPacked(st *stor.Stor, off uint64) *SchemaPacked {
	r := st.Reader(off)
	nitems := r.Get2()
	nfingers := 1 + nitems/itemsPerFinger
	fingers := make([]finger, nfingers)
	for i := 0; i < nfingers; i++ {
		fingers[i].pos = r.Get3()
	}
	for i := 0; i < nfingers; i++ {
		fingers[i].table = r.Pos(fingers[i].pos).GetStr()
	}
	return &SchemaPacked{packed{r: r, fingers: fingers}}
}

func (p SchemaPacked) Get(table string) *TableSchema {
	p.r.Pos(p.binarySearch(table))
	count := 0
	for {
		ti := ReadTableSchema(p.r)
		if ti.Table == table {
			return ti
		}
		count++
		if count > 20 {
			panic("linear search too long")
		}
	}
}
