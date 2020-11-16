// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package fbtree

import "github.com/apmckinlay/gsuneido/db19/stor"

// builder is used to bulk load an fbtree.
// Keys must be added in order.
// The fbtree is built bottom up with no splitting or inserting.
// All nodes will be "full" except for the right hand edge.
type builder struct {
	levels []*level // leaf is [0]
	prev   string
	store  *stor.Stor
	count  int
}

type level struct {
	splitKey string
	builder  fNodeBuilder
}

func Builder(store *stor.Stor) *builder {
	return &builder{store: store, levels: []*level{{}}}
}

func (fb *builder) Add(key string, off uint64) {
	if fb.count > 0 {
		if key == fb.prev {
			panic("fbtreeBuilder keys must not have duplicates")
		}
		if key < fb.prev {
			panic("fbtreeBuilder keys must be inserted in order")
		}
	}
	fb.add(0, key, off)
	fb.prev = key
	fb.count++
}

func (fb *builder) add(li int, key string, off uint64) {
	if li >= len(fb.levels) {
		fb.levels = append(fb.levels, &level{})
	}
	lev := fb.levels[li]
	if len(lev.builder.fe) > (MaxNodeSize * 3 / 4) {
		// split full node to stor
		offNode, splitKey := lev.builder.Split(fb.store)
		fb.add(li+1, lev.splitKey, offNode) // RECURSE
		lev.splitKey = splitKey
	}
	embedLen := 1
	if li > 0 /*|| fb.count == 1*/ {
		embedLen = 255
	}
	lev.builder.Add(key, off, embedLen)
}

func (fb *builder) Finish() *fbtree {
	var key string
	var off uint64
	for li := 0; li < len(fb.levels); li++ {
		if li > 0 {
			// allow node to slightly exceed max size
			fb.levels[li].builder.Add(key, off, 255)
		}
		key = fb.levels[li].splitKey
		off = fb.levels[li].builder.fe.putNode(fb.store)
	}
	treeLevels := len(fb.levels) - 1
	return OpenFbtree(fb.store, off, treeLevels)
}

//-------------------------------------------------------------------

type fNodeBuilder struct {
	fe       fnode
	notFirst bool
	fi       int
	prev     string
	known    string
	offset   uint64
	fi2      int
	known2   string
	offset2  uint64
}

func (fb *fNodeBuilder) Add(key string, offset uint64, embedLen int) {
	if fb.notFirst {
		if key <= fb.prev {
			panic("fBuilder keys must be inserted in order, without duplicates")
		}
	} else {
		fb.notFirst = true
	}
	if len(fb.fe) == 0 {
		fb.fe = fb.fe.append(offset, 0, "")
		fb.known = ""
	} else {
		fb.fi2 = fb.fi
		fb.fi = len(fb.fe)
		npre, diff, known := addone(key, fb.prev, fb.known, embedLen)
		fb.fe = fb.fe.append(offset, npre, diff)
		fb.known2 = fb.known
		fb.known = known
		fb.offset2 = fb.offset
		fb.offset = offset
	}
	fb.prev = key
}

// Split saves all but the last two entries as the left node
// and initializes fb.fe with the last two entries
func (fb *fNodeBuilder) Split(store *stor.Stor) (leftOff uint64, splitKey string) {
	splitKey = fb.known2 // known of second last entry
	left := fb.fe[:fb.fi2]
	leftOff = left.putNode(store)
	// first entry becomes 0, ""
	right := fb.fe[:0].append(fb.offset2, 0, "") // offset of second last entry
	// second entry becomes 0, known
	right = right.append(fb.offset, 0, fb.known) // offset,known of last entry
	fb.fe = right
	return
}

func (fb *fNodeBuilder) Entries() fnode {
	return fb.fe
}
