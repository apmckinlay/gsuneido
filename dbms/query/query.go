// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package query implements query parsing, optimization, and execution.
/*
	Query
		Table
		Nothing
		ProjectNone
		schemaTable
			Tables
			TablesLookup
			Columns
			Indexes
		Query1
			Extend
			Project / Remove
			Rename
			Sort
			Summarize
			TempIndex
			Where
		Query2
			Compatible
				Union
				Compatible1
					Intersect
					Minus
			joinLike
				Join
					LeftJoin
				Times
*/
package query

import (
	"fmt"
	"math"
	"strings"

	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/trace"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/ord"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"golang.org/x/exp/slices"
)

type Query interface {
	// Columns is all the available columns, including derived
	Columns() []string

	// Transform refactors the query for more efficient execution.
	// This stage is not cost based, transforms are applied when possible.
	//
	// Transform methods MUST ensure they call Transform on their children.
	Transform() Query

	// SetTran is used for cursors
	SetTran(tran QueryTran)

	Ordering() []string

	// Fixed returns the field values that are constant from Extend or Where
	Fixed() []Fixed

	// Updateable returns the table name if the rows from the query can be updated
	// else ""
	Updateable() string

	// SingleTable is used by TempIndex.
	// It is true if the query Get returns a single record stored in the database
	SingleTable() bool

	// Indexes returns all the indexes
	Indexes() [][]string

	// Keys returns sets of fields that are unique keys.
	// On a table this will be the key indexes, but on other operations
	// this is logical, there may not necessarily be an index.
	Keys() [][]string

	// Nrows returns n the number of expected result rows from the query,
	// and p the "population" it was drawn from.
	// For example, a Where on a key selects a single row (n = 1),
	// from the entire table with p rows.
	//
	// Nrows should be the same regardless of the strategy.
	// For symmetrical operations e.g. join or union
	// it should give the same result both ways.
	//
	// Nrows does *not* incorporate frac
	Nrows() (n, p int)

	// rowSize returns the average number of bytes per row
	rowSize() int

	Rewind()

	Get(th *runtime.Thread, dir runtime.Dir) runtime.Row

	// Lookup returns the row matching the given key value, or nil if not found.
	// It is used by Compatible (Intersect, Minus, Union). See also: Select
	// It is valid (although not necessarily the most efficient)
	// to implement Lookup with Select and Get
	// in which case it should leave the select cleared.
	// Lookup should rewind.
	Lookup(th *runtime.Thread, cols, vals []string) runtime.Row

	// Select restricts the query to records matching the given packed values.
	// It is used by Join and LeftJoin. See also: Lookup
	// To clear the select, use Select(nil, nil)
	// Select should rewind.
	Select(cols, vals []string)

	Header() *runtime.Header
	Output(th *runtime.Thread, rec runtime.Record)

	String() string

	cacheAdd(index []string, frac float64, fixcost, varcost Cost, approach any)

	// cacheGet returns the cost and approach associated with an index
	// or -1 if the index has not been added.
	cacheGet(index []string, frac float64) (fixcost, varcost Cost, approach any)

	cacheSetCost(frac float64, fixcost, varcost Cost)
	cacheCost() (frac float64, fixcost, varcost Cost)

	// optimize determines the minimum cost strategy based on estimates.
	//
	// index is what is required for order or lookup
	//
	// frac is the estimated fraction of the rows that will be read.
	// It affects the variable cost.
	// frac = 0 means only Lookup, else frac < 1 means Select
	//
	// varcost should already incorporate frac
	optimize(mode Mode, index []string, frac float64) (
		fixcost, varcost Cost, approach any)
	setApproach(index []string, frac float64, approach any, tran QueryTran)

	// lookupCost returns the cost of one Lookup
	lookupCost() Cost

	// fastSingle returns whether it's a fast singleton.
	// This is mostly equivalent to whether it has an empty key().
	// Join, Intersect, and Union return false because it depends on strategy.
	fastSingle() bool
}

// Mode is the transaction context - cursor, read, or update.
// It affects the use of temporary indexes.
type Mode int

const (
	CursorMode Mode = iota
	ReadMode
	UpdateMode
)

func (mode Mode) String() string {
	switch mode {
	case CursorMode:
		return "cursorMode"
	case ReadMode:
		return "readMode"
	case UpdateMode:
		return "updateMode"
	default:
		panic("invalid mode")
	}
}

type Cost = int

type QueryTran interface {
	GetSchema(table string) *schema.Schema
	GetInfo(table string) *meta.Info
	GetAllInfo() []*meta.Info
	GetAllSchema() []*meta.Schema
	GetAllViews() []string
	GetView(string) string
	GetStore() *stor.Stor
	RangeFrac(table string, iIndex int, org, end string) float64
	Lookup(table string, iIndex int, key string) *runtime.DbRec
	Output(th *runtime.Thread, table string, rec runtime.Record)
	GetIndexI(table string, iIndex int) *index.Overlay
	GetRecord(off uint64) runtime.Record
	MakeLess(is *ixkey.Spec) func(x, y uint64) bool
	Read(string, int, string, string)
}

// Setup prepares a parsed query for execution.
// It calls Transform, Optimize, and SetApproach.
// The resulting Query is ready for execution.
func Setup(q Query, mode Mode, t QueryTran) (Query, Cost, Cost) {
	q = q.Transform()
	return setup(q, mode, 1, t)
}

// Setup1 is the same as Setup except it passes a frac of 1/nrows
// which will minimize fixed cost e.g. by avoiding temp indexes.
// It is used by DbmsLocal for Query1, QueryFirst, QueryLast
func Setup1(q Query, mode Mode, t QueryTran) (Query, Cost, Cost) {
	q = q.Transform()
	nrows, _ := q.Nrows()
	nrows = ord.Max(1, nrows) // avoid divide by zero
	return setup(q, mode, 1/float64(nrows), t)
}

func setup(q Query, mode Mode, frac float64, t QueryTran) (Query, Cost, Cost) {
	fixcost, varcost := Optimize(q, mode, nil, frac)
	if fixcost+varcost >= impossible {
		panic("invalid query: " + q.String())
	}
	q = SetApproach(q, nil, frac, t)
	return q, fixcost, varcost
}

// SetupKey is like Setup but it ensures a key index
func SetupKey(q Query, mode Mode, t QueryTran) Query {
	q = q.Transform()
	best := newBestIndex()
	for _, key := range q.Keys() {
		b := bestGrouped(q, mode, nil, 1, key)
		best.update(b.index, b.fixcost, b.varcost)
	}
	if best.fixcost+best.varcost >= impossible {
		panic("invalid query: " + q.String())
	}
	q = SetApproach(q, best.index, 1, t)
	return q
}

const outOfOrder = 10 // minimal penalty for executing out of order

const impossible = Cost(math.MaxInt / 64) // allow for adding impossible's

// Optimize determines the best (lowest estimated cost) query execution approach
func Optimize(q Query, mode Mode, index []string, frac float64) (
	fixcost, varcost Cost) {
	assert.That(!math.IsNaN(frac) && !math.IsInf(frac, 0))
	if fastSingle(q, index) {
		index = nil
	}
	if fixcost, varcost, _ := q.cacheGet(index, frac); varcost >= 0 {
		return fixcost, varcost
	}
	fixcost, varcost, app := optTempIndex(q, mode, index, frac)
	assert.That(fixcost >= 0 && varcost >= 0)
	q.cacheAdd(index, frac, fixcost, varcost, app)
	return fixcost, varcost
}

func fastSingle(q Query, index []string) bool {
	return q.fastSingle() && set.Subset(q.Columns(), index)
}

func optTempIndex(q Query, mode Mode, index []string, frac float64) (
	fixcost, varcost Cost, approach any) {
	traceQO := func(more ...any) {
		if trace.QueryOpt.On() {
			args := append([]any{index, frac, "="}, more...)
			trace.QueryOpt.Println(mode, args...)
			trace.Println(format(q, 1))
		}
	}
	if !set.Subset(q.Columns(), index) {
		traceQO("impossible index not a subset of columns")
		return impossible, impossible, nil
	}
	if len(index) == 0 || !tempIndexable(mode) {
		fixcost, varcost, approach = q.optimize(mode, index, frac)
		traceQO(fixcost + varcost)
		return fixcost, varcost, approach
	}
	noIndexFixCost, noIndexVarCost, noIndexApp := q.optimize(mode, nil, 1)
	assert.That(noIndexFixCost >= 0)
	assert.That(noIndexVarCost >= 0)
	noIndexCost := noIndexFixCost + noIndexVarCost
	if noIndexCost >= impossible {
		traceQO("impossible even without index")
		return impossible, impossible, nil
	}

	indexedFixCost, indexedVarCost, indexedApp := q.optimize(mode, index, frac)
	assert.That(indexedFixCost >= 0)
	assert.That(indexedVarCost >= 0)
	indexedCost := indexedFixCost + indexedVarCost

	nrows, _ := q.Nrows()
	assert.That(nrows >= 0)
	tempindexFixCost := noIndexCost + 1000 // ???
	tempindexFixCost += 100 * len(index)   // prefer fewer fields
	if nrows > 0 {
		fnrows := float64(nrows)
		tempindexFixCost += Cost(265 * fnrows * math.Log(fnrows)) // empirical
	}
	tempindexVarCost := Cost(frac * float64(nrows*100)) // ???
	if !q.SingleTable() {
		tempindexVarCost *= 2 // ???
	}
	tempindexCost := tempindexFixCost + tempindexVarCost

	if indexedCost <= tempindexCost {
		traceQO("indexed", indexedCost, "<=", tempindexCost)
		return indexedFixCost, indexedVarCost, indexedApp
	}
	traceQO("tempindex", tempindexCost, "<", indexedCost)
	// trace.Println("    noIndex", noIndexFixCost, noIndexVarCost,
	// 	"indexed", indexedFixCost, indexedVarCost)
	return tempindexFixCost, tempindexVarCost,
		&tempIndex{approach: noIndexApp, index: index,
			srcfixcost: noIndexFixCost, srcvarcost: noIndexVarCost}
}

type tempIndex struct {
	approach   any
	index      []string
	srcfixcost Cost
	srcvarcost Cost
}

func tempIndexable(mode Mode) bool {
	if mode == ReadMode {
		return true
	}
	if mode == CursorMode {
		return false
	}
	// else updateMode
	return true
	// BUG this matches jSuneido, but it is not correct.
	// A temp index allows reading deleted or old versions of records.
	// But there is a big performance penalty
	// especially from the key sort added by QueryApply.
}

func min3(fixcost1, varcost1 Cost, app1 any, fixcost2, varcost2 Cost, app2 any,
	fixcost3, varcost3 Cost, app3 any) (Cost, Cost, any) {
	fixcost, varcost, app := fixcost1, varcost1, app1
	if fixcost2+varcost2 < fixcost+varcost {
		fixcost, varcost, app = fixcost2, varcost2, app2
	}
	if fixcost3+varcost3 < fixcost+varcost {
		fixcost, varcost, app = fixcost3, varcost3, app3
	}
	return fixcost, varcost, app
}

func LookupCost(q Query, mode Mode, index []string, nrows int) (
	Cost, Cost) {
	if fastSingle(q, index) {
		index = nil
	}
	fixcost, varcost := Optimize(q, mode, index, 0)
	if fixcost+varcost >= impossible {
		return impossible, impossible
	}
	var lookupCost Cost
	_, _, approach := q.cacheGet(index, 0)
	if _, ok := approach.(*tempIndex); ok {
		if q.SingleTable() {
			lookupCost = 200 // ???
		} else {
			lookupCost = 400 // ???
		}
	} else {
		lookupCost = q.lookupCost()
	}
	lookupCost *= nrows
	// trace.Println("LookupCost", fixcost, "+", lookupCost, "=", fixcost+lookupCost)
	return fixcost, lookupCost
}

// SetApproach finalizes the chosen approach.
// It also adds temp indexes where required.
func SetApproach(q Query, index []string, frac float64, tran QueryTran) Query {
	if fastSingle(q, index) {
		index = nil
	}
	fixcost, varcost, approach := q.cacheGet(index, frac)
	assert.That(fixcost >= 0)
	assert.That(varcost >= 0)
	if app, ok := approach.(*tempIndex); ok {
		q.cacheSetCost(1, app.srcfixcost, app.srcvarcost)
		q.setApproach(nil, 1, app.approach, tran)
		ti := NewTempIndex(q, app.index, tran)
		ti.cacheSetCost(frac, fixcost, varcost)
		return ti
	}
	q.cacheSetCost(frac, fixcost, varcost)
	q.setApproach(index, frac, approach, tran)
	return q
}

// Query1 -----------------------------------------------------------

type Query1 struct {
	cache
	source Query
}

func (q1 *Query1) Columns() []string {
	return q1.source.Columns()
}

func (q1 *Query1) Keys() [][]string {
	return q1.source.Keys()
}

func (q1 *Query1) fastSingle() bool {
	return q1.source.fastSingle()
}

func (q1 *Query1) Indexes() [][]string {
	return q1.source.Indexes()
}

func (q1 *Query1) Nrows() (int, int) {
	return q1.source.Nrows()
}

func (q1 *Query1) rowSize() int {
	return q1.source.rowSize()
}

func (*Query1) Ordering() []string {
	return nil
}

func (q1 *Query1) Fixed() []Fixed {
	return q1.source.Fixed()
}

func (q1 *Query1) Updateable() string {
	return q1.source.Updateable()
}

func (q1 *Query1) SingleTable() bool {
	return q1.source.SingleTable()
}

func (q1 *Query1) SetTran(t QueryTran) {
	q1.source.SetTran(t)
}

func (q1 *Query1) optimize(mode Mode, index []string, frac float64) (
	Cost, Cost, any) {
	fixcost, varcost := Optimize(q1.source, mode, index, frac)
	return fixcost, varcost, nil
}

func (q1 *Query1) lookupCost() Cost {
	return q1.source.lookupCost()
}

// Lookup default applies to Summarize and Sort
func (*Query1) Lookup(*runtime.Thread, []string, []string) runtime.Row {
	panic("Lookup not implemented")
}

func (q1 *Query1) Header() *runtime.Header {
	return q1.source.Header()
}

func (q1 *Query1) Output(th *runtime.Thread, rec runtime.Record) {
	q1.source.Output(th, rec)
}

func (q1 *Query1) Rewind() {
	q1.source.Rewind()
}

type q1i interface {
	Source() Query
	stringOp() string
}

func (q1 *Query1) Source() Query {
	return q1.source
}

// Query2 -----------------------------------------------------------

type Query2 struct {
	cache
	source1 Query
	source2 Query
}

func (q2 *Query2) String2(op string) string {
	return parenQ2(q2.source1) + " " + op + " " + paren(q2.source2)
}

func (q2 *Query2) SetTran(t QueryTran) {
	q2.source1.SetTran(t)
	q2.source2.SetTran(t)
}

func (q2 *Query2) Header() *runtime.Header {
	return runtime.JoinHeaders(q2.source1.Header(), q2.source2.Header())
}

func (q2 *Query2) Updateable() string {
	return ""
}

func (q2 *Query2) SingleTable() bool {
	return false // not single
}

func (*Query2) Ordering() []string {
	return nil
}

func (*Query2) Output(*runtime.Thread, runtime.Record) {
	panic("can't output to this query")
}

func (q2 *Query2) keypairs() [][]string {
	var keys [][]string
	for _, k1 := range q2.source1.Keys() {
		for _, k2 := range q2.source2.Keys() {
			keys = set.AddUniqueFn(keys, set.Union(k1, k2), set.Equal[string])
		}
	}
	assert.That(len(keys) != 0)
	return keys
}

type q2i interface {
	q1i
	Source2() Query
}

func (q2 *Query2) Source() Query {
	return q2.source1
}

func (q2 *Query2) Source2() Query {
	return q2.source2
}

//-------------------------------------------------------------------

// paren is a helper for Query String methods
func paren(q Query) string {
	switch q.(type) {
	case *Table, *Tables, *Columns, *Indexes, *Nothing:
		return q.String()
	}
	return "(" + q.String() + ")"
}

func parenQ2(q Query) string {
	if _, ok := q.(q2i); ok {
		return "(" + q.String() + ")"
	}
	return q.String()
}

//-------------------------------------------------------------------

type bestIndex struct {
	index   []string
	fixcost Cost
	varcost Cost
}

func newBestIndex() bestIndex {
	return bestIndex{fixcost: impossible, varcost: impossible}
}

func (bi *bestIndex) update(index []string, fixcost, varcost Cost) {
	if fixcost+varcost < bi.fixcost+bi.varcost {
		*bi = bestIndex{index: index, fixcost: fixcost, varcost: varcost}
	}
}

func (bi *bestIndex) cost() int {
	return bi.fixcost + bi.varcost
}

func (bi *bestIndex) String() string {
	if bi.cost() >= impossible {
		return "impossible"
	}
	return fmt.Sprint("{", bi.index, " ",
		trace.Number(bi.fixcost), " + ", trace.Number(bi.varcost),
		" = ", trace.Number(bi.cost()), "}")
}

// bestGrouped finds the best index with cols (in any order) as a prefix
// taking fixed into consideration.
// It is used by Project, Summarize, and Join.
func bestGrouped(source Query, mode Mode, index []string, frac float64, cols []string) bestIndex {
	var indexes [][]string
	if index == nil {
		indexes = source.Indexes()
	} else {
		indexes = [][]string{index}
	}
	fixed := source.Fixed()
	nColsUnfixed := countUnfixed(cols, fixed)
	best := newBestIndex()
	for _, idx := range indexes {
		if grouped(idx, cols, nColsUnfixed, fixed) {
			fixcost, varcost := Optimize(source, mode, idx, frac)
			best.update(idx, fixcost, varcost)
		}
	}
	if index == nil {
		fixcost, varcost := Optimize(source, mode, cols, frac)
		best.update(cols, fixcost, varcost)
	}
	return best
}

func countUnfixed(cols []string, fixed []Fixed) int {
	nunfixed := 0
	for _, col := range cols {
		if !isFixed(fixed, col) {
			nunfixed++
		}
	}
	return nunfixed
}

// grouped returns whether an index has cols (in any order) as a prefix
// taking fixed into consideration
func grouped(index []string, cols []string, nColsUnfixed int, fixed []Fixed) bool {
	if len(index) < nColsUnfixed {
		return false
	}
	n := 0
	for _, col := range index {
		if isFixed(fixed, col) {
			continue
		}
		if !slices.Contains(cols, col) {
			return false
		}
		n++
		if n == nColsUnfixed {
			return true
		}
	}
	return false
}

// ordered returns whether an index supplies an order
// taking fixed into consideration.
func ordered(index []string, order []string, fixed []Fixed) bool {
	i := 0
	o := 0
	in := len(index)
	on := len(order)
	for i < in && o < on {
		if index[i] == order[o] {
			o++
			i++
		} else if isFixed(fixed, index[i]) {
			i++
		} else if isFixed(fixed, order[o]) {
			o++
		} else {
			return false
		}
	}
	for o < on && isFixed(fixed, order[o]) {
		o++
	}
	return o >= on
}

func withoutDupsOrSupersets(keys [][]string) [][]string {
	om := newOptMod(keys)
outer:
	for _, k1 := range keys {
		for _, k2 := range keys {
			if len(k1) > len(k2) && set.Subset(k1, k2) {
				continue outer // skip/exclude k1 - superset
			}
			if !slc.ContainsFn(om.result(), k1, set.Equal[string]) { // exclude duplicates
				om.add(k1)
			}
		}
	}
	return om.result()
}

// optmod is useful when building a new version
// which is likely to be the same as the original.
// It avoids constructing a new version unless there are changes,
// without having to redundantly check in advance.
type optmod struct {
	orig [][]string
	i    int
	mod  [][]string
}

func newOptMod(orig [][]string) *optmod {
	return &optmod{orig: orig}
}

func (b *optmod) add(x []string) {
	if b.mod == nil {
		if b.i < len(b.orig) && set.Equal(x, b.orig[b.i]) {
			b.i++ // same as orig
			return
		}
		b.mod = append(b.mod, b.orig[:b.i]...)
	}
	b.mod = append(b.mod, x)
}

func (b *optmod) result() [][]string {
	if b.mod == nil {
		return b.orig[:b.i:b.i]
	}
	return slices.Clip(b.mod)
}

// ------------------------------------------------------------------

func Format(q Query) string {
	return format(q, 0)
}

const indent1 = "    "

func format(q Query, indent int) string { // recursive
	in := strings.Repeat(indent1, indent)
	nrows, pop := q.Nrows()
	frac, fixcost, varcost := q.cacheCost()
	cost := "{"
	if frac != 1 {
		cost += fmt.Sprintf("%.3fx ", frac)
	}
	cost += trace.Number(nrows)
	if nrows != pop {
		cost += "/" + trace.Number(pop)
	}
	if fixcost+varcost > 0 {
		cost += " " + trace.Number(fixcost) + "+" + trace.Number(varcost)
	}
	cost += "} "
	switch q := q.(type) {
	case q2i:
		return format(q.Source(), indent+1) + "\n" +
			in + cost + q.stringOp() + "\n" +
			format(q.Source2(), indent+1)
	case q1i:
		return format(q.Source(), indent) + "\n" +
			in + cost + q.stringOp()
	default:
		return in + cost + q.String()
	}
}

//lint:ignore U1000 for debugging
func unpack(packed []string) []runtime.Value {
	vals := make([]runtime.Value, len(packed))
	for i, p := range packed {
		if p == ixkey.Max {
			vals[i] = runtime.SuStr("<max>")
		} else {
			vals[i] = runtime.Unpack(p)
		}
	}
	return vals
}
