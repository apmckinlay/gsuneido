// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"github.com/apmckinlay/gsuneido/database/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/hash"
)

type Schema struct {
	schema.Schema
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

func (ts *Schema) storSize() int {
	size := stor.LenStr(ts.Table) +
		stor.LenStrs(ts.Columns) + stor.LenStrs(ts.Derived) + 1
	for i := range ts.Indexes {
		idx := ts.Indexes[i]
		size += 1 + stor.Len2Ints(idx.Fields) +
			stor.LenStr(idx.Fktable) + 1 + stor.LenStrs(idx.Fkcolumns)
	}
	return size
}

func (ts *Schema) Write(w *stor.Writer) {
	w.PutStr(ts.Table)
	w.PutStrs(ts.Columns)
	w.PutStrs(ts.Derived)
	w.Put1(len(ts.Indexes))
	for _, ix := range ts.Indexes {
		w.Put1(ix.Mode).Put2Ints(ix.Fields)
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
			Fields:    r.Get2Ints(),
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
		ix.Ixspec.Fields = ix.Fields
		switch ts.Indexes[i].Mode {
		case 'u':
			ix.Ixspec.Fields2 = key
		case 'i':
			ix.Ixspec.Fields = append(ix.Fields, key...)
		}
	}
}

func (ts *Schema) firstShortestKey() []int {
	var key []int
	for i := range ts.Indexes {
		ix := &ts.Indexes[i]
		if usableKey(ix) &&
			(key == nil || len(ix.Fields) < len(key)) {
			key = ix.Fields
		}
	}
	return key
}

func usableKey(ix *schema.Index) bool {
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
