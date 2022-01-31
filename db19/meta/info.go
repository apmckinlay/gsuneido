// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/hash"
)

type Info struct {
	Table string
	Nrows int
	Size  uint64
	// origNrows and origSize are used to determine the changes (delta)
	// made by a transaction. They are not used outside transactions.
	origNrows int
	origSize  uint64
	Indexes   []*index.Overlay
	// lastMod must be set to Meta.infoClock on new or modified items.
	// It is used for persist meta chaining/flattening.
	lastMod int
}

//go:generate genny -in ../../genny/hamt/hamt.go -out infohamt.go -pkg meta gen "Item=*Info KeyType=string"

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
	if ni := r.Get1(); ni > 0 {
		ti.Indexes = make([]*index.Overlay, ni)
		for i := 0; i < ni; i++ {
			ti.Indexes[i] = index.ReadOverlay(st, r)
		}
	}
	return &ti
}

func (m *Meta) newInfoTomb(table string) *Info {
	return &Info{Table: table}
}

func (ti *Info) isTomb() bool {
	return ti.Indexes == nil
}

func (ht InfoHamt) MustGet(key string) *Info {
	it, ok := ht.Get(key)
	if !ok || it.isTomb() {
		panic("info MustGet failed for " + key)
	}
	return it
}

// GetCopy returns a copy of the Info for a table, or nil if not found
func (ht InfoHamt) GetCopy(table string) *Info {
	ti, ok := ht.Get(table)
	if !ok || ti.isTomb() {
		return nil
	}
	cp := *ti
	return &cp
}

//-------------------------------------------------------------------

type btOver = *index.Overlay
type MergeResult = index.MergeResult

type MergeUpdate struct {
	table   string
	nmerged int
	results []MergeResult // per index
}

// Merge collects the updates which are then applied by ApplyMerge.
// It is called by concur merger.
// WARNING: must not modify meta.
func (m *Meta) Merge(table string, nmerge int) MergeUpdate {
	ti := m.info.MustGet(table)
	results := make([]MergeResult, len(ti.Indexes))
	for i, ov := range ti.Indexes {
		results[i] = ov.Merge(nmerge)
	}
	return MergeUpdate{table: table, nmerged: nmerge, results: results}
}

// ApplyMerge applies the updates collected by Merge.
// It is called by state.go Database.Merge, inside UpdateState.
func (m *Meta) ApplyMerge(updates []MergeUpdate) {
	// NOTE: ApplyMerge and ApplyPersist have almost identical code.
	// Any changes probably apply to both.
	// TODO use generics to eliminate duplication
	info := m.info.Mutable()
	for _, up := range updates {
		ti := info.GetCopy(up.table)
		ti.Indexes = append(ti.Indexes[:0:0], ti.Indexes...) // copy
		for i, ov := range ti.Indexes {
			ti.Indexes[i] = ov.WithMerged(up.results[i], up.nmerged)
		}
		ti.lastMod = m.info.clock
		info.Put(ti)
	}
	m.info.InfoHamt = info.Freeze()
}

//-------------------------------------------------------------------

type SaveResult = index.SaveResult

type PersistUpdate struct {
	table   string
	results []SaveResult // per index
}

// Persist is called by state.Persist to write the state to the database.
// It collects the new btree roots which are then applied by ApplyPersist.
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

// ApplyPersist takes the new btree roots from Persist
// and updates the state with them.
func (m *Meta) ApplyPersist(updates []PersistUpdate) {
	// NOTE: ApplyMerge and ApplyPersist have almost identical code.
	// Any changes probably apply to both.
	// TODO use generics to eliminate duplication
	info := m.info.Mutable()
	for _, up := range updates {
		ti := info.GetCopy(up.table)
		ti.Indexes = append(ti.Indexes[:0:0], ti.Indexes...) // copy
		for i, ov := range ti.Indexes {
			if up.results[i] != nil {
				ti.Indexes[i] = ov.WithSaved(up.results[i])
			}
		}
		ti.lastMod = m.info.clock
		info.Put(ti)
	}
	m.info.InfoHamt = info.Freeze()
}

func (ti *Info) Cksum() uint32 {
	cksum := hash.HashString(ti.Table) + uint32(ti.Nrows) + uint32(ti.Size)
	for _, ov := range ti.Indexes {
		cksum += ov.Cksum()
	}
	return cksum
}
