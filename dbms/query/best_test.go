// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestMinimizeKeys(t *testing.T) {
	test := func(keys, expected [][]string) {
		result := minimizeKeys(keys)
		assert.T(t).This(result).Is(expected)
	}
	test([][]string{}, [][]string{})
	test([][]string{{"a"}}, [][]string{{"a"}})
	test([][]string{{"a"}, {"b", "c"}}, [][]string{{"a"}, {"b", "c"}})
	test([][]string{{"a", "b"}, {"b", "a"}}, [][]string{{"a", "b"}})
	test([][]string{{"a", "b"}, {"b", "a", "c"}}, [][]string{{"a", "b"}})
	test([][]string{{"a", "b", "c"}, {"a", "b"}}, [][]string{{"a", "b"}})
	test([][]string{{"a"}, {"a", "b"}, {"a", "c"}}, [][]string{{"a"}})
	test([][]string{{"a", "b"}, {"a", "b"}}, [][]string{{"a", "b"}})
	test([][]string{{"a"}, {"b"}, {"c"}}, [][]string{{"a"}, {"b"}, {"c"}})
	test([][]string{{}, {"a"}}, [][]string{{}})
}

var result [][]string

func BenchmarkNoOptMod(b *testing.B) {
	orig := [][]string{{"a"}, {"b"}, {"c"}, {"d"}, {"e"}, {"f"}}
	for b.Loop() {
		result = make([][]string, len(orig))
		for _, o := range orig { //nolint
			result = append(result, o)
		}
	}
}

func BenchmarkOptMod(b *testing.B) {
	orig := [][]string{{"a"}, {"b"}, {"c"}, {"d"}, {"e"}, {"f"}}
	for b.Loop() {
		om := newOptMod(orig)
		for _, o := range orig {
			om.add(o)
		}
		om.result()
	}
}

func TestOrderedN(t *testing.T) {
	test := func(index []string, order []string, fixed Fixed, expected int) {
		t.Helper()
		result := orderedn(index, order, fixed)
		assert.T(t).This(result).Is(expected)
	}

	// Basic matching - all fields match
	test([]string{"a", "b", "c"}, []string{"a", "b", "c"}, nil, 3)

	// Partial match
	test([]string{"a", "b", "c"}, []string{"a", "b"}, nil, 2)

	// No match at first field
	test([]string{"a", "b", "c"}, []string{"x", "y"}, nil, 0)

	// Index shorter than order
	test([]string{"a", "b"}, []string{"a", "b", "c"}, nil, 2)

	// Order shorter than index
	test([]string{"a", "b", "c"}, []string{"a"}, nil, 1)

	// Empty index
	test([]string{}, []string{"a", "b"}, nil, 0)

	// Empty order
	test([]string{"a", "b"}, []string{}, nil, 0)

	// Both empty
	test([]string{}, []string{}, nil, 0)

	// Fixed allows skipping in index - fixed 'b' allows index to skip 'b'
	fixed := Fixed{{col: "b", values: fixvals("1")}}
	test([]string{"a", "b", "c"}, []string{"a", "c"}, fixed, 2)

	// Fixed allows skipping in order - fixed 'b' allows order to skip 'b'
	test([]string{"a", "c"}, []string{"a", "b", "c"}, fixed, 3)

	// Fixed in both index and order
	test([]string{"a", "b", "c"}, []string{"a", "b", "c"}, fixed, 3)

	// Multiple fixed values
	fixed2 := Fixed{{col: "a", values: fixvals("1")}, {col: "c", values: fixvals("2")}}
	test([]string{"a", "b", "c"}, []string{"b"}, fixed2, 1)

	// Fixed doesn't help when fields don't match
	test([]string{"x", "y"}, []string{"a", "b"}, fixed, 0)

	// Order has fixed field that can be skipped
	fixed3 := Fixed{{col: "b", values: fixvals("1")}, {col: "c", values: fixvals("2")}}
	test([]string{"a"}, []string{"a", "b", "c"}, fixed3, 3)

	// Index exhausted before order
	test([]string{"a"}, []string{"a", "b", "c"}, nil, 1)

	// Mismatch after some matches
	test([]string{"a", "x", "c"}, []string{"a", "b", "c"}, nil, 1)

	// Single field match
	test([]string{"a"}, []string{"a"}, nil, 1)

	// Single field no match
	test([]string{"a"}, []string{"b"}, nil, 0)
}

func TestGrouped(t *testing.T) {
	oneval := []string{""}
	fixed := Fixed{{col: "f1", values: oneval}, {col: "f2", values: oneval}}
	test := func(sidx, scols string) {
		t.Helper()
		idx := strings.Fields(sidx)
		cols := strings.Fields(scols)
		nu := countUnfixed(cols, fixed)
		assert.T(t).That(grouped(idx, cols, nu, fixed))
		idx = append(idx, "x")
		assert.T(t).That(grouped(idx, cols, nu, fixed))
		cols = append(cols, "y")
		assert.T(t).That(!grouped(idx, cols, nu+1, fixed))
	}
	test("a", "a")
	test("a b", "a")
	test("a b", "b a")
	test("a f1", "f2 a")
	test("a f1 b f2", "a f1")
	test("a f1 b f2", "f1 b f2 a")

	// index too short - only has one unfixed column but need two
	idx := []string{"a"}
	cols := []string{"a", "b"}
	nu := countUnfixed(cols, fixed)
	assert.T(t).That(!grouped(idx, cols, nu, fixed))

	// missing required column in index
	idx = []string{"a", "c"}
	cols = []string{"a", "b"}
	nu = countUnfixed(cols, fixed)
	assert.T(t).That(!grouped(idx, cols, nu, fixed))

	// index starts with fixed column, then has required unfixed columns
	fixed2 := Fixed{{col: "f1", values: oneval}, {col: "f2", values: oneval}, {col: "f3", values: oneval}}
	idx = []string{"f3", "a", "b"}
	cols = []string{"a", "b"}
	nu = countUnfixed(cols, fixed2)
	assert.T(t).That(grouped(idx, cols, nu, fixed2))

	// empty index but non-zero unfixed columns should return false
	idx = []string{}
	cols = []string{"a"}
	nu = countUnfixed(cols, fixed)
	assert.T(t).That(!grouped(idx, cols, nu, fixed))
}

func TestIndexCovered(t *testing.T) {
	oneval := []string{""}
	fixed := Fixed{{col: "f1", values: oneval}, {col: "f2", values: oneval}}
	test := func(sidx, scols string, expected bool) {
		t.Helper()
		idx := strings.Fields(sidx)
		cols := strings.Fields(scols)
		result := indexCovered(idx, cols, fixed)
		assert.T(t).This(result).Is(expected)
	}

	// no fixed
	test("", "", true)
	test("a", "a", true)
	test("a b", "a b", true)
	test("a b", "b a", true)
	test("a", "b", false)
	test("a b", "a", false)
	test("a b c", "a b", false)

	// with fixed - fixed columns in index are ignored
	test("a f1", "a", true)
	test("f1 a", "a", true)
	test("a f1 b", "a b", true)
	test("f1 f2 a", "a", true)
	test("a f1 b f2", "a b", true)

	// fixed columns in cols don't help cover unfixed index columns
	test("a b", "a f1", false)

	// all index columns fixed - always covered
	test("f1 f2", "", true)
	test("f1 f2", "a", true)
}
