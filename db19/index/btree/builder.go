// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import "github.com/apmckinlay/gsuneido/db19/stor"

// builder is used to bulk load an btree.
// Keys must be added in order.
// The btree is built bottom up with no splitting or inserting.
// All nodes will be "full" except for the right hand edge.
type builder struct {
	levels []*level // leaf is [0]
	prev   string
	stor   *stor.Stor
	count  int
}

type level struct {
	splitKey string
	nb       nodeBuilder
}

func Builder(st *stor.Stor) *builder {
	return &builder{stor: st, levels: []*level{{}}}
}

func (b *builder) Add(key string, off uint64) {
	if b.count > 0 {
		if key == b.prev {
			panic("btreeBuilder keys must not have duplicates")
		}
		if key < b.prev {
			panic("btreeBuilder keys must be inserted in order")
		}
	}
	b.add(0, key, off)
	b.prev = key
	b.count++
}

func (b *builder) add(li int, key string, off uint64) {
	if li >= len(b.levels) {
		b.levels = append(b.levels, &level{})
	}
	lev := b.levels[li]
	if len(lev.nb.node) > (MaxNodeSize * 3 / 4) {
		// split full node to stor
		offNode, splitKey := lev.nb.Split(b.stor)
		b.add(li+1, lev.splitKey, offNode) // RECURSE
		lev.splitKey = splitKey
	}
	embedLen := 1
	if li > 0 {
		embedLen = embedAll
	}
	lev.nb.Add(key, off, embedLen)
}

func (b *builder) Finish() *btree {
	var key string
	var off uint64
	for li := 0; li < len(b.levels); li++ {
		if li > 0 {
			// allow node to slightly exceed max size
			b.levels[li].nb.Add(key, off, embedAll)
		}
		key = b.levels[li].splitKey
		off = b.levels[li].nb.node.putNode(b.stor)
	}
	treeLevels := len(b.levels) - 1
	return OpenBtree(b.stor, off, treeLevels)
}

//-------------------------------------------------------------------

type nodeBuilder struct {
	node     node
	notFirst bool
	pos      int
	prev     string
	known    string
	offset   uint64
	pos2     int
	known2   string
	offset2  uint64
}

func (b *nodeBuilder) Add(key string, offset uint64, embedLen int) {
	if b.notFirst {
		if key <= b.prev {
			panic("fBuilder keys must be inserted in order, without duplicates")
		}
	} else {
		b.notFirst = true
	}
	if len(b.node) == 0 {
		b.node = b.node.append(offset, 0, "")
		b.known = ""
	} else {
		b.pos2 = b.pos
		b.pos = len(b.node)
		npre, diff, known := addone(key, b.prev, b.known, embedLen)
		b.node = b.node.append(offset, npre, diff)
		b.known2 = b.known
		b.known = known
		b.offset2 = b.offset
		b.offset = offset
	}
	b.prev = key
}

// Split saves all but the last two entries as the left node
// and sets the builder node to the last two entries
func (b *nodeBuilder) Split(st *stor.Stor) (leftOff uint64, splitKey string) {
	splitKey = b.known2 // known of second last entry
	left := b.node[:b.pos2]
	leftOff = left.putNode(st)
	// first entry becomes 0, ""
	right := b.node[:0].append(b.offset2, 0, "") // offset of second last entry
	// second entry becomes 0, known
	right = right.append(b.offset, 0, b.known) // offset,known of last entry
	b.node = right
	return
}

func (b *nodeBuilder) Entries() node {
	return b.node
}
