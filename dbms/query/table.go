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
	useTempIndex
	name   string
	t      QueryTran
	schema *Schema
	info   *meta.Info
	// index is the index that will be used to access the data.
	// It is set by optimize.
	index     []string
	tempIndex []string
}

func (tbl *Table) String() string {
	if tbl.index == nil {
		return tbl.name
	}
	s := tbl.name + "^" + str.Join("(,)", tbl.index)
	if tbl.tempIndex != nil {
		s += " TEMPINDEX" + str.Join("(,)", tbl.tempIndex)
	}
	return s
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
}

func (tbl *Table) Columns() []string {
	allcols := make([]string, 0, len(tbl.schema.Columns)+len(tbl.schema.Derived))
	allcols = append(allcols, tbl.schema.Columns...)
	allcols = append(allcols, tbl.schema.Derived...)
	return allcols
}

func (tbl *Table) Indexes() [][]string {
	idxs := make([][]string, 0, 1)
	for _, ix := range tbl.schema.Indexes {
		idxs = append(idxs, ix.Columns)
	}
	return idxs
}

func (tbl *Table) Keys() [][]string {
	keys := make([][]string, 0, 1)
	for _, ix := range tbl.schema.Indexes {
		if ix.Mode == 'k' {
			keys = append(keys, ix.Columns)
		}
	}
	return keys
}

func (tbl *Table) nrows() int {
	return tbl.info.Nrows
}

func (tbl *Table) dataSize() int {
	return int(tbl.info.Size)
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

func (tbl *Table) optimize(mode Mode, index []string, act action) Cost {
	i := 0
	if index == nil {
		index = tbl.schema.Indexes[0].Columns
	} else if i = tbl.findIndex(index); i < 0 && mode == cursorMode {
		return impossible
	}
	if act == freeze {
		if i >= 0 {
			tbl.index = index
		} else {
			tbl.index = tbl.schema.Indexes[0].Columns
			tbl.tempIndex = index
		}
	}
	indexReadCost := tbl.info.Nrows * btree.EntrySize
	dataReadCost := int(tbl.info.Size)
	cost := indexReadCost + dataReadCost
	if i < 0 {
		indexWriteCost := indexReadCost // ???
		cost += indexReadCost + dataReadCost + indexWriteCost
	}
	return cost
}

func (tbl *Table) findIndex(index []string) int {
	for i := range tbl.schema.Indexes {
		if str.Equal(index, tbl.schema.Indexes[i].Columns) {
			return i
		}
	}
	return -1
}

func (tbl *Table) addTempIndex(tran QueryTran) Query {
	return addTempIndex(tbl, tran)
}
