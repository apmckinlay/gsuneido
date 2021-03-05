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

	// Lookup returns the row matching the given key, or nil if not found
	// Lookup(index []string, key string) runtime.Row

	// Rewind()

	// Get(dir runtime.Dir) runtime.Row

	String() string

	cacheAdd(index []string, cost Cost)
	cacheGet(index []string) Cost

	optimize(index []string, readonly, freeze bool) Cost

	setTempIndex(index []string)
	getTempIndex() []string
}

type QueryTran interface {
	GetSchema(table string) *schema.Schema
	GetInfo(table string) *meta.Info
}

// Setup prepares a parsed query for execution
func Setup(q Query, readonly bool, t QueryTran) Query {
	q.SetTran(t)
	q.Init()
	q = q.Transform()
	cost := q.optimize(nil, readonly, true)
	if cost >= impossible {
		panic("invalid query " + q.String())
	}
	// q = q.addTempIndex(t) //TODO
	return q
}

type Cost = int

const impossible Cost = ints.MaxInt / 64 // allow for adding IMPOSSIBLE's

// Optimize is the start of optimization.
// It handles whether or not to use a temp index.
func Optimize(q Query, index []string, readonly, freeze bool) Cost {
	if !sset.Subset(q.Columns(), index) {
		return impossible
	}
	if !readonly || !q.Updateable() {
		return optimize1(q, index, readonly, freeze)
	}
	cost1 := optimize1(q, index, readonly, false)
	noIndexCost := optimize1(q, nil, readonly, false)
	tempIndexCost := 0 //TODO
	cost2 := noIndexCost + Cost(tempIndexCost)
	assert.That(cost2 >= 0)

	cost := ints.Min(cost1, cost2)
	if cost >= impossible {
		return impossible
	}
	if freeze {
		if cost2 < cost2 {
			q.setTempIndex(index)
			index = nil
		}
		optimize1(q, index, readonly, true)
	}
	return cost
}

// optimize1 handles caching
func optimize1(q Query, index []string, readonly, freeze bool) Cost {
	var cost Cost
	if !freeze {
		cost = q.cacheGet(index)
		if cost >= 0 {
			return cost
		}
	}
	cost = q.optimize(index, readonly, freeze)
	assert.That(cost >= 0)
	if !freeze {
		q.cacheAdd(index, cost)
	}
	return cost
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

func (q1 *Query1) Fixed() []Fixed {
	return q1.source.Fixed()
}

func (q1 *Query1) Updateable() bool {
	return q1.source.Updateable()
}

func (q1 *Query1) SetTran(t QueryTran) {
	q1.source.SetTran(t)
}

func (q1 *Query1) optimize(index []string, readonly, freeze bool) Cost {
	return 0
} //TODO remove (temporary)

// Query2 -----------------------------------------------------------

type Query2 struct {
	Query1
	disallow
	source2 Query
}

type disallow struct{}

func (disallow) Keys() [][]string {
	return nil
}

func (disallow) Indexes() [][]string {
	return nil
}

func (q2 *Query2) String(op string) string {
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

func paren(q Query) string {
	if tbl, ok := q.(*Table); ok {
		return tbl.String()
	}
	return "(" + q.String() + ")"
}
