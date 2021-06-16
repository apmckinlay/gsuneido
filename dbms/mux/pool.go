// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mux

import "time"

// pool creates worker goroutines on demand.
// Worker goroutines terminate if idle for longer than timeout.
type pool struct {
	ch chan task
	h  workfn
}

type task struct {
	c    *conn
	id   int
	data []byte
}

type workfn func(wb *writeBuffer, id int, data []byte)

// NewPool creates a new goroutine pool
func NewPool(h workfn) *pool {
	ch := make(chan task) // intentionally unbuffered
	return &pool{ch: ch, h: h}
}

// Submit passes a message to a worker
func (p *pool) Submit(c *conn, id int, data []byte) {
	select {
	// try to use an existing worker goroutine
	case p.ch <- task{c: c, id: id, data: data}:
	default:
		// otherwise start a new one
		go p.worker(c, id, data)
	}
}

const timeout = 1 * time.Minute

func (p *pool) worker(c *conn, id int, data []byte) {
	wb := newWriteBuffer(c) // each worker has its own writeBuffer
	p.h(wb, id, data) // do the task
	timer := time.NewTimer(timeout)
	for {
		select {
		case t := <-p.ch:
			if !timer.Stop() {
				<-timer.C
			}
			wb.conn = t.c
			p.h(wb, t.id, t.data) // do the task
			timer.Reset(timeout)
		case <-timer.C:
			return // idle timeout, worker terminates
		}
	}
}
