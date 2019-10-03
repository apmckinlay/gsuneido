// Package heapstack provides a heap for win32 dll arguments
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
package heapstack

import (
	"fmt"
	"unsafe"
)

const debug = true

func init() {
	if debug {
		fmt.Println("heap debug enabled")
	}
}

const align uintptr = 8 // Alloc assumes power of two
const heapsize = 64 * 1024

var heap = [heapsize]byte{248, 249, 250, 251, 252, 253, 254, 255}
var heapnext = align

func Alloc(n uintptr) unsafe.Pointer {
	n = ((n - 1) | (align - 1)) + 1
	heapcheck("alloc")
	// zero out memory
	// probably not required ???
	for i := uintptr(0); i < n; i++ {
		heap[heapnext+i] = 0
	}
	p := &heap[heapnext]
	heapnext += n
	if debug {
		heapnext += align
		for i := align; i > 0; i-- {
			heap[heapnext-i] = byte(256 - i)
		}
	}
	return unsafe.Pointer(p)
}

func CurSize() uintptr {
	return heapnext
}

func FreeTo(prevSize uintptr) {
	heapcheck("free1")
	heapnext = prevSize
	heapcheck("free2")
}

func heapcheck(s string) {
	if debug {
		for i := align; i > 0; i-- {
			if heap[heapnext-i] != byte(256-i) {
				panic("heap corrupt " + s)
			}
		}
	}
}
