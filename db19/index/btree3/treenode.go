// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"log"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
)

/*
treeNode is:
- 1 byte array count
- array of 2 byte field offset + 5 byte db offset
- 2 byte end offset (node size)
- 5 byte final db offset
- field data, contiguous, sorted

The node size and final offset can be treated as an additional array element.

offsets are stored most significant byte first.

There are nkeys fields, and noffs = nkeys + 1 record offsets.
*/

type treeNode []byte

const emptyTree = "\x00\x00\x03"

func (nd treeNode) nkeys() int {
	return int(nd[0])
}

func (nd treeNode) noffs() int {
	n := nd.nkeys()
	if n == 0 && nd.size() == 3 {
		return 0
	}
	return n + 1
}

// size returns the length of a treeNode
func (nd treeNode) size() int {
	count := nd.nkeys()
	pos := 1 + count*7
	return int(uint16(nd[pos])<<8 | uint16(nd[pos+1]))
}

// key returns the i'th field
func (nd treeNode) key(i int) []byte {
	base := 1 + i*7
	fieldPos := uint16(nd[base])<<8 | uint16(nd[base+1])
	endPos := uint16(nd[base+7])<<8 | uint16(nd[base+8])
	return nd[fieldPos:endPos]
}

// offset returns the i'th offset.
// There are nkeys + 1 offsets.
func (nd treeNode) offset(i int) uint64 {
	pos := 1 + i*7 + 2
	return uint64(nd[pos])<<32 |
		uint64(nd[pos+1])<<24 |
		uint64(nd[pos+2])<<16 |
		uint64(nd[pos+3])<<8 |
		uint64(nd[pos+4])
}

// search returns the position and offset of the first entry that is >= key
func (nd treeNode) search(key string) (int, uint64) {
	// binary search
	lo, hi := 0, nd.nkeys()-1
	for lo <= hi {
		mid := (lo + hi) / 2
		if string(nd.key(mid)) <= key {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return lo, nd.offset(lo)
}

func (nd treeNode) seek(key string) treeIter {
	lo, hi := 0, nd.nkeys()-1
	for lo <= hi {
		mid := (lo + hi) / 2
		if string(nd.key(mid)) <= key {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return treeIter{nd: nd, i: lo}
}

// write writes a tree node to storage
func (nd treeNode) write(st *stor.Stor) uint64 {
	n := len(nd)
	if n > 8192 {
		log.Println("ERROR: btree node too large")
	}
	off, buf := st.Alloc(n + cksum.Len)
	copy(buf, nd)
	cksum.Update(buf)
	return off
}

// splitTo splits a tree node into two nodes, writing them directly to storage.
// Returns the storage offsets of the left and right nodes, and the split key.
// Unlike leafNode.splitTo, this does NOT shorten the split key.
func (nd treeNode) splitTo(st *stor.Stor) (leftOff, rightOff uint64, splitKey string) {
	nkeys := nd.nkeys()
	splitPos := nkeys / 2

	// The split key is the key at the split position (NOT shortened)
	splitKey = string(nd.key(splitPos))

	// Calculate sizes for both nodes
	leftSize := 1 + splitPos*7 + 2 + 5 // count + offsets + size + final offset
	for i := 0; i < splitPos; i++ {
		leftSize += len(nd.key(i))
	}

	rightCount := nkeys - splitPos - 1
	rightSize := 1 + rightCount*7 + 2 + 5 // count + offsets + size + final offset
	for i := splitPos + 1; i < nkeys; i++ {
		rightSize += len(nd.key(i))
	}

	// Allocate storage for both nodes
	leftOff, leftBuf := st.Alloc(leftSize + cksum.Len)
	rightOff, rightBuf := st.Alloc(rightSize + cksum.Len)

	// Build left node
	leftBuf[0] = byte(splitPos)

	fieldPos := 1 + splitPos*7 + 2 + 5
	pos := 1
	for i := 0; i < splitPos; i++ {
		off := nd.offset(i)
		leftBuf[pos] = byte(fieldPos >> 8)
		leftBuf[pos+1] = byte(fieldPos)
		leftBuf[pos+2] = byte(off >> 32)
		leftBuf[pos+3] = byte(off >> 24)
		leftBuf[pos+4] = byte(off >> 16)
		leftBuf[pos+5] = byte(off >> 8)
		leftBuf[pos+6] = byte(off)
		pos += 7
		fieldPos += len(nd.key(i))
	}
	// Add final offset for left node (offset at splitPos)
	leftBuf[pos] = byte(fieldPos >> 8)
	leftBuf[pos+1] = byte(fieldPos)
	finalOff := nd.offset(splitPos)
	leftBuf[pos+2] = byte(finalOff >> 32)
	leftBuf[pos+3] = byte(finalOff >> 24)
	leftBuf[pos+4] = byte(finalOff >> 16)
	leftBuf[pos+5] = byte(finalOff >> 8)
	leftBuf[pos+6] = byte(finalOff)
	pos += 7

	for i := 0; i < splitPos; i++ {
		key := nd.key(i)
		copy(leftBuf[pos:], key)
		pos += len(key)
	}

	cksum.Update(leftBuf)

	// Build right node
	rightBuf[0] = byte(rightCount)

	fieldPos = 1 + rightCount*7 + 2 + 5
	pos = 1
	for i := splitPos + 1; i < nkeys; i++ {
		off := nd.offset(i)
		rightBuf[pos] = byte(fieldPos >> 8)
		rightBuf[pos+1] = byte(fieldPos)
		rightBuf[pos+2] = byte(off >> 32)
		rightBuf[pos+3] = byte(off >> 24)
		rightBuf[pos+4] = byte(off >> 16)
		rightBuf[pos+5] = byte(off >> 8)
		rightBuf[pos+6] = byte(off)
		pos += 7
		fieldPos += len(nd.key(i))
	}
	// Add final offset for right node (offset at nkeys)
	rightBuf[pos] = byte(fieldPos >> 8)
	rightBuf[pos+1] = byte(fieldPos)
	finalOff = nd.offset(nkeys)
	rightBuf[pos+2] = byte(finalOff >> 32)
	rightBuf[pos+3] = byte(finalOff >> 24)
	rightBuf[pos+4] = byte(finalOff >> 16)
	rightBuf[pos+5] = byte(finalOff >> 8)
	rightBuf[pos+6] = byte(finalOff)
	pos += 7

	for i := splitPos + 1; i < nkeys; i++ {
		key := nd.key(i)
		copy(rightBuf[pos:], key)
		pos += len(key)
	}

	cksum.Update(rightBuf)

	return leftOff, rightOff, splitKey
}

// readTree reads a leaf node from storage
func readTree(st *stor.Stor, off uint64) treeNode {
	node := treeNode(st.Data(off))
	return node[:node.size()]
}

func (nd treeNode) String() string {
	if len(nd) == 0 {
		return "tree{}"
	}
	var sb strings.Builder
	sb.WriteString("tree{")
	sep := ""
	n := nd.nkeys()
	for i := range n {
		sb.WriteString(sep)
		sep = " "
		sb.WriteString(strconv.FormatUint(nd.offset(i), 10))
		sb.WriteString(" ")
		sb.WriteString("<" + string(nd.key(i)) + ">")
	}
	if nd.noffs() > 0 {
		sb.WriteString(" ")
		sb.WriteString(strconv.FormatUint(nd.offset(n), 10))
	}
	sb.WriteString("}")
	return sb.String()
}

// ------------------------------------------------------------------

type treeIter struct {
	nd treeNode
	i  int
}

// iter returns an iterator over the offsets (not the separators)
func (nd treeNode) iter() treeIter {
	return treeIter{nd: nd, i: -1}
}

func (it *treeIter) next() bool {
	if it.i >= it.nd.noffs() {
		return false
	}
	it.i++
	return it.i < it.nd.noffs()
}

func (it *treeIter) prev() bool {
	if it.i < 0 {
		return false
	}
	it.i--
	return it.i >= 0
}

func (it *treeIter) off() uint64 {
	return it.nd.offset(it.i)
}

// ------------------------------------------------------------------

type treeBuilder struct {
	keys      []string
	offsets   []uint64
	entrySize int
}

func (b *treeBuilder) add(offset uint64, key string) {
	b.keys = append(b.keys, key)
	b.offsets = append(b.offsets, offset)
	b.entrySize += len(key) + 7
}

func (b *treeBuilder) size() int {
	return b.entrySize + 8
}

func (b *treeBuilder) nkeys() int {
	return len(b.keys)
}

// finishInto builds the tree node directly into the provided buffer with a final offset
// buf must be large enough (at least b.size() bytes)
func (b *treeBuilder) finishInto(buf []byte, offset uint64) treeNode {
	n := len(b.keys)
	if n > 255 {
		panic("too many keys for tree node")
	}

	buf[0] = byte(n)

	fieldPos := 1 + (n+1)*7 // position where field data starts
	pos := 1

	for i, off := range b.offsets {
		buf[pos] = byte(fieldPos >> 8)
		buf[pos+1] = byte(fieldPos)
		buf[pos+2] = byte(off >> 32)
		buf[pos+3] = byte(off >> 24)
		buf[pos+4] = byte(off >> 16)
		buf[pos+5] = byte(off >> 8)
		buf[pos+6] = byte(off)
		pos += 7
		fieldPos += len(b.keys[i])
	}

	// Add final offset
	buf[pos] = byte(fieldPos >> 8)
	buf[pos+1] = byte(fieldPos)
	buf[pos+2] = byte(offset >> 32)
	buf[pos+3] = byte(offset >> 24)
	buf[pos+4] = byte(offset >> 16)
	buf[pos+5] = byte(offset >> 8)
	buf[pos+6] = byte(offset)
	pos += 7

	for _, key := range b.keys {
		copy(buf[pos:], key)
		pos += len(key)
	}

	return treeNode(buf[:pos])
}

// finish adds a final offset and then builds the tree node
func (b *treeBuilder) finish(offset uint64) treeNode {
	size := b.size()
	buf := make([]byte, size)
	return b.finishInto(buf, offset)
}

// finishTo builds the tree node with a final offset and writes it to storage
func (b *treeBuilder) finishTo(st *stor.Stor, offset uint64) uint64 {
	size := b.size()
	off, buf := st.Alloc(size + cksum.Len)
	b.finishInto(buf[:size], offset)
	cksum.Update(buf)
	return off
}

// reset clears the builder, keeping its capacity to save allocation
func (b *treeBuilder) reset() {
	b.keys = b.keys[:0]
	b.offsets = b.offsets[:0]
	b.entrySize = 0
}

// ------------------------------------------------------------------

// insert inserts an entry, maintaining order
func (nd treeNode) insert(i int, newoff uint64, key string) treeNode {
	// Handle empty node case
	if len(nd) == 0 {
		var b treeBuilder
		b.add(newoff, key)
		return b.finish(0) // Add final offset as 0 for now
	}

	n := nd.nkeys()
	if n == 255 {
		panic("too many keys for tree node")
	}

	i = min(i, n)

	fieldLen := len(key)
	oldSize := nd.size()

	// Calculate where field data will be inserted (before any shifts)
	oldOffsetArrayEnd := 1 + n*7 + 2 + 5 // count + offsets + size + final offset
	fieldInsertPos := oldOffsetArrayEnd
	for j := range i {
		fieldInsertPos += len(nd.key(j))
	}

	// Calculate total size increase: 7 bytes for offset entry + field data length
	sizeIncrease := 7 + fieldLen
	newSize := oldSize + sizeIncrease

	// Grow the slice to accommodate the new data
	nd = slc.Grow(nd, sizeIncrease)

	// Update count
	nd[0] = byte(n + 1)

	// REVERSE ORDER OPERATIONS (work backwards to avoid overwriting data):

	// 1. First move field data that comes AFTER the insertion point
	if i < n {
		// Source: field data after insertion point in old layout
		srcStart := fieldInsertPos
		srcEnd := oldSize
		// Destination: same data after making room for both offset entry and new field data
		dstStart := fieldInsertPos + 7 + fieldLen
		copy(nd[dstStart:], nd[srcStart:srcEnd])
	}

	// 2. Move offset array entries at position i and later PLUS size field + final offset + field data before insertion point
	// These are physically contiguous in memory
	entryInsertPos := 1 + i*7
	copy(nd[entryInsertPos+7:], nd[entryInsertPos:fieldInsertPos])

	// 3. Insert the new field data in its correct position (adjusted for offset array growth)
	newFieldPos := fieldInsertPos + 7
	copy(nd[newFieldPos:], key)

	// 4. Insert the new offset entry
	nd[entryInsertPos] = byte(newFieldPos >> 8)
	nd[entryInsertPos+1] = byte(newFieldPos)
	nd[entryInsertPos+2] = byte(newoff >> 32)
	nd[entryInsertPos+3] = byte(newoff >> 24)
	nd[entryInsertPos+4] = byte(newoff >> 16)
	nd[entryInsertPos+5] = byte(newoff >> 8)
	nd[entryInsertPos+6] = byte(newoff)

	// 5. Update field positions in all offset entries
	for j := 0; j < n+1; j++ {
		if j == i {
			continue // already set above
		}
		pos := 1 + j*7
		oldFieldPos := int(uint16(nd[pos])<<8 | uint16(nd[pos+1]))

		// All field positions need to account for offset array growth (+7)
		newFieldPos := oldFieldPos + 7

		// Field positions after insertion point also need to account for new field data
		if j > i {
			newFieldPos += fieldLen
		}

		nd[pos] = byte(newFieldPos >> 8)
		nd[pos+1] = byte(newFieldPos)
	}

	// 6. Update the end offset
	endPos := 1 + (n+1)*7
	nd[endPos] = byte(newSize >> 8)
	nd[endPos+1] = byte(newSize)

	return nd[:newSize]
}

// update modifies the offset for a key, in place
func (nd treeNode) update(i int, off uint64) {
	assert.That(i < nd.noffs())
	pos := 1 + i*7 + 2
	nd[pos] = byte(off >> 32)
	nd[pos+1] = byte(off >> 24)
	nd[pos+2] = byte(off >> 16)
	nd[pos+3] = byte(off >> 8)
	nd[pos+4] = byte(off)
}

// delete needs to return the updated treeNode because the size will change
// even though it alters in place and does not allocate
// and always returns the same memory it was given
func (nd treeNode) delete(i int) treeNode {
	newNkeys := nd.nkeys() - 1
	if newNkeys == -1 {
		return append(nd[:0], emptyTree...)
	}
	oldSize := nd.size()

	// handle deleting the final offset
	// which deletes the last key and moves its offset to the final position
	if i == nd.nkeys() {
		nd.update(i, nd.offset(i-1))
		i--
	}

	// Get the field positions for the entry to delete
	base := 1 + i*7
	fieldStart := int(uint16(nd[base])<<8 | uint16(nd[base+1]))
	fieldEnd := int(uint16(nd[base+7])<<8 | uint16(nd[base+8]))
	fieldLen := fieldEnd - fieldStart

	nd[0] = byte(newNkeys) // Update count

	// shift the middle section
	entryStart := base
	entryEnd := base + 7
	copy(nd[entryStart:], nd[entryEnd:fieldStart])

	// Update field positions in all remaining entries
	// All field positions need to be reduced by 7 (removed offset entry)
	// Plus, entries after the deleted field need reduction by fieldLen
	for j := 0; j < newNkeys; j++ {
		pos := 1 + j*7
		oldFieldPos := int(uint16(nd[pos])<<8 | uint16(nd[pos+1]))
		newFieldPos := oldFieldPos - 7 // account for removed offset entry
		if j >= i {
			// This entry was after the deleted field
			newFieldPos -= fieldLen
		}
		nd[pos] = byte(newFieldPos >> 8)
		nd[pos+1] = byte(newFieldPos)
	}

	// Update the end offset (now at position 1 + (n-1)*7)
	endPos := 1 + (newNkeys)*7
	newSize := oldSize - 7 - fieldLen
	nd[endPos] = byte(newSize >> 8)
	nd[endPos+1] = byte(newSize)

	// shift the field data after the deleted field
	copy(nd[fieldStart-7:], nd[fieldEnd:])

	// Return the truncated slice
	return nd[:newSize]
}
