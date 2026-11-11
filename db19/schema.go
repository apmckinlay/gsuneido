// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"slices"
	"sort"

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

// This is a new implementation of the schema tables.
// Instead of implementing the Query interface,
// it intercepts at the transaction - primarily Lookup and IndexIter
// This has the advantage of using the existing [query.Table]
// which means we get the Where-on-Table optimizations.
// The "_" suffix names are temporary during development and testing.

func IsSchemaTable(table string) bool {
	return schemaTableSchema(table) != nil
}

var tablesSchema = &schema.Schema{
	Table:   "tables_",
	Columns: []string{"table", "nrows", "totalsize"},
	Indexes: []schema.Index{{Mode: 'k', Columns: []string{"table"}, Fields: []string{"table"}}},
}

var columnsSchema = &schema.Schema{
	Table:   "columns_",
	Columns: []string{"table", "column", "field"},
	Indexes: []schema.Index{{Mode: 'k', Columns: []string{"table", "column"},
		Fields: []string{"table", "column"}}},
}

var indexesSchema = &schema.Schema{
	Table:   "indexes_",
	Columns: []string{"table", "columns", "key", "fktable", "fkcolumns", "fkmode"},
	Indexes: []schema.Index{{Mode: 'k', Columns: []string{"table", "columns"},
		Fields: []string{"table", "columns"}}},
}

var viewsSchema = &schema.Schema{
	Table:   "views_",
	Columns: []string{"view_name", "view_definition"},
	Indexes: []schema.Index{{Mode: 'k', Columns: []string{"view_name"}, Fields: []string{"view_name"}}},
}

var historySchema = &schema.Schema{
	Table:   "history_",
	Columns: []string{"asof"},
	Indexes: []schema.Index{{Mode: 'k', Columns: []string{"asof"}, Fields: []string{"asof"}}},
}

var statsSchema = &schema.Schema{
	Table:   "dbstats_",
	Columns: []string{"table", "column", "distinct", "quantiles", "common"},
	Indexes: []schema.Index{{Mode: 'k', Columns: []string{"table", "column"},
		Fields: []string{"table", "column"}}},
}

// historyNrows is a constant estimate; counting database states would be too slow.
const historyNrows = 1000

// schemaTableInfo returns the Info for a "_" schema table, or nil if not one.
func (t *tran) schemaTableInfo(table string) *meta.Info {
	nrows := 0
	switch table {
	case "tables_":
		for range t.meta.Infos() {
			nrows++
		}
	case "columns_":
		for sc := range t.meta.Tables() {
			for _, col := range sc.Columns {
				if col != "-" {
					nrows++
				}
			}
			nrows += len(sc.Derived)
		}
	case "indexes_":
		for sc := range t.meta.Tables() {
			nrows += len(sc.Indexes)
		}
	case "views_":
		for range t.meta.Views() {
			nrows++
		}
	case "history_":
		nrows = historyNrows
	case "dbstats_":
		for _, ts := range t.Stats() {
			nrows += len(ts.Columns)
		}
	default:
		return nil
	}
	return &meta.Info{Table: table, Nrows: nrows, Size: nrows * 32}
}

// schemaTableSchema returns the Schema for a "_" schema table, or nil if not one.
func schemaTableSchema(table string) *schema.Schema {
	switch table {
	case "tables_":
		return tablesSchema
	case "columns_":
		return columnsSchema
	case "indexes_":
		return indexesSchema
	case "views_":
		return viewsSchema
	case "history_":
		return historySchema
	case "dbstats_":
		return statsSchema
	}
	return nil
}

// schemaTableLookup performs a key lookup on a "_" schema table.
// It returns the Record and true if table is such a table,
// or "" and false otherwise.
func (t *ReadTran) schemaTableLookup(table string, iIndex int, key string) (core.Record, bool) {
	switch table {
	case "tables_":
		return t.tablesLookup(iIndex, key), true
	case "columns_":
		return t.columnsLookup(iIndex, key), true
	case "indexes_":
		return t.indexesLookup(iIndex, key), true
	case "views_":
		return t.viewsLookup(iIndex, key), true
	case "history_":
		panic("history_ lookup not implemented")
	case "dbstats_":
		return t.statsLookup(iIndex, key), true
	}
	return "", false
}

// SchemaTableAlias maps a "_" schema table name to the name of the schema
// table it represents, so it reuses that table's permissions.
// TEMPORARY to keep new "_" tables separate.
func SchemaTableAlias(table string) string {
	switch table {
	case "tables_":
		return "tables"
	case "columns_":
		return "columns"
	case "indexes_":
		return "indexes"
	case "views_":
		return "views"
	case "history_":
		return "history"
	case "dbstats_":
		return "dbstats"
	}
	return table
}

// unpackStr decodes a packed field value as a string.
// It returns "" if the value is not a string.
func unpackStr(p string) string {
	if p == "" || p[0] != core.PackString {
		return ""
	}
	return p[1:]
}

// infoEntry adapts *meta.Info to the schemaEntry interface.
type infoEntry struct{ *meta.Info }

func (e infoEntry) keyAt() string   { return core.Pack(core.SuStr(e.Table)) }
func (e infoEntry) boundAt() string { return e.Table }
func (e infoEntry) row() core.Record {
	return tablesRow(e.Info)
}

// schemaTableIndexIter returns an IndexIter for a "_" schema table.
func (t *ReadTran) schemaTableIndexIter(table string, iIndex int) index.IndexIter {
	switch table {
	case "tables_":
		assert.That(iIndex == 0)
		return newTablesIter(t)
	case "views_":
		assert.That(iIndex == 0)
		return newViewsIter(t)
	case "columns_":
		assert.That(iIndex == 0)
		return newColumnsIter(t)
	case "indexes_":
		assert.That(iIndex == 0)
		return newIndexesIter(t)
	case "history_":
		assert.That(iIndex == 0)
		return newHistoryIter(t)
	case "dbstats_":
		assert.That(iIndex == 0)
		return newStatsIter(t)
	}
	return nil
}

// schemaTableRangeFrac returns the estimated fraction of rows in the range
// for a "_" schema table, or false if table is not one.
// The schema tables have no btrees, so the fraction is computed from the
// rows themselves. history_ is an exception: enumerating database states
// would be too slow, so it returns a rough estimate.
func (t *ReadTran) schemaTableRangeFrac(table string, iIndex int, org, end string) (float64, bool) {
	switch table {
	case "tables_":
		assert.That(iIndex == 0)
		it := newTablesIter(t)
		lo, hi := it.bounds(org, end)
		return rangeFrac(lo, hi, len(it.items)), true
	case "columns_":
		assert.That(iIndex == 0)
		it := newColumnsIter(t)
		lo, hi := it.bounds(org, end)
		return rangeFrac(lo, hi, len(it.items)), true
	case "indexes_":
		assert.That(iIndex == 0)
		it := newIndexesIter(t)
		lo, hi := it.bounds(org, end)
		return rangeFrac(lo, hi, len(it.items)), true
	case "views_":
		assert.That(iIndex == 0)
		it := newViewsIter(t)
		lo, hi := it.bounds(org, end)
		return rangeFrac(lo, hi, len(it.items)), true
	case "history_":
		return .5, true
	case "dbstats_":
		assert.That(iIndex == 0)
		it := newStatsIter(t)
		lo, hi := it.bounds(org, end)
		return rangeFrac(lo, hi, len(it.items)), true
	}
	return 0, false
}

// rangeFrac returns the fraction of n rows in the [lo, hi) range.
func rangeFrac(lo, hi, n int) float64 {
	if n == 0 || hi <= lo {
		return 0
	}
	return float64(hi-lo) / float64(n)
}

// nameBounds returns the [lo, hi) positions of the names matching the
// packed key range [org, end) within a sorted list of names.
// Non-string bounds sort before or after all names.
func nameBounds(n int, nameAt func(int) string, org, end string) (int, int) {
	return nameBound(n, nameAt, org), nameBound(n, nameAt, end)
}

func nameBound(n int, nameAt func(int) string, key string) int {
	if key == "" || key[0] < core.PackString {
		return 0
	}
	if key[0] > core.PackString {
		return n
	}
	name := unpackStr(key)
	return sort.Search(n, func(i int) bool { return nameAt(i) >= name })
}

// keyBounds returns the [lo, hi) positions of the encoded keys matching the
// key range [org, end) within a sorted list of keys.
func keyBounds(n int, keyAt func(int) string, org, end string) (int, int) {
	bound := func(key string) int {
		return sort.Search(n, func(i int) bool { return keyAt(i) >= key })
	}
	return bound(org), bound(end)
}

// tables -----------------------------------------------------------

func (t *ReadTran) tablesLookup(iIndex int, key string) core.Record {
	assert.That(iIndex == 0)
	name := unpackStr(key)
	if name == "" {
		return ""
	}
	info := t.meta.GetRoInfo(name)
	if info == nil {
		return ""
	}
	return tablesRow(info)
}

func tablesRow(info *meta.Info) core.Record {
	var rb core.RecordBuilder
	rb.Add(core.SuStr(info.Table))
	rb.Add(core.IntVal(info.Nrows))
	rb.Add(core.IntVal(info.Size))
	return rb.Build()
}

func newTablesIter(t *ReadTran) *schemaIter[infoEntry] {
	infos := t.GetAllInfo()
	items := make([]infoEntry, len(infos))
	for i, info := range infos {
		items[i] = infoEntry{info}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Table < items[j].Table })
	return &schemaIter[infoEntry]{
		items: items,
		bound: nameBounds,
		name:  "tablesIter",
		hi:    len(items),
	}
}

// columns ----------------------------------------------------------

func (t *ReadTran) columnsLookup(iIndex int, key string) core.Record {
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
	ts := t.meta.GetRoSchema(table)
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

// columnEntry is one row of the columns_ table with its encoded {table, column} key.
type columnEntry struct {
	table  string
	column string
	field  int
	key    string
}

func (e columnEntry) keyAt() string   { return e.key }
func (e columnEntry) boundAt() string { return e.key }
func (e columnEntry) row() core.Record {
	return columnsRow(e.table, e.column, e.field)
}

func newColumnsIter(t *ReadTran) *schemaIter[columnEntry] {
	entries := make([]columnEntry, 0, 32)
	for sc := range t.meta.Tables() {
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
	sort.Slice(entries, func(i, j int) bool { return entries[i].key < entries[j].key })
	return &schemaIter[columnEntry]{
		items:  entries,
		bound:  keyBounds,
		name:   "columnsIter",
		skipOK: true,
		hi:     len(entries),
	}
}

// indexes ----------------------------------------------------------

func (t *ReadTran) indexesLookup(iIndex int, key string) core.Record {
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
	ts := t.meta.GetRoSchema(table)
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

// indexEntry is one row of the indexes_ table with its encoded {table, columns} key.
type indexEntry struct {
	table string
	idx   *schema.Index
	key   string
}

func (e indexEntry) keyAt() string   { return e.key }
func (e indexEntry) boundAt() string { return e.key }
func (e indexEntry) row() core.Record {
	return indexesRow(e.table, e.idx)
}

func newIndexesIter(t *ReadTran) *schemaIter[indexEntry] {
	entries := make([]indexEntry, 0, 32)
	for sc := range t.meta.Tables() {
		for i := range sc.Indexes {
			idx := &sc.Indexes[i]
			entries = append(entries, indexEntry{
				table: sc.Table, idx: idx,
				key: ixkey.CompKey(core.Pack(core.SuStr(sc.Table)),
					core.Pack(core.SuStr(str.Join(",", idx.Columns)))),
			})
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].key < entries[j].key })
	return &schemaIter[indexEntry]{
		items:  entries,
		bound:  keyBounds,
		name:   "indexesIter",
		skipOK: true,
		hi:     len(entries),
	}
}

// views ------------------------------------------------------------

func (t *ReadTran) viewsLookup(iIndex int, key string) core.Record {
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
	name string
	def  string
}

func (v viewEntry) keyAt() string   { return core.Pack(core.SuStr(v.name)) }
func (v viewEntry) boundAt() string { return v.name }
func (v viewEntry) row() core.Record {
	return viewsRow(v.name, v.def)
}

func newViewsIter(t *ReadTran) *schemaIter[viewEntry] {
	views := make([]viewEntry, 0, 16)
	for name, def := range t.Views() {
		views = append(views, viewEntry{name, def})
	}
	sort.Slice(views, func(i, j int) bool { return views[i].name < views[j].name })
	return &schemaIter[viewEntry]{
		items: views,
		bound: nameBounds,
		name:  "viewsIter",
		hi:    len(views),
	}
}

// history ----------------------------------------------------------

// historyIter is an index.IndexIter over the history_ table.
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

func (t *ReadTran) statsLookup(iIndex int, key string) core.Record {
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

func (e statsEntry) keyAt() string   { return e.key }
func (e statsEntry) boundAt() string { return e.key }
func (e statsEntry) row() core.Record {
	return statsRow(e.table, e.column, e.cs)
}

func newStatsIter(t *ReadTran) *schemaIter[statsEntry] {
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
	return &schemaIter[statsEntry]{
		items:  entries,
		bound:  keyBounds,
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
	keyAt() string   // encoded key for Cur
	boundAt() string // value compared against range bounds
	row() core.Record
}

// schemaIter is an index.IndexIter over a materialized, sorted slice of
// schema table rows. It uses row ordinals instead of storage offsets.
type schemaIter[T schemaEntry] struct {
	items  []T
	bound  func(int, func(int) string, string, string) (int, int)
	name   string
	skipOK bool
	skip   schemaSkipScan
	lo     int
	hi     int
	i      int
	state  iface.State
}

func (it *schemaIter[T]) Eof() bool    { return it.state.Eof() }
func (it *schemaIter[T]) HasCur() bool { return it.state.Within() }

func (it *schemaIter[T]) Cur() (string, uint64) {
	it.checkHasCur()
	return it.items[it.i].keyAt(), uint64(it.i) + 1
}

func (it *schemaIter[T]) CurOff() uint64 {
	it.checkHasCur()
	return uint64(it.i) + 1
}

func (it *schemaIter[T]) GetRecord(index.OiTran) core.Record {
	it.checkHasCur()
	return it.items[it.i].row()
}

func (it *schemaIter[T]) checkHasCur() {
	if it.state.Eof() {
		panic(it.name + " eof")
	}
	if it.state.Rewound() {
		panic(it.name + " rewound")
	}
}

func (it *schemaIter[T]) Next(index.OiTran) {
	switch it.state {
	case iface.Eof:
		return // stick at eof
	case iface.Rewound:
		it.i = it.lo
	default:
		it.i++
	}
	if it.skip.start != 0 {
		for it.i < it.hi && !it.skip.matches(it.items[it.i].keyAt()) {
			it.i++
		}
	}
	if it.i >= it.hi {
		it.state = iface.Eof
		return
	}
	it.state = iface.Within
}

func (it *schemaIter[T]) Prev(index.OiTran) {
	switch it.state {
	case iface.Eof:
		return // stick at eof
	case iface.Rewound:
		it.i = it.hi - 1
	default:
		it.i--
	}
	if it.skip.start != 0 {
		for it.i >= it.lo && !it.skip.matches(it.items[it.i].keyAt()) {
			it.i--
		}
	}
	if it.i < it.lo {
		it.state = iface.Eof
		return
	}
	it.state = iface.Within
}

func (it *schemaIter[T]) Rewind() {
	it.state = iface.Rewound
	it.i = 0
}

// bounds returns the [lo, hi) positions matching the key range [org, end).
func (it *schemaIter[T]) bounds(org, end string) (int, int) {
	return it.bound(len(it.items),
		func(i int) string { return it.items[i].boundAt() }, org, end)
}

func (it *schemaIter[T]) Range(rng index.Range) {
	it.lo, it.hi = it.bounds(rng.Org, rng.End)
	it.skip = schemaSkipScan{}
	it.state = iface.Rewound
	it.i = 0
}

func (it *schemaIter[T]) SkipScan(prefixRng, suffixRng index.Range, skipStart int) {
	if !it.skipOK {
		assert.ShouldNotReachHere()
	}
	assert.That(skipStart > 0)
	it.Range(prefixRng)
	it.skip = schemaSkipScan{rng: suffixRng, start: skipStart}
}
