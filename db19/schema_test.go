// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestTablesNames(t *testing.T) {
	assert := assert.T(t)
	m := &meta.Meta{}
	add := func(name string) {
		m = m.Put(&meta.Schema{Schema: schema.Schema{
			Table: name, Columns: []string{"a"},
			Indexes: []schema.Index{{Mode: 'k', Columns: []string{"a"}}},
		}}, nil)
	}
	add("z")
	add("a")
	// Table names require only Schema, not Info.
	rt := &ReadTran{meta: m}
	assert.This(tablesSchema.Columns).Is([]string{"table"})
	check := func(rec core.Record, name string) {
		assert.This(rec.GetStr(0)).Is(name)
		assert.This(rec.Count()).Is(1)
	}
	st := tablesTable{}
	assert.This(st.nrows(&rt.tran)).Is(2 + len(schemaTables))
	check(st.lookup(rt, 0, core.Pack(core.SuStr("a"))), "a")
	check(st.lookup(rt, 0, core.Pack(core.SuStr("tables"))), "tables")
	assert.This(st.lookup(rt, 0, core.Pack(core.SuStr("missing")))).Is(core.Record(""))
	it := newTablesIter(rt)
	for _, name := range []string{"a", "columns", "dbstats", "history", "indexes", "tables", "views", "z"} {
		it.Next(rt)
		check(it.GetRecord(rt), name)
	}
	it.Next(rt)
	assert.That(it.Eof())
}

func addTable(m *meta.Meta, name string) *meta.Meta {
	return m.Put(&meta.Schema{Schema: schema.Schema{
		Table: name, Columns: []string{"a"},
		Indexes: []schema.Index{{Mode: 'k', Columns: []string{"a"}}},
	}}, nil)
}

func TestSchemaCacheReuse(t *testing.T) {
	assert := assert.T(t)
	db := &Database{}
	m1 := addTable(&meta.Meta{}, "a")

	rt1 := &ReadTran{db: db, meta: m1}
	e1 := rt1.schemaEntry()
	assert.That(e1 != nil)
	assert.That(rt1.schemaEntry() == e1) // captured

	// a new Meta with the same schema root reuses the entry
	rt2 := &ReadTran{db: db, meta: m1.Mutable()}
	assert.That(rt2.schemaEntry() == e1)

	// a schema change gets a new entry
	rt3 := &ReadTran{db: db, meta: addTable(m1, "b")}
	assert.That(rt3.schemaEntry() != e1)
}

func TestSchemaCacheBounded(t *testing.T) {
	assert := assert.T(t)
	db := &Database{}
	m := &meta.Meta{}
	entries := make([]*schemaCacheEntry, 0, schemaCacheSize+2)
	for i := range schemaCacheSize + 2 {
		m = addTable(m, "t"+strconv.Itoa(i))
		rt := &ReadTran{db: db, meta: m}
		entries = append(entries, rt.schemaEntry())
	}
	assert.This(len(db.schemaCache.entries)).Is(schemaCacheSize)
	// most recent retained, oldest evicted
	assert.That(db.schemaCache.entries[0] == entries[len(entries)-1])
	assert.That(!slices.Contains(db.schemaCache.entries, entries[0]))

	retained := slices.Clone(db.schemaCache.entries)
	// Accessing retained old versions does not change eviction order.
	assert.That(db.schemaCache.get(retained[len(retained)-1].meta) == retained[len(retained)-1])
	assert.This(db.schemaCache.entries).Is(retained)
	// An evicted version is usable without being retained again.
	old := &ReadTran{db: db, meta: entries[0].meta}
	assert.This(len(old.tablesItems())).Is(1 + len(schemaTables))
	assert.That(old.schemaEntry() == old.scache)
	assert.This(db.schemaCache.entries).Is(retained)
	assert.This(len(entries[0].tables())).Is(1 + len(schemaTables)) // eviction does not invalidate entries
}

func TestSchemaCacheVersionOrder(t *testing.T) {
	assert := assert.T(t)
	db := &Database{}
	versions := make([]*meta.Meta, schemaCacheSize+2)
	m := &meta.Meta{}
	for i := range versions {
		m = addTable(m, "t"+strconv.Itoa(i))
		versions[i] = m
	}
	// Transactions may first access their schema in a different order.
	for _, i := range []int{2, 5, 0, 4, 1, 3} {
		db.schemaCache.get(versions[i])
	}
	assert.This(len(db.schemaCache.entries)).Is(schemaCacheSize)
	for i, e := range db.schemaCache.entries {
		assert.That(e.meta == versions[len(versions)-1-i])
	}
}

func TestSchemaCacheExpiration(t *testing.T) {
	assert := assert.T(t)
	c := &schemaCache{}
	m1 := addTable(&meta.Meta{}, "a")
	m2 := addTable(m1, "b")
	m3 := addTable(m2, "c")
	old := c.get(m1)
	active := c.get(m2)
	newest := c.get(m3)
	aged := time.Now().Add(-schemaCacheIdle - time.Second)
	old.lastUsed = aged
	newest.lastUsed = aged
	active.lastUsed = time.Now().Add(-schemaCacheIdle / 2)
	previous := active.lastUsed
	assert.That(c.get(m2) == active)
	assert.That(active.lastUsed.After(previous))
	assert.This(len(c.entries)).Is(2)
	assert.That(c.entries[0] == newest) // retained despite its age
	assert.That(c.entries[1] == active)
	assert.This(len(old.tables())).Is(1 + len(schemaTables)) // expiration does not invalidate references
	rebuilt := c.get(m1)
	assert.That(rebuilt != old)
	assert.This(len(rebuilt.tables())).Is(1 + len(schemaTables))
	assert.That(c.entries[0] == newest)
	assert.This(len(c.entries)).Is(3)
	assert.That(c.get(m3) == newest)
	assert.That(newest.lastUsed.After(aged))
}

func TestSchemaCacheIterators(t *testing.T) {
	assert := assert.T(t)
	db := &Database{}
	m := &meta.Meta{}
	for _, name := range []string{"z", "a", "m"} {
		m = addTable(m, name)
	}
	rt := &ReadTran{db: db, meta: m}
	it1 := newTablesIter(rt)
	it2 := newTablesIter(rt)
	it1.Next(rt)
	assert.This(it1.GetRecord(rt).GetStr(0)).Is("a")
	it1.Next(rt)
	assert.This(it1.GetRecord(rt).GetStr(0)).Is("columns")
	// iterators have independent positions
	it2.Next(rt)
	assert.This(it2.GetRecord(rt).GetStr(0)).Is("a")
	assert.This(tablesTable{}.nrows(&rt.tran)).Is(3 + len(schemaTables))
}

func TestSchemaCacheOtherTables(t *testing.T) {
	assert := assert.T(t)
	db := &Database{}
	m := &meta.Meta{}
	m = addTable(m, "t")
	m = m.AddView("v", "t")
	rt := &ReadTran{db: db, meta: m}
	assert.This(columnsTable{}.nrows(&rt.tran)).Is(1 + 18)
	assert.This(indexesTable{}.nrows(&rt.tran)).Is(1 + 6)
	assert.This(viewsTable{}.nrows(&rt.tran)).Is(1)
}

func TestSchemaCacheLazy(t *testing.T) {
	assert := assert.T(t)
	db := &Database{}
	m := &meta.Meta{}
	for _, name := range []string{"a", "b", "c"} {
		m = addTable(m, name)
	}
	m = m.AddView("v", "a")
	rt := &ReadTran{db: db, meta: m}

	tables := rt.tablesItems()
	for i := range tables {
		assert.That(tables[i].cached.Load() == nil)
	}
	it := newTablesIter(rt)
	it.Next(rt)
	assert.This(it.GetRecord(rt).GetStr(0)).Is("a")
	assert.That(tables[0].cached.Load() != nil)
	assert.That(tables[1].cached.Load() == nil)
	assert.That(tables[2].cached.Load() == nil)

	columns := rt.columnsItems()
	assert.That(len(columns) > 0)
	for i := range columns {
		assert.That(columns[i].cached.Load() == nil)
	}
	assert.This(columns[0].row().Count()).Is(3)
	assert.That(columns[0].cached.Load() != nil)

	indexes := rt.indexesItems()
	assert.That(len(indexes) > 0)
	for i := range indexes {
		assert.That(indexes[i].cached.Load() == nil)
	}
	assert.This(indexes[0].row().Count()).Is(3)
	assert.That(indexes[0].cached.Load() != nil)

	views := rt.viewsItems()
	assert.That(views[0].cached.Load() == nil)
	assert.This(views[0].keyAt()).Is(core.Pack(core.SuStr("v")))
	assert.That(views[0].cached.Load() != nil)
}

func TestSchemaCacheAsofReset(t *testing.T) {
	assert := assert.T(t)
	db := CreateDb(stor.HeapStor(8192))
	db.CheckerSync()
	createTbl(db)
	db.persist(&execPersistSingle{}, false)

	rt := db.NewReadTran(core.AllPerms)
	assert.That(rt.schemaEntry() != nil)
	assert.That(rt.scache != nil)
	assert.That(rt.Asof(-1) != 0)
	assert.That(rt.scache == nil) // re-resolved on next access
	retained := slices.Clone(db.schemaCache.entries)
	historical := rt.schemaEntry()
	assert.That(!slices.Contains(db.schemaCache.entries, historical))
	assert.That(rt.schemaEntry() == historical) // still reused within the transaction
	assert.This(db.schemaCache.entries).Is(retained)

	// Even Asof selecting the current state bypasses the shared cache.
	rt.Asof(1 << 62)
	assert.That(rt.meta == db.GetState().Meta)
	current := rt.schemaEntry()
	assert.That(current != historical)
	assert.That(!slices.Contains(db.schemaCache.entries, current))
	assert.This(db.schemaCache.entries).Is(retained)
}
