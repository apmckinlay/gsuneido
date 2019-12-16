// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package str

type Queue struct {
	list []string
}

// Add adds a string to the end of the queue.
func (q *Queue) Add(s string) {
	q.list = append(q.list, s)
}

// Take removes and returns the first string in the queue (FIFO).
// Will panic if queue is empty.
func (q *Queue) Take() string {
	s := q.list[0]
	n := len(q.list)
	copy(q.list, q.list[1:])
	q.list[n-1] = "" // for gc
	q.list = q.list[:n-1]
	return s
}

// Empty returns true is the queue is empty, otherwise false.
func (q *Queue) Empty() bool {
	return len(q.list) == 0
}
