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
	"golang.org/x/exp/constraints"
)

func TestRandom(t *testing.T) {
	// use an explicit seed and NewMapInt for reproducibility
	r := rand.New(rand.NewPCG(12345, 67890))
	m := NewMapInt[int, int]()
	gm := make(map[int]int)
	const max = 1400
	const n = 1_000_000
	for i := range n {
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
		if i%(n/10) == 0 {
			m.check()
			assert.This(m.Size()).Is(len(gm))
			m.Clear()
			clear(gm)
		}
	}
	// m.summary()
	// fmt.Println("nPuts", nPuts, "nProbes", nPutProbes,
	// 	"average", float64(nPutProbes)/float64(nPuts))
	// fmt.Println("nGets", nGets, "nProbes", nGetProbes,
	// 	"average", float64(nGetProbes)/float64(nGets))
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

	const n = 10_000
	for i := 4; i < n; i++ {
		// fmt.Println(">>> Put", i)
		m.Put(i, i*10)
	}
	// m.summary()
	// fmt.Println("nPuts", nPuts, "nProbes", nPutProbes,
	// 	"average", float64(nPutProbes)/float64(nPuts))
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
	// m.print()
}

func TestAnyEmpty(t *testing.T) {
	assert.That(anyEmpty(0x8888888888888888) == false)
	assert.That(anyEmpty(0x8888887f88888888) == false) // deleted
	assert.That(anyEmpty(0x8888008888888888) == true)
}

func TestBits(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
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
	if testing.Short() {
		t.SkipNow()
	}
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
	if testing.Short() {
		t.SkipNow()
	}
	hash := uint64(0x1234567890abcdef)
	seq := makeProbeSeq(hash, 0x0f)
	for range 16 {
		fmt.Println(seq.offset)
		seq = seq.next()
	}
}

var H uint64

func BenchmarkMaphash(b *testing.B) {
	seed := maphash.MakeSeed()
	for i := 0; b.Loop(); i++ {
		H = maphash.Comparable(seed, uint64(i))
	}
}

func BenchmarkPut(b *testing.B) {
	for b.Loop() {
		m := NewMapInt[int, int]()
		for i := range 1000 {
			m.Put(i, i)
		}
		for i := range 1000 {
			m.Put(i, i)
		}
	}
}

func Benchmark_grow(b *testing.B) {
	m := NewMapInt[int, int]()
	for i := range 7000 {
		m.Put(i, i)
	}
	for b.Loop() {
		m2 := *m
		m2.grow()
	}
}

func BenchmarkGet(b *testing.B) {
	m := NewMapInt[int, int]()
	for i := range 1000 {
		m.Put(i, i)
	}
	for b.Loop() {
		for i := range 1000 {
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
	if testing.Short() {
		t.SkipNow()
	}
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

func Test_emptyOrDeleted(t *testing.T) {
	w := uint64(0x8080808080808080)
	assert.This(emptyOrDeleted(w)).Is(bitset(0))
	w = uint64(0x80800080807f8080)
	assert.This(emptyOrDeleted(w)).Is(bitset(0x10000010000))
}

//-------------------------------------------------------------------

// func (m *Map[K, V, H]) summary() {
// 	fmt.Println("MAP", " ngroups", len(m.groups), "count", m.count,
// 		"growthLeft", m.growthLeft)
// }

// func (m *Map[K, V, H]) print() {
// 	fmt.Println("==============================================")
// 	m.summary()
// 	for gi := range m.groups {
// 		fmt.Println("group", gi)
// 		grp := &m.groups[gi]
// 		ctrls := grp.control
// 		for i := range groupSize {
// 			c := uint8(ctrls)
// 			ctrls >>= 8
// 			fmt.Printf("%d %2x: ", i, c)
// 			if c == empty {
// 				fmt.Println("empty")
// 			} else if c == deleted {
// 				fmt.Println("deleted")
// 			} else {
// 				fmt.Printf("%v => %v\n", grp.keys[i], grp.vals[i])
// 			}
// 		}
// 	}
// }

func (m *Map[K, V, H]) check() {
	entries := 0
	deletes := 0
	hasEmpty := false
	for gi := range m.groups {
		grp := &m.groups[gi]
		ctrls := grp.control
		for i := range groupSize {
			c := uint8(ctrls)
			ctrls >>= 8
			if c == deleted {
				deletes++
			} else if c == empty {
				hasEmpty = true
			} else {
				entries++
				k := grp.keys[i]
				h := m.help.Hash(k)
				h2 := uint8(h & 0x7f)
				assert.Msg("control").This(c).Is(h2 | 0x80)
				v, ok := m.Get(k)
				assert.Msg("Get", k).That(ok)
				assert.Msg("Get", k).This(v).Is(grp.vals[i])
			}
		}
	}
	assert.Msg("entries").This(entries).Is(int(m.count))
	assert.Msg("has empty").That(hasEmpty)
	growthLeft := len(m.groups)*loadFactor - (entries + deletes)
	assert.Msg("growthLeft").This(growthLeft).Is(int(m.growthLeft))
}

//-------------------------------------------------------------------

// NewMapInt returns a map for integer keys
// NOTE: this is intended for testing, to give a consistent result
func NewMapInt[K constraints.Integer, V any]() *Map[K, V, integer[K]] {
	return &Map[K, V, integer[K]]{help: integer[K]{}}
}

type integer[K constraints.Integer] struct{}

func (m integer[K]) Hash(k K) uint64 {
	x := uint64(k)
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	x = x ^ (x >> 31)
	return x
}
func (m integer[K]) Equal(x, y K) bool {
	return x == y
}

// NewMapCmpable returns a map for comparable keys
func NewMapCmpable[K comparable, V any]() *Map[K, V, Cmpable[K]] {
	return &Map[K, V, Cmpable[K]]{help: Cmpable[K]{maphash.MakeSeed()}}
}

type Cmpable[K comparable] struct{ seed maphash.Seed }

func (c Cmpable[K]) Hash(k K) uint64 {
	return maphash.Comparable(c.seed, k)
}
func (Cmpable[K]) Equal(x, y K) bool {
	return x == y
}
