// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestParseQuery(t *testing.T) {
	test := func(s string) {
		t.Helper()
		q := ParseQuery(s)
		assert.T(t).This(q.String()).Is(s)
	}
	test("table")
	test("table sort one")
	test("table sort reverse one, two")
	test("table project one")
	test("table project one, two, three")
	test("table remove one")
	test("table remove one, two, three")
	test("table rename one to two")
	test("table rename one to two, three to four")
	test("left intersect right")
	test("left minus right")
	test("left times right")
	test("left union right")
	test("left join right")
	test("left join by(a,b) right")
	test("left leftjoin right")
	test("left leftjoin by(a,b) right")
	test("table summarize count")
	test("table summarize n = count")
	test("table summarize total one")
	test("table summarize t1 = total one")
	test("table summarize total one, count, max two")
	test("table summarize one, two, count")

	test("one union two join three")
	test("one join two sort a, b")
	test("one join two project a, b, c rename b to bb sort a, c")

	xtest := func(s, err string) {
		fn := func() { ParseQuery(s) }
		assert.T(t).This(fn).Panics(err)
	}
	xtest("table project", "expecting identifier")
	xtest("table remove", "expecting identifier")
	xtest("table rename", "expecting identifier")
	xtest("left join by() right", "invalid empty join by")
	xtest("table summarize one, two", "expecting Comma")
	xtest("table summarize total", "expecting identifier")
}

func TestParseQuery2(t *testing.T) {
	test := func(s, expected string) {
		t.Helper()
		q := ParseQuery(s)
		assert.T(t).This(q.String()).Is(expected)
	}

	test("table extend one",
		"table extend one")
	test("table extend one, two = a + b",
		"table extend one, two = Nary(Add a b)")

	test("table where x > 1",
		"table where Binary(Gt x 1)")
	test("table where a and b and c",
		"table where Nary(And a b c)")
	test("table where n in (1,2,3)",
		"table where In(n [1 2 3])")
}
