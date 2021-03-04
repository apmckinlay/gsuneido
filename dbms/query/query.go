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
*/
package query

import (
	"math"

	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/util/assert"
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
}

type QueryTran interface {
	GetSchema(table string) *schema.Schema
}

// Setup prepares a parsed query for execution
func Setup(q Query, readonly bool, t QueryTran) Query {
	return q //TODO
}

type Cost float64

const IMPOSSIBLE = math.MaxFloat64 / 100 // allow for adding IMPOSSIBLE's

func Optimize(q Query, index []string, readonly, freeze bool) Cost {
	if !sset.Subset(q.Columns(), index) {
		return IMPOSSIBLE
	}
	return 0 //TODO
}

type Query1 struct {
	cache
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
