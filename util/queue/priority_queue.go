// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package queue

import (
	"slices"
	"sync"
)

const bufSize = 8

type element struct {
	priority int
	tran     int
	value    any
}

// PriorityQueue is a bounded, blocking priority queue.
// Put blocks when the buffer is full; Get blocks when empty.
// Among the oldest items per transaction, Get returns the one with the highest priority.
// This maintains FIFO order within each transaction.
type PriorityQueue struct {
	items    []element
	lock     sync.Mutex
	notFull  sync.Cond
	notEmpty sync.Cond
}

// NewPriorityQueue creates a new priority queue.
func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{
		items: make([]element, 0, bufSize),
	}
	pq.notFull.L = &pq.lock
	pq.notEmpty.L = &pq.lock
	return pq
}

// Put adds an element to the queue.
// It blocks if the queue is full.
func (pq *PriorityQueue) Put(priority, tran int, value any) {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	for len(pq.items) >= bufSize {
		pq.notFull.Wait()
	}

	pq.items = append(pq.items, element{priority, tran, value})
	pq.notEmpty.Signal()
}

// Get returns the highest priority element among the oldest elements per transaction.
// It blocks if the queue is empty.
func (pq *PriorityQueue) Get() any {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	for len(pq.items) == 0 {
		pq.notEmpty.Wait()
	}

	// Initialize with first item (guaranteed oldest)
	bestIdx := 0
	bestPriority := pq.items[0].priority
	for i := 1; i < len(pq.items); i++ {
		e := &pq.items[i]
		if e.priority > bestPriority && pq.isOldest(i, e) {
			bestIdx = i
			bestPriority = e.priority
		}
	}

	result := pq.items[bestIdx].value
	pq.items = slices.Delete(pq.items, bestIdx, bestIdx+1)
	pq.notFull.Signal()

	return result
}

func (pq *PriorityQueue) isOldest(i int, e *element) bool {
	for j := range i {
		if pq.items[j].tran == e.tran {
			return false
		}
	}
	return true
}
