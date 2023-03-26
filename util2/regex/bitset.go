// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package regex

import "github.com/apmckinlay/gsuneido/util/generic/slc"

// BitSet is a set of integers stored in bits.
// Zero value is ready to use.
// Memory usages depends on the maximum integer value
// so it should only be used for small values.
type BitSet struct {
	bits []int64
}

func (bs *BitSet) Add(n int16) {
	i := int(n / 64)
	bs.bits = slc.Allow(bs.bits, i+1)
	bs.bits[i] |= 1 << uint(n%64)
}

func (bs *BitSet) Has(n int16) bool {
	i := int(n / 64)
	if i >= len(bs.bits) {
		return false
	}
	return (bs.bits[i]>>(n%64))&1 == 1
}

func (bs *BitSet) AddNew(n int16) bool {
	if bs.Has(n) {
		return false
	}
	bs.Add(n)
	return true
}

func (bs *BitSet) Clear() {
	bs.bits = bs.bits[:0]
}
