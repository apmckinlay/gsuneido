// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"github.com/apmckinlay/gsuneido/database/db19/btree"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/hash"
)

//go:generate genny -in ../../../genny/flathash/flathash.go -out infohtbl.go -pkg meta gen "Key=string Item=Info"

func (*InfoHtbl) hash(key string) uint32 {
	return hash.HashString(key)
}

func (*InfoHtbl) keyOf(ti *Info) string {
	return ti.Table
}

type Info struct {
	Table   string
	Nrows   int
	Size    uint64
	Indexes []*btree.Overlay
	// mutable is used to know whether to persist
	mutable bool
}

//-------------------------------------------------------------------

func (ti *Info) Write(w *stor.Writer) {
	w.PutStr(ti.Table).
		Put4(ti.Nrows).
		Put5(ti.Size).
		Put1(len(ti.Indexes))
	for _, ii := range ti.Indexes {
		ii.Write(w)
	}
}

func ReadInfo(st *stor.Stor, r *stor.Reader) *Info {
	var ti Info
	ti.Table = r.GetStr()
	ti.Nrows = r.Get4()
	ti.Size = r.Get5()
	ni := r.Get1()
	ti.Indexes = make([]*btree.Overlay, ni)
	for i := 0; i < ni; i++ {
		ti.Indexes[i] = btree.ReadOverlay(st, r)
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

func (t *InfoHtbl) process(fn func(btOver) btOver) []update {
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

func (t *InfoHtbl) withUpdates(updates []update, fn func(btOver, btOver) btOver) *InfoHtbl {
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

//TODO Merge an InfoHtbl and an InfoPacked to make new base
