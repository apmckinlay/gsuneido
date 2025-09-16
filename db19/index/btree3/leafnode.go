// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"log"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/cksum"
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

// search returns the offset of key, or 0 if not found
func (nd leafNode) search(key string) uint64 {
	// NOTE: search should not do any allocation
	prefix := nd.prefix()
	prefixLen := len(prefix)
	if prefixLen > 0 {
		prefixStr := string(prefix)
		if key < prefixStr || !strings.HasPrefix(key, prefixStr) {
			return 0
		}
	}
	// key starts with prefix, binary search for suffix
	keySuffix := key[prefixLen:]
	low, high := 0, nd.nkeys()-1
	for low <= high {
		mid := (low + high) / 2
		midSuffix := string(nd.suffix(mid))
		if midSuffix < keySuffix {
			low = mid + 1
		} else if midSuffix > keySuffix {
			high = mid - 1
		} else {
			return nd.offset(mid) // found
		}
	}
	return 0 // not found
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

type leafIter struct {
	nd leafNode
	i  int
}

// iter returns an iterator over the offsets (not the separators)
func (nd leafNode) iter() *leafIter {
	return &leafIter{nd: nd, i: -1}
}

func (it *leafIter) next() bool {
	if it.i+1 >= it.nd.nkeys() {
		return false
	}
	it.i++
	return true
}

func (it *leafIter) prev() bool {
	if it.i <= 0 {
		return false
	}
	it.i--
	return true
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

func (b *leafBuilder) empty() bool {
	return len(b.keys) == 0
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

func (b *leafBuilder) reset() {
	b.keys = b.keys[:0]
	b.offsets = b.offsets[:0]
	b.fieldsLen = 0
	b.prefix = ""
}

// ------------------------------------------------------------------

// insert inserts an entry, maintaining order
// func (nd leafNode) insert(key string, off uint64) leafNode {
// 	pos, prev, found := nd.search(key)
// 	if string(found) == key {
// 		panic("duplicate key")
// 	}
// 	n := leafEntryLen(prev, key)
// 	nd = slc.Grow(nd, n)

// 	copy(nd[pos+n:], nd[pos:])
// 	gap := nd[pos : pos : pos+n]
// 	assert.That(unsafe.SliceData(nd[pos:]) == unsafe.SliceData(gap))
// 	result := gap.add(prev, key, off)
// 	assert.That(unsafe.SliceData(gap) == unsafe.SliceData(result))

// 	return nd
// }
