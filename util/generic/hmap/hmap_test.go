// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hmap

import (
	"math/rand"
	"sort"
	"testing"
	"time"
	"unsafe"

	"github.com/apmckinlay/gsuneido/util/assert"
)

type intkey int

func (ik intkey) Hash() uint64 {
	return uint64(ik)
}

func (ik intkey) Equal(ik2 any) bool {
	return ik == ik2.(intkey)
}

type int2int = Hmap[intkey, int, Meth[intkey]]
type int2string = Hmap[intkey, string, Meth[intkey]]

func TestHmap(t *testing.T) {
	data := map[int]string{}
	hmap := int2string{}
	assert.T(t).This(unsafe.Sizeof(hmap)).Is(uintptr(32))
	check := func() {
		check(t, data, &hmap)
	}
	put := func(n int, s string) {
		t.Helper()
		hmap.Put(intkey(n), s)
		data[n] = s
		check()
	}
	del := func(n int) {
		t.Helper()
		s, ok := data[n]
		v := s
		if !ok {
			v = ""
		}
		assert.T(t).This(hmap.Del(intkey(n))).Is(v)
		delete(data, n)
		check()
	}

	check()

	del(33) // delete when empty

	// insert
	put(20, "twenty")
	put(30, "thirty")

	// update
	put(20, "Twenty")

	// collision, add to chain
	put(31, "thirtyone")
	put(32, "thirtytwo")

	// move chain
	put(40, "forty")

	// delete
	del(10) // nonexistent
	del(39) // nonexistent chain
	del(20) // simplest case
	del(20) // should be gone
	del(32) // end of chain
	put(32, "thirtytwo")
	del(30) // beginning of chain
	del(32) // beginning of chain

	// trigger grow
	for i := 0; i < 10; i++ { // long chain
		put(i, "x")
	}
	for i := 15; i < 100; i += 10 {
		put(i, "y")
	}
}

func check(t *testing.T, data map[int]string, hmap *int2string) {
	assert := assert.T(t).This
	t.Helper()
	//fmt.Println(data)
	//fmt.Println(hmap)
	assert(hmap.Size()).Msg("size").Is(len(data))
	for i := 0; i < 100; i++ {
		s, ok := data[i]
		if ok {
			assert(hmap.Get(intkey(i))).Is(s)
		} else {
			assert(hmap.Get(intkey(i))).Is("")
		}
	}
}

func TestHmap_full(*testing.T) {
	h := int2string{}
	for i := 0; i < 80; i++ {
		h.Put(intkey(i), "")
	}
}

func TestHmap_random(t *testing.T) {
	assert := assert.This
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	const N = 10000
	hm := int2int{}
	assert(hm.Size()).Is(0)
	nums := map[int]int{}
	for i := 0; i < N; i++ {
		n := r.Intn(N)
		hm.Put(intkey(n), i)
		nums[n] = i
	}
	r.Seed(seed)
	for i := 0; i < N; i++ {
		n := r.Intn(N)
		assert(hm.Get(intkey(n))).Is(nums[n])
	}
	r.Seed(seed)
	for i := 0; i < N; i++ {
		n := r.Intn(N)
		v := hm.Del(intkey(n))
		if nums[n] == -1 {
			assert(v).Is(0)
		} else {
			assert(v).Is(nums[n])
		}
		nums[n] = -1
	}
	assert(hm.Size()).Is(0)
}

func TestHmap_funcs(t *testing.T) {
	hm := NewHmapFuncs[int, int](func(i int) uint64 { return uint64(i) },
		func(x, y int) bool { return x == y })
	hm.Put(123, 456)
	assert.That(hm.Has(123))
	assert.This(hm.Get(123)).Is(456)
	assert.That(!hm.Has(999))
}

func TestHmap_Copy(t *testing.T) {
	assert := assert.T(t).This
	h1 := int2string{}
	h2 := h1.Copy()
	assert(h2.Size()).Is(0)
	h1.Put(intkey(123), "foo")
	assert(h1.Size()).Is(1)
	assert(h2.Size()).Is(0)
	h2 = h1.Copy()
	assert(h2.Size()).Is(1)
	assert(h2.Get(intkey(123))).Is("foo")
	h1.Put(intkey(123), "bar")
	assert(h1.Get(intkey(123))).Is("bar")
	assert(h2.Get(intkey(123))).Is("foo")
}

func TestHmap_Iter(t *testing.T) {
	assert := assert.T(t).This
	hm := int2int{}
	test := func(n int) {
		it := hm.Iter()
		var nums []int
		for {
			k, v := it()
			if k == 0 {
				assert(v).Is(0)
				break
			}
			ki := int(k)
			assert(v).Is(-ki)
			nums = append(nums, ki)
		}
		assert(len(nums)).Is(n)
		sort.Ints(nums)
		for i := 0; i < n; i++ {
			assert(nums[i]).Is(i + 1)
		}
	}
	test(0)
	hm.Put(intkey(1), -1)
	test(1)
	hm.Put(intkey(2), -2)
	test(2)
	for i := 3; i < 50; i++ {
		hm.Put(intkey(i), -i)
		test(i)
	}
}

func TestHmap_Iter_modified(t *testing.T) {
	hm := int2string{}
	it := hm.Iter()
	hm.Put(intkey(123), "foo")
	assert.T(t).This(func() { it() }).Panics("hmap modified during iteration")
	it = hm.Iter()
	hm.Del(intkey(999)) // non-existent
	it()                // shouldn't panic
}

type testKey struct {
	hash int
}

func (tk *testKey) Equal(other any) bool {
	return tk == other.(*testKey)
}

func (tk *testKey) Hash() uint64 {
	return uint64(tk.hash)
}

type tkType = Hmap[*testKey, int, Meth[*testKey]]

func TestHmapSameHash(t *testing.T) {
	const N = 100000
	hm := &tkType{}
	for i := 0; i < N; i++ {
		hm.Put(&testKey{i}, i)
	}
}

func FuzzHmap(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		hm := &tkType{}
		for _, n := range b {
			hm.Put(&testKey{int(n)}, int(n))
		}
	})
}

// to run: go test -fuzz=FuzzHmap -run=FuzzHmap

//-------------------------------------------------------------------

func BenchmarkHmap_Get(b *testing.B) {
	hm := int2int{}
	for i := 0; i < 100; i++ {
		hm.Put(intkey(i), i)
	}
	for n := 0; n < b.N; n++ {
		hm.Get(intkey(n % 100))
	}
}

func BenchmarkHmap_Put(b *testing.B) {
	for n := 0; n < b.N; n++ {
		hm := int2int{}
		for i := 0; i < 100; i++ {
			hm.Put(mix(i), i)
		}
	}
}

func mix(n int) intkey {
	n = ^n + (n << 15)
	n = n ^ (n >> 12)
	n = n + (n << 2)
	n = n ^ (n >> 4)
	n = n * 2057
	n = n ^ (n >> 16)
	return intkey(n)
}

func BenchmarkHmap_chainIter(b *testing.B) {
	h := int2int{}
	h.grow()
	k := intkey(123)
	var iter chainIter[intkey, int, Meth[intkey]]
	for n := 0; n < b.N; n++ {
		iter = h.iterFromKey(k)
	}
	iter.next()
}
