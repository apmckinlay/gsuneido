// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"log"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/cksum"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/str"
)

/*
leafNode is:
- 1 byte array count
- 1 byte prefix length
- array of 2 byte field offset + 5 byte db offset
- 2 byte end offset (node size)
- prefix
- field data, contiguous, sorted

The node size can be treated as an additional array element.

offsets are stored most significant byte first.
*/

type leafNode []byte

func (nd leafNode) nkeys() int {
	return int(nd[0])
}

// size returns the length of a treeNode
func (nd leafNode) size() int {
	count := nd.nkeys()
	pos := 2 + count*7
	return int(uint16(nd[pos])<<8 | uint16(nd[pos+1]))
}

// offset returns the i'th offset
func (nd leafNode) offset(i int) uint64 {
	pos := 2 + i*7 + 2
	return uint64(nd[pos])<<32 |
		uint64(nd[pos+1])<<24 |
		uint64(nd[pos+2])<<16 |
		uint64(nd[pos+3])<<8 |
		uint64(nd[pos+4])
}

// key returns the i'th field.
// It allocates.
func (nd leafNode) key(i int) string {
	base := 2 + i*7
	fieldPos := uint16(nd[base])<<8 | uint16(nd[base+1])
	endPos := uint16(nd[base+7])<<8 | uint16(nd[base+8])
	return cat(nd.prefix(), nd[fieldPos:endPos])
}

func cat(x, y []byte) string {
	buf := make([]byte, len(x)+len(y))
	copy(buf[copy(buf, x):], y)
	return hacks.BStoS(buf)
}

// suffix returns the i'th field without the shared prefix.
// It does not allocate.
func (nd leafNode) suffix(i int) []byte {
	base := 2 + i*7
	fieldPos := uint16(nd[base])<<8 | uint16(nd[base+1])
	endPos := uint16(nd[base+7])<<8 | uint16(nd[base+8])
	return nd[fieldPos:endPos]
}

// prefix returns the shared prefix for the node.
// It does not allocate.
func (nd leafNode) prefix() []byte {
	prelen := int(nd[1])
	if prelen == 0 {
		return nil
	}
	pos := 4 + 7*nd.nkeys()
	end := pos + prelen
	return nd[pos:end]
}

// search returns the position and whether the key was found
// if the key isn't found, the returned position is where it would be
func (nd leafNode) search(key string) (int, bool) {
	// NOTE: search should not do any allocation
	prefix := nd.prefix()
	prefixLen := len(prefix)
	if prefixLen > 0 {
		if key < string(prefix) {
			return 0, false // key would be before all entries
		}
		if !str.HasPrefix(key, prefix) {
			return nd.nkeys(), false // key would be after all entries
		}
	}
	// key starts with prefix, binary search for suffix
	keySuffix := key[prefixLen:]
	lo, hi := 0, nd.nkeys()-1
	for lo <= hi {
		mid := (lo + hi) / 2
		midSuffix := string(nd.suffix(mid))
		if midSuffix < keySuffix {
			lo = mid + 1
		} else if midSuffix > keySuffix {
			hi = mid - 1
		} else {
			return mid, true // found
		}
	}
	return lo, false // not found, lo is where it would be inserted
}

// seek returns a leafIter positioned at the first key >= key
func (nd leafNode) seek(key string) *leafIter {
	prefix := nd.prefix()
	prefixLen := len(prefix)
	if prefixLen > 0 {
		if key < string(prefix) {
			// key is less than prefix, no entries <= key
			return &leafIter{nd: nd, i: 0}
		} else if !str.HasPrefix(key, prefix) {
			// key doesn't start with prefix but key >= prefix,
			// so all entries are < key
			return &leafIter{nd: nd, i: nd.nkeys()}
		}
	}
	// key starts with prefix (or no prefix), binary search for suffix
	keySuffix := key[prefixLen:]
	lo, hi := 0, nd.nkeys()
	for lo < hi {
		mid := (lo + hi) / 2
		if string(nd.suffix(mid)) < keySuffix {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return &leafIter{nd: nd, i: lo}
}

// write writes a leaf node to storage
func (nd leafNode) write(st *stor.Stor) uint64 {
	n := len(nd)
	if n > 8192 {
		log.Println("ERROR: btree node too large")
	}
	off, buf := st.Alloc(n + cksum.Len)
	copy(buf, nd)
	cksum.Update(buf)
	return off
}

// readLeaf reads a leaf node from storage
func readLeaf(st *stor.Stor, off uint64) leafNode {
	node := leafNode(st.Data(off))
	return node[:node.size()]
}

func (nd leafNode) String() string {
	if len(nd) == 0 {
		return "leaf{}"
	}
	var sb strings.Builder
	sb.WriteString("leaf{")
	if pre := nd.prefix(); len(pre) > 0 {
		sb.WriteByte('|')
		sb.Write(pre)
		sb.WriteString("| ")
	}
	sep := ""
	n := nd.nkeys()
	for i := range n {
		sb.WriteString(sep)
		sep = " "
		sb.Write(nd.suffix(i))
		sb.WriteString(" ")
		sb.WriteString(strconv.FormatUint(nd.offset(i), 10))
	}
	sb.WriteString("}")
	return sb.String()
}

// ------------------------------------------------------------------

// leafIter iterates over a leaf node.
type leafIter struct {
	nd leafNode
	i  int
}

// iter starts at -1 so next goes to the first entry
func (nd leafNode) iter() *leafIter {
	return &leafIter{nd: nd, i: -1}
}

// next moves to the next entry.
// It returns false if there are no more entries.
func (it *leafIter) next() bool {
	if it.i >= it.nd.nkeys() {
		return false
	}
	it.i++
	return it.i < it.nd.nkeys()
}

// prev moves to the previous entry.
// It returns false if there are no more entries.
func (it *leafIter) prev() bool {
	if it.i < 0 {
		return false
	}
	it.i--
	return it.i >= 0
}

// eof returns true if past the beginning or end
func (it *leafIter) eof() bool {
	return it.i < 0 || it.i >= it.nd.nkeys()
}

func (it *leafIter) key(buf []byte) []byte {
	return append(append(buf[:0], it.nd.prefix()...), it.nd.suffix(it.i)...)
}

func (it *leafIter) suffix() []byte {
	return it.nd.suffix(it.i)
}

func (it *leafIter) offset() uint64 {
	return it.nd.offset(it.i)
}

// ------------------------------------------------------------------

type leafBuilder struct {
	keys      []string
	offsets   []uint64
	fieldsLen int
	prefix    string
}

func (b *leafBuilder) add(key string, offset uint64) {
	if len(b.keys) == 0 {
		b.prefix = key
	} else {
		b.prefix = str.CommonPrefix(b.prefix, key)
	}
	b.keys = append(b.keys, key)
	b.offsets = append(b.offsets, offset)
	b.fieldsLen += len(key)
}

// size returns the size the leaf node would be if finished now
func (b *leafBuilder) size() int {
	n := len(b.keys)
	prelen := len(b.prefix)
	fieldsLen := b.fieldsLen - n*prelen
	return 4 + 7*n + prelen + fieldsLen
}

func (b *leafBuilder) nkeys() int {
	return len(b.keys)
}

func (b *leafBuilder) finish() leafNode {
	n := len(b.keys)
	if n > 255 {
		panic("too many keys for leaf node")
	}
	if n == 1 {
		b.prefix = ""
	}
	nodeSize := b.size()
	prelen := len(b.prefix)
	result := make([]byte, 0, nodeSize)
	result = append(result, byte(n), byte(prelen))
	fieldPos := 4 + n*7 + prelen // position where field data starts
	for i, off := range b.offsets {
		result = append(result,
			byte(fieldPos>>8),
			byte(fieldPos),
			byte(off>>32),
			byte(off>>24),
			byte(off>>16),
			byte(off>>8),
			byte(off))
		fieldPos += len(b.keys[i]) - prelen
	}
	result = append(result, byte(fieldPos>>8), byte(fieldPos))
	result = append(result, b.prefix...)
	for _, key := range b.keys {
		result = append(result, key[prelen:]...)
	}
	return leafNode(result)
}

// reset clears the builder, keeping its capacity to save allocation
func (b *leafBuilder) reset() {
	b.keys = b.keys[:0]
	b.offsets = b.offsets[:0]
	b.fieldsLen = 0
	b.prefix = ""
}

// ------------------------------------------------------------------

// modify modifies a node, inserting, updating, or deleting an entry
// based on the tag on the offset
// func (nd leafNode) modify(key string, off uint64) leafNode {
// 	i, found := nd.search(key)
// 	if off&ixbuf.Update != 0 {
// 		assert.That(found)
// 		return nd.update(i, off&ixbuf.Mask)
// 	}
// 	if off&ixbuf.Delete != 0 {
// 		assert.That(found)
// 		return nd.delete(i)
// 	}
// 	return nd.insert(i, key, off&ixbuf.Mask)
// }

// insert inserts an entry, maintaining order
func (nd leafNode) insert(i int, key string, newoff uint64) leafNode {
	// Handle empty node case
	if len(nd) == 0 {
		var b leafBuilder
		b.add(key, newoff)
		return b.finish()
	}

	n := nd.nkeys()
	if n == 255 {
		panic("too many keys for leaf node")
	}

	// If the new key doesn't start with the current prefix, use leafBuilder
	prefix := nd.prefix()
	if n > 0 && !str.HasPrefix(key, prefix) {
		var b leafBuilder
		for j := range i {
			b.add(nd.key(j), nd.offset(j))
		}
		b.add(key, newoff)
		for j := i; j < n; j++ {
			b.add(nd.key(j), nd.offset(j))
		}
		return b.finish()
	}

	// Prefix won't change, proceed with in-place insertion
	keySuffix := key[len(prefix):]
	fieldLen := len(keySuffix)
	oldSize := nd.size()

	// Calculate where field data will be inserted (before any shifts)
	oldOffsetArrayEnd := 2 + n*7 + 2
	fieldInsertPos := oldOffsetArrayEnd + len(prefix)
	for j := 0; j < i; j++ {
		fieldInsertPos += len(nd.suffix(j))
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

	// 2. Move offset array entries at position i and later PLUS prefix and field data before insertion point
	// These are physically contiguous in memory
	entryInsertPos := 2 + i*7
	copy(nd[entryInsertPos+7:], nd[entryInsertPos:fieldInsertPos])

	// 3. Insert the new field data in its correct position (adjusted for offset array growth)
	newFieldPos := fieldInsertPos + 7
	copy(nd[newFieldPos:], keySuffix)

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
		pos := 2 + j*7
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
	endPos := 2 + (n+1)*7
	nd[endPos] = byte(newSize >> 8)
	nd[endPos+1] = byte(newSize)

	return nd[:newSize]
}

// update modifies the offset for a key, in place
func (nd leafNode) update(i int, off uint64) leafNode {
	pos := 2 + i*7 + 2 // skip to the 5-byte offset field
	nd[pos] = byte(off >> 32)
	nd[pos+1] = byte(off >> 24)
	nd[pos+2] = byte(off >> 16)
	nd[pos+3] = byte(off >> 8)
	nd[pos+4] = byte(off)
	return nd
}

func (nd leafNode) delete(i int) leafNode {
	n := nd.nkeys()
	if n == 1 {
		return nd[:0] // Deleting the only key results in an empty node
	}
	oldSize := nd.size()

	// Get the field positions for the entry to delete
	base := 2 + i*7
	fieldStart := int(uint16(nd[base])<<8 | uint16(nd[base+1]))
	fieldEnd := int(uint16(nd[base+7])<<8 | uint16(nd[base+8]))
	fieldLen := fieldEnd - fieldStart

	nd[0] = byte(n - 1) // Update count

	// shift the middle section
	entryStart := base
	entryEnd := base + 7
	copy(nd[entryStart:], nd[entryEnd:fieldStart])

	// Update field positions in all remaining entries
	// All field positions need to be reduced by 7 (removed offset entry)
	// Plus, entries after the deleted field need reduction by fieldLen
	for j := 0; j < n-1; j++ {
		pos := 2 + j*7
		oldFieldPos := int(uint16(nd[pos])<<8 | uint16(nd[pos+1]))
		newFieldPos := oldFieldPos - 7 // account for removed offset entry
		if j >= i {
			// This entry was after the deleted field
			newFieldPos -= fieldLen
		}
		nd[pos] = byte(newFieldPos >> 8)
		nd[pos+1] = byte(newFieldPos)
	}

	// Update the end offset (now at position 2 + (n-1)*7)
	endPos := 2 + (n-1)*7
	newSize := oldSize - 7 - fieldLen
	nd[endPos] = byte(newSize >> 8)
	nd[endPos+1] = byte(newSize)

	// shift the field data after the deleted field
	copy(nd[fieldStart-7:], nd[fieldEnd:])

	// Handle special case: going from 2 keys to 1 key
	// Need to clear prefix and adjust field position
	if n-1 == 1 {
		nd[1] = 0  // clear prefix length
		nd[2] = 0  // field position high byte
		nd[3] = 11 // field position low byte (4 + 1*7)
	}

	// Return the truncated slice
	return nd[:newSize]
}
