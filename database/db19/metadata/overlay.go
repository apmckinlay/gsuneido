// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package metadata

// Overlay handles layering a per-transaction mutable TableInfoHtbl,
// over an immutable TableInfoHtbl,
// over immutable InfoPacked and SchemaPacked
type Overlay struct {
	rwInfo       *TableInfoHtbl
	roInfo       *TableInfoHtbl
	infoPacked   *InfoPacked
	schemaPacked *SchemaPacked
}

func (ov *Overlay) GetReadonly(table string) *TableInfo {
	if ti := ov.rwInfo.Get(table); ti != nil {
		return ti
	}
	if ti := ov.roInfo.Get(table); ti != nil {
		return ti
	}
	ti := ov.infoPacked.Get(table)
	ti.schema = ov.schemaPacked.Get(table)
	ov.rwInfo.Put(ti) // cache in memory
	return ti
}

func (ov *Overlay) GetMutable(table string) *TableInfo {
	if ti := ov.rwInfo.Get(table); ti != nil {
		return ti
	}
	if ti := ov.roInfo.Get(table); ti != nil {
		ti2 := *ti
		schema := *ti.schema
		ti2.schema = &schema
		ti2.mutable = true
		ov.rwInfo.Put(&ti2)
		return &ti2
	}
	ti := ov.infoPacked.Get(table)
	ti.schema = ov.schemaPacked.Get(table)
	ti.mutable = true
	ov.rwInfo.Put(ti) // cache in memory
	return ti
}
