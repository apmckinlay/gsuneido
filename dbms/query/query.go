// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package query implements query parsing, optimization, and execution.
/*
	Query
		Table
		Query1
			Where
			Extend
			Rename
			Project / Remove
			Summarize
			TempIndex
			Sort
			Query2
				Compatible
					Union
					Intersect
					Minus
				Join
					LeftJoin
				Times

The cost model is based on the number of bytes read.
*/
package query

import (
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

	// dataSize returns the number of bytes of data, estimated for operations
	dataSize() int

	// Lookup returns the row matching the given key, or nil if not found
	// Lookup(index []string, key string) runtime.Row

	// Rewind()

	// Get(dir runtime.Dir) runtime.Row

	String() string

	cacheAdd(index []string, cost Cost)
	cacheGet(index []string) Cost

	optimize(mode Mode, index []string, act action) Cost
	addTempIndex(tran QueryTran) Query
	lookupCost() Cost //TODO temporary, combine with optimize

	setTempIndex(index []string)
	getTempIndex() []string
}

// Mode is the transaction context - cursor, read, or update.
// It affects the use of temporary indexes.
type Mode int

const (
	cursorMode Mode = iota
	readMode
	updateMode
)

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
	cost := Optimize(q, mode, nil, freeze)
	if cost >= impossible {
		panic("invalid query " + q.String())
	}
	q = q.addTempIndex(t)
	return q
}

const outOfOrder = 10 // minimal penalty for executing out of order

const impossible = Cost(ints.MaxInt / 64) // allow for adding IMPOSSIBLE's

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

// Optimize is the start of optimization.
// It handles whether or not to use a temp index.
func Optimize(q Query, mode Mode, index []string, act action) Cost {
	defer be(gin("Optimize", q, index, act))
	if !sset.Subset(q.Columns(), index) {
		return impossible
	}
	if index == nil || !tempIndexable(q, mode) {
		return optimize1(q, mode, index, act)
	}
	cost1 := optimize1(q, mode, index, assess)
	noIndexCost := optimize1(q, mode, nil, assess)
	tempIndexReadCost := q.nrows() * btree.EntrySize
	tempIndexWriteCost := tempIndexReadCost // ???
	tempIndexCost := tempIndexWriteCost + tempIndexReadCost + q.dataSize()
	cost2 := noIndexCost + tempIndexCost
	assert.That(cost2 >= 0)
	trace("cost1", cost1, "cost2 (temp index)", cost2)

	cost := ints.Min(cost1, cost2)
	if cost >= impossible {
		return impossible
	}
	if act == freeze {
		if cost2 < cost1 {
			q.setTempIndex(index)
			index = nil
		}
		optimize1(q, mode, index, freeze)
	}
	return cost
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

type action byte

const (
	assess action = 0
	freeze action = 1
)

func (act action) String() string {
	if act == freeze {
		return "freeze"
	}
	return "assess"
}

// optimize1 handles caching
func optimize1(q Query, mode Mode, index []string, act action) (cost Cost) {
	defer be(gin("optimize1", q, index, act))
	defer func() { trace("=>", cost) }()
	if act == freeze {
		return q.optimize(mode, index, freeze)
	}
	if cost = q.cacheGet(index); cost >= 0 {
		return cost
	}
	cost = q.optimize(mode, index, assess)
	assert.That(cost >= 0)
	q.cacheAdd(index, cost)
	return cost
}

func addTempIndex(q Query, tran QueryTran) Query {
	if ti := q.getTempIndex(); ti != nil {
		return &TempIndex{Query1: Query1{source: q}, order: ti, tran: tran}
	}
	return q
}

// Query1 -----------------------------------------------------------

type Query1 struct {
	cache
	useTempIndex
	source Query
}

func (q1 *Query1) String() string {
	return q1.source.String()
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

func (q1 *Query1) dataSize() int {
	return q1.source.dataSize()
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

func (q1 *Query1) Transform() Query {
	panic("should be overridden")
}

func (q1 *Query1) optimize(mode Mode, index []string, act action) Cost {
	return Optimize(q1.source, mode, index, act)
}

func (q1 *Query1) addTempIndex(tran QueryTran) Query {
	q1.source = addTempIndex(q1.source, tran)
	return addTempIndex(q1, tran)
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
			// optimize1 to bypass tempindex
			cost := optimize1(q1.source, mode, ix, assess)
			best.update(ix, cost)
		}
	}
	return best
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
	return paren(q2.source) + " " + op + " " + paren(q2.source2)
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

func (q2 *Query2) Transform() Query {
	panic("should be overridden")
}

func (q2 *Query2) optimize(Mode, []string, action) Cost {
	panic("should be overridden")
}

func (q2 *Query2) addTempIndex(tran QueryTran) Query {
	q2.source = addTempIndex(q2.source, tran)
	q2.source2 = addTempIndex(q2.source2, tran)
	return addTempIndex(q2, tran)
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

// temp index -------------------------------------------------------

type useTempIndex struct {
	tempIndex []string
}

func (uti *useTempIndex) setTempIndex(index []string) {
	uti.tempIndex = index
}

func (uti *useTempIndex) getTempIndex() []string {
	return uti.tempIndex
}

//-------------------------------------------------------------------

// paren is a helper to String methods
func paren(q Query) string {
	if tbl, ok := q.(*Table); ok {
		return tbl.String()
	}
	return "(" + q.String() + ")"
}

func bestKey(q Query, mode Mode) []string {
	var best []string
	bestCost := impossible
	for _, key := range q.Keys() {
		cost := Optimize(q, mode, key, assess)
		cost += (len(key) - 1) * cost / 20 // ??? prefer shorter keys
		if cost < bestCost {
			best = key
		}
	}
	return best
}
