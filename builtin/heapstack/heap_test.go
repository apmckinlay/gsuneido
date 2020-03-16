// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package heapstack

import "testing"

func TestHeap(*testing.T) {
	f1()
}

func f1() {
	defer FreeTo(CurSize())
	p := Alloc(64)
	*(*[64]byte)(p) = [64]byte{}
	f2()
}

func f2() {
	defer FreeTo(CurSize())
	q := Alloc(8)
	*(*int64)(q) = 0
}
