// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/ssset"
	"github.com/apmckinlay/gsuneido/util/str"
)

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
	// It is set by optimize.
	index       []string
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
	return tbl.name + "^" + str.Join("(,)", tbl.index)
}

func (tbl *Table) Init() {
}

func (tbl *Table) SetTran(t QueryTran) {
	tbl.tran = t
	tbl.schema = t.GetSchema(tbl.name)
	if tbl.schema == nil {
		panic("nonexistent table: " + tbl.name)
	}
	tbl.info = t.GetInfo(tbl.name)

	cols := make([]string, 0, len(tbl.schema.Columns)+len(tbl.schema.Derived))
	cols = append(cols, tbl.schema.Columns...)
	cols = append(cols, tbl.schema.Derived...)
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

func (tbl *Table) nrows() int {
	return tbl.info.Nrows
}

func (tbl *Table) rowSize() int {
	return int(tbl.info.Size) / tbl.info.Nrows
}

func (tbl *Table) Transform() Query {
	return tbl
}

func (tbl *Table) Fixed() []Fixed {
	return nil
}

func (tbl *Table) Updateable() bool {
	return true
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
		if str.List(ix).HasPrefix(order) {
			return i
		}
	}
	return -1 // not found
}

func (tbl *Table) findIndex(index []string) int {
	for i, ix := range tbl.indexes {
		if str.List(index).Equal(ix) {
			return i
		}
	}
	return -1 // not found
}

func (tbl *Table) setApproach(_ []string, approach interface{}, _ QueryTran) {
	tbl.index = approach.(tableApproach).index
	tbl.indexEncode = len(tbl.index) > 1 || !ssset.Contains(tbl.keys, tbl.index)
}

// lookupCost returns the cost of one lookup
func (tbl *Table) lookupCost() Cost {
	// average node size is 2/3 of max, on average we read half = 1/3
	nodeScan := btree.MaxNodeSize / 3
	return (nodeScan * btree.TreeHeight) + tbl.rowSize()
}

// execution --------------------------------------------------------

func (tbl *Table) Lookup(key string) runtime.Row {
	iIndex := tbl.findIndex(tbl.index) //TODO cache
	rec := tbl.tran.Lookup(tbl.name, iIndex, key)
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

func (tbl *Table) Select(org, end string) {
	if end == "" {
		if tbl.indexEncode {
			end = org + ixkey.Sep + ixkey.Max
		} else {
			end = org + "\x00"
		}
	}
	tbl.ensureIter()
	tbl.iter.Range(index.Range{Org: org, End: end})
}

func (tbl *Table) ensureIter() {
	if tbl.iter == nil {
		tbl.iter = index.NewOverIter(tbl.name, tbl.index)
	}
}
