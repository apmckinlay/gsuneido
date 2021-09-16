// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"errors"
	"log"
	"math"
	"math/bits"

	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
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

// Mutable returns a mutable copy of a Meta. Used by UpdateTran.
func (m *Meta) Mutable() *Meta {
	assert.That(m.difInfo.IsNil())
	ov2 := *m // copy
	ov2.difInfo = InfoHamt{}.Mutable()
	return &ov2
}

func (m *Meta) SameSchemaAs(m2 *Meta) bool {
	return m.schema.root == m2.schema.root
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

// copyInfo is for debugging
func copyInfo(ti *Info) *Info {
	cp := *ti
	cp.Indexes = append(cp.Indexes[:0:0], cp.Indexes...) // copy
	for i, ov := range cp.Indexes {
		cp.Indexes[i] = ov.Copy()
	}
	return &cp
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
	if !ok || !ts.isTable() {
		return nil
	}
	return ts
}

func (m *Meta) GetView(name string) string {
	ts, ok := m.schema.Get("=" + name)
	if !ok || !ts.isView() {
		return ""
	}
	return ts.Columns[0]
}

func (m *Meta) ForEachSchema(fn func(*Schema)) {
	m.schema.ForEach(func(schema *Schema) {
		if schema.isTable() {
			fn(schema)
		}
	})
}

func (m *Meta) ForEachView(fn func(name, def string)) {
	m.schema.ForEach(func(schema *Schema) {
		if schema.isView() {
			fn(schema.Table[1:], schema.Columns[0])
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

type metaUpdate struct {
	meta   *Meta      // original
	schema SchemaHamt // mutable
	info   InfoHamt   // mutable
}

func newMetaUpdate(m *Meta) *metaUpdate {
	return &metaUpdate{meta: m}
}

func (mu *metaUpdate) getSchema(table string) *Schema {
	if ti, ok := mu.schema.Get(table); ok {
		cp := *ti // copy
		return &cp
	}
	return nil
}

func (mu *metaUpdate) putSchema(ts *Schema) {
	if mu.schema == (SchemaHamt{}) {
		mu.schema = mu.meta.schema.Mutable()
	}
	ts.lastmod = mu.meta.schemaClock
	mu.schema.Put(ts)
}

func (mu *metaUpdate) putInfo(ti *Info) {
	if mu.info == (InfoHamt{}) {
		mu.info = mu.meta.info.Mutable()
	}
	ti.lastmod = mu.meta.infoClock
	mu.info.Put(ti)
}

func (mu *metaUpdate) freeze() *Meta {
	cp := *mu.meta
	if mu.schema != (SchemaHamt{}) {
		cp.schema = mu.schema.Freeze()
	}
	if mu.info != (InfoHamt{}) {
		cp.info = mu.info.Freeze()
	}
	return &cp
}

// admin schema changes ---------------------------------------------

//TODO Derived

func (m *Meta) Ensure(a *schema.Schema, store *stor.Stor) ([]schema.Index, *Meta, error) {
	ts, ok := m.schema.Get(a.Table)
	if !ok || ts.isTomb() {
		panic("ensure: couldn't find " + a.Table)
	}
	ts, ti := m.alterGet(a.Table)
	newCols := sset.Difference(a.Columns, ts.Columns)
	newIdxs := []schema.Index{}
	for i := range a.Indexes {
		if nil == ts.FindIndex(a.Indexes[i].Columns) {
			newIdxs = append(newIdxs, a.Indexes[i])
		}
	}
	if err := createColumns(ts, newCols); err != nil {
		return nil, nil, err
	}
	if err := createIndexes(ts, ti, newIdxs, store); err != nil {
		return nil, nil, err
	}
	ac := &schema.Schema{Table: a.Table, Indexes: newIdxs}
	return newIdxs, m.PutNew(ts, ti, ac), nil
}

func (m *Meta) RenameTable(from, to string) *Meta {
	ts, ok := m.schema.Get(from)
	if !ok || ts.isTomb() {
		panic("can't rename nonexistent table: " + from)
	}
	tsNew := *ts // copy
	tsNew.Table = to
	if tmp, ok := m.schema.Get(to); ok && !tmp.isTomb() {
		panic("can't rename to existing table: " + to)
	}
	ti, ok := m.info.Get(from)
	assert.That(ok && ti != nil)
	tiNew := *ti // copy
	tiNew.Table = to

	mu := newMetaUpdate(m)
	mu.putSchema(m.newSchemaTomb(from))
	mu.putSchema(&tsNew)
	mu.putInfo(m.newInfoTomb(from))
	mu.putInfo(&tiNew)
	m.dropFkeys(mu, &ts.Schema)
	m.createFkeys(mu, &tsNew.Schema)
	return mu.freeze()
}

// Drop removes a table or view
func (m *Meta) Drop(name string) *Meta {
	//TODO delete without tombstone if not persisted e.g. tests
	if m.GetView(name) != "" {
		return m.Put(m.newSchemaTomb("="+name), nil)
	}
	ts, ok := m.schema.Get(name)
	if !ok || ts.isTomb() {
		return nil // nonexistent
	}
	mu := newMetaUpdate(m)
	mu.putSchema(m.newSchemaTomb(name))
	mu.putInfo(m.newInfoTomb(name))
	m.dropFkeys(mu, &ts.Schema)
	return mu.freeze()
}

func (m *Meta) AlterRename(table string, from, to []string) *Meta {
	ts, ok := m.schema.Get(table)
	if !ok || ts.isTomb() {
		panic("can't alter nonexistent table: " + table)
	}
	missing := sset.Difference(from, ts.Columns)
	if len(missing) > 0 {
		panic("can't rename nonexistent column(s): " + strs.Join(", ", missing))
	}
	existing := sset.Intersect(to, ts.Columns)
	if len(existing) > 0 {
		panic("can't rename to existing column(s): " + strs.Join(", ", existing))
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
	// ixspecs are ok since they are field indexes, not names
	mu := newMetaUpdate(m)
	mu.putSchema(&tsNew)
	m.dropFkeys(mu, &ts.Schema)
	m.createFkeys(mu, &tsNew.Schema)
	return mu.freeze()
}

func (m *Meta) AlterCreate(ac *schema.Schema, store *stor.Stor) (*Meta, error) {
	ts, ti := m.alterGet(ac.Table)
	if err := createColumns(ts, ac.Columns); err != nil {
		return nil, err
	}
	if err := createIndexes(ts, ti, ac.Indexes, store); err != nil {
		return nil, err
	}
	return m.PutNew(ts, ti, ac), nil
}

func (m *Meta) PutNew(ts *Schema, ti *Info, ac *schema.Schema) *Meta {
	mu := newMetaUpdate(m)
	mu.putSchema(ts)
	mu.putInfo(ti)
	m.createFkeys(mu, ac)
	return mu.freeze()
}

func (m *Meta) alterGet(table string) (*Schema, *Info) {
	ts, ok := m.schema.Get(table)
	if !ok || ts.isTomb() {
		panic("can't alter nonexistent table: " + table)
	}
	tsNew := *ts // copy
	ti, ok := m.info.Get(table)
	assert.That(ok && ti != nil)
	tiNew := *ti // copy
	return &tsNew, &tiNew
}

func createColumns(ts *Schema, cols []string) error {
	existing := sset.Intersect(cols, ts.Columns)
	if len(existing) > 0 {
		return errors.New("can't create existing column(s): " +
			strs.Join(", ", existing))
	}
	ts.Columns = append(strs.Cow(ts.Columns), cols...)
	return nil
}

func createIndexes(ts *Schema, ti *Info, idxs []schema.Index, store *stor.Stor) error {
	if len(idxs) == 0 {
		return nil
	}
	for i := range idxs {
		missing := sset.Difference(idxs[i].Columns, ts.Columns)
		if len(missing) > 0 {
			return errors.New("can't create index on nonexistent column(s): " +
				strs.Join(", ", missing))
		}
	}
	ts.Ixspecs(idxs)
	n := len(ts.Indexes)
	ts.Indexes = append(ts.Indexes[:n:n], idxs...)
	n = len(ti.Indexes)
	ti.Indexes = ti.Indexes[:n:n] // copy on write
	for i := range idxs {
		bt := btree.CreateBtree(store, &ts.Indexes[n+i].Ixspec)
		ti.Indexes = append(ti.Indexes, index.OverlayFor(bt))
	}
	return nil
}

func (m *Meta) createFkeys(mu *metaUpdate, ac *schema.Schema) {
	idxs := ac.Indexes
	for i := range idxs {
		fk := &idxs[i].Fk
		if fk.Table != "" {
			fkCols := fk.Columns
			if len(fkCols) == 0 {
				fkCols = idxs[i].Columns
			}
			target := mu.getSchema(fk.Table)
			if target == nil {
				log.Println("foreign key: can't find", fk.Table, "(from "+ac.Table+")")
				continue
			}
			target.Indexes = append(target.Indexes[:0:0], target.Indexes...) // copy
			for j := range target.Indexes {
				ix := target.Indexes[j] // copy
				if strs.Equal(fkCols, ix.Columns) {
					fk.IIndex = j
					n := len(ix.FkToHere)
					ix.FkToHere = append(ix.FkToHere[:n:n], // copy on write
						Fkey{Table: ac.Table,
							Columns: idxs[i].Columns, IIndex: i, Mode: fk.Mode})
					target.Indexes[j] = ix
				}
			}
			mu.putSchema(target)
		}
	}
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
	mu := newMetaUpdate(m)
	mu.putSchema(ts)
	mu.putInfo(ti)
	m.dropFkeys(mu, ad)
	return mu.freeze()
}

func dropIndexes(ts *Schema, ti *Info, idxs []schema.Index) {
loop:
	for j := range idxs {
		for i := range ts.Indexes {
			if strs.Equal(ts.Indexes[i].Columns, idxs[j].Columns) {
				continue loop
			}
		}
		panic("can't drop nonexistent index: " + strs.Join(",", idxs[j].Columns))
	}
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
	missing := sset.Difference(cols, ts.Columns)
	if len(missing) > 0 {
		panic("can't drop nonexistent column(s): " + strs.Join(", ", missing))
	}
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

func (m *Meta) dropFkeys(mu *metaUpdate, ad *schema.Schema) {
	// unlike createFkeys
	// we need to get the actual schema to get the foreign key information
	schema := &m.GetRoSchema(ad.Table).Schema
	idxs := ad.Indexes
	for i := range idxs {
		idx := schema.FindIndex(idxs[i].Columns)
		fk := idx.Fk
		if fk.Table == "" || fk.Table == ad.Table {
			continue
		}
		fkCols := fk.Columns
		if len(fkCols) == 0 {
			fkCols = idx.Columns
		}
		target, ok := m.schema.Get(fk.Table)
		if !ok {
			log.Println("foreign key: can't find", fk.Table, "(from "+ad.Table+")")
			continue
		}
		target.Indexes = append(target.Indexes[:0:0], target.Indexes...) // copy
		for j := range target.Indexes {
			ix := target.Indexes[j] // copy
			if strs.Equal(fkCols, ix.Columns) {
				fk.IIndex = j
				fkToHere := make([]Fkey, 0, len(ix.FkToHere))
				for k := range ix.FkToHere {
					fk := ix.FkToHere[k]
					if ad.Table != fk.Table ||
						!strs.Equal(idx.Columns, fk.Columns) {
						fkToHere = append(fkToHere, fk)
					}
				}
				ix.FkToHere = fkToHere
				target.Indexes[j] = ix
			}
		}
		mu.putSchema(target)
	}
}

func (m *Meta) AddView(name, def string) *Meta {
	if m.GetView(name) != "" {
		panic("view: '" + name + "' already exists")
	}
	return m.Put(m.newSchemaView(name, def), nil)
}

// TouchTable is for tests
func (m *Meta) TouchTable(table string) *Meta {
	schema := *m.GetRoSchema(table) // copy
	mu := newMetaUpdate(m)
	mu.putSchema(&schema)
	return mu.freeze()
}

// TouchIndexes is for tests
func (m *Meta) TouchIndexes(table string) *Meta {
	schema := *m.GetRoSchema(table)                                  // copy
	schema.Indexes = append(schema.Indexes[:0:0], schema.Indexes...) // copy
	mu := newMetaUpdate(m)
	mu.putSchema(&schema)
	return mu.freeze()
}

//-------------------------------------------------------------------

// LayeredOnto is done by transaction commit.
// It layers the mutable ixbuf's from transactions
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
		assert.That(ti.Nrows >= 0)
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
	if offInfo != 0 {
		// fmt.Println("replace", m.infoOffs, npersists, offInfo)
		m.infoOffs = replace(m.infoOffs, npersists, offInfo)
		// fmt.Println("    =>", m.infoOffs)
		m.infoClock++
		if len(m.infoOffs) == 1 {
			m.infoClock = delayMerge
		}
	} else if len(m.infoOffs) > 0 {
		offInfo = m.infoOffs[len(m.infoOffs)-1]
	}

	return offSchema, offInfo
}

// mergeSize returns the number of persists to merge.
// 1 means lastmod == m.clock, 2 means lastmod >= m.clock-1, etc.
func mergeSize(clock int, flatten bool) (npersists, timespan int) {
	if flatten {
		clock = math.MaxInt
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

type Fkey = schema.Fkey

func ReadMeta(store *stor.Stor, offSchema, offInfo uint64) *Meta {
	schema, schemaOffs := ReadSchemaChain(store, offSchema)
	info, infoOffs := ReadInfoChain(store, offInfo)
	m := Meta{
		schema: schema, schemaOffs: schemaOffs, schemaClock: clock(schemaOffs),
		info: info, infoOffs: infoOffs, infoClock: clock(infoOffs)}
	// copy Ixspec to Info from Schema (constructed by ReadSchema)
	m.info.ForEach(func(ti *Info) {
		if ti.isTomb() {
			return
		}
		ts := m.schema.MustGet(ti.Table)
		for i := range ti.Indexes {
			ti.Indexes[i].SetIxspec(&ts.Indexes[i].Ixspec)
		}
	})
	// link foreign keys to targets
	m.schema.ForEach(func(s *Schema) {
		if s.isTomb() {
			return
		}
		for i := range s.Indexes {
			fk := &s.Indexes[i].Fk
			if fk.Table != "" {
				fkCols := fk.Columns
				if len(fkCols) == 0 {
					fkCols = s.Indexes[i].Columns
				}
				// ok to modify actual schema because it's not shared yet
				target, ok := m.schema.Get(fk.Table)
				if !ok {
					log.Println("foreign key: can't find", fk.Table, "(from "+s.Table+")")
					continue
				}
				for j := range target.Indexes {
					ix := &target.Indexes[j]
					if strs.Equal(fkCols, ix.Columns) {
						fk.IIndex = j
						ix.FkToHere = append(ix.FkToHere,
							Fkey{Table: s.Table, Mode: fk.Mode,
								Columns: s.Indexes[i].Columns, IIndex: i})
					}
				}
			}
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
		return math.MaxInt
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
