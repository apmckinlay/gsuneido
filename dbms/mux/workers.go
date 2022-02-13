// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mux

import "time"

// workers creates worker goroutines on demand.
// It is not used by mux directly, but can be used by server handlers.
// Worker goroutines terminate if idle for longer than timeout.
type workers struct {
	ch chan task
	h  workfn
}

type task struct {
	c    *conn
	id   uint64
	data []byte
}

type workfn func(wb *writeBuf, rb []byte)

// NewWorkers creates a new goroutine pool
func NewWorkers(h workfn) *workers {
	ch := make(chan task) // intentionally unbuffered
	return &workers{ch: ch, h: h}
}

// Submit passes a task to a worker
func (ws *workers) Submit(c *conn, id uint64, data []byte) {
	select {
	// try to use an existing worker goroutine
	case ws.ch <- task{c: c, id: id, data: data}:
	default:
		// otherwise start a new one
		go ws.worker(c, id, data)
	}
}

const timeout = 1 * time.Minute

func (ws *workers) worker(c *conn, id uint64, data []byte) {
	wb := newWriteBuf(c, uint32(id)) // each worker has its own writeBuffer
	ws.h(wb, data)       // do the task
	timer := time.NewTimer(timeout)
	for {
		select {
		case t := <-ws.ch:
			if !timer.Stop() {
				<-timer.C
			}
			wb.conn = t.c
			wb.id = uint32(t.id)
			ws.h(wb, t.data) // do the task
			timer.Reset(timeout)
		case <-timer.C:
			return // idle timeout, worker terminates
		}
	}
}
