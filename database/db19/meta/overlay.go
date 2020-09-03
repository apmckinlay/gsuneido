// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// Overlay provides access to the layered meta info and schema.
type Overlay struct {
	rwInfo      InfoHamt
	roInfo      InfoHamt
	roInfoOff   uint64
	rwSchema    SchemaHamt
	roSchema    SchemaHamt
	roSchemaOff uint64
	baseInfo    *InfoPacked
	baseSchema  *SchemaPacked
}

func NewOverlay(baseSchema *SchemaPacked, baseInfo *InfoPacked,
	roSchema SchemaHamt, roSchemaOff uint64,
	roInfo InfoHamt, roInfoOff uint64) *Overlay {
	return &Overlay{
		baseSchema:  baseSchema,
		baseInfo:    baseInfo,
		roInfo:      roInfo,
		roInfoOff:   roInfoOff,
		rwSchema:    SchemaHamt{},
		roSchema:    roSchema,
		roSchemaOff: roSchemaOff,
		rwInfo:      InfoHamt{},
	}
}

func CreateOverlay(store *stor.Stor) *Overlay {
	return &Overlay{
		baseInfo:   &InfoPacked{stor: store},
		baseSchema: &SchemaPacked{stor: store},
		roInfo: InfoHamt{},
		roSchema: SchemaHamt{},
	}
}

// NewOverlay returns a new Overlay based on an existing one
func (ov *Overlay) NewOverlay() *Overlay {
	assert.That(ov.rwInfo.IsNil())
	assert.That(ov.rwSchema.IsNil())
	ov2 := *ov // copy
	ov2.rwInfo = InfoHamt{}.Mutable()
	ov2.rwSchema = SchemaHamt{}.Mutable()
	return &ov2
}

func (ov *Overlay) GetRoInfo(table string) *Info {
	if ti, ok := ov.rwInfo.Get(table); ok {
		return ti
	}
	if ti, ok := ov.roInfo.Get(table); ok {
		return ti
	}
	if ti, ok := ov.baseInfo.Get(table); ok {
		if !ov.rwInfo.IsNil() {
			ov.rwInfo.Put(ti) // cache in memory
		}
		return ti
	}
	return nil
}

func (ov *Overlay) GetRwInfo(table string, tranNum int) *Info {
	if pti, ok := ov.rwInfo.Get(table); ok {
		return pti // already have mutable
	}
	var ti Info
	if pti, ok := ov.roInfo.Get(table); ok {
		ti = *pti // copy
	} else if pti, ok := ov.baseInfo.Get(table); ok {
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
			ts := ov.GetRoSchema(table)
			is := &ts.Indexes[i].Ixspec
			ti.Indexes[i].SetIxspec(is)
		}
	}

	ov.rwInfo.Put(&ti) // cache in memory
	return &ti
}

func (ov *Overlay) GetRoSchema(table string) *Schema {
	if ts, ok := ov.rwSchema.Get(table); ok {
		return ts
	}
	if ts, ok := ov.roSchema.Get(table); ok {
		return ts
	}
	if ts, ok := ov.baseSchema.Get(table); ok {
		return ts
	}
	return nil
}

func (ov *Overlay) GetRwSchema(table string) *Schema {
	if ts, ok := ov.rwSchema.Get(table); ok {
		return ts
	}
	var ts Schema
	if pts, ok := ov.roSchema.Get(table); ok {
		ts = *pts // copy
	} else if pts, ok := ov.baseSchema.Get(table); ok {
		ts = *pts // copy
	} else {
		return nil
	}
	ts.Columns = append(ts.Columns[:0:0], ts.Columns...) // copy
	ts.Indexes = append(ts.Indexes[:0:0], ts.Indexes...) // copy
	ov.rwSchema.Put(&ts)
	return &ts
}

func (ov *Overlay) Add(ts *Schema, ti *Info) *Overlay {
	roSchema := ov.roSchema.Mutable()
	roSchema.Put(ts)
	roInfo := ov.roInfo.Mutable()
	roInfo.Put(ti)
	ov2 := *ov // copy
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
func (ov *Overlay) LayeredOnto(latest *Overlay) *Overlay {
	// start with a copy of the latest hash table because it may have more
	assert.That(latest.rwInfo.IsNil())
	roInfo2 := latest.roInfo.Mutable()
	ov.rwInfo.ForEach(func(ti *Info) {
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
	result := *latest
	result.roInfo = roInfo2.Freeze()
	return &result
}

//-------------------------------------------------------------------

const Noffsets = 4

type offsets = [Noffsets]uint64

func (ov *Overlay) Write(st *stor.Stor) offsets {
	assert.That(ov.rwInfo.IsNil())
	assert.That(ov.rwSchema.IsNil())
	return offsets{
		ov.baseSchema.Offset(),
		ov.baseInfo.Offset(),
		ov.roSchema.Write(st),
		ov.roInfo.Write(st),
	}
}

func ReadOverlay(st *stor.Stor, offs offsets) *Overlay {
	ov := Overlay{
		baseSchema: NewSchemaPacked(st, offs[0]),
		baseInfo:   NewInfoPacked(st, offs[1]),
		roSchema:   ReadSchemaHamt(st, offs[2]),
		roInfo:     ReadInfoHamt(st, offs[3]),
	}
	ov.roSchemaOff = offs[2]
	ov.roInfoOff = offs[3]
	return &ov
}

//-------------------------------------------------------------------

// Merge is called by state.Merge to collect updates
// which are then applied by ApplyMerge
func (ov *Overlay) Merge(tranNum int) []update {
	return ov.roInfo.process(func(bto btOver) btOver {
		return bto.Merge(tranNum)
	})
}

func (ov *Overlay) ApplyMerge(updates []update) {
	ov.roInfo = ov.roInfo.withUpdates(updates, btOver.WithMerged)
}

//-------------------------------------------------------------------

//TODO schema

func (ov *Overlay) Persist() []update {
	return ov.roInfo.process(btOver.Save)
}

func (ov *Overlay) ApplyPersist(updates []update) {
	ov.roInfo = ov.roInfo.withUpdates(updates, btOver.WithSaved)
}
