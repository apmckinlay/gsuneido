// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/slc"
)

func TestGetQuery(t *testing.T) {
	test := func(ob *SuObject, expected string) {
		t.Helper()
		assert.T(t).This(getQuery(ob)).Is(expected)
	}

	// From list position 0
	ob := &SuObject{}
	ob.Add(SuStr("mytable"))
	test(ob, "mytable")

	ob2 := &SuObject{}
	ob2.Add(SuStr("customer join orders"))
	test(ob2, "customer join orders")

	// From named argument
	obj := &SuObject{}
	obj.Set(SuStr("query"), SuStr("table"))
	test(obj, "table")

	// Empty object
	test(&SuObject{}, "")

	// List takes precedence
	obj2 := &SuObject{}
	obj2.Add(SuStr("list_query"))
	obj2.Set(SuStr("query"), SuStr("named_query"))
	test(obj2, "list_query")
}

func TestGetWhere(t *testing.T) {
	// Empty object
	assert.T(t).This(getWhere(&SuObject{})).Is("")

	// Single field
	obj := &SuObject{}
	obj.Set(SuStr("name"), SuStr("John"))
	assert.T(t).This(getWhere(obj)).Is(`where name is "John"`)

	// With query argument (should be excluded)
	obj3 := &SuObject{}
	obj3.Set(SuStr("query"), SuStr("customer"))
	obj3.Set(SuStr("city"), SuStr("Calgary"))
	assert.T(t).This(getWhere(obj3)).Is(`where city is "Calgary"`)

	// Multiple fields - check structure
	obj4 := &SuObject{}
	obj4.Set(SuStr("a"), SuInt(1))
	obj4.Set(SuStr("b"), SuInt(2))
	w := getWhere(obj4)
	assert.T(t).True(strings.HasPrefix(w, "where "))
	assert.T(t).True(strings.Contains(w, "\nand "))
}

func TestFindKey(t *testing.T) {
	keys := [][]string{
		{"id"},
		{"name", "date"},
		{"a", "b", "c"},
	}

	test := func(flds []string, expected []string) {
		t.Helper()
		result := findKey(keys, flds)
		assert.T(t).This(result).Is(expected)
	}

	// Exact match
	test([]string{"id"}, []string{"id"})

	// Superset match (flds must contain all of key)
	test([]string{"name", "date", "extra"}, []string{"name", "date"})
	test([]string{"name", "date"}, []string{"name", "date"})

	// No match (flds doesn't contain all of key)
	test([]string{"name"}, nil)
	test([]string{"date"}, nil)
	test([]string{"xyz"}, nil)
	test([]string{}, nil)

	// Multiple keys possible (should return first match)
	test([]string{"a", "b", "c"}, []string{"a", "b", "c"})
}

func TestFindAll(t *testing.T) {
	indexes := [][]string{
		{"a", "b", "c"},
		{"x", "y"},
		{"p", "q", "r"},
	}

	test := func(flds []string, expected []string) {
		t.Helper()
		result := findAll(indexes, flds)
		assert.T(t).This(result).Is(expected)
	}

	// Exact match on first index
	test([]string{"a", "b", "c"}, []string{"a", "b", "c"})

	// Prefix match
	test([]string{"a", "b"}, []string{"a", "b", "c"})
	test([]string{"a"}, []string{"a", "b", "c"})

	// Match on second index
	test([]string{"x", "y"}, []string{"x", "y"})
	test([]string{"x"}, []string{"x", "y"})

	// No match
	test([]string{"z"}, nil)
	test([]string{"a", "x"}, nil)
	test([]string{"b", "c"}, nil) // must start from beginning
}

func TestHasPrefix(t *testing.T) {
	test := func(idx, flds []string, expected bool) {
		t.Helper()
		result := hasPrefix(idx, flds)
		assert.T(t).This(result).Is(expected)
	}

	test([]string{"a", "b"}, []string{"a"}, true)
	test([]string{"a", "b"}, []string{"b"}, false)
	test([]string{"a", "b"}, []string{"a", "b"}, true)
	test([]string{"a", "b"}, []string{"x"}, false)
	test([]string{"a"}, []string{"a"}, true)
}

func TestUsableIndexes(t *testing.T) {
	indexes := [][]string{
		{"a", "b"},
		{"b", "c"},
		{"c", "d"},
		{"a", "x"},
	}

	test := func(flds []string, expected [][]string) {
		t.Helper()
		result := usableIndexes(indexes, flds)
		assert.T(t).This(result).Is(expected)
	}

	// Fields match first elements (any index whose first field is in flds)
	test([]string{"a"}, [][]string{{"a", "b"}, {"a", "x"}})
	test([]string{"b"}, [][]string{{"b", "c"}})
	test([]string{"c"}, [][]string{{"c", "d"}})

	// Multiple fields - includes indexes starting with "a" or "b"
	test([]string{"a", "b"}, [][]string{{"a", "b"}, {"b", "c"}, {"a", "x"}})

	// No matches
	test([]string{"x"}, nil)
	test([]string{"d"}, nil)
}

func TestFormatFieldsVals(t *testing.T) {
	test := func(flds []string, vals []Value, expected string) {
		t.Helper()
		packed := slc.MapFn(vals, func(v Value) string { return Pack(v.(Packable)) })
		result := formatFieldsVals(flds, packed)
		assert.T(t).This(result).Is(expected)
	}

	// Single field
	test([]string{"name"}, []Value{SuStr("John")}, `name: "John"`)

	// Multiple fields
	test([]string{"name", "age"}, []Value{SuStr("John"), SuInt(30)}, `name: "John", age: 30`)

	// Empty
	test([]string{}, []Value{}, "")
}
