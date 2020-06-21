// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package metadata

import "github.com/apmckinlay/gsuneido/database/stor"

type TableSchema struct {
	table   string
	columns []ColumnSchema
	indexes []IndexSchema
	//TODO foreign key target stuff
}

type ColumnSchema struct {
	name  string
	field int
}

type IndexSchema struct {
	fields []int
	// mode is 'k' for key, 'i' for index, 'u' for unique index
	mode     int
	fktable  string
	fkmode   int
	fkfields []int
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
	ts := ti.schema
	w.PutStr(ti.table)
	w.Put2(len(ts.columns))
	for _, col := range ts.columns {
		w.Put2(col.field).PutStr(col.name)
	}
	w.Put1(len(ts.indexes))
	for _, idx := range ts.indexes {
		w.Put1(idx.mode).PutInts(idx.fields)
		w.PutStr(idx.fktable).Put1(idx.fkmode).PutInts(idx.fkfields)
	}
}

func (t *TableInfoHtbl) ReadSchema(st *stor.Stor, off uint64) {
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
		t.Get(ts.table).schema = ts
	}
}

func ReadTableSchema(r *stor.Reader) *TableSchema {
	ts := TableSchema{}
	ts.table = r.GetStr()
	n := r.Get2()
	ts.columns = make([]ColumnSchema, n)
	for i := 0; i < n; i++ {
		ts.columns[i] = ColumnSchema{field: r.Get2(), name: r.GetStr()}
	}
	n = r.Get1()
	ts.indexes = make([]IndexSchema, n)
	for i := 0; i < n; i++ {
		ts.indexes[i] = IndexSchema{
			mode:     r.Get1(),
			fields:   r.GetInts(),
			fktable:  r.GetStr(),
			fkmode:   r.Get1(),
			fkfields: r.GetInts(),
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
		if ti.table == table {
			return ti
		}
		count++
		if count > 20 {
			panic("linear search too long")
		}
	}
}
