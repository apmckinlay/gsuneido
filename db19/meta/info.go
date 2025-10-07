// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/hamt"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/hash"
)

type InfoHamt = hamt.Hamt[string, *Info]

type Info struct {
	Table      string
	Indexes    []*index.Overlay
	Nrows      int
	Size       int64
	BtreeNrows int
	BtreeSize  int64
	// Deltas tracks the count & size changes per layer
	// parallel to the Indexes Overlay layers.
	// Deltas + BtreeNrows/Size should equal Nrows/Size
	Deltas []Delta
	// lastMod must be set to Meta.infoClock on new or modified items.
	// It is used for persist meta chaining/flattening.
	lastMod int
	// created is used to avoid tombstones (and persisting them)
	// for temporary tables (e.g. from tests)
	created int
}

type Delta struct {
	Nrows int
	Size  int64
}

func NewInfo(table string, indexes []*index.Overlay, nrows int, size int64) *Info {
	return &Info{Table: table, Indexes: indexes, Deltas: []Delta{{}},
		Nrows: nrows, Size: size, BtreeNrows: nrows, BtreeSize: size}
}

func (ti *Info) Key() string {
	return ti.Table
}

func (*Info) Hash(key string) uint64 {
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
		Put4(ti.BtreeNrows).
		Put5(ti.BtreeSize).
		Put1(len(ti.Indexes))
	for i := range ti.Indexes {
		ti.Indexes[i].Write(w)
	}
}

func ReadInfo(st *stor.Stor, r *stor.Reader) *Info {
	table := r.GetStr()
	nrows := r.Get4()
	size := r.Get5()
	var indexes []*index.Overlay
	if ni := r.Get1(); ni > 0 {
		indexes = make([]*index.Overlay, ni)
		for i := range ni {
			indexes[i] = index.ReadOverlay(st, r)
		}
	}
	return NewInfo(table, indexes, nrows, size)
}

func (m *Meta) newInfoTomb(table string) *Info {
	return &Info{Table: table}
}

func (ti *Info) IsTomb() bool {
	return ti.Indexes == nil
}

func (ti *Info) Cksum() uint32 {
	cksum := hash.HashString(ti.Table) + uint32(ti.BtreeNrows) + uint32(ti.BtreeSize)
	for _, ov := range ti.Indexes {
		cksum += ov.Cksum()
	}
	return cksum
}

func (ti *Info) Check() {
	for i := range ti.Indexes {
		assert.That(ti.Indexes[i].Nlayers() == len(ti.Deltas))
	}
	sum := Delta{Nrows: ti.BtreeNrows, Size: ti.BtreeSize}
	for _, d := range ti.Deltas {
		sum.Nrows += d.Nrows
		sum.Size += d.Size
	}
	assert.That(sum.Nrows == ti.Nrows)
	assert.That(sum.Size == ti.Size)
}

//-------------------------------------------------------------------

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

func (mu MergeUpdate) Apply1(ti *Info) {
	var sum Delta
	for _, d := range ti.Deltas[:1+mu.nmerged] {
		sum.Nrows += d.Nrows
		sum.Size += d.Size
	}
	deltas := make([]Delta, len(ti.Deltas)-mu.nmerged)
	deltas[0] = sum
	copy(deltas[1:], ti.Deltas[1+mu.nmerged:])
	ti.Deltas = deltas
}

func (mu MergeUpdate) Apply2(ov *index.Overlay, i int) *index.Overlay {
	return ov.WithMerged(mu.results[i], mu.nmerged)
}

//-------------------------------------------------------------------

type PersistUpdate struct {
	table   string
	results []iface.Btree // per index
}

// Persist is called by state.Persist to write the index updates.
// It collects the new btree roots which are then applied by Apply.
// WARNING: must not modify meta.
func (m *Meta) Persist(exec func(func() PersistUpdate)) {
	for ti := range m.info.All() {
		if len(ti.Indexes) >= 1 && ti.Indexes[0].Modified() {
			exec(func() PersistUpdate {
				results := make([]iface.Btree, len(ti.Indexes))
				for i, ov := range ti.Indexes {
					results[i] = ov.Save()
				}
				return PersistUpdate{table: ti.Table, results: results}
			})
		}
	}
}

func (pu PersistUpdate) Table() string {
	return pu.table
}

func (mu PersistUpdate) Apply1(ti *Info) {
	ti.BtreeNrows += ti.Deltas[0].Nrows
	assert.That(ti.BtreeNrows >= 0)
	ti.BtreeSize += ti.Deltas[0].Size
	assert.That(ti.BtreeSize >= 0)
	ti.Deltas = slc.Clone(ti.Deltas)
	ti.Deltas[0] = Delta{}
}

func (pu PersistUpdate) Apply2(ov *index.Overlay, i int) *index.Overlay {
	return ov.WithSaved(pu.results[i])
}

//-------------------------------------------------------------------

type applyable interface {
	Table() string
	Apply1(ti *Info)
	Apply2(*index.Overlay, int) *index.Overlay
}

// Apply applies the updates collected by Merge or Persist
// It is called by state.go Database.Merge/Persist, inside UpdateState.
func Apply[U applyable](m *Meta, updates []U) {
	info := m.info.Mutable()
	for _, up := range updates {
		ti := *info.MustGet(up.Table()) // shallow copy
		up.Apply1(&ti)
		ti.Indexes = slc.Clone(ti.Indexes)
		for i, ov := range ti.Indexes {
			ti.Indexes[i] = up.Apply2(ov, i)
		}
		ti.lastMod = m.info.Clock
		info.Put(&ti)
	}
	m.info.Hamt = info.Freeze()
}
