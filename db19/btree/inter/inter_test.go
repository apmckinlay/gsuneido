// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package inter

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestInsert(t *testing.T) {
	r := str.UniqueRandom(4, 8)
	const nkeys = 16000
	x := &T{}
	for i := 0; i < nkeys; i++ {
		x.Insert(r(), uint64(i))
	}
	assert.T(t).This(x.size).Is(nkeys)
	// x.stats()
	x.check()
}

func BenchmarkInsert(b *testing.B) {
	const nkeys = 100
	keys := make([]string, nkeys)
	r := str.UniqueRandom(4, 32)
	for i := 0; i < nkeys; i++ {
		keys[i] = r()
	}

	for i := 0; i < b.N; i++ {
		X = &T{}
		for j := 0; j < nkeys; j++ {
			X.Insert(keys[j], uint64(j))
		}
	}
}

var X *T

func TestMerge(t *testing.T) {
	assert := assert.T(t).This

	x := Merge(&T{}, &T{}, &T{})
	assert(x.size).Is(0)

	a := &T{}
	a.Insert("a", 1)
	x = Merge(&T{}, a, &T{})
	assert(x.size).Is(1)

	b := &T{}
	b.Insert("b", 2)
	x = Merge(b, &T{}, a)
	assert(x.size).Is(2)
	assert(len(x.chunks)).Is(1)
	// x.print()
	x.check()

	c := &T{}
	for i := 0; i < 25; i++ {
		c.Insert(strconv.Itoa(i), uint64(i))
	}
	x = Merge(b, c, a)
	assert(x.size).Is(a.size + b.size + c.size)
	// x.print()
	x.check()

	a.Insert("c", 3)
	b.Insert("d", 4)
	x = Merge(b, a)
	// x.print()
	assert(x.size).Is(4)
	assert(len(x.chunks)).Is(1)
	x.check()

	r := str.UniqueRandom(4, 8)
	gen := func(nkeys int) *T {
		t := &T{}
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
	x = Merge(a, b, c)
	// x.print()
	assert(x.size).Is(a.size + b.size + c.size)
	x.check()
	a.check()
	b.check()
	c.check()
}

func TestMergeBug(*testing.T) {
	a := &T{}
	a.Insert("a", 1)
	a.Insert("d", 1)
	b := &T{}
	b.Insert("b", 1)
	b.Insert("c", 1)
	c := &T{}
	c.Insert("e", 1)
	c.Insert("f", 1)
	x := Merge(a, b, c)
	// x.print()
	x.check()
}

func TestMergeRandom(*testing.T) {
	n := 1000
	if testing.Short() {
		n = 100
	}
	for i := 0; i < n; i++ {
		r := str.UniqueRandom(4, 8)
		nin := 2 + rand.Intn(11)
		in := make([]*T, nin)
		for j := range in {
			in[j] = &T{}
			size := rand.Intn(1000)
			for k := 0; k < size; k++ {
				in[j].Insert(r(), 1)
			}
		}
		Merge(in...).check()
	}
}

func TestMergeUneven(*testing.T) {
	r := str.UniqueRandom(4, 8)
	gen := func(nkeys int) *T {
		t := &T{}
		for i := 0; i < nkeys; i++ {
			t.Insert(r(), 1)
		}
		return t
	}
	x := gen(1000)
	y := gen(1)
	Merge(x, y)
}

func BenchmarkMerge(b *testing.B) {
	r := str.UniqueRandom(4, 8)
	gen := func(nkeys int) *T {
		t := &T{}
		for i := 0; i < nkeys; i++ {
			t.Insert(r(), 1)
		}
		return t
	}
	x := gen(1000)
	y := gen(1)
	b.Run("bench", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			X = Merge(x, y)
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
	x := &T{}
	for i := 0; i < nkeys; i++ {
		x.Insert(r(), 1)
	}
	r = str.UniqueRandom(8, 12)
	for i := 0; i < nkeys; i++ {
		assert.That(!x.Delete(r()))
	}
	rand.Seed(12345)
	r = str.UniqueRandom(4, 8)
	for i := 0; i < nkeys; i++ {
		assert.That(x.Delete(r()))
		x.check()
	}
	assert.T(t).This(len(x.chunks)).Is(0)
}

func TestIter(t *testing.T) {
	x := &T{}
	iter := x.Iter(false)
	_, _, ok := iter()
	assert.That(!ok)
	const nkeys = 1000
	for i := nkeys; i < nkeys*2; i++ {
		x.Insert(strconv.Itoa(i), 1)
	}
	iter = x.Iter(false)
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
	x := &T{}
	for i := nkeys; i < nkeys*2; i++ {
		x.Insert(strconv.Itoa(i), 1)
	}
	i := nkeys
	x.ForEach(func (key string, _ uint64) {
		assert.T(t).This(key).Is(strconv.Itoa(i))
		i++
	})
	assert.T(t).This(i).Is(nkeys*2)
}

//-------------------------------------------------------------------

func (t *T) stats() {
	fmt.Println("size", t.size, "chunks", len(t.chunks), "avg size", t.size/len(t.chunks), "goal", goal(t.size))
}

func (t *T) print() {
	fmt.Println("<<<------------------------")
	t.stats()
	for i, c := range t.chunks {
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

func (t *T) check() {
	n := 0
	prev := ""
	for _, c := range t.chunks {
		assert.That(len(c) > 0)
		for _, s := range c {
			if s.key <= prev {
				panic("out of order " + prev + " " + s.key)
			}
			prev = s.key
			n++
		}
	}
	assert.This(t.size).Is(n)
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
