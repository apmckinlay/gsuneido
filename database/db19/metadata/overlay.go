// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package metadata

import "github.com/apmckinlay/gsuneido/database/db19/btree"

// Overlay handles layering a per-transaction mutable TableInfoHtbl,
// over an immutable TableInfoHtbl,
// over immutable InfoPacked and SchemaPacked
type Overlay struct {
	rwMeta     *TableInfoHtbl
	roMeta     *TableInfoHtbl
	baseInfo   *InfoPacked
	baseSchema *SchemaPacked
}

func NewOverlay(schemaPacked *SchemaPacked, infoPacked *InfoPacked,
	roMeta *TableInfoHtbl, rwMeta *TableInfoHtbl) *Overlay {
	return &Overlay{
		baseInfo:   infoPacked,
		baseSchema: schemaPacked,
		roMeta:     roMeta,
		rwMeta:     rwMeta,
	}
}

func (ov *Overlay) GetReadonly(table string) *TableInfo {
	if ti := ov.rwMeta.Get(table); ti != nil {
		return ti
	}
	if ti := ov.roMeta.Get(table); ti != nil {
		return ti
	}
	ti := ov.baseInfo.Get(table)
	ti.Schema = ov.baseSchema.Get(table)
	ov.rwMeta.Put(ti) // cache in memory
	return ti
}

func (ov *Overlay) GetMutable(table string, tranNum int) *TableInfo {
	var ti *TableInfo
	if ti = ov.rwMeta.Get(table); ti != nil {
		return ti // already have mutable
	}
	if ti = ov.roMeta.Get(table); ti != nil {
		ti2 := *ti
		ti = &ti2
	} else {
		ti = ov.baseInfo.Get(table)
		ti.Schema = ov.baseSchema.Get(table)
	}
	ti.mutable = true

	// setup up index overlays
	ti.Indexes = append([]*btree.Overlay(nil), ti.Indexes...) // copy
	for i := range ti.Indexes {
		ti.Indexes[i] = ti.Indexes[i].Mutable(tranNum)
	}

	ov.rwMeta.Put(ti) // cache in memory
	return ti
}

//-------------------------------------------------------------------

// LayeredOnto takes the mbtree's from a transaction
// and applies them to the latest DbState
// reusing the structs from the transaction
func (ov *Overlay) LayeredOnto(latest *TableInfoHtbl) *TableInfoHtbl {
	// have to copy the latest hash table because it may have more stuff
	result := latest.Dup() // shallow copy
	iter := ov.rwMeta.Iter()
	for ti := iter(); ti != nil; ti = iter() {
		if ti.mutable {
			if lti := latest.Get(ti.Table); lti != nil {
				ti.Schema = lti.Schema
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
			result.Put(ti)
		}
	}
	return result
}
