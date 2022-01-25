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

The cost model is based on the number of bytes read.
*/
package query

import (
	"fmt"
	"math"

	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/setset"
	"github.com/apmckinlay/gsuneido/util/sset"
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

	// Nrows returns the number of rows in the query, estimated for operations
	Nrows() int

	// rowSize returns the average number of bytes per row
	rowSize() int

	// Lookup returns the row matching the given key value, or nil if not found.
	// It is used by Compatible (Intersect, Minus, Union). See also: Select
	Lookup(cols, vals []string) runtime.Row

	Rewind()

	Get(dir runtime.Dir) runtime.Row

	// Select restricts the query to records matching the given packed values.
	// It is used by Join and LeftJoin. See also: Lookup
	Select(cols, vals []string)

	Header() *runtime.Header
	Output(rec runtime.Record)

	String() string

	cacheAdd(index []string, cost Cost, approach interface{})

	// cacheGet returns the cost and approach associated with an index
	// or -1 if the index has not been added.
	cacheGet(index []string) (Cost, interface{})

	optimize(mode Mode, index []string) (cost Cost, approach interface{})
	setApproach(index []string, approach interface{}, tran QueryTran)

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
	RangeFrac(table string, iIndex int, org, end string) float64
	Lookup(table string, iIndex int, key string) *runtime.DbRec
	Output(table string, rec runtime.Record)
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
func Setup(q Query, mode Mode, t QueryTran) (Query, Cost) {
	q = q.Transform()
	cost := Optimize(q, mode, nil)
	if cost >= impossible {
		panic("invalid query: " + q.String())
	}
	q = SetApproach(q, nil, t)
	return q, cost
}

// SetupKey is like Setup but it ensures a key index
func SetupKey(q Query, mode Mode, t QueryTran) Query {
	q = q.Transform()
	best := newBestIndex()
	for _, key := range q.Keys() {
		b := bestGrouped(q, mode, nil, key)
		best.update(b.index, b.cost)
	}
	if best.cost >= impossible {
		panic("invalid query: " + q.String())
	}
	q = SetApproach(q, best.index, t)
	return q
}

const outOfOrder = 10 // minimal penalty for executing out of order

const impossible = Cost(math.MaxInt / 64) // allow for adding IMPOSSIBLE's

// gin is used with be e.g. defer be(gin(...))
func gin(args ...interface{}) string {
	trace(args...)
	indent++
	return args[0].(string)
}
func trace(args ...interface{}) {
	// fmt.Print(strings.Repeat(" ", 4*indent))
	// fmt.Println(args...)
}
func be(arg string) {
	// fmt.Println(strings.Repeat(" ", 4*indent)+"end", arg)
	indent--
}

var indent = 0

// Optimize determines the best (lowest estimated cost) query execution approach
func Optimize(q Query, mode Mode, index []string) (cost Cost) {
	if cost, _ = q.cacheGet(index); cost >= 0 {
		return cost
	}
	cost, app := optTempIndex(q, mode, index)
	q.cacheAdd(index, cost, app)
	return cost
}

func optTempIndex(q Query, mode Mode, index []string) (
	cost Cost, approach interface{}) {
	defer be(gin("Optimize", q, mode, index))
	defer func() { trace("=>", cost) }()
	if !sset.Subset(q.Columns(), index) {
		return impossible, nil
	}
	if index == nil || !tempIndexable(q, mode) {
		return q.optimize(mode, index)
	}
	cost1, app1 := q.optimize(mode, index)
	noIndexCost, app2 := q.optimize(mode, nil)
	tempIndexReadCost := q.Nrows() * btree.EntrySize
	tempIndexWriteCost := tempIndexReadCost * 2 // ??? surcharge (for memory)
	dataReadCost := q.Nrows() * q.rowSize()
	tempIndexCost := tempIndexWriteCost + tempIndexReadCost + dataReadCost
	cost2 := noIndexCost + tempIndexCost
	assert.That(cost2 >= 0) // ???
	trace("cost1", cost1, "noIndexCost",
		noIndexCost, "tempIndexCost", tempIndexCost, "cost2", cost2)
	cost, approach = min(cost1, app1, cost2, app2)
	if cost >= impossible {
		return impossible, nil
	}
	if cost2 < cost1 {
		approach = &tempIndex{approach: approach, index: index}
	}
	return cost, approach
}

type tempIndex struct {
	approach interface{}
	index    []string
}

func tempIndexable(q Query, mode Mode) bool {
	if mode == ReadMode {
		return true
	}
	if mode == CursorMode {
		return false
	}
	// else updateMode, tempIndex only directly on Table
	_, ok := q.(*Table)
	return ok
}

func min(cost1 Cost, app1 interface{}, cost2 Cost, app2 interface{}) (
	Cost, interface{}) {
	if cost1 <= cost2 {
		return cost1, app1
	}
	return cost2, app2
}

func min3(cost1 Cost, app1 interface{}, cost2 Cost, app2 interface{},
	cost3 Cost, app3 interface{}) (Cost, interface{}) {
	cost, app := cost1, app1
	if cost2 < cost {
		cost, app = cost2, app2
	}
	if cost3 < cost {
		cost, app = cost3, app3
	}
	return cost, app
}

func LookupCost(q Query, mode Mode, index []string, nrows int) Cost {
	if Optimize(q, mode, index) >= impossible {
		return impossible
	}
	return q.lookupCost() * nrows
}

// SetApproach locks in the best approach.
// It also adds temp indexes where required.
func SetApproach(q Query, index []string, tran QueryTran) Query {
	cost, approach := q.cacheGet(index)
	if cost < 0 {
		panic(fmt.Sprintln("NOT IN CACHE:", q, index))
	}
	assert.That(cost >= 0)
	if app, ok := approach.(*tempIndex); ok {
		q.setApproach(nil, app.approach, tran)
		return &TempIndex{Query1: Query1{source: q}, order: app.index, tran: tran}
	}
	q.setApproach(index, approach, tran)
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

func (q1 *Query1) Nrows() int {
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

func (q1 *Query1) optimize(mode Mode, index []string) (Cost, interface{}) {
	return Optimize(q1.source, mode, index), nil
}

func (q1 *Query1) setApproach(_ []string, _ interface{}, _ QueryTran) {
	panic("shouldn't reach here")
}

func (q1 *Query1) lookupCost() Cost {
	return q1.source.lookupCost()
}

func (*Query1) Lookup([]string, []string) runtime.Row {
	panic("Lookup not implemented")
}

func (q1 *Query1) Header() *runtime.Header {
	return q1.source.Header()
}

func (q1 *Query1) Output(rec runtime.Record) {
	q1.source.Output(rec)
}

func (*Query1) Get(runtime.Dir) runtime.Row {
	panic("Get not implemented")
}

func (q1 *Query1) Rewind() {
	q1.source.Rewind()
}

// Query2 -----------------------------------------------------------

type Query2 struct {
	Query1
	source2 Query
}

func (q2 *Query2) Keys() [][]string {
	panic("should be overridden")
}

func (q2 *Query2) Indexes() [][]string {
	panic("should be overridden")
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

func (*Query2) Output(runtime.Record) {
	panic("can't output to this query")
}

func (q2 *Query2) optimize(Mode, []string) Cost {
	panic("should be overridden")
}

func (q2 *Query2) keypairs() [][]string {
	var keys [][]string
	for _, k1 := range q2.source.Keys() {
		for _, k2 := range q2.source2.Keys() {
			keys = setset.AddUnique(keys, sset.Union(k1, k2))
		}
	}
	assert.That(len(keys) != 0)
	return keys
}

type q2i interface {
	tagQuery2()
}

func (*Query2) tagQuery2() {
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
	index []string
	cost  Cost
}

func newBestIndex() bestIndex {
	return bestIndex{cost: impossible}
}

func (bi *bestIndex) update(index []string, cost Cost) {
	if cost < bi.cost {
		bi.index = index
		bi.cost = cost
	}
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
			cost := Optimize(source, mode, idx)
			best.update(idx, cost)
		}
	}
	if best.index == nil && index == nil {
		best.cost = Optimize(source, mode, cols)
		if best.cost < impossible {
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
		if !sset.Contains(cols, col) {
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
