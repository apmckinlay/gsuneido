package hmap2

import (
	"fmt"
	"math/rand"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func (h *Hmap) String() string {
	s := ""
	for _, b := range h.blocks {
		for ib := 0; ib < blockSize; ib++ {
			meta := b.meta[ib]
			if meta == metaEmpty {
				s += "(-)"
			} else {
				s += "("
				if meta&metaDirect == metaDirect {
					s += "*"
				}
				s += fmt.Sprint(b.key[ib]) //+ " " + fmt.Sprint(b.val[ib])
				if meta.jump() != noJump {
					s += " " + fmt.Sprint(jumpSize[meta.jump()])
				}
				s += ")"
			}
		}
	}
	return s
}

type ik int

func (x ik) Hash() uint32 {
	return uint32(x)
}

func (x ik) Equal(y interface{}) bool {
	return x == y.(ik)
}

func TestHmap(t *testing.T) {
	data := map[int]string{}
	hmap := &Hmap{}
	check := func() {
		check(t, data, hmap)
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
		Assert(t).That(hmap.Del(ik(n)), Equals(v))
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
	Assert(t).That(hmap.Size(), Equals(len(data)).Comment("size"))
	for i := 0; i < 100; i++ {
		s, ok := data[i]
		if ok {
			Assert(t).That(hmap.Get(ik(i)), Equals(s))
		} else {
			Assert(t).That(hmap.Get(ik(i)), Equals(nil))
		}
	}
}

func TestHmapFull(*testing.T) {
	h := &Hmap{}
	for i := 0; i < 80; i++ {
		h.Put(ik(i), nil)
	}
}

func TestRandom(t *testing.T) {
	const N = 10000
	hm := &Hmap{}
	Assert(t).That(hm.Size(), Equals(0))
	nums := map[int32]int{}
	for i := 0; i < N; i++ {
		n := rand.Int31n(N)
		hm.Put(ik(n), i)
		nums[n] = i
	}
	rand.Seed(1)
	for i := 0; i < N; i++ {
		n := rand.Int31n(N)
		Assert(t).That(hm.Get(ik(n)), Equals(nums[n]))
	}
	rand.Seed(1)
	for i := 0; i < N; i++ {
		n := rand.Int31n(N)
		v := hm.Del(ik(n))
		if nums[n] == -1 {
			Assert(t).That(v, Equals(nil))
		} else {
			Assert(t).That(v, Equals(nums[n]))
		}
		nums[n] = -1
	}
	Assert(t).That(hm.Size(), Equals(0))
}

func BenchmarkGet(b *testing.B) {
	hm := &Hmap{}
	for i := 0; i < 100; i++ {
		hm.Put(ik(i), i)
	}
	for n := 0; n < b.N; n++ {
		hm.Get(ik(n % 100))
	}
}

func BenchmarkPut(b *testing.B) {
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

func BenchmarkIter(b *testing.B) {
	h := &Hmap{}
	h.grow()
	var k Key = ik(123)
	var iter chainIter
	for n := 0; n < b.N; n++ {
		iter = h.iterFromKey(k)
	}
	iter.next()
}
