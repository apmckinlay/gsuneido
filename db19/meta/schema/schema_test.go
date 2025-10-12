// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package schema

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestCheckDuplicateColumns(t *testing.T) {
	assert := assert.T(t)

	// Test duplicate in Columns
	assert.This(func() {
		s := Schema{
			Table:   "test",
			Columns: []string{"a", "b", "a"},
			Indexes: []Index{{Mode: 'k', Columns: []string{"a"}}},
		}
		s.Check()
	}).Panics("duplicate column")

	// Test duplicate in Derived
	assert.This(func() {
		s := Schema{
			Table:   "test",
			Columns: []string{"a", "b"},
			Derived: []string{"C", "D", "C"},
			Indexes: []Index{{Mode: 'k', Columns: []string{"a"}}},
		}
		s.Check()
	}).Panics("duplicate derived")

	// Test that deleted columns ("-") can repeat
	s := Schema{
		Table:   "test",
		Columns: []string{"a", "-", "-", "b"},
		Indexes: []Index{{Mode: 'k', Columns: []string{"a"}}},
	}
	s.Check() // Should not panic

	// Test valid schema with no duplicates
	s = Schema{
		Table:   "test",
		Columns: []string{"a", "b", "c"},
		Derived: []string{"D", "E"},
		Indexes: []Index{{Mode: 'k', Columns: []string{"a"}}},
	}
	s.Check() // Should not panic
}

func TestCheckIndexes(t *testing.T) {
	assert := assert.T(t)

	// Test invalid index column
	assert.This(func() {
		CheckIndexes("test", []string{"a", "b"},
			[]Index{{Mode: 'k', Columns: []string{"c"}}})
	}).Panics("invalid index column: c")

	// Test duplicate indexes
	assert.This(func() {
		CheckIndexes("test", []string{"a", "b"},
			[]Index{
				{Mode: 'k', Columns: []string{"a"}},
				{Mode: 'i', Columns: []string{"a"}},
			})
	}).Panics("duplicate index")

	// Test empty index columns for non-key
	assert.This(func() {
		CheckIndexes("test", []string{"a", "b"},
			[]Index{{Mode: 'i', Columns: []string{}}})
	}).Panics("index columns must not be empty")

	// Test valid indexes
	CheckIndexes("test", []string{"a", "b"},
		[]Index{
			{Mode: 'k', Columns: []string{"a"}},
			{Mode: 'i', Columns: []string{"b"}},
		}) // Should not panic
}

func TestCheckForKey(t *testing.T) {
	assert := assert.T(t)

	// Test schema without a key
	assert.This(func() {
		s := Schema{
			Table:   "test",
			Columns: []string{"a", "b"},
			Indexes: []Index{{Mode: 'i', Columns: []string{"a"}}},
		}
		s.Check()
	}).Panics("key required")

	// Test schema with a key
	s := Schema{
		Table:   "test",
		Columns: []string{"a", "b"},
		Indexes: []Index{{Mode: 'k', Columns: []string{"a"}}},
	}
	s.Check() // Should not panic
}

func TestCheckLower(t *testing.T) {
	assert := assert.T(t)

	// Test _lower! with nonexistent column
	assert.This(func() {
		s := Schema{
			Table:   "test",
			Columns: []string{"a", "b"},
			Derived: []string{"name_lower!"},
			Indexes: []Index{{Mode: 'k', Columns: []string{"a"}}},
		}
		s.Check()
	}).Panics("_lower! nonexistent column: name")

	// Test valid _lower! reference
	s := Schema{
		Table:   "test",
		Columns: []string{"a", "name"},
		Derived: []string{"name_lower!"},
		Indexes: []Index{{Mode: 'k', Columns: []string{"a"}}},
	}
	s.Check() // Should not panic
}
