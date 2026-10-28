// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"
	"sort"
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stats"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/tsc"
)

// schema implements virtual tables for tables, columns, indexes, and views

type schemaTable struct {
	cache
	tran  QueryTran
	state iface.State
	metrics
}

func (*schemaTable) Init() {
}

func (*schemaTable) Indexes() [][]string {
	return nil
}

func (*schemaTable) Fixed() Fixed {
	return nil
}

func (*schemaTable) Order() []string {
	return nil
}

func (*schemaTable) Updateable() string {
	return ""
}

func (*schemaTable) SingleTable() bool {
	return false // not a physical table
}

func (*schemaTable) rowSize() int {
	return 32 // ???
}

func (*schemaTable) Output(th *Thread, _ Record) {
	panic("can't output to schema table")
}

func (*schemaTable) optimize(mode Mode, req Require) (Cost, Cost, any) {
	if mode != CursorMode && len(req.cols) == 0 {
		return 0, 1000, nil
	}
	return impossible, impossible, nil
}

func (*schemaTable) setApproach(Require, any, QueryTran) {
}

func (*schemaTable) Select(Sels) {
	assert.ShouldNotReachHere()
}

func (*schemaTable) fastSingle() bool {
	return false
}

func (*schemaTable) knowExactNrows() bool {
	return true
}

func (*schemaTable) Simple(*Thread) []Row {
	panic("Simple not implemented for schema tables")
}

func (st *schemaTable) Metrics() *metrics {
	return &st.metrics
}

func (*schemaTable) Clone() *Table {
	panic("can't clone schema table")
}

func (*schemaTable) SetIndex(index []string, mode Mode) {
}

func (*schemaTable) UniqueIndexes() [][]string {
	return nil
}

// schemaTableNames are the virtual schema tables.
// Note: indexes does not include them.
var schemaTableNames = []string{"tables", "columns", "indexes", "views", "history", "dbstats"}

// schemaTableCols returns the columns of a virtual schema table,
// or nil if table is not one.
func schemaTableCols(table string) []string {
	switch table {
	case "tables":
		return tablesFields[0]
	case "columns":
		return columnsFields[0]
	case "indexes":
		return indexesFields[0]
	case "views":
		return viewsFields[0]
	case "history":
		return historyFields[0]
	case "dbstats":
		return statsFields[0]
	}
	return nil
}

// selValue decodes a packed selection value as a string.
// It returns "" if the value is not a string, which no table or column
// name matches, so Select and Lookup (including non-string join keys)
// tolerate it instead of panicking. This also covers the invalid packed
// values the fuzzer uses for "non-existent" selects.
func selValue(v string) string {
	if v == "" {
		return "" // empty string
	}
	if v[0] != PackString {
		return ""
	}
	return v[1:]
}

//-------------------------------------------------------------------

type Tables struct {
	schemaTable
	info     []*meta.Info
	selTable string
	selLevel int // 0 = no select, 1 = table
	i        int
}

func (*Tables) String() string {
	return "tables"
}

func (*Tables) Name() string {
	return "tables"
}

func (*Tables) optimize(mode Mode, req Require) (Cost, Cost, any) {
	if mode != CursorMode && req.SatisfiedBy(tablesIndexes[0]) {
		return 0, 1000, nil
	}
	return impossible, impossible, nil
}

func (ts *Tables) Transform() Query {
	return ts
}

var tablesIndexes = [][]string{{"table"}}

func (*Tables) Indexes() [][]string {
	return tablesIndexes
}

func (*Tables) Keys() [][]string {
	return tablesIndexes
}

var tablesFields = [][]string{{"table", "nrows", "totalsize"}}

func (*Tables) Columns() []string {
	return tablesFields[0]
}

func (*Tables) Header() *Header {
	return NewHeader(tablesFields, tablesFields[0])
}

func (ts *Tables) Nrows() (int, int) {
	ts.ensure()
	n := len(ts.info)
	return n, n
}

func (ts *Tables) SetTran(tran QueryTran) {
	tran.CheckPerm("tables", PermRead)
	ts.tran = tran
	ts.info = nil
}

// Select restricts Get to the {table} index prefix.
// Only the "table" column is used; other columns are ignored.
// Select(nil) clears the restriction.
func (ts *Tables) Select(sels Sels) {
	ts.selLevel = 0
	ts.selTable = ""
	if sels != nil {
		if v, ok := sels.Get("table"); ok {
			// a non-string value decodes to "" which no table name matches
			ts.selTable = selValue(v)
			ts.selLevel = 1
		}
	}
	ts.Rewind()
}

func (ts *Tables) Rewind() {
	ts.i = -1
	ts.state = iface.Rewound
}

func (ts *Tables) Get(_ *Thread, dir Dir) Row {
	defer func(t uint64) { ts.tget += tsc.Read() - t }(tsc.Read())
	ts.ensure()
	if ts.state.Eof() {
		return Row{}
	}
	if ts.selLevel > 0 && ts.state.Rewound() {
		// position at the selected row
		i, found := ts.seek(ts.selTable)
		if !found {
			ts.state = iface.Eof
			return Row{}
		}
		ts.i = i
		ts.state = iface.Within
		ts.ngets++
		return ts.row(ts.info[ts.i])
	}
	if dir == Next {
		if ts.state.Rewound() {
			ts.i = -1
		}
		ts.i++
	} else { // Prev
		if ts.state.Rewound() {
			ts.i = len(ts.info)
		}
		ts.i--
	}
	if ts.i < 0 || len(ts.info) <= ts.i {
		ts.state = iface.Eof
		return Row{}
	}
	ts.state = iface.Within
	if !ts.inRange(ts.info[ts.i].Table) {
		ts.state = iface.Eof
		return Row{}
	}
	ts.ngets++
	return ts.row(ts.info[ts.i])
}

// seek returns the position of table in ts.info and whether it was found.
func (ts *Tables) seek(table string) (int, bool) {
	return slices.BinarySearchFunc(ts.info, table,
		func(info *meta.Info, table string) int {
			return strings.Compare(info.Table, table)
		})
}

// inRange reports whether table is within the current select.
func (ts *Tables) inRange(table string) bool {
	if ts.selLevel == 1 {
		return table == ts.selTable
	}
	return true
}

func (*Tables) row(info *meta.Info) Row {
	var rb RecordBuilder
	rb.Add(SuStr(info.Table))
	rb.Add(IntVal(info.Nrows))
	rb.Add(IntVal(info.Size))
	rec := rb.Build()
	return NewRow(rec)
}

func (ts *Tables) ensure() {
	if ts.info != nil {
		return
	}
	ts.info = ts.tran.GetAllInfo()

	// count the virtual schema tables instead of constructing them
	var ncols, nidxs int
	for s := range ts.tran.Schemas() {
		for _, col := range s.Columns {
			if col != "-" {
				ncols++
			}
		}
		ncols += len(s.Derived)
		nidxs += len(s.Indexes)
	}
	// the columns table also includes the virtual schema tables themselves
	for _, name := range schemaTableNames {
		ncols += len(schemaTableCols(name))
	}
	var nviews int
	for range ts.tran.Views() {
		nviews++
	}
	nstats := 0
	for _, tbl := range ts.tran.Stats() {
		nstats += len(tbl.Columns)
	}

	ts.info = append(ts.info,
		&meta.Info{Table: "tables", Nrows: len(ts.info) + len(schemaTableNames)},
		&meta.Info{Table: "columns", Nrows: ncols},
		&meta.Info{Table: "indexes", Nrows: nidxs},
		&meta.Info{Table: "views", Nrows: nviews},
		&meta.Info{Table: "history", Nrows: historyNrows},
		&meta.Info{Table: "dbstats", Nrows: nstats},
	)
	sort.Slice(ts.info,
		func(i, j int) bool { return ts.info[i].Table < ts.info[j].Table })
}

func (ts *Tables) Lookup(th *Thread, sels Sels) Row {
	table := selValue(sels.MustGet("table"))
	if schemaTableCols(table) == nil {
		// real table: keyed lookup. Get builds rows from GetAllInfo, so
		// GetInfo must agree with it (see sizeTran) to keep Lookup consistent.
		if ti := ts.tran.GetInfo(table); ti != nil {
			ts.ngets++
			return ts.row(ti)
		}
		return Row{}
	}
	// virtual schema table: its nrows comes from ts.info (see ensure)
	ts.ensure()
	i, found := ts.seek(table)
	if !found {
		return Row{}
	}
	ts.ngets++
	return ts.row(ts.info[i])
}

//-------------------------------------------------------------------

type Columns struct {
	schemaTable
	schema    []*meta.Schema
	order     []int // reusable, sorted column indexes for the current table
	selTable  string
	selColumn string
	si        int
	ci        int
	selLevel  int // 0 = no select, 1 = table, 2 = table and column
	nrows     int
}

func (*Columns) String() string {
	return "columns"
}

func (*Columns) Name() string {
	return "columns"
}

func (*Columns) optimize(mode Mode, req Require) (Cost, Cost, any) {
	if mode != CursorMode && req.SatisfiedBy(columnsIndexes[0]) {
		return 0, 1000, nil
	}
	return impossible, impossible, nil
}

func (cs *Columns) Transform() Query {
	return cs
}

var columnsIndexes = [][]string{{"table", "column"}}

func (*Columns) Indexes() [][]string {
	return columnsIndexes
}

func (*Columns) Keys() [][]string {
	return columnsIndexes
}

var columnsFields = [][]string{{"table", "column", "field"}}

func (cs *Columns) Columns() []string {
	return columnsFields[0]
}

func (*Columns) Header() *Header {
	return NewHeader(columnsFields, columnsFields[0])
}

func (cs *Columns) Nrows() (int, int) {
	cs.ensure()
	return cs.nrows, cs.nrows
}

func (cs *Columns) getNrows() int {
	n := 0
	for _, schema := range cs.schema {
		for _, col := range schema.Columns {
			if col != "-" {
				n++
			}
		}
		n += len(schema.Derived)
	}
	return n
}

func (cs *Columns) SetTran(tran QueryTran) {
	tran.CheckPerm("columns", PermRead)
	cs.tran = tran
	cs.schema = nil
}

func (cs *Columns) Rewind() {
	cs.si = -1
	cs.state = iface.Rewound
}

// Select restricts Get to a prefix of the emulated {table, column} index.
// The prefix is either "table" or "table, column"; other columns are ignored.
// Select(nil) clears the restriction.
func (cs *Columns) Select(sels Sels) {
	cs.selLevel = 0
	cs.selTable = ""
	cs.selColumn = ""
	if sels != nil {
		if v, ok := sels.Get("table"); ok {
			cs.selTable = selValue(v)
			cs.selLevel = 1
			if v, ok := sels.Get("column"); ok {
				cs.selColumn = selValue(v)
				cs.selLevel = 2
			}
		}
	}
	cs.Rewind()
}

func (cs *Columns) Get(_ *Thread, dir Dir) Row {
	defer func(t uint64) { cs.tget += tsc.Read() - t }(tsc.Read())
	cs.ensure()
	if cs.state.Eof() {
		return Row{}
	}
	if cs.selLevel > 0 && cs.state.Rewound() {
		// position at the start of the selected range
		if !cs.seek(dir) {
			cs.state = iface.Eof
			return Row{}
		}
		cs.state = iface.Within
		cs.ngets++
		schema := cs.schema[cs.si]
		col, fld := columnOrDerived(schema, cs.order[cs.ci])
		return cs.row(schema.Table, col, fld)
	}
	if dir == Next {
		if cs.state.Rewound() {
			cs.si, cs.ci = 0, -1
			cs.sortCols(cs.si)
		}
		cs.ci++
		for cs.ci >= len(cs.order) {
			cs.si++
			if cs.si >= len(cs.schema) {
				cs.state = iface.Eof
				return Row{}
			}
			cs.ci = 0
			cs.sortCols(cs.si)
		}
	} else { // Prev
		if cs.state.Rewound() {
			cs.si, cs.ci = len(cs.schema), 0
		}
		cs.ci--
		for cs.ci < 0 {
			cs.si--
			if cs.si < 0 {
				cs.state = iface.Eof
				return Row{}
			}
			cs.sortCols(cs.si)
			cs.ci = len(cs.order) - 1
		}
	}
	cs.state = iface.Within
	schema := cs.schema[cs.si]
	col, fld := columnOrDerived(schema, cs.order[cs.ci])
	if !cs.inRange(schema.Table, col) {
		cs.state = iface.Eof
		return Row{}
	}
	cs.ngets++
	return cs.row(schema.Table, col, fld)
}

// seek positions to the first (Next) or last (Prev) row matching the select,
// returning false if the range is empty.
func (cs *Columns) seek(dir Dir) bool {
	si := sort.Search(len(cs.schema), func(i int) bool {
		return cs.schema[i].Table >= cs.selTable
	})
	if si >= len(cs.schema) || cs.schema[si].Table != cs.selTable {
		return false
	}
	cs.sortCols(si)
	cs.si = si
	if cs.selLevel == 1 {
		if len(cs.order) == 0 {
			return false
		}
		if dir == Next {
			cs.ci = 0
		} else {
			cs.ci = len(cs.order) - 1
		}
		return true
	}
	s := cs.schema[si]
	ci, found := slices.BinarySearchFunc(cs.order, cs.selColumn,
		func(i int, col string) int {
			name, _ := columnOrDerived(s, i)
			return strings.Compare(name, col)
		})
	if !found {
		return false
	}
	cs.ci = ci
	return true
}

// sortCols loads cs.order with the indexes (into Columns followed by Derived)
// of the given schema's columns, sorted by column name and skipping "-".
// It reuses cs.order to avoid allocating for each table.
func (cs *Columns) sortCols(si int) {
	s := cs.schema[si]
	cs.order = cs.order[:0]
	for i, col := range s.Columns {
		if col != "-" {
			cs.order = append(cs.order, i)
		}
	}
	for i := range s.Derived {
		cs.order = append(cs.order, len(s.Columns)+i)
	}
	slices.SortFunc(cs.order, func(a, b int) int {
		ac, _ := columnOrDerived(s, a)
		bc, _ := columnOrDerived(s, b)
		return strings.Compare(ac, bc)
	})
}

// inRange reports whether (table, col) is within the current select.
// Assumes the iteration order (table, column) matches the index order.
func (cs *Columns) inRange(table, col string) bool {
	switch cs.selLevel {
	case 1:
		return table == cs.selTable
	case 2:
		return table == cs.selTable && col == cs.selColumn
	}
	return true
}

func (*Columns) row(table, col string, fld int) Row {
	var rb RecordBuilder
	rb.Add(SuStr(table))
	rb.Add(SuStr(col))
	rb.Add(IntVal(fld))
	rec := rb.Build()
	return NewRow(rec)
}

func columnOrDerived(schema *meta.Schema, i int) (string, int) {
	if i >= len(schema.Columns) {
		i -= len(schema.Columns)
		return str.UnCapitalize(schema.Derived[i]), -1
	}
	return schema.Columns[i], i
}

func (cs *Columns) ensure() {
	if cs.schema != nil {
		return
	}
	cs.schema = slices.AppendSeq(make([]*meta.Schema, 0, 32), cs.tran.Schemas())
	for _, name := range schemaTableNames {
		cs.schema = append(cs.schema,
			&meta.Schema{Table: name, Columns: schemaTableCols(name)})
	}
	sort.Slice(cs.schema,
		func(i, j int) bool { return cs.schema[i].Table < cs.schema[j].Table })
	cs.nrows = cs.getNrows()
}

func (cs *Columns) Lookup(th *Thread, sels Sels) Row {
	table := selValue(sels.MustGet("table"))
	column := selValue(sels.MustGet("column"))
	if cols := schemaTableCols(table); cols != nil {
		// virtual schema table, not in the meta schema
		if i := slices.Index(cols, column); i >= 0 {
			return cs.row(table, column, i)
		}
		return Row{}
	}
	if !cs.tran.HasTable(table) {
		return Row{}
	}
	if column == "-" {
		return Row{} // not included in the columns table
	}
	ts := cs.tran.GetSchema(table)
	fld := -1
	if i := slices.Index(ts.Columns, column); i >= 0 {
		fld = i
	} else if !slices.ContainsFunc(ts.Derived,
		// Derived are stored capitalized but reported uncapitalized
		func(d string) bool { return str.EqualUncap(d, column) }) {
		return Row{}
	}
	return cs.row(table, column, fld)
}

//-------------------------------------------------------------------
// note: indexes does not include tables, columns, indexes, views

type Indexes struct {
	schemaTable
	schema []*meta.Schema
	si     int
	ci     int
	nrows  int
}

func (*Indexes) String() string {
	return "indexes"
}

func (*Indexes) Name() string {
	return "indexes"
}

func (is *Indexes) Transform() Query {
	return is
}

func (*Indexes) Keys() [][]string {
	return [][]string{{"table", "columns"}}
}

var indexesFields = [][]string{{"table", "columns", "key",
	"fktable", "fkcolumns", "fkmode"}}

func (*Indexes) Columns() []string {
	return indexesFields[0]
}

func (*Indexes) Header() *Header {
	return NewHeader(indexesFields, indexesFields[0])
}

func (is *Indexes) Nrows() (int, int) {
	is.ensure()
	return is.nrows, is.nrows
}

func (is *Indexes) getNrows() int {
	n := 0
	for _, schema := range is.schema {
		n += len(schema.Indexes)
	}
	return n
}

func (is *Indexes) SetTran(tran QueryTran) {
	tran.CheckPerm("indexes", PermRead)
	is.tran = tran
	is.schema = nil
}

func (is *Indexes) Rewind() {
	is.state = iface.Rewound
}

func (is *Indexes) Get(_ *Thread, dir Dir) Row {
	defer func(t uint64) { is.tget += tsc.Read() - t }(tsc.Read())
	is.ensure()
	if is.state.Eof() {
		return Row{}
	}
	if dir == Next {
		if is.state.Rewound() {
			is.si, is.ci = 0, -1
		}
		is.ci++
		for is.si >= len(is.schema) || is.ci >= len(is.schema[is.si].Indexes) {
			is.si++
			if is.si >= len(is.schema) {
				is.state = iface.Eof
				return Row{}
			}
			is.ci = 0
		}
	} else { // Prev
		if is.state.Rewound() {
			is.si, is.ci = len(is.schema), 0
		}
		is.ci--
		for is.ci < 0 {
			is.si--
			if is.si < 0 {
				is.state = iface.Eof
				return Row{}
			}
			is.ci = len(is.schema[is.si].Indexes) - 1
		}
	}
	is.state = iface.Within
	schema := is.schema[is.si]
	is.ngets++
	return is.row(schema.Table, &schema.Indexes[is.ci])
}

func (is *Indexes) row(table string, idx *schema.Index) Row {
	var rb RecordBuilder
	rb.Add(SuStr(table))
	rb.Add(SuStr(str.Join(",", idx.Columns)))
	switch idx.Mode {
	case 'k':
		rb.Add(True)
	case 'i':
		rb.Add(False)
	case 'u':
		rb.Add(SuStr("u"))
	default:
		assert.ShouldNotReachHere()
	}
	if idx.Fk.Table != "" {
		rb.Add(SuStr(idx.Fk.Table))
		rb.Add(SuStr(str.Join(",", idx.Fk.Columns)))
		rb.Add(SuInt16(int(idx.Fk.Mode)))
	}
	rec := rb.Build()
	return NewRow(rec)
}

func (is *Indexes) ensure() {
	if is.schema != nil {
		return
	}
	is.schema = slices.AppendSeq(make([]*meta.Schema, 0, 32), is.tran.Schemas())
	sort.Slice(is.schema,
		func(i, j int) bool { return is.schema[i].Table < is.schema[j].Table })
	is.nrows = is.getNrows()
}

func (is *Indexes) Lookup(th *Thread, sels Sels) Row {
	table := ToStr(Unpack(sels.MustGet("table")))
	cols := ToStr(Unpack(sels.MustGet("columns")))
	if !is.tran.HasTable(table) {
		return Row{} // note: indexes does not include the virtual schema tables
	}
	ts := is.tran.GetSchema(table)
	if idx := ts.FindIndex(str.Split(cols, ",")); idx != nil {
		return is.row(table, idx)
	}
	return Row{}
}

//-------------------------------------------------------------------

type Views struct {
	views []view
	schemaTable
	selName  string
	i        int
	selLevel int // 0 = no select, 1 = view_name
}

type view struct {
	name string
	def  string
}

func (*Views) String() string {
	return "views"
}

func (*Views) Name() string {
	return "views"
}

func (*Views) optimize(mode Mode, req Require) (Cost, Cost, any) {
	if mode != CursorMode && req.SatisfiedBy(viewsIndexes[0]) {
		return 0, 1000, nil
	}
	return impossible, impossible, nil
}

func (vs *Views) Transform() Query {
	return vs
}

var viewsIndexes = [][]string{{"view_name"}}

func (*Views) Indexes() [][]string {
	return viewsIndexes
}

func (*Views) Keys() [][]string {
	return viewsIndexes
}

var viewsFields = [][]string{{"view_name", "view_definition"}}

func (*Views) Columns() []string {
	return viewsFields[0]
}

func (*Views) Header() *Header {
	return NewHeader(viewsFields, viewsFields[0])
}

func (vs *Views) Nrows() (int, int) {
	vs.ensure()
	n := len(vs.views)
	return n, n
}

func (vs *Views) SetTran(tran QueryTran) {
	tran.CheckPerm("views", PermRead)
	vs.tran = tran
	vs.views = nil
}

// Select restricts Get to the {view_name} index prefix.
// Only the "view_name" column is used; other columns are ignored.
// Select(nil) clears the restriction.
func (vs *Views) Select(sels Sels) {
	vs.selLevel = 0
	vs.selName = ""
	if sels != nil {
		if v, ok := sels.Get("view_name"); ok {
			// a non-string value decodes to "" which no view name matches
			vs.selName = selValue(v)
			vs.selLevel = 1
		}
	}
	vs.Rewind()
}

func (vs *Views) Rewind() {
	vs.i = -1
	vs.state = iface.Rewound
}

func (vs *Views) Get(_ *Thread, dir Dir) Row {
	vs.ensure()
	if vs.state.Eof() {
		return Row{}
	}
	if vs.selLevel > 0 && vs.state.Rewound() {
		// position at the selected view
		if !vs.seek() {
			vs.state = iface.Eof
			return Row{}
		}
		vs.state = iface.Within
		return vs.row(vs.views[vs.i].name, vs.views[vs.i].def)
	}
	if dir == Next {
		if vs.state.Rewound() {
			vs.i = -1
		}
		vs.i++
	} else { // Prev
		if vs.state.Rewound() {
			vs.i = len(vs.views)
		}
		vs.i--
	}
	if vs.i < 0 || len(vs.views) <= vs.i {
		vs.state = iface.Eof
		return Row{}
	}
	vs.state = iface.Within
	if !vs.inRange(vs.views[vs.i].name) {
		vs.state = iface.Eof
		return Row{}
	}
	return vs.row(vs.views[vs.i].name, vs.views[vs.i].def)
}

// seek positions to the selected view, returning false if it is not present.
func (vs *Views) seek() bool {
	n := vs.Len()
	idx := sort.Search(n, func(i int) bool {
		return vs.views[i].name >= vs.selName
	})
	if idx >= n || vs.views[idx].name != vs.selName {
		return false
	}
	vs.i = idx
	return true
}

// inRange reports whether name is within the current select.
func (vs *Views) inRange(name string) bool {
	if vs.selLevel == 1 {
		return name == vs.selName
	}
	return true
}

func (vs *Views) row(name, def string) Row {
	var rb RecordBuilder
	rb.Add(SuStr(name))
	rb.Add(SuStr(def))
	rec := rb.Build()
	return NewRow(rec)
}

func (vs *Views) ensure() {
	if vs.views != nil {
		return
	}
	vs.views = make([]view, 0, 16)
	for name, def := range vs.tran.Views() {
		vs.views = append(vs.views, view{name, def})
	}
	sort.Sort(vs)
}

func (vs *Views) Len() int {
	return len(vs.views)
}
func (vs *Views) Less(i, j int) bool {
	return vs.views[i].name < vs.views[j].name
}
func (vs *Views) Swap(i, j int) {
	vs.views[i], vs.views[j] = vs.views[j], vs.views[i]
}

func (vs *Views) Lookup(th *Thread, sels Sels) Row {
	name := selValue(sels.MustGet("view_name"))
	if def := vs.tran.GetView(name); def != "" {
		return vs.row(name, def)
	}
	return Row{}
}

//-------------------------------------------------------------------

type History struct {
	schemaTable
	off uint64
}

func (*History) String() string {
	return "history"
}

// Name deliberately not implemented so dbms.get doesn't try to do Lookup

func (his *History) Transform() Query {
	return his
}

func (*History) Keys() [][]string {
	return [][]string{{"asof"}}
}

var historyFields = [][]string{{"asof"}}

func (*History) Columns() []string {
	return historyFields[0]
}

func (*History) Header() *Header {
	return NewHeader(historyFields, historyFields[0])
}

// historyNrows is a constant estimate for the history data.
// Counting the database states would be too slow.
const historyNrows = 1000

func (his *History) Nrows() (int, int) {
	return historyNrows, historyNrows
}

func (his *History) SetTran(tran QueryTran) {
	tran.CheckPerm("history", PermRead)
	his.tran = tran
}

func (his *History) Rewind() {
	his.state = iface.Rewound
}

func (his *History) Get(_ *Thread, dir Dir) Row {
	defer func(t uint64) { his.tget += tsc.Read() - t }(tsc.Read())
	if his.state.Eof() {
		return Row{}
	}
	if his.state.Rewound() {
		his.off = 0
	}
	var state *db19.DbState
	if dir == Next {
		state = db19.NextState(his.tran.GetStore(), his.off)
	} else {
		state = db19.PrevState(his.tran.GetStore(), his.off)
	}
	if state == nil {
		his.state = iface.Eof
		return Row{}
	}
	his.state = iface.Within
	his.off = state.Off
	var rb RecordBuilder
	rb.Add(SuDateFromUnixMilli(int(state.Asof)))
	rec := rb.Build()
	his.ngets++
	return NewRow(rec)
}

func (his *History) Lookup(th *Thread, sels Sels) Row {
	panic("not implemented")
}

//-------------------------------------------------------------------

type StatsTable struct {
	schemaTable
	stats  stats.Stats
	tables []string
	si     int
	ci     int
	nrows  int
}

func (*StatsTable) String() string {
	return "dbstats"
}

func (*StatsTable) Name() string {
	return "dbstats"
}

func (st *StatsTable) Transform() Query {
	return st
}

func (*StatsTable) Keys() [][]string {
	return [][]string{{"table", "column"}}
}

var statsFields = [][]string{{"table", "column", "distinct", "quantiles", "common"}}

func (*StatsTable) Columns() []string {
	return statsFields[0]
}

func (*StatsTable) Header() *Header {
	return NewHeader(statsFields, statsFields[0])
}

func (st *StatsTable) Nrows() (int, int) {
	st.ensure()
	return st.nrows, st.nrows
}

func (st *StatsTable) SetTran(tran QueryTran) {
	tran.CheckPerm("dbstats", PermRead)
	st.tran = tran
	st.stats = nil
}

func (st *StatsTable) Rewind() {
	st.si = -1
	st.state = iface.Rewound
}

func (st *StatsTable) Get(_ *Thread, dir Dir) Row {
	defer func(t uint64) { st.tget += tsc.Read() - t }(tsc.Read())
	st.ensure()
	if st.state.Eof() {
		return Row{}
	}
	var table, col string
	var cs stats.ColStats
	if dir == Next {
		if st.state.Rewound() {
			st.si, st.ci = 0, -1
		}
		st.ci++
		for st.si >= len(st.tables) ||
			st.ci >= len(st.stats[st.tables[st.si]].Columns) {
			st.si++
			if st.si >= len(st.tables) {
				st.state = iface.Eof
				return Row{}
			}
			st.ci = 0
		}
	} else { // Prev
		if st.state.Rewound() {
			st.si = len(st.tables)
			st.ci = 0
		}
		st.ci--
		for st.ci < 0 {
			st.si--
			if st.si < 0 {
				st.state = iface.Eof
				return Row{}
			}
			st.ci = len(st.stats[st.tables[st.si]].Columns) - 1
		}
	}
	table = st.tables[st.si]
	cols := sortedCols(st.stats[table])
	col = cols[st.ci]
	cs = st.stats[table].Columns[col]
	st.state = iface.Within
	st.ngets++
	return st.row(table, col, cs)
}

func (*StatsTable) row(table string, col string, cs stats.ColStats) Row {
	var rb RecordBuilder
	rb.Add(SuStr(table))
	rb.Add(SuStr(col))
	rb.Add(IntVal(cs.Cardinality))
	rb.Add(formatQuantiles(cs.Quantiles))
	rb.Add(formatTops(cs.Tops))
	rec := rb.Build()
	return NewRow(rec)
}

func sortedCols(ts stats.TableStats) []string {
	cols := make([]string, 0, len(ts.Columns))
	for col := range ts.Columns {
		cols = append(cols, col)
	}
	sort.Strings(cols)
	return cols
}

func formatQuantiles(quantiles []string) *SuObject {
	parts := &SuObject{}
	for _, q := range quantiles {
		parts.Add(Unpack(q))
	}
	return parts
}

func formatTops(tops []stats.Top) *SuObject {
	parts := &SuObject{}
	for _, t := range tops {
		part := &SuObject{}
		part.Add(SuDnum{Dnum: dnum.FromFloat(t.Frac)})
		part.Add(Unpack(t.Value))
		parts.Add(part)
	}
	return parts
}

func (st *StatsTable) ensure() {
	if st.stats != nil {
		return
	}
	st.stats = st.tran.Stats()
	st.tables = make([]string, 0, len(st.stats))
	for table := range st.stats {
		st.tables = append(st.tables, table)
	}
	sort.Strings(st.tables)
	st.nrows = 0
	for _, table := range st.tables {
		st.nrows += len(st.stats[table].Columns)
	}
}

func (st *StatsTable) Lookup(th *Thread, sels Sels) Row {
	table := ToStr(Unpack(sels.MustGet("table")))
	column := ToStr(Unpack(sels.MustGet("column")))
	stats := st.tran.Stats()
	if ts, ok := stats[table]; ok {
		if cs, ok := ts.Columns[column]; ok {
			return st.row(table, column, cs)
		}
	}
	return Row{}
}
