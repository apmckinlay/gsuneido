// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/generic/hamt"
	"github.com/apmckinlay/gsuneido/util/hash"
	"golang.org/x/exp/slices"
)

type InfoHamt = hamt.Hamt[string, *Info]

type Info struct {
	Table   string
	Indexes []*index.Overlay
	Nrows   int
	Size    uint64
	// origNrows and origSize are used to determine the changes (delta)
	// made by a transaction. They are not used outside transactions.
	origNrows int
	origSize  uint64
	// lastMod must be set to Meta.infoClock on new or modified items.
	// It is used for persist meta chaining/flattening.
	lastMod int
	// created is used to avoid tombstones (and persisting them)
	// for temporary tables (e.g. from tests)
	created int
}

func (ti *Info) Key() string {
	return ti.Table
}

func (*Info) Hash(key string) uint32 {
	return hash.String(key)
}

func (ti *Info) LastMod() int {
	return ti.lastMod
}

func (ti *Info) SetLastMod(mod int) {
	ti.lastMod = mod
}

func (ti *Info) StorSize() int {
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

func (ti *Info) IsTomb() bool {
	return ti.Indexes == nil
}

//-------------------------------------------------------------------

type btOver = *index.Overlay
type MergeResult = index.MergeResult

type MergeUpdate struct {
	table   string
	results []MergeResult // per index
	nmerged int
}

// Merge collects the updates which are then applied by Apply.
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

func (mu MergeUpdate) Table() string {
	return mu.table
}

func (mu MergeUpdate) Apply(ov *index.Overlay, i int) *index.Overlay {
	return ov.WithMerged(mu.results[i], mu.nmerged)
}

//-------------------------------------------------------------------

type SaveResult = index.SaveResult

type PersistUpdate struct {
	table   string
	results []SaveResult // per index
}

// Persist is called by state.Persist to write the state to the database.
// It collects the new btree roots which are then applied by Apply.
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

func (pu PersistUpdate) Table() string {
	return pu.table
}

func (pu PersistUpdate) Apply(ov *index.Overlay, i int) *index.Overlay {
	return ov.WithSaved(pu.results[i])
}

func (ti *Info) Cksum() uint32 {
	cksum := hash.HashString(ti.Table) + uint32(ti.Nrows) + uint32(ti.Size)
	for _, ov := range ti.Indexes {
		cksum += ov.Cksum()
	}
	return cksum
}

//-------------------------------------------------------------------

type applyable interface {
	Table() string
	Apply(*index.Overlay, int) *index.Overlay
}

// Apply applies the updates collected by Merge or Persist
// It is called by state.go Database.Merge/Persist, inside UpdateState.
func Apply[U applyable](m *Meta, updates []U) {
	info := m.info.Mutable()
	for _, up := range updates {
		ti := *info.MustGet(up.Table()) // copy
		ti.Indexes = slices.Clone(ti.Indexes)
		for i, ov := range ti.Indexes {
			ti.Indexes[i] = up.Apply(ov, i)
		}
		ti.lastMod = m.info.Clock
		info.Put(&ti)
	}
	m.info.Hamt = info.Freeze()
}
