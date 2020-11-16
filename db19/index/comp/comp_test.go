// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package comp

import (
	"math/rand"
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestKey(t *testing.T) {
	assert := assert.T(t).This
	fields := []int{0, 1, 2}
	for _, flds2 := range [][]int{nil, {1, 2}} {
		assert(Key(mkrec("a", "b"), []int{}, flds2)).Is("")
		assert(Key(mkrec("a", "b"), []int{0}, flds2)).Is("a")
		assert(Key(mkrec("a", "b"), []int{1}, flds2)).Is("b")
		assert(Key(mkrec("a", "b"), []int{0, 1}, flds2)).Is("a\x00\x00b")
		assert(Key(mkrec("a", "b"), []int{1, 0}, flds2)).Is("b\x00\x00a")

		// omit trailing empty fields
		assert(Key(mkrec("a", "b", "c"), fields, flds2)).Is("a\x00\x00b\x00\x00c")
		assert(Key(mkrec("a", "", "c"), fields, flds2)).Is("a\x00\x00\x00\x00c")
		assert(Key(mkrec("", "", "c"), fields, flds2)).Is("\x00\x00\x00\x00c")
		assert(Key(mkrec("a", "b", ""), fields, flds2)).Is("a\x00\x00b")
		assert(Key(mkrec("a", "", ""), fields, flds2)).Is("a")

		// no escape for single field
		assert(Key(mkrec("a\x00b"), []int{0}, flds2)).Is("a\x00b")

		// escaping
		first := []int{0, 1}
		assert(Key(mkrec("ab"), first, flds2)).Is("ab")
		assert(Key(mkrec("a\x00b"), first, flds2)).Is("a\x00\x01b")
		assert(Key(mkrec("\x00ab"), first, flds2)).Is("\x00\x01ab")
		assert(Key(mkrec("a\x00\x00b"), first, flds2)).Is("a\x00\x01\x00\x01b")
		assert(Key(mkrec("a\x00\x01b"), first, flds2)).Is("a\x00\x01\x01b")
		assert(Key(mkrec("ab\x00"), first, flds2)).Is("ab\x00\x01")
		assert(Key(mkrec("ab\x00\x00"), first, flds2)).Is("ab\x00\x01\x00\x01")
	}

	// fields2
	fields2 := []int{3, 4}
	assert(Key(mkrec("", "", ""), fields, nil)).Is("")
	assert(Key(mkrec("", "", "", "a", "b"), fields, fields2)).
		Is("\x00\x00\x00\x00\x00\x00a\x00\x00b")
	assert(Key(mkrec("x", "", "", "a", "b"), fields, fields2)).
		Is("x")
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
	assert := assert.T(t).This
	var n = 100000
	if testing.Short() {
		n = 10000
	}
	fields := []int{0, 1, 2}
	for i := 0; i < n; i++ {
		x := gen()
		y := gen()
		yenc := Key(y, fields, nil)
		xenc := Key(x, fields, nil)
		assert(xenc < yenc).Is(lt(x, y))
		assert(strings.Compare(xenc, yenc)).Is(Compare(x, y, fields, nil))
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
