// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"strings"

	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Schema struct {
	schema.Schema
	//TODO foreign key target stuff
	// mutable is used to know whether to persist
	mutable bool
}

//go:generate genny -in ../../genny/hamt/hamt.go -out schemahamt.go -pkg meta gen "Item=*Schema KeyType=string"
//go:generate genny -in ../../genny/hamt/hamt2.go -out schemahamt2.go -pkg meta gen "Item=*Schema KeyType=string"

func SchemaKey(ti *Schema) string {
	return ti.Table
}

func SchemaHash(key string) uint32 {
	return hash.HashString(key)
}

func (ts *Schema) storSize() int {
	size := stor.LenStr(ts.Table) +
		stor.LenStrs(ts.Columns) + stor.LenStrs(ts.Derived) + 1
	for i := range ts.Indexes {
		idx := ts.Indexes[i]
		size += 1 + stor.LenStrs(idx.Columns) +
			stor.LenStr(idx.Fktable) + 1 + stor.LenStrs(idx.Fkcolumns)
	}
	return size
}

func (ts *Schema) Write(w *stor.Writer) {
	assert.That(!ts.isTomb())
	w.PutStr(ts.Table)
	w.PutStrs(ts.Columns)
	w.PutStrs(ts.Derived)
	w.Put1(len(ts.Indexes))
	for _, ix := range ts.Indexes {
		w.Put1(ix.Mode).PutStrs(ix.Columns)
		w.PutStr(ix.Fktable).Put1(ix.Fkmode).PutStrs(ix.Fkcolumns)
	}
}

func ReadSchema(_ *stor.Stor, r *stor.Reader) *Schema {
	ts := Schema{}
	ts.Table = r.GetStr()
	ts.Columns = r.GetStrs()
	ts.Derived = r.GetStrs()
	n := r.Get1()
	ts.Indexes = make([]schema.Index, n)
	for i := 0; i < n; i++ {
		ts.Indexes[i] = schema.Index{
			Mode:      r.Get1(),
			Columns:   r.GetStrs(),
			Fktable:   r.GetStr(),
			Fkmode:    r.Get1(),
			Fkcolumns: r.GetStrs(),
		}
	}
	ts.Ixspecs()
	return &ts
}

func (ts *Schema) Ixspecs() {
	key := ts.firstShortestKey()
	for i := range ts.Indexes {
		ix := &ts.Indexes[i]
		ix.Ixspec.Fields = ts.colsToFlds(ix.Columns)
		switch ts.Indexes[i].Mode {
		case 'u':
			ix.Ixspec.Fields2 = key
		case 'i':
			ix.Ixspec.Fields = append(ix.Ixspec.Fields, key...)
		}
	}
}

func (ts *Schema) colsToFlds(cols []string) []int {
	flds := make([]int, len(cols))
	for i, col := range cols {
		c := str.List(ts.Columns).Index(col)
		if strings.HasSuffix(col, "_lower!") {
			if c = str.List(ts.Columns).Index(col[:len(col)-7]); c != -1 {
				c = -c - 2
			}
		}
		assert.That(c != -1)
		flds[i] = c
	}
	return flds
}

func (ts *Schema) firstShortestKey() []int {
	var key []string
	for i := range ts.Indexes {
		ix := &ts.Indexes[i]
		if usableKey(ix) &&
			(key == nil || len(ix.Columns) < len(key)) {
			key = ix.Columns
		}
	}
	return ts.colsToFlds(key)
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

func newSchemaTomb(table string) *Schema {
	return &Schema{Schema: schema.Schema{Table: table}}
}

func (ts *Schema) isTomb() bool {
	return ts.Indexes == nil
}
