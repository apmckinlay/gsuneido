// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows

// Package heap provides a heap for win32 dll arguments
// that are allocated and freed stack-wise.
// It must not move and is therefore fixed size and statically declared.
// It is NOT thread safe. (Since win32 should be single threaded.)
// Normal usage is:
// 		defer FreeTo(CurSize())
//		...Alloc...
//
// Originally, data was just declared locally in the functions.
// But it crashed a lot (although intermittently and randomly).
// I think the crashes were caused by Go moving (resizing) stacks.
// Although I thought escape analysis would be putting the data on the heap.
// The crashes were mostly (always?) with nested callbacks.
// see: https://github.com/lxn/walk/pull/493
//
// Allocated memory is zero filled by FreeTo
// since that is less overhead than alloc zeroing.
package heap

import (
	"bytes"
	"unsafe"

	"github.com/apmckinlay/gsuneido/util/assert"
)

const align = 8 // Alloc assumes power of two
const heapsize = 64 * 1024

var heap [heapsize]byte
var heapnext = 0

// Alloc returns an unsafe.Pointer to n bytes of zero heap space.
func Alloc(n uintptr) unsafe.Pointer {
	i := alloc(int(n))
	return unsafe.Pointer(&heap[i])
}

func alloc(n int) int {
	n = ((n - 1) | (align - 1)) + 1 // requires align is power of 2
	if heapnext + n > heapsize {
		panic("Windows dll interface argument space limit exceeded")
	}
	i := heapnext
	heapnext += n
	return i
}

// CopyStr copies the string into the heap with a nul terminator (the +1).
// Due to alignment, "" is guaranteed to have two nuls as required by UTF16.
func CopyStr(s string) unsafe.Pointer {
	return Copy(s, len(s)+1)
}

// Copy allocates n bytes on the heap and copies the string into it.
// If n > len(s) the excess will be zero.
// If len(s) > n only n bytes will be copied.
func Copy(s string, n int) unsafe.Pointer {
	i := alloc(n)
	copy(heap[i:i+n], s)
	return unsafe.Pointer(&heap[i])
}

func CurSize() int {
	return heapnext
}

// GetStrN return a string containing a copy of a slice of the heap
// starting at p and n bytes long
func GetStrN(p unsafe.Pointer, n int) string {
	return string(get(p)[:n])
}

// GetStrZ return a string containing a copy of a slice of the heap
// starting at p, up to the first nul or n bytes, whichever comes first.
func GetStrZ(p unsafe.Pointer, n int) string {
	buf := get(p)[:n]
	if i := bytes.IndexByte(buf, 0); i != -1 {
		buf = buf[:i]
	}
	return string(buf)
}

// get returns a byte slice of the heap from p to end of allocation.
// The slice is temporary and will be overwritten. Copy it if you want to keep it.
func get(p unsafe.Pointer) []byte {
	h := unsafe.Pointer(&heap[0])
	i := uintptr(p) - uintptr(h)
	return heap[i:heapnext]
}

func FreeTo(prevSize int) {
	assert.That(prevSize <= heapnext)
	zero(heap[prevSize:heapnext])
	heapnext = prevSize
}

func zero(buf []byte) {
	// Zero memory, the compiler should optimize this to memclr.
	// github.com/golang/go/commit/f03c9202c43e0abb130669852082117ca50aa9b1
	for i := range buf {
		buf[i] = 0
	}
}
