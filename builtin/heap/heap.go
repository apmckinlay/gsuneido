// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

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
package heap

import (
	"bytes"
	"log"
	"runtime/debug"
	"unsafe"

	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/verify"
)

const align uintptr = 8 // Alloc assumes power of two
const heapsize = 64 * 1024

var heap = [heapsize]byte{248, 249, 250, 251, 252, 253, 254, 255}
var heapnext = align

var lastAlloc uintptr

// Alloc returns an unsafe.Pointer to n byts of heap space.
func Alloc(n uintptr) unsafe.Pointer {
	i := alloc(n)
	// zero out memory
	for j := uintptr(0); j < n; j++ {
		heap[i + j] = 0
	}
	return unsafe.Pointer(&heap[i])
}

func alloc(n uintptr) uintptr {
	lastAlloc = n
	n = ((n - 1) | (align - 1)) + 1
	if heapnext+n > heapsize {
		panic("Windows dll interface argument space limit exceeded")
	}
	heapcheck("alloc")
	i := heapnext
	heapnext += n
	if options.HeapDebug {
		heapnext += align
		for j := align; j > 0; j-- {
			heap[heapnext-j] = byte(256 - j)
		}
	}
	return i
}

// Copy allocates n bytes on the heap and copies the string into it.
func Copy(s string, n int) unsafe.Pointer {
	i := alloc(uintptr(n))
	copy(heap[i:i+uintptr(n)], s)
	// zero out any remainder (handles nul string terminator)
	for j := uintptr(len(s)); j < uintptr(n); j++ {
		heap[i + j] = 0
	}
	return unsafe.Pointer(&heap[i])
}

func CurSize() uintptr {
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

func FreeTo(prevSize uintptr) {
	heapcheck("free1")
	verify.That(prevSize <= heapnext)
	heapnext = prevSize
	heapcheck("free2")
}

func heapcheck(s string) {
	if options.HeapDebug {
		for i := align; i > 0; i-- {
			if heap[heapnext-i] != byte(256-i) {
				debug.PrintStack()
				log.Fatalln("heap corrupt", s, "lastAlloc", lastAlloc)
			}
		}
	}
}
