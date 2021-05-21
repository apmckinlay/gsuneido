// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package sset

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestContains(*testing.T) {
	test := func(x, y string, expected bool) {
		assert.Msg(x + " : " + y).
			That(expected == Contains(strings.Fields(x), y))
	}
	test("", "", false)
	test("", "x", false)
	test("a b c", "", false)
	test("a b c", "x", false)
	test("a b c", "a", true)
	test("a b c", "b", true)
	test("a b c", "c", true)
}

func TestEqual(*testing.T) {
	test := func(x, y string, expected bool) {
		assert.Msg(x + " : " + y).
			That(expected == Equal(strings.Fields(x), strings.Fields(y)))
	}
	test("", "", true)
	test("a b c", "c b a", true)
	test("a b a", "a b c", true) // failure from duplicates

	test("", "a b c", false)
	test("a b c", "", false)
	test("a b c", "a B c", false)
	test("a b c", "a b a", false) // duplicates on right side
	x := randList(100)
	assert.That(Equal(x, x))
	y := Copy(x)
	assert.That(Equal(x, y))
	y[99] = "~"
	assert.That(!Equal(x, y))
	y[0] = "~" // not sorted now
	assert.That(!Equal(x, y))
}

func TestSubset(*testing.T) {
	test := func(x, y string, expected bool) {
		assert.Msg(x + " : " + y).
			That(expected == Subset(strings.Fields(x), strings.Fields(y)))
	}
	test("", "", true)
	test("a b c", "a b c", true)
	test("a b c", "c a b", true)
	test("a b c", "c a", true)
	test("c a", "a b c", false)
}

func TestDisjoint(*testing.T) {
	test := func(x, y string, expected bool) {
		assert.Msg(x + " : " + y).
			That(expected == Disjoint(strings.Fields(x), strings.Fields(y)))
	}
	test("", "", true)
	test("a b c", "a b c", false)
	test("a b c", "c a b", false)
	test("a b c", "d e f", true)
}

func BenchmarkEqual(b *testing.B) {
	const n = 100
	x := randList(n)
	y := append([]string{}, x...)
	x[n-1] = "~" // differ at the end
	for i := 0; i < b.N; i++ {
		BM = Equal(x, y)
	}
}

var BM bool

func randList(n int) []string {
	r := str.UniqueRandom(4, 16)
	x := make([]string, n)
	for i := 0; i < n; i++ {
		x[i] = r()
	}
	return x
}

func TestUnion(*testing.T) {
	test := func(x, y, expected string) {
		assert.Msg(x + " union " + y).
			This(Union(strings.Fields(x), strings.Fields(y))).
			Is(strings.Fields(expected))
	}
	test("", "", "")
	test("a b c", "", "a b c")
	test("", "a b c", "a b c")
	test("a b c", "a b c", "a b c")
	test("a b c d", "e f", "a b c d e f")
	test("e f", "a b c d", "e f a b c d")
	test("a b c d", "c d e f", "a b c d e f")

	x := randList(100)
	assert.This(Union(x, x)).Is(x)
	y := Copy(x)
	assert.This(Union(x, y)).Is(x)
	y = y[2:98]
	assert.This(Union(x, y)).Is(x)
	assert.That(Equal(x, Union(y, x)))
}

func TestIntersect(*testing.T) {
	test := func(x, y, expected string) {
		assert.Msg(x + " intersect " + y).
			This(Intersect(strings.Fields(x), strings.Fields(y))).
			Is(strings.Fields(expected))
	}
	test("", "", "")
	test("a b c", "", "")
	test("a b c", "d e f", "")
	test("a b c d", "c d e f", "c d")
	test("a b c d", "c", "c")
	test("a b c", "a b c", "a b c")

	x := randList(100)
	assert.This(Intersect(x, x)).Is(x)
	y := Copy(x)
	assert.This(Intersect(x, y)).Is(y)
	y = y[2:98]
	assert.This(Intersect(x, y)).Is(y)
}

func TestDifference(*testing.T) {
	test := func(x, y, expected string) {
		assert.Msg(x + " difference " + y).
			This(Difference(strings.Fields(x), strings.Fields(y))).
			Is(strings.Fields(expected))
	}
	test("", "", "")
	test("a b c", "", "a b c")
	test("a b c", "d e f", "a b c")
	test("a b c d", "c d e f", "a b")
	test("a b c d", "c", "a b d")
	test("a b c", "a b c", "")

	x := randList(100)
	assert.This(Difference(x, []string{})).Is(x)
	assert.This(Difference(x, x)).Is([]string{})
	assert.This(Difference([]string{}, x)).Is([]string{})
	y := Copy(x)
	assert.This(Difference(x, y)).Is([]string{})
	assert.This(Difference(x, y[:50])).Is(x[50:])
	assert.This(Difference(x, y[50:])).Is(x[:50])
}

func TestStartsWithSet(t *testing.T) {
	assert := assert.T(t).That
	assert(StartsWithSet([]string{"a", "b", "c"}, []string{}))
	assert(StartsWithSet([]string{"a", "b", "c"}, []string{"a"}))
	assert(!StartsWithSet([]string{"a", "b", "c"}, []string{"b"}))
	assert(!StartsWithSet([]string{"a", "b", "c"}, []string{"c"}))
	assert(StartsWithSet([]string{"a", "b", "c"}, []string{"a", "b"}))
	assert(StartsWithSet([]string{"a", "b", "c"}, []string{"b", "a"}))
	assert(!StartsWithSet([]string{"a", "b", "c"}, []string{"c", "a"}))
	assert(StartsWithSet([]string{"a", "b", "c"}, []string{"a", "b", "c"}))
	assert(StartsWithSet([]string{"a", "b", "c"}, []string{"c", "a", "b"}))
	assert(!StartsWithSet([]string{"a", "b", "c"}, []string{"c", "a", "d"}))
	assert(!StartsWithSet([]string{"a"}, []string{"b"}))
	assert(!StartsWithSet([]string{"a"}, []string{"b", "a"}))
}

func TestUnique(*testing.T) {
	uniq := func(list ...string) list {
		return list
	}
	uniq().is()
	uniq("a").is()
	uniq("a", "b", "c").is()
	uniq("a", "a", "b").is("a", "b")
	uniq("a", "b", "b").is("a", "b")
	uniq("a", "b", "b", "c").is("a", "b", "c")
}

type list []string

func (input list) is(expected ...string) {
	if len(expected) == 0 {
		assert.This(Unique(input)).Is([]string(input))
	} else {
		assert.This(Unique(input)).Is(expected)
	}
}
