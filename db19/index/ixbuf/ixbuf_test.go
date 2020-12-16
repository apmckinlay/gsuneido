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

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestInsert(t *testing.T) {
	r := str.UniqueRandom(4, 8)
	const nkeys = 16000
	ib := &ixbuf{}
	for i := 0; i < nkeys; i++ {
		ib.Insert(r(), uint64(i))
	}
	assert.T(t).This(ib.size).Is(nkeys)
	// x.stats()
	ib.check()
}

func BenchmarkInsert(b *testing.B) {
	const nkeys = 100
	keys := make([]string, nkeys)
	r := str.UniqueRandom(4, 32)
	for i := 0; i < nkeys; i++ {
		keys[i] = r()
	}

	for i := 0; i < b.N; i++ {
		Ib = &ixbuf{}
		for j := 0; j < nkeys; j++ {
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
	for i := 0; i < 25; i++ {
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
		for i := 0; i < nkeys; i++ {
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
	for i := 0; i < n; i++ {
		nacts := rand.Intn(11)
		x := &ixbuf{}
		for j := 0; j < nacts; j++ {
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
	iter := ib.Iter(false)
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
		for i := 0; i < nkeys; i++ {
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
	iter := ib.Iter(false)
	for k, o, ok := iter(); ok; k, o, ok = iter() {
		fmt.Fprint(&sb, ", ", k, " ", o)
	}
	return sb.String()
}

func BenchmarkMerge(b *testing.B) {
	r := str.UniqueRandom(4, 8)
	gen := func(nkeys int) *ixbuf {
		ib := &ixbuf{}
		for i := 0; i < nkeys; i++ {
			ib.Insert(r(), 1)
		}
		return ib
	}
	x := gen(1000)
	y := gen(1)
	b.Run("bench", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
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
	rand.Seed(12345)
	r := str.UniqueRandom(4, 8)
	ib := &ixbuf{}
	for i := 0; i < nkeys; i++ {
		ib.Insert(r(), 1)
	}
	rand.Seed(12345)
	r = str.UniqueRandom(4, 8)
	for i := 0; i < nkeys; i++ {
		ib.Delete(r(), 1)
		ib.check()
	}
	assert.T(t).This(len(ib.chunks)).Is(0)
}

func TestIter(t *testing.T) {
	ib := &ixbuf{}
	iter := ib.Iter(false)
	_, _, ok := iter()
	assert.That(!ok)
	const nkeys = 1000
	for i := nkeys; i < nkeys*2; i++ {
		ib.Insert(strconv.Itoa(i), 1)
	}
	iter = ib.Iter(false)
	for i := nkeys; i < nkeys*2; i++ {
		key, _, ok := iter()
		assert.That(ok)
		assert.T(t).This(key).Is(strconv.Itoa(i))
	}
	_, _, ok = iter()
	assert.That(!ok)
}

func TestForEach(t *testing.T) {
	const nkeys = 1000
	ib := &ixbuf{}
	for i := nkeys; i < nkeys*2; i++ {
		ib.Insert(strconv.Itoa(i), 1)
	}
	i := nkeys
	ib.ForEach(func(key string, _ uint64) {
		assert.T(t).This(key).Is(strconv.Itoa(i))
		i++
	})
	assert.T(t).This(i).Is(nkeys * 2)
}

func TestIterator(t *testing.T) {
	assert := assert.T(t)
	ib := &ixbuf{}
	it := ib.Iterator()
	test := func(expected int) {
		t.Helper()
		if expected == -1 {
			assert.That(it.Eof())
		} else {
			key, off := it.Cur()
			assert.This(key).Is(strconv.Itoa(expected))
			assert.This(off).Is(uint64(expected))
		}
	}
	testNext := func(expected int) { t.Helper(); it.Next(); test(expected) }
	testPrev := func(expected int) { t.Helper(); it.Prev(); test(expected) }

	assert.That(it.Eof())
	testNext(-1)
	it.Rewind()
	testPrev(-1)
	it.Rewind()
	testNext(-1)
	testPrev(-1)

	for i := 0; i < 10; i++ {
		ib.Insert(strconv.Itoa(i), uint64(i))
	}
	it.Rewind()
	for i := 0; i < 10; i++ {
		testNext(i)
		if i == 5 {
			it.modCount-- // invalidate
		}
	}
	testNext(-1)

	it.Rewind()
	for i := 9; i >= 0; i-- {
		testPrev(i)
		if i == 5 {
			it.modCount-- // invalidate
		}
	}
	testPrev(-1)

	it.Rewind()
	testNext(0)
	testPrev(-1) // stick at eof
	testPrev(-1)
	testNext(-1)

	it.Rewind()
	testPrev(9)
	testPrev(8)
	testPrev(7)
	testNext(8)
	testNext(9) // last
	testPrev(8)

	it.Rewind()
	testNext(0)
	testNext(1)
	ib.Insert("11", uint64(11))
	testNext(11)
	testNext(2)
	ib.Delete("11", uint64(11))
	testPrev(1)

	it.Rewind()
	testPrev(9)
	testPrev(8)
	ib.Insert("77", uint64(77))
	testPrev(77)
	ib.modCount++
	testPrev(7)
}

//-------------------------------------------------------------------

func (ib *ixbuf) stats() {
	fmt.Println("size", ib.size, "chunks", len(ib.chunks),
		"avg size", int(ib.size)/len(ib.chunks), "goal", goal(ib.size))
}

func (ib *ixbuf) print() {
	fmt.Println("<<<------------------------")
	ib.stats()
	for i, c := range ib.chunks {
		if i > 0 {
			fmt.Println("+++")
		}
		c.print()
	}
	fmt.Println("------------------------>>>")
}

func (c chunk) print() {
	for _, s := range c {
		fmt.Println(s.key, s.off)
	}
}

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
