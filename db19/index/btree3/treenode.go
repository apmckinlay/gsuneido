// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"log"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/cksum"
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

There are count fields, and count + 1 record offsets.
*/

type treeNode []byte

func (nd treeNode) nkeys() int {
	return int(nd[0])
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

// search returns the offset of the first entry that is >= key
func (nd treeNode) search(key string) uint64 {
	// binary search
	low, high := 0, nd.nkeys()-1
	for low <= high {
		mid := (low + high) / 2
		midKey := nd.key(mid)
		if string(midKey) <= key {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return nd.offset(low)
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
	sb.WriteString(" ")
	sb.WriteString(strconv.FormatUint(nd.offset(n), 10))
	sb.WriteString("}")
	return sb.String()
}

// ------------------------------------------------------------------

type treeIter struct {
	nd treeNode
	i  int
}

// iter returns an iterator over the offsets (not the separators)
func (nd treeNode) iter() *treeIter {
	return &treeIter{nd: nd, i: -1}
}

func (it *treeIter) next() bool {
	if it.i+1 > it.nd.nkeys() {
		return false
	}
	it.i++
	return true
}

func (it *treeIter) prev() bool {
	if it.i <= 0 {
		return false
	}
	it.i--
	return true
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

func (b *treeBuilder) finish(offset uint64) treeNode {
	b.offsets = append(b.offsets, offset)
	b.keys = append(b.keys, "")
	n := len(b.keys)
	if n > 256 {
		panic("too many keys for tree node")
	}
	totalKeyLen := 0
	for _, key := range b.keys {
		totalKeyLen += len(key)
	}
	nodeSize := 1 + n*7 + totalKeyLen
	result := make([]byte, 0, nodeSize)
	result = append(result, byte(n-1))
	fieldPos := 1 + n*7 // position where field data starts
	for i, off := range b.offsets {
		result = append(result,
			byte(fieldPos>>8),
			byte(fieldPos),
			byte(off>>32),
			byte(off>>24),
			byte(off>>16),
			byte(off>>8),
			byte(off))
		fieldPos += len(b.keys[i])
	}
	for _, key := range b.keys {
		result = append(result, key...)
	}
	return treeNode(result)
}

func (b *treeBuilder) reset() {
	b.keys = b.keys[:0]
	b.offsets = b.offsets[:0]
	b.entrySize = 0
}

// ------------------------------------------------------------------

// insert inserts an entry, maintaining order
// func (nd treeNode) insert(key string, offset uint64) treeNode {
// 	count := nd.nkeys()

// 	// Find insertion position using binary search
// 	pos := nd.search(key)

// 	// Check if key already exists (should not happen for tree nodes)
// 	if pos < count && nd.getKey(pos) == key {
// 		panic("duplicate key")
// 	}

// 	// Calculate new count
// 	newCount := count + 1

// 	// Grow the node to accommodate new entry
// 	// We need space for: new count byte + (newCount * 7 bytes for entries) + all keys
// 	totalKeyLen := 0
// 	for i := 0; i < count; i++ {
// 		totalKeyLen += len(nd.getKey(i))
// 	}
// 	totalKeyLen += len(key)

// 	newSize := 1 + newCount*7 + totalKeyLen
// 	nd = slc.Grow(nd, newSize-len(nd))

// 	// Shift data after insertion point to make room for new entry
// 	// The data to shift includes: entries after pos, and all keys
// 	entryStart := 1 + pos*7
// 	keyDataStart := 1 + count*7 + 2 // after all entries and size field

// 	// Shift the entries after insertion point
// 	copy(nd[entryStart+7:], nd[entryStart:1+count*7])

// 	// Shift all keys to make room for the new key
// 	copy(nd[keyDataStart+len(key):], nd[keyDataStart:])

// 	// Insert new entry data
// 	// Calculate field position for the new key
// 	fieldPos := 1 + newCount*7
// 	for i := 0; i < pos; i++ {
// 		fieldPos += len(nd.getKey(i))
// 	}
// 	fieldPos += len(key)
// 	for i := pos; i < count; i++ {
// 		fieldPos += len(nd.getKey(i))
// 	}

// 	// Set the new entry
// 	nd[entryStart] = byte(fieldPos >> 8)
// 	nd[entryStart+1] = byte(fieldPos)
// 	nd[entryStart+2] = byte(offset >> 32)
// 	nd[entryStart+3] = byte(offset >> 24)
// 	nd[entryStart+4] = byte(offset >> 16)
// 	nd[entryStart+5] = byte(offset >> 8)
// 	nd[entryStart+6] = byte(offset)

// 	// Insert the key data
// 	keyPos := 1 + newCount*7
// 	for i := 0; i < pos; i++ {
// 		keyPos += len(nd.getKey(i))
// 	}
// 	copy(nd[keyPos:], key)

// 	// Update count
// 	nd[0] = byte(newCount)

// 	return nd
// }
