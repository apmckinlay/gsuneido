// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/hash"
)

type Info struct {
	Table   string
	Nrows   int
	Size    uint64
	Indexes []*index.Overlay
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
	ti.Indexes = make([]*index.Overlay, ni)
	for i := 0; i < ni; i++ {
		ti.Indexes[i] = index.ReadOverlay(st, r)
	}
	return &ti
}

func (m *Meta) newInfoTomb(table string) *Info {
	return &Info{Table: table, lastmod: m.infoClock}
}

func (ti *Info) isTomb() bool {
	return len(ti.Indexes) == 0
}

//-------------------------------------------------------------------

type btOver = *index.Overlay
type MergeResult = index.MergeResult

type MergeUpdate struct {
	table   string
	nmerged int
	results []MergeResult // per index
}

// Merge collects the updates which are then applied by applyMerge.
// WARNING: must not modify meta.
func (m *Meta) Merge(table string, nmerge int) MergeUpdate {
	// fmt.Println("Merge", table, tns)
	ti := m.info.MustGet(table)
	results := make([]MergeResult, len(ti.Indexes))
	for j, ov := range ti.Indexes {
		results[j] = ov.Merge(nmerge)
	}
	return MergeUpdate{table: table, nmerged: nmerge, results: results}
}

// ApplyMerge applies the updates collected by Merge
func (m *Meta) ApplyMerge(updates []MergeUpdate) {
	t2 := m.info.Mutable()
	for _, up := range updates {
		// fmt.Println("applyMerge", up.table)
		ti := *t2.MustGet(up.table)                          // copy
		ti.Indexes = append(ti.Indexes[:0:0], ti.Indexes...) // copy
		for i, ov := range ti.Indexes {
			ti.Indexes[i] = ov.WithMerged(up.results[i], up.nmerged)
		}
		t2.Put(&ti)
	}
	m.info = t2.Freeze()
}

//-------------------------------------------------------------------

type SaveResult = index.SaveResult

type PersistUpdate struct {
	table   string
	results []SaveResult // per index
}


// Persist is called by state.Persist to write the state to the database.
// It collects the new fbtree roots which are then applied ApplyPersist.
// WARNING: must not modify meta.
func (m *Meta) Persist(exec func(func() PersistUpdate)) {
	m.info.ForEach(func(ti *Info) {
		if len(ti.Indexes) >= 1 && ti.Indexes[0].Modified() {
			exec(func() PersistUpdate {
				results := make([]SaveResult, len(ti.Indexes))
				for i, ov := range ti.Indexes {
					results[i] = ov.Save()
				}
				return PersistUpdate{table: ti.Table, results: results}
			})
		}
	})
}

// ApplyPersist takes the new fbtree roots from Persist
// and updates the state with them.
func (m *Meta) ApplyPersist(updates []PersistUpdate) {
	t2 := m.info.Mutable()
	for _, up := range updates {
		ti := *t2.MustGet(up.table)                          // copy
		ti.Indexes = append(ti.Indexes[:0:0], ti.Indexes...) // copy
		for i, ov := range ti.Indexes {
			if up.results[i] != nil {
				ti.Indexes[i] = ov.WithSaved(up.results[i])
			}
		}
		t2.Put(&ti)
	}
	m.info = t2.Freeze()
}
