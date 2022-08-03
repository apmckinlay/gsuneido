// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"math/rand"
	"sort"
	"testing"
	"unsafe"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestHmap(t *testing.T) {
	data := map[int]string{}
	hmap := Hmap{}
	assert.T(t).This(unsafe.Sizeof(hmap)).Is(uintptr(32))
	check := func() {
		check(t, data, &hmap)
	}
	put := func(n int, s string) {
		t.Helper()
		hmap.Put(SuInt(n), SuStr(s))
		data[n] = s
		check()
	}
	del := func(n int) {
		t.Helper()
		s, ok := data[n]
		v := Value(SuStr(s))
		if !ok {
			v = nil
		}
		assert.T(t).This(hmap.Del(SuInt(n))).Is(v)
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
	assert := assert.T(t).This
	t.Helper()
	//fmt.Println(data)
	//fmt.Println(hmap)
	assert(hmap.Size()).Msg("size").Is(len(data))
	for i := 0; i < 100; i++ {
		s, ok := data[i]
		if ok {
			assert(hmap.Get(SuInt(i))).Is(SuStr(s))
		} else {
			assert(hmap.Get(SuInt(i))).Is(nil)
		}
	}
}

func TestHmap_full(*testing.T) {
	h := Hmap{}
	for i := 0; i < 80; i++ {
		h.Put(SuInt(i), nil)
	}
}

func TestHmap_random(t *testing.T) {
	assert := assert.T(t).This
	const N = 10000
	hm := Hmap{}
	assert(hm.Size()).Is(0)
	nums := map[int]int{}
	for i := 0; i < N; i++ {
		n := rand.Intn(N)
		hm.Put(SuInt(n), SuInt(i))
		nums[n] = i
	}
	rand.Seed(1)
	for i := 0; i < N; i++ {
		n := rand.Intn(N)
		assert(hm.Get(SuInt(n))).Is(SuInt(nums[n]))
	}
	rand.Seed(1)
	for i := 0; i < N; i++ {
		n := rand.Intn(N)
		v := hm.Del(SuInt(n))
		if nums[n] == -1 {
			assert(v).Is(nil)
		} else {
			assert(v).Is(SuInt(nums[n]))
		}
		nums[n] = -1
	}
	assert(hm.Size()).Is(0)
}

func TestHmap_Copy(t *testing.T) {
	assert := assert.T(t).This
	h1 := Hmap{}
	h2 := h1.Copy()
	assert(h2.Size()).Is(0)
	h1.Put(SuInt(123), SuStr("foo"))
	assert(h1.Size()).Is(1)
	assert(h2.Size()).Is(0)
	h2 = h1.Copy()
	assert(h2.Size()).Is(1)
	assert(h2.Get(SuInt(123))).Is(SuStr("foo"))
	h1.Put(SuInt(123), SuStr("bar"))
	assert(h1.Get(SuInt(123))).Is(SuStr("bar"))
	assert(h2.Get(SuInt(123))).Is(SuStr("foo"))
}

func TestHmap_Iter(t *testing.T) {
	assert := assert.T(t).This
	hm := Hmap{}
	test := func(n int) {
		it := hm.Iter()
		var nums []int
		for {
			k, v := it()
			if k == nil {
				assert(v).Is(nil)
				break
			}
			ki := ToInt(k)
			assert(v).Is(SuInt(-ki))
			nums = append(nums, ki)
		}
		assert(len(nums)).Is(n)
		sort.Ints(nums)
		for i := 0; i < n; i++ {
			assert(nums[i]).Is(i)
		}
	}
	test(0)
	hm.Put(SuInt(0), SuInt(0))
	test(1)
	hm.Put(SuInt(1), SuInt(-1))
	test(2)
	for i := 2; i < 50; i++ {
		hm.Put(SuInt(i), SuInt(-i))
		test(i + 1)
	}
}

func TestHmap_Iter_modified(t *testing.T) {
	hm := Hmap{}
	it := hm.Iter()
	hm.Put(SuInt(123), SuStr("foo"))
	assert.T(t).This(func() { it() }).Panics("hmap modified during iteration")
	it = hm.Iter()
	hm.Del(SuInt(999)) // non-existent
	it()               // shouldn't panic
}

type testVal struct {
	ValueBase[testVal]
	n int
}

func (sh *testVal) Equal(other interface{}) bool {
	return sh == other
}

func (sh *testVal) Hash() uint32 {
	return uint32(sh.n)
}

func TestHmapSameHash(t *testing.T) {
	const N = 100000
	hm := Hmap{}
	for i := 0; i < N; i++ {
		hm.Put(&testVal{n: i}, IntVal(i))
	}
}

func FuzzHmap(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		hm := Hmap{}
		for _, n := range b {
			hm.Put(&testVal{n: int(n)}, False)
		}
	})
}

// to run: go test -fuzz=FuzzHmap -run=FuzzHmap

//-------------------------------------------------------------------

func BenchmarkHmap_Get(b *testing.B) {
	hm := &Hmap{}
	for i := 0; i < 100; i++ {
		hm.Put(SuInt(i), SuInt(i))
	}
	for n := 0; n < b.N; n++ {
		hm.Get(SuInt(n % 100))
	}
}

func BenchmarkHmap_Put(b *testing.B) {
	for n := 0; n < b.N; n++ {
		hm := &Hmap{}
		for i := 0; i < 100; i++ {
			hm.Put(mix(i), SuInt(i))
		}
	}
}

func mix(n int) Value {
	n = ^n + (n << 15)
	n = n ^ (n >> 12)
	n = n + (n << 2)
	n = n ^ (n >> 4)
	n = n * 2057
	n = n ^ (n >> 16)
	return IntVal(n)
}

func BenchmarkHmap_chainIter(b *testing.B) {
	h := &Hmap{}
	h.grow()
	var k Value = SuInt(123)
	var iter chainIter
	for n := 0; n < b.N; n++ {
		iter = h.iterFromKey(k)
	}
	iter.next()
}
