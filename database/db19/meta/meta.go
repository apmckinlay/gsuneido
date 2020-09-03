// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// Meta is the layered info and schema metadata
type Meta struct {
	rwInfo      InfoHamt
	roInfo      InfoHamt
	rwSchema    SchemaHamt
	roSchema    SchemaHamt
	baseInfo    *InfoPacked
	baseSchema  *SchemaPacked
}

func NewMeta(baseSchema *SchemaPacked, baseInfo *InfoPacked,
	roSchema SchemaHamt, roInfo InfoHamt) *Meta {
	return &Meta{
		baseSchema:  baseSchema,
		baseInfo:    baseInfo,
		roInfo:      roInfo,
		rwSchema:    SchemaHamt{},
		roSchema:    roSchema,
		rwInfo:      InfoHamt{},
	}
}

func CreateMeta(store *stor.Stor) *Meta {
	return &Meta{
		baseInfo:   &InfoPacked{stor: store},
		baseSchema: &SchemaPacked{stor: store},
		roInfo: InfoHamt{},
		roSchema: SchemaHamt{},
	}
}

// Mutable returns a mutable copy of a Meta
func (m *Meta) Mutable() *Meta {
	assert.That(m.rwInfo.IsNil())
	assert.That(m.rwSchema.IsNil())
	ov2 := *m // copy
	ov2.rwInfo = InfoHamt{}.Mutable()
	ov2.rwSchema = SchemaHamt{}.Mutable()
	return &ov2
}

func (m *Meta) GetRoInfo(table string) *Info {
	if ti, ok := m.rwInfo.Get(table); ok {
		return ti
	}
	if ti, ok := m.roInfo.Get(table); ok {
		return ti
	}
	if ti, ok := m.baseInfo.Get(table); ok {
		if !m.rwInfo.IsNil() {
			m.rwInfo.Put(ti) // cache in memory
		}
		return ti
	}
	return nil
}

func (m *Meta) GetRwInfo(table string, tranNum int) *Info {
	if pti, ok := m.rwInfo.Get(table); ok {
		return pti // already have mutable
	}
	var ti Info
	if pti, ok := m.roInfo.Get(table); ok {
		ti = *pti // copy
	} else if pti, ok := m.baseInfo.Get(table); ok {
		ti = *pti // copy
	} else {
		return nil
	}
	// start at 0 since these are deltas
	ti.Nrows = 0
	ti.Size = 0
	ti.mutable = true

	// set up index overlays and ixspecs
	ti.Indexes = append(ti.Indexes[:0:0], ti.Indexes...) // copy
	for i := range ti.Indexes {
		ti.Indexes[i] = ti.Indexes[i].Mutable(tranNum)
		if ti.Indexes[i].GetIxspec() == nil {
			ts := m.GetRoSchema(table)
			is := &ts.Indexes[i].Ixspec
			ti.Indexes[i].SetIxspec(is)
		}
	}

	m.rwInfo.Put(&ti) // cache in memory
	return &ti
}

func (m *Meta) GetRoSchema(table string) *Schema {
	if ts, ok := m.rwSchema.Get(table); ok {
		return ts
	}
	if ts, ok := m.roSchema.Get(table); ok {
		return ts
	}
	if ts, ok := m.baseSchema.Get(table); ok {
		return ts
	}
	return nil
}

func (m *Meta) GetRwSchema(table string) *Schema {
	if ts, ok := m.rwSchema.Get(table); ok {
		return ts
	}
	var ts Schema
	if pts, ok := m.roSchema.Get(table); ok {
		ts = *pts // copy
	} else if pts, ok := m.baseSchema.Get(table); ok {
		ts = *pts // copy
	} else {
		return nil
	}
	ts.Columns = append(ts.Columns[:0:0], ts.Columns...) // copy
	ts.Indexes = append(ts.Indexes[:0:0], ts.Indexes...) // copy
	m.rwSchema.Put(&ts)
	return &ts
}

func (m *Meta) Add(ts *Schema, ti *Info) *Meta {
	roSchema := m.roSchema.Mutable()
	roSchema.Put(ts)
	roInfo := m.roInfo.Mutable()
	roInfo.Put(ti)
	ov2 := *m // copy
	ov2.roSchema = roSchema.Freeze()
	ov2.roInfo = roInfo.Freeze()
	return &ov2
}

//-------------------------------------------------------------------

// LayeredOnto layers the mutable mbtree's from a transaction
// onto the latest/current state and returns a new state.
// Also, the nrows and size deltas are applied.
// Note: this does not merge the btrees, that is done later by merge.
// Nor does it save the changes to disk, that is done later by persist.
func (m *Meta) LayeredOnto(latest *Meta) *Meta {
	// start with a copy of the latest hash table because it may have more
	assert.That(latest.rwInfo.IsNil())
	roInfo2 := latest.roInfo.Mutable()
	m.rwInfo.ForEach(func(ti *Info) {
		if ti.mutable {
			if lti, ok := roInfo2.Get(ti.Table); ok {
				ti.Nrows += lti.Nrows
				ti.Size += lti.Size
				for i := range ti.Indexes {
					ti.Indexes[i].UpdateWith(lti.Indexes[i])
				}
			} else {
				// latest doesn't have this table, i.e. first update
				for i := range ti.Indexes {
					ti.Indexes[i].Freeze()
				}
			}
			roInfo2.Put(ti)
		}
	})
	//TODO handle rwSchema
	result := *latest // copy
	result.roInfo = roInfo2.Freeze()
	return &result
}

//-------------------------------------------------------------------

const Noffsets = 4

type offsets = [Noffsets]uint64

func (m *Meta) Write(st *stor.Stor) offsets {
	assert.That(m.rwInfo.IsNil())
	assert.That(m.rwSchema.IsNil())
	offs := offsets{
		m.baseSchema.Offset(),
		m.baseInfo.Offset(),
		m.roSchema.Write(st),
		m.roInfo.Write(st),
	}
	if m.baseSchema.Offset() == 0 {
		offs[0] = offs[2]
		offs[2] = 0
	}
	if m.baseInfo.Offset() == 0 {
		offs[1] = offs[3]
		offs[3] = 0
	}
	return offs
}

func ReadOverlay(st *stor.Stor, offs offsets) *Meta {
	m := Meta{
		baseSchema: NewSchemaPacked(st, offs[0]),
		baseInfo:   NewInfoPacked(st, offs[1]),
		roSchema:   ReadSchemaHamt(st, offs[2]),
		roInfo:     ReadInfoHamt(st, offs[3]),
	}
	return &m
}

//-------------------------------------------------------------------

// Merge is called by state.Merge to collect updates
// which are then applied by ApplyMerge
func (m *Meta) Merge(tranNum int) []update {
	return m.roInfo.process(func(bto btOver) btOver {
		return bto.Merge(tranNum)
	})
}

func (m *Meta) ApplyMerge(updates []update) {
	m.roInfo = m.roInfo.withUpdates(updates, btOver.WithMerged)
}

//-------------------------------------------------------------------

//TODO schema

func (m *Meta) Persist() []update {
	return m.roInfo.process(btOver.Save)
}

func (m *Meta) ApplyPersist(updates []update) {
	m.roInfo = m.roInfo.withUpdates(updates, btOver.WithSaved)
}
