// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"slices"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stats"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/str"
)

// This file implements the schema tables.
// It intercepts at the transaction - primarily Lookup and IndexIter
// This has the advantage of using the existing [query.Table]
// which means we get the Where-on-Table optimizations.

func IsSchemaTable(table string) bool {
	return getSchemaTable(table) != nil
}

// IsSystemTable includes virtual schema tables and the persisted stats table.
func IsSystemTable(table string) bool {
	return table == StatsTableName || IsSchemaTable(table)
}

type schemaTable interface {
	schema() *schema.Schema
	nrows(*tran) int
	lookup(*ReadTran, int, string) core.Record
	indexIter(*ReadTran, int) index.IndexIter
	rangeFrac(*ReadTran, int, string, string) float64
}

func getSchemaTable(table string) schemaTable {
	switch table {
	case "tables":
		return tablesTable{}
	case "columns":
		return columnsTable{}
	case "indexes":
		return indexesTable{}
	case "views":
		return viewsTable{}
	case "history":
		return historyTable{}
	case "dbstats":
		return statsTable{}
	}
	return nil
}

// The schema-table cache is kept "to the side" of metadata so that
// materializations can be reused across queries and transactions.
// Entries are keyed by the schema key (Meta.SchemaKey), which is
// stable across ordinary data commits. A small number of newest versions is
// retained; transactions and iterators keep older entries alive through Go
// references, so eviction only drops the database's reference.
// Older entries expire lazily after two minutes without a shared-cache hit.
// The newest cached schema is always retained.

const schemaCacheSize = 4
const schemaCacheIdle = 2 * time.Minute

type schemaCache struct {
	mu      sync.Mutex
	entries []*schemaCacheEntry // newest schema version first
}

func (c *schemaCache) get(m *meta.Meta) *schemaCacheEntry {
	key := m.SchemaKey()
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	if len(c.entries) > 1 {
		newest := c.entries[0]
		c.entries = slices.DeleteFunc(c.entries, func(e *schemaCacheEntry) bool {
			return e != newest && now.Sub(e.lastUsed) >= schemaCacheIdle
		})
	}
	if i := slices.IndexFunc(c.entries, func(e *schemaCacheEntry) bool {
		return e.meta.SchemaKey() == key
	}); i >= 0 {
		c.entries[i].lastUsed = now
		return c.entries[i]
	}
	e := newSchemaCacheEntry(m)
	e.lastUsed = now
	i := sort.Search(len(c.entries), func(i int) bool {
		return c.entries[i].meta.SchemaVersion() <= m.SchemaVersion()
	})
	if i >= schemaCacheSize {
		return e // old transactions must not displace newer schemas
	}
	c.entries = slices.Insert(c.entries, i, e)
	if len(c.entries) > schemaCacheSize {
		c.entries = slices.Delete(c.entries, schemaCacheSize, len(c.entries))
	}
	return e
}

// schemaCacheEntry holds lazily-built materializations for one schema version.
// Each materialization is built at most once, independently of the others.
type schemaCacheEntry struct {
	meta     *meta.Meta
	lastUsed time.Time // protected by schemaCache.mu
	tables   func() []tableEntry
	columns  func() []columnEntry
	indexes  func() []indexEntry
	views    func() []viewEntry
}

func newSchemaCacheEntry(m *meta.Meta) *schemaCacheEntry {
	return &schemaCacheEntry{
		meta:    m,
		tables:  sync.OnceValue(func() []tableEntry { return buildTables(m) }),
		columns: sync.OnceValue(func() []columnEntry { return buildColumns(m) }),
		indexes: sync.OnceValue(func() []indexEntry { return buildIndexes(m) }),
		views:   sync.OnceValue(func() []viewEntry { return buildViews(m) }),
	}
}

// schemaEntry resolves (once) the transaction's captured schema cache entry.
// Update transactions share entries because their local writes affect Info,
// not Schema.
func (t *tran) schemaEntry() *schemaCacheEntry {
	if t.scache == nil {
		if t.db != nil && !t.noSchemaCache {
			t.scache = t.db.schemaCache.get(t.meta)
		} else {
			t.scache = newSchemaCacheEntry(t.meta) // Asof or tests
		}
	}
	return t.scache
}

func (t *tran) tablesItems() []tableEntry   { return t.schemaEntry().tables() }
func (t *tran) columnsItems() []columnEntry { return t.schemaEntry().columns() }
func (t *tran) indexesItems() []indexEntry  { return t.schemaEntry().indexes() }
func (t *tran) viewsItems() []viewEntry     { return t.schemaEntry().views() }

var tablesSchema = &schema.Schema{
	Table:   "tables",
	Columns: []string{"table"},
	Indexes: []schema.Index{{Mode: 'k', Columns: []string{"table"}, Fields: []string{"table"}}},
}

var columnsSchema = &schema.Schema{
	Table:   "columns",
	Columns: []string{"table", "column", "field"},
	Indexes: []schema.Index{{Mode: 'k', Columns: []string{"table", "column"},
		Fields: []string{"table", "column"}}},
}

var indexesSchema = &schema.Schema{
	Table:   "indexes",
	Columns: []string{"table", "columns", "key", "fktable", "fkcolumns", "fkmode"},
	Indexes: []schema.Index{{Mode: 'k', Columns: []string{"table", "columns"},
		Fields: []string{"table", "columns"}}},
}

var viewsSchema = &schema.Schema{
	Table:   "views",
	Columns: []string{"view_name", "view_definition"},
	Indexes: []schema.Index{{Mode: 'k', Columns: []string{"view_name"}, Fields: []string{"view_name"}}},
}

var historySchema = &schema.Schema{
	Table:   "history",
	Columns: []string{"asof"},
	Indexes: []schema.Index{{Mode: 'k', Columns: []string{"asof"}, Fields: []string{"asof"}}},
}

var statsSchema = &schema.Schema{
	Table:   "dbstats",
	Columns: []string{"table", "column", "distinct", "quantiles", "common"},
	Indexes: []schema.Index{{Mode: 'k', Columns: []string{"table", "column"},
		Fields: []string{"table", "column"}}},
}

var schemaTables = []*schema.Schema{
	tablesSchema,
	columnsSchema,
	indexesSchema,
	viewsSchema,
	historySchema,
	statsSchema,
}

func schemaForTable(t *ReadTran, table string) *schema.Schema {
	if st := getSchemaTable(table); st != nil {
		return st.schema()
	}
	if ts := t.meta.GetRoSchema(table); ts != nil {
		return &ts.Schema
	}
	return nil
}

// unpackStr decodes a packed field value as a string.
// It returns "" if the value is not a string.
func unpackStr(p string) string {
	if p == "" || p[0] != core.PackString {
		return ""
	}
	return p[1:]
}

// nameGE reports whether name sorts at or after the packed key.
// Non-string keys sort before or after all names.
func nameGE(name, key string) bool {
	if key == "" || key[0] < core.PackString {
		return true
	}
	if key[0] > core.PackString {
		return false
	}
	return name >= unpackStr(key)
}

// tables -----------------------------------------------------------

type tablesTable struct{}

func (tablesTable) schema() *schema.Schema { return tablesSchema }

func (tablesTable) nrows(t *tran) int { return len(t.tablesItems()) }

func (tablesTable) indexIter(t *ReadTran, iIndex int) index.IndexIter {
	assert.That(iIndex == 0)
	return newTablesIter(t)
}

func (tablesTable) rangeFrac(t *ReadTran, iIndex int, org, end string) float64 {
	assert.That(iIndex == 0)
	return newTablesIter(t).rangeFrac(org, end)
}

func (tablesTable) lookup(t *ReadTran, iIndex int, key string) core.Record {
	assert.That(iIndex == 0)
	name := unpackStr(key)
	if name == "" {
		return ""
	}
	if schemaForTable(t, name) == nil {
		return ""
	}
	return tablesRow(name)
}

func tablesRow(table string) core.Record {
	var rb core.RecordBuilder
	rb.Add(core.SuStr(table))
	return rb.Build()
}

func newTablesIter(t *ReadTran) *schemaIter[tableEntry, *tableEntry] {
	items := t.tablesItems()
	return &schemaIter[tableEntry, *tableEntry]{
		items: items,
		name:  "tablesIter",
		hi:    len(items),
	}
}

func buildTables(m *meta.Meta) []tableEntry {
	items := make([]tableEntry, 0, 32)
	for sc := range m.Tables() {
		items = append(items, tableEntry{table: sc.Table})
	}
	for _, sc := range schemaTables {
		items = append(items, tableEntry{table: sc.Table})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].table < items[j].table })
	return items
}

type tableEntry struct {
	table  string
	cached atomic.Pointer[tableRow]
}

// tableRow is built on first access so that narrow schema queries do not
// construct keys and records for every table.
type tableRow struct {
	key string
	rec core.Record
}

func (e *tableEntry) keyAt() string         { return e.load().key }
func (e *tableEntry) geKey(key string) bool { return nameGE(e.table, key) }
func (e *tableEntry) row() core.Record      { return e.load().rec }

func (e *tableEntry) load() *tableRow {
	if row := e.cached.Load(); row != nil {
		return row
	}
	// Rows are deterministic, so concurrent builders may safely overwrite.
	row := &tableRow{
		key: core.Pack(core.SuStr(e.table)),
		rec: tablesRow(e.table),
	}
	e.cached.Store(row)
	return row
}

// columns ----------------------------------------------------------

type columnsTable struct{}

func (columnsTable) schema() *schema.Schema { return columnsSchema }

func (columnsTable) nrows(t *tran) int { return len(t.columnsItems()) }

func (columnsTable) indexIter(t *ReadTran, iIndex int) index.IndexIter {
	assert.That(iIndex == 0)
	return newColumnsIter(t)
}

func (columnsTable) rangeFrac(t *ReadTran, iIndex int, org, end string) float64 {
	assert.That(iIndex == 0)
	return newColumnsIter(t).rangeFrac(org, end)
}

func (columnsTable) lookup(t *ReadTran, iIndex int, key string) core.Record {
	assert.That(iIndex == 0)
	// the {table, column} key is encoded (see ixkey)
	parts := ixkey.Decode(key)
	if len(parts) != 2 {
		return ""
	}
	table := unpackStr(parts[0])
	col := unpackStr(parts[1])
	if table == "" || col == "" {
		return ""
	}
	ts := schemaForTable(t, table)
	if ts == nil || col == "-" { // "-" columns are not in the schema table
		return ""
	}
	fld := -1
	if i := slices.Index(ts.Columns, col); i >= 0 {
		fld = i
	} else if !slices.ContainsFunc(ts.Derived,
		// Derived are stored capitalized but reported uncapitalized
		func(d string) bool { return str.EqualUncap(d, col) }) {
		return ""
	}
	return columnsRow(table, col, fld)
}

func columnsRow(table, column string, field int) core.Record {
	var rb core.RecordBuilder
	rb.Add(core.SuStr(table))
	rb.Add(core.SuStr(column))
	rb.Add(core.IntVal(field))
	return rb.Build()
}

// columnEntry is one row of the columns table with its encoded {table, column} key.
// The key is eager because it is needed for sorting and range bounds.
type columnEntry struct {
	table  string
	column string
	field  int
	key    string
	cached atomic.Pointer[core.Record]
}

func (e *columnEntry) keyAt() string         { return e.key }
func (e *columnEntry) geKey(key string) bool { return e.key >= key }
func (e *columnEntry) row() core.Record {
	if rec := e.cached.Load(); rec != nil {
		return *rec
	}
	rec := columnsRow(e.table, e.column, e.field)
	e.cached.Store(&rec)
	return rec
}

func newColumnsIter(t *ReadTran) *schemaIter[columnEntry, *columnEntry] {
	entries := t.columnsItems()
	return &schemaIter[columnEntry, *columnEntry]{
		items:  entries,
		name:   "columnsIter",
		skipOK: true,
		hi:     len(entries),
	}
}

func buildColumns(m *meta.Meta) []columnEntry {
	entries := make([]columnEntry, 0, 32)
	add := func(sc *schema.Schema) {
		for field, col := range sc.Columns {
			if col == "-" {
				continue
			}
			entries = append(entries, columnEntry{
				table: sc.Table, column: col, field: field,
				key: ixkey.CompKey(core.Pack(core.SuStr(sc.Table)),
					core.Pack(core.SuStr(col))),
			})
		}
		for _, derived := range sc.Derived {
			// Derived are stored capitalized but reported uncapitalized
			col := str.Uncapitalize(derived)
			entries = append(entries, columnEntry{
				table: sc.Table, column: col, field: -1,
				key: ixkey.CompKey(core.Pack(core.SuStr(sc.Table)),
					core.Pack(core.SuStr(col))),
			})
		}
	}
	for sc := range m.Tables() {
		add(&sc.Schema)
	}
	for _, sc := range schemaTables {
		add(sc)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].key < entries[j].key })
	return entries
}

// indexes ----------------------------------------------------------

type indexesTable struct{}

func (indexesTable) schema() *schema.Schema { return indexesSchema }

func (indexesTable) nrows(t *tran) int { return len(t.indexesItems()) }

func (indexesTable) indexIter(t *ReadTran, iIndex int) index.IndexIter {
	assert.That(iIndex == 0)
	return newIndexesIter(t)
}

func (indexesTable) rangeFrac(t *ReadTran, iIndex int, org, end string) float64 {
	assert.That(iIndex == 0)
	return newIndexesIter(t).rangeFrac(org, end)
}

func (indexesTable) lookup(t *ReadTran, iIndex int, key string) core.Record {
	assert.That(iIndex == 0)
	// the {table, columns} key is encoded (see ixkey)
	parts := ixkey.Decode(key)
	if len(parts) == 0 || len(parts) > 2 {
		return ""
	}
	table := unpackStr(parts[0])
	if table == "" {
		return ""
	}
	cols := "" // a one-part key means an empty columns (e.g. key())
	if len(parts) == 2 {
		if parts[1] == "" || parts[1][0] != core.PackString {
			return "" // non-string columns
		}
		cols = parts[1][1:]
	}
	ts := schemaForTable(t, table)
	if ts == nil {
		return ""
	}
	idx := ts.FindIndex(str.Split(cols, ","))
	if idx == nil {
		return ""
	}
	return indexesRow(table, idx)
}

func indexesRow(table string, idx *schema.Index) core.Record {
	var rb core.RecordBuilder
	rb.Add(core.SuStr(table))
	rb.Add(core.SuStr(str.Join(",", idx.Columns)))
	switch idx.Mode {
	case 'k':
		rb.Add(core.True)
	case 'i':
		rb.Add(core.False)
	case 'u':
		rb.Add(core.SuStr("u"))
	default:
		assert.ShouldNotReachHere()
	}
	if idx.Fk.Table != "" {
		rb.Add(core.SuStr(idx.Fk.Table))
		rb.Add(core.SuStr(str.Join(",", idx.Fk.Columns)))
		rb.Add(core.SuInt16(int(idx.Fk.Mode)))
	}
	return rb.Build()
}

// indexEntry is one row of the indexes table with its encoded {table, columns} key.
// The key is eager because it is needed for sorting and range bounds.
type indexEntry struct {
	table  string
	idx    *schema.Index
	key    string
	cached atomic.Pointer[core.Record]
}

func (e *indexEntry) keyAt() string         { return e.key }
func (e *indexEntry) geKey(key string) bool { return e.key >= key }
func (e *indexEntry) row() core.Record {
	if rec := e.cached.Load(); rec != nil {
		return *rec
	}
	rec := indexesRow(e.table, e.idx)
	e.cached.Store(&rec)
	return rec
}

func newIndexesIter(t *ReadTran) *schemaIter[indexEntry, *indexEntry] {
	entries := t.indexesItems()
	return &schemaIter[indexEntry, *indexEntry]{
		items:  entries,
		name:   "indexesIter",
		skipOK: true,
		hi:     len(entries),
	}
}

func buildIndexes(m *meta.Meta) []indexEntry {
	entries := make([]indexEntry, 0, 32)
	add := func(sc *schema.Schema) {
		for i := range sc.Indexes {
			idx := &sc.Indexes[i]
			entries = append(entries, indexEntry{
				table: sc.Table, idx: idx,
				key: ixkey.CompKey(core.Pack(core.SuStr(sc.Table)),
					core.Pack(core.SuStr(str.Join(",", idx.Columns)))),
			})
		}
	}
	for sc := range m.Tables() {
		add(&sc.Schema)
	}
	for _, sc := range schemaTables {
		add(sc)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].key < entries[j].key })
	return entries
}

// views ------------------------------------------------------------

type viewsTable struct{}

func (viewsTable) schema() *schema.Schema { return viewsSchema }

func (viewsTable) nrows(t *tran) int { return len(t.viewsItems()) }

func (viewsTable) indexIter(t *ReadTran, iIndex int) index.IndexIter {
	assert.That(iIndex == 0)
	return newViewsIter(t)
}

func (viewsTable) rangeFrac(t *ReadTran, iIndex int, org, end string) float64 {
	assert.That(iIndex == 0)
	return newViewsIter(t).rangeFrac(org, end)
}

func (viewsTable) lookup(t *ReadTran, iIndex int, key string) core.Record {
	assert.That(iIndex == 0)
	name := unpackStr(key)
	if name == "" {
		return ""
	}
	def := t.meta.GetView(name)
	if def == "" {
		return ""
	}
	return viewsRow(name, def)
}

func viewsRow(name, def string) core.Record {
	var rb core.RecordBuilder
	rb.Add(core.SuStr(name))
	rb.Add(core.SuStr(def))
	return rb.Build()
}

type viewEntry struct {
	name   string
	def    string
	cached atomic.Pointer[viewRow]
}

type viewRow struct {
	key string
	rec core.Record
}

func (v *viewEntry) keyAt() string         { return v.load().key }
func (v *viewEntry) geKey(key string) bool { return nameGE(v.name, key) }
func (v *viewEntry) row() core.Record {
	return v.load().rec
}

func (v *viewEntry) load() *viewRow {
	if row := v.cached.Load(); row != nil {
		return row
	}
	row := &viewRow{
		key: core.Pack(core.SuStr(v.name)),
		rec: viewsRow(v.name, v.def),
	}
	v.cached.Store(row)
	return row
}

func newViewsIter(t *ReadTran) *schemaIter[viewEntry, *viewEntry] {
	views := t.viewsItems()
	return &schemaIter[viewEntry, *viewEntry]{
		items: views,
		name:  "viewsIter",
		hi:    len(views),
	}
}

func buildViews(m *meta.Meta) []viewEntry {
	views := make([]viewEntry, 0, 16)
	for name, def := range m.Views() {
		views = append(views, viewEntry{name: name, def: def})
	}
	sort.Slice(views, func(i, j int) bool { return views[i].name < views[j].name })
	return views
}

// history ----------------------------------------------------------

type historyTable struct{}

func (historyTable) schema() *schema.Schema { return historySchema }
func (historyTable) nrows(*tran) int        { return historyNrows }

// historyNrows is a constant estimate; counting database states would be too slow.
const historyNrows = 1000

func (historyTable) lookup(*ReadTran, int, string) core.Record {
	panic("history lookup not implemented")
}

func (historyTable) indexIter(t *ReadTran, iIndex int) index.IndexIter {
	assert.That(iIndex == 0)
	return newHistoryIter(t)
}

// Enumerating database states would be too slow, so use a rough estimate.
func (historyTable) rangeFrac(*ReadTran, int, string, string) float64 {
	return .5
}

// historyIter is an index.IndexIter over the history table.
// It scans database states directly instead of materializing a list,
// since a database may have a very large number of states.
type historyIter struct {
	store *stor.Stor
	off   uint64 // offset of the current state; 0 when rewound
	key   string // current packed asof key
	asof  int64
	org   string
	end   string
	state iface.State
}

func newHistoryIter(t *ReadTran) *historyIter {
	return &historyIter{store: t.GetStore(),
		org: ixkey.Min, end: ixkey.Max, state: iface.Rewound}
}

func (it *historyIter) Eof() bool    { return it.state.Eof() }
func (it *historyIter) HasCur() bool { return it.state.Within() }

func (it *historyIter) Cur() (string, uint64) {
	it.checkHasCur()
	return it.key, it.off
}

func (it *historyIter) CurOff() uint64 {
	it.checkHasCur()
	return it.off
}

func (it *historyIter) GetRecord(index.OiTran) core.Record {
	it.checkHasCur()
	var rb core.RecordBuilder
	rb.Add(core.SuDateFromUnixMilli(int(it.asof)))
	return rb.Build()
}

func (it *historyIter) checkHasCur() {
	if it.state.Eof() {
		panic("historyIter eof")
	}
	if it.state.Rewound() {
		panic("historyIter rewound")
	}
}

func (it *historyIter) setCur(state *DbState, key string) {
	it.off = state.Off
	it.key = key
	it.asof = state.Asof
	it.state = iface.Within
}

func (it *historyIter) Next(index.OiTran) {
	if it.state.Eof() {
		return // stick at eof
	}
	off := it.off
	if it.state.Rewound() {
		off = 0
	}
	for {
		state := NextState(it.store, off)
		if state == nil {
			it.state = iface.Eof
			return
		}
		key := historyKey(state)
		if key < it.org {
			off = state.Off
			continue
		}
		if key >= it.end {
			it.state = iface.Eof
			return
		}
		it.setCur(state, key)
		return
	}
}

func (it *historyIter) Prev(index.OiTran) {
	if it.state.Eof() {
		return // stick at eof
	}
	off := it.off
	if it.state.Rewound() {
		off = 0
	}
	for {
		state := PrevState(it.store, off)
		if state == nil {
			it.state = iface.Eof
			return
		}
		key := historyKey(state)
		if key >= it.end {
			off = state.Off
			continue
		}
		if key < it.org {
			it.state = iface.Eof
			return
		}
		it.setCur(state, key)
		return
	}
}

func historyKey(state *DbState) string {
	return core.Pack(core.SuDateFromUnixMilli(int(state.Asof)))
}

func (it *historyIter) Rewind() {
	it.state = iface.Rewound
	it.off = 0
	it.key = ""
	it.asof = 0
}

func (it *historyIter) Range(rng index.Range) {
	it.org = rng.Org
	it.end = rng.End
	it.Rewind()
}

func (it *historyIter) SkipScan(prefixRng, _ index.Range, _ int) {
	assert.ShouldNotReachHere()
}

// dbstats ----------------------------------------------------------

type statsTable struct{}

func (statsTable) schema() *schema.Schema { return statsSchema }

func (statsTable) nrows(t *tran) (n int) {
	for _, ts := range t.Stats() {
		n += len(ts.Columns)
	}
	return n
}

func (statsTable) indexIter(t *ReadTran, iIndex int) index.IndexIter {
	assert.That(iIndex == 0)
	return newStatsIter(t)
}

func (statsTable) rangeFrac(t *ReadTran, iIndex int, org, end string) float64 {
	assert.That(iIndex == 0)
	return newStatsIter(t).rangeFrac(org, end)
}

func (statsTable) lookup(t *ReadTran, iIndex int, key string) core.Record {
	assert.That(iIndex == 0)
	parts := ixkey.Decode(key)
	if len(parts) != 2 {
		return ""
	}
	table := unpackStr(parts[0])
	column := unpackStr(parts[1])
	if table == "" || column == "" {
		return ""
	}
	if ts, ok := t.Stats()[table]; ok {
		if cs, ok := ts.Columns[column]; ok {
			return statsRow(table, column, cs)
		}
	}
	return ""
}

func statsRow(table, column string, cs stats.ColStats) core.Record {
	var rb core.RecordBuilder
	rb.Add(core.SuStr(table))
	rb.Add(core.SuStr(column))
	rb.Add(core.IntVal(cs.Cardinality))
	rb.Add(formatQuantiles(cs.Quantiles))
	rb.Add(formatTops(cs.Tops))
	return rb.Build()
}

func formatQuantiles(quantiles []string) *core.SuObject {
	parts := &core.SuObject{}
	for _, q := range quantiles {
		parts.Add(core.Unpack(q))
	}
	return parts
}

func formatTops(tops []stats.Top) *core.SuObject {
	parts := &core.SuObject{}
	for _, t := range tops {
		part := &core.SuObject{}
		part.Add(core.SuDnum{Dnum: dnum.FromFloat(t.Frac)})
		part.Add(core.Unpack(t.Value))
		parts.Add(part)
	}
	return parts
}

type statsEntry struct {
	table  string
	column string
	key    string
	cs     stats.ColStats
}

func (e statsEntry) keyAt() string         { return e.key }
func (e statsEntry) geKey(key string) bool { return e.key >= key }
func (e statsEntry) row() core.Record {
	return statsRow(e.table, e.column, e.cs)
}

func newStatsIter(t *ReadTran) *schemaIter[statsEntry, *statsEntry] {
	entries := make([]statsEntry, 0, 32)
	for table, ts := range t.Stats() {
		for column, cs := range ts.Columns {
			entries = append(entries, statsEntry{
				table: table, column: column, cs: cs,
				key: ixkey.CompKey(core.Pack(core.SuStr(table)),
					core.Pack(core.SuStr(column))),
			})
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].key < entries[j].key })
	return &schemaIter[statsEntry, *statsEntry]{
		items:  entries,
		name:   "statsIter",
		skipOK: true,
		hi:     len(entries),
	}
}

//-------------------------------------------------------------------

// schemaSkipScan filters out keys whose suffix does not match the suffix range.
type schemaSkipScan struct {
	rng   index.Range
	start int
}

func (ss schemaSkipScan) matches(key string) bool {
	if ss.start == 0 {
		return true
	}
	_, suffix := ixkey.SplitPrefixSuffix(key, ss.start)
	return ss.rng.Org <= suffix && suffix < ss.rng.End
}

// schemaEntry is the interface implemented by schemaIter items.
type schemaEntry interface {
	keyAt() string         // encoded key for Cur
	geKey(key string) bool // entry sorts at or after the range key
	row() core.Record
}

// schemaIter is an index.IndexIter over a materialized, sorted slice of
// schema table rows. It uses row ordinals instead of storage offsets.
// PT is the pointer type for T; entries use pointer methods so they can
// publish lazy values atomically.
type schemaIter[T any, PT interface {
	*T
	schemaEntry
}] struct {
	items  []T
	name   string
	skipOK bool
	skip   schemaSkipScan
	lo     int
	hi     int
	i      int
	state  iface.State
}

func (it *schemaIter[T, PT]) entry(i int) PT {
	return PT(&it.items[i])
}

func (it *schemaIter[T, PT]) Eof() bool    { return it.state.Eof() }
func (it *schemaIter[T, PT]) HasCur() bool { return it.state.Within() }

func (it *schemaIter[T, PT]) Cur() (string, uint64) {
	it.checkHasCur()
	return it.entry(it.i).keyAt(), uint64(it.i) + 1
}

func (it *schemaIter[T, PT]) CurOff() uint64 {
	it.checkHasCur()
	return uint64(it.i) + 1
}

func (it *schemaIter[T, PT]) GetRecord(index.OiTran) core.Record {
	it.checkHasCur()
	return it.entry(it.i).row()
}

func (it *schemaIter[T, PT]) checkHasCur() {
	if it.state.Eof() {
		panic(it.name + " eof")
	}
	if it.state.Rewound() {
		panic(it.name + " rewound")
	}
}

func (it *schemaIter[T, PT]) Next(index.OiTran) { it.step(+1) }
func (it *schemaIter[T, PT]) Prev(index.OiTran) { it.step(-1) }

// step moves one row in direction dir (+1 or -1), skipping rows that
// do not match the skip scan, and sets Eof when leaving [lo, hi).
func (it *schemaIter[T, PT]) step(dir int) {
	switch it.state {
	case iface.Eof:
		return // stick at eof
	case iface.Rewound:
		if dir > 0 {
			it.i = it.lo
		} else {
			it.i = it.hi - 1
		}
	default:
		it.i += dir
	}
	for it.i >= it.lo && it.i < it.hi &&
		!it.skip.matches(it.entry(it.i).keyAt()) {
		it.i += dir
	}
	if it.i < it.lo || it.i >= it.hi {
		it.state = iface.Eof
		return
	}
	it.state = iface.Within
}

func (it *schemaIter[T, PT]) Rewind() {
	it.state = iface.Rewound
	it.i = 0
}

// bounds returns the [lo, hi) positions matching the key range [org, end).
func (it *schemaIter[T, PT]) bounds(org, end string) (int, int) {
	n := len(it.items)
	at := func(key string) int {
		return sort.Search(n, func(i int) bool { return it.entry(i).geKey(key) })
	}
	return at(org), at(end)
}

// rangeFrac computes the fraction from the rows, since schema tables have no btrees.
func (it *schemaIter[T, PT]) rangeFrac(org, end string) float64 {
	lo, hi := it.bounds(org, end)
	if len(it.items) == 0 || hi <= lo {
		return 0
	}
	return float64(hi-lo) / float64(len(it.items))
}

func (it *schemaIter[T, PT]) Range(rng index.Range) {
	it.lo, it.hi = it.bounds(rng.Org, rng.End)
	it.skip = schemaSkipScan{}
	it.state = iface.Rewound
	it.i = 0
}

func (it *schemaIter[T, PT]) SkipScan(prefixRng, suffixRng index.Range, skipStart int) {
	if !it.skipOK {
		assert.ShouldNotReachHere()
	}
	assert.That(skipStart > 0)
	it.Range(prefixRng)
	it.skip = schemaSkipScan{rng: suffixRng, start: skipStart}
}
