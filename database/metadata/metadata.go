// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package metadata

import "sort"

//go:generate genny -in ../../genny/flathash/flathash.go -out tables.go -pkg metadata gen "Key=int Item=TableInfo"

type TableInfo struct {
	table   int
	nrows   int
	size    uint64
	indexes []IndexInfo
}

type IndexInfo struct {
	root       uint64
	treeLevels int
}

func (*TableInfoHtbl) hash(k int) uint32 {
	return uint32(k)
}

func (*TableInfoHtbl) keyOf(ti *TableInfo) int {
	return ti.table
}

//-------------------------------------------------------------------

// Merge combines two TableInfoHtbl into a new one.
// t2 takes precedence.
func (t *TableInfoHtbl) Merge(t2 *TableInfoHtbl) *TableInfoHtbl {
	// Important - bulk copy rather than inserting individually
	t3 := t.Dup()
	iter := t2.Iter()
	for ti2 := iter(); ti2 != nil; ti2 = iter() {
		ti := t.Get(ti2.table)
		if ti == nil {
			t3.Put(ti2)
		} else {
			t3.Put(ti.Merge(ti2))
		}
	}
	return t3
}

// Merge combines two TableInfo into a new one.
// ti2 takes precedence.
func (ti *TableInfo) Merge(ti2 *TableInfo) *TableInfo {
	return &TableInfo{table: ti2.table,
		nrows:   ti.nrows + ti2.nrows,
		size:    ti.size + ti2.size,
		indexes: append([]IndexInfo(nil), ti2.indexes...),
	}
}

//-------------------------------------------------------------------

// Write converts a TableInfoHtbl to external packed format in a byte slice.
// Tables are sorted by table number.
// TODO binary search fingers
func (t *TableInfoHtbl) Write() []byte {
	keys := t.List()
	sort.Ints(keys)
	buf := make(buffer, 0, 32*len(keys)) // guesstimate
	buf.put2(t.nitems)
	for _, k := range keys {
		t.Get(k).Write(&buf)
	}
	return buf
}

// Write converts a TableInfo to external packed format in a byte slice
func (ti *TableInfo) Write(b *buffer) {
	b.put3(ti.table).
		put4(ti.nrows).
		put5(ti.size).
		put1(len(ti.indexes))
	for _, ii := range ti.indexes {
		b.put5(ii.root).put1(ii.treeLevels)
	}
}

func ReadTablesInfo(buf []byte) *TableInfoHtbl {
	b := buffer(buf)
	nitems := b.get2()
	t := NewTableInfoHtbl(nitems)
	for len(b) > 0 {
		t.Put(ReadTableInfo(&b))
	}
	return t
}

func ReadTableInfo(b *buffer) *TableInfo {
	var ti TableInfo
	ti.table = b.get3()
	ti.nrows = b.get4()
	ti.size = b.get5()
	ni := b.get1()
	ti.indexes = make([]IndexInfo, ni)
	for i := 0; i < ni; i++ {
		ti.indexes[i].Read(b)
	}
	return &ti
}

func (ii *IndexInfo) Read(b *buffer) {
	ii.root = b.get5()
	ii.treeLevels = b.get1()
}
