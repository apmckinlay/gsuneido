// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hamttest

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/verify"
)

func TestRandom(t *testing.T) {
	ht := FooHamt{}.Mutable()
	_, ok := ht.Get(123)
	Assert(t).False(ok)
	var n = 100000
	if testing.Short() {
		n = 1000
	}
	rand.Seed(123456)
	for i := 0; i < n; i++ {
		f := int(rand.Int31())
		ht.Put(&Foo{f, strconv.Itoa(f)})
		if i%100 == 0 {
			ht = ht.Freeze().Mutable()
		}
	}
	rand.Seed(123456)
	for i := 0; i < n; i++ {
		f := int(rand.Int31())
		ht.Put(&Foo{f, strconv.Itoa(f)})
	}
	nums := map[int]bool{}
	ht = ht.Freeze()
	rand.Seed(123456)
	for i := 0; i < n; i++ {
		f := int(rand.Int31())
		foo, ok := ht.Get(f)
		Assert(t).True(ok)
		Assert(t).That(foo.key, Is(f))
		Assert(t).That(foo.data, Is(strconv.Itoa(f)))
		nums[f] = true
	}

	ht = ht.Mutable()
	// ht.print()
	for f := range nums {
		// fmt.Println("======================= del", f)
		verify.That(ht.Delete(f))
		// ht.print()
	}
	ht.ForEach(func(*Foo) { panic("should be empty") })
}

func TestPersistent(t *testing.T) {
	var ht FooHamt
	Assert(t).That(ht.string(), Is("{}"))
	h2 := ht.Mutable()
	h2.Put(&Foo{12, "12"})
	h2.Put(&Foo{34, "34"})
	h2 = h2.Freeze()
	Assert(t).That(ht.string(), Is("{}"))
	Assert(t).That(h2.string(), Is("{12,34}"))
	h3 := h2.Mutable()
	Assert(t).That(h3.string(), Is("{12,34}"))
	h3.Put(&Foo{56, "56"})
	h3.Put(&Foo{78, "78"})
	h3 = h3.Freeze()
	Assert(t).That(ht.string(), Is("{}"))
	Assert(t).That(h2.string(), Is("{12,34}"))
	Assert(t).That(h3.string(), Is("{12,34,56,78}"))
}

func (ht FooHamt) string() string {
	var list []string
	ht.ForEach(func(f *Foo) {
		list = append(list, f.data)
	})
	sort.Strings(list)
	return str.Join("{,}", list...)
}

func TestDelete(*testing.T) {
	data := []int{
		0, 2, 4, 6, 8, 16, 30, // should all go in root
		32, 34, 62, // collisions => child nodes
		0x10000, 0x20000, 0x30000, // same hash => overflow
	}
	ht := FooHamt{}.Mutable()
	for _, d := range data {
		// fmt.Printf("------------------------------ put %#x\n", d)
		ht.Put(&Foo{key: d})
		// ht.print()
	}
	ht = ht.Freeze()
	const nShuffles = 1000
	for i := 0; i < nShuffles; i++ {
		dht := ht.Mutable()
		rand.Shuffle(len(data), func(i, j int) { data[i], data[j] = data[j], data[i] })
		for i, d := range data {
			// fmt.Printf("------------------------------ del %#x\n", d)
			verify.That(dht.Delete(d))
			// dht.print()
			for j, d := range data {
				x, ok := dht.Get(data[j])
				verify.That(ok == (j > i))
				if ok {
					verify.That(x.key == d)
				}
			}
		}
		// fmt.Println("SHUFFLE =============================")
	}
}

func (ht FooHamt) print() {
	ht.root.print1(0)
}

func (nd *nodeFoo) print1(depth int) {
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
			p.print1(depth + 1)
		}
	}
}
