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

	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/ssset"
)

type Query interface {
	Init()

	// Columns is all the available columns, including derived
	Columns() []string

	// Transform refactors the query for more efficient execution.
	// This stage is not cost based, transforms are applied whenever possible.
	Transform() Query

	// SetTran sets the transaction to be used by the query
	SetTran(tran QueryTran)

	// Header() runtime.Header

	// Order() []string

	// Fixed returns the field values that are constant from Extend or Where
	Fixed() []Fixed

	// Updateable returns whether the rows from the query can be updated
	Updateable() bool

	// Indexes returns all the indexes
	Indexes() [][]string

	// Keys returns the indexes that are keys
	Keys() [][]string

	// nrows returns the number of rows in the query, estimated for operations
	nrows() int

	// rowSize returns the average number of bytes per row
	rowSize() int

	// Lookup returns the row matching the given key, or nil if not found
	// Lookup(index []string, key string) runtime.Row

	// Rewind()

	// Get(dir runtime.Dir) runtime.Row

	String() string

	cacheAdd(index []string, cost Cost, approach interface{})

	// cacheGet returns the cost and approach associated with an index
	// or -1 if the index as not been added.
	cacheGet(index []string) (Cost, interface{})

	optimize(mode Mode, index []string) (cost Cost, approach interface{})
	setApproach(index []string, approach interface{}, tran QueryTran)

	lookupCost() Cost
}

// Mode is the transaction context - cursor, read, or update.
// It affects the use of temporary indexes.
type Mode int

const (
	cursorMode Mode = iota
	readMode
	updateMode
)

func (mode Mode) String() string {
	switch mode {
	case cursorMode:
		return "cursorMode"
	case readMode:
		return "readMode"
	case updateMode:
		return "updateMode"
	default:
		panic("invalid mode")
	}
}

type Cost = int

type QueryTran interface {
	GetSchema(table string) *schema.Schema
	GetInfo(table string) *meta.Info
}

// Setup prepares a parsed query for execution.
// It calls transform and optimize.
// The resulting Query is ready for execution.
func Setup(q Query, mode Mode, t QueryTran) Query {
	q.SetTran(t)
	q.Init()
	q = q.Transform()
	cost := Optimize(q, mode, nil)
	if cost >= impossible {
		panic("invalid query " + q.String())
	}
	q = SetApproach(q, nil, t)
	return q
}

const outOfOrder = 10 // minimal penalty for executing out of order

const impossible = Cost(ints.MaxInt / 64) // allow for adding IMPOSSIBLE's

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
	tempIndexReadCost := q.nrows() * btree.EntrySize
	tempIndexWriteCost := tempIndexReadCost * 2 // surcharge for memory ???
	dataReadCost := q.nrows() * q.rowSize()
	tempIndexCost := tempIndexWriteCost + tempIndexReadCost + dataReadCost
	cost2 := noIndexCost + tempIndexCost
	assert.That(cost2 >= 0)
	trace("cost1", cost1, "cost2 (temp index)", cost2)
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
	if mode == readMode {
		return true
	}
	if mode == cursorMode {
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
		panic(fmt.Sprintln("MISSING", q, index))
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

func (q1 *Query1) Init() {
	q1.source.Init()
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

func (q1 *Query1) nrows() int {
	return q1.source.nrows()
}

func (q1 *Query1) rowSize() int {
	return q1.source.rowSize()
}

func (q1 *Query1) Fixed() []Fixed {
	return q1.source.Fixed()
}

func (q1 *Query1) Updateable() bool {
	return q1.source.Updateable()
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

// bestPrefixed returns the best index that supplies the required order
// taking fixed into account
func (q1 *Query1) bestPrefixed(indexes [][]string, order []string,
	mode Mode) bestIndex {
	var best bestIndex
	best.cost = impossible
	fixed := q1.source.Fixed()
	for _, ix := range indexes {
		if q1.prefixed(ix, order, fixed) {
			cost := Optimize(q1.source, mode, ix)
			best.update(ix, cost)
		}
	}
	return best
}

type bestIndex struct {
	index []string
	cost  Cost
}

func (bi *bestIndex) update(index []string, cost Cost) {
	if bi.index == nil || cost < bi.cost {
		bi.index = index
		bi.cost = cost
	}
}

// prefixed returns whether an index supplies an order, given what's fixed
func (*Query1) prefixed(index []string, order []string, fixed []Fixed) bool {
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

func (q2 *Query2) Init() {
	q2.source.Init()
	q2.source2.Init()
}

func (q2 *Query2) SetTran(t QueryTran) {
	q2.source.SetTran(t)
	q2.source2.SetTran(t)
}

func (q2 *Query2) Updateable() bool {
	return false
}

func (q2 *Query2) optimize(Mode, []string) Cost {
	panic("should be overridden")
}

func (q2 *Query2) keypairs() [][]string {
	var keys [][]string
	for _, k1 := range q2.source.Keys() {
		for _, k2 := range q2.source2.Keys() {
			keys = ssset.AddUnique(keys, sset.Union(k1, k2))
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

// paren is a helper to String methods
func paren(q Query) string {
	if tbl, ok := q.(*Table); ok {
		return tbl.String()
	}
	return "(" + q.String() + ")"
}

func parenQ2(q Query) string {
	if _, ok := q.(q2i); ok {
		return "(" + q.String() + ")"
	}
	return q.String()
}

func bestKey(q Query, mode Mode) []string {
	var best []string
	bestCost := impossible
	for _, key := range q.Keys() {
		cost := Optimize(q, mode, key)
		cost += (len(key) - 1) * cost / 20 // ??? prefer shorter keys
		if cost < bestCost {
			best = key
		}
	}
	return best
}
