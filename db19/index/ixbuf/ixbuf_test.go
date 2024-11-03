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

	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestInsert(t *testing.T) {
	r := str.UniqueRandom(4, 8)
	const nkeys = 16000
	ib := &ixbuf{}
	for i := range nkeys {
		ib.Insert(r(), uint64(i))
	}
	assert.T(t).This(ib.size).Is(nkeys)
	// ib.stats()
	ib.check()
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
			ib.Insert(r(), uint64(i))
		}
		big = Merge(big, ib)
	}
	assert.T(t).This(big.size).Is(n * m)
	// big.stats()
	big.check()
}

func BenchmarkInsert(b *testing.B) {
	const nkeys = 100
	keys := make([]string, nkeys)
	r := str.UniqueRandom(4, 32)
	for i := range nkeys {
		keys[i] = r()
	}

	for range b.N {
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
	ib.check()

	c := &ixbuf{}
	for i := range 25 {
		c.Insert(strconv.Itoa(i), uint64(i))
	}
	ib = Merge(b, c, a)
	assert(ib.size).Is(a.size + b.size + c.size)
	// x.print()
	ib.check()

	a.Insert("c", 3)
	b.Insert("d", 4)
	ib = Merge(b, a)
	// x.print()
	assert(ib.size).Is(4)
	assert(len(ib.chunks)).Is(1)
	ib.check()

	r := str.UniqueRandom(4, 8)
	gen := func(nkeys int) *ixbuf {
		t := &ixbuf{}
		for range nkeys {
			t.Insert(r(), 1)
		}
		// t.print()
		t.check()
		return t
	}
	a = gen(1000)
	b = gen(100)
	c = gen(10)
	ib = Merge(a, b, c)
	// x.print()
	assert(ib.size).Is(a.size + b.size + c.size)
	ib.check()
	a.check()
	b.check()
	c.check()
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
	x.check()
}

func TestMergeRandom(*testing.T) {
	n := 10_000
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
	assert.T(t).This(x.String()).Is("3, a 1, b 22, d 4")
}

func (ib *ixbuf) String() string {
	var sb strings.Builder
	fmt.Fprint(&sb, ib.size)
	iter := ib.Iter()
	for k, o, ok := iter(); ok; k, o, ok = iter() {
		fmt.Fprint(&sb, ", ", k, " ", o)
	}
	return sb.String()
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
	b.Run("bench", func(b *testing.B) {
		for range b.N {
			Ib = Merge(x, y)
		}
	})
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
		ib.check()
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

	for i := range 10 {
		ib.Insert(strconv.Itoa(i), uint64(i))
	}
	it.Rewind()
	for i := range 10 {
		testNext(i)
	}
	testNext(eof)

	it.Rewind()
	for i := 9; i >= 0; i-- {
		testPrev(i)
	}
	testPrev(eof)

	it.Rewind()
	testNext(0)
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
	data := strings.Fields("a b c d e f g h")
	for _, d := range data {
		ib.Insert(d, 1)
	}
	it := ib.Iterator()
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

//-------------------------------------------------------------------

//lint:ignore U1000 for debugging
func (ib *ixbuf) stats() {
	fmt.Println("size", ib.size, "chunks", len(ib.chunks),
		"avg size", int(ib.size)/len(ib.chunks), "goal", goal(ib.size)*2/3)
}

//lint:ignore U1000 for debugging
func (ib *ixbuf) check() {
	n := 0
	prev := ""
	for _, c := range ib.chunks {
		assert.That(len(c) > 0)
		for _, s := range c {
			if s.key <= prev {
				panic("out of order " + prev + " " + s.key)
			}
			prev = s.key
			n++
		}
	}
	assert.This(ib.size).Is(n)
}

//lint:ignore U1000 for debugging
func chunkstr(c chunk) string {
	switch len(c) {
	case 0:
		return "empty"
	case 1:
		return fmt.Sprint(c[0].key)
	default:
		return fmt.Sprint(c[0].key, " -> ", c.lastKey(), " (", len(c), ")")
	}
}

// func TestCombine(t *testing.T) {
// 	Combine(123 | Delete, 456 | Update)
// }
