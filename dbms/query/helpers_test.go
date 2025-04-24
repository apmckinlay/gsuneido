// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestStripSort(t *testing.T) {
	test := func(query, expected string) {
		t.Helper()
		q := StripSort(query)
		assert.T(t).This(q).Is(expected)
	}
	test("table", "table")
	test("table sort a", "table")
	test("table sort a, b, c", "table")
	test("table sort reverse a", "table")
	test("table sort reverse a, b, c", "table")
	test("table sort a where x > 5", "table where x > 5")
	test("customer join trans where date > '2020-01-01' sort id",
		"customer join trans where date > '2020-01-01'")
	test("(table where x > 5) sort a", "(table where x > 5)")
	test("customer extend total = sum(amount) sort total",
		"customer extend total = sum(amount)")
	test("table project sort sort sort",
		"table project sort")
	test("table project sort",
		"table project sort")
}

func TestGetSort(t *testing.T) {
	test := func(query, expected string) {
		t.Helper()
		sort := GetSort(query)
		assert.T(t).This(sort).Is(expected)
	}
	test("table", "")
	test("table sort a", "a")
	test("table sort a, b, /* comment */ c", "a,b,c")
	test("table sort reverse a", "reverse a")
	test("table sort reverse a, b, c", "reverse a,b,c")
	test("table sort a where x > 5", "a")
	test("customer join trans where date > '2020-01-01' sort id", "id")
	test("(table where x > 5) sort a", "a")
	test("customer extend total = sum(amount) sort total", "total")
	test("table project sort sort sort", "sort")
	test("table project sort", "")
}

func TestJustTable(t *testing.T) {
	test := func(query string, expected string) {
		t.Helper()
		assert.T(t).This(JustTable(query)).Is(expected)
	}

	// Valid table names
	test("customer", "customer")
	test("  customer  ", "customer")
	test("/* comment */ customer // comment", "customer")
	test("Customer", "Customer")

	// Invalid - not just a table
	test("customer join supplier", "")
	test("customer where city = 'Calgary'", "")
	test("", "")
	test("customer,supplier", "")
	test("customer sort name", "")
	test("'string'", "")
	test("123", "")
	test("customer.field", "")
}
