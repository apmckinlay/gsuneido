// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"strings"

	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/strs"
)

// Note: views are stored with the name in Schema.Table prefixed by '='
// and the definition in Schema.Columns[0]

type Schema struct {
	schema.Schema
	// lastMod must be set to Meta.infoClock on new or modified items.
	// It is used for persist meta chaining/flattening.
	lastMod int
}

//go:generate genny -in ../../genny/hamt/hamt.go -out schemahamt.go -pkg meta gen "Item=*Schema KeyType=string"

func SchemaKey(ti *Schema) string {
	return ti.Table
}

func SchemaHash(key string) uint32 {
	return hash.HashString(key)
}

func (ht SchemaHamt) MustGet(key string) *Schema {
	it, ok := ht.Get(key)
	if !ok || it.isTomb() {
		panic("schema MustGet failed for " + key)
	}
	return it
}

func (ts *Schema) storSize() int {
	size := stor.LenStr(ts.Table) +
		stor.LenStrs(ts.Columns) + stor.LenStrs(ts.Derived) + 1
	for i := range ts.Indexes {
		idx := ts.Indexes[i]
		size += 1 + stor.LenStrs(idx.Columns) +
			stor.LenStr(idx.Fk.Table) + 1 + stor.LenStrs(idx.Fk.Columns)
	}
	return size
}

func (ts *Schema) Write(w *stor.Writer) {
	w.PutStr(ts.Table)
	w.PutStrs(ts.Columns)
	w.PutStrs(ts.Derived)
	w.Put1(len(ts.Indexes))
	for _, ix := range ts.Indexes {
		w.Put1(ix.Mode).PutStrs(ix.Columns)
		w.PutStr(ix.Fk.Table).Put1(ix.Fk.Mode).PutStrs(ix.Fk.Columns)
	}
}

func ReadSchema(_ *stor.Stor, r *stor.Reader) *Schema {
	ts := Schema{}
	ts.Table = r.GetStr()
	ts.Columns = r.GetStrs()
	ts.Derived = r.GetStrs()
	if n := r.Get1(); n > 0 {
		ts.Indexes = make([]schema.Index, n)
		for i := 0; i < n; i++ {
			ts.Indexes[i] = schema.Index{
				Mode:    r.Get1(),
				Columns: r.GetStrs(),
				Fk: schema.Fkey{
					Table:   r.GetStr(),
					Mode:    r.Get1(),
					Columns: r.GetStrs()},
			}
		}
		ts.Ixspecs(ts.Indexes)
	}
	return &ts
}

// Ixspecs sets up the ixspecs for a table's indexes.
func (ts *Schema) Ixspecs(idxs []schema.Index) {
	key := ts.firstShortestKey()
	for i := range idxs {
		ix := &idxs[i]
		switch ix.Mode {
		case 'u':
			cols := sset.Difference(key, ix.Columns)
			ix.Ixspec.Fields2 = ts.colsToFlds(cols)
			fallthrough
		case 'k':
			ix.Ixspec.Fields = ts.colsToFlds(ix.Columns)
		case 'i':
			cols := sset.Union(ix.Columns, key)
			ix.Ixspec.Fields = ts.colsToFlds(cols)
		default:
			panic("Ixspecs invalid mode")
		}
	}
}

func (ts *Schema) firstShortestKey() []string {
	var key []string
	for i := range ts.Indexes {
		ix := &ts.Indexes[i]
		if usableKey(ix) &&
			(key == nil || len(ix.Columns) < len(key)) {
			key = ix.Columns
		}
	}
	return key
}

func (ts *Schema) colsToFlds(cols []string) []int {
	flds := make([]int, len(cols))
	for i, col := range cols {
		c := strs.Index(ts.Columns, col)
		if strings.HasSuffix(col, "_lower!") {
			if c = strs.Index(ts.Columns, col[:len(col)-7]); c != -1 {
				c = -c - 2
			}
		}
		assert.That(c != -1)
		flds[i] = c
	}
	return flds
}

func usableKey(ix *schema.Index) bool {
	return ix.Mode == 'k' && len(ix.Columns) > 0 && !hasSpecial(ix.Columns)
}

func hasSpecial(cols []string) bool {
	for _, col := range cols {
		if strings.HasSuffix(col, "_lower!") {
			return true
		}
	}
	return false
}

func (m *Meta) newSchemaTomb(table string) *Schema {
	return &Schema{Schema: schema.Schema{Table: table}}
}

func (m *Meta) newSchemaView(name, def string) *Schema {
	return &Schema{Schema: schema.Schema{Table: "=" + name, Columns: []string{def}}}
}

func (ts *Schema) isTomb() bool {
	return ts.Columns == nil && ts.Indexes == nil
}

func (ts *Schema) isView() bool {
	return !ts.isTomb() && ts.Table[0] == '='
}

// isTable returns true if not a view and not a tombstone
func (ts *Schema) isTable() bool {
	return !ts.isTomb() && !ts.isView()
}
