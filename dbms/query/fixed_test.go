// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestFixed(t *testing.T) {
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	DefaultSingleQuotes = true
	defer func() { DefaultSingleQuotes = false }()
	test := func(query, expected string) {
		t.Helper()
		q := ParseQuery(query, testTran{}, nil)
		assert.T(t).Msg(q).This(fixedStr(q.Fixed())).Is(expected)
	}
	test("table", "[]")

	test("table extend f=1", "[f=(1)]")
	test("table extend f=1, g='s'", "[f=(1), g=('s')]")
	test("table extend f=1 extend g=2", "[f=(1), g=(2)]")

	test("table extend f=1 where f is 1", "[f=(1)]")
	test("table extend f=1, g=2 where a=5",
		"[a=(5), f=(1), g=(2)]")
	test("table extend f=1, g=2 where g in (1,2,3) and a=5",
		"[a=(5), f=(1), g=(2)]")
	test("table extend x=a where x is 5",
		"[x=(5)]")
	test("table where a is 5 extend x=a",
		"[a=(5), x=(5)]")
	test("table extend x=a, y=x where y is 5",
		"[y=(5)]")
	test("table where a is 5 extend x=a, y=x",
		"[a=(5), x=(5), y=(5)]")

	test("table where a is 1", "[a=(1)]")
	test("table where a is ''", "[a=('')]")
	test("table where a <= ''", "[a=('')]")
	test("table where a <= 'x'", "[]")
	test("table where a is 1 and b is 's' and a is b", "[a=(1), b=('s')]")

	test("table union (table extend f=1)", "[f=(1,'')]")
	test("(table extend f=2) union (table extend f=1)", "[f=(2,1)]")

	test("(table extend f=1, g=2) project a,b", "[]")
	test("(table extend f=1, g=2) project a,g", "[g=(2)]")
	test("(table extend f=1, g=2) project a,f,g", "[f=(1), g=(2)]")

	test("(table extend f=1) join (table extend f=1, g=2)", "[f=(1), g=(2)]")

	test("table extend f=1, g=2 rename g to h", "[f=(1), h=(2)]")

	test("table extend x=1, y=2 leftjoin (table extend z=3)",
		"[x=(1), y=(2), z=(3,'')]")

	test("table extend x=1, y=2 leftjoin (table extend z=3, zz=4)",
		"[x=(1), y=(2), z=(3,''), zz=(4,'')]")

	test("table extend x=1, y=2 summarize x, count",
		"[x=(1)]")

	test2 := func(query string, expected string) {
		t.Helper()
		q := ParseQuery(query, testTran{}, nil)
		q, _, _ = Setup(q, ReadMode, testTran{})
		assert.T(t).This(String(q)).Is(expected)
	}
	test2("table extend f=1 where f is 2",
		"nothing(table)")
	test2("table extend f=1, g=2 where f is 1",
		"table^(a) where true extend f = 1, g = 2") //TODO remove where true
	test2("table extend f=1, g=2 where f is 3",
		"nothing(table)")
	test2("tables extend x=1 join (columns extend x=2)",
		"nothing")
	test2("tables extend x=1 leftjoin (columns extend x=2)",
		"tables extend x = 1, column = '', field = ''")
	test2("table where a = 1 and a = 2",
		"nothing(table)")
	test2("(table union table2) where a = 1",
		"table^(a) where*1 a is 1 extend d = '', e = ''")
	test2("(table minus table2) where a = 1",
		"table^(a) where*1 a is 1")
}

func TestCombineFixed(t *testing.T) {
	test := func(src, oth string, expected string) {
		result, none := combineFixed(makeFixed(src), makeFixed(oth))
		if expected == "none" {
			assert.T(t).That(none == true)
		} else {
			assert.T(t).This(result).Is(makeFixed(expected))
		}
	}
	test("a 1; b 2 3; c 4", "", "a 1; b 2 3; c 4")
	test("", "a 1; b 2 3; c 4", "a 1; b 2 3; c 4")
	test("a 1; c 1 2", "c 2 3; d 4", "a 1; c 2; d 4")
	test("a 1; b 8", "z 9; a 2", "none")
	test("a 1 2 3", "a 4 5 6", "none")
}

func makeFixed(f string) []Fixed {
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

func Test_selectFixed(t *testing.T) {
	test := func(fixedStr string, expected string) {
		cols, vals := []string{"a", "b"}, []string{"1", "2"}
		fixed := makeFixed(fixedStr)
		satisfied, conflict := selectFixed(cols, vals, fixed)
		assert.That(!(satisfied && conflict))
		actual := ""
		if satisfied {
			actual = "satisfied"
		} else if conflict {
			actual = "conflict"
		}
		assert.T(t).This(actual).Is(expected)
	}
	test("", "")
	test("x 1; y 2 3 4", "")
	test("a 1; b 2", "satisfied")
	test("a 1; b 2 3 4", "")
	test("a 2", "conflict")
	test("a 2 3 4", "conflict")
	test("b 9", "conflict")
}
