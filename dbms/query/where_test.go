// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestExtractCompares(t *testing.T) {
	test := func(query string, expected string) *Where {
		t.Helper()
		q := ParseQuery("table where " + query)
		q.SetTran(testTran{})
		q.Init()
		w := q.(*Where)
		before := w.expr.String()
		cmps := w.extractCompares()
		if expected == "" {
			after := w.expr.String()
			assert.T(t).This(after).Is(before)
		}
		assert.T(t).This(fmt.Sprint(cmps)).Is(expected)
		return w
	}
	test("Foo()", "[]")
	test("a isnt 1", "[]")
	test("a is 1", "[a Is 1]")
	test("a > 2", "[a Gt 2]")
	test("a or b", "[]")
	test("a is 1 or a is 2", "[a In (1, 2)]")
	test("a is 1 and b is 2", "[a Is 1 b Is 2]")
	test("a in (1, 2, 3)", "[a In (1, 2, 3)]")
	w := test("1 is 2", "[]")
	assert.T(t).That(w.conflict)
}

func TestComparesToColSelects(t *testing.T) {
	test := func(query string, expected string) {
		t.Helper()
		q := ParseQuery("table where " + query)
		q.SetTran(testTran{})
		q.Init()
		w := q.(*Where)
		cmps := w.extractCompares()
		colSels := w.comparesToFilters(cmps)
		assert.T(t).This(fmt.Sprint(colSels)).Is("map[" + expected + "]")
	}
	test("a > 2", "a:(2..<max>)")
	test("a is 1 and b is 2", "a:[1] b:[2]")
	test("a in (1,2,3)", "a:[1,2,3]")
	test("a in (1,2,3,4) and a in (3,4,5,6)", "a:[3,4]")
	test("a in (1,2,3,4) and a > 2", "a:[3,4]")
	test("a in (1,2,3,4) and a < 3", "a:[1,2]")
	test("a > 1 and a < 3", "a:(1..3)")
	test("a > 1 and a > 2", "a:(2..<max>)")
	test("a > 5 and a < 3", "") // conflict
}

func TestColSelsToIdxSel(t *testing.T) {
	idx := []string{"a", "b", "c"}
	test := func(query string, expected string) {
		t.Helper()
		q := ParseQuery("comp where " + query)
		q.SetTran(testTran{})
		q.Init()
		w := q.(*Where)
		cmps := w.extractCompares()
		colSels := w.comparesToFilters(cmps)
		idxSels, _ := colSelsToIdxSel(colSels, idx)
		assert.T(t).This(fmt.Sprint(idxSels)).Is("[" + expected + "]")
	}
	test("a is 1", "[1]")
	test("a is 1 and c is 2", "[1]")
	test("a is 1 and b is 2", "[1] [2]")
	test("a is 1 and b is 2 and c is 3", "[1] [2] [3]")
	test("a > 4", "(4..<max>)")
	test("a > 4 and b is 2", "(4..<max>)")
	test("a is 2 and b > 4", "[2] (4..<max>)")
	test("a in (1,2) and b in (3,4)", "[1,2] [3,4]")
}

func TestExplodeFilters(t *testing.T) {
	idx := []string{"a", "b", "c"}
	test := func(query string, expected string) {
		t.Helper()
		q := ParseQuery("comp where " + query)
		q.SetTran(testTran{})
		q.Init()
		w := q.(*Where)
		cmps := w.extractCompares()
		colSels := w.comparesToFilters(cmps)
		idxSels, _ := colSelsToIdxSel(colSels, idx)
		filters := explodeFilters(idxSels, [][]filter{nil})
		assert.T(t).This(fmt.Sprint(filters)).Is("[" + expected + "]")
	}
	test("a is 1", "[[1]]")
	test("a is 1 and b is 2", "[[1] [2]]")
	test("a is 1 and b is 2 and c is 3", "[[1] [2] [3]]")
	test("a > 4", "[(4..<max>)]")
	test("a is 2 and b > 4", "[[2] (4..<max>)]")
	test("a in (1,2) and b in (3,4)", "[[1] [3]] [[1] [4]] [[2] [3]] [[2] [4]]")
	test("a in (1,2) and b > 4", "[[1] (4..<max>)] [[2] (4..<max>)]")
}

func TestCompositeFilters(t *testing.T) {
	test := func(query string, expected string) {
		t.Helper()
		q := ParseQuery("comp where " + query)
		q.SetTran(testTran{})
		q.Init()
		w := q.(*Where)
		cmps := w.extractCompares()
		colSels := w.comparesToFilters(cmps)
		idxSels := w.colSelsToIdxSels(colSels)
		assert.T(t).This(fmt.Sprint(idxSels)).Is("[" + expected + "]")
	}
	test("a is 1", "a,b,c: 1")
	test("a is 1 and b is 2", "a,b,c: 1,2")
	test("a is 1 and b is 2 and c is 3", "a,b,c: 1,2,3")
	test("a > 4", "a,b,c: 4..<max>")
	test("a is 2 and b > 4", "a,b,c: 2,4..2,<max>")
	test("a in (1,2) and b in (3,4)", "a,b,c: 1,3 | 1,4 | 2,3 | 2,4")
	test("a in (1,2) and b > 4", "a,b,c: 1,4..1,<max> | 2,4..2,<max>")
}
