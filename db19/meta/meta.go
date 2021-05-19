// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"math/bits"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/ints"
)

// Meta is the schema and info metadata
// difInfo is per transaction, overrides info
type Meta struct {
	schema      SchemaHamt
	info        InfoHamt
	difInfo     InfoHamt
	schemaOffs  []uint64
	infoOffs    []uint64
	schemaClock int
	infoClock   int
}

// Mutable returns a mutable copy of a Meta
func (m *Meta) Mutable() *Meta {
	assert.That(m.difInfo.IsNil())
	ov2 := *m // copy
	ov2.difInfo = InfoHamt{}.Mutable()
	return &ov2
}

// GetRoInfo returns read-only Info for the table or nil if not found
func (m *Meta) GetRoInfo(table string) *Info {
	if ti, ok := m.difInfo.Get(table); ok {
		return ti
	}
	if ti, ok := m.info.Get(table); ok && !ti.isTomb() {
		return ti
	}
	return nil
}

func (m *Meta) GetRwInfo(table string) *Info {
	if pti, ok := m.difInfo.Get(table); ok {
		return pti // already have mutable
	}
	pti, ok := m.info.Get(table)
	if !ok || pti.isTomb() {
		return nil
	}
	ti := *pti // copy
	ti.origNrows = ti.Nrows
	ti.origSize = ti.Size

	// set up index overlays
	ti.Indexes = append(ti.Indexes[:0:0], ti.Indexes...) // copy
	for i := range ti.Indexes {
		ti.Indexes[i] = ti.Indexes[i].Mutable()
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
	ts.lastmod = m.schemaClock
	ti.lastmod = m.infoClock
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
	schema.Put(m.newSchemaTomb(table))
	info := m.info.Mutable()
	info.Put(m.newInfoTomb(table))
	ov2 := *m // copy
	ov2.schema = schema.Freeze()
	ov2.info = info.Freeze()
	return &ov2
}

func (m *Meta) ForEachSchema(fn func(*Schema)) {
	m.schema.ForEach(func(schema *Schema) {
		if !schema.isTomb() {
			fn(schema)
		}
	})
}

func (m *Meta) ForEachInfo(fn func(*Info)) {
	m.info.ForEach(func(info *Info) {
		if !info.isTomb() {
			fn(info)
		}
	})
}

//-------------------------------------------------------------------

// LayeredOnto layers the mutable ixbuf's from a transaction
// onto the latest/current state and returns a new state.
// Also, the nrows and size deltas are applied.
// Note: this does not merge the ixbuf's, that is done later by merge.
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
		ti.Nrows = lti.Nrows + (ti.Nrows - ti.origNrows)
		ti.Size = lti.Size + (ti.Size - ti.origSize)
		ti.origNrows = 0
		ti.origSize = 0
		for i := range ti.Indexes {
			ti.Indexes[i].UpdateWith(lti.Indexes[i])
		}
		ti.lastmod = m.infoClock
		info.Put(ti)
	})
	result := *latest // copy
	result.info = info.Freeze()
	return &result
}

//-------------------------------------------------------------------

func (m *Meta) Write(store *stor.Stor, flatten bool) (offSchema, offInfo uint64) {
	assert.That(m.difInfo.IsNil())

	// schema
	npersists, timespan := mergeSize(m.schemaClock, flatten)
	// fmt.Printf("clock %d = %b npersists %d timespan %d\n", m.schemaClock, m.schemaClock, npersists, timespan)
	sfilter := func(ts *Schema) bool { return ts.lastmod >= m.schemaClock-timespan }
	if flatten || npersists >= len(m.schemaOffs) {
		sfilter = func(ts *Schema) bool { return !ts.isTomb() }
	}
	offSchema = m.schema.Write(store, nth(m.schemaOffs, npersists), sfilter)
	if offSchema != 0 {
		// fmt.Println("replace", m.schemaOffs, npersists, offSchema)
		m.schemaOffs = replace(m.schemaOffs, npersists, offSchema)
		// fmt.Println("    =>", m.schemaOffs)
		m.schemaClock++
		if len(m.schemaOffs) == 1 {
			m.schemaClock = delayMerge
		}
	} else if len(m.schemaOffs) > 0 {
		offSchema = m.schemaOffs[len(m.schemaOffs)-1]
	}

	// info
	npersists, timespan = mergeSize(m.infoClock, flatten)
	// fmt.Printf("clock %d = %b npersists %d timespan %d\n", m.infoClock, m.infoClock, npersists, timespan)
	ifilter := func(ti *Info) bool { return ti.lastmod >= m.infoClock-timespan }
	if flatten || npersists >= len(m.infoOffs) {
		ifilter = func(ti *Info) bool { return !ti.isTomb() }
	}
	offInfo = m.info.Write(store, nth(m.infoOffs, npersists), ifilter)
	// fmt.Println("replace", m.infoOffs, npersists, offInfo)
	m.infoOffs = replace(m.infoOffs, npersists, offInfo)
	// fmt.Println("    =>", m.infoOffs)
	m.infoClock++
	if len(m.infoOffs) == 1 {
		m.infoClock = delayMerge
	}

	return offSchema, offInfo
}

// mergeSize returns the number of persists to merge.
// 1 means lastmod == m.clock, 2 means lastmod >= m.clock-1, etc.
func mergeSize(clock int, flatten bool) (npersists, timespan int) {
	if flatten {
		clock = ints.MaxInt
	}
	trailingOnes := bits.TrailingZeros(^uint(clock))
	return trailingOnes, (1 << trailingOnes) - 1
}

func nth(v []uint64, n int) uint64 {
	if len(v) <= n {
		return 0
	}
	return v[n]
}

// replace replaces the first n elements with x
func replace(v []uint64, n int, x uint64) []uint64 {
	if n == 0 {
		if len(v) > 0 && v[0] == x {
			return v
		}
		v = append(v, 0)
		copy(v[1:], v)
	} else if n < len(v) {
		copy(v[1:], v[n:])
		v = v[:len(v)-(n-1)]
	} else if len(v) == 0 {
		return []uint64{x}
	} else {
		v = v[:1]
	}
	v[0] = x
	return v
}

func ReadMeta(store *stor.Stor, offSchema, offInfo uint64) *Meta {
	schema, schemaOffs := ReadSchemaChain(store, offSchema)
	info, infoOffs := ReadInfoChain(store, offInfo)
	m := Meta{
		schema: schema, schemaOffs: schemaOffs, schemaClock: clock(schemaOffs),
		info: info, infoOffs: infoOffs, infoClock: clock(infoOffs)}
	// set up ixspecs
	m.info.ForEach(func(ti *Info) {
		ts := m.schema.MustGet(ti.Table)
		for i := range ti.Indexes {
			ti.Indexes[i].SetIxspec(&ts.Indexes[i].Ixspec)
		}
	})
	return &m
}

const delayMerge = 0b1000000 // = 64 = approx 1 hour at 1 persist per minute

func clock(offs []uint64) int {
	switch len(offs) {
	case 0:
		return 0
	case 1:
		return delayMerge
	default:
		return ints.MaxInt
	}
}

//-------------------------------------------------------------------

func (m *Meta) CheckAllMerged() {
	m.info.ForEach(func(ti *Info) {
		for _, ov := range ti.Indexes {
			ov.CheckFlat()
		}
	})
}
