// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/runtime"
)

type Query interface {
	Init()
	Transform() Query

	// Header() runtime.Header
	// Order() []string
	// Fixed() []Fixed
	// Updateable() bool
	// Columns() []string
	// Keys() [][]string
	// Indexes() [][]string

	// Lookup returns the row matching the given key, or nil if not found
	// Lookup(index []string, key string) runtime.Row

	// Rewind()

	// Get(dir runtime.Dir) runtime.Row

	String() string
}

// Setup prepares a parsed query for execution
func Setup(q Query, isCursor bool, t runtime.ITran) Query {
	return q //TODO
}

type Cost float64

func Optimize(index, needs []string, isCursor, freeze bool) Cost {
	return 0 //TODO
}

type Fixed struct {
	//TODO
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

func (q1 *Query1) Transform() Query { //TODO remove - temporary
	// q1.source = q1.source.Transform()
	panic("not implemented")
}

type Query2 struct {
	Query1
	source2 Query
}

func (q2 *Query2) String(op string) string {
	return q2.Query1.String() + " " + op + " " + q2.source2.String()
}

func (q2 *Query2) Init() {
	q2.source.Init()
	q2.source2.Init()
}

func (q2 *Query2) Transform() Query { //TODO remove - temporary
	// 	q2.source = q2.source.Transform()
	// 	q2.source2 = q2.source2.Transform()
	panic("not implemented")
}
