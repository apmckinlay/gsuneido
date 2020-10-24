// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"github.com/apmckinlay/gsuneido/db19/btree"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hash"
)

type Info struct {
	Table   string
	Nrows   int
	Size    uint64
	Indexes []*btree.Overlay
	lastmod int
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
	assert.That(!ti.isTomb())
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

func (m *Meta) newInfoTomb(table string) *Info {
	return &Info{Table: table, Nrows: -1, lastmod: m.infoClock}
}

func (ti *Info) isTomb() bool {
	return ti.Nrows == -1
}

//-------------------------------------------------------------------

type btOver = *btree.Overlay
type Result = btree.Result

type Update struct {
	table   string
	results []Result // per index
}

// merge is used by meta.Merge.
// It collects the updates which are then applied by withUpdates.
func (t InfoHamt) merge(tn int, table string) Update {
	ti := t.MustGet(table)
	results := make([]Result, len(ti.Indexes))
	for j, ov := range ti.Indexes {
		results[j] = ov.Merge(tn)
	}
	return Update{table: table, results: results}
}

// process is used by meta.Persist.
// process collects the updates which are then applied by withUpdates.
func (t InfoHamt) process(fn func(btOver) Result) []Update {
	var updates []Update
	t.ForEach(func(ti *Info) {
		results := make([]Result, len(ti.Indexes))
		for i, ov := range ti.Indexes {
			r := fn(ov)
			if r == nil {
				assert.That(i == 0)
				return
			}
			results[i] = r
		}
		updates = append(updates, Update{table: ti.Table, results: results})
	})
	return updates
}

func (t InfoHamt) withUpdates(updates []Update, fn func(btOver, Result) btOver) InfoHamt {
	t2 := t.Mutable()
	for _, up := range updates {
		ti := *t2.MustGet(up.table)                          // copy
		ti.Indexes = append(ti.Indexes[:0:0], ti.Indexes...) // copy
		for i, ov := range ti.Indexes {
			if up.results[i] != nil {
				ti.Indexes[i] = fn(ov, up.results[i])
			}
		}
		t2.Put(&ti)
	}
	return t2.Freeze()
}
