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
// The true count is in the range [Count-Error, Count].
type Entry[T comparable] struct {
	Value T
	idx   int
	Count int
	Error int
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
		panic("ss capacity must be greater than zero")
	}
	return &Sketch[T]{
		capacity: capacity * 2, // for more accuracy
		entries:  make(minHeap[T], 0, capacity),
		index:    make(map[T]*Entry[T], capacity),
	}
}

// Add inserts one value from the stream (weight 1).
func (ss *Sketch[T]) Add(value T) {
	ss.AddWeight(value, 1)
}

// AddWeight inserts a value with the given weight.
func (ss *Sketch[T]) AddWeight(value T, weight int) {
	if weight <= 0 {
		panic("ss weight must be greater than zero")
	}
	ss.total += weight
	if e, ok := ss.index[value]; ok {
		e.Count += weight
		heap.Fix(&ss.entries, e.idx)
		return
	}

	if len(ss.entries) < ss.capacity {
		e := &Entry[T]{Value: value, Count: weight}
		heap.Push(&ss.entries, e)
		ss.index[value] = e
		return
	}

	min := ss.entries[0]
	delete(ss.index, min.Value)

	min.Value = value
	min.Error = min.Count
	min.Count += weight
	ss.index[value] = min
	heap.Fix(&ss.entries, min.idx)
}

// Count returns the total weight added.
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
// When ok is true, the true count is in the range [count-error, count].
func (ss *Sketch[T]) Estimate(value T) (count int, error int, ok bool) {
	e, ok := ss.index[value]
	if !ok {
		return 0, 0, false
	}
	return e.Count, e.Error, true
}

// Top returns the top entries
// pruned by (count - error), sorted by descending count.
func (ss *Sketch[T]) Top() []*Entry[T] {
	return ss.TopMin(0)
}

// TopMin returns the top entries with frequency >= minFreq,
// pruned by (count - error), sorted by descending count.
//
// It must not modify the sketch - results have to be consistent
// if it is called more than once (e.g. by PackSize and then Pack).
func (ss *Sketch[T]) TopMin(minFreq float64) []*Entry[T] {
	// sort a clone rather than reordering (and having to restore) the heap
	entries := slices.Clone(ss.entries)
	slices.SortFunc(entries, func(a, b *Entry[T]) int {
		return -cmp.Compare(a.minCount(), b.minCount())
	})
	minCount := int(minFreq * float64(ss.total))
	maxLen := min(ss.capacity/2, len(entries))
	topLen := 0
	for topLen < maxLen && entries[topLen].minCount() >= minCount {
		topLen++
	}
	top := entries[:topLen]
	slices.SortFunc(top, func(a, b *Entry[T]) int {
		return -cmp.Compare(a.Count, b.Count)
	})
	return top
}

func (e *Entry[T]) minCount() int {
	// Count should always be >= Error
	return e.Count - e.Error
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
	*h = (*h)[:n-1]
	return e
}
