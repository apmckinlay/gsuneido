// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"sort"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/ord"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Summarize struct {
	Query1
	by       []string
	cols     []string
	ops      []string
	ons      []string
	wholeRow bool
	summarizeApproach
	rewound bool
	srcHdr  *Header
	get     func(th *Thread, su *Summarize, dir Dir) Row
	t       QueryTran
}

type summarizeApproach struct {
	strategy sumStrategy
	index    []string
}

type sumStrategy int

const (
	// sumSeq reads in order of 'by', and groups consecutive, like projSeq
	sumSeq sumStrategy = iota + 1
	// sumMap uses a map to accumulate results. It is not incremental -
	// it must process all the source before producing any results.
	sumMap
	// sumIdx uses an index to get an overall min or max
	sumIdx
	// sumTbl optimizes <table> summarize count
	sumTbl
)

func NewSummarize(src Query, by, cols, ops, ons []string) *Summarize {
	if !set.Subset(src.Columns(), by) {
		panic("summarize: nonexistent columns: " +
			str.Join(", ", set.Difference(by, src.Columns())))
	}
	check(by)
	check(ons)
	for i := 0; i < len(cols); i++ {
		if cols[i] == "" {
			if ons[i] == "" {
				cols[i] = "count"
			} else {
				cols[i] = ops[i] + "_" + ons[i]
			}
		}
	}
	su := &Summarize{Query1: Query1{source: src},
		by: by, cols: cols, ops: ops, ons: ons}
	// if single min or max, and on is a key, then we can give the whole row
	su.wholeRow = su.minmax1() && slc.ContainsFn(src.Keys(), ons, set.Equal[string])
	return su
}

func check(cols []string) {
	for _, c := range cols {
		if strings.HasSuffix(c, "_lower!") {
			panic("can't summarize _lower! fields")
		}
	}
}

func (su *Summarize) minmax1() bool {
	return len(su.by) == 0 &&
		len(su.ops) == 1 && (su.ops[0] == "min" || su.ops[0] == "max")
}

func (su *Summarize) SetTran(t QueryTran) {
	su.t = t
	su.source.SetTran(t)
}

func (su *Summarize) String() string {
	s := parenQ2(su.source) + " SUMMARIZE"
	switch su.strategy {
	case sumSeq:
		s += "-SEQ"
	case sumMap:
		s += "-MAP"
	case sumIdx:
		s += "-IDX"
	case sumTbl:
		s += "-TBL"
	}
	if su.wholeRow {
		s += "*"
	}
	if len(su.by) > 0 {
		s += " " + str.Join(", ", su.by) + ","
	}
	sep := " "
	for i := range su.cols {
		s += sep
		sep = ", "
		if su.cols[i] != "" {
			s += su.cols[i] + " = "
		}
		s += su.ops[i]
		if su.ops[i] != "count" {
			s += " " + su.ons[i]
		}
	}
	return s
}

func (su *Summarize) Columns() []string {
	if su.wholeRow {
		return set.Union(su.source.Columns(), su.cols)
	}
	return set.Union(su.by, su.cols)
}

func (su *Summarize) Keys() [][]string {
	return projectKeys(su.source.Keys(), su.by)
}

func (su *Summarize) Indexes() [][]string {
	if len(su.by) == 0 || containsKey(su.by, su.source.Keys()) {
		return su.source.Indexes()
	}
	var idxs [][]string
	for _, src := range su.source.Indexes() {
		if set.StartsWithSet(src, su.by) {
			idxs = append(idxs, src)
		}
	}
	return idxs
}

// containsKey returns true if a set of columns contain one of the keys
func containsKey(cols []string, keys [][]string) bool {
	for _, key := range keys {
		if set.Subset(cols, key) {
			return true
		}
	}
	return false
}

func (su *Summarize) Nrows() int {
	nr := su.source.Nrows()
	if len(su.by) == 0 {
		nr = 1
	} else if !containsKey(su.by, su.source.Keys()) {
		nr /= 2 // ???
	}
	return nr
}

func (su *Summarize) rowSize() int {
	return su.source.rowSize() + len(su.cols)*8
}

func (su *Summarize) Updateable() string {
	return "" // override Query1 source.Updateable
}

func (su *Summarize) SingleTable() bool {
	return false
}

func (*Summarize) Output(*Thread, Record) {
	panic("can't output to this query")
}

func (su *Summarize) Transform() Query {
	su.source = su.source.Transform()
	return su
}

func (su *Summarize) optimize(mode Mode, index []string) (Cost, interface{}) {
	if _, ok := su.source.(*Table); ok &&
		len(su.by) == 0 && len(su.ops) == 1 && su.ops[0] == "count" {
		Optimize(su.source, mode, nil)
		return 1, &summarizeApproach{strategy: sumTbl}
	}
	seqCost, seqApp := su.seqCost(mode, index)
	idxCost, idxApp := su.idxCost(mode)
	mapCost, mapApp := su.mapCost(mode, index)
	return min3(seqCost, seqApp, idxCost, idxApp, mapCost, mapApp)
}

func (su *Summarize) seqCost(mode Mode, index []string) (Cost, interface{}) {
	approach := &summarizeApproach{strategy: sumSeq}
	if len(su.by) == 0 || containsKey(su.by, su.source.Keys()) {
		if len(su.by) != 0 {
			approach.index = index
		}
		cost := Optimize(su.source, mode, approach.index)
		return cost, approach
	}
	best := bestGrouped(su.source, mode, index, su.by)
	if best.index == nil {
		return impossible, nil
	}
	approach.index = best.index
	return best.cost, approach
}

func (su *Summarize) idxCost(Mode) (Cost, interface{}) {
	if len(su.by) > 0 || !su.minmax1() {
		return impossible, nil
	}
	// dividing by nrows since we're only reading one row
	// NOTE: this is not correct if there is any fixed cost
	nr := ord.Max(1, su.source.Nrows())
	// cursorMode to bypass temp index
	cost := Optimize(su.source, CursorMode, su.ons) / nr
	approach := &summarizeApproach{strategy: sumIdx, index: su.ons}
	return cost, approach
}

func (su *Summarize) mapCost(mode Mode, index []string) (Cost, interface{}) {
	if index != nil {
		return impossible, nil
	}
	cost := Optimize(su.source, mode, nil)
	cost += cost / 2 // add 50% for map overhead
	approach := &summarizeApproach{strategy: sumMap}
	return cost, approach
}

func (su *Summarize) setApproach(_ []string, approach interface{}, tran QueryTran) {
	su.summarizeApproach = *approach.(*summarizeApproach)
	su.source = SetApproach(su.source, su.index, tran)
	switch su.strategy {
	case sumTbl:
		su.get = getTbl
	case sumIdx:
		su.get = getIdx
	case sumMap:
		t := sumMapT{}
		su.get = t.getMap
	case sumSeq:
		t := sumSeqT{}
		su.get = t.getSeq
	}
	su.rewound = true
	su.srcHdr = su.source.Header()
}

// execution --------------------------------------------------------

const sumMaxSize = 10_000

func (su *Summarize) Header() *Header {
	if su.wholeRow {
		flds := su.source.Header().Fields
		n := len(flds)
		flds = append(flds[:n:n], su.cols)
		return NewHeader(flds, su.Columns())
	}
	flds := append(su.by, su.cols...)
	return NewHeader([][]string{flds}, flds)
}

func (su *Summarize) Rewind() {
	su.source.Rewind()
	su.rewound = true
}

func (su *Summarize) Get(th *Thread, dir Dir) Row {
	defer func() { su.rewound = false }()
	return su.get(th, su, dir)
}

func getTbl(_ *Thread, su *Summarize, _ Dir) Row {
	if !su.rewound {
		return nil
	}
	tbl := su.source.(*Table)
	var rb RecordBuilder
	rb.Add(IntVal(tbl.Nrows()).(Packable))
	return Row{DbRec{Record: rb.Build()}}
}

func getIdx(th *Thread, su *Summarize, dir Dir) Row {
	if !su.rewound {
		return nil
	}
	if str.EqualCI(su.ops[0], "min") {
		dir = Next
	} else { // max
		dir = Prev
	}
	row := su.source.Get(th, dir)
	if row == nil {
		return nil
	}
	var rb RecordBuilder
	rb.AddRaw(row.GetRaw(su.srcHdr, su.ons[0]))
	rec := rb.Build()
	if su.wholeRow {
		return append(row, DbRec{Record: rec})
	}
	return Row{DbRec{Record: rec}}
}

//-------------------------------------------------------------------

type sumMapT struct {
	mapList []mapPair
	mapPos  int
}

type mapPair struct {
	key Record
	ops []sumOp
}

func (t *sumMapT) getMap(th *Thread, su *Summarize, dir Dir) Row {
	if su.rewound {
		assert.That(!su.wholeRow)
		t.mapList = su.buildMap(th)
		if dir == Next {
			t.mapPos = -1
		} else { // Prev
			t.mapPos = len(t.mapList)
		}
	}
	if dir == Next {
		t.mapPos++
	} else { // Prev
		t.mapPos--
	}
	if t.mapPos < 0 || len(t.mapList) <= t.mapPos {
		return nil
	}
	key := t.mapList[t.mapPos].key
	var rb RecordBuilder
	for i := 0; i < len(su.by); i++ {
		rb.AddRaw(key.GetRaw(i))
	}
	ops := t.mapList[t.mapPos].ops
	for i := range ops {
		val, _ := ops[i].result()
		rb.Add(val.(Packable))
	}
	return Row{DbRec{Record: rb.Build()}}
}

func (su *Summarize) buildMap(th *Thread) []mapPair {
	hdr := su.source.Header()
	sumMap := make(map[Record][]sumOp)
	var thread Thread
	for {
		row := su.source.Get(th, Next)
		if row == nil {
			break
		}
		key := keyRec(row, hdr, su.by)
		sums, ok := sumMap[key]
		if !ok {
			sums = su.newSums()
			sumMap[key] = sums
			if len(sumMap) > sumMaxSize {
				panic("summarize too large")
			}
		}
		for i := range sums {
			x := row.GetVal(hdr, su.ons[i], &thread, MakeSuTran(su.t))
			sums[i].add(x, nil)
		}
	}
	i := 0
	list := make([]mapPair, len(sumMap))
	for key, ops := range sumMap {
		list[i] = mapPair{key: key, ops: ops}
		i++
	}
	if len(list) <= 3 { // for tests
		sort.Slice(list,
			func(i, j int) bool { return list[i].key < list[j].key })
	}
	return list
}

func keyRec(row Row, hdr *Header, cols []string) Record {
	var rb RecordBuilder
	for _, fld := range cols {
		rb.AddRaw(row.GetRaw(hdr, fld))
	}
	return rb.Build()
}

//-------------------------------------------------------------------

type sumSeqT struct {
	curDir  Dir
	curRow  Row
	nextRow Row
	sums    []sumOp
}

func (t *sumSeqT) getSeq(th *Thread, su *Summarize, dir Dir) Row {
	if su.rewound {
		t.sums = su.newSums()
		t.curDir = dir
		t.curRow = nil
		t.nextRow = su.source.Get(th, dir)
	}

	// if direction changes, have to skip over previous result
	if dir != t.curDir {
		if t.nextRow == nil {
			su.source.Rewind()
		}
		for {
			t.nextRow = su.source.Get(th, dir)
			if t.nextRow == nil || !su.sameBy(t.curRow, t.nextRow) {
				break
			}
		}
		t.curDir = dir
	}

	if t.nextRow == nil {
		return nil
	}
	t.curRow = t.nextRow
	for i := range t.sums {
		t.sums[i].reset()
	}
	var thread Thread
	for {
		for i := range t.sums {
			x := t.nextRow.GetVal(su.srcHdr, su.ons[i], &thread, MakeSuTran(su.t))
			t.sums[i].add(x, su.sumRow(t.nextRow))
		}
		t.nextRow = su.source.Get(th, dir)
		if t.nextRow == nil || !su.sameBy(t.curRow, t.nextRow) {
			break
		}
	}
	// output after each group
	return su.seqRow(t.curRow, t.sums)
}

func (su *Summarize) sumRow(row Row) Row {
	if su.wholeRow {
		return row
	}
	return nil
}

func (su *Summarize) sameBy(row1, row2 Row) bool {
	for _, f := range su.by {
		if row1.GetRaw(su.srcHdr, f) != row2.GetRaw(su.srcHdr, f) {
			return false
		}
	}
	return true
}

func (su *Summarize) seqRow(curRow Row, sums []sumOp) Row {
	var rb RecordBuilder
	if !su.wholeRow {
		for _, fld := range su.by {
			rb.AddRaw(curRow.GetRaw(su.srcHdr, fld))
		}
	}
	for _, sum := range sums {
		val, _ := sum.result()
		rb.Add(val.(Packable))
	}
	row := Row{DbRec{Record: rb.Build()}}
	if su.wholeRow {
		_, wholeRow := sums[0].result()
		row = append(wholeRow, row[0])
	}
	return row
}

func (su *Summarize) Select(cols, vals []string) {
	su.source.Select(cols, vals)
	su.rewound = true
}

// operations -------------------------------------------------------

func (su *Summarize) newSums() []sumOp {
	sums := make([]sumOp, len(su.ops))
	for i, op := range su.ops {
		sums[i] = newSumOp(op)
	}
	return sums
}

func newSumOp(op string) sumOp {
	switch op {
	case "count":
		return &sumCount{}
	case "total":
		return &sumTotal{total: Zero}
	case "average":
		return &sumAverage{total: Zero}
	case "min":
		return &sumMin{}
	case "max":
		return &sumMax{}
	case "list":
		return &sumList{set: &SuObject{}}
	}
	panic("shouldn't reach here")
}

type sumOp interface {
	add(val Value, row Row)
	result() (Value, Row)
	reset()
}

type sumCount struct {
	count int
}

func (sum *sumCount) add(_ Value, _ Row) {
	sum.count++
}
func (sum *sumCount) result() (Value, Row) {
	return IntVal(sum.count), nil
}
func (sum *sumCount) reset() {
	sum.count = 0
}

type sumTotal struct {
	total Value
}

func (sum *sumTotal) add(val Value, _ Row) {
	sum.total = OpAdd(sum.total, val)
}
func (sum *sumTotal) result() (Value, Row) {
	return sum.total, nil
}
func (sum *sumTotal) reset() {
	sum.total = Zero
}

type sumAverage struct {
	count int
	total Value
}

func (sum *sumAverage) add(val Value, _ Row) {
	sum.count++
	sum.total = OpAdd(sum.total, val)
}
func (sum *sumAverage) result() (Value, Row) {
	return OpDiv(sum.total, IntVal(sum.count)), nil
}
func (sum *sumAverage) reset() {
	sum.count = 0
	sum.total = Zero
}

type sumMin struct {
	val Value
	row Row
}

func (sum *sumMin) add(val Value, row Row) {
	if sum.val == nil || val.Compare(sum.val) < 0 {
		sum.val, sum.row = val, row
	}
}
func (sum *sumMin) result() (Value, Row) {
	return sum.val, sum.row
}
func (sum *sumMin) reset() {
	sum.val = nil
	sum.row = nil
}

type sumMax struct {
	val Value
	row Row
}

func (sum *sumMax) add(val Value, row Row) {
	if sum.val == nil || val.Compare(sum.val) > 0 {
		sum.val, sum.row = val, row
	}
}
func (sum *sumMax) result() (Value, Row) {
	return sum.val, sum.row
}
func (sum *sumMax) reset() {
	sum.val = nil
	sum.row = nil
}

type sumList struct {
	set *SuObject
}

func (sum *sumList) add(val Value, _ Row) {
	sum.set.Set(val, True)
	if sum.set.Size() > sumMaxSize {
		panic("summarize list too large")
	}
}
func (sum *sumList) result() (Value, Row) {
	list := make([]Value, sum.set.Size())
	iter := sum.set.Iter2(true, true)
	for i := range list {
		x, _ := iter()
		list[i] = x
	}
	if len(list) <= 3 { // for tests
		sort.Slice(list,
			func(i, j int) bool { return list[i].Compare(list[j]) < 0 })
	}
	return NewSuObject(list), nil
}
func (sum *sumList) reset() {
	sum.set = &SuObject{}
}
