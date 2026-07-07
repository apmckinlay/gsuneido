// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"math/rand"
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/options"
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
	test := func(selsIn, selsOut Sels) {
		assert.This(removeNonexistentEmpty(srccols, selsIn)).Is(selsOut)
	}
	test(nil, nil)
	test(Sels{}, Sels{})
	test(Sels{{"a", "1"}, {"c", "2"}, {"x", "3"}},
		Sels{{"a", "1"}, {"c", "2"}, {"x", "3"}})
	test(Sels{{"a", "1"}, {"n", ""}, {"c", "2"}, {"x", "3"}},
		Sels{{"a", "1"}, {"c", "2"}, {"x", "3"}})
	test(Sels{{"x", ""}, {"y", ""}}, nil)
}

func TestUnion_DisjointRequiredIndexNoKey(t *testing.T) {
	index := []string{"a"}
	src1 := &QueryMock{
		ColumnsResult: []string{"a", "k", "d"},
		HeaderResult:  SimpleHeader([]string{"a", "k", "d"}),
		IndexesResult: [][]string{index},
		KeysResult:    [][]string{{"k"}},
		FixedResult:   Fixed{NewFix("d", SuInt(1))},
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
		FixedResult:   Fixed{NewFix("d", SuInt(2))},
		NrowsN:        1,
		NrowsP:        1,
		RowSizeResult: 1,
		LookupLevels:  1,
	}

	u := NewUnion(src1, src2)
	assert.T(t).This(u.disjoint).Is("d")

	fixcost, varcost := Optimize(u, CursorMode, OrderReq(index, 1))
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

func TestKeyPrefixOfIndex(t *testing.T) {
	assert := assert.T(t)
	// last key field is at position 1
	assert.This(keyPrefixOfIndex(
		[]string{"a", "b", "c"},
		[]string{"b", "d"},
	)).Is([]string{"a", "b"})
	// all index fields are in key
	assert.This(keyPrefixOfIndex(
		[]string{"a", "b", "c"},
		[]string{"a", "b", "c"},
	)).Is([]string{"a", "b", "c"})
	// no key fields in index
	assert.This(keyPrefixOfIndex(
		[]string{"a", "b"},
		[]string{"x", "y"},
	)).Is(nil)
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

func TestUnion_StrictCompareDb(t *testing.T) {
	defer func(sc bool) { options.StrictCompareDb = sc }(options.StrictCompareDb)
	options.StrictCompareDb = true
	db := heapDb()
	db.adm("create one (k, i) key(k)")
	db.act("insert { k: 1, i: '' } into one")
	db.adm("create two (k, j) key(k)")
	db.act("insert { k: 1, j: 2 } into two")

	queryAll(db.Database, "(one union two) where i isnt '' where i > 0")

	queryAll(db.Database, "(one union two) where i isnt '' and i > 0")

	queryAll(db.Database, "(one union two) where Number?(i) where i > 0")

	queryAll(db.Database, "(one union two) where Number?(i) and i > 0")
}

func TestUnionLookupBug(t *testing.T) {
	db := heapDb()
	db.adm("create cus (ck, c1, c2) key(ck)")
	db.act(`insert {ck: "0", c1: '6'} into cus`)
	db.act(`insert {ck: "10", c2: '2'} into cus`)
	db.act(`insert {ck: "11", c1: "16", c2: '3'} into cus`)
	db.act(`insert {ck: "12", c1: '1', c2: "10"} into cus`)
	db.act(`insert {ck: "15", c1: '0', c2: "12"} into cus`)
	db.act(`insert {ck: "16", c1: '4', c2: "13"} into cus`)
	db.act(`insert {ck: "18", c1: '2', c2: "11"} into cus`)
	db.act(`insert {ck: '3', c1: '5', c2: '4'} into cus`)
	db.act(`insert {ck: '7', c1: "16", c2: "11"} into cus`)
	db.act(`insert {ck: '9', c1: '2', c2: "18"} into cus`)
	db.adm("create cus2 (ck, c2) key(ck)")
	db.act(`insert {ck: "0" } into cus2`)
	db.act(`insert {ck: "10", c2: '2'} into cus2`)
	db.act(`insert {ck: "11", c2: '3'} into cus2`)
	db.act(`insert {ck: "12", c2: "10"} into cus2`)
	db.act(`insert {ck: "15", c2: "12"} into cus2`)
	db.act(`insert {ck: "16", c2: "13"} into cus2`)
	db.act(`insert {ck: "18", c2: "11"} into cus2`)
	db.act(`insert {ck: '3', c2: '4'} into cus2`)
	db.act(`insert {ck: '7', c2: "11"} into cus2`)
	db.act(`insert {ck: '9', c2: "18"} into cus2`)

	queryHashAll(db.Database, `cus union (cus union cus)`)
	queryHashAll(db.Database, `cus2 union (cus union cus)`)
}

func queryHashAll(db *db19.Database, query string) {
	tran := db.NewReadTran()
	q := ParseQuery(query, tran, nil)
	th := &Thread{}

	h1 := NewQueryHasher(q.Header()).CheckDups()
	for _, row := range q.Simple(th) {
		h1.Row(row)
	}

	q, _, _ = Setup(q, ReadMode, tran)
	h2 := NewQueryHasher(q.Header()).CheckDups()
	for row := q.Get(th, Next); row != nil; row = q.Get(th, Next) {
		h2.Row(row)
	}
	assert.This(h2.Result(true)).Is(h1.Result(true))
}
