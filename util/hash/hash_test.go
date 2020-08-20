// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hash

import (
	"hash/fnv"
	"hash/maphash"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestHash(t *testing.T) {
	test := func(s string, expected uint32) {
		t.Helper()
		assert.T(t).This(HashString(s)).Is(expected)
		assert.T(t).This(HashBytes([]byte(s))).Is(expected)
	}
	test("", 0x811c9dc5)
	test("foobar", 0xbf9cf968)
}

var Sum = uint32(0)
var S = "now is the time for all good"

func BenchmarkHashString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sum += HashString(S)
	}
}

func BenchmarkMaphash(b *testing.B) {
	h := maphash.Hash{}
	for i := 0; i < b.N; i++ {
		h.Reset()
		h.WriteString(S)
		Sum += uint32(h.Sum64())
	}
}

func BenchmarkFnv(b *testing.B) {
	h := fnv.New32()
	for i := 0; i < b.N; i++ {
		h.Reset()
		h.Write([]byte(S))
		Sum += h.Sum32()
	}
}
