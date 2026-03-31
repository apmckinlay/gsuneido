// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ixbuf

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestInsert(t *testing.T) {
	r := str.UniqueRandom(4, 8)
	const nkeys = 16000
	ib := &ixbuf{}
	for i := range nkeys {
		ib.Insert(r(), uint64(i+1))
	}
	assert.T(t).This(ib.size).Is(nkeys)
	// ib.stats()
	ib.Check()
}

func TestBig(t *testing.T) {
	big := &ixbuf{}
	r := str.UniqueRandom(4, 8)
	n := 256
	if testing.Short() {
		n = 64
	}
	const m = 1000
	for range n {
		ib := &ixbuf{}
		for i := range m {
			ib.Insert(r(), uint64(i+1))
		}
		big = Merge(big, ib)
	}
	assert.T(t).This(big.size).Is(n * m)
	// big.stats()
	big.Check()
}

func BenchmarkInsert(b *testing.B) {
	const nkeys = 100
	keys := make([]string, nkeys)
	r := str.UniqueRandom(4, 32)
	for i := range nkeys {
		keys[i] = r()
	}

	for b.Loop() {
		Ib = &ixbuf{}
		for j := range nkeys {
			Ib.Insert(keys[j], uint64(j))
		}
	}
}

var Ib *ixbuf

func TestMerge(t *testing.T) {
	assert := assert.T(t).This

	ib := Merge(&ixbuf{}, &ixbuf{}, &ixbuf{})
	assert(ib.size).Is(0)

	a := &ixbuf{}
	a.Insert("a", 1)
	ib = Merge(&ixbuf{}, a, &ixbuf{})
	assert(ib.size).Is(1)

	b := &ixbuf{}
	b.Insert("b", 2)
	ib = Merge(b, &ixbuf{}, a)
	assert(ib.size).Is(2)
	assert(len(ib.chunks)).Is(1)
	// x.print()
	ib.Check()

	c := &ixbuf{}
	for i := range 25 {
		c.Insert(strconv.Itoa(i), uint64(i+1))
	}
	ib = Merge(b, c, a)
	assert(ib.size).Is(a.size + b.size + c.size)
	// x.print()
	ib.Check()

	a.Insert("c", 3)
	b.Insert("d", 4)
	ib = Merge(b, a)
	// x.print()
	assert(ib.size).Is(4)
	assert(len(ib.chunks)).Is(1)
	ib.Check()

	r := str.UniqueRandom(4, 8)
	gen := func(nkeys int) *ixbuf {
		t := &ixbuf{}
		for range nkeys {
			t.Insert(r(), 1)
		}
		// t.print()
		t.Check()
		return t
	}
	a = gen(1000)
	b = gen(100)
	c = gen(10)
	ib = Merge(a, b, c)
	// x.print()
	assert(ib.size).Is(a.size + b.size + c.size)
	ib.Check()
	a.Check()
	b.Check()
	c.Check()
}

func TestMergeBug(*testing.T) {
	a := &ixbuf{}
	a.Insert("a", 1)
	a.Insert("d", 1)
	b := &ixbuf{}
	b.Insert("b", 1)
	b.Insert("c", 1)
	c := &ixbuf{}
	c.Insert("e", 1)
	c.Insert("f", 1)
	x := Merge(a, b, c)
	// x.print()
	x.Check()
}

func TestMergeRandom(*testing.T) {
	n := 100_000
	if testing.Short() {
		n = 1000
	}
	var data chunk
	ib := &ixbuf{}
	var s slot
	r := str.UniqueRandom(4, 8)
	for range n {
		nacts := rand.Intn(11)
		x := &ixbuf{}
		for range nacts {
			k := rand.Intn(4)
			switch {
			case k == 0 || k == 1 || len(data) == 0: // add
				s = slot{key: r(), off: uint64(rand.Uint32())}
				// fmt.Println("add", s)
				data = append(data, s)
				x.Insert(s.key, s.off)
			case k == 2: // update
				i := rand.Intn(len(data))
				data[i].off = uint64(rand.Uint32())
				s = data[i]
				// fmt.Println("update", s)
				x.Update(s.key, s.off)
			case k == 3: // delete
				i := rand.Intn(len(data))
				s = data[i]
				// fmt.Println("delete", s)
				data[i] = data[len(data)-1]
				data = data[:len(data)-1]
				x.Delete(s.key, s.off)
			}
		}
		// fmt.Println(x)
		ib = Merge(ib, x)
		// fmt.Println("=", ib)
		// fmt.Println(len(data), data)
		assert.This(ib.Len()).Is(len(data))
	}
	assert.This(ib.Len()).Is(len(data))
	sort.Sort(data)
	i := 0
	iter := ib.Iter()
	for k, o, ok := iter(); ok; k, o, ok = iter() {
		assert.This(k).Is(data[i].key)
		assert.This(o).Is(data[i].off)
		i++
	}
}

func (c chunk) Len() int           { return len(c) }
func (c chunk) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c chunk) Less(i, j int) bool { return c[i].key < c[j].key }

func TestMergeMore(t *testing.T) {
	const n = 10
	keys := make([]string, n)
	adrs := make([]uint64, n)
	nextadr := 1
	loops := 1000000
	if testing.Short() {
		loops = 10000
	}

	gen := func() *ixbuf {
		size := rand.Intn(10)
		ib := &ixbuf{}
		for range size {
			i := rand.Intn(n)
			if keys[i] == "" { // insert
				keys[i] = strconv.Itoa(i)
				adrs[i] = uint64(nextadr)
				nextadr++
				ib.Insert(keys[i], adrs[i])
			} else if rand.Intn(2) == 0 { // update
				adrs[i] = uint64(nextadr)
				nextadr++
				ib.Update(keys[i], adrs[i])
			} else { // delete
				ib.Delete(keys[i], adrs[i])
				keys[i] = ""
				adrs[i] = 0
			}
		}
		return ib
	}
	for range loops {
		nextadr = 1
		clear(keys)
		clear(adrs)
		a := gen()
		b := gen()
		c := gen()
		d := gen()
		ib := Merge(a, b, c, d)
		iter := ib.Iter()
		i := 0
		for {
			key, adr, ok := iter()
			if !ok {
				break
			}
			for i < n && keys[i] == "" {
				i++
			}
			assert.This(key).Is(keys[i])
			assert.This(adr).Is(adrs[i])
			i++
		}
	}
}

func TestMergeUneven(*testing.T) {
	r := str.UniqueRandom(4, 8)
	gen := func(nkeys int) *ixbuf {
		ib := &ixbuf{}
		for range nkeys {
			ib.Insert(r(), 1)
		}
		return ib
	}
	x := gen(1000)
	y := gen(1)
	Merge(x, y)
}

func TestMergeUpdate(t *testing.T) {
	a := &ixbuf{}
	a.Insert("a", 1)
	a.Insert("b", 2)
	a.Insert("c", 3)
	a.Insert("d", 4)
	b := &ixbuf{}
	b.Update("b", 22)
	b.Delete("c", 3)
	x := Merge(a, b)
	assert.T(t).This(x.String()).Is("a+1 b+22 d+4")
}

func BenchmarkMerge(b *testing.B) {
	r := str.UniqueRandom(4, 8)
	gen := func(nkeys int) *ixbuf {
		ib := &ixbuf{}
		for range nkeys {
			ib.Insert(r(), 1)
		}
		return ib
	}
	x := gen(1000)
	y := gen(1)
	for b.Loop() {
		Ib = Merge(x, y)
	}
}

func TestGoal(t *testing.T) {
	assert.T(t).This(goal(0)).Is(24) // min
	assert.T(t).This(goal(100)).Is(24)
	assert.T(t).This(goal(1000)).Is(48)
	assert.T(t).This(goal(4000)).Is(96)
}

func TestDelete(t *testing.T) {
	const nkeys = 1000
	r := str.UniqueRandom(4, 8, 12345)
	ib := &ixbuf{}
	for range nkeys {
		ib.Insert(r(), 1)
	}
	r = str.UniqueRandom(4, 8, 12345)
	for range nkeys {
		ib.Delete(r(), 1)
		ib.Check()
	}
	assert.T(t).This(len(ib.chunks)).Is(0)
}

func TestLookup(*testing.T) {
	const nkeys = 1000
	r := str.UniqueRandom(4, 8, 12345)
	ib := &ixbuf{}
	for i := 1; i < nkeys; i++ {
		ib.Insert(r(), uint64(i))
	}
	r = str.UniqueRandom(4, 8, 12345)
	for i := 1; i < nkeys; i++ {
		k := r()
		assert.This(ib.Lookup(k)).Is(i)
		assert.This(ib.Lookup(k + " ")).Is(0)
		assert.This(ib.Lookup(k + "~")).Is(0)
	}
	for range nkeys {
		assert.This(ib.Lookup(r())).Is(0) // nonexistent
	}
}

func TestIter(t *testing.T) {
	ib := &ixbuf{}
	iter := ib.Iter()
	_, _, ok := iter()
	assert.That(!ok)
	const nkeys = 1000
	for i := nkeys; i < nkeys*2; i++ {
		ib.Insert(strconv.Itoa(i), 1)
	}
	iter = ib.Iter()
	for i := nkeys; i < nkeys*2; i++ {
		key, _, ok := iter()
		assert.That(ok)
		assert.T(t).This(key).Is(strconv.Itoa(i))
	}
	_, _, ok = iter()
	assert.That(!ok)
}

func TestIterator(t *testing.T) {
	assert := assert.T(t)
	const eof = -1
	ib := &ixbuf{}
	it := ib.Iterator()
	test := func(expected int) {
		t.Helper()
		if expected == eof {
			assert.That(it.Eof())
		} else {
			key, off := it.Cur()
			assert.This(key).Is(strconv.Itoa(expected))
			assert.This(off).Is(uint64(expected))
		}
	}
	testNext := func(expected int) { t.Helper(); it.Next(); test(expected) }
	testPrev := func(expected int) { t.Helper(); it.Prev(); test(expected) }

	test(eof)
	testNext(eof)
	it.Rewind()
	testPrev(eof)
	it.Rewind()
	testNext(eof)
	testPrev(eof)

	for i := 1; i < 10; i++ {
		ib.Insert(strconv.Itoa(i), uint64(i))
	}
	it.Rewind()
	for i := 1; i < 10; i++ {
		testNext(i)
	}
	testNext(eof)

	it.Rewind()
	for i := 9; i >= 1; i-- {
		testPrev(i)
	}
	testPrev(eof)

	it.Rewind()
	testNext(1)
	testPrev(eof) // stick at eof
	testPrev(eof)
	testNext(eof)

	it.Rewind()
	testPrev(9)
	testPrev(8)
	testPrev(7)
	testNext(8)
	testNext(9) // last
	testPrev(8)

	// Seek to nonexistent
	it.Seek("00")
	test(1) // leaves us on next
	it.Seek("99")
	test(9) // or last
}

func TestIterRange(t *testing.T) {
	ib := &ixbuf{}
	data := strings.FieldsSeq("a b c d e f g h")
	for d := range data {
		ib.Insert(d, 1)
	}
	it := ib.Iterator().(*Iterator)
	test := func(fn func(), expected string) {
		fn()
		assert.That(it.state == within)
		assert.This(it.cur.key).Is(expected)
	}
	test(it.Next, "a")
	it.Rewind()
	test(it.Prev, "h")

	it.Range(Range{Org: "c", End: ixkey.Max})
	test(it.Next, "c")
	it.Range(Range{Org: "c+", End: ixkey.Max})
	test(it.Next, "d")

	it.Range(Range{End: "f"})
	test(it.Prev, "e")
	it.Range(Range{End: "f+"})
	test(it.Prev, "f")

	it.Range(Range{Org: "c", End: "g"})
	test(it.Next, "c")
	test(it.Next, "d")
	test(it.Next, "e")
	test(it.Next, "f")
	it.Next()
	assert.T(t).That(it.Eof())

	it.Rewind()
	test(it.Prev, "f")
	test(it.Prev, "e")
	test(it.Prev, "d")
	test(it.Prev, "c")
	it.Prev()
	assert.T(t).That(it.Eof())

	it.Range(Range{Org: "c", End: "g"})
	it.Seek("c")
	assert.T(t).This(it.cur.key).Is("c")
	it.Seek("b")
	assert.T(t).That(it.Eof())
	it.Seek("f")
	assert.T(t).This(it.cur.key).Is("f")
	it.Seek("g")
	assert.T(t).That(it.Eof())
}

func TestIxbufSearch(t *testing.T) {
	ib := &ixbuf{}
	ib.Insert("a\x00\x001", 11)
	ib.Insert("b\x00\x002", 22)
	_, _, i := ib.search("a")
	assert.T(t).This(i).Is(0)
	_, _, i = ib.search("a\x00\x00\xff")
	assert.T(t).This(i).Is(1)
}

func TestSkipScanNext(t *testing.T) {
	ib := &ixbuf{}
	rows := [][2]string{
		{"a", "01"}, {"a", "03"}, {"a", "05"},
		{"b", "02"}, {"b", "04"},
		{"c", "00"}, {"c", "03"},
		{"d", "02"},
	}
	for i, r := range rows {
		ib.Insert(ixkey.CompKey(r[0], r[1]), uint64(i+1))
	}
	it := ib.Iterator().(*Iterator)
	it.SkipScan(iface.All, Range{Org: "02", End: "04"}, 1)

	var got []string
	for it.Next(); !it.Eof(); it.Next() {
		f, s := ixkey.SplitPrefixSuffix(it.Key(), 1)
		got = append(got, f+":"+s)
	}
	assert.T(t).This(got).Is([]string{"a:03", "b:02", "c:03", "d:02"})
}

// TestSkipScanEmptyPrefix tests skip scan when the key's first field is empty string.
// This triggered a bug because skipGroup is initialized to "" which is also a valid prefix.
func TestSkipScanEmptyPrefix(t *testing.T) {
	ib := &ixbuf{}
	ib.Insert(ixkey.CompKey("", "One"), 1)
	ib.Insert(ixkey.CompKey("", "Two"), 2)
	ib.Insert(ixkey.CompKey("", "Three"), 3)
	it := ib.Iterator().(*Iterator)
	it.SkipScan(iface.All, Range{Org: "Two", End: "Two\x00"}, 1)
	it.Next()
	assert.T(t).That(!it.Eof())
	_, s := ixkey.SplitPrefixSuffix(it.Key(), 1)
	assert.T(t).This(s).Is("Two")
	it.Next()
	assert.T(t).That(it.Eof())
}

// TestSkipScanEmptyPrefixPrev tests backward skip scan with empty first field.
// Exercises the suffix >= End path in skipRetreatToMatch with skipGroup="" collision.
func TestSkipScanEmptyPrefixPrev(t *testing.T) {
	ib := &ixbuf{}
	ib.Insert(ixkey.CompKey("", "01"), 1)
	ib.Insert(ixkey.CompKey("", "03"), 2)
	ib.Insert(ixkey.CompKey("", "05"), 3)
	it := ib.Iterator().(*Iterator)
	it.SkipScan(iface.All, Range{Org: "01", End: "04"}, 1)
	it.Prev()
	assert.T(t).That(!it.Eof())
	_, s := ixkey.SplitPrefixSuffix(it.Key(), 1)
	assert.T(t).This(s).Is("03")
	it.Prev()
	assert.T(t).That(!it.Eof())
	_, s = ixkey.SplitPrefixSuffix(it.Key(), 1)
	assert.T(t).This(s).Is("01")
	it.Prev()
	assert.T(t).That(it.Eof())
}

func TestSkipScanNoMatches(t *testing.T) {
	ib := &ixbuf{}
	ib.Insert(ixkey.CompKey("a", "01"), 1)
	ib.Insert(ixkey.CompKey("b", "01"), 2)
	it := ib.Iterator().(*Iterator)
	it.SkipScan(iface.All, Range{Org: "02", End: "03"}, 1)
	it.Next()
	assert.T(t).That(it.Eof())
}

func TestSkipScanPrev(t *testing.T) {
	ib := &ixbuf{}
	rows := [][2]string{
		{"a", "01"}, {"a", "03"}, {"a", "05"},
		{"b", "02"}, {"b", "04"},
		{"c", "00"}, {"c", "03"},
		{"d", "02"},
	}
	for i, r := range rows {
		ib.Insert(ixkey.CompKey(r[0], r[1]), uint64(i+1))
	}
	it := ib.Iterator().(*Iterator)
	it.SkipScan(iface.All, Range{Org: "02", End: "04"}, 1)

	var got []string
	for it.Prev(); !it.Eof(); it.Prev() {
		f, s := ixkey.SplitPrefixSuffix(it.Key(), 1)
		got = append(got, f+":"+s)
	}
	assert.T(t).This(got).Is([]string{"d:02", "c:03", "b:02", "a:03"})
}

func TestSkipScanNextMultiPrefix(t *testing.T) {
	ib := &ixbuf{}
	rows := [][3]string{
		{"a", "x", "01"}, {"a", "x", "03"}, {"a", "x", "05"},
		{"a", "y", "02"}, {"a", "y", "04"},
		{"b", "x", "03"},
		{"b", "y", "00"}, {"b", "y", "02"},
	}
	for i, r := range rows {
		ib.Insert(ixkey.CompKey(r[0], r[1], r[2]), uint64(i+1))
	}
	it := ib.Iterator().(*Iterator)
	it.SkipScan(iface.All, Range{Org: "02", End: "04"}, 2)

	var got []string
	for it.Next(); !it.Eof(); it.Next() {
		p, s := ixkey.SplitPrefixSuffix(it.Key(), 2)
		got = append(got, p+":"+s)
	}
	assert.T(t).This(got).Is([]string{
		ixkey.CompKey("a", "x") + ":03",
		ixkey.CompKey("a", "y") + ":02",
		ixkey.CompKey("b", "x") + ":03",
		ixkey.CompKey("b", "y") + ":02",
	})
}

func TestSkipScanRangeDisablesSkipScan(t *testing.T) {
	ib := &ixbuf{}
	ib.Insert(ixkey.CompKey("a", "01"), 1)
	ib.Insert(ixkey.CompKey("a", "02"), 2)
	ib.Insert(ixkey.CompKey("b", "01"), 3)
	it := ib.Iterator().(*Iterator)

	it.SkipScan(iface.All, Range{Org: "01", End: "03"}, 1)
	it.Next()
	assert.T(t).That(!it.Eof())

	// Range() disables skip scan and returns to normal range mode
	it.Range(Range{Org: ixkey.CompKey("a", "02"), End: ixkey.CompKey("a", "03")})
	it.Next()
	assert.T(t).That(!it.Eof())
	f, s := ixkey.SplitPrefixSuffix(it.Key(), 1)
	assert.T(t).This(f).Is("a")
	assert.T(t).This(s).Is("02")
	it.Next()
	assert.T(t).That(it.Eof())
}

func TestSkipScanPrevGroupOutOfRange(t *testing.T) {
	// group "a" has only keys with suffix >= End; Prev should return eof
	// (exercises the eof path in skipRetreatToMatch after skipSeekGroupEnd)
	ib := &ixbuf{}
	ib.Insert(ixkey.CompKey("a", "20"), 1)
	ib.Insert(ixkey.CompKey("a", "25"), 2)
	it := ib.Iterator().(*Iterator)
	it.SkipScan(iface.All, Range{Org: "08", End: "19"}, 1)
	it.Prev()
	assert.T(t).That(it.Eof())
}

func TestSkipScanPrefixRangeNext(t *testing.T) {
	ib := &ixbuf{}
	rows := [][2]string{
		{"a", "01"}, {"a", "03"},
		{"b", "02"},
		{"c", "03"},
		{"d", "02"},
	}
	for i, r := range rows {
		ib.Insert(ixkey.CompKey(r[0], r[1]), uint64(i+1))
	}
	it := ib.Iterator().(*Iterator)
	it.SkipScan(Range{Org: "b", End: "d"}, Range{Org: "02", End: "04"}, 1)

	var got []string
	for it.Next(); !it.Eof(); it.Next() {
		f, s := ixkey.SplitPrefixSuffix(it.Key(), 1)
		got = append(got, f+":"+s)
	}
	assert.T(t).This(got).Is([]string{"b:02", "c:03"})
}

func TestSkipScanPrefixRangePrev(t *testing.T) {
	ib := &ixbuf{}
	rows := [][2]string{
		{"a", "01"}, {"a", "03"},
		{"b", "02"},
		{"c", "03"},
		{"d", "02"},
	}
	for i, r := range rows {
		ib.Insert(ixkey.CompKey(r[0], r[1]), uint64(i+1))
	}
	it := ib.Iterator().(*Iterator)
	it.SkipScan(Range{Org: "b", End: "d"}, Range{Org: "02", End: "04"}, 1)

	var got []string
	for it.Prev(); !it.Eof(); it.Prev() {
		f, s := ixkey.SplitPrefixSuffix(it.Key(), 1)
		got = append(got, f+":"+s)
	}
	assert.T(t).This(got).Is([]string{"c:03", "b:02"})
}

func TestSkipScanRandomParallelWithSubset(t *testing.T) {
	const (
		org   = "08"
		end   = "19"
		steps = 30000
	)

	// groups with suffixes designed to exercise edge cases:
	//   "a","b" have only suffixes below org
	//   "c","d" have only suffixes within [org,end)
	//   "e","f" have only suffixes >= end
	//   "g","h","i","j" span the range boundary
	// Extra filler groups "p".."z" ensure the ixbuf uses multiple chunks.
	type groupSpec struct {
		first    string
		suffixes []string
	}
	groups := []groupSpec{
		{"a", []string{"01", "03", "05", "07"}},
		{"b", []string{"02", "04", "06"}},
		{"c", []string{"08", "10", "12", "15", "18"}},
		{"d", []string{"09", "11", "14", "17"}},
		{"e", []string{"19", "21", "25"}},
		{"f", []string{"20", "22", "28"}},
		{"g", []string{"05", "08", "12", "19", "22"}},
		{"h", []string{"06", "09", "13", "18", "20"}},
		{"i", []string{"07", "10", "18"}},
		{"j", []string{"11", "16", "19", "24"}},
	}
	// Add filler groups to push total entries above one chunk (goal >= 24)
	for _, letter := range "pqrstuvwxyz" {
		first := string(letter)
		var suffixes []string
		for s := 1; s <= 5; s++ {
			suffixes = append(suffixes, fmt.Sprintf("%02d", s))
		}
		groups = append(groups, groupSpec{first, suffixes})
	}

	fullIb := &ixbuf{}
	subIb := &ixbuf{}

	var off uint64 = 1
	for _, g := range groups {
		for _, suffix := range g.suffixes {
			k := ixkey.CompKey(g.first, suffix)
			fullIb.Insert(k, off)
			if org <= suffix && suffix < end {
				subIb.Insert(k, off)
			}
			off++
		}
	}
	// Verify multiple chunks are used
	assert.T(t).That(len(fullIb.chunks) > 1)

	firsts := make([]string, len(groups))
	for i, g := range groups {
		firsts[i] = g.first
	}

	fit := fullIb.Iterator().(*Iterator)
	fit.SkipScan(iface.All, Range{Org: org, End: end}, 1)
	sit := subIb.Iterator()

	assertSame := func(step int, op string) {
		t.Helper()
		fhc := fit.HasCur()
		shc := sit.HasCur()
		assert.T(t).Msg(fmt.Sprintf("step %d op %s hascur", step, op)).
			This(fhc).Is(shc)
		assert.T(t).Msg(fmt.Sprintf("step %d op %s eof", step, op)).
			This(fit.Eof()).Is(sit.Eof())
		if !fhc {
			return
		}
		fk, fo := fit.Cur()
		sk, so := sit.Cur()
		assert.T(t).Msg(fmt.Sprintf("step %d op %s key", step, op)).
			This(fk).Is(sk)
		assert.T(t).Msg(fmt.Sprintf("step %d op %s off", step, op)).
			This(fo).Is(so)
	}

	seek := func(first, s string) {
		t.Helper()
		k := ixkey.CompKey(first, s)
		fit.Seek(k)
		sit.Seek(k)
	}

	seekSuffixes := []string{
		"00",
		"01", "02",
		"07", "075",
		"08",
		"085",
		"10", "11",
		"135",
		"18",
		"185",
		"19",
		"20", "21",
		"29", "30",
	}

	r := rand.New(rand.NewSource(1))
	assertSame(0, "init")
	for step := 1; step <= steps; step++ {
		switch r.Intn(4) {
		case 0:
			fit.Next()
			sit.Next()
			assertSame(step, "Next")
		case 1:
			fit.Prev()
			sit.Prev()
			assertSame(step, "Prev")
		case 2:
			fit.Rewind()
			sit.Rewind()
			assertSame(step, "Rewind")
		case 3:
			first := firsts[r.Intn(len(firsts))]
			s := seekSuffixes[r.Intn(len(seekSuffixes))]
			seek(first, s)
			assertSame(step, "Seek")
		}
	}
}

func BenchmarkSkipScanBreakevenVsFullScan(b *testing.B) {
	// 32 groups × 500 records = 16k total, goal ≈ 192 => each group spans ~2-3 chunks
	const (
		groups          = 32
		recordsPerGroup = 500
		suffixWidth     = 4
		startSuffix     = 50
	)

	ib := &ixbuf{}
	for g := range groups {
		first := fmt.Sprintf("g%04d", g)
		for s := range recordsPerGroup {
			suffix := fmt.Sprintf("%0*d", suffixWidth, s)
			off := uint64(g*recordsPerGroup + s + 1)
			ib.Insert(ixkey.CompKey(first, suffix), off)
		}
	}

	widths := []int{50, 100, 150, 200, 250, 300, 350, 400, 450}
	maxWidth := widths[len(widths)-1]
	fullOrg := fmt.Sprintf("%0*d", suffixWidth, startSuffix)
	fullEnd := fmt.Sprintf("%0*d", suffixWidth, startSuffix+maxWidth)
	fullExpected := groups * maxWidth

	for _, width := range widths {
		width := width
		endSuffix := startSuffix + width
		org := fmt.Sprintf("%0*d", suffixWidth, startSuffix)
		end := fmt.Sprintf("%0*d", suffixWidth, endSuffix)
		expected := groups * width

		b.Run(fmt.Sprintf("skip_w%03d", width), func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				it := ib.Iterator().(*Iterator)
				it.SkipScan(iface.All, Range{Org: org, End: end}, 1)
				n := 0
				for it.Next(); !it.Eof(); it.Next() {
					n++
				}
				if n != expected {
					b.Fatalf("skip scan width %d expected %d got %d", width, expected, n)
				}
			}
		})
	}

	b.Run("full_once", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			it := ib.Iterator()
			n := 0
			for it.Next(); !it.Eof(); it.Next() {
				_, suffix := ixkey.SplitPrefixSuffix(it.Key(), 1)
				if fullOrg <= suffix && suffix < fullEnd {
					n++
				}
			}
			if n != fullExpected {
				b.Fatalf("full scan expected %d got %d", fullExpected, n)
			}
		}
	})
}

//-------------------------------------------------------------------

// func (ib *ixbuf) stats() {
// 	fmt.Println("size", ib.size, "chunks", len(ib.chunks),
// 		"avg size", int(ib.size)/len(ib.chunks), "goal", goal(ib.size)*2/3)
// }

// func chunkstr(c chunk) string {
// 	switch len(c) {
// 	case 0:
// 		return "empty"
// 	case 1:
// 		return fmt.Sprint(c[0].key)
// 	default:
// 		return fmt.Sprint(c[0].key, " -> ", c.lastKey(), " (", len(c), ")")
// 	}
// }

// func TestCombine(t *testing.T) {
// 	Combine(123 | Delete, 456 | Update)
// }
