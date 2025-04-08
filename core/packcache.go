// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"bytes"

	"github.com/apmckinlay/gsuneido/util/shmap"
)

type packCache struct {
	lru     []uint8 // uint8 means max size of 256
	entries []entry
	hm      shmap.Map[Value, uint8, shmap.Meth[Value]]
}

type entry struct {
	key Value
	val int
}

const packCacheSize = 220 // ~ 7/8 * 256

func (lc *packCache) Get(key Value) (int, bool) {
	ei, ok := lc.hm.Get(key)
	if !ok {
		// not in cache
		return 0, false
	}
	// in cache
	li := bytes.IndexByte(lc.lru, uint8(ei))
	// don't move it if it's already near the end
	if li < packCacheSize-packCacheSize/4 {
		// move to the newest (the end)
		copy(lc.lru[li:], lc.lru[li+1:])
		lc.lru[len(lc.lru)-1] = uint8(ei)
	}
	return lc.entries[ei].val, true
}

func (lc *packCache) Put(key Value, val int) {
	ei := len(lc.entries)
	if ei < packCacheSize {
		lc.entries = append(lc.entries, entry{key: key, val: val})
		lc.lru = append(lc.lru, uint8(ei))
	} else { // full
		// replace oldest entry lru[0]
		ei = int(lc.lru[0])
		lc.hm.Del(lc.entries[ei].key)
		lc.entries[ei] = entry{key: key, val: val}
		copy(lc.lru, lc.lru[1:])
		lc.lru[packCacheSize-1] = uint8(ei)
	}
	lc.hm.Put(key, uint8(ei))
}
