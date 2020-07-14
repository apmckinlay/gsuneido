// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hamttest

import (
	"math/rand"
	"sort"
	"strconv"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestRandom(t *testing.T) {
	ht := FooHamt{}.Mutable()
	_, ok := ht.Get(123)
	Assert(t).False(ok)
	var n = 1000000
	if testing.Short() {
		n = 1000
	}
	rand.Seed(123456)
	for i := 0; i < n; i++ {
		f := int(rand.Int31())
		ht.Put(&Foo{f, strconv.Itoa(f)})
	}
	rand.Seed(123456)
	for i := 0; i < n; i++ {
		f := int(rand.Int31())
		ht.Put(&Foo{f, strconv.Itoa(f)})
	}
	ht = ht.Freeze()
	rand.Seed(123456)
	for i := 0; i < n; i++ {
		f := int(rand.Int31())
		foo, ok := ht.Get(f)
		Assert(t).True(ok)
		Assert(t).That(foo.key, Equals(f))
		Assert(t).That(foo.data, Equals(strconv.Itoa(f)))
	}
}

func TestPersistent(t *testing.T) {
	var ht FooHamt
	Assert(t).That(ht.string(), Equals("{}"))
	h2 := ht.Mutable()
	h2.Put(&Foo{12, "12"})
	h2.Put(&Foo{34, "34"})
	h2 = h2.Freeze()
	Assert(t).That(ht.string(), Equals("{}"))
	Assert(t).That(h2.string(), Equals("{12,34}"))
	h3 := h2.Mutable()
	Assert(t).That(h3.string(), Equals("{12,34}"))
	h3.Put(&Foo{56, "56"})
	h3.Put(&Foo{78, "78"})
	h3 = h3.Freeze()
	Assert(t).That(ht.string(), Equals("{}"))
	Assert(t).That(h2.string(), Equals("{12,34}"))
	Assert(t).That(h3.string(), Equals("{12,34,56,78}"))
}

func (ht FooHamt) string() string {
	var list []string
	ht.ForEach(func(f *Foo) {
		list = append(list, f.data)
	})
	sort.Strings(list)
	return "{" + str.FromList(list) + "}"
}
