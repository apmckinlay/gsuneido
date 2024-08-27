// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"log"

	"github.com/apmckinlay/gsuneido/db19/stor"
)

// builder is used to bulk load an btree.
// Keys must be added in order.
// The btree is built bottom up with no splitting or inserting.
// All nodes will be "full" except for the right hand edge.
type builder struct {
	stor   *stor.Stor
	prev   string
	levels []*level // leaf is [0]
	count  int
}

type level struct {
	splitKey string
	nb       nodeBuilder
}

func Builder(st *stor.Stor) *builder {
	return &builder{stor: st, levels: []*level{{}}}
}

// Add returns false for duplicate keys and panics for out of order
func (b *builder) Add(key string, off uint64) bool {
	if b.count > 0 {
		if key == b.prev {
			return false
		}
		if key < b.prev {
			panic("btree builder keys must be inserted in order")
		}
	}
	b.add(0, key, off)
	b.prev = key
	b.count++
	return true
}

func (b *builder) add(li int, key string, off uint64) {
	if li >= len(b.levels) {
		b.levels = append(b.levels, &level{})
	}
	lev := b.levels[li]
	if shouldSplit(lev.nb.node, lev.nb.count) {
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

func shouldSplit(nd node, nodeCount int) bool {
	if len(nd) > MaxNodeSize && nodeCount >= MinSplitSize {
		return true
	}
	if len(nd) > 8*MaxNodeSize {
		log.Println("ERROR: btree node too large", len(nd), "count", nodeCount)
	}
	return false
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
	prev     string
	known    string
	known2   string
	node     node
	pos      int
	offset   uint64
	pos2     int
	offset2  uint64
	count    int
	notFirst bool
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
	b.count++
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
	b.count = 2
	return
}

func (b *nodeBuilder) Entries() node {
	return b.node
}
