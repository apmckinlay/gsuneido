package hmap

import (
	"math/rand"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

type ik int

func (x ik) Hash() uint32 {
	return uint32(x)
}

func (x ik) Equals(y interface{}) bool {
	return x == y.(ik)
}

func TestRandom(t *testing.T) {
	const N = 1000
	hm := NewHmap(0)
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
		found := hm.Del(ik(n))
		Assert(t).That(found, Equals(nums[n] != -1))
		nums[n] = -1
	}
	Assert(t).That(hm.Size(), Equals(0))
}
