// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hamt

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Foo struct {
	key  int
	data string
}

func (f *Foo) Key() int {
	return f.key
}

func (*Foo) Hash(key int) uint64 {
	return uint64(key) & 0xffff // reduce bits to force overflows

}

func (f *Foo) StorSize() int {
	return 0
}

func (f *Foo) Cksum() uint32 {
	return 0
}

func (f *Foo) IsTomb() bool {
	return false
}

func (f *Foo) LastMod() int {
	return 0
}

func (f *Foo) SetLastMod(int) {
}

func (f *Foo) Write(*stor.Writer) {
}

func (f *Foo) Read(*stor.Stor, *stor.Reader) {
}

type FooHamt = Hamt[int, *Foo]
type FooNode = node[int, *Foo]

func TestRandom(t *testing.T) {
	assert := assert.T(t)
	ht := FooHamt{}.Mutable()
	_, ok := ht.Get(123)
	assert.False(ok)
	var n = 100000
	if testing.Short() {
		n = 1000
	}
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	for i := range n {
		f := int(r.Int31())
		ht.Put(&Foo{f, strconv.Itoa(f)})
		// ht.check()
		if i%100 == 0 {
			ht = ht.Freeze().Mutable()
		}
	}
	r.Seed(seed)
	for range n {
		f := int(r.Int31())
		ht.Put(&Foo{f, strconv.Itoa(f)})
		// ht.check()
	}
	nums := map[int]bool{}
	ht = ht.Freeze()
	r.Seed(seed)
	for range n {
		f := int(r.Int31())
		foo, ok := ht.Get(f)
		assert.True(ok)
		assert.This(foo.key).Is(f)
		assert.This(foo.data).Is(strconv.Itoa(f))
		nums[f] = true
	}

	ht = ht.Mutable()
	// ht.print()
	for f := range nums {
		// fmt.Println("======================= del", f)
		assert.That(ht.Delete(f))
		// ht.check()
		// ht.print()
	}
	ht.ForEach(func(*Foo) { panic("should be empty") })
}

func TestPersistent(t *testing.T) {
	assert := assert.T(t).This
	var ht FooHamt
	assert(tostr(ht)).Is("{}")
	h2 := ht.Mutable()
	h2.Put(&Foo{12, "12"})
	h2.Put(&Foo{34, "34"})
	h2 = h2.Freeze()
	assert(tostr(ht)).Is("{}")
	assert(tostr(h2)).Is("{12,34}")
	h3 := h2.Mutable()
	assert(tostr(h3)).Is("{12,34}")
	h3.Put(&Foo{56, "56"})
	h3.Put(&Foo{78, "78"})
	h3 = h3.Freeze()
	assert(tostr(ht)).Is("{}")
	assert(tostr(h2)).Is("{12,34}")
	assert(tostr(h3)).Is("{12,34,56,78}")
}

func tostr(ht FooHamt) string {
	var list []string
	ht.ForEach(func(f *Foo) {
		list = append(list, f.data)
	})
	sort.Strings(list)
	return str.Join("{,}", list)
}

func TestDelete(*testing.T) {
	data := []int{
		0, 2, 4, 6, 8, 16, 30, // should all go in root
		32, 34, 62, // collisions => child nodes
		0x10000, 0x20000, 0x30000, 0x40000, // same hash => child nodes
		0x50000, 0x60000, 0x70000, 0x80000, 0x90000, // => overflow
	}
	nShuffles := 10000
	if testing.Short() {
		nShuffles = 1000
	}
	for range nShuffles {
		rand.Shuffle(len(data), func(i, j int) { data[i], data[j] = data[j], data[i] })
		ht := FooHamt{}.Mutable()
		for _, d := range data {
			// fmt.Printf("------------------------------ put %#x\n", d)
			ht.Put(&Foo{key: d})
			// ht.print()
		}
		for i, d := range data {
			// fmt.Printf("------------------------------ del %#x\n", d)
			// fmt.Printf("delete %#x\n", d)
			assert.That(ht.Delete(d))
			// dht.print()
			for j, d := range data {
				x, ok := ht.Get(data[j])
				assert.That(ok == (j > i))
				if ok {
					assert.That(x.key == d)
				}
			}
		}
		// fmt.Println("SHUFFLE =============================")
	}
}

func TestDeleteInsertBug(*testing.T) {
	ht := FooHamt{}.Mutable()
	ht.Put(&Foo{key: 0x10000})
	ht.Put(&Foo{key: 0x20000}) // same hash, collision goes in child
	ht.Put(&Foo{key: 0x30000}) // same hash, collision goes in child
	ht.Put(&Foo{key: 0x40000}) // same hash, collision goes in child
	// print(ht)
	ht.Delete(0x10000) // will pull up 0x40000
	// print(ht)
	ht.Put(&Foo{key: 0x20000})
	// print(ht)
	check(ht)
}

//lint:ignore U1000 for debugging
func print(ht FooHamt) {
	print1(0, ht.root)
}

//lint:ignore U1000 for debugging
func print1(depth int, nd *FooNode) {
	indent := strings.Repeat("    ", depth)

	if depth > 6 {
		fmt.Print(indent + "overflow")
		for i := range nd.vals {
			fmt.Printf(" %#x", nd.vals[i].key)
		}
		fmt.Println()
		return
	}

	if nd.bmVal != 0 {
		fmt.Printf(indent+"vals %032b ", nd.bmVal)
		for i := range nd.vals {
			fmt.Printf("%#x ", nd.vals[i].key)
		}
		fmt.Println()
	}
	if nd.bmPtr != 0 {
		fmt.Printf(indent+"ptrs %032b\n", nd.bmPtr)
		for _, p := range nd.ptrs {
			print1(depth+1, p)
		}
	}
}

func check(ht FooHamt) {
	keys := make(map[int]bool)
	ht.ForEach(func(foo *Foo) {
		if _, ok := keys[foo.key]; ok {
			panic("duplicate key")
		}
		keys[foo.key] = true
	})
}
