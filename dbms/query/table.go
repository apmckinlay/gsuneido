// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/slc"
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
	queryBase
	tran        QueryTran
	iter        index.IndexIter
	info        *meta.Info
	schema      *Schema
	name        string
	allKeys     [][]string
	index       []string // index that will be used to access the data
	sels        Sels
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
		idxs = append(idxs, ix.Fields)
		if ix.Mode == 'k' {
			keys = append(keys, ix.Columns)
			if len(ix.Columns) == 0 {
				tbl.singleton = true
			}
		}
	}
	tbl.indexes = idxs
	tbl.allKeys = keys
	tbl.keys = minimizeKeys(keys)
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

func (tbl *Table) Fixed() Fixed {
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

const ( // ???
	tableFast  = 150
	tableSlow  = 300
	tableLarge = 4_000_000
	colsBias   = 5
)

func (tbl *Table) optimize(mode Mode, req Require) (Cost, Cost, any) {
	if tbl.singleton {
		return tbl.costFor(tbl.indexes[0], req.frac, 0)
	}
	switch req.Use() {
	case ReqUnordered:
		return tbl.costFor(tbl.indexes[0], req.frac, 0)
	case ReqOrdered:
		idxi := tbl.indexFor(req.cols)
		if idxi == -1 {
			return impossible, impossible, nil
		}
		return tbl.costFor(tbl.indexes[idxi], req.frac, 0)
	case ReqGrouped:
		best := newBestIndex()
		for _, idx := range tbl.indexes {
			if set.StartsWithSet(idx, req.cols) {
				f, v, _ := tbl.costFor(idx, req.frac, req.nlookups)
				best.update(idx, f, v)
			}
		}
		if best.index == nil {
			return impossible, impossible, nil
		}
		return best.fixcost, best.varcost, tableApproach{index: best.index}
	case ReqLookup:
		best := newBestIndex()
		for idxi, idx := range tbl.indexes {
			if set.Subset(req.cols, idx) {
				varcost := Cost(req.nlookups) * tbl.lookupCostFor(idxi)
				best.update(idx, 0, varcost)
			}
		}
		if best.index == nil {
			return impossible, impossible, nil
		}
		return 0, best.varcost, tableApproach{index: best.index}
	}
	panic("unreachable")
}

func (tbl *Table) costFor(index []string, frac float32, nlookups int32) (Cost, Cost, any) {
	rowCost := tableFast
	if tbl.info.Size > tableLarge && !slices.Equal(index, tbl.indexes[0]) {
		rowCost = tableSlow
	}
	rowCost += len(index) * colsBias
	varcost := tbl.info.Nrows * rowCost
	result := Cost(float64(frac) * float64(varcost))
	if nlookups > 0 {
		idxi := slc.IndexFn(tbl.indexes, index, slices.Equal)
		if idxi != -1 {
			result += Cost(nlookups) * tbl.lookupCostFor(idxi)
		}
	}
	return 0, result, tableApproach{index: index}
}

// find an index that satisfies the required order
func (tbl *Table) indexFor(order []string) int {
	for i, index := range tbl.indexes {
		if slc.HasPrefix(index, order) {
			return i
		}
	}
	return -1 // not found
}

func (tbl *Table) setApproach(_ Require, approach any, _ QueryTran) {
	tbl.SetIndex(approach.(tableApproach).index)
}

func (tbl *Table) SetIndex(index []string) {
	if tbl.singleton {
		index = tbl.allKeys[0]
	}
	tbl.index = index
	tbl.iIndex = tbl.indexi(index)
	tbl.indexEncode = tbl.IndexEncodes(index)
	IdxUse(tbl.name, tbl.index)
}

// IndexEncodes returns whether the index key is encoded
// (multi-field or unique with Fields2)
func (tbl *Table) IndexEncodes(index []string) bool {
	return len(index) > 1 ||
		!slc.ContainsFn(tbl.allKeys, index, set.Equal[string])
}

func (tbl *Table) indexi(index []string) int {
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

func (tbl *Table) lookupCostFor(i int) Cost {
	var levels int
	if tbl.info.Indexes == nil { // tests
		levels = 1
		if tbl.info.Nrows > 100 {
			levels = 2
		}
	} else {
		levels = tbl.info.Indexes[i].BtreeLevels()
	}
	return lookupCost(levels)
}

func lookupCost(levels int) Cost {
	return levels*800 - 300 // empirical
}

// execution --------------------------------------------------------

func (tbl *Table) Lookup(_ *Thread, sels Sels) Row {
	assert.That(!selConflict(tbl.header.Columns, sels))
	tbl.nlooks++
	key := ""
	if !tbl.singleton {
		ix := &tbl.schema.Indexes[tbl.iIndex]
		key = selOrg(tbl.indexEncode, ix.Fields, sels, true)
		if len(ix.Ixspec.Fields2) > 0 && key == "" {
			fullFields := set.Union(ix.Fields, ix.BestKey)
			key = selOrg(true, fullFields, sels, true)
		}
	}
	row := tbl.LookupRaw(key)
	if row == nil || !singletonFilter(tbl.header, row, sels) {
		return nil
	}
	return row
}

func (tbl *Table) LookupRaw(key string) Row {
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
		if tbl.singleton &&
			!singletonFilter(tbl.header, row, tbl.sels) {
			return nil
		}
		tbl.ngets++
		return row
	}
}

func (tbl *Table) Select(sels Sels) {
	tbl.nsels++
	if tbl.singleton {
		tbl.sels = sels
		tbl.ensureIter().Range(iface.All)
		return
	}
	if sels == nil { // clear select
		tbl.ensureIter().Range(iface.All)
		return
	}
	assert.That(!selConflict(tbl.header.Columns, sels))
	org, end := selKeys(tbl.indexEncode, tbl.index, sels)
	tbl.SelectRaw(org, end)
}

func selKeys(encode bool, dstCols []string, sels Sels) (string, string) {
	if len(dstCols) == 0 {
		return ixkey.Min, ixkey.Max
	}
	if !encode {
		assert.That(len(dstCols) == 1)
		org := sels.MustGet(dstCols[0])
		end := org + "\x00"
		return org, end
	}
	end := selEnd(dstCols, sels)
	org := trim(end) // selOrg(true, dstCols, srcCols, vals)
	return org, end
}

func selEnd(dstCols []string, sels Sels) string {
	enc := ixkey.Encoder{}
	data := false
	for _, col := range dstCols {
		val, ok := sels.Get(col)
		if !ok {
			break
		}
		enc.Add(val)
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

func selOrg(encode bool, dstCols []string, sels Sels, full bool) string {
	if !encode {
		if len(dstCols) == 0 {
			return ""
		}
		assert.That(len(dstCols) == 1)
		return sels.MustGet(dstCols[0])
	}
	enc := ixkey.Encoder{}
	data := false
	for _, col := range dstCols {
		val, ok := sels.Get(col)
		if !ok {
			if full {
				panic("selOrg not full")
			}
			break
		}
		enc.Add(val)
		data = true
	}
	assert.That(data)
	return enc.String()
}

func (tbl *Table) SelectRaw(org, end string) {
	tbl.ensureIter().Range(index.Range{Org: org, End: end})
}

func (tbl *Table) SelectSkipScan(prefixRng, suffixRng iface.Range, skipStart int) {
	tbl.ensureIter().SkipScan(prefixRng, suffixRng, skipStart)
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
	tbl.iter.Rewind()
	return rows
}

func (tbl *Table) RangeFrac(index []string, sels Sels) float64 {
	iIndex := tbl.indexi(index)
	encode := len(index) > 1 ||
		!slc.ContainsFn(tbl.allKeys, index, set.Equal[string])
	org, end := selKeys(encode, index, sels)
	return tbl.tran.RangeFrac(tbl.name, iIndex, org, end)
}
