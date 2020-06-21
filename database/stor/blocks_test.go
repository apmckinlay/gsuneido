// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestBlocks(t *testing.T) {
	st := HeapStor(1024)
	w := st.Writer(32)
	w.Put1(0x11)
	w.Put2(0x2222)
	w.Put3(0x333333)
	w.Put4(0x44444444)
	pos := w.Pos()
	w.Put5(0x5555555555)
	w.PutStr("hello world")
	w.PutInts([]int{1,2,3})
	off := w.Close()

	r := st.Reader(off)
	Assert(t).That(r.Get1(), Equals(0x11))
	Assert(t).That(r.Get2(), Equals(0x2222))
	Assert(t).That(r.Get3(), Equals(0x333333))
	Assert(t).That(r.Get4(), Equals(0x44444444))
	Assert(t).That(r.Get5(), Equals(0x5555555555))
	Assert(t).That(r.GetStr(), Equals("hello world"))
	Assert(t).That(r.GetInts(), Equals([]int{1,2,3}))
	r.Pos(pos)
	Assert(t).That(r.Get5(), Equals(0x5555555555))
}
