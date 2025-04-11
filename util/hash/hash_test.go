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

var Sum = uint64(0)
var S = "now is the time"

func BenchmarkString(b *testing.B) {
	for b.Loop() {
		Sum += String(S)
	}
}

func BenchmarkHashString(b *testing.B) {
	for b.Loop() {
		Sum += uint64(HashString(S))
	}
}

func BenchmarkMaphash(b *testing.B) {
	for b.Loop() {
		h := maphash.Hash{}
		h.WriteString(S)
		Sum += h.Sum64()
	}
}

func BenchmarkFnv(b *testing.B) {
	h := fnv.New64()
	for b.Loop() {
		h.Reset()
		h.Write([]byte(S))
		Sum += h.Sum64()
	}
}
