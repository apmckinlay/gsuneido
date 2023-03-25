// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mux

import (
	"time"

	"github.com/apmckinlay/gsuneido/runtime"
)

// Workers creates worker goroutines on demand.
// It is not used by mux directly, but can be used by server handlers.
// Worker goroutines terminate if idle for longer than timeout.
type Workers struct {
	ch chan task
	h  workfn
}

type task struct {
	c    *conn
	data []byte
	id   uint64
}

type workfn func(wb *WriteBuf, th *runtime.Thread, id uint64, rb []byte)

// NewWorkers creates a new goroutine pool
func NewWorkers(h workfn) *Workers {
	ch := make(chan task) // intentionally unbuffered
	return &Workers{ch: ch, h: h}
}

// Submit passes a task to a worker
func (ws *Workers) Submit(c *conn, id uint64, data []byte) {
	t := task{c: c, id: id, data: data}
	select {
	// try to use an existing worker goroutine
	case ws.ch <- t:
	default:
		// otherwise start a new one
		go ws.worker(t)
	}
}

const timeout = 5 * time.Second //1 * time.Minute // ???

func (ws *Workers) worker(t task) {
	// each worker has its own WriteBuf and Thread
	wb := newWriteBuf(nil, 0)
	th := &runtime.Thread{}
	timer := time.NewTimer(timeout)
	for {
		wb.conn = t.c
		wb.id = uint32(t.id)
		ws.h(wb, th, t.id, t.data) // do the task
		select {
		case t = <-ws.ch:
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(timeout)
		case <-timer.C:
			return // idle timeout, worker terminates
		}
		th.Reset()
	}
}
