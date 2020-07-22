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
	Assert(t).That(Key(mkrec("a", "b"), []int{}), Equals(""))
	Assert(t).That(Key(mkrec("a", "b"), []int{0}), Equals("a"))
	Assert(t).That(Key(mkrec("a", "b"), []int{1}), Equals("b"))
	Assert(t).That(Key(mkrec("a", "b"), []int{0, 1}), Equals("a\x00\x00b"))
	Assert(t).That(Key(mkrec("a", "b"), []int{1, 0}), Equals("b\x00\x00a"))

	first := []int{0}
	Assert(t).That(Key(mkrec("ab"), first), Equals("ab"))
	Assert(t).That(Key(mkrec("a\x00b"), first), Equals("a\x00\x01b"))
	Assert(t).That(Key(mkrec("\x00ab"), first), Equals("\x00\x01ab"))
	Assert(t).That(Key(mkrec("a\x00\x00b"), first), Equals("a\x00\x01\x00\x01b"))
	Assert(t).That(Key(mkrec("a\x00\x01b"), first), Equals("a\x00\x01\x01b"))
	Assert(t).That(Key(mkrec("ab\x00"), first), Equals("ab\x00\x01"))
	Assert(t).That(Key(mkrec("ab\x00\x00"), first), Equals("ab\x00\x01\x00\x01"))
}

func mkrec(args ...string) Record {
	var b RecordBuilder
	for _, a := range args {
		b.AddRaw(a)
	}
	return b.Build()
}

const m = 3

func TestEncodeRandom(t *testing.T) {
	var n = 100000
	if testing.Short() {
		n = 10000
	}
	fields := []int{0, 1, 2}
	for i := 0; i < n; i++ {
		x := gen()
		y := gen()
		yenc := Key(y, fields)
		xenc := Key(x, fields)
		Assert(t).That(xenc < yenc, Equals(lt(x, y)))
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
