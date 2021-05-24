// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"math/bits"

	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/strs"
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

// Put is used by Database.LoadedTable and admin schema changes
func (m *Meta) Put(ts *Schema, ti *Info) *Meta {
	cp := *m // copy
	ts.lastmod = m.schemaClock
	schema := m.schema.Mutable()
	schema.Put(ts)
	cp.schema = schema.Freeze()
	if ti != nil {
		ti.lastmod = m.infoClock
		info := m.info.Mutable()
		info.Put(ti)
		cp.info = info.Freeze()
	}
	return &cp
}

// admin schema changes ---------------------------------------------

//TODO Derived

func (m *Meta) Ensure(a *schema.Schema, store *stor.Stor) *Meta {
	ts, ti := m.alterGet(a.Table)
	newCols := sset.Difference(a.Columns, ts.Columns)
	newIdxs := []schema.Index{}
outer:
	for i := range a.Indexes {
		for j := range ts.Indexes {
			if strs.Equal(a.Indexes[i].Columns, ts.Indexes[j].Columns) {
				continue outer
			}
		}
		newIdxs = append(newIdxs, a.Indexes[i])
	}
	if ti.Nrows > 0 && len(newIdxs) > 0 {
		panic("creating indexes on tables with data not implemented") //TODO
	}
	if !createColumns(ts, newCols) ||
		!createIndexes(ts, ti, newIdxs, store) {
		return nil
	}
	return m.Put(ts, ti)
}

func (m *Meta) RenameTable(from, to string) *Meta {
	ts, ok := m.schema.Get(from)
	if !ok || ts.isTomb() {
		return nil // from doesn't exist
	}
	tsNew := *ts // copy
	if tmp, ok := m.schema.Get(to); ok && !tmp.isTomb() {
		return nil // to exists
	}
	ti, ok := m.info.Get(from)
	assert.That(ok && ti != nil)
	tiNew := *ti // copy

	schema := m.schema.Mutable()
	schema.Put(m.newSchemaTomb(from))
	tsNew.Table = to
	tsNew.lastmod = m.schemaClock
	schema.Put(&tsNew)

	info := m.info.Mutable()
	info.Put(m.newInfoTomb(from))
	tiNew.Table = to
	tiNew.lastmod = m.infoClock
	info.Put(&tiNew)

	cp := *m // copy
	cp.schema = schema.Freeze()
	cp.info = info.Freeze()
	return &cp
}

func (m *Meta) DropTable(table string) *Meta {
	if ts, ok := m.schema.Get(table); !ok || ts.isTomb() {
		return nil // nonexistent
	}
	//TODO delete without tombstone if not persisted
	return m.Put(m.newSchemaTomb(table), m.newInfoTomb(table))
}

func (m *Meta) AlterRename(table string, from, to []string) *Meta {
	ts, ok := m.schema.Get(table)
	if !ok || ts.isTomb() {
		return nil // nonexistent
	}
	tsNew := *ts // copy
	tsNew.Columns = strs.Replace(ts.Columns, from, to)
	tsNew.Derived = strs.Replace(ts.Derived, from, to)
	tsNew.Indexes = make([]schema.Index, len(ts.Indexes))
	copy(tsNew.Indexes, ts.Indexes)
	for i := range tsNew.Indexes {
		ix := &tsNew.Indexes[i]
		ix.Columns = strs.Replace(ix.Columns, from, to)
	}
	return m.Put(&tsNew, nil)
}

func (m *Meta) AlterCreate(ac *schema.Schema, store *stor.Stor) *Meta {
	ts, ti := m.alterGet(ac.Table)
	if ti.Nrows > 0 && len(ac.Indexes) > 0 {
		panic("creating indexes on tables with data not implemented") //TODO
	}
	if !createColumns(ts, ac.Columns) ||
		!createIndexes(ts, ti, ac.Indexes, store) {
		return nil
	}
	return m.Put(ts, ti)
}

func (m *Meta) alterGet(table string) (*Schema, *Info) {
	ts, ok := m.schema.Get(table)
	if !ok || ts.isTomb() {
		return nil, nil // nonexistent
	}
	tsNew := *ts // copy
	ti, ok := m.info.Get(table)
	assert.That(ok && ti != nil)
	tiNew := *ti // copy
	return &tsNew, &tiNew
}

func createColumns(ts *Schema, cols []string) bool {
	if !sset.Disjoint(ts.Columns, cols) {
		return false
	}
	ts.Columns = append(strs.Cow(ts.Columns), cols...)
	return true
}

func createIndexes(ts *Schema, ti *Info, idxs []schema.Index, store *stor.Stor) bool {
	if len(idxs) == 0 {
		return true
	}
	for i := range idxs {
		if !sset.Subset(ts.Columns, idxs[i].Columns) {
			return false
		}
	}
	ts.Ixspecs(idxs)
	n := len(ts.Indexes)
	ts.Indexes = append(ts.Indexes[:n:n], idxs...)
	n = len(ti.Indexes)
	ti.Indexes = ti.Indexes[:n:n]
	for i := range idxs {
		bt := btree.CreateBtree(store, &ts.Indexes[i].Ixspec)
		ti.Indexes = append(ti.Indexes, index.OverlayFor(bt))
	}
	return true
}

func (m *Meta) AlterDrop(ad *schema.Schema) *Meta {
	ts, ti := m.alterGet(ad.Table)
	// need to drop indexes before columns
	// in case we drop a column and an index that contains it
	if len(ad.Indexes) > 0 {
		dropIndexes(ts, ti, ad.Indexes)
	}
	if len(ad.Columns) > 0 {
		if !dropColumns(ts, ad.Columns) {
			return nil
		}
	}
	return m.Put(ts, ti)
}

func dropIndexes(ts *Schema, ti *Info, idxs []schema.Index) {
	tsIdxs := make([]schema.Index, 0, len(ts.Indexes))
	tiIdxs := make([]*index.Overlay, 0, len(ti.Indexes))
outer:
	for i := range ts.Indexes {
		for j := range idxs {
			if strs.Equal(ts.Indexes[i].Columns, idxs[j].Columns) {
				continue outer // i.e. don't copy deletion
			}
		}
		tsIdxs = append(tsIdxs, ts.Indexes[i])
		tiIdxs = append(tiIdxs, ti.Indexes[i])
	}
	ts.Indexes = tsIdxs
	ti.Indexes = tiIdxs
}

func dropColumns(ts *Schema, cols []string) bool {
	for i := range ts.Indexes {
		if !sset.Disjoint(ts.Indexes[i].Columns, cols) {
			return false // can't drop if used by index
		}
	}
	for _, col := range cols {
		ts.Columns = strs.Replace1(ts.Columns, col, "-")
	}
	return true
}

//-------------------------------------------------------------------

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
