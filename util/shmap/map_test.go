// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package shmap

import (
	"fmt"
	"hash/maphash"
	"math/bits"
	"math/rand/v2"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestRandom(t *testing.T) {
	// use an explicit seed and NewMapInt for reproducibility
	r := rand.New(rand.NewPCG(12345, 67890))
	m := NewMapInt[int, int]()
	gm := make(map[int]int)
	const max = 20000
	var i int
	for i = range 1_000_000 {
		k := r.IntN(max)
		switch r.IntN(3) {
		case 0:
			v := r.IntN(max)
			m.Put(k, v)
			gm[k] = v
		case 1:
			m.Del(k)
			delete(gm, k)
		case 2:
			v, _ := m.Get(k)
			assert.This(v).Is(gm[k])
		}
		if i%1000 == 0 {
			m.check()
			assert.This(m.Len()).Is(len(gm))
		}
	}
	m.summary()
	fmt.Println("nPuts", nPuts, "nProbes", nPutProbes,
		"average", float64(nPutProbes)/float64(nPuts))
	fmt.Println("nGets", nGets, "nProbes", nGetProbes,
		"average", float64(nGetProbes)/float64(nGets))
	m2 := m.Copy()
	m2.check()
	// m2.summary()
}

func TestIter(t *testing.T) {
	m := NewMapInt[int, int]()
	gm := make(map[int]int)
	for i := range 1000 {
		m.Put(i, i*10)
		gm[i] = i * 10
	}
	it := m.Iter()
	for {
		k, v, ok := it()
		if !ok {
			break
		}
		assert.That(v == k*10)
		_, ok = gm[k]
		assert.That(ok)
		delete(gm, k)
	}
	assert.This(len(gm)).Is(0)
}

func TestPut(t *testing.T) {
	m := NewMapCmpable[int, int]()
	m.Put(0, 0)
	m.Put(1, 10)
	m.Put(2, 22)
	m.Put(3, 30)
	m.Put(2, 20) // update
	// m.print()
	m.check()

	const n = 128000
	for i := 4; i < n; i++ {
		// fmt.Println(">>> Put", i)
		m.Put(i, i*10)
	}
	m.summary()
	fmt.Println("nPuts", nPuts, "nProbes", nPutProbes,
		"average", float64(nPutProbes)/float64(nPuts))
	m.check()
	// m.print()
	// for i := 4; i < n; i++ {
	// 	m.Put(i, i*10)
	// }
	// m.check()
}

func TestDel(t *testing.T) {
	const n = 14 // gives one full group
	m := NewMapInt[int, int]()
	for i := range n {
		m.Put(i, i)
	}
	m.check()
	for i := range n {
		m.Del(i)
		m.check()
	}
	m.print()
}

func TestAnyEmpty(t *testing.T) {
	assert.That(anyEmpty(0x8888888888888888) == false)
	assert.That(anyEmpty(0x8888887f88888888) == false) // deleted
	assert.That(anyEmpty(0x8888008888888888) == true)
}

func TestBits(t *testing.T) {
	a := uint64(0x1144334455664400)
	fmt.Printf("%016x\n", a)
	fmt.Printf("%016x find\n", findByte(a, 0x44))
	b := uint64(0x4444444444444444)
	fmt.Printf("%016x\n", b)
	c := a ^ ^b
	fmt.Printf("%016x\n", c)
	d := c
	d = d & (d << 1)
	d = d & (d << 2)
	d = d & (d << 4)
	d = d & 0x80808080_80808080
	fmt.Printf("%016x\n", d)
	fmt.Println(bits.LeadingZeros64(d) / 8)

	// 1122334455667700
	// 4444444444444444
	// aa9988ffeeddccbb
	// 	     f000000000
	// 3

	x := uint64(0x48)
	// x = x | (x << 8)
	// x = x | (x << 16)
	// x = x | (x << 32)
	x *= 0x01010101_01010101
	fmt.Printf("%016x\n", x)
}

func TestFindByte(t *testing.T) {
	a := uint64(0x4444444444444444)
	fmt.Printf("%016x\n", a)
	b := findByte(a, 0x44)
	fmt.Printf("%016x\n", b)
	for b != 0 {
		i := b.first()
		fmt.Println(i)
		b = b.dropFirst()
	}
	fmt.Printf("%016x\n", b)
}

func TestProbe(t *testing.T) {
	hash := uint64(0x1234567890abcdef)
	seq := makeProbeSeq(hash, 0x0f)
	for range 16 {
		fmt.Println(seq.offset)
		seq = seq.next()
	}
}

func TestTableIndex(t *testing.T) {
	m := &Map[int, int, integer[int]]{}
	m.depth = 0
	m.dir = make([]*table[int, int], 1)
	h := uint64(0)
	assert.This(m.tableIndex(h)).Is(0)
	h = 0x80_00_00_00_00_00_00_00
	assert.This(m.tableIndex(h)).Is(0)

	m.depth = 1
	m.dir = make([]*table[int, int], 2)
	h = uint64(0)
	assert.This(m.tableIndex(h)).Is(0)
	h = 0x80_00_00_00_00_00_00_00
	assert.This(m.tableIndex(h)).Is(1)

	m.depth = 2
	m.dir = make([]*table[int, int], 4)
	h = uint64(0)
	assert.This(m.tableIndex(h)).Is(0)
	h = 0x80_00_00_00_00_00_00_00
	assert.This(m.tableIndex(h)).Is(2)
}

var H uint64

func BenchmarkMaphash(b *testing.B) {
	seed := maphash.MakeSeed()
	for i := 0; i < b.N; i++ {
		H = maphash.Comparable(seed, uint64(i))
	}
}

func BenchmarkPut(b *testing.B) {
	for b.Loop() {
		m := NewMapInt[int, int]()
		for i := range 10000 {
			m.Put(i, i)
		}
	}
}

func BenchmarkGet(b *testing.B) {
	m := NewMapInt[int, int]()
	for i := range 100000 {
		m.Put(i, i)
	}
	for b.Loop() {
		for i := range 100000 {
			m.Get(i)
		}
	}
}

func BenchmarkGetString(b *testing.B) {
	const n = 1000
	rs := make([]string, n)
	m := NewMapCmpable[string, string]()
	for i := range n {
		s := str.Random(5, 10)
		m.Put(s, s)
		rs[i] = s
	}
	for b.Loop() {
		for i := range n {
			m.Get(rs[i])
		}
	}
}

func TestGet(t *testing.T) { // for profiling
	m := NewMapInt[int, int]()
	for i := range 1000 {
		m.Put(i, i)
	}
	for range 100000 {
		for i := range 1000 {
			m.Get(i)
		}
	}
}

func BenchmarkTmp(b *testing.B) {
	for b.Loop() {
		m := NewMapInt[int, int]()
		for i := range 8 {
			m.Put(i, i)
		}
	}
}
