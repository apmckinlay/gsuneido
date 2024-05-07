// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"log"

	"slices"

	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/hamt"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
)

// Meta is the schema and info metadata.
// See also: schema.go - Schema, info.go - Info
type Meta struct {
	// difInfo is per transaction updates, overrides info
	difInfo map[string]*Info
	// schema and info are immutable persistent hash tables.
	// WARNING: do NOT modify the items from the hash tables.
	// To modify an item, Get it, copy it, modify it, then Put the new version.
	// Unfortunately, Go does not have any way to enforce this.
	schema hamt.Chain[string, *Schema]
	info   hamt.Chain[string, *Info]
}

func (m *Meta) Cksum() uint32 {
	return m.schema.Cksum() + m.info.Cksum()
}

func (m *Meta) CksumData() uint32 {
	return m.schema.Hamt.Cksum() + m.info.Hamt.Cksum()
}

func (m *Meta) ResetClock() { // for testing only
	m.schema.Clock = 0
	m.info.Clock = 0
}

// Mutable returns a mutable copy of a Meta. Used by UpdateTran.
func (m *Meta) Mutable() *Meta {
	assert.That(m.difInfo == nil)
	ov2 := *m // copy
	ov2.difInfo = make(map[string]*Info)
	return &ov2
}

func (m *Meta) SameSchemaAs(m2 *Meta) bool {
	return m.schema.Hamt.SameAs(m2.schema.Hamt)
}

// GetRoInfo returns read-only Info for the table or nil if not found
func (m *Meta) GetRoInfo(table string) *Info {
	if ti, ok := m.difInfo[table]; ok {
		return ti
	}
	if ti, ok := m.info.Get(table); ok && !ti.IsTomb() {
		return ti
	}
	return nil
}

//lint:ignore U1000 for debugging
func copyInfo(ti *Info) *Info {
	cp := *ti
	cp.Indexes = slices.Clone(cp.Indexes)
	for i, ov := range cp.Indexes {
		cp.Indexes[i] = ov.Copy()
	}
	return &cp
}

func (m *Meta) GetRwInfo(table string) *Info {
	if pti, ok := m.difInfo[table]; ok {
		return pti // already have mutable
	}
	pti, ok := m.info.Get(table)
	if !ok || pti.IsTomb() {
		return nil
	}
	ti := *pti // copy

	ti.Indexes = slices.Clone(ti.Indexes)
	for i := range ti.Indexes {
		ti.Indexes[i] = ti.Indexes[i].Mutable()
	}

	m.difInfo[table] = &ti
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
		if !info.IsTomb() {
			fn(info)
		}
	})
}

// Put is used by Database.LoadedTable and admin schema changes
func (m *Meta) Put(ts *Schema, ti *Info) *Meta {
	cp := *m // copy
	ts.lastMod = m.schema.Clock
	schema := m.schema.Mutable()
	schema.Put(ts)
	cp.schema.Hamt = schema.Freeze()
	if ti != nil {
		ti.lastMod = m.info.Clock
		info := m.info.Mutable()
		info.Put(ti)
		cp.info.Hamt = info.Freeze()
	}
	return &cp
}

// PutNew sets created so drop knows it doesn't need a tombstone
func (m *Meta) PutNew(ts *Schema, ti *Info, ac *schema.Schema) *Meta {
	if _, ok := m.schema.Get(ts.Table); !ok {
		ts.created = m.schema.Clock
	}
	if _, ok := m.info.Get(ti.Table); !ok {
		ti.created = m.info.Clock
	}
	mu := newMetaUpdate(m)
	mu.putSchema(ts)
	mu.putInfo(ti)
	m.createFkeys(mu, &ts.Schema, ac)
	return mu.freeze()
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
	assert.That(len(ts.Indexes) > 0 || ts.IsTomb())
	if mu.schema == (SchemaHamt{}) {
		mu.schema = mu.meta.schema.Mutable()
	}
	ts.lastMod = mu.meta.schema.Clock
	mu.schema.Put(ts)
}

func (mu *metaUpdate) putInfo(ti *Info) {
	if mu.info == (InfoHamt{}) {
		mu.info = mu.meta.info.Mutable()
	}
	ti.lastMod = mu.meta.info.Clock
	mu.info.Put(ti)
}

func (mu *metaUpdate) freeze() *Meta {
	cp := *mu.meta
	if mu.schema != (SchemaHamt{}) {
		cp.schema.Hamt = mu.schema.Freeze()
	}
	if mu.info != (InfoHamt{}) {
		cp.info.Hamt = mu.info.Freeze()
	}
	return &cp
}

// admin schema changes ---------------------------------------------

// Ensure returns nil newIdxs if there is nothing more to be done
// i.e. if there are no new indexes or if there is no data yet.
func (m *Meta) Ensure(a *schema.Schema, store *stor.Stor) ([]schema.Index, *Meta) {
	ts, ok := m.schema.Get(a.Table)
	if !ok || ts.IsTomb() {
		panic("ensure: couldn't find " + a.Table)
	}
	ts, ti := m.alterGet(a.Table)
	var newIdxs []schema.Index
	for i := range a.Indexes {
		if nil == ts.FindIndex(a.Indexes[i].Columns) {
			newIdxs = append(newIdxs, a.Indexes[i])
		}
	}
	newCols := set.Difference(a.Columns, ts.Columns)
	createColumns(ts, newCols)
	newDer := set.Difference(a.Derived, ts.Derived)
	createDerived(ts, newDer)
	createIndexes(ts, ti, newIdxs, store)
	ac := &schema.Schema{Table: a.Table, Indexes: newIdxs}
	if ti.Nrows == 0 {
		newIdxs = nil
	}
	mu := newMetaUpdate(m)
	mu.putSchema(ts)
	mu.putInfo(ti)
	m.createFkeys(mu, &ts.Schema, ac)
	return newIdxs, mu.freeze()
}

func (m *Meta) RenameTable(from, to string) *Meta {
	ts, ok := m.schema.Get(from)
	if !ok || ts.IsTomb() {
		panic("can't rename nonexistent table: " + from)
	}
	tsNew := *ts // copy
	tsNew.Table = to
	if tmp, ok := m.schema.Get(to); ok && !tmp.IsTomb() {
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
	m.createFkeys(mu, &tsNew.Schema, &tsNew.Schema)
	return mu.freeze()
}

// Drop removes a table or view
func (m *Meta) Drop(name string) *Meta {
	if m.GetView(name) != "" {
		// view
		return m.Put(m.newSchemaTomb("="+name), nil)
	}
	// table
	ts, ok := m.schema.Get(name)
	if !ok || ts.IsTomb() {
		return nil // nonexistent
	}
	if list := fkToHere(&ts.Schema); list != nil {
		panic("can't drop table used by foreign keys: " +
			name + " <- " + str.Join(",", list))
	}
	mu := newMetaUpdate(m)
	if ts.created != 0 && ts.created == m.schema.Clock {
		// not persisted so no need for tombstone
		mu.schema = mu.meta.schema.Mutable()
		mu.schema.Delete(ts.Table)
	} else {
		mu.putSchema(m.newSchemaTomb(name))
	}
	ti := m.schema.MustGet(ts.Table)
	if ti.created != 0 && ti.created == m.info.Clock {
		// not persisted so no need for tombstone
		mu.info = mu.meta.info.Mutable()
		mu.info.Delete(ti.Table)
	} else {
		mu.putInfo(m.newInfoTomb(name))
	}
	m.dropFkeys(mu, &ts.Schema)
	return mu.freeze()
}

func fkToHere(ts *schema.Schema) []string {
	var list []string
	for i := range ts.Indexes {
		ix := ts.Indexes[i]
		for j := range ix.FkToHere {
			fk := ix.FkToHere[j]
			if fk.Table != ts.Table {
				list = append(list, fk.Table)
			}
		}
	}
	return list
}

func (m *Meta) AlterRename(table string, from, to []string) *Meta {
	ts, ok := m.schema.Get(table)
	if !ok || ts.IsTomb() {
		panic("can't alter nonexistent table: " + table)
	}
	tsNew := *ts // copy
	tsNew.Columns = replaceUnique(ts.Columns, from, to)
	tsNew.Derived = replace(ts.Derived, from, to)
	tsNew.Indexes = slices.Clone(ts.Indexes)
	for i := range tsNew.Indexes {
		ix := &tsNew.Indexes[i]
		cols := replace(ix.Columns, from, to)
		if !slc.Same(cols, ix.Columns) && tsNew.FindIndex(cols) != nil {
			panic("rename causes duplicate index: " + str.Join("(,)", cols))
		}
		ix.Columns = cols
		ix.BestKey = replace(ix.BestKey, from, to)
	}
	// ixspecs are ok since they are field indexes, not names
	mu := newMetaUpdate(m)
	mu.putSchema(&tsNew)
	m.dropFkeys(mu, &ts.Schema)
	m.createFkeys(mu, &tsNew.Schema, &tsNew.Schema)
	return mu.freeze()
}

// replaceUnique replaces occurrences of from with to.
// Assumes list values are unique (no duplicates).
// from values must exist.
// Replacements must preserve uniqueness.
// Replacements are done in from/to order.
func replaceUnique(list, from, to []string) []string {
	list = slices.Clone(list)
	for i, f := range from {
		j := slices.Index(list, f)
		if j == -1 {
			panic("can't rename nonexistent column: " + f)
		}
		if slices.Contains(list, to[i]) {
			panic("can't rename to existing column: " + to[i])
		}
		list[j] = to[i]
	}
	return list
}

// replace replaces occurrences of from with to.
// Replacements are done in from/to order.
func replace(list, from, to []string) []string {
	cloned := false
	for i, f := range from {
		for j := range list {
			if list[j] == f {
				if !cloned {
					list = slices.Clone(list)
					cloned = true
				}
				list[j] = to[i]
			}
		}
	}
	return list
}

func (m *Meta) AlterCreate(ac *schema.Schema, store *stor.Stor) *Meta {
	ts, ti := m.alterGet(ac.Table)
	createColumns(ts, ac.Columns)
	createDerived(ts, ac.Derived)
	createIndexes(ts, ti, ac.Indexes, store)
	mu := newMetaUpdate(m)
	mu.putSchema(ts)
	mu.putInfo(ti)
	m.createFkeys(mu, &ts.Schema, ac)
	return mu.freeze()
}

func (m *Meta) alterGet(table string) (*Schema, *Info) {
	ts, ok := m.schema.Get(table)
	if !ok || ts.IsTomb() {
		panic("can't alter nonexistent table: " + table)
	}
	tsNew := *ts // copy
	ti, ok := m.info.Get(table)
	assert.That(ok && ti != nil)
	tiNew := *ti // copy
	return &tsNew, &tiNew
}

func createColumns(ts *Schema, cols []string) {
	existing := set.Intersect(cols, ts.Columns)
	if len(existing) > 0 {
		panic("can't create existing column(s): " + str.Join(", ", existing))
	}
	ts.Columns = slc.With(ts.Columns, cols...)
}

func createDerived(ts *Schema, cols []string) {
	existing := set.Intersect(cols, ts.Derived)
	if len(existing) > 0 {
		panic("can't create existing column(s): " + str.Join(", ", existing))
	}
	ts.Derived = slc.With(ts.Derived, cols...)
}

// createIndexes appends the new indexes to ts.Indexes
// and appends empty overlays for them to ti.Indexes
// It does not build the btrees, that's done by buildIndexes.
func createIndexes(ts *Schema, ti *Info, idxs []schema.Index, store *stor.Stor) {
	if len(idxs) == 0 {
		return
	}
	ts.Indexes = slices.Clip(ts.Indexes) // copy on write
	nold := len(ts.Indexes)
	for i := range idxs {
		ix := &idxs[i]
		if ts.FindIndex(ix.Columns) != nil {
			panic("duplicate index: " +
				str.Join("(,)", ix.Columns) + " in " + ts.Table)
		}
		ts.Indexes = append(ts.Indexes, *ix)
	}
	idxs = ts.SetupNewIndexes(nold)
	n := len(ti.Indexes)
	ti.Indexes = slices.Clip(ti.Indexes) // copy on write
	for i := range idxs {
		bt := btree.CreateBtree(store, &ts.Indexes[n+i].Ixspec)
		ti.Indexes = append(ti.Indexes, index.OverlayFor(bt))
	}
}

func (*Meta) createFkeys(mu *metaUpdate, ts, ac *schema.Schema) {
	idxs := ac.Indexes
	for i := range idxs {
		fk := &idxs[i].Fk
		if fk.Table == "" {
			continue
		}
		tsi := ts.IIndex(idxs[i].Columns)
		fk = &ts.Indexes[tsi].Fk
		fkCols := fk.Columns
		if len(fkCols) == 0 {
			fkCols = idxs[i].Columns
		}
		target := mu.getSchema(fk.Table)
		if target == nil {
			panic("can't create foreign key to nonexistent table: " +
				ac.Table + " -> " + fk.Table)
		}
		found := false
		target.Indexes = slices.Clone(target.Indexes)
		for j := range target.Indexes {
			ix := &target.Indexes[j]
			if slices.Equal(fkCols, ix.Columns) {
				if ix.Mode != 'k' {
					panic("foreign key must point to key: " +
						ac.Table + " -> " + fk.Table + str.Join("(,)", fkCols))
				}
				found = true
				fk.IIndex = j
				ix.FkToHere = slc.With(ix.FkToHere,
					Fkey{Table: ac.Table,
						Columns: idxs[i].Columns, IIndex: tsi, Mode: fk.Mode})
			}
		}
		if !found {
			panic("can't create foreign key to nonexistent index: " +
				ac.Table + " -> " + fk.Table + str.Join("(,)", fkCols))
		}
		mu.putSchema(target)
	}
}

// CheckFkeys checks the targets of the foreign keys of a table.
// It panics on error.
// It is basically a read-only version of createFkeys.
func (m *Meta) CheckFkeys(ts *schema.Schema) {
	idxs := ts.Indexes
	for i := range idxs {
		fk := &idxs[i].Fk
		if fk.Table == "" {
			continue
		}
		tsi := ts.IIndex(idxs[i].Columns)
		fk = &ts.Indexes[tsi].Fk
		fkCols := fk.Columns
		if len(fkCols) == 0 {
			fkCols = idxs[i].Columns
		}
		target, ok := m.schema.Get(fk.Table)
		if !ok {
			panic("can't create foreign key to nonexistent table: " +
				ts.Table + " -> " + fk.Table)
		}
		found := false
		for j := range target.Indexes {
			ix := &target.Indexes[j]
			if slices.Equal(fkCols, ix.Columns) {
				if ix.Mode != 'k' {
					panic("foreign key must point to key: " +
						ts.Table + " -> " + fk.Table + str.Join("(,)", fkCols))
				}
				found = true
			}
		}
		if !found {
			panic("can't create foreign key to nonexistent index: " +
				ts.Table + " -> " + fk.Table + str.Join("(,)", fkCols))
		}
	}
}

func updateFkeysIIndex(mu *metaUpdate, sch *schema.Schema) {
	for i := range sch.Indexes {
		ix := &sch.Indexes[i]
		if ix.Fk.Table != "" {
			updateOtherFkToHere(mu, sch.Table, &ix.Fk, i)
		}
		for j := range ix.FkToHere {
			updateOtherFk(mu, sch.Table, &ix.FkToHere[j], i)
		}
	}
}

func updateOtherFkToHere(mu *metaUpdate, table string, fk *Fkey, iindex int) {
	ts := mu.getSchema(fk.Table)
	ts.Indexes = slices.Clone(ts.Indexes)
	for i := range ts.Indexes {
		ix := &ts.Indexes[i]
		for j := range ix.FkToHere {
			ix.FkToHere = slices.Clone(ix.FkToHere)
			fk2 := &ix.FkToHere[j]
			if fk2.Table == table && slices.Equal(ix.Columns, fk.Columns) {
				fk2.IIndex = iindex
			}
		}
	}
	mu.putSchema(ts)
}

func updateOtherFk(mu *metaUpdate, table string, fk *Fkey, iindex int) {
	ts := mu.getSchema(fk.Table)
	ts.Indexes = slices.Clone(ts.Indexes)
	for i := range ts.Indexes {
		ix := &ts.Indexes[i]
		if ix.Fk.Table == table && slices.Equal(ix.Columns, fk.Columns) {
			ix.Fk.IIndex = iindex
		}
	}
	mu.putSchema(ts)
}

func (m *Meta) AlterDrop(ad *schema.Schema) *Meta {
	ts, ti := m.alterGet(ad.Table)
	// need to drop indexes before columns
	// in case we drop a column and an index that contains it
	dropIndexes(ts, ti, ad.Indexes)
	if !dropColumns(ts, ad) {
		return nil
	}
	mu := newMetaUpdate(m)
	mu.putSchema(ts)
	mu.putInfo(ti)
	m.dropFkeys(mu, ad)
	updateFkeysIIndex(mu, &ts.Schema)
	return mu.freeze()
}

func dropIndexes(ts *Schema, ti *Info, idxs []schema.Index) {
	if len(idxs) == 0 {
		return
	}
	for j := range idxs {
		exists := false
		for i := range ts.Indexes {
			if slices.Equal(ts.Indexes[i].Columns, idxs[j].Columns) {
				if 0 != len(ts.Indexes[i].FkToHere) {
					panic("can't drop index used by foreign keys: " +
						ts.Table + " " + str.Join("(,)", idxs[j].Columns))
				}
				exists = true
			} else if slices.Equal(ts.Indexes[i].BestKey, idxs[j].Columns) {
				panic("can't drop key used to make index unique: " +
					ts.Table + " " + str.Join("(,)", idxs[j].Columns))
			}
		}
		if !exists {
			panic("can't drop nonexistent index: " +
				ts.Table + " " + str.Join("(,)", idxs[j].Columns))
		}
	}
	tsIdxs := make([]schema.Index, 0, len(ts.Indexes))
	tiIdxs := make([]*index.Overlay, 0, len(ti.Indexes))
outer:
	for i := range ts.Indexes {
		for j := range idxs {
			if slices.Equal(ts.Indexes[i].Columns, idxs[j].Columns) {
				continue outer // i.e. don't copy deletion
			}
		}
		tsIdxs = append(tsIdxs, ts.Indexes[i])
		tiIdxs = append(tiIdxs, ti.Indexes[i])
	}
	mustHaveKey(tsIdxs, ts)
	ts.Indexes = tsIdxs
	ts.Ixspecs(len(ts.Indexes)) // need to run setPrimary and setContainsKey
	ti.Indexes = tiIdxs
}

func mustHaveKey(tsIdxs []schema.Index, ts *Schema) {
	for i := range tsIdxs {
		if tsIdxs[i].Mode == 'k' {
			return
		}
	}
	panic("can't drop all keys: " + ts.Table)
}

func dropColumns(ts *Schema, ad *schema.Schema) bool {
	for _, col := range ad.Columns {
		if !dropColumn(ts, col) {
			return false
		}
	}
	for _, col := range ad.Derived {
		if !dropDerived(ts, col) {
			return false
		}
	}
	return true
}

func dropColumn(ts *Schema, col string) bool {
	if inIndex(ts, col) {
		return false // can't drop if used by index
	}
	ucol := str.UnCapitalize(col)
	ccol := str.Capitalize(col)
	if slices.Contains(ts.Columns, ucol) {
		ts.Columns = slc.Replace1(ts.Columns, ucol, "-")
	} else if slices.Contains(ts.Derived, ccol) {
		ts.Derived = slc.Without(ts.Derived, ccol)
	} else {
		panic("can't drop nonexistent column: " + col)
	}
	return true
}

func dropDerived(ts *Schema, col string) bool {
	if !slices.Contains(ts.Derived, col) {
		panic("can't drop nonexistent column: " + col)
	}
	if inIndex(ts, col) {
		return false // can't drop if used by index
	}
	ts.Derived = slc.Without(ts.Derived, col)
	return true
}

func inIndex(ts *Schema, col string) bool {
	for i := range ts.Indexes {
		if slices.Contains(ts.Indexes[i].Columns, col) {
			return true // can't drop if used by index
		}
	}
	return false
}

func (m *Meta) dropFkeys(mu *metaUpdate, drop *schema.Schema) {
	// unlike createFkeys
	// we need to get the actual schema to get the foreign key information
	schema := &m.GetRoSchema(drop.Table).Schema
	idxs := drop.Indexes
	for i := range idxs {
		idx := schema.FindIndex(idxs[i].Columns)
		fk := idx.Fk
		if fk.Table == "" || fk.Table == drop.Table {
			continue
		}
		fkCols := fk.Columns
		if len(fkCols) == 0 {
			fkCols = idx.Columns
		}
		t, ok := mu.schema.Get(fk.Table)
		if !ok {
			log.Println("foreign key: can't find", fk.Table, "(from "+drop.Table+")")
			continue
		}
		target := *t // copy
		target.Indexes = slices.Clone(target.Indexes)
		for j := range target.Indexes {
			ix := &target.Indexes[j]
			if slices.Equal(fkCols, ix.Columns) {
				fk.IIndex = j
				fkToHere := make([]Fkey, 0, len(ix.FkToHere))
				for k := range ix.FkToHere {
					fk2 := &ix.FkToHere[k]
					if drop.Table != fk2.Table ||
						!slices.Equal(idx.Columns, fk2.Columns) {
						fkToHere = append(fkToHere, *fk2)
					}
				}
				ix.FkToHere = fkToHere
			}
		}
		mu.putSchema(&target)
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
	schema := *m.GetRoSchema(table) // copy
	schema.Indexes = slices.Clone(schema.Indexes)
	mu := newMetaUpdate(m)
	mu.putSchema(&schema)
	return mu.freeze()
}

//-------------------------------------------------------------------

// LayeredOnto is called by transaction commit inside UpdateState.
// It layers the mutable ixbuf's from transactions
// onto the latest/current state and returns a new state.
// Also, the nrows and size deltas are applied.
// Note: this does not merge the ixbuf's, that is done later by merge.
// Nor does it save the changes to disk, that is done later by persist.
func (m *Meta) LayeredOnto(latest *Meta) *Meta {
	// start with a snapshot of the latest hash table because it may have more
	assert.That(latest.difInfo == nil)
	info := latest.info.Mutable()
	for _, ti := range m.difInfo {
		tiOrig, _ := m.info.Get(ti.Table)

		lti, ok := info.Get(ti.Table)
		if !ok || lti.IsTomb() {
			continue
		}
		ti.Nrows = lti.Nrows + (ti.Nrows - tiOrig.Nrows)
		assert.That(ti.Nrows >= 0)
		d := int64(ti.Size) - int64(tiOrig.Size)
		ti.Size = uint64(int64(lti.Size) + d)
		for i := range ti.Indexes {
			ti.Indexes[i].UpdateWith(lti.Indexes[i])
		}
		ti.lastMod = m.info.Clock
		info.Put(ti)
	}
	result := *latest // copy
	result.info.Hamt = info.Freeze()
	return &result
}

//-------------------------------------------------------------------

func (m *Meta) Write(store *stor.Stor) (schemaOff uint64, infoOff uint64) {
	assert.That(m.difInfo == nil)
	schemaOff, m.schema = m.schema.WriteChain(store)
	infoOff, m.info = m.info.WriteChain(store)
	return schemaOff, infoOff
}

func ReadMeta(store *stor.Stor, offSchema, offInfo uint64) *Meta {
	m := Meta{
		schema: hamt.ReadChain[string](store, offSchema, ReadSchema),
		info:   hamt.ReadChain[string](store, offInfo, ReadInfo)}
	// copy Ixspec to Info from Schema (constructed by ReadSchema)
	// Ok to modify since it's not in use yet.
	m.info.ForEach(func(ti *Info) {
		if ti.IsTomb() {
			return
		}
		ts := m.schema.MustGet(ti.Table)
		for i := range ti.Indexes {
			ti.Indexes[i].SetIxspec(&ts.Indexes[i].Ixspec)
		}
	})
	linkFkeys(&m)
	return &m
}

type Fkey = schema.Fkey

// linkFkeys links foreign keys to targets (Fk and FkToHere[])
func linkFkeys(m *Meta) {
	m.schema.ForEach(func(s *Schema) {
		if s.IsTomb() {
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
					if slices.Equal(fkCols, ix.Columns) {
						fk.IIndex = j
						ix.FkToHere = append(ix.FkToHere,
							Fkey{Table: s.Table, Mode: fk.Mode,
								Columns: s.Indexes[i].Columns, IIndex: i})
					}
				}
			}
		}
	})
}

//-------------------------------------------------------------------

func (m *Meta) CheckAllMerged() {
	m.info.ForEach(func(ti *Info) {
		for _, ov := range ti.Indexes {
			ov.CheckMerged()
		}
	})
}

func (m *Meta) Offsets() (schemaOff, infoOff uint64) {
	if ns := len(m.schema.Offs); ns > 0 {
		schemaOff = m.schema.Offs[ns-1]
	}
	if ni := len(m.info.Offs); ni > 0 {
		infoOff = m.info.Offs[ni-1]
	}
	return
}
