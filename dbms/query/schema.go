// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"sort"

	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

// schema implements virtual tables for tables, columns, indexes, and views

type schemaTable struct {
	cache
	tran QueryTran
	state
}

type state int

const (
	rewound state = iota
	within
	eof
)

func (*schemaTable) Init() {
}

func (*schemaTable) Indexes() [][]string {
	return nil
}

func (*schemaTable) Fixed() []Fixed {
	return nil
}

func (*schemaTable) Ordering() []string {
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

func (*schemaTable) optimize(_ Mode, index []string) (Cost, any) {
	if index == nil {
		return 1000, nil
	}
	return impossible, nil
}

func (*schemaTable) setApproach([]string, any, QueryTran) {
}

func (st *schemaTable) lookupCost() Cost {
	// no indexes, so lookups only from TempIndex
	return lookupCost(st.rowSize())
}

func (*schemaTable) Lookup(*Thread, []string, []string) Row {
	assert.ShouldNotReachHere()
	return nil
}

func (*schemaTable) Select([]string, []string) {
	assert.ShouldNotReachHere()
}

//-------------------------------------------------------------------

type Tables struct {
	schemaTable
	info []*meta.Info
	i    int
}

func (*Tables) String() string {
	return "tables"
}

func (ts *Tables) Transform() Query {
	return ts
}

func (*Tables) Keys() [][]string {
	return [][]string{{"table"}, {"tablename"}}
}

var tablesFields = [][]string{{"table", "tablename", "nrows", "totalsize"}}

func (*Tables) Columns() []string {
	return tablesFields[0]
}

func (*Tables) Header() *Header {
	return NewHeader(tablesFields, tablesFields[0])
}

func (ts *Tables) Nrows() int {
	ts.ensure()
	return len(ts.info)
}

func (ts *Tables) SetTran(tran QueryTran) {
	ts.tran = tran
	ts.info = nil
}

func (ts *Tables) Rewind() {
	ts.i = -1
	ts.state = rewound
}

func (ts *Tables) Get(_ *Thread, dir Dir) Row {
	ts.ensure()
	if ts.state == eof {
		return nil
	}
	if dir == Next {
		if ts.state == rewound {
			ts.i = -1
		}
		ts.i++
	} else { // Prev
		if ts.state == rewound {
			ts.i = len(ts.info)
		}
		ts.i--
	}
	if ts.i < 0 || len(ts.info) <= ts.i {
		ts.state = eof
		return nil
	}
	ts.state = within
	return ts.row(ts.info[ts.i])
}

func (*Tables) row(info *meta.Info) Row {
	var rb RecordBuilder
	rb.Add(SuStr(info.Table))
	rb.Add(SuStr(info.Table)) // tablename
	rb.Add(IntVal(info.Nrows).(Packable))
	rb.Add(Int64Val(int64(info.Size)).(Packable))
	rec := rb.Build()
	return Row{DbRec{Record: rec}}
}

func (ts *Tables) ensure() {
	if ts.info != nil {
		return
	}
	ts.info = ts.tran.GetAllInfo()

	cols := Columns{}
	cols.SetTran(ts.tran)
	ncols := cols.Nrows()

	idxs := Indexes{}
	idxs.SetTran(ts.tran)
	nidxs := idxs.Nrows()

	views := Views{}
	views.SetTran(ts.tran)
	nviews := views.Nrows()

	ts.info = append(ts.info,
		// +4 for tables, columns, indexes, views
		&meta.Info{Table: "tables", Nrows: len(ts.info) + 4},
		&meta.Info{Table: "columns", Nrows: ncols},
		&meta.Info{Table: "indexes", Nrows: nidxs},
		&meta.Info{Table: "views", Nrows: nviews},
	)
	sort.Slice(ts.info,
		func(i, j int) bool { return ts.info[i].Table < ts.info[j].Table })
}

//-------------------------------------------------------------------

// TablesLookup is used to optimize lookups on tables.
// It is inserted by Where Transform.
type TablesLookup struct {
	Tables
	table string
}

func NewTablesLookup(tran QueryTran, table string) *TablesLookup {
	tl := TablesLookup{table: table}
	tl.tran = tran
	return &tl
}

func (tl *TablesLookup) String() string {
	return "tables(" + tl.table + ")"
}

func (*TablesLookup) Nrows() int {
	return 1
}

func (tl *TablesLookup) Get(*Thread, Dir) Row {
	if tl.state != eof {
		tl.state = eof
		switch tl.table {
		case "tables", "columns", "indexes", "views":
			var rb RecordBuilder
			rb.Add(SuStr(tl.table))
			rb.Add(SuStr(tl.table)) // tablename
			rec := rb.Build()
			return Row{DbRec{Record: rec}}
		default:
			ti := tl.tran.GetInfo(tl.table)
			if ti != nil {
				return tl.row(ti)
			}
		}
	}
	return nil
}

//-------------------------------------------------------------------

type Columns struct {
	schemaTable
	schema []*meta.Schema
	si     int
	ci     int
}

func (*Columns) String() string {
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

func (cs *Columns) Nrows() int {
	cs.ensure()
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
	cs.tran = tran
	cs.schema = nil
}

func (cs *Columns) Rewind() {
	cs.si = -1
	cs.state = rewound
}

func (cs *Columns) Get(_ *Thread, dir Dir) Row {
	cs.ensure()
	if cs.state == eof {
		return nil
	}
	var col string
	var fld int
	for {
		if dir == Next {
			if cs.state == rewound {
				cs.si, cs.ci = 0, -1
			}
			cs.ci++
			for cs.ci >= len(cs.schema[cs.si].Columns)+len(cs.schema[cs.si].Derived) {
				cs.si++
				if cs.si >= len(cs.schema) {
					cs.state = eof
					return nil
				}
				cs.ci = 0
			}
		} else { // Prev
			if cs.state == rewound {
				cs.si = len(cs.schema)
				cs.ci = 0
			}
			cs.ci--
			for cs.ci < 0 {
				cs.si--
				if cs.si < 0 {
					cs.state = eof
					return nil
				}
				cs.ci = len(cs.schema[cs.si].Columns) + len(cs.schema[cs.si].Derived) - 1
			}
		}
		col, fld = columnOrDerived(cs.schema[cs.si], cs.ci)
		if col != "-" {
			break
		}
	}
	cs.state = within
	schema := cs.schema[cs.si]
	var rb RecordBuilder
	rb.Add(SuStr(schema.Table))
	rb.Add(SuStr(col))
	rb.Add(IntVal(fld).(Packable))
	rec := rb.Build()
	return Row{DbRec{Record: rec}}
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
	cs.schema = append(cs.schema,
		&meta.Schema{Schema: schema.Schema{Table: "tables", Columns: tablesFields[0]}},
		&meta.Schema{Schema: schema.Schema{Table: "columns", Columns: columnsFields[0]}},
		&meta.Schema{Schema: schema.Schema{Table: "indexes", Columns: indexesFields[0]}},
		&meta.Schema{Schema: schema.Schema{Table: "views", Columns: viewsFields[0]}},
	)
	sort.Slice(cs.schema,
		func(i, j int) bool { return cs.schema[i].Table < cs.schema[j].Table })
}

//-------------------------------------------------------------------
// note: indexes does not include tables, columns, indexes, views

type Indexes struct {
	schemaTable
	schema []*meta.Schema
	si     int
	ci     int
}

func (*Indexes) String() string {
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

func (is *Indexes) Nrows() int {
	is.ensure()
	n := 0
	for _, schema := range is.schema {
		n += len(schema.Indexes)
	}
	return n
}

func (is *Indexes) SetTran(tran QueryTran) {
	is.tran = tran
	is.schema = nil
}

func (is *Indexes) Rewind() {
	is.state = rewound
}

func (is *Indexes) Get(_ *Thread, dir Dir) Row {
	is.ensure()
	if is.state == eof {
		return nil
	}
	if dir == Next {
		if is.state == rewound {
			is.si, is.ci = 0, 0
		} else {
			is.ci++
			if is.ci >= len(is.schema[is.si].Indexes) {
				is.si++
				if is.si >= len(is.schema) {
					is.state = eof
					return nil
				}
				is.ci = 0
			}
		}
	} else { // Prev
		if is.state == rewound {
			is.si = len(is.schema) - 1
			is.ci = len(is.schema[is.si].Indexes) - 1
		} else {
			is.ci--
			if is.ci < 0 {
				is.si--
				if is.si < 0 {
					is.state = eof
					return nil
				}
				is.ci = len(is.schema[is.si].Indexes) - 1
			}
		}
	}
	is.state = within
	schema := is.schema[is.si]
	var rb RecordBuilder
	rb.Add(SuStr(schema.Table))
	idx := schema.Indexes[is.ci]
	rb.Add(SuStr(str.Join(",", idx.Columns)))
	switch idx.Mode {
	case 'k':
		rb.Add(True.(Packable))
	case 'i', 'I':
		rb.Add(False.(Packable))
	case 'u', 'U':
		rb.Add(SuStr("u"))
	default:
		assert.ShouldNotReachHere()
	}
	if idx.Fk.Table != "" {
		rb.Add(SuStr(idx.Fk.Table))
		rb.Add(SuStr(str.Join(",", idx.Fk.Columns)))
		rb.Add(SuInt(int(idx.Fk.Mode)))
	}
	rec := rb.Build()
	return Row{DbRec{Record: rec}}
}

func (is *Indexes) ensure() {
	if is.schema != nil {
		return
	}
	is.schema = is.tran.GetAllSchema()
	sort.Slice(is.schema,
		func(i, j int) bool { return is.schema[i].Table < is.schema[j].Table })
}

//-------------------------------------------------------------------

type Views struct {
	schemaTable
	views []string
	i     int
}

func (*Views) String() string {
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

func (vs *Views) Nrows() int {
	vs.ensure()
	return len(vs.views) / 2
}

func (vs *Views) SetTran(tran QueryTran) {
	vs.tran = tran
	vs.views = nil
}

func (vs *Views) Rewind() {
	vs.i = -2
	vs.state = rewound
}

func (vs *Views) Get(_ *Thread, dir Dir) Row {
	vs.ensure()
	if vs.state == eof {
		return nil
	}
	if dir == Next {
		if vs.state == rewound {
			vs.i = -2
		}
		vs.i += 2
	} else { // Prev
		if vs.state == rewound {
			vs.i = len(vs.views)
		}
		vs.i -= 2
	}
	if vs.i < 0 || len(vs.views) <= vs.i {
		vs.state = eof
		return nil
	}
	vs.state = within
	var rb RecordBuilder
	rb.Add(SuStr(vs.views[vs.i]))   // name
	rb.Add(SuStr(vs.views[vs.i+1])) // definition
	rec := rb.Build()
	return Row{DbRec{Record: rec}}
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

//-------------------------------------------------------------------

type History struct {
	schemaTable
	off uint64
}

func (*History) String() string {
	return "History"
}

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

func (his *History) Nrows() int {
	return 1000 // ???
}

func (his *History) SetTran(tran QueryTran) {
	his.tran = tran
}

func (his *History) Rewind() {
	his.state = rewound
}

func (his *History) Get(_ *Thread, dir Dir) Row {
	if his.state == eof {
		return nil
	}
	if his.state == rewound {
		his.off = 0
	}
	var state *db19.DbState
	if dir == Next {
		state = db19.NextState(his.tran.GetStore(), his.off)
	} else {
		state = db19.PrevState(his.tran.GetStore(), his.off)
	}
	if state == nil {
		his.state = eof
		return nil
	}
	his.state = within
	his.off = state.Off
	var rb RecordBuilder
	rb.Add(SuDateFromUnixMilli(state.Asof))
	rec := rb.Build()
	return Row{DbRec{Record: rec}}
}
