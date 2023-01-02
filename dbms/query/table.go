// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/iterator"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/trace"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/exp/slices"
)

func NewTable(t QueryTran, name string) Query {
	var tbl Query
	switch name {
	case "tables":
		tbl = &Tables{}
	case "columns":
		tbl = &Columns{}
	case "indexes":
		tbl = &Indexes{}
	case "views":
		tbl = &Views{}
	case "history":
		tbl = &History{}
	default:
		tbl = &Table{name: name}
	}
	tbl.SetTran(t)
	return tbl
}

type Table struct {
	cache
	name      string
	hdr       *runtime.Header
	indexes   [][]string
	keys      [][]string
	primary   [][]string
	singleton bool
	tran      QueryTran
	schema    *Schema
	info      *meta.Info
	// index is the index that will be used to access the data.
	index       []string
	iIndex      int
	indexEncode bool
	iter        *index.OverIter
	selcols     []string
	selvals     []string
}

type tableApproach struct {
	index []string
}

func (tbl *Table) String() string {
	if tbl.index == nil {
		return tbl.name
	}
	return tbl.name + "^" + str.Join("(,)", tbl.index)
}

func (tbl *Table) SetTran(t QueryTran) {
	tbl.tran = t
	tbl.schema = t.GetSchema(tbl.name)
	if tbl.schema == nil {
		panic("nonexistent table: " + tbl.name)
	}
	tbl.info = t.GetInfo(tbl.name)

	cols := make([]string, 0, len(tbl.schema.Columns)+len(tbl.schema.Derived))
	for _, col := range tbl.schema.Columns {
		if col != "-" {
			cols = append(cols, col)
		}
	}
	for _, col := range tbl.schema.Derived {
		cols = append(cols, str.UnCapitalize(col))
	}
	tbl.hdr = runtime.NewHeader([][]string{tbl.schema.Columns}, cols)

	idxs := make([][]string, 0, len(tbl.schema.Indexes))
	keys := make([][]string, 0, 1)
	for _, ix := range tbl.schema.Indexes {
		idxs = append(idxs, ix.Columns)
		if ix.Mode == 'k' {
			keys = append(keys, ix.Columns)
			if len(ix.Columns) == 0 {
				tbl.singleton = true
			}
		}
	}
	tbl.indexes = idxs
	tbl.keys = keys
}

func (tbl *Table) Columns() []string {
	return tbl.hdr.Columns
}

func (tbl *Table) Indexes() [][]string {
	return tbl.indexes
}

func (tbl *Table) Keys() [][]string {
	if tbl.primary == nil {
		tbl.primary = withoutDupsOrSupersets(tbl.keys)
	}
	return tbl.primary
}

func (*Table) Ordering() []string {
	return nil
}

func (tbl *Table) Nrows() (int, int) {
	return tbl.info.Nrows, tbl.info.Nrows
}

func (tbl *Table) rowSize() int {
	if tbl.info.Nrows == 0 {
		return 0
	}
	return int(tbl.info.Size) / tbl.info.Nrows
}

func (tbl *Table) Transform() Query {
	return tbl
}

func (tbl *Table) Fixed() []Fixed {
	return nil
}

func (tbl *Table) Updateable() string {
	return tbl.name
}

func (tbl *Table) SingleTable() bool {
	switch tbl.name {
	case "tables", "columns", "indexes":
		return false
	}
	return true
}

func (tbl *Table) optimize(_ Mode, index []string, frac float64) (Cost, Cost, any) {
	if index == nil {
		index = tbl.schema.Indexes[0].Columns
	} else if !tbl.singleton {
		i := tbl.indexFor(index)
		if i < 0 {
			return impossible, impossible, nil
		}
		index = tbl.indexes[i]
	}
	varcost := tbl.info.Nrows * 250 // empirical
	trace.QueryOpt.Println("optimize", tbl.name, index, frac, "=", Cost(frac*float64(varcost)))
	return 0, Cost(frac * float64(varcost)), tableApproach{index: index}
}

// find an index that satisfies the required order
func (tbl *Table) indexFor(order []string) int {
	for i, ix := range tbl.indexes {
		if slc.HasPrefix(ix, order) {
			return i
		}
	}
	return -1 // not found
}

func (tbl *Table) setApproach(_ []string, _ float64, approach any, _ QueryTran) {
	tbl.setIndex(approach.(tableApproach).index)
}

func (tbl *Table) setIndex(index []string) {
	if tbl.singleton {
		index = tbl.keys[0]
	}
	tbl.index = index
	tbl.iIndex = slc.IndexFn(tbl.indexes, tbl.index, slices.Equal[string])
	assert.Msg("setIndex", tbl.name, index).That(tbl.iIndex >= 0)
	tbl.indexEncode = len(tbl.index) > 1 ||
		!slc.ContainsFn(tbl.keys, tbl.index, set.Equal[string])
}

func (tbl *Table) lookupCost() Cost {
	var levels int
	if tbl.info.Indexes == nil { // tests
		levels = 1
		if tbl.info.Nrows > 100 {
			levels = 2
		}
	} else {
		levels = tbl.info.Indexes[tbl.iIndex].BtreeLevels()
	}
	return levels*800 - 300
}

// execution --------------------------------------------------------

func (tbl *Table) Lookup(_ *runtime.Thread, cols, vals []string) runtime.Row {
	key := selOrg(tbl.indexEncode, tbl.index, cols, vals, true)
	return tbl.lookup(key)
}

func (tbl *Table) lookup(key string) runtime.Row {
	rec := tbl.tran.Lookup(tbl.name, tbl.iIndex, key)
	if rec == nil {
		return nil
	}
	return runtime.Row{*rec}
}

func (tbl *Table) Header() *runtime.Header {
	return tbl.hdr
}

func (tbl *Table) Output(th *runtime.Thread, rec runtime.Record) {
	tbl.tran.Output(th, tbl.name, rec)
}

func (tbl *Table) Rewind() {
	if tbl.iter != nil {
		tbl.iter.Rewind()
	}
}

func (tbl *Table) Get(_ *runtime.Thread, dir runtime.Dir) runtime.Row {
	tbl.ensureIter()
	if dir == runtime.Prev {
		tbl.iter.Prev(tbl.tran)
	} else {
		tbl.iter.Next(tbl.tran)
	}
	if tbl.iter.Eof() {
		return nil
	}
	_, off := tbl.iter.Cur()
	rec := tbl.tran.GetRecord(off)
	row := runtime.Row{runtime.DbRec{Record: rec, Off: off}}
	if tbl.singleton && !singletonFilter(tbl.hdr, row, tbl.selcols, tbl.selvals) {
		return nil
	}
	return row
}

func (tbl *Table) Select(cols, vals []string) {
	if tbl.singleton {
		tbl.selcols, tbl.selvals = cols, vals
		tbl.ensureIter().Range(iterator.All)
		return
	}
	if cols == nil && vals == nil { // clear select
		tbl.iter.Range(iterator.All)
		return
	}
	org, end := selKeys(tbl.indexEncode, tbl.index, cols, vals)
	tbl.SelectRaw(org, end)
}

func selKeys(encode bool, dstCols, srcCols, vals []string) (string, string) {
	if len(dstCols) == 0 {
		return ixkey.Min, ixkey.Max
	}
	if !encode {
		assert.That(len(dstCols) == 1)
		org := selGet(dstCols[0], srcCols, vals)
		end := org + "\x00"
		return org, end
	}
	end := selEnd(dstCols, srcCols, vals)
	org := trim(end) // selOrg(true, dstCols, srcCols, vals)
	return org, end
}

func selGet(col string, cols, vals []string) string {
	i := slices.Index(cols, col)
	assert.Msg("selGet", col, "NOT IN", cols).That(i != -1)
	return vals[i]
}

func selEnd(dstCols, srcCols, vals []string) string {
	enc := ixkey.Encoder{}
	data := false
	for _, col := range dstCols {
		i := slices.Index(srcCols, col)
		if i == -1 {
			break
		}
		enc.Add(vals[i])
		data = true
	}
	assert.Msg("selEnd no data").That(data)
	enc.Add(ixkey.Max)
	return enc.String()
}

func trim(end string) string {
	n := len(end) - len(ixkey.Max)
	org := end[:n]
	for n >= 2 && org[n-2:n] == ixkey.Sep {
		n -= 2
	}
	return org[:n]
}

func selOrg(encode bool, dstCols, srcCols, vals []string, full bool) string {
	if !encode {
		assert.That(len(dstCols) == 1)
		return selGet(dstCols[0], srcCols, vals)
	}
	enc := ixkey.Encoder{}
	data := false
	for _, col := range dstCols {
		i := slices.Index(srcCols, col)
		if i == -1 {
			if full {
				panic("selOrg not full")
			}
			break
		}
		enc.Add(vals[i])
		data = true
	}
	assert.Msg("selOrg no data").That(data)
	return enc.String()
}

func (tbl *Table) SelectRaw(org, end string) {
	tbl.ensureIter().Range(index.Range{Org: org, End: end})
}

func (tbl *Table) ensureIter() *index.OverIter {
	if tbl.iter == nil {
		tbl.iter = index.NewOverIter(tbl.name, tbl.iIndex)
	}
	return tbl.iter
}

// func packedListStr(sel []string) string { // for debugging
// 	p := ""
// 	sep := ""
// 	for _, s := range sel {
// 		p += sep + runtime.Display(&runtime.Thread{}, runtime.Unpack(s))
// 		sep = ","
// 	}
// 	return p
// }
