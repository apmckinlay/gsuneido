// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"github.com/apmckinlay/gsuneido/db19/btree"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// Meta is the schema and info metadata
// difInfo is per transaction, overrides info
type Meta struct {
	schema  SchemaHamt
	info    InfoHamt
	difInfo InfoHamt
}

// Mutable returns a mutable copy of a Meta
func (m *Meta) Mutable() *Meta {
	assert.That(m.difInfo.IsNil())
	ov2 := *m // copy
	ov2.difInfo = InfoHamt{}.Mutable()
	return &ov2
}

func (m *Meta) GetRoInfo(table string) *Info {
	if ti, ok := m.difInfo.Get(table); ok {
		return ti
	}
	if ti, ok := m.info.Get(table); ok && !ti.isTomb() {
		return ti
	}
	return nil
}

func (m *Meta) GetRwInfo(table string, tranNum int) *Info {
	if pti, ok := m.difInfo.Get(table); ok {
		return pti // already have mutable
	}
	pti, ok := m.info.Get(table)
	if !ok || pti.isTomb() {
		return nil
	}
	ti := *pti // copy
	// start at 0 since these are deltas
	ti.Nrows = 0
	ti.Size = 0

	// set up index overlays
	ti.Indexes = append(ti.Indexes[:0:0], ti.Indexes...) // copy
	for i := range ti.Indexes {
		ti.Indexes[i] = ti.Indexes[i].Mutable(tranNum)
	}

	m.difInfo.Put(&ti)
	return &ti
}

func (m *Meta) GetRoSchema(table string) *Schema {
	ts, ok := m.schema.Get(table)
	if !ok || ts.isTomb() {
		return nil
	}
	return ts
}

// Put is used by Database.LoadedTable
func (m *Meta) Put(ts *Schema, ti *Info) *Meta {
	schema := m.schema.Mutable()
	schema.Put(ts)
	info := m.info.Mutable()
	info.Put(ti)
	ov2 := *m // copy
	ov2.schema = schema.Freeze()
	ov2.info = info.Freeze()
	return &ov2
}

func (m *Meta) DropTable(table string) *Meta {
	assert.That(m.difInfo.IsNil())
	if ts, ok := m.schema.Get(table); !ok || ts.isTomb() {
		return nil // nonexistent
	}
	//TODO delete without tombstone if not persisted
	schema := m.schema.Mutable()
	schema.Put(newSchemaTomb(table))
	info := m.info.Mutable()
	info.Put(newInfoTomb(table))
	ov2 := *m // copy
	ov2.schema = schema.Freeze()
	ov2.info = info.Freeze()
	return &ov2
}

func (m *Meta) ForEachSchema(fn func(*Schema)) {
	m.schema.ForEach(fn)
}

//-------------------------------------------------------------------

// LayeredOnto layers the mutable mbtree's from a transaction
// onto the latest/current state and returns a new state.
// Also, the nrows and size deltas are applied.
// Note: this does not merge the btrees, that is done later by merge.
// Nor does it save the changes to disk, that is done later by persist.
func (m *Meta) LayeredOnto(latest *Meta) *Meta {
	// start with a snapshot of the latest hash table because it may have more
	assert.That(latest.difInfo.IsNil())
	info := latest.info.Mutable()
	m.difInfo.ForEach(func(ti *Info) {
		lti, ok := info.Get(ti.Table)
		if !ok || lti.isTomb() {
			return
		}
		ti.Nrows += lti.Nrows
		ti.Size += lti.Size
		for i := range ti.Indexes {
			ti.Indexes[i].UpdateWith(lti.Indexes[i])
		}
		info.Put(ti)
	})
	result := *latest // copy
	result.info = info.Freeze()
	return &result
}

//-------------------------------------------------------------------

func (m *Meta) Write(store *stor.Stor) (offSchema, offInfo uint64) {
	assert.That(m.difInfo.IsNil())
	return m.schema.Write(store), m.info.Write(store)
}

func ReadMeta(store *stor.Stor, offSchema, offInfo uint64) *Meta {
	m := Meta{
		schema: SchemaHamt{}.Mutable().Read(store, offSchema).Freeze(),
		info:   InfoHamt{}.Mutable().Read(store, offInfo).Freeze(),
	}
	// set up ixspecs
	m.info.ForEach(func(ti *Info) {
		ts := m.schema.MustGet(ti.Table)
		for i := range ti.Indexes {
			ti.Indexes[i].SetIxspec(&ts.Indexes[i].Ixspec)
		}
	})
	return &m
}

//-------------------------------------------------------------------

// Merge is called by state.Merge
// to merge the mbtree's for tranNum into the fbtree's.
// It collect updates which are then applied by ApplyMerge
func (m *Meta) Merge(tranNum int) []update {
	return m.info.process(func(bto btOver) btOver {
		return bto.Merge(tranNum)
	})
}

// ApplyMerge applies the updates collected by meta.Merge
func (m *Meta) ApplyMerge(updates []update) {
	m.info = m.info.withUpdates(updates, btOver.WithMerged)
}

//-------------------------------------------------------------------

// Persist is called by state.Persist to write the state to the database.
// It collects the new fbtree roots which are then applied ApplyPersist.
func (m *Meta) Persist(flatten bool) []update {
	return m.info.process(func(ov *btree.Overlay) *btree.Overlay {
		return ov.Save(flatten)
	})
}

// ApplyPersist takes the new fbtree roots from meta.Persist
// and updates the state with them.
func (m *Meta) ApplyPersist(updates []update) {
	m.info = m.info.withUpdates(updates, btOver.WithSaved)
}
