// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package queue is a concurrent size limited fifo queue.
// It is equivalent to a buffered channel
// except that we have access to the elements.
package queue

import (
	"fmt"
	"log"
	"slices"
	"sync"
)

type Queue[T any] struct {
	size  int
	limit int
	items []T
	lock  sync.Mutex
	cond  sync.Cond
}

// New creates a new queue with the given size.
func New[T any](size, limit int) *Queue[T] {
	q := &Queue[T]{size: size, limit: limit}
	q.cond.L = &q.lock
	return q
}

// Put adds an item to the queue.
// It blocks if the queue is full.
func (q *Queue[T]) Put(item T) {
	q.lock.Lock()
	defer q.lock.Unlock()
	for len(q.items) >= q.limit {
		// queue is full
		q.cond.Wait()
	}
	q.items = append(q.items, item)
}

// MustPut adds an item to the queue.
// It logs if the queue reaches size and panics if it reaches limit.
func (q *Queue[T]) MustPut(item T) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.items) == q.size {
		log.Println("WARNING Queue: over", q.size, "items")
	}
	if len(q.items) > q.limit {
		panic(fmt.Sprint("ERROR Queue over ", q.limit, " items"))
	}
	q.items = append(q.items, item)
}

// TryGet returns the first item in the queue if one is available.
// It does not block.
func (q *Queue[T]) TryGet() (item T, ok bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.items) == 0 {
		return
	}
	item = popfirst(&q.items)
	if len(q.items) < q.limit {
		q.cond.Signal()
	}
	return item, true
}

// Remove removes the first (oldest) item in the queue that satisfies f.
func (q *Queue[T]) Remove(f func(T) bool) bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	for i, it := range q.items {
		if f(it) {
			q.items = slices.Delete(q.items, i, i+1)
			q.cond.Signal()
			return true
		}
	}
	return false
}

func popfirst[T any](x *[]T) T {
	it := (*x)[0]
	*x = (*x)[:copy(*x, (*x)[1:])]
	return it
}
