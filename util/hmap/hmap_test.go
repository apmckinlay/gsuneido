// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hmap

import (
	"math/rand"
	"sort"
	"testing"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

type ik int

func (x ik) Hash() uint32 {
	return uint32(x)
}

func (x ik) Equal(y interface{}) bool {
	return x == y.(ik)
}

func TestHmap(t *testing.T) {
	data := map[int]string{}
	hmap := Hmap{}
	Assert(t).That(unsafe.Sizeof(hmap), Is(uintptr(32)))
	check := func() {
		check(t, data, &hmap)
	}
	put := func(n int, s string) {
		t.Helper()
		hmap.Put(ik(n), s)
		data[n] = s
		check()
	}
	del := func(n int) {
		t.Helper()
		var v Val
		var ok bool
		v, ok = data[n]
		if !ok {
			v = nil
		}
		Assert(t).That(hmap.Del(ik(n)), Is(v))
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

func check(t *testing.T, data map[int]string, hmap *Hmap) {
	t.Helper()
	//fmt.Println(data)
	//fmt.Println(hmap)
	Assert(t).That(hmap.Size(), Is(len(data)).Comment("size"))
	for i := 0; i < 100; i++ {
		s, ok := data[i]
		if ok {
			Assert(t).That(hmap.Get(ik(i)), Is(s))
		} else {
			Assert(t).That(hmap.Get(ik(i)), Is(nil))
		}
	}
}

func TestHmap_full(*testing.T) {
	h := Hmap{}
	for i := 0; i < 80; i++ {
		h.Put(ik(i), nil)
	}
}

func TestHmap_random(t *testing.T) {
	const N = 10000
	hm := Hmap{}
	Assert(t).That(hm.Size(), Is(0))
	nums := map[int32]int{}
	for i := 0; i < N; i++ {
		n := rand.Int31n(N)
		hm.Put(ik(n), i)
		nums[n] = i
	}
	rand.Seed(1)
	for i := 0; i < N; i++ {
		n := rand.Int31n(N)
		Assert(t).That(hm.Get(ik(n)), Is(nums[n]))
	}
	rand.Seed(1)
	for i := 0; i < N; i++ {
		n := rand.Int31n(N)
		v := hm.Del(ik(n))
		if nums[n] == -1 {
			Assert(t).That(v, Is(nil))
		} else {
			Assert(t).That(v, Is(nums[n]))
		}
		nums[n] = -1
	}
	Assert(t).That(hm.Size(), Is(0))
}

func TestHmap_Copy(t *testing.T) {
	h1 := Hmap{}
	h2 := h1.Copy()
	Assert(t).That(h2.Size(), Is(0))
	h1.Put(ik(123), "foo")
	Assert(t).That(h1.Size(), Is(1))
	Assert(t).That(h2.Size(), Is(0))
	h2 = h1.Copy()
	Assert(t).That(h2.Size(), Is(1))
	Assert(t).That(h2.Get(ik(123)), Is("foo"))
	h1.Put(ik(123), "bar")
	Assert(t).That(h1.Get(ik(123)), Is("bar"))
	Assert(t).That(h2.Get(ik(123)), Is("foo"))
}

func TestHmap_Iter(t *testing.T) {
	hm := Hmap{}
	test := func(n int) {
		it := hm.Iter()
		var nums []int
		for {
			k, v := it()
			if k == nil {
				Assert(t).That(v, Is(nil))
				break
			}
			ki := int(k.(ik))
			Assert(t).That(v, Is(-ki))
			nums = append(nums, ki)
		}
		Assert(t).That(len(nums), Is(n))
		sort.Ints(nums)
		for i := 0; i < n; i++ {
			Assert(t).That(nums[i], Is(i))
		}
	}
	test(0)
	hm.Put(ik(0), 0)
	test(1)
	hm.Put(ik(1), -1)
	test(2)
	for i := 2; i < 50; i++ {
		hm.Put(ik(i), -i)
		test(i + 1)
	}
}

func TestHmap_Iter_modified(t *testing.T) {
	hm := Hmap{}
	it := hm.Iter()
	hm.Put(ik(123), "foo")
	Assert(t).That(func() { it() }, Panics("hmap modified during iteration"))
	it = hm.Iter()
	hm.Del(ik(999)) // non-existent
	it()            // shouldn't panic
}

func BenchmarkHmap_Get(b *testing.B) {
	hm := &Hmap{}
	for i := 0; i < 100; i++ {
		hm.Put(ik(i), i)
	}
	for n := 0; n < b.N; n++ {
		hm.Get(ik(n % 100))
	}
}

func BenchmarkHmap_Put(b *testing.B) {
	for n := 0; n < b.N; n++ {
		hm := &Hmap{}
		for i := 0; i < 100; i++ {
			hm.Put(mix(i), i)
		}
	}
}

func mix(n int) ik {
	n = ^n + (n << 15)
	n = n ^ (n >> 12)
	n = n + (n << 2)
	n = n ^ (n >> 4)
	n = n * 2057
	n = n ^ (n >> 16)
	return ik(n)
}

func BenchmarkHmap_chainIter(b *testing.B) {
	h := &Hmap{}
	h.grow()
	var k Key = ik(123)
	var iter chainIter
	for n := 0; n < b.N; n++ {
		iter = h.iterFromKey(k)
	}
	iter.next()
}
