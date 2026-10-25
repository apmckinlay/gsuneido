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

func (*schemaTable) optimize(_ Mode, req Require) (Cost, Cost, any) {
	if len(req.cols) == 0 {
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
var schemaTableNames = []string{"tables", "columns", "indexes", "views", "dbstats"}

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
	case "dbstats":
		return statsFields[0]
	}
	return nil
}

//-------------------------------------------------------------------

type Tables struct {
	schemaTable
	info     []*meta.Info
	posTable string
	i        int
	gap      bool
}

func (*Tables) String() string {
	return "tables"
}

func (*Tables) Name() string {
	return "tables"
}

func (ts *Tables) Transform() Query {
	return ts
}

func (*Tables) Keys() [][]string {
	return [][]string{{"table"}}
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
	if ts.state.Within() && !ts.gap {
		ts.posTable = ts.info[ts.i].Table
	}
	ts.tran = tran
	ts.info = nil
	if ts.state.Within() {
		ts.ensure()
		var found bool
		ts.i, found = slices.BinarySearchFunc(ts.info, ts.posTable,
			func(info *meta.Info, table string) int {
				return strings.Compare(info.Table, table)
			})
		ts.gap = !found
	}
}

func (ts *Tables) Rewind() {
	ts.i = -1
	ts.state = iface.Rewound
	ts.gap = false
}

func (ts *Tables) Get(_ *Thread, dir Dir) Row {
	defer func(t uint64) { ts.tget += tsc.Read() - t }(tsc.Read())
	ts.ensure()
	if ts.state.Eof() {
		return Row{}
	}
	if dir == Next {
		if ts.state.Rewound() {
			ts.i = -1
		}
		if !ts.gap {
			ts.i++
		}
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
	ts.gap = false
	ts.ngets++
	return ts.row(ts.info[ts.i])
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

	cols := Columns{tran: ts.tran}
	ncols, _ := cols.Nrows()

	idxs := Indexes{tran: ts.tran}
	nidxs, _ := idxs.Nrows()

	views := Views{tran: ts.tran}
	nviews, _ := views.Nrows()

	stats := StatsTable{tran: ts.tran}
	nstats, _ := stats.Nrows()

	ts.info = append(ts.info,
		// +5 for tables, columns, indexes, views, dbstats
		&meta.Info{Table: "tables", Nrows: len(ts.info) + 5},
		&meta.Info{Table: "columns", Nrows: ncols},
		&meta.Info{Table: "indexes", Nrows: nidxs},
		&meta.Info{Table: "views", Nrows: nviews},
		&meta.Info{Table: "dbstats", Nrows: nstats},
	)
	sort.Slice(ts.info,
		func(i, j int) bool { return ts.info[i].Table < ts.info[j].Table })
}

func (ts *Tables) Lookup(th *Thread, sels Sels) Row {
	table := ToStr(Unpack(sels.MustGet("table")))
	if schemaTableCols(table) != nil {
		var rb RecordBuilder
		rb.Add(SuStr(table))
		rec := rb.Build()
		ts.ngets++
		return NewRow(rec)
	}
	if ti := ts.tran.GetInfo(table); ti != nil {
		ts.ngets++
		return ts.row(ti)
	}
	return Row{}
}

//-------------------------------------------------------------------

type Columns struct {
	schemaTable
	schema    []*meta.Schema
	order     []int // reusable, sorted column indexes for the current table
	posTable  string
	posColumn string
	si        int
	ci        int
	nrows     int
	gap       bool // ci is an insertion point rather than a returned column
}

func (*Columns) String() string {
	return "columns"
}

func (*Columns) Name() string {
	return "columns"
}

func (cs *Columns) Transform() Query {
	return cs
}

func (*Columns) Keys() [][]string {
	return [][]string{{"table", "column"}}
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
	var prev *meta.Schema
	if cs.state.Within() && cs.si < len(cs.schema) {
		prev = cs.schema[cs.si]
	}
	cs.tran = tran
	cs.schema = nil
	if !cs.state.Within() {
		return
	}
	cs.ensure()
	if cs.si < len(cs.schema) && cs.schema[cs.si] == prev {
		return
	}
	// Resume by name, since both table positions and column indexes may change.
	if !cs.gap {
		cs.posTable = prev.Table
		cs.posColumn, _ = columnOrDerived(prev, cs.order[cs.ci])
	}
	cs.si = sort.Search(len(cs.schema), func(i int) bool {
		return cs.schema[i].Table >= cs.posTable
	})
	cs.ci, cs.gap = 0, true
	if cs.si < len(cs.schema) {
		cs.sortCols(cs.si)
		if cs.schema[cs.si].Table == cs.posTable {
			var found bool
			cs.ci, found = slices.BinarySearchFunc(cs.order, cs.posColumn,
				func(i int, col string) int {
					name, _ := columnOrDerived(cs.schema[cs.si], i)
					return strings.Compare(name, col)
				})
			cs.gap = !found
		}
	} else {
		cs.order = cs.order[:0]
	}
}

func (cs *Columns) Rewind() {
	cs.si = -1
	cs.state = iface.Rewound
	cs.gap = false
}

func (cs *Columns) Select(sels Sels) {
	if sels != nil {
		panic(assert.ShouldNotReachHere())
	}
	cs.Rewind()
}

func (cs *Columns) Get(_ *Thread, dir Dir) Row {
	defer func(t uint64) { cs.tget += tsc.Read() - t }(tsc.Read())
	cs.ensure()
	if cs.state.Eof() {
		return Row{}
	}
	if dir == Next {
		if cs.state.Rewound() {
			cs.si, cs.ci = 0, -1
			cs.sortCols(cs.si)
		}
		if !cs.gap {
			cs.ci++
		}
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
	cs.gap = false
	cs.ngets++
	schema := cs.schema[cs.si]
	col, fld := columnOrDerived(schema, cs.order[cs.ci])
	return cs.row(schema.Table, col, fld)
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
	cs.schema = cs.tran.GetAllSchema()
	for _, name := range schemaTableNames {
		cs.schema = append(cs.schema,
			&meta.Schema{Table: name, Columns: schemaTableCols(name)})
	}
	sort.Slice(cs.schema,
		func(i, j int) bool { return cs.schema[i].Table < cs.schema[j].Table })
	cs.nrows = cs.getNrows()
}

func (cs *Columns) Lookup(th *Thread, sels Sels) Row {
	table := ToStr(Unpack(sels.MustGet("table")))
	column := ToStr(Unpack(sels.MustGet("column")))
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
		func(d string) bool { return str.UnCapitalize(d) == column }) {
		return Row{}
	}
	return cs.row(table, column, fld)
}

func (cs *Columns) Simple(*Thread) []Row {
	cs.ensure()
	rows := make([]Row, 0, cs.nrows)
	var order []int // reused across schemas to avoid allocating for each
	for _, schema := range cs.schema {
		order = order[:0]
		for i, col := range schema.Columns {
			if col != "-" {
				order = append(order, i)
			}
		}
		for i := range schema.Derived {
			order = append(order, len(schema.Columns)+i)
		}
		slices.SortFunc(order, func(a, b int) int {
			ac, _ := columnOrDerived(schema, a)
			bc, _ := columnOrDerived(schema, b)
			return strings.Compare(ac, bc)
		})
		for _, i := range order {
			col, fld := columnOrDerived(schema, i)
			rows = append(rows, cs.row(schema.Table, col, fld))
		}
	}
	return rows
}

//-------------------------------------------------------------------
// note: indexes does not include tables, columns, indexes, views

type Indexes struct {
	schemaTable
	schema     []*meta.Schema
	posIndexes []schema.Index // old prefix through the current index, retained across gaps
	posTable   string
	si         int
	ci         int
	nrows      int
	gap        bool
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
	var prev *meta.Schema
	if is.state.Within() && is.si < len(is.schema) {
		prev = is.schema[is.si]
	}
	is.tran = tran
	is.schema = nil
	if !is.state.Within() {
		return
	}
	is.ensure()
	if is.si < len(is.schema) && is.schema[is.si] == prev {
		return
	}
	if !is.gap {
		is.posTable = prev.Table
		is.posIndexes = prev.Indexes[:is.ci+1]
	}
	is.si = sort.Search(len(is.schema), func(i int) bool {
		return is.schema[i].Table >= is.posTable
	})
	is.ci, is.gap = 0, true
	if is.si < len(is.schema) && is.schema[is.si].Table == is.posTable {
		idxs := is.schema[is.si].Indexes
		// Additions append and removals preserve physical order. If the
		// current index is gone, resume after its last surviving predecessor.
		for p := len(is.posIndexes) - 1; p >= 0; p-- {
			i := slices.IndexFunc(idxs, func(idx schema.Index) bool {
				return slices.Equal(idx.Columns, is.posIndexes[p].Columns)
			})
			if i >= 0 {
				is.ci = i
				if p == len(is.posIndexes)-1 {
					is.gap = false
				} else {
					is.ci++
				}
				break
			}
		}
	}
}

func (is *Indexes) Rewind() {
	is.state = iface.Rewound
	is.gap = false
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
		if !is.gap {
			is.ci++
		}
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
	is.gap = false
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
	is.schema = is.tran.GetAllSchema()
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
	views []string
	schemaTable
	posName string
	i       int
	gap     bool
}

func (*Views) String() string {
	return "views"
}

func (*Views) Name() string {
	return "views"
}

func (vs *Views) Transform() Query {
	return vs
}

func (*Views) Keys() [][]string {
	return [][]string{{"view_name"}}
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
	n := len(vs.views) / 2
	return n, n
}

func (vs *Views) SetTran(tran QueryTran) {
	tran.CheckPerm("views", PermRead)
	if vs.state.Within() && !vs.gap {
		vs.posName = vs.views[vs.i]
	}
	vs.tran = tran
	vs.views = nil
	if vs.state.Within() {
		vs.ensure()
		vs.i = 2 * sort.Search(vs.Len(), func(i int) bool {
			return vs.views[2*i] >= vs.posName
		})
		vs.gap = vs.i == len(vs.views) || vs.views[vs.i] != vs.posName
	}
}

func (vs *Views) Rewind() {
	vs.i = -2
	vs.state = iface.Rewound
	vs.gap = false
}

func (vs *Views) Get(_ *Thread, dir Dir) Row {
	vs.ensure()
	if vs.state.Eof() {
		return Row{}
	}
	if dir == Next {
		if vs.state.Rewound() {
			vs.i = -2
		}
		if !vs.gap {
			vs.i += 2
		}
	} else { // Prev
		if vs.state.Rewound() {
			vs.i = len(vs.views)
		}
		vs.i -= 2
	}
	if vs.i < 0 || len(vs.views) <= vs.i {
		vs.state = iface.Eof
		return Row{}
	}
	vs.state = iface.Within
	vs.gap = false
	return vs.row(vs.views[vs.i], vs.views[vs.i+1])
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
	vs.views = vs.tran.GetAllViews()
	sort.Sort(vs)
}

func (vs *Views) Len() int {
	return len(vs.views) / 2
}
func (vs *Views) Less(i, j int) bool {
	return vs.views[i*2] < vs.views[j*2]
}
func (vs *Views) Swap(i, j int) {
	i *= 2
	j *= 2
	vs.views[i], vs.views[j] = vs.views[j], vs.views[i]
	vs.views[i+1], vs.views[j+1] = vs.views[j+1], vs.views[i+1]
}

func (vs *Views) Lookup(th *Thread, sels Sels) Row {
	name := ToStr(Unpack(sels.MustGet("view_name")))
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

var HistoryFields = [][]string{{"asof"}}

func (*History) Columns() []string {
	return HistoryFields[0]
}

func (*History) Header() *Header {
	return NewHeader(HistoryFields, HistoryFields[0])
}

func (his *History) Nrows() (int, int) {
	return 1000, 1000 // ???
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
