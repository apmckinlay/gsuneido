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
	s := "hello world"
	b := []byte("hello world")
	assert.T(t).That(String(s) == Bytes(b))
}

var Sum = uint32(0)
var S = "now is the time"

func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sum += String(S)
	}
}

func BenchmarkHashString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sum += HashString(S)
	}
}

func BenchmarkMaphash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := maphash.Hash{}
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
