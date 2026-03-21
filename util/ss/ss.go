// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package ss implements the Space-Saving algorithm for string streams.
//
// It tracks approximate frequent items using fixed memory.
// See: Efficient Computation of Frequent and Top-k Elements in Data Streams
// https://www.cs.ucsb.edu/sites/default/files/documents/2005-23.pdf
package ss

import (
	"container/heap"
	"slices"
)

// Entry is one tracked item.
//
// The true count is in [Count-Error, Count].
type Entry struct {
	Value string
	Count int
	Error int
	idx   int
}

// Sketch tracks approximate frequencies of strings using fixed capacity.
type Sketch struct {
	capacity int
	total    int
	entries  minHeap
	index    map[string]*Entry
}

// New creates a Space-Saving sketch that tracks up to capacity items.
func New(capacity int) *Sketch {
	if capacity <= 0 {
		panic("ss capacity must be > 0")
	}
	return &Sketch{
		capacity: capacity,
		entries:  make(minHeap, 0, capacity),
		index:    make(map[string]*Entry, capacity),
	}
}

// Add inserts one value from the stream.
func (ss *Sketch) Add(value string) {
	ss.total++
	if e, ok := ss.index[value]; ok {
		e.Count++
		heap.Fix(&ss.entries, e.idx)
		return
	}

	if len(ss.entries) < ss.capacity {
		e := &Entry{Value: value, Count: 1}
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
func (ss *Sketch) Count() int {
	return ss.total
}

// Capacity returns the maximum number of tracked items.
func (ss *Sketch) Capacity() int {
	return ss.capacity
}

// Len returns the current number of tracked items.
func (ss *Sketch) Len() int {
	return len(ss.entries)
}

// Estimate returns the tracked estimate for value.
//
// When ok is true, the true count is in [count-error, count].
func (ss *Sketch) Estimate(value string) (count int, error int, ok bool) {
	e, ok := ss.index[value]
	if !ok {
		return 0, 0, false
	}
	return e.Count, e.Error, true
}

// Top returns the tracked entries sorted by descending count.
func (ss *Sketch) Top() []Entry {
	top := make([]Entry, len(ss.entries))
	for i, e := range ss.entries {
		top[i] = Entry{Value: e.Value, Count: e.Count, Error: e.Error}
	}
	slices.SortFunc(top, func(a, b Entry) int {
		if a.Count > b.Count {
			return -1
		}
		if a.Count < b.Count {
			return 1
		}
		if a.Value < b.Value {
			return -1
		}
		if a.Value > b.Value {
			return 1
		}
		return 0
	})
	return top
}

type minHeap []*Entry

func (h minHeap) Len() int { return len(h) }

func (h minHeap) Less(i, j int) bool {
	if h[i].Count != h[j].Count {
		return h[i].Count < h[j].Count
	}
	return h[i].Value < h[j].Value
}

func (h minHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].idx = i
	h[j].idx = j
}

func (h *minHeap) Push(x any) {
	e := x.(*Entry)
	e.idx = len(*h)
	*h = append(*h, e)
}

func (h *minHeap) Pop() any {
	n := len(*h)
	e := (*h)[n-1]
	e.idx = -1
	*h = (*h)[:n-1]
	return e
}
