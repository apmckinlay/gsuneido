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
	"slices"
)

// Entry is one tracked item.
//
// The true count is in [Count-Error, Count].
type Entry[T comparable] struct {
	Value T
	Count int
	Error int
	idx   int
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
	if capacity <= 0 {
		panic("ss capacity must be > 0")
	}
	return &Sketch[T]{
		capacity: capacity,
		entries:  make(minHeap[T], 0, capacity),
		index:    make(map[T]*Entry[T], capacity),
	}
}

// Add inserts one value from the stream.
func (ss *Sketch[T]) Add(value T) {
	ss.total++
	if e, ok := ss.index[value]; ok {
		e.Count++
		heap.Fix(&ss.entries, e.idx)
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
	heap.Fix(&ss.entries, min.idx)
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
	return e.Count, e.Error, true
}

// Top returns the tracked entries sorted by descending count.
func (ss *Sketch[T]) Top() []Entry[T] {
	top := make([]Entry[T], len(ss.entries))
	for i, e := range ss.entries {
		top[i] = Entry[T]{Value: e.Value, Count: e.Count, Error: e.Error}
	}
	slices.SortFunc(top, func(a, b Entry[T]) int {
		return -cmp.Compare(a.Count, b.Count)
	})
	return top
}

type minHeap[T comparable] []*Entry[T]

func (h minHeap[T]) Len() int { return len(h) }

func (h minHeap[T]) Less(i, j int) bool {
	return h[i].Count < h[j].Count
}

func (h minHeap[T]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].idx = i
	h[j].idx = j
}

func (h *minHeap[T]) Push(x any) {
	e := x.(*Entry[T])
	e.idx = len(*h)
	*h = append(*h, e)
}

func (h *minHeap[T]) Pop() any {
	n := len(*h)
	e := (*h)[n-1]
	e.idx = -1
	*h = (*h)[:n-1]
	return e
}
