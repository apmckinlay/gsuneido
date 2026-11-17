// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dbg"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/slc"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/tsc"
)

func NewTable(t QueryTran, name string) *Table {
	tbl := &Table{name: name}
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
	req         Require
}

func (tbl *Table) Clone() *Table {
	tbl2 := *tbl
	return &tbl2
}

func (tbl *Table) isSingleton() bool {
	return tbl.singleton
}

func (tbl *Table) schemaIndexes() []Index {
	return tbl.schema.Indexes
}

// UniqueIndexes returns the columns of the unique ('u') indexes
func (tbl *Table) UniqueIndexes() [][]string {
	var result [][]string
	for i := range tbl.schema.Indexes {
		if tbl.schema.Indexes[i].Mode == 'u' {
			result = append(result, tbl.schema.Indexes[i].Columns)
		}
	}
	return result
}

type tableApproach struct {
	index []string
	mode  Mode
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
	schema := t.GetSchema(tbl.name)
	if schema == nil {
		panic("nonexistent table: " + tbl.name)
	}
	if tbl.tran != t {
		t.CheckPerm(tbl.name, PermRead)
	}
	tbl.tran = t
	// Update transactions can mutate info without changing its pointer.
	tbl.info = t.GetInfo(tbl.name)
	tbl.rowSiz.Set(tbl.getRowSize())
	if schema == tbl.schema {
		return
	}
	tbl.schema = schema

	cols := make([]string, 0, len(tbl.schema.Columns)+len(tbl.schema.Derived))
	for _, col := range tbl.schema.Columns {
		if col != "-" {
			cols = append(cols, col)
		}
	}
	for _, col := range tbl.schema.Derived {
		cols = append(cols, str.Uncapitalize(col))
	}
	tbl.header = NewHeader([][]string{tbl.schema.Columns}, cols)

	idxs := make([][]string, 0, len(tbl.schema.Indexes))
	keys := make([][]string, 0, 1)
	tbl.singleton = false
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
	return tbl.info.Size / tbl.info.Nrows
}

func (tbl *Table) Transform() Query {
	return tbl
}

func (tbl *Table) Fixed() Fixed {
	return nil
}

func (tbl *Table) Updateable() string {
	if db19.IsSchemaTable(tbl.name) {
		return "" // schema tables are virtual, so they can't be updated
	}
	return tbl.name
}

// modeNotSupported reports whether the table is a virtual schema table.
// These are materialized per transaction, so a cursor cannot use them.
func (tbl *Table) modeNotSupported(mode Mode) bool {
	return mode == CursorMode && db19.IsSchemaTable(tbl.name)
}

func (tbl *Table) optimize(mode Mode, req Require) (Cost, Cost, any) {
	if tbl.modeNotSupported(mode) {
		return impossible, impossible, nil
	}
	if tbl.singleton || req.use == ReqNone {
		index := tbl.indexes[0]
		return 0, tbl.costFor(index, req),
			tableApproach{index: index, mode: mode}
	}
	// trackIndexUsage is deliberately called in optimize
	// (i.e. for each approach strategy considered)
	// to make it insensitive to the current physical index
	// and its resulting approach strategy.
	// Only tracking the chosen approach strategy e.g. in SetApproach
	// tends to lock in the current physical index.
	tbl.trackIndexUsage(req)
	best := newBest[[]string]()
	for i, idx := range tbl.indexes {
		if req.use == ReqUnique && !tbl.uniqueForLookup(i) {
			continue
		}
		if req.SatisfiedBy(idx) {
			v := tbl.costFor(idx, req)
			best.update(0, v, idx)
		}
	}
	if best.none() {
		return impossible, impossible, nil
	}
	return 0, best.varcost, tableApproach{index: best.data, mode: mode}
}

const smallTable = 16 * 1024

func (tbl *Table) trackIndexUsage(req Require) {
	if req.use == ReqUnique ||
		tbl.info.Size < smallTable ||
		len(tbl.indexes) == 1 ||
		len(req.cols) == 0 {
		return
	}
	nrows := int(req.frac * float32(tbl.info.Nrows))
	if nrows <= 2 {
		return
	}
	for i, idx := range tbl.indexes {
		if req.SatisfiedBy(idx) {
			tbl.tran.BusyAddScan(tbl.name, i, nrows)
		}
	}
}

// uniqueForLookup returns whether a Lookup on this index yields at most one row.
// Unique 'u' indexes allow multiple all-empty entries
// so they are not usable for Lookup.
func (tbl *Table) uniqueForLookup(i int) bool {
	return tbl.schema.Indexes[i].Mode != 'u'
}

func (tbl *Table) costFor(index []string, req Require) Cost {
	result := 0
	if req.use == ReqUnique {
		assert.That(req.frac == 0)
	} else {
		isPhysicalOrder := slices.Equal(index, tbl.indexes[0])
		cost := scanCost(tbl.info.Nrows, tbl.info.Size, isPhysicalOrder)
		// prefer indexes with fewer columns
		cost += tbl.info.Nrows * len(index) * tableColCost
		result = Cost(float64(req.frac) * float64(cost))
	}
	if req.nseeks > 0 {
		if i := slc.IndexFn(tbl.indexes, index, slices.Equal); i != -1 {
			result += Cost(req.nseeks) * tbl.lookupCostI(i)
		}
	}
	return result
}

// scanCost is the cost to iterate an entire table.
// It is based on BenchmarkTableGetCost
func scanCost(tableRows, tableBytes int, isPhysicalOrder bool) Cost {
	if tableRows <= 0 || tableBytes <= 0 {
		return 0
	}

	recordSize := float64(tableBytes) / float64(tableRows)

	// 1. Sequential base cost per row (scales with record size).
	rowCost := scanBaseCost + (scanRowCost * recordSize)

	// 2. Cold random access penalties
	if !isPhysicalOrder {
		// Cold baseline: un-prefetched DRAM line pulls on small tables
		// (~2.5x physical)
		coldFloor := scanColdBaseCost + (scanColdRowCost * recordSize)

		// Cache spill penalty:
		// TLB page-walks and DRAM row conflicts as table outgrows 16 MiB
		mm := float64(tableBytes * tableBytes)
		penaltyFactor := mm / (mm + (cacheMidpointBytes * cacheMidpointBytes))
		dramSpill := penaltyFactor *
			(scanSpillBaseCost + (scanSpillRowCost * recordSize))
		rowCost += coldFloor + dramSpill
	}

	return tableRows * int(rowCost)
}

func (tbl *Table) setApproach(req Require, approach any, _ QueryTran) {
	tbl.req = req
	ap := approach.(tableApproach)
	tbl.setIndex(ap.index, ap.mode)
}

func (tbl *Table) setIndex(index []string, mode Mode) {
	tbl.cursorMode = (mode == CursorMode)
	if tbl.singleton {
		index = tbl.allKeys[0]
	}
	tbl.index = index
	tbl.iIndex = tbl.indexi(index)
	tbl.indexEncode = tbl.IndexEncodes(index)
	IdxUse(tbl.name, tbl.index)
}

// SetIndex sets the index used to access the table.
// It also sets req to ReqAny so that Select and Lookup can be used.
func (tbl *Table) SetIndex(index []string, mode Mode) {
	tbl.setIndex(index, mode)
	tbl.req = Require{use: ReqAny}
}

// IndexEncodes returns whether the index key is encoded
// (multi-field or non-key, i.e. has BestKey appended or Fields2)
func (tbl *Table) IndexEncodes(index []string) bool {
	return len(index) > 1 ||
		!slc.ContainsFn(tbl.allKeys, index, set.Equal[string])
}

func (tbl *Table) indexi(index []string) int {
	i := slc.IndexFn(tbl.indexes, index, slices.Equal)
	assert.That(i >= 0)
	return i
}

func (tbl *Table) lookupCost(idx []string) Cost {
	return tbl.lookupCostI(tbl.indexi(idx))

}

func (tbl *Table) lookupCostI(i int) Cost {
	var levels int
	if db19.IsSchemaTable(tbl.name) {
		// schema tables are materialized in memory, not btrees
		levels = 1
	} else if tbl.info.Indexes == nil { // tests
		levels = 1
		if tbl.info.Nrows > 100 {
			levels = 2
		}
	} else {
		levels = tbl.info.Indexes[i].BtreeLevels()
	}
	return tableLookupCost(levels)
}

// execution --------------------------------------------------------

func (tbl *Table) Lookup(_ *Thread, sels Sels) Row {
	dbg.Assert(func() bool { return checkSels(sels, tbl.header.Columns) })
	tbl.nlooks++
	key := ""
	if !tbl.singleton {
		assert.That(tbl.req.use == ReqUnique || tbl.req.use == ReqAny)
		ix := &tbl.schema.Indexes[tbl.iIndex]
		key = selOrg(tbl.indexEncode, ix.Fields, sels, true)
		// unique indexes ('u') allow multiple empty entries (via Fields2)
		// so a Lookup is only valid when the key is non-empty
		assert.That(ix.Mode != 'u' || key != "")
	}
	return tbl.LookupRaw(key)
}

func checkSels(sels Sels, srcCols []string) bool {
	for _, sel := range sels {
		if !slices.Contains(srcCols, sel.col) {
			return false
		}
	}
	return true
}

func (tbl *Table) LookupRaw(key string) Row {
	// unique indexes ('u') allow multiple empty entries (via Fields2)
	// so an empty key can't be used to find them by value alone;
	// a non-empty key is fine since Fields2 is only appended to empty entries
	assert.That(tbl.schema.Indexes[tbl.iIndex].Mode != 'u' || key != "")
	rec := tbl.tran.Lookup(tbl.name, tbl.iIndex, key)
	if rec == "" {
		return Row{}
	}
	return NewRow(rec)
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
	return tbl.GetFilter(dir, nil, nil)
}

func (tbl *Table) GetFilter(dir Dir, offFilter func(uint64) bool,
	keyFilter func(key string, off uint64) bool) Row {
	defer func(t uint64) { tbl.tget += tsc.Read() - t }(tsc.Read())
	tbl.ensureIter()
	for {
		if dir == Prev {
			tbl.iter.Prev(tbl.tran)
		} else {
			tbl.iter.Next(tbl.tran)
		}
		if tbl.iter.Eof() {
			return Row{}
		}
		off := tbl.iter.CurOff()
		if offFilter != nil && !offFilter(off) {
			continue
		}
		if keyFilter != nil {
			key, _ := tbl.iter.Cur() // allocates key string
			if !keyFilter(key, off) {
				continue
			}
		}
		rec := tbl.iter.GetRecord(tbl.tran)
		row := NewRow(rec)
		if tbl.singleton &&
			!singletonFilter(tbl.header, row, tbl.sels) {
			return Row{}
		}
		tbl.ngets++
		return row
	}
}

func (tbl *Table) Select(sels Sels) {
	// singleton doesn't use an index range - it filters (via singletonFilter
	// in GetFilter) after a full scan, so it supports Select regardless of
	// tbl.req.use.
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
	assert.That(tbl.req.use == ReqAny ||
		tbl.req.use == ReqGroup || tbl.req.use == ReqOrder)
	dbg.Assert(func() bool { return checkSels(sels, tbl.header.Columns) })
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
	tbl.Select(nil)
	tbl.Rewind()
	rows := make([]Row, 0, 8)
	for {
		tbl.iter.Next(tbl.tran)
		if tbl.iter.Eof() {
			break
		}
		rec := tbl.iter.GetRecord(tbl.tran)
		row := NewRow(rec)
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
