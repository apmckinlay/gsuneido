// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"math"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestExtractCompares(t *testing.T) {
	test := func(query string, expected string) *Where {
		t.Helper()
		w := ParseQuery("table where "+query, testTran{}, nil).(*Where)
		w.optInit()
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
}

func TestComparesToColSelects(t *testing.T) {
	test := func(query string, expected string) {
		t.Helper()
		w := ParseQuery("table where "+query, testTran{}, nil).(*Where)
		w.optInit()
		cmps := w.extractCompares()
		colSels := w.comparesToFilters(cmps)
		assert.T(t).This(fmt.Sprint(colSels)).Is("map[" + expected + "]")
	}
	test("a >= 2", "a:(2..<max>)")
	test("a > 2", "a:(2+..<max>)")
	test("a is 1 and b is 2", "a:1 b:2")
	test("a in (1,2,3)", "a:[1,2,3]")
	test("a in (1,2,3,4) and a in (3,4,5,6)", "a:[3,4]")
	test("a in (1,2,3,4) and a > 2", "a:[3,4]")
	test("a in (1,2,3,4) and a < 3", "a:[1,2]")
	test("a >= 1 and a < 3", "a:(1..3)")
	test("a > 1 and a >= 2", "a:(2..<max>)")
	test("a > 5 and a < 3", "") // conflict
}

func TestColSelsToIdxFilter(t *testing.T) {
	idx := []string{"a", "b", "c"}
	test := func(query string, expected string) {
		t.Helper()
		w := ParseQuery(query, testTran{}, nil).(*Where)
		w.optInit()
		cmps := w.extractCompares()
		colSels := w.comparesToFilters(cmps)
		filters := colSelsToIdxFilters(colSels, idx)
		assert.T(t).This(fmt.Sprint(filters)).Is("[" + expected + "]")
	}
	test("comp where a is 1", "1")
	test("comp where a is 1 and c is 2", "1")
	test("comp where a is 1 and b is 2", "1 2")
	test("comp where a is 1 and b is 2 and c is 3", "1 2 3")
	test("comp where a >= 4", "(4..<max>)")
	test("comp where a >= 4 and b is 2", "(4..<max>)")
	test("comp where a is 2 and b >= 4", "2 (4..<max>)")
	test("comp where a in (1,2) and b in (3,4)", "[1,2] [3,4]")
	idx = []string{"id"}
	test("customer where id is 'e'", "'e'")
}

func TestExplodeFilters(t *testing.T) {
	idx := []string{"a", "b", "c"}
	test := func(query string, expected string) {
		t.Helper()
		w := ParseQuery("comp where "+query, testTran{}, nil).(*Where)
		w.optInit()
		cmps := w.extractCompares()
		colSels := w.comparesToFilters(cmps)
		filters := colSelsToIdxFilters(colSels, idx)
		exploded := explodeFilters(filters, [][]filter{nil})
		assert.T(t).This(fmt.Sprint(exploded)).Is("[" + expected + "]")
	}
	test("a is 1", "[1]")
	test("a is 1 and b is 2", "[1 2]")
	test("a is 1 and b is 2 and c is 3", "[1 2 3]")
	test("a >= 4", "[(4..<max>)]")
	test("a is 2 and b >= 4", "[2 (4..<max>)]")
	test("a in (1,2) and b in (3,4)", "[1 3] [1 4] [2 3] [2 4]")
	test("a in (1,2) and b >= 4", "[1 (4..<max>)] [2 (4..<max>)]")
}

func TestColSelsToIdxSels(t *testing.T) {
	test := func(query string, expected string) {
		t.Helper()
		w := ParseQuery("comp where "+query, testTran{}, nil).(*Where)
		w.optInit()
		cmps := w.extractCompares()
		colSels := w.comparesToFilters(cmps)
		idxSels := w.colSelsToIdxSels(colSels)
		assert.T(t).This(fmt.Sprint(idxSels)).Is("[" + expected + "]")
	}
	test("a is 1", "a,b,c: 1..1,<max> = 0.1")
	test("a is 1 and b is 2", "a,b,c: 1,2..1,2,<max> = 0.01")
	test("a is 1 and b is 2 and c is 3", "a,b,c: 1,2,3 = 0.001")
	test("a > 4", "a,b,c: 4,<max>..<max> = 0.5")
	test("a <= 4", "a,b,c: ..4,<max> = 0.5")
	test("a is 2 and b >= 4", "a,b,c: 2,4..2,<max> = 0.06")
	test("a in (1,2) and b in (3,4)", "a,b,c: 1,3..1,3,<max> | 1,4..1,4,<max> | "+
		"2,3..2,3,<max> | 2,4..2,4,<max> = 0.04")
	test("a in (1,2) and b > 4",
		"a,b,c: 1,4,<max>..1,<max> | 2,4,<max>..2,<max> = 0.1")
}

func TestFracPos(t *testing.T) {
	tt := testTran{}
	test := func(expected float64, digits ...int) {
		t.Helper()
		var enc ixkey.Encoder
		for _, d := range digits {
			enc.Add(runtime.Pack(runtime.SuInt(d)))
		}
		key := enc.String()
		f := tt.fracPos(key, true)
		assert.T(t).That(math.Abs(f-expected) < .0001)
	}
	test(0)
	test(.5, 5)
	test(.234, 2, 3, 4)
}

func TestWhereNrows(t *testing.T) {
	test := func(query string, nrows int) {
		t.Helper()
		var tran testTran
		w := ParseQuery(query, tran, nil)
		Setup(w, ReadMode, tran)
		n, p := w.Nrows()
		assert.T(t).This(n).Is(nrows)
		assert.T(t).This(p).Is(100)
	}
	test("table where F()", 50)
	test("inven where item >= 5", 50)
	test("inven where item < 3 and item > 3", 0) // conflict
	test("inven where item is 1", 1)
	test("inven where item in (1,2,3)", 3)
	test("inven where item > 2 and item < 4", 20)
	test("inven where item > 2 and item < 4 and qty", 10)
	test("hist where date is 3", 5) // half of 1/10 (not 1 since not key)
	test("inven extend x where x > 5", 50) // not on table
}
