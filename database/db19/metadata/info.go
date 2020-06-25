// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package metadata

import (
	"sort"

	"github.com/apmckinlay/gsuneido/database/db19/btree"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/verify"
)

//go:generate genny -in ../../../genny/flathash/flathash.go -out tables.go -pkg metadata gen "Key=string Item=TableInfo"

type TableInfo struct {
	Table   string
	Nrows   int
	Size    uint64
	Indexes []*btree.Overlay
	// Schema is separate because it changes much less often
	Schema *TableSchema
	// mutable is used to know whether to persist
	mutable bool
}

func (*TableInfoHtbl) hash(key string) uint32 {
	return hash.HashString(key)
}

func (*TableInfoHtbl) keyOf(ti *TableInfo) string {
	return ti.Table
}

//-------------------------------------------------------------------

// // Merge combines two TableInfoHtbl into a new one.
// // t2 takes precedence.
// func (t *TableInfoHtbl) Merge(t2 *TableInfoHtbl) *TableInfoHtbl {
// 	// Important - bulk copy rather than inserting individually
// 	t3 := t.Dup()
// 	iter := t2.Iter()
// 	for ti2 := iter(); ti2 != nil; ti2 = iter() {
// 		ti := t.Get(ti2.Table)
// 		if ti == nil {
// 			t3.Put(ti2)
// 		} else {
// 			t3.Put(ti.Merge(ti2))
// 		}
// 	}
// 	return t3
// }

// // Merge combines two TableInfo into a new one.
// // ti2 takes precedence.
// func (ti *TableInfo) Merge(ti2 *TableInfo) *TableInfo {
// 	return &TableInfo{Table: ti2.Table,
// 		Nrows:   ti.Nrows + ti2.Nrows,
// 		Size:    ti.Size + ti2.Size,
// 		Indexes: append([]*btree.Overlay(nil), ti2.Indexes...),
// 		Schema:  ti2.Schema,
// 	}
// }

//-------------------------------------------------------------------

const blockSize = 4 * 1024
const itemsPerFinger = 16

// WriteInfo saves a TableInfoHtbl to external packed format in a stor.
// Tables are sorted by table number.
func (t *TableInfoHtbl) WriteInfo(st *stor.Stor) uint64 {
	return t.Write(st, (*TableInfo).WriteInfo)
}

func (t *TableInfoHtbl) Write(st *stor.Stor, write func(*TableInfo, *stor.Writer)) uint64 {
	w := st.Writer(blockSize)
	keys := t.List()
	sort.Strings(keys)
	w.Put2(t.nitems)
	nfingers := 1 + t.nitems/itemsPerFinger
	w2 := *w
	for i := 0; i < nfingers; i++ {
		w.Put3(0) // leave room
	}
	fingers := make([]int, 0, nfingers)
	for i, k := range keys {
		if i%16 == 0 {
			fingers = append(fingers, w.Pos())
		}
		write(t.Get(k), w)
	}
	verify.That(len(fingers) == nfingers)
	for _, f := range fingers {
		w2.Put3(f) // update with actual values
	}
	return w.Close()
}

// Write saves a TableInfo to external packed format in a stor
func (ti *TableInfo) WriteInfo(w *stor.Writer) {
	w.PutStr(ti.Table).
		Put4(ti.Nrows).
		Put5(ti.Size).
		Put1(len(ti.Indexes))
	for _, ii := range ti.Indexes {
		ii.Write(w)
	}
}

func ReadInfo(st *stor.Stor, off uint64) *TableInfoHtbl {
	r := st.Reader(off)
	nitems := r.Get2()
	nfingers := 1 + nitems/itemsPerFinger
	for i := 0; i < nfingers; i++ {
		r.Get3() // skip the fingers
	}
	t := NewTableInfoHtbl(nitems)
	for i := 0; i < nitems; i++ {
		t.Put(ReadTableInfo(r))
	}
	return t
}

func ReadTableInfo(r *stor.Reader) *TableInfo {
	var ti TableInfo
	ti.Table = r.GetStr()
	ti.Nrows = r.Get4()
	ti.Size = r.Get5()
	ni := r.Get1()
	ti.Indexes = make([]*btree.Overlay, ni)
	for i := 0; i < ni; i++ {
		ti.Indexes[i] = btree.ReadOverlay(r)
	}
	return &ti
}

//-------------------------------------------------------------------

type update struct {
	table string
	overlays
}
type overlays []*btree.Overlay

type btOver = *btree.Overlay

func (t *TableInfoHtbl) process(fn func(btOver) btOver) []update {
	var updates []update
	iter := t.Iter()
	for ti := iter(); ti != nil; ti = iter() {
		if ti.mutable {
			updated := make(overlays, len(ti.Indexes))
			for i, ov := range ti.Indexes {
				updated[i] = fn(ov)
			}
			updates = append(updates, update{table: ti.Table, overlays: updated})
		}
	}
	return updates
}

func (t *TableInfoHtbl) withUpdates(updates []update, fn func(btOver, btOver) btOver) *TableInfoHtbl {
	t2 := t.Dup()
	for _, up := range updates {
		ti := *t2.Get(up.table)                           // copy
		ti.Indexes = append(overlays(nil), ti.Indexes...) // copy
		for i, ov := range ti.Indexes {
			if up.overlays[i] != nil {
				ti.Indexes[i] = fn(ov, up.overlays[i])
			}
		}
		t2.Put(&ti)
	}
	return t2
}

//-------------------------------------------------------------------

func (t *TableInfoHtbl) Merge(tranNum int) []update {
	return t.process(func(ov btOver) btOver {
		return ov.Merge(tranNum)
	})
}
func (t *TableInfoHtbl) WithMerged(updates []update) *TableInfoHtbl {
	return t.withUpdates(updates, btOver.WithMerged)
}

//-------------------------------------------------------------------

func (t *TableInfoHtbl) SaveIndexes() []update {
	return t.process(btOver.Save)
}

func (t *TableInfoHtbl) WithSaved(updates []update) *TableInfoHtbl {
	return t.withUpdates(updates, btOver.WithSaved)
}

//-------------------------------------------------------------------

type InfoPacked struct {
	packed
}

type packed struct {
	off     uint64
	r       *stor.Reader
	fingers []finger
}

type finger struct {
	table string
	pos   int
}

func NewInfoPacked(st *stor.Stor, off uint64) *InfoPacked {
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
	return &InfoPacked{packed{off: off, r: r, fingers: fingers}}
}

func (p InfoPacked) Get(table string) *TableInfo {
	p.r.Pos(p.binarySearch(table))
	count := 0
	for {
		ti := ReadTableInfo(p.r)
		if ti.Table == table {
			return ti
		}
		count++
		if count > 20 {
			panic("linear search too long")
		}
	}
}

func (p packed) binarySearch(table string) int {
	i, j := 0, len(p.fingers)
	count := 0
	for i < j {
		h := int(uint(i+j) >> 1) // i ≤ h < j
		if table >= p.fingers[h].table {
			i = h + 1
		} else {
			j = h
		}
		count++
		if count > 20 {
			panic("binary search too long")
		}
	}
	// i is first one greater, so we want i-1
	return int(p.fingers[i-1].pos)
}

func (p packed) Offset() uint64 {
	return p.off
}

//-------------------------------------------------------------------

//TODO Merge a TableInfoHtbl and an InfoPacked to make new base
