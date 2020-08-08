// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestBlocks(t *testing.T) {
	buf := make([]byte, 64)
	w := NewWriter(buf)
	w.Put1(0x11)
	w.Put2(0x2222)
	w.Put3(0x333333)
	w.Put4(0x44444444)
	pos := w.Len()
	w.Put5(0x5555555555)
	w.PutStr("hello world")
	w.PutInts([]int{1, 2, 3})

	r := NewReader(buf)
	Assert(t).That(r.Get1(), Is(0x11))
	Assert(t).That(r.Get2(), Is(0x2222))
	Assert(t).That(r.Get3(), Is(0x333333))
	Assert(t).That(r.Get4(), Is(0x44444444))
	Assert(t).That(r.Get5(), Is(0x5555555555))
	Assert(t).That(r.GetStr(), Is("hello world"))
	Assert(t).That(r.GetInts(), Is([]int{1, 2, 3}))
	r = NewReader(buf[pos:])
	Assert(t).That(r.Get5(), Is(0x5555555555))
}

func BenchmarkBlocksCopy(b *testing.B) {
	buf := make([]byte, 50)
	for i := 0; i < b.N; i++ {
		w := write1{buf}
		for j := 0; j < 10; j++ {
			w.put5copy(i)
		}
	}
}

type write1 struct {
	buf []byte
}

func (w *write1) put5copy(n int) {
	copy(w.buf, []byte{
		byte(n),
		byte(n >> 8),
		byte(n >> 16),
		byte(n >> 24),
		byte(n >> 32)})
	w.buf = w.buf[5:]
}

func BenchmarkBlocksCopy2(b *testing.B) {
	buf := make([]byte, 50)
	for i := 0; i < b.N; i++ {
		w := write2{buf: buf}
		for j := 0; j < 10; j++ {
			w.put5copy(i)
		}
	}
}

type write2 struct {
	buf []byte
	i   int
}

func (w *write2) put5copy(n int) {
	copy(w.buf[w.i:], []byte{
		byte(n),
		byte(n >> 8),
		byte(n >> 16),
		byte(n >> 24),
		byte(n >> 32)})
	w.i += 5
}

func BenchmarkBlocksAppend(b *testing.B) {
	buf := make([]byte, 50)
	for i := 0; i < b.N; i++ {
		w := write1{buf[:0]}
		for j := 0; j < 10; j++ {
			w.put5append(i)
		}
	}
}

func (w *write1) put5append(n int) {
	w.buf = append(w.buf,
		byte(n),
		byte(n>>8),
		byte(n>>16),
		byte(n>>24),
		byte(n>>32))
}
