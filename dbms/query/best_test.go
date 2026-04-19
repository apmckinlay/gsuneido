// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestWithoutDupsOrSupersets(t *testing.T) {
	test := func(keys, expected [][]string) {
		result := withoutDupsOrSupersets(keys)
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
	test := func(index []string, order []string, fixed []Fixed, expected int) {
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
	fixed := []Fixed{{col: "b", values: fixvals("1")}}
	test([]string{"a", "b", "c"}, []string{"a", "c"}, fixed, 2)

	// Fixed allows skipping in order - fixed 'b' allows order to skip 'b'
	test([]string{"a", "c"}, []string{"a", "b", "c"}, fixed, 3)

	// Fixed in both index and order
	test([]string{"a", "b", "c"}, []string{"a", "b", "c"}, fixed, 3)

	// Multiple fixed values
	fixed2 := []Fixed{{col: "a", values: fixvals("1")}, {col: "c", values: fixvals("2")}}
	test([]string{"a", "b", "c"}, []string{"b"}, fixed2, 1)

	// Fixed doesn't help when fields don't match
	test([]string{"x", "y"}, []string{"a", "b"}, fixed, 0)

	// Order has fixed field that can be skipped
	fixed3 := []Fixed{{col: "b", values: fixvals("1")}, {col: "c", values: fixvals("2")}}
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
	fixed := []Fixed{{col: "f1", values: oneval}, {col: "f2", values: oneval}}
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
	fixed2 := []Fixed{{col: "f1", values: oneval}, {col: "f2", values: oneval}, {col: "f3", values: oneval}}
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

func TestBestLookupIndex(t *testing.T) {
	test := func(q Query, expected []string) {
		t.Helper()
		best := bestLookupIndex(q, ReadMode, 100, 1, nil)
		assert.T(t).This(best.index).Is(expected)
	}

	// ineligible index should be ignored even if cheaper
	test(&bestLookupIndexMock{
		QueryMock: QueryMock{
			ColumnsResult:     []string{"date", "item", "id", "cost"},
			FixedResult:       []Fixed{},
			IndexesResult:     [][]string{{"item"}, {"date", "item", "id"}},
			KeysResult:        [][]string{{"date", "item", "id"}},
			LookupLevels:      1,
			SingleTableResult: true,
		},
		costs: costs{
			"item":         {fix: 1, varc: 1},
			"date,item,id": {fix: 10, varc: 10},
		},
	}, []string{"date", "item", "id"})

	// among eligible indexes, lower cost should win
	test(&bestLookupIndexMock{
		QueryMock: QueryMock{
			ColumnsResult:     []string{"a", "b", "c"},
			FixedResult:       []Fixed{},
			IndexesResult:     [][]string{{"a", "b"}, {"b", "a"}},
			KeysResult:        [][]string{{"a", "b"}},
			LookupLevels:      1,
			SingleTableResult: true,
		},
		costs: costs{
			"a,b": {fix: 40, varc: 30},
			"b,a": {fix: 10, varc: 10},
		},
	}, []string{"b", "a"})

	// fixed key columns make any index eligible (nColsUnfixed == 0)
	test(&bestLookupIndexMock{
		QueryMock: QueryMock{
			ColumnsResult:     []string{"a", "b", "x", "y"},
			FixedResult:       []Fixed{NewFixed("a", SuInt(1)), NewFixed("b", SuInt(1))},
			IndexesResult:     [][]string{{"x"}, {"y"}},
			KeysResult:        [][]string{{"a", "b"}},
			LookupLevels:      1,
			SingleTableResult: true,
		},
		costs: costs{
			"x": {fix: 30, varc: 10},
			"y": {fix: 5, varc: 5},
		},
	}, []string{"y"})

	// fallback to logical keys when no physical index is eligible
	test(&bestLookupIndexMock{
		QueryMock: QueryMock{
			ColumnsResult:     []string{"a", "b", "x"},
			FixedResult:       []Fixed{},
			IndexesResult:     [][]string{{"x"}},
			KeysResult:        [][]string{{"a", "b"}, {"a"}},
			LookupLevels:      1,
			SingleTableResult: true,
		},
		costs: costs{
			"x":   {fix: 1, varc: 1},
			"a,b": {fix: 50, varc: 10},
			"a":   {fix: 10, varc: 10},
		},
	}, []string{"a"})
}

type bestLookupIndexMock struct {
	QueryMock
	costs
}

type costs map[string]struct {
	fix  Cost
	varc Cost
}

func (m *bestLookupIndexMock) optimize(_ Mode, index []string, _ float64) (Cost, Cost, any) {
	c, ok := m.costs[strings.Join(index, ",")]
	if !ok {
		return impossible, impossible, nil
	}
	return c.fix, c.varc, nil
}

func TestBestLookupIndexGrouped(t *testing.T) {
	test := func(q Query, cols []string, expected []string) {
		t.Helper()
		best := bestLookupIndex(q, ReadMode, 100, 1, cols)
		assert.T(t).This(best.index).Is(expected)
	}

	// grouped index must also be lookup-eligible; otherwise fallback to cols
	test(&bestLookupIndexMock{
		QueryMock: QueryMock{
			ColumnsResult:     []string{"a", "b", "k"},
			FixedResult:       []Fixed{},
			IndexesResult:     [][]string{{"a", "b"}},
			KeysResult:        [][]string{{"k"}},
			LookupLevels:      1,
			SingleTableResult: true,
		},
		costs: costs{
			"a,b": {fix: 20, varc: 20},
		},
	}, []string{"a", "b"}, []string{"a", "b"})

	// among eligible grouped indexes, lower cost should win
	test(&bestLookupIndexMock{
		QueryMock: QueryMock{
			ColumnsResult:     []string{"a", "b", "c", "d"},
			FixedResult:       []Fixed{},
			IndexesResult:     [][]string{{"a", "b", "c"}, {"b", "a", "d"}},
			KeysResult:        [][]string{{"a", "b"}},
			LookupLevels:      1,
			SingleTableResult: true,
		},
		costs: costs{
			"a,b,c": {fix: 40, varc: 20},
			"b,a,d": {fix: 5, varc: 5},
			"a,b":   {fix: 30, varc: 30},
		},
	}, []string{"a", "b"}, []string{"b", "a", "d"})

	// if no physical index groups by cols, fallback to cols
	test(&bestLookupIndexMock{
		QueryMock: QueryMock{
			ColumnsResult:     []string{"a", "b", "x", "y"},
			FixedResult:       []Fixed{},
			IndexesResult:     [][]string{{"x"}, {"y"}},
			KeysResult:        [][]string{{"a", "b"}},
			LookupLevels:      1,
			SingleTableResult: true,
		},
		costs: costs{
			"x":   {fix: 1, varc: 1},
			"y":   {fix: 1, varc: 1},
			"a,b": {fix: 10, varc: 10},
		},
	}, []string{"a", "b"}, []string{"a", "b"})

	// fixed columns reduce required grouped columns (nColsUnfixed)
	test(&bestLookupIndexMock{
		QueryMock: QueryMock{
			ColumnsResult:     []string{"a", "b", "x", "y"},
			FixedResult:       []Fixed{NewFixed("a", SuInt(1))},
			IndexesResult:     [][]string{{"b", "x"}, {"b", "y"}},
			KeysResult:        [][]string{{"a", "b"}},
			LookupLevels:      1,
			SingleTableResult: true,
		},
		costs: costs{
			"b,x": {fix: 50, varc: 10},
			"b,y": {fix: 10, varc: 10},
			"a,b": {fix: 30, varc: 30},
		},
	}, []string{"a", "b"}, []string{"b", "y"})
}
