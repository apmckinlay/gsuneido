// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/tsc"
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
	tran    QueryTran
	iter    index.IndexIter
	info    *meta.Info
	schema  *Schema
	name    string
	allKeys [][]string
	// index is the index that will be used to access the data.
	index   []string
	selcols []string
	selvals []string
	queryBase
	iIndex      int
	singleton   bool
	indexEncode bool
	cursorMode  bool
}

func (tbl *Table) isSingleton() bool {
	return tbl.singleton
}

func (tbl *Table) schemaIndexes() []Index {
	return tbl.schema.Indexes
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

func (tbl *Table) Name() string {
	return tbl.name
}

func (tbl *Table) SetTran(t QueryTran) {
	tbl.tran = t
	tbl.schema = t.GetSchema(tbl.name)
	if tbl.schema == nil {
		panic("nonexistent table: " + tbl.name)
	}
	tbl.info = t.GetInfo(tbl.name)
	tbl.rowSiz.Set(tbl.getRowSize())

	cols := make([]string, 0, len(tbl.schema.Columns)+len(tbl.schema.Derived))
	for _, col := range tbl.schema.Columns {
		if col != "-" {
			cols = append(cols, col)
		}
	}
	for _, col := range tbl.schema.Derived {
		cols = append(cols, str.UnCapitalize(col))
	}
	tbl.header = NewHeader([][]string{tbl.schema.Columns}, cols)

	idxs := make([][]string, 0, len(tbl.schema.Indexes))
	keys := make([][]string, 0, 1)
	for i := range tbl.schema.Indexes {
		ix := &tbl.schema.Indexes[i]
		idxs = append(idxs, ix.Columns)
		if ix.Mode == 'k' {
			keys = append(keys, ix.Columns)
			if len(ix.Columns) == 0 {
				tbl.singleton = true
			}
		}
	}
	tbl.indexes = idxs
	tbl.allKeys = keys
	tbl.keys = withoutDupsOrSupersets(keys)
	tbl.fast1.Set(isEmptyKey(tbl.keys))
}

func (tbl *Table) Nrows() (int, int) {
	return tbl.info.Nrows, tbl.info.Nrows
}

func (*Table) knowExactNrows() bool {
	return true
}

func (tbl *Table) getRowSize() int {
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
		index = tbl.indexes[0]
	} else if !tbl.singleton {
		i := tbl.indexFor(index)
		if i < 0 {
			return impossible, impossible, nil
		}
		index = tbl.indexes[i]
	}
	varcost := tbl.info.Nrows * 250 // empirical
	trace.QueryOpt.Println("Table optimize", tbl.name, index, frac, "=",
		Cost(frac*float64(varcost)))
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
	tbl.SetIndex(approach.(tableApproach).index)
}

func (tbl *Table) SetIndex(index []string) {
	if tbl.singleton {
		index = tbl.allKeys[0]
	}
	tbl.index = index
	tbl.iIndex = tbl.indexi(index)
	tbl.indexEncode = len(tbl.index) > 1 ||
		!slc.ContainsFn(tbl.allKeys, tbl.index, set.Equal[string])
	IdxUse(tbl.name, tbl.index)
}

// IndexCols returns the columns available in index entried.
// i.e. on non-unique indexes it includes the key fields added for uniqueness
// WARNING: for unique indexes, these columns are NOT sufficient for lookup
// because they exclude the key fields added when the index fields are empty
func (tbl *Table) IndexCols(index []string) []string {
	iIndex := tbl.indexi(index)
	return tbl.schema.Indexes[iIndex].Fields
}

func (tbl *Table) indexi(index []string) int {
	// WARNING: assumes tbl.indexes is parallel to schema.Indexes
	i := slc.IndexFn(tbl.indexes, index, slices.Equal)
	assert.That(i >= 0)
	return i
}

func (tbl *Table) lookupCost() Cost {
	var levels int
	if tbl.info.Indexes == nil { // tests
		levels = 1
		if tbl.info.Nrows > 100 {
			levels = 2
		}
	} else {
		levels = tbl.info.Indexes[0].BtreeLevels()
		//TODO should be the actual index but we don't have it
	}
	return lookupCost(levels)
}

func lookupCost(levels int) Cost {
	return levels*800 - 300 // empirical
}

// execution --------------------------------------------------------

func (tbl *Table) Lookup(_ *Thread, cols, vals []string) (row Row) {
	assert.That(!selConflict(tbl.header.Columns, cols, vals))
	tbl.nlooks++
	if tbl.singleton {
		// For singleton tables (empty key), any columns are acceptable
		// Singleton tables have at most one row, so we can use GetFilter
		// which already handles singleton filtering
		tbl.selcols, tbl.selvals = cols, vals
		tbl.ensureIter().Range(iface.All)
		tbl.Rewind()
		return tbl.GetFilter(Next, nil)
	}
	// For non-singleton tables, the lookup columns must match the index.
	ix := &tbl.schema.Indexes[tbl.iIndex]
	key := selOrg(tbl.indexEncode, ix.Fields, cols, vals, true)
	if len(ix.Ixspec.Fields2) > 0 && key == "" {
		// For unique indexes, if all index field values are empty,
		// Fields2 (from BestKey) is added to make the key unique
		fullFields := set.Union(ix.Fields, ix.BestKey)
		assert.That(set.Equal(fullFields, cols))
		key = selOrg(true, fullFields, cols, vals, true)
	} else {
		assert.That(set.Equal(ix.Fields, cols))
	}
	return tbl.lookup(key)
}

func (tbl *Table) lookup(key string) Row {
	rec := tbl.tran.Lookup(tbl.name, tbl.iIndex, key)
	if rec == nil {
		return nil
	}
	return Row{*rec}
}

func (tbl *Table) Output(th *Thread, rec Record) {
	tbl.tran.Output(th, tbl.name, rec)
}

func (tbl *Table) Rewind() {
	if tbl.iter != nil {
		tbl.iter.Rewind()
	}
}

func (tbl *Table) Get(_ *Thread, dir Dir) Row {
	return tbl.GetFilter(dir, nil)
}

func (tbl *Table) GetFilter(dir Dir, filter func(key string) bool) Row {
	defer func(t uint64) { tbl.tget += tsc.Read() - t }(tsc.Read())
	tbl.ensureIter()
	for {
		if dir == Prev {
			tbl.iter.Prev(tbl.tran)
		} else {
			tbl.iter.Next(tbl.tran)
		}
		if tbl.iter.Eof() {
			return nil
		}
		var off uint64
		if filter != nil {
			key, o := tbl.iter.Cur()
			if !filter(key) {
				continue
			}
			off = o
		} else {
			off = tbl.iter.CurOff()
		}
		rec := tbl.tran.GetRecord(off)
		row := Row{DbRec{Record: rec, Off: off}}
		if tbl.singleton && !singletonFilter(tbl.header, row, tbl.selcols, tbl.selvals) {
			return nil
		}
		tbl.ngets++
		return row
	}
}

func (tbl *Table) Select(cols, vals []string) {
	tbl.nsels++
	if tbl.singleton {
		tbl.selcols, tbl.selvals = cols, vals
		tbl.ensureIter().Range(iface.All)
		return
	}
	if cols == nil && vals == nil { // clear select
		tbl.ensureIter().Range(iface.All)
		return
	}
	assert.That(!selConflict(tbl.header.Columns, cols, vals))
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
	assert.That(i != -1)
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
	assert.That(data)
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
		if len(dstCols) == 0 {
			return ""
		}
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
	assert.That(data)
	return enc.String()
}

func (tbl *Table) SelectRaw(org, end string) {
	tbl.ensureIter().Range(index.Range{Org: org, End: end})
}

func (tbl *Table) ensureIter() index.IndexIter {
	if tbl.iter == nil {
		if tbl.cursorMode {
			tbl.iter = index.NewOverIter(tbl.name, tbl.iIndex)
		} else {
			tbl.iter = tbl.tran.IndexIter(tbl.name, tbl.iIndex)
		}
	}
	return tbl.iter
}

const maxSimple = 2000

func (tbl *Table) Simple(*Thread) []Row {
	// can't use info.Nrows because sizeTran overrides it for tests
	tbl.ensureIter()
	rows := make([]Row, 0, 8)
	for {
		tbl.iter.Next(tbl.tran)
		if tbl.iter.Eof() {
			break
		}
		off := tbl.iter.CurOff()
		rec := tbl.tran.GetRecord(off)
		row := Row{DbRec{Record: rec, Off: off}}
		rows = append(rows, row)
		assert.That(len(rows) < maxSimple)
	}
	return rows
}

func (tbl *Table) RangeFrac(index, cols, vals []string) float64 {
	iIndex := tbl.indexi(index)
	encode := len(index) > 1 ||
		!slc.ContainsFn(tbl.allKeys, index, set.Equal[string])
	org, end := selKeys(encode, index, cols, vals)
	return tbl.tran.RangeFrac(tbl.name, iIndex, org, end)
}
