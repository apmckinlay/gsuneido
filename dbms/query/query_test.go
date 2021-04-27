// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestKeys(t *testing.T) {
	test := func(query, expected string) {
		t.Helper()
		q := ParseQuery(query)
		q.SetTran(testTran{})
		q.Init()
		// q = q.Transform()
		assert.T(t).This(idxsToS(q.Keys())).Is(expected)
	}
	test("tables", "table")
	test("columns", "table+column")
	test("columns rename column to col", "table+col")
	test("tables extend x=1,b=2", "table")
	test("tables intersect tables", "table")
	test("abc intersect bcd", "b, c")
	test("hist project item, cost", "item+cost")
	test("hist2 project date, item", "date")
	test("abc times inven", "b+item, c+item")
	test("abc union abc", "a+b+c") // not disjoint
	test("(bcd where b is 1) union (bcd where b is 2)", "b")
}

func TestByContainsKey(t *testing.T) {
	test := func(by string, keys string, expected bool) {
		t.Helper()
		result := containsKey(strings.Fields(by), sToIdxs(keys))
		assert.T(t).This(result).Is(expected)
	}
	test("", "a b, c", false)
	test("a b", "a", true)
	test("a b", "b", true)
	test("a b", "x, a+b, y", true)
}

func TestProjectIndexes(t *testing.T) {
	test := func(idxs string, cols string, expected string) {
		t.Helper()
		result := projectIndexes(sToIdxs(idxs), strings.Fields(cols))
		assert.T(t).This(idxsToS(result)).Is(expected)
	}
	test("a, b+c, d+e+f", "a b c d", "a, b+c")
	test("a, b+c, d+e+f", "c e f", "")
}

// sToIdxs splits strings like: "a+b, c, d+e+f"
func sToIdxs(s string) [][]string {
	var idxs [][]string
	for _, ix := range strings.Split(s, ", ") {
		idxs = append(idxs, strings.Split(ix, "+"))
	}
	return idxs
}

// idxsToS converts [][]string to a string like: "a+b, c, d+e+f"
func idxsToS(idxs [][]string) string {
	tmp := make([]string, len(idxs))
	for i, ix := range idxs {
		tmp[i] = strings.Join(ix, "+")
	}
	return strings.Join(tmp, ", ")
}
