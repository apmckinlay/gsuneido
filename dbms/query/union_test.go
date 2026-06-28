// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
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

	fixcost, varcost := Optimize(u, CursorMode, OrderedReq(index, 1))
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

func TestUnionDuplicateBug(t *testing.T) {
	assert.TestOnlyIndividually(t)

	// QueryFuzz seed: 6218445892
	db, err := db19.OpenDb("../../suneido.db", stor.Read, true)
	if err != nil {
		panic(err.Error())
	}
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	tran := db.NewReadTran()
	const qstr = `(((cus union (cus union cus)) union (cus union (cus union ((cus extend x1 = "1") union cus)))) where c2 <= "77")`

	minimal := []string{
		`((cus union (cus union cus)) union (cus union (cus union ((cus extend x1 = "1") union cus)))) where c2 <= "77"`,
	}
	for _, mq := range minimal {
		q := ParseQuery(mq, tran, nil)
		q = q.Transform()
		fmt.Println("=== transformed ===")
		fmt.Println(format(0, q, 0))
		req := UnorderedReq(1)
		// capture the top union's chosen approach
		top := q.(*Union)
		fc, vc, app := top.optimize(ReadMode, req)
		fmt.Printf("top union optimize: cost=%d+%d app=%T %+v\n", fc, vc, app, app)
		Optimize(q, ReadMode, req)
		q = SetApproach(q, req, tran)
		fmt.Println("=== optimized ===")
		fmt.Println(format(0, q, 0))
		walkUnion(t, q, 0)
		hdr := q.Header()
		seen := map[string]int{}
		th := &Thread{}
		var n, dups int
		for row := q.Get(th, Next); row != nil; row = q.Get(th, Next) {
			n++
			s := row2str(hdr, row)
			if seen[s] > 0 {
				dups++
			}
			seen[s]++
		}
		fmt.Printf("rows=%d unique=%d dups=%d\n", n, len(seen), dups)
		for s, c := range seen {
			if c > 1 {
				fmt.Printf("  DUP(%d): %s\n", c, s)
			}
		}
	}
	// V1 for comparison
	q1 := ParseQuery(minimal[0], tran, nil)
	q1 = q1.Transform()
	Optimize(q1, ReadMode, UnorderedReq(1))
	q1 = SetApproach(q1, UnorderedReq(1), tran)
	hdr1 := q1.Header()
	seen1 := map[string]int{}
	th1 := &Thread{}
	var n1, dups1 int
	for row := q1.Get(th1, Next); row != nil; row = q1.Get(th1, Next) {
		n1++
		s := row2str(hdr1, row)
		if seen1[s] > 0 {
			dups1++
		}
		seen1[s]++
	}
	fmt.Printf("V1 rows=%d unique=%d dups=%d\n", n1, len(seen1), dups1)
	fmt.Println("=== V1 plan ===")
	fmt.Println(format(0, q1, 0))
	walkUnion(t, q1, 0)
}

func walkUnion(t *testing.T, q Query, depth int) {
	if u, ok := q.(*Union); ok {
		fmt.Printf("%*sunion strat=%v keyIndex=%v disjoint=%q\n", depth*2, "", u.strat, u.keyIndex, u.disjoint)
		fmt.Printf("%*s  allCols=%v\n", depth*2, "", u.allCols)
		fmt.Printf("%*s  keys=%v\n", depth*2, "", u.Keys())
		fmt.Printf("%*s  src1 cols=%v keys=%v\n", depth*2, "", u.source1.Columns(), u.source1.Keys())
		fmt.Printf("%*s  src2 cols=%v keys=%v\n", depth*2, "", u.source2.Columns(), u.source2.Keys())
	}
	if q1, ok := q.(q1i); ok {
		walkUnion(t, q1.Source(), depth+1)
	}
	if q2, ok := q.(q2i); ok {
		walkUnion(t, q2.Source2(), depth+1)
	}
}
