// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"github.com/apmckinlay/gsuneido/db19/btree"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/hash"
)

type Info struct {
	Table   string
	Nrows   int
	Size    uint64
	Indexes []*btree.Overlay
	// mutable is used to know whether to persist
	mutable bool
}

//go:generate genny -in ../../genny/hamt/hamt.go -out infohamt.go -pkg meta gen "Item=*Info KeyType=string"
//go:generate genny -in ../../genny/hamt/hamt2.go -out infohamt2.go -pkg meta gen "Item=*Info KeyType=string"

func InfoKey(ti *Info) string {
	return ti.Table
}

func InfoHash(key string) uint32 {
	return hash.HashString(key)
}

func (ti *Info) storSize() int {
	size := 2 + len(ti.Table) + 4 + 5 + 1
	for i := range ti.Indexes {
		size += ti.Indexes[i].StorSize()
	}
	return size
}

func (ti *Info) Write(w *stor.Writer) {
	w.PutStr(ti.Table).
		Put4(ti.Nrows).
		Put5(ti.Size).
		Put1(len(ti.Indexes))
	for i := range ti.Indexes {
		ti.Indexes[i].Write(w)
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

// process is used by overlay.Merge and overlay.Persist.
// process collects the updates which are then applied by withUpdates.
func (t InfoHamt) process(fn func(btOver) btOver) []update {
	var updates []update
	t.ForEach(func(ti *Info) {
		if ti.mutable {
			updated := make(overlays, len(ti.Indexes))
			for i, ov := range ti.Indexes {
				updated[i] = fn(ov)
			}
			updates = append(updates, update{table: ti.Table, overlays: updated})
		}
	})
	return updates
}

func (t InfoHamt) withUpdates(updates []update, fn func(btOver, btOver) btOver) InfoHamt {
	t2 := t.Mutable()
	for _, up := range updates {
		ti := *t2.MustGet(up.table)                       // copy
		ti.Indexes = append(overlays(nil), ti.Indexes...) // copy
		for i, ov := range ti.Indexes {
			if up.overlays[i] != nil {
				ti.Indexes[i] = fn(ov, up.overlays[i])
			}
		}
		t2.Put(&ti)
	}
	return t2.Freeze()
}

//-------------------------------------------------------------------

//TODO merge an InfoHamt and an InfoPacked to make new base
