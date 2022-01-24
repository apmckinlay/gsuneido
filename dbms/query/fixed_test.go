// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestFixed(t *testing.T) {
	runtime.DefaultSingleQuotes = true
	defer func() { runtime.DefaultSingleQuotes = false }()
	test := func(query, expected string) {
		t.Helper()
		q := ParseQuery(query, testTran{})
		assert.T(t).This(fixedStr(q.Fixed())).Is(expected)
	}
	test("table", "[]")

	test("table extend f=1", "[f=(1)]")
	test("table extend f=1, g='s'", "[f=(1), g=('s')]")
	test("table extend f=1 extend g=2", "[f=(1), g=(2)]")
	test("table extend f=1 where f is 1", "[f=(1)]")
	test("table extend f=1, g=2 where g in (1,2,3) and a=5",
		"[f=(1), g=(2), a=(5)]")

	test("table where a is 1", "[a=(1)]")
	test("table where a is 1 and b is 's' and a is b", "[a=(1), b=('s')]")

	test("table union (table extend f=1)", "[f=(1,'')]")
	test("(table extend f=2) union (table extend f=1)", "[f=(2,1)]")

	test("(table extend f=1, g=2) project a,b", "[]")
	test("(table extend f=1, g=2) project a,g", "[g=(2)]")
	test("(table extend f=1, g=2) project a,f,g", "[f=(1), g=(2)]")

	test("(table extend f=1) join (table extend f=1, g=2)", "[f=(1), g=(2)]")

	test("table extend f=1, g=2 rename g to h", "[f=(1), h=(2)]")

	test2 := func(query string, expected string) {
		t.Helper()
		q := ParseQuery(query, testTran{})
		q, _ = Setup(q, ReadMode, testTran{})
		assert.T(t).This(q.String()).Is(expected)
	}
	test2("table extend f=1 where f is 2",
		"NOTHING")
	test2("table extend f=1, g=2 where f is 1",
		"table^(a) WHERE true EXTEND f = 1, g = 2") //TODO remove WHERE true
	test2("table extend f=1, g=2 where f is 3",
		"NOTHING")
	test2("tables extend x=1 join (columns extend x=2)",
		"NOTHING")
	test2("tables extend x=1 leftjoin (columns extend x=2)",
		"tables EXTEND x = 1")
	test2("table where a = 1 and a = 2",
		"table WHERE nothing") //TODO should be just: table
	test2("(table union table2) where a = 1",
		"table^(a) WHERE*1 a is 1")
	test2("(table minus table2) where a = 1",
		"table^(a) WHERE*1 a is 1")
}

func TestCombineFixed(t *testing.T) {
	f := func(f string) []Fixed {
		result := []Fixed{}
		if f == "" {
			return result
		}
		for _, g := range strings.Split(f, "; ") {
			h := strings.Split(g, " ")
			result = append(result, Fixed{col: h[0], values: h[1:]})
		}
		return result
	}
	test := func(src, oth string, expected string) {
		result, none := combineFixed(f(src), f(oth))
		if expected == "none" {
			assert.T(t).That(none == true)
		} else {
			assert.T(t).This(result).Is(f(expected))
		}
	}
	test("a 1; b 2 3; c 4", "", "a 1; b 2 3; c 4")
	test("", "a 1; b 2 3; c 4", "a 1; b 2 3; c 4")
	test("a 1; c 1 2", "c 2 3; d 4", "a 1; c 2; d 4")
	test("a 1; b 8", "z 9; a 2", "none")
	test("a 1 2 3", "a 4 5 6", "none")
}
