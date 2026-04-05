// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package ss implements the Space-Saving algorithm for streams.
//
// It tracks approximate frequent items using fixed memory.
// See: Efficient Computation of Frequent and Top-k Elements in Data Streams
// https://www.cs.ucsb.edu/sites/default/files/documents/2005-23.pdf
package ss

import (
	"cmp"
	"container/heap"
	"math"
	"slices"
)

// Entry is one tracked item.
//
// The true count is in the range [Count-Error, Count].
type Entry[T comparable] struct {
	// layout and types are optimized for T being uint16 or string
	Value T
	idx   uint16
	Count uint32
	Error uint32
}

// Sketch tracks approximate frequencies using fixed capacity.
type Sketch[T comparable] struct {
	capacity int
	total    int
	entries  minHeap[T]
	index    map[T]*Entry[T]
}

// New creates a Space-Saving sketch that tracks up to capacity items.
func New[T comparable](capacity int) *Sketch[T] {
	if capacity <= 0 || capacity > math.MaxUint16 {
		panic("ss capacity must be from 0 to 64k")
	}
	return &Sketch[T]{
		capacity: capacity * 2, // for more accuracy
		entries:  make(minHeap[T], 0, capacity),
		index:    make(map[T]*Entry[T], capacity),
	}
}

// Add inserts one value from the stream.
func (ss *Sketch[T]) Add(value T) {
	ss.total++
	if e, ok := ss.index[value]; ok {
		e.Count++
		heap.Fix(&ss.entries, int(e.idx))
		return
	}

	if len(ss.entries) < ss.capacity {
		e := &Entry[T]{Value: value, Count: 1}
		heap.Push(&ss.entries, e)
		ss.index[value] = e
		return
	}

	min := ss.entries[0]
	delete(ss.index, min.Value)

	min.Value = value
	min.Error = min.Count
	min.Count++
	ss.index[value] = min
	heap.Fix(&ss.entries, int(min.idx))
}

// Count returns the number of Add calls.
func (ss *Sketch[T]) Count() int {
	return ss.total
}

// Capacity returns the maximum number of tracked items.
func (ss *Sketch[T]) Capacity() int {
	return ss.capacity
}

// Len returns the current number of tracked items.
func (ss *Sketch[T]) Len() int {
	return len(ss.entries)
}

// Estimate returns the tracked estimate for value.
//
// When ok is true, the true count is in [count-error, count].
func (ss *Sketch[T]) Estimate(value T) (count int, error int, ok bool) {
	e, ok := ss.index[value]
	if !ok {
		return 0, 0, false
	}
	return int(e.Count), int(e.Error), true
}

// Top returns the top entries
// pruned by (count - error), sorted by descending count.
func (ss *Sketch[T]) Top() []*Entry[T] {
	// temporarily reorder the entries (breaking the heap)
	slices.SortFunc(ss.entries, func(a, b *Entry[T]) int {
		return -cmp.Compare(a.Count-a.Error, b.Count-a.Error)
	})
	topLen := min(ss.capacity/2, len(ss.entries))
	top := slices.Clone(ss.entries[:topLen])
	slices.SortFunc(top, func(a, b *Entry[T]) int {
		return -cmp.Compare(a.Count, b.Count)
	})
	heap.Init(&ss.entries) // restore heap
	return top
}

type minHeap[T comparable] []*Entry[T]

func (h minHeap[T]) Len() int { return len(h) }

func (h minHeap[T]) Less(i, j int) bool {
	return h[i].Count < h[j].Count
}

func (h minHeap[T]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].idx = uint16(i)
	h[j].idx = uint16(j)
}

func (h *minHeap[T]) Push(x any) {
	e := x.(*Entry[T])
	e.idx = uint16(len(*h))
	*h = append(*h, e)
}

func (h *minHeap[T]) Pop() any {
	n := len(*h)
	e := (*h)[n-1]
	*h = (*h)[:n-1]
	return e
}
