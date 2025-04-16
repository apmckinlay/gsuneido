// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// package intern stores a set of strings.
// NOTE: strings are never removed or garbage collected (other than Clear)
package intern

import (
	"strings"
	"sync"
	"unsafe"

	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/shmap"
)

const chunkSize = 64 * 1024 // 64 kb to match int16 offset

var (
	htbl shmap.Map[entry, struct{}, helper]
	// chunks holds the actual strings
	chunks []*chunkType = []*chunkType{{}}
	// next offset in current chunk
	next int
	lock sync.Mutex
)

type chunkType [chunkSize]byte

type entry struct {
	size   uint8
	chunk  uint8
	offset uint16
}

func (e entry) str() string {
	return unsafe.String(&chunks[e.chunk][e.offset], e.size)
}

type helper struct{}

func (helper) Hash(e entry) uint64 {
	return hash.String(e.str())
}

func (helper) Equal(e1, e2 entry) bool {
	return e1.str() == e2.str()
}

// String returns a copy of the string in its private storage.
// If the string doesn't exist yet, it is added.
// Strings longer than 256 bytes are just cloned every time, not stored.
func String(s string) string {
	if len(s) >= 256 {
		return strings.Clone(s)
	}
	lock.Lock()
	defer lock.Unlock()
	e := add(s)
	k, exists := htbl.GetInit(e)
	if !exists {
		return e.str()
	}
	next -= int(e.size) // undo the add
	return k.str()
}

func add(s string) entry {
	if int(next)+len(s) > chunkSize {
		chunks = append(chunks, &chunkType{})
		next = 0
	}
	ci := len(chunks) - 1
	copy(chunks[ci][next:], s)
	offset := next
	next += len(s)
	return entry{size: uint8(len(s)), chunk: uint8(ci), offset: uint16(offset)}
}

// Count returns the number of unique strings stored
func Count() int {
	lock.Lock()
	defer lock.Unlock()
	return htbl.Size()
}

// Bytes returns the total allocated size of the strings
func Bytes() int {
	return (len(chunks)-1)*chunkSize + next
}

// Clear discards all the strings and starts fresh
func Clear() {
	lock.Lock()
	defer lock.Unlock()
	htbl.Clear()
	chunks = []*chunkType{{}}
	next = 0

}

// Recent returns the most recently added strings (for debugging)
func Recent() string {
	lock.Lock()
	defer lock.Unlock()
	return string(chunks[len(chunks)-1][max(0, next-256):next])
}
