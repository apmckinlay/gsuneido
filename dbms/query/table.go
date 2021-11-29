// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/index/iterator"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/setset"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/strs"
	"github.com/apmckinlay/gsuneido/util/strss"
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
	default:
		tbl = &Table{name: name}
	}
	tbl.SetTran(t)
	return tbl
}

type Table struct {
	cache
	name      string
	columns   []string
	indexes   [][]string
	keys      [][]string
	singleton bool
	tran      QueryTran
	schema    *Schema
	info      *meta.Info
	// index is the index that will be used to access the data.
	index       []string
	iIndex      int
	indexEncode bool
	iter        *index.OverIter
}

type tableApproach struct {
	index []string
}

func (tbl *Table) String() string {
	if tbl.index == nil {
		return tbl.name
	}
	return tbl.name + "^" + strs.Join("(,)", tbl.index)
}

func (tbl *Table) SetTran(t QueryTran) {
	tbl.tran = t
	tbl.schema = t.GetSchema(tbl.name)
	if tbl.schema == nil {
		panic("nonexistent table: " + tbl.name)
	}
	tbl.info = t.GetInfo(tbl.name)

	cols := make([]string, 0, len(tbl.schema.Columns)+len(tbl.schema.Derived))
	// cols = append(cols, tbl.schema.Columns...)
	for _, col := range tbl.schema.Columns {
		if col != "-" {
			cols = append(cols, col)
		}
	}
	for _, col := range tbl.schema.Derived {
		cols = append(cols, str.UnCapitalize(col))
	}
	tbl.columns = cols

	idxs := make([][]string, 0, len(tbl.schema.Indexes)-1)
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
	return tbl.columns
}

func (tbl *Table) Indexes() [][]string {
	return tbl.indexes
}

func (tbl *Table) Keys() [][]string {
	return tbl.keys
}

func (*Table) Ordering() []string {
	return nil
}

func (tbl *Table) Nrows() int {
	return tbl.info.Nrows
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

func (tbl *Table) optimize(_ Mode, index []string) (Cost, interface{}) {
	if index == nil {
		index = tbl.schema.Indexes[0].Columns
	} else if !tbl.singleton {
		i := tbl.indexFor(index)
		if i < 0 {
			return impossible, nil
		}
		index = tbl.indexes[i]
	}
	indexReadCost := tbl.info.Nrows * btree.EntrySize
	dataReadCost := int(tbl.info.Size)
	return indexReadCost + dataReadCost, tableApproach{index: index}
}

// find an index that satisfies the required order
func (tbl *Table) indexFor(order []string) int {
	for i, ix := range tbl.indexes {
		if strs.HasPrefix(ix, order) {
			return i
		}
	}
	return -1 // not found
}

func (tbl *Table) setApproach(_ []string, approach interface{}, _ QueryTran) {
	tbl.setIndex(approach.(tableApproach).index)
}

func (tbl *Table) setIndex(index []string) {
	tbl.index = index
	tbl.iIndex = strss.Index(tbl.indexes, tbl.index)
	tbl.indexEncode = len(tbl.index) > 1 || !setset.Contains(tbl.keys, tbl.index)
}

// lookupCost returns the cost of one lookup
func (tbl *Table) lookupCost() Cost {
	return lookupCost(tbl.rowSize())
}

func lookupCost(rowSize int) Cost {
	// average node size is 2/3 of max, on average we read half = 1/3
	nodeScan := btree.MaxNodeSize / 3
	return (nodeScan * btree.TreeHeight) + rowSize
}

// execution --------------------------------------------------------

func (tbl *Table) Lookup(cols, vals []string) runtime.Row {
	key := selOrg(tbl.indexEncode, tbl.index, cols, vals)
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
	physical := [][]string{tbl.schema.Columns}
	return runtime.NewHeader(physical, tbl.columns)
}

func (tbl *Table) Output(rec runtime.Record) {
	tbl.tran.Output(tbl.name, rec)
}

func (tbl *Table) Rewind() {
	if tbl.iter != nil {
		tbl.iter.Rewind()
	}
}

func (tbl *Table) Get(dir runtime.Dir) runtime.Row {
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
	return runtime.Row{runtime.DbRec{Record: rec, Off: off}}
}

func (tbl *Table) Select(cols, vals []string) {
	if cols == nil && vals == nil { // clear select
		tbl.iter.Range(iterator.All)
		return
	}
	org, end := selKeys(tbl.indexEncode, tbl.index, cols, vals)
	tbl.SelectRaw(org, end)
}

func selKeys(encode bool, dstCols, srcCols, vals []string) (string, string) {
	if !encode {
		org := selGet(dstCols[0], srcCols, vals)
		end := org + "\x00"
		return org, end
	}
	end := selEnd(dstCols, srcCols, vals)
	org := trim(end)
	return org, end
}

func selGet(col string, cols, vals []string) string {
	i := strs.Index(cols, col)
	assert.Msg("selGet", col, "NOT IN", cols).That(i != -1)
	return vals[i]
}

func selEnd(dstCols, srcCols, vals []string) string {
	enc := ixkey.Encoder{}
	prefix := true
	for _, col := range dstCols {
		i := strs.Index(srcCols, col)
		if i != -1 {
			assert.Msg("selEncode").That(prefix)
			enc.Add(vals[i])
		} else {
			prefix = false
			enc.Add("")
		}
	}
	enc.Add(ixkey.Max)
	return enc.String()
}

func trim(end string) string {
	n := len(end)
	org := end[:n-len(ixkey.Max)]
	n -= len(ixkey.Max)
	for n >= 2 && org[n-2:n] == ixkey.Sep {
		n -= 2
	}
	return org[:n]
}

func selOrg(encode bool, dstCols, srcCols, vals []string) string {
	if !encode {
		return selGet(dstCols[0], srcCols, vals)
	}
	enc := ixkey.Encoder{}
	prefix := true
	for _, col := range dstCols {
		i := strs.Index(srcCols, col)
		if i != -1 {
			assert.Msg("selEncode").That(prefix)
			enc.Add(vals[i])
		} else {
			prefix = false
		}
	}
	return enc.String()
}

func (tbl *Table) SelectRaw(org, end string) {
	tbl.ensureIter()
	tbl.iter.Range(index.Range{Org: org, End: end})
}

func (tbl *Table) ensureIter() {
	if tbl.iter == nil {
		tbl.iter = index.NewOverIter(tbl.name, tbl.iIndex)
	}
}
