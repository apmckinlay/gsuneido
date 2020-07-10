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
	var ht FooHamt
	Assert(t).That(ht.Get(123), Equals(nil))
	var n = 100000
	if testing.Short() {
		n = 1000
	}
	hu := ht.Update()
	rand.Seed(123456)
	for i := 0; i < n; i++ {
		f := int(rand.Int31())
		hu.Put(&Foo{f, strconv.Itoa(f)})
	}
	rand.Seed(123456)
	for i := 0; i < n; i++ {
		f := int(rand.Int31())
		hu.Put(&Foo{f, strconv.Itoa(f)})
	}
	ht = hu.Finish()
	rand.Seed(123456)
	for i := 0; i < n; i++ {
		f := int(rand.Int31())
		foo := ht.Get(f)
		Assert(t).That(foo, NotEquals(nil))
		Assert(t).That(foo.key, Equals(f))
		Assert(t).That(foo.data, Equals(strconv.Itoa(f)))
	}
}

func TestPersistent(t *testing.T) {
	var ht FooHamt
	Assert(t).That(ht.string(), Equals("{}"))
	ht2 := ht.Update().Put(&Foo{12, "12"}).Put(&Foo{34, "34"}).Finish()
	Assert(t).That(ht.string(), Equals("{}"))
	Assert(t).That(ht2.string(), Equals("{12,34}"))
	ht3 := ht2.Update().Put(&Foo{56, "56"}).Put(&Foo{78, "78"}).Finish()
	Assert(t).That(ht.string(), Equals("{}"))
	Assert(t).That(ht2.string(), Equals("{12,34}"))
	Assert(t).That(ht3.string(), Equals("{12,34,56,78}"))
}

func (ht FooHamt) string() string {
	var list []string
	ht.ForEach(func(f *Foo) {
		list = append(list, f.data)
	})
	sort.Strings(list)
	return "{" + str.FromList(list) + "}"
}
