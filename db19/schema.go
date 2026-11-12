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
	return getSchemaTable(table) != nil
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
	case "tables_":
		return tablesTable{}
	case "columns_":
		return columnsTable{}
	case "indexes_":
		return indexesTable{}
	case "views_":
		return viewsTable{}
	case "history_":
		return historyTable{}
	case "dbstats_":
		return statsTable{}
	}
	return nil
}

var tablesSchema = &schema.Schema{
	Table:   "tables_",
	Columns: []string{"table"},
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

// SchemaTableAlias maps a "_" schema table name to the name of the schema
// table it represents, so it reuses that table's permissions.
// TEMPORARY to keep new "_" tables separate.
func SchemaTableAlias(table string) string {
	if IsSchemaTable(table) {
		return table[:len(table)-1]
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

type tablesTable struct{}

func (tablesTable) schema() *schema.Schema { return tablesSchema }

func (tablesTable) nrows(t *tran) (n int) {
	for range t.meta.Tables() {
		n++
	}
	return n
}

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
	if t.meta.GetRoSchema(name) == nil {
		return ""
	}
	return tablesRow(name)
}

func tablesRow(table string) core.Record {
	var rb core.RecordBuilder
	rb.Add(core.SuStr(table))
	return rb.Build()
}

func newTablesIter(t *ReadTran) *schemaIter[tableEntry] {
	items := make([]tableEntry, 0, 32)
	for sc := range t.meta.Tables() {
		items = append(items, tableEntry{sc.Table})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].table < items[j].table })
	return &schemaIter[tableEntry]{
		items: items,
		bound: nameBounds,
		name:  "tablesIter",
		hi:    len(items),
	}
}

type tableEntry struct{ table string }

func (e tableEntry) keyAt() string    { return core.Pack(core.SuStr(e.table)) }
func (e tableEntry) boundAt() string  { return e.table }
func (e tableEntry) row() core.Record { return tablesRow(e.table) }

// columns ----------------------------------------------------------

type columnsTable struct{}

func (columnsTable) schema() *schema.Schema { return columnsSchema }

func (columnsTable) nrows(t *tran) (n int) {
	for sc := range t.meta.Tables() {
		for _, col := range sc.Columns {
			if col != "-" {
				n++
			}
		}
		n += len(sc.Derived)
	}
	return n
}

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

type indexesTable struct{}

func (indexesTable) schema() *schema.Schema { return indexesSchema }

func (indexesTable) nrows(t *tran) (n int) {
	for sc := range t.meta.Tables() {
		n += len(sc.Indexes)
	}
	return n
}

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

type viewsTable struct{}

func (viewsTable) schema() *schema.Schema { return viewsSchema }

func (viewsTable) nrows(t *tran) (n int) {
	for range t.meta.Views() {
		n++
	}
	return n
}

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

type historyTable struct{}

func (historyTable) schema() *schema.Schema { return historySchema }
func (historyTable) nrows(*tran) int        { return historyNrows }

func (historyTable) lookup(*ReadTran, int, string) core.Record {
	panic("history_ lookup not implemented")
}

func (historyTable) indexIter(t *ReadTran, iIndex int) index.IndexIter {
	assert.That(iIndex == 0)
	return newHistoryIter(t)
}

// Enumerating database states would be too slow, so use a rough estimate.
func (historyTable) rangeFrac(*ReadTran, int, string, string) float64 {
	return .5
}

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

// rangeFrac computes the fraction from the rows, since schema tables have no btrees.
func (it *schemaIter[T]) rangeFrac(org, end string) float64 {
	lo, hi := it.bounds(org, end)
	if len(it.items) == 0 || hi <= lo {
		return 0
	}
	return float64(hi-lo) / float64(len(it.items))
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
