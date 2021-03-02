// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/runtime"
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

	SetTran(tran QueryTran)

	// Header() runtime.Header

	// Order() []string

	Fixed() []Fixed

	// Updateable() bool

	// Indexes returns all the indexes
	Indexes() [][]string

	// Keys returns the indexes that are keys
	Keys() [][]string

	// Lookup returns the row matching the given key, or nil if not found
	// Lookup(index []string, key string) runtime.Row

	// Rewind()

	// Get(dir runtime.Dir) runtime.Row

	String() string
}

type QueryTran interface {
	GetSchema(table string) *schema.Schema
}

// Setup prepares a parsed query for execution
func Setup(q Query, isCursor bool, t runtime.ITran) Query {
	return q //TODO
}

type Cost float64

func Optimize(index, needs []string, isCursor, freeze bool) Cost {
	return 0 //TODO
}

type Query1 struct {
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
