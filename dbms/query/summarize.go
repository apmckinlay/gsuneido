// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"log"
	"sort"
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/shmap"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/tsc"
)

type Summarize struct {
	Query1
	t   QueryTran
	get func(th *Thread, su *Summarize, dir Dir) Row
	st  *SuTran
	by  []string
	// cols, ops, and ons are parallel
	cols []string
	ops  []string
	ons  []string
	summarizeApproach
	wholeRow bool
	rewound  bool
	unique   bool
	hint     sumHint
	th       *Thread
}

type summarizeApproach struct {
	index []string
	strat sumStrategy
	frac  float64
}

type sumStrategy int

type sumHint string

const (
	sumSmall sumHint = "small"
	sumLarge sumHint = "large"
)

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

func NewSummarize(src Query, hint sumHint, by, cols, ops, ons []string) *Summarize {
	if !set.Subset(src.Columns(), by) {
		panic("summarize: nonexistent columns: " +
			str.Join(", ", set.Difference(by, src.Columns())))
	}
	check(by)
	check(ons)
	for i := range len(cols) {
		if cols[i] == "" {
			cols[i] = defaultColName(ops[i], ons[i])
		}
	}
	su := &Summarize{hint: hint, by: by, cols: cols, ops: ops, ons: ons}
	su.source = src
	sort.Stable(su)
	su.unique = hasKey(cols, src.Keys(), src.Fixed())
	// if single min or max, and on is a key, then we can give the whole row
	su.wholeRow = su.minmax1() && slc.ContainsFn(src.Keys(), ons, set.Equal[string])
	su.header = su.getHeader()
	su.keys = projectKeys(src.Keys(), su.by)
	su.indexes = projectIndexes(src.Indexes(), su.by)
	su.fixed = projectFixed(src.Fixed(), by)
	su.setNrows(su.getNrows())
	su.rowSiz.Set(su.source.rowSize() + len(su.cols)*8) // ???
	su.fast1.Set(src.fastSingle())
	su.lookCost.Set(su.getLookupCost())
	return su
}

func defaultColName(op, on string) string {
	if op == "count" {
		return "count"
	}
	return op + "_" + on
}

// Len / Less / Swap implement sort.Interface
func (su *Summarize) Len() int {
	return len(su.cols)
}
func (su *Summarize) Less(i, j int) bool {
	return su.ons[i] < su.ons[j]
}
func (su *Summarize) Swap(i, j int) {
	su.ons[i], su.ons[j] = su.ons[j], su.ons[i]
	su.cols[i], su.cols[j] = su.cols[j], su.cols[i]
	su.ops[i], su.ops[j] = su.ops[j], su.ops[i]
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
	su.st = MakeSuTran(su.t)
	su.source.SetTran(t)
}

func (su *Summarize) String() string {
	s := "summarize"
	switch su.strat {
	case 0:
		s += str.Opt("/*", string(su.hint), "*/")
	case sumSeq:
		s += "-seq"
	case sumMap:
		s += "-map"
	case sumIdx:
		s += "-idx"
	case sumTbl:
		s += "-tbl"
	default:
		assert.ShouldNotReachHere()
	}
	if su.strat != 0 && su.wholeRow {
		s += "*"
	}
	return s + su.string2()
}

func (su *Summarize) string2() string {
	s := ""
	if len(su.by) > 0 {
		s += " " + str.Join(", ", su.by) + ","
	}
	sep := " "
	for i := range su.cols {
		s += sep
		sep = ", "
		if su.cols[i] != defaultColName(su.ops[i], su.ons[i]) {
			s += su.cols[i] + " = "
		}
		s += su.ops[i]
		if su.ops[i] != "count" {
			s += " " + su.ons[i]
		}
	}
	return s
}

func (su *Summarize) getNrows() (int, int) {
	nr, pop := su.source.Nrows()
	if len(su.by) == 0 {
		nr = 1
	} else if !su.unique {
		nr /= 10 // ??? (matches lookupCost)
	}
	return nr, pop
}

func (su *Summarize) Updateable() string {
	return ""
}

func (su *Summarize) SingleTable() bool {
	return false
}

func (*Summarize) Output(*Thread, Record) {
	panic("can't output to this query")
}

func (su *Summarize) Transform() Query {
	src := su.source.Transform()
	if _, ok := src.(*Nothing); ok {
		return NewNothing(su)
	}
	if p, ok := src.(*Project); ok && p.unique {
		// remove project-copy
		return NewSummarize(p.source, su.hint, su.by, su.cols, su.ops, su.ons)
	}
	if src != su.source {
		return NewSummarize(src, su.hint, su.by, su.cols, su.ops, su.ons)
	}
	return su
}

func (su *Summarize) optimize(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	if _, ok := su.source.(*Table); ok &&
		len(su.by) == 0 && len(su.ops) == 1 && su.ops[0] == "count" {
		Optimize(su.source, mode, nil, 0)
		return 0, 1, &summarizeApproach{strat: sumTbl}
	}
	seqFixCost, seqVarCost, seqApp := su.seqCost(mode, index, frac)
	idxFixCost, idxVarCost, idxApp := su.idxCost(mode)
	mapFixCost, mapVarCost, mapApp := su.mapCost(mode, index, frac)
	// trace.Println("summarize seq", seqFixCost, "+", seqVarCost, "=", seqFixCost+seqVarCost)
	// trace.Println("summarize idx", idxFixCost, "+", idxVarCost, "=", idxFixCost+idxVarCost)
	// trace.Println("summarize map", mapFixCost, "+", mapVarCost, "=", mapFixCost+mapVarCost)
	return min3(seqFixCost, seqVarCost, seqApp, idxFixCost, idxVarCost, idxApp,
		mapFixCost, mapVarCost, mapApp)
}

func (su *Summarize) seqCost(mode Mode, index []string, frac float64) (Cost, Cost, any) {
	if len(su.by) == 0 {
		frac = min(1, frac)
	}
	approach := &summarizeApproach{strat: sumSeq, frac: frac}
	if len(su.by) == 0 || hasKey(su.by, su.source.Keys(), su.source.Fixed()) {
		if len(su.by) != 0 {
			approach.index = index
		}
		fixcost, varcost := Optimize(su.source, mode, approach.index, frac)
		return fixcost, varcost, approach
	}
	best := bestGrouped(su.source, mode, index, frac, su.by)
	if best.index == nil {
		return impossible, impossible, nil
	}
	approach.index = best.index
	return best.fixcost, best.varcost, approach
}

func (su *Summarize) idxCost(mode Mode) (Cost, Cost, any) {
	if !su.minmax1() {
		return impossible, impossible, nil
	}
	nrows, _ := su.source.Nrows()
	frac := 1.
	if nrows > 0 {
		frac = 1 / float64(nrows)
	}
	fixcost, varcost := Optimize(su.source, mode, su.ons, frac)
	return fixcost, varcost,
		&summarizeApproach{strat: sumIdx, index: su.ons, frac: frac}
}

func (su *Summarize) mapCost(mode Mode, index []string, _ float64) (Cost, Cost, any) {
	// WARNING technically, map should only be allowed in ReadMode
	nrows, _ := su.Nrows()
	if index != nil || su.hint == sumLarge ||
		(nrows > mapThreshold && su.hint != sumSmall) {
		return impossible, impossible, nil
	}
	fixcost, varcost := Optimize(su.source, mode, nil, 1)
	fixcost += nrows * 20 // ???
	return fixcost + varcost, 0, &summarizeApproach{strat: sumMap, frac: 1}
}

func (su *Summarize) setApproach(_ []string, frac float64, approach any, tran QueryTran) {
	su.summarizeApproach = *approach.(*summarizeApproach)
	switch su.strat {
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
	default:
		assert.ShouldNotReachHere()
	}
	su.source = SetApproach(su.source, su.index, su.frac, tran)
	su.rewound = true
	su.header = su.getHeader()
}

// execution --------------------------------------------------------

func (su *Summarize) getHeader() *Header {
	if su.wholeRow {
		flds := su.source.Header().Fields
		flds = slc.With(flds, su.cols)
		return NewHeader(flds, su.getColumns())
	}
	flds := append(su.by, su.cols...)
	return NewHeader([][]string{flds}, flds)
}

func (su *Summarize) getColumns() []string {
	if su.wholeRow {
		return set.Union(su.source.Columns(), su.cols)
	}
	return set.Union(su.by, su.cols)
}

func (su *Summarize) Rewind() {
	su.source.Rewind()
	su.rewound = true
}

func (su *Summarize) Get(th *Thread, dir Dir) Row {
	defer func(t uint64) { su.tget += tsc.Read() - t }(tsc.Read())
	defer func() { su.rewound = false }()
	row := su.get(th, su, dir)
	if row != nil {
		su.ngets++
	}
	return row
}

func getTbl(_ *Thread, su *Summarize, _ Dir) Row {
	if !su.rewound {
		return nil
	}
	tbl := su.source.(*Table)
	var rb RecordBuilder
	nr, _ := tbl.Nrows()
	rb.Add(IntVal(nr))
	return Row{DbRec{Record: rb.Build()}}
}

func getIdx(th *Thread, su *Summarize, _ Dir) Row {
	if !su.rewound {
		return nil
	}
	dir := Prev // max
	if str.EqualCI(su.ops[0], "min") {
		dir = Next
	}
	row := su.source.Get(th, dir)
	if row == nil {
		return nil
	}
	var rb RecordBuilder
	rb.AddRaw(row.GetRawVal(su.source.Header(), su.ons[0], th, su.st))
	rec := rb.Build()
	if su.wholeRow {
		return append(row, DbRec{Record: rec})
	}
	return Row{DbRec{Record: rec}}
}

func (su *Summarize) Lookup(th *Thread, cols, vals []string) Row {
	su.nlooks++
	su.Select(cols, vals)
	defer su.Select(nil, nil) // clear
	return su.Get(th, Next)
}

func (su *Summarize) getLookupCost() Cost {
	srcCost := su.source.lookupCost()
	if su.unique {
		return srcCost
	}
	return 10 * srcCost // ??? (matches Nrows)
	//TODO should be 1 lookup + 10 gets
}

//-------------------------------------------------------------------

type sumMapT struct {
	mapList []mapPair
	mapPos  int
}

type mapPair struct {
	row Row
	ops []sumOp
}

func (t *sumMapT) getMap(th *Thread, su *Summarize, dir Dir) Row {
	su.th = th
	defer func() { su.th = nil }()
	if su.rewound {
		assert.That(!su.wholeRow)
		t.mapList = su.buildMap()
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
	row := t.mapList[t.mapPos].row
	var rb RecordBuilder
	for _, col := range su.by {
		rb.AddRaw(row.GetRawVal(su.source.Header(), col, th, su.st))
	}
	ops := t.mapList[t.mapPos].ops
	for i := range ops {
		val, _ := ops[i].result()
		rb.Add(val.(Packable))
	}
	return Row{DbRec{Record: rb.Build()}}
}

func (su *Summarize) buildMap() []mapPair {
	hdr := su.source.Header()
	hfn := func(k rowHash) uint64 { return k.hash }
	eqfn := func(x, y rowHash) bool {
		return x.hash == y.hash &&
			equalCols(x.row, y.row, hdr, su.by, su.th, su.st)
	}
	sumMap := shmap.NewMapFuncs[rowHash, []sumOp](hfn, eqfn)
	warned := false
	for {
		row := su.source.Get(su.th, Next)
		if row == nil {
			break
		}
		rh := rowHash{hash: hashCols(row, hdr, su.by, su.th, su.st), row: row}
		sums, ok := sumMap.Get(rh)
		if !ok {
			sums = su.newSums()
			sumMap.Put(rh, sums)
			if !warned && sumMap.Size() > mapWarn {
				// log inside loop in case we run out of memory
				warned = true
				Warning("summarize-map large >", mapWarn)
			}
		}
		su.addToSums(sums, row, su.th, su.st)
	}
	if sumMap.Size() > 2*mapWarn {
		log.Println("summarize-map large =", sumMap.Size())
	}
	i := 0
	list := make([]mapPair, sumMap.Size())
	iter := sumMap.Iter()
	for rh, ops, ok := iter(); ok; rh, ops, ok = iter() {
		list[i] = mapPair{row: rh.row, ops: ops}
		i++
	}
	return list
}

//-------------------------------------------------------------------

type sumSeqT struct {
	curRow  Row
	nextRow Row
	sums    []sumOp
	curDir  Dir
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
			if t.nextRow == nil || !su.sameBy(th, su.st, t.curRow, t.nextRow) {
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
	for {
		su.addToSums(t.sums, t.nextRow, th, su.st)
		t.nextRow = su.source.Get(th, dir)
		if t.nextRow == nil || !su.sameBy(th, su.st, t.curRow, t.nextRow) {
			break
		}
	}
	// output after each group
	return su.seqRow(th, t.curRow, t.sums)
}

func (su *Summarize) addToSums(sums []sumOp, row Row, th *Thread, st *SuTran) {
	for i := 0; i < len(su.ons); {
		raw := "*uninit*"
		var val Value
		for col := su.ons[i]; i < len(su.ons) && su.ons[i] == col; i++ {
			switch su.ops[i] {
			case "count":
				sums[i].add("", nil, row)
			case "list", "min", "max":
				if raw == "*uninit*" {
					raw = row.GetRawVal(su.source.Header(), col, th, st)
				}
				sums[i].add(raw, nil, row)
			default: // total, average
				if val == nil {
					val = row.GetVal(su.source.Header(), col, th, st)
				}
				sums[i].add("", val, row)
			}
		}
	}
}

func (su *Summarize) sameBy(th *Thread, st *SuTran, row1, row2 Row) bool {
	for _, f := range su.by {
		if row1.GetRawVal(su.source.Header(), f, th, st) !=
			row2.GetRawVal(su.source.Header(), f, th, st) {
			return false
		}
	}
	return true
}

func (su *Summarize) seqRow(th *Thread, curRow Row, sums []sumOp) Row {
	var rb RecordBuilder
	if !su.wholeRow {
		for _, fld := range su.by {
			rb.AddRaw(curRow.GetRawVal(su.source.Header(), fld, th, su.st))
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
	su.nsels++
	su.source.Select(cols, vals)
	su.rewound = true
}

func (*Summarize) Simple(*Thread) []Row {
	panic("Simple not implemented for summarize")
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
		return sumList(make(map[string]struct{}))
	}
	panic(assert.ShouldNotReachHere())
}

type sumOp interface {
	add(raw string, val Value, row Row)
	result() (Value, Row)
	reset()
}

type sumCount struct {
	count int
}

func (sum *sumCount) add(_ string, _ Value, _ Row) {
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

func (sum *sumTotal) add(_ string, val Value, _ Row) {
	defer func() { recover() }() // ignore panics
	sum.total = OpAdd(sum.total, val)
}
func (sum *sumTotal) result() (Value, Row) {
	return sum.total, nil
}
func (sum *sumTotal) reset() {
	sum.total = Zero
}

type sumAverage struct {
	total Value
	count int
}

func (sum *sumAverage) add(_ string, val Value, _ Row) {
	sum.count++
	defer func() { recover() }()
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
	raw string
	row Row
}

func (sum *sumMin) add(raw string, _ Value, row Row) {
	if sum.row == nil || raw < sum.raw {
		sum.raw, sum.row = raw, row
	}
}
func (sum *sumMin) result() (Value, Row) {
	return Unpack(sum.raw), sum.row
}
func (sum *sumMin) reset() {
	sum.raw = ""
	sum.row = nil
}

type sumMax struct {
	raw string
	row Row
}

func (sum *sumMax) add(raw string, _ Value, row Row) {
	if sum.row == nil || raw > sum.raw {
		sum.raw, sum.row = raw, row
	}
}
func (sum *sumMax) result() (Value, Row) {
	return Unpack(sum.raw), sum.row
}
func (sum *sumMax) reset() {
	sum.raw = ""
	sum.row = nil
}

type sumList map[string]struct{}

const sumListLimit = 16384 // ???

func (sum sumList) add(raw string, _ Value, _ Row) {
	sum[raw] = struct{}{}
	if len(sum) > sumListLimit {
		panic(fmt.Sprintf("summarize list too large (> %d)", sumListLimit))
	}
}
func (sum sumList) result() (Value, Row) {
	list := make([]Value, len(sum))
	i := 0
	for raw := range sum {
		list[i] = Unpack(raw)
		i++
	}
	if len(list) <= 3 { // for tests
		sort.Slice(list,
			func(i, j int) bool { return list[i].Compare(list[j]) < 0 })
	}
	return NewSuObject(list), nil
}
func (sum sumList) reset() {
	for k := range sum {
		delete(sum, k)
	}
}
