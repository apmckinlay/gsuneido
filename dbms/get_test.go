// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	qry "github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"
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

	test := func(sels Sels, expected []string) {
		t.Helper()
		result := findKey(keys, sels)
		assert.T(t).This(result).Is(expected)
	}

	// Exact match
	test(selCols("id"), []string{"id"})

	// Superset match (sels must contain all of key)
	test(selCols("name", "date", "extra"), []string{"name", "date"})
	test(selCols("name", "date"), []string{"name", "date"})

	// No match (sels doesn't contain all of key)
	test(selCols("name"), nil)
	test(selCols("date"), nil)
	test(selCols("xyz"), nil)
	test(Sels{}, nil)

	// Multiple keys possible (should return first match)
	test(selCols("a", "b", "c"), []string{"a", "b", "c"})
}

func selCols(cols ...string) Sels {
	sels := make(Sels, len(cols))
	for i, col := range cols {
		sels[i] = qry.NewSel(col, "")
	}
	return sels
}

func TestFindAll(t *testing.T) {
	indexes := [][]string{
		{"a", "b", "c"},
		{"x", "y"},
		{"p", "q", "r"},
	}

	test := func(sels Sels, expected []string) {
		t.Helper()
		result := findAll(indexes, sels)
		assert.T(t).This(result).Is(expected)
	}

	// Exact match on first index
	test(selCols("a", "b", "c"), []string{"a", "b", "c"})

	// Prefix match
	test(selCols("a", "b"), []string{"a", "b", "c"})
	test(selCols("a"), []string{"a", "b", "c"})

	// Match on second index
	test(selCols("x", "y"), []string{"x", "y"})
	test(selCols("x"), []string{"x", "y"})

	// No match
	test(selCols("z"), nil)
	test(selCols("a", "x"), nil)
	test(selCols("b", "c"), nil) // must start from beginning
}

func TestHasPrefix(t *testing.T) {
	test := func(idx []string, sels Sels, expected bool) {
		t.Helper()
		result := hasPrefix(idx, sels)
		assert.T(t).This(result).Is(expected)
	}

	test([]string{"a", "b"}, selCols("a"), true)
	test([]string{"a", "b"}, selCols("b"), false)
	test([]string{"a", "b"}, selCols("a", "b"), true)
	test([]string{"a", "b"}, selCols("x"), false)
	test([]string{"a"}, selCols("a"), true)
}

func TestUsableIndexes(t *testing.T) {
	indexes := [][]string{
		{"a", "b"},
		{"b", "c"},
		{"c", "d"},
		{"a", "x"},
	}

	test := func(sels Sels, expected [][]string) {
		t.Helper()
		result := usableIndexes(indexes, sels)
		assert.T(t).This(result).Is(expected)
	}

	// Fields match first elements (any index whose first field is in sels)
	test(selCols("a"), [][]string{{"a", "b"}, {"a", "x"}})
	test(selCols("b"), [][]string{{"b", "c"}})
	test(selCols("c"), [][]string{{"c", "d"}})

	// Multiple fields - includes indexes starting with "a" or "b"
	test(selCols("a", "b"), [][]string{{"a", "b"}, {"b", "c"}, {"a", "x"}})

	// No matches
	test(selCols("x"), nil)
	test(selCols("d"), nil)
}

func TestFormatFieldsVals(t *testing.T) {
	test := func(sels Sels, expected string) {
		t.Helper()
		result := formatFieldsVals(sels)
		assert.T(t).This(result).Is(expected)
	}

	// Single field
	test(Sels{qry.NewSel("name", Pack(SuStr("John")))}, `name: "John"`)

	// Multiple fields
	test(Sels{qry.NewSel("name", Pack(SuStr("John"))), qry.NewSel("age", Pack(SuInt(30)))}, `name: "John", age: 30`)

	// Empty
	test(Sels{}, "")
}
