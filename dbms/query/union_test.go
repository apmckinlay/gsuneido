// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"math/rand"
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

type r struct {
	key    []string
	i1, i2 int
}
type s []string

func TestUnion_MergeIndexes(t *testing.T) {
	var list []r
	callback := func(key []string, i1, i2 int) {
		list = append(list, r{key, i1, i2})
	}
	mergeIndexes(
		[][]string{{"a", "b"}},
		[][]string{{"b", "a", "x"}, {"a", "b"}},
		[][]string{{"b", "a", "y"}, {"b", "z", "a"}, {"b", "a", "z"}},
		callback)
	assert.This(list).Is([]r{
		{s{"a", "b"}, -1, -1},
		{s{"a", "b"}, 0, 0},
		{s{"a", "b"}, 0, 2}})
}

func TestUnion_MergeSwitchDir(t *testing.T) {
	db := heapDb()
	db.adm("create one (a) key(a)")
	db.adm("create two (a) key(a)")
	db.act("insert { a: 1 } into one")
	db.act("insert { a: 2 } into two")
	db.act("insert { a: 3 } into two")
	db.act("insert { a: 4 } into one")
	db.act("insert { a: 4 } into two")
	db.act("insert { a: 5 } into two")
	db.act("insert { a: 6 } into one")
	db.act("insert { a: 7 } into one")
	db.act("insert { a: 8 } into two")
	db.act("insert { a: 9 } into two")
	tran := db.NewReadTran()
	q := ParseQuery("one union two", tran, nil)
	q, _, _ = Setup(q, ReadMode, tran)
	// fmt.Println(Format(q))
	assert.That(strings.Contains(q.String(), "merge"))

	get := func(dir Dir) Row {
		// fmt.Println(dir)
		row := q.Get(nil, dir)
		// fmt.Println("=>", row)
		return row
	}
	hdr := q.Header()
	a := func(row Row) int {
		if row == nil {
			return 0
		}
		return ToInt(row.GetVal(hdr, "a", nil, nil))
	}

	// random walk
	cur := 0 // rewound
	for range 1000 {
		dir := Next
		if rand.Int()%2 == 0 {
			dir = Prev
			if cur--; cur < 0 {
				cur = 9
			}
		} else {
			cur = (cur + 1) % 10
		}
		row := get(dir)
		assert.This(a(row)).Is(cur)
		if row == nil {
			q.Rewind()
		}
	}
}

func TestUnion_removeNonexistentEmpty(t *testing.T) {
	srccols := []string{"a", "b", "c"}
	test := func(colsIn, valsIn, colsOut, valsOut []string) {
		cols, vals := removeNonexistentEmpty(srccols, colsIn, valsIn)
		assert.This(cols).Is(colsOut)
		assert.This(vals).Is(valsOut)
	}
	test(nil, nil, nil, nil)
	test([]string{}, []string{}, []string{}, []string{})
	test([]string{"a", "c", "x"}, []string{"1", "2", "3"},
		[]string{"a", "c", "x"}, []string{"1", "2", "3"})
	test([]string{"a", "n", "c", "x"}, []string{"1", "", "2", "3"},
		[]string{"a", "c", "x"}, []string{"1", "2", "3"})
	test([]string{"x", "y"}, []string{"", ""}, nil, nil)
}

func TestBestMergeIndexes(t *testing.T) {
	test := func(order []string, idx1, idx2, keys1, keys2 [][]string,
		expIdx1, expIdx2, expKey []string, expImpossible bool) {
		t.Helper()
		src1 := &QueryMock{
			ColumnsResult: []string{"a", "b", "c", "d", "e", "f"},
			IndexesResult: idx1,
			KeysResult:    keys1,
			FixedResult:   []Fixed{},
			MetricsResult: &metrics{},
		}
		src2 := &QueryMock{
			ColumnsResult: []string{"a", "b", "c", "d", "x", "y"},
			IndexesResult: idx2,
			KeysResult:    keys2,
			FixedResult:   []Fixed{},
			MetricsResult: &metrics{},
		}
		u := &Union{}
		u.source1 = src1
		u.source2 = src2

		resIdx1, resIdx2, resKey, fixcost, _ := u.bestMergeIndexes(order, ReadMode, 1.0)
		assert.T(t).This(resIdx1).Is(expIdx1)
		assert.T(t).This(resIdx2).Is(expIdx2)
		assert.T(t).This(resKey).Is(expKey)
		if expImpossible {
			assert.T(t).That(fixcost == impossible)
		}
	}

	// nil order
	test(nil,
		[][]string{{"a", "b"}},
		[][]string{{"a", "b"}},
		[][]string{{"a", "b"}},
		[][]string{{"a", "b"}},
		[]string{"a", "b"},
		[]string{"a", "b"},
		[]string{"a", "b"},
		false)

	// simple case with keys and order equal
	test([]string{"a", "b"},
		[][]string{{"a", "b"}},
		[][]string{{"a", "b"}},
		[][]string{{"a", "b"}},
		[][]string{{"a", "b"}},
		[]string{"a", "b"},
		[]string{"a", "b"},
		[]string{"a", "b"},
		false)

	// example from spec:
	// required order(c,b)
	// source1 key(a,c,b) index(c,b,a)
	// source2 key(b,a,c) index(c,b,d,a)
	test([]string{"c", "b"},
		[][]string{{"c", "b", "a"}},
		[][]string{{"c", "b", "d", "a"}},
		[][]string{{"a", "c", "b"}},
		[][]string{{"b", "a", "c"}},
		[]string{"c", "b", "a"},
		[]string{"c", "b", "d", "a"},
		[]string{"a", "c", "b"},
		false)

	test(nil,
		[][]string{{"c", "b", "a"}},
		[][]string{{"c", "b", "d", "a"}},
		[][]string{{"a", "c", "b"}},
		[][]string{{"b", "a", "c"}},
		[]string{"c", "b", "a"},
		[]string{"c", "b", "d", "a"},
		[]string{"a", "c", "b"},
		false)

	// no matching indexes
	test([]string{"c", "b"},
		[][]string{{"x", "y"}},
		[][]string{{"c", "b", "d"}},
		[][]string{{"x", "y"}},
		[][]string{{"c", "b"}},
		nil,
		nil,
		nil,
		true)

	// both sources have empty keys - don't need order
	emptyKey := [][]string{{}}
	test([]string{"x"},
		[][]string{{"a", "b"}},
		[][]string{{"a", "b"}},
		emptyKey,
		emptyKey,
		nil,
		nil,
		nil,
		false)

	// only source1 has empty key - need index on source2 that includes a key
	test([]string{"a"},
		[][]string{{"a", "b"}},
		[][]string{{"a", "b"}},
		emptyKey,
		[][]string{{"a", "b"}},
		nil,
		[]string{"a", "b"},
		[]string{"a", "b"},
		false)

	// only source2 has empty key - need index on source1 that includes a key
	test([]string{"a"},
		[][]string{{"a", "b"}},
		[][]string{{"a", "b"}},
		[][]string{{"a", "b"}},
		emptyKey,
		[]string{"a", "b"},
		nil,
		[]string{"a", "b"},
		false)
}

func TestUnion_DisjointRequiredIndexNoKey(t *testing.T) {
	index := []string{"a"}
	src1 := &QueryMock{
		ColumnsResult: []string{"a", "k", "d"},
		HeaderResult:  SimpleHeader([]string{"a", "k", "d"}),
		IndexesResult: [][]string{index},
		KeysResult:    [][]string{{"k"}},
		FixedResult:   []Fixed{NewFixed("d", SuInt(1))},
		NrowsN:        1,
		NrowsP:        1,
		RowSizeResult: 1,
		LookupLevels:  1,
	}
	src2 := &QueryMock{
		ColumnsResult: []string{"a", "k", "d"},
		HeaderResult:  SimpleHeader([]string{"a", "k", "d"}),
		IndexesResult: [][]string{index},
		KeysResult:    [][]string{{"k"}},
		FixedResult:   []Fixed{NewFixed("d", SuInt(2))},
		NrowsN:        1,
		NrowsP:        1,
		RowSizeResult: 1,
		LookupLevels:  1,
	}

	u := NewUnion(src1, src2)
	assert.T(t).This(u.disjoint).Is("d")

	fixcost, varcost := Optimize(u, CursorMode, index, 1)
	assert.T(t).That(fixcost+varcost < impossible)
}

func TestIndexContainsKey(t *testing.T) {
	assert := assert.T(t)

	// Empty keys list
	assert.This(indexContainsKey([]string{"a", "b"}, nil)).Is(nil)

	// Index contains key
	assert.This(indexContainsKey(
		[]string{"a", "b", "c"},
		[][]string{{"a", "b"}},
	)).Is([]string{"a", "b"})

	// Index doesn't contain key
	assert.This(indexContainsKey(
		[]string{"a", "b"},
		[][]string{{"a", "b", "c"}},
	)).Is(nil)

	// Multiple keys, first match returned
	assert.This(indexContainsKey(
		[]string{"a", "b", "c"},
		[][]string{{"d", "e"}, {"a", "c"}},
	)).Is([]string{"a", "c"})
}

func TestKeyFieldOrder(t *testing.T) {
	assert := assert.T(t)

	// Basic case
	assert.This(keyFieldOrder(
		[]string{"c", "b", "a"},
		[]string{"a", "c", "b"},
	)).Is([]string{"c", "b", "a"})

	// Key fields in different order in index
	assert.This(keyFieldOrder(
		[]string{"a", "x", "b"},
		[]string{"a", "b"},
	)).Is([]string{"a", "b"})

	// Empty key
	assert.This(keyFieldOrder(
		[]string{"a", "b"},
		[]string{},
	)).Is([]string{})
}

func TestSameKeyFieldOrder(t *testing.T) {
	assert := assert.T(t)

	// Same order
	assert.That(sameKeyFieldOrder(
		[]string{"c", "b", "d", "a"},
		[]string{"b", "a", "c"},
		[]string{"c", "b", "a"},
	))

	// Different order
	assert.That(!sameKeyFieldOrder(
		[]string{"a", "b", "c"},
		[]string{"a", "b", "c"},
		[]string{"c", "b", "a"},
	))
}
