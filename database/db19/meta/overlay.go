// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"github.com/apmckinlay/gsuneido/database/db19/btree"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/verify"
)

type Overlay struct {
	//TODO rwSchema
	rwInfo      *InfoHtbl
	roInfo      *InfoHtbl
	roInfoOff   uint64
	roSchema    *SchemaHtbl
	roSchemaOff uint64
	baseInfo    *InfoPacked
	baseSchema  *SchemaPacked
}

func NewOverlay(baseSchema *SchemaPacked, baseInfo *InfoPacked,
	roSchema *SchemaHtbl, roSchemaOff uint64,
	roInfo *InfoHtbl, roInfoOff uint64,
	rwInfo *InfoHtbl) *Overlay {
	return &Overlay{
		baseSchema:  baseSchema,
		baseInfo:    baseInfo,
		roInfo:      roInfo,
		roInfoOff:   roInfoOff,
		roSchema:    roSchema,
		roSchemaOff: roSchemaOff,
		rwInfo:      rwInfo,
	}
}

// NewOverlay returns a new Overlay based on an existing one
func (ov *Overlay) NewOverlay() *Overlay {
	verify.That(ov.rwInfo == nil)
	ov2 := *ov
	ov2.rwInfo = NewInfoHtbl(0)
	return &ov2
}

func (ov *Overlay) GetRoInfo(table string) *Info {
	if ti := ov.rwInfo.Get(table); ti != nil {
		return ti
	}
	if ti := ov.roInfo.Get(table); ti != nil {
		return ti
	}
	ti := ov.baseInfo.Get(table)
	ov.rwInfo.Put(ti) // cache in memory
	return ti
}

func (ov *Overlay) GetRwInfo(table string, tranNum int) *Info {
	var ti *Info
	if ti = ov.rwInfo.Get(table); ti != nil {
		return ti // already have mutable
	}
	if ti = ov.roInfo.Get(table); ti != nil {
		ti2 := *ti // copy
		ti = &ti2
	} else {
		ti = ov.baseInfo.Get(table) // this will be a copy
	}
	// start at 0 since these are deltas
	ti.Nrows = 0
	ti.Size = 0
	ti.mutable = true

	// set up index overlays
	ti.Indexes = append([]*btree.Overlay(nil), ti.Indexes...) // copy
	for i := range ti.Indexes {
		ti.Indexes[i] = ti.Indexes[i].Mutable(tranNum)
	}

	ov.rwInfo.Put(ti) // cache in memory
	return ti
}

//-------------------------------------------------------------------

// LayeredOnto takes the mutable mbtree's from a transaction
// and applies them to the latest DbState
// reusing the structs from the transaction
func (ov *Overlay) LayeredOnto(latest *Overlay) *Overlay {
	// start with a copy of the latest hash table because it may have more
	verify.That(latest.rwInfo == nil)
	roInfo2 := latest.roInfo.Dup() // shallow copy
	iter := ov.rwInfo.Iter()
	for ti := iter(); ti != nil; ti = iter() {
		if ti.mutable {
			if lti := roInfo2.Get(ti.Table); lti != nil {
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
	}
	result := *latest
	result.roInfo = roInfo2
	return &result
}

//-------------------------------------------------------------------

const Noffsets = 4

type offsets = [Noffsets]uint64

func (ov *Overlay) Write(st *stor.Stor) offsets {
	verify.That(ov.rwInfo == nil)
	return offsets{
		ov.baseSchema.Offset(),
		ov.baseInfo.Offset(),
		ov.roSchemaOff,
		ov.roInfo.Write(st),
	}
}

func FromOffsets(st *stor.Stor, offs offsets) *Overlay {
	ov := Overlay{
		baseSchema: NewSchemaPacked(st, offs[0]),
		baseInfo:   NewInfoPacked(st, offs[1]),
		roSchema:   ReadSchemaHtbl(st, offs[2]),
		roInfo:     ReadInfoHtbl(st, offs[3]),
	}
	ov.roSchemaOff = offs[2]
	ov.roInfoOff = offs[3]
	return &ov
}

//-------------------------------------------------------------------

func (ov *Overlay) Merge(tranNum int) []update {
	return ov.roInfo.process(func(bto btOver) btOver {
		return bto.Merge(tranNum)
	})
}
func (ov *Overlay) ApplyMerge(updates []update) {
	ov.roInfo = ov.roInfo.withUpdates(updates, btOver.WithMerged)
}

//-------------------------------------------------------------------

func (ov *Overlay) SaveIndexes() []update {
	return ov.roInfo.process(btOver.Save)
}

func (ov *Overlay) ApplySave(updates []update) {
	ov.roInfo = ov.roInfo.withUpdates(updates, btOver.WithSaved)
}
