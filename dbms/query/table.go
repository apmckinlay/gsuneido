// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Table struct {
	cache
	name      string
	columns   []string
	indexes   [][]string
	keys      [][]string
	singleton bool
	t         QueryTran
	schema    *Schema
	info      *meta.Info
	// index is the index that will be used to access the data.
	// It is set by optimize.
	index []string
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
	tbl.t = t
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

func (tbl *Table) optimize(_ Mode, index []string) (Cost, interface{}) {
	index = tbl.findIndex(index)
	if index == nil {
		return impossible, nil
	}
	indexReadCost := tbl.info.Nrows * btree.EntrySize
	dataReadCost := int(tbl.info.Size)
	return indexReadCost + dataReadCost, tableApproach{index: index}
}

func (tbl *Table) findIndex(index []string) []string {
	if index == nil {
		return tbl.schema.Indexes[0].Columns
	}
	if tbl.singleton {
		return index
	}
	for i := range tbl.schema.Indexes {
		if str.Equal(index, tbl.schema.Indexes[i].Columns) {
			return index
		}
	}
	return nil
}

func (tbl *Table) setApproach(_ []string, approach interface{}, _ QueryTran) {
	tbl.index = approach.(tableApproach).index
}

// lookupCost returns the cost of one lookup
func (tbl *Table) lookupCost() Cost {
	// average node size is 2/3 of max, on average we read half = 1/3
	nodeScan := btree.MaxNodeSize / 3
	return (nodeScan * btree.TreeHeight) + tbl.rowSize()
}
