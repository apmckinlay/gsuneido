// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package comp

import (
	"math/rand"
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/runtime"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestKey(t *testing.T) {
	Assert(t).That(Key(mkrec("a", "b"), []int{}, false), Is(""))
	Assert(t).That(Key(mkrec("a", "b"), []int{0}, false), Is("a"))
	Assert(t).That(Key(mkrec("a", "b"), []int{1}, false), Is("b"))
	Assert(t).That(Key(mkrec("a", "b"), []int{0, 1}, false), Is("a\x00\x00b"))
	Assert(t).That(Key(mkrec("a", "b"), []int{1, 0}, false), Is("b\x00\x00a"))

	// omit trailing empty fields
	fields := []int{0, 1, 2}
	Assert(t).That(Key(mkrec("a", "b", "c"), fields, false),
		Is("a\x00\x00b\x00\x00c"))
	Assert(t).That(Key(mkrec("a", "", "c"), fields, false),
		Is("a\x00\x00\x00\x00c"))
	Assert(t).That(Key(mkrec("", "", "c"), fields, false),
		Is("\x00\x00\x00\x00c"))
	Assert(t).That(Key(mkrec("a", "b", ""), fields, false),
		Is("a\x00\x00b"))
	Assert(t).That(Key(mkrec("a", "", ""), fields, false),
		Is("a"))
	Assert(t).That(Key(mkrec("", "", ""), fields, false),
		Is(""))

	// no escape for single field
	Assert(t).That(Key(mkrec("a\x00b"), []int{0}, false), Is("a\x00b"))

	// escaping
	first := []int{0,1}
	Assert(t).That(Key(mkrec("ab"), first, false), Is("ab"))
	Assert(t).That(Key(mkrec("a\x00b"), first, false), Is("a\x00\x01b"))
	Assert(t).That(Key(mkrec("\x00ab"), first, false), Is("\x00\x01ab"))
	Assert(t).That(Key(mkrec("a\x00\x00b"), first, false), Is("a\x00\x01\x00\x01b"))
	Assert(t).That(Key(mkrec("a\x00\x01b"), first, false), Is("a\x00\x01\x01b"))
	Assert(t).That(Key(mkrec("ab\x00"), first, false), Is("ab\x00\x01"))
	Assert(t).That(Key(mkrec("ab\x00\x00"), first, false), Is("ab\x00\x01\x00\x01"))

	// ts
	Assert(t).That(Key(mkrec("a", "b"), []int{0, 1}, false), Is("a\x00\x00b"))
	Assert(t).That(Key(mkrec("a", "b"), []int{0, 1}, true), Is("a"))
	Assert(t).That(Key(mkrec("", "b"), []int{0, 1}, true), Is("\x00\x00b"))

}

func mkrec(args ...string) Record {
	var b RecordBuilder
	for _, a := range args {
		b.AddRaw(a)
	}
	return b.Build()
}

const m = 3

func TestRandom(t *testing.T) {
	var n = 100000
	if testing.Short() {
		n = 10000
	}
	fields := []int{0, 1, 2}
	for i := 0; i < n; i++ {
		x := gen()
		y := gen()
		yenc := Key(y, fields, false)
		xenc := Key(x, fields, false)
		Assert(t).That(xenc < yenc, Is(lt(x, y)))
		Assert(t).That(strings.Compare(xenc, yenc), Is(Compare(x, y, fields)))
	}
}

func gen() Record {
	var b RecordBuilder
	for i := 0; i < m; i++ {
		x := make([]byte, rand.Intn(6)+1)
		for j := range x {
			x[j] = byte(rand.Intn(4)) // 25% zeros
		}
		b.AddRaw(string(x))
	}
	return b.Build()
}

func lt(x Record, y Record) bool {
	for i := 0; i < x.Len() && i < y.Len(); i++ {
		if cmp := strings.Compare(x.GetRaw(i), y.GetRaw(i)); cmp != 0 {
			return cmp < 0
		}
	}
	return x.Len() < y.Len()
}
