// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package query implements query parsing, optimization, and execution.
/*
	Query
		Table
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
					Intersect
					Minus
					Union
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
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"golang.org/x/exp/slices"
)

type Query interface {
	// Columns is all the available columns, including derived
	Columns() []string

	// Transform refactors the query for more efficient execution.
	// This stage is not cost based, transforms are applied when possible.
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
	// from the entire table with p rows
	Nrows() (n, p int)

	// rowSize returns the average number of bytes per row
	rowSize() int

	// Lookup returns the row matching the given key value, or nil if not found.
	// It is used by Compatible (Intersect, Minus, Union). See also: Select
	Lookup(th *runtime.Thread, cols, vals []string) runtime.Row

	Rewind()

	Get(th *runtime.Thread, dir runtime.Dir) runtime.Row

	// Select restricts the query to records matching the given packed values.
	// It is used by Join and LeftJoin. See also: Lookup
	Select(cols, vals []string)

	Header() *runtime.Header
	Output(th *runtime.Thread, rec runtime.Record)

	String() string

	cacheAdd(mode Mode, index []string, fixcost, varcost Cost, approach any)

	// cacheGet returns the cost and approach associated with an index
	// or -1 if the index has not been added.
	cacheGet(mode Mode, index []string) (fixcost, varcost Cost, approach any)

	cacheSetCost(fixcost, varcost Cost)
	cacheCost() (fixcost, varcost Cost)

	optimize(mode Mode, index []string) (fixcost, varcost Cost, approach any)
	setApproach(mode Mode, index []string, approach any, tran QueryTran)

	// lookupCost returns the cost of one Lookup
	lookupCost() Cost
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
//
// NOTE: Correct usage is: q, cost = Setup(q, mode, t)
func Setup(q Query, mode Mode, t QueryTran) (Query, Cost, Cost) {
	q = q.Transform()
	fixcost, varcost := Optimize(q, mode, nil)
	if fixcost+varcost >= impossible {
		panic("invalid query: " + q.String())
	}
	q = SetApproach(q, mode, nil, t)
	return q, fixcost, varcost
}

// SetupKey is like Setup but it ensures a key index
func SetupKey(q Query, mode Mode, t QueryTran) Query {
	q = q.Transform()
	best := newBestIndex()
	for _, key := range q.Keys() {
		b := bestGrouped(q, mode, nil, key)
		best.update(b.index, b.fixcost, b.varcost)
	}
	if best.fixcost+best.varcost >= impossible {
		panic("invalid query: " + q.String())
	}
	q = SetApproach(q, mode, best.index, t)
	return q
}

const outOfOrder = 10 // minimal penalty for executing out of order

const impossible = Cost(math.MaxInt / 64) // allow for adding impossible's

// Optimize determines the best (lowest estimated cost) query execution approach
func Optimize(q Query, mode Mode, index []string) (fixcost, varcost Cost) {
	if isSingleton(q) {
		index = nil
	}
	if fixcost, varcost, _ := q.cacheGet(mode, index); varcost >= 0 {
		return fixcost, varcost
	}
	fixcost, varcost, app := optTempIndex(q, mode, index)
	assert.That(fixcost >= 0)
	assert.That(varcost >= 0)
	q.cacheAdd(mode, index, fixcost, varcost, app)
	return fixcost, varcost
}

func isSingleton(q Query) bool {
	keys := q.Keys()
	return len(keys) == 1 && len(keys[0]) == 0
}

func optTempIndex(q Query, mode Mode, index []string) (
	fixcost, varcost Cost, approach any) {
	traceQO := func(more ...any) {
		if trace.QueryOpt.On() {
			args := append([]any{index, "="}, more...)
			trace.QueryOpt.Println(mode, args...)
			trace.Println(format(q, 1))
		}
	}
	if !set.Subset(q.Columns(), index) {
		traceQO("impossible index not a subset of columns")
		return impossible, impossible, nil
	}
	if len(index) == 0 || !tempIndexable(mode) {
		fixcost, varcost, approach = q.optimize(mode, index)
		traceQO(fixcost + varcost)
		assert.That(fixcost >= 0)
		assert.That(varcost >= 0)
		return fixcost, varcost, approach
	}
	noIndexFixCost, noIndexVarCost, noIndexApp := q.optimize(mode, nil)
	assert.That(noIndexFixCost >= 0)
	assert.That(noIndexVarCost >= 0)
	noIndexCost := noIndexFixCost + noIndexVarCost
	if noIndexCost >= impossible {
		traceQO("impossible even without index")
		return impossible, impossible, nil
	}

	indexedFixCost, indexedVarCost, indexedApp := q.optimize(mode, index)
	assert.That(indexedFixCost >= 0)
	assert.That(indexedVarCost >= 0)
	indexedCost := indexedFixCost + indexedVarCost

	nrows, _ := q.Nrows()
	assert.That(nrows >= 0)
	tempindexFixCost := 0
	if nrows > 0 {
		fn := float64(nrows)
		tempindexFixCost = noIndexCost + int(265*fn*math.Log(fn)) // empirical
		assert.Msg(fn, math.Log(fn)).That(tempindexFixCost >= 0)
	}
	tempindexVarCost := 250 * nrows
	if !q.SingleTable() {
		tempindexVarCost *= 2 // ???
	}
	tempindexCost := tempindexFixCost + tempindexVarCost

	if indexedCost <= tempindexCost {
		traceQO("indexed", indexedCost, "<=", tempindexCost)
		return indexedFixCost, indexedVarCost, indexedApp
	}
	traceQO("tempindex", tempindexCost, "<", indexedCost)
	return tempindexFixCost, tempindexVarCost,
		&tempIndex{approach: noIndexApp, index: index,
			fixcost: tempindexFixCost, varcost: tempindexVarCost}
}

type tempIndex struct {
	approach any
	index    []string
	fixcost  Cost
	varcost  Cost
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

func LookupCost(q Query, mode Mode, index []string, nrows int) (Cost, Cost) {
	if isSingleton(q) {
		index = nil
	}
	fixcost, varcost := Optimize(q, mode, index)
	if fixcost+varcost >= impossible {
		return impossible, impossible
	}
	lookupCost := q.lookupCost() * nrows
	// trace.Println("LookupCost", fixcost, "+", varcost, "+", lookupCost)
	return fixcost, varcost + lookupCost
}

// SetApproach locks in the best approach.
// It also adds temp indexes where required.
func SetApproach(q Query, mode Mode, index []string, tran QueryTran) Query {
	if isSingleton(q) {
		index = nil
	}
	fixcost, varcost, approach := q.cacheGet(mode, index)
	assert.That(fixcost >= 0)
	assert.That(varcost >= 0)
	q.cacheSetCost(fixcost, varcost)
	if app, ok := approach.(*tempIndex); ok {
		q.setApproach(mode, nil, app.approach, tran)
		return &TempIndex{Query1: Query1{source: q,
			cache: cache{fixcost: app.fixcost, varcost: app.varcost}},
			order: app.index, tran: tran}
	}
	q.setApproach(mode, index, approach, tran)
	return q
}

// Query1 -----------------------------------------------------------

type Query1 struct {
	cache
	source Query
}

func (q1 *Query1) String() string {
	panic("should be overridden")
}

func (q1 *Query1) Columns() []string {
	return q1.source.Columns()
}

func (q1 *Query1) Keys() [][]string {
	return q1.source.Keys()
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

func (q1 *Query1) optimize(mode Mode, index []string) (Cost, Cost, any) {
	fixcost, varcost := Optimize(q1.source, mode, index)
	return fixcost, varcost, nil
}

func (q1 *Query1) setApproach(Mode, []string, any, QueryTran) {
	assert.ShouldNotReachHere()
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
	source  Query
	source2 Query
}

func (q2 *Query2) String2(op string) string {
	return parenQ2(q2.source) + " " + op + " " + paren(q2.source2)
}

func (q2 *Query2) SetTran(t QueryTran) {
	q2.source.SetTran(t)
	q2.source2.SetTran(t)
}

func (q2 *Query2) Header() *runtime.Header {
	return runtime.JoinHeaders(q2.source.Header(), q2.source2.Header())
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
	for _, k1 := range q2.source.Keys() {
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
	return q2.source
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
func bestGrouped(source Query, mode Mode, index, cols []string) bestIndex {
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
			fixcost, varcost := Optimize(source, mode, idx)
			best.update(idx, fixcost, varcost)
		}
	}
	if best.index == nil && index == nil {
		best.fixcost, best.varcost = Optimize(source, mode, cols)
		if best.fixcost+best.varcost < impossible {
			best.index = cols
		}
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
	fixcost, varcost := q.cacheCost()
	cost := "{" + trace.Number(nrows)
	if nrows != pop {
		cost += "/" + trace.Number(pop)
	}
	cost += " " + trace.Number(fixcost) + "+" + trace.Number(varcost) + "} "
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
