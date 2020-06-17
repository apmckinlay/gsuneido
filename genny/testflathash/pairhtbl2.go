// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package testflathash

type Pair struct {
	key int
	val int
}

func (*PairHtbl) hash(key int) uint32 {
	return uint32(key)
}

func (*PairHtbl) keyOf(item *Pair) int {
	return item.key
}
