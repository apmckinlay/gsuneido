// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mux

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/dbg"
)

// Workers creates worker goroutines on demand.
// It is not used by mux directly, but can be used by server handlers.
// The number of workers is not limited.
// Unnecessary workers will be terminated.
type Workers struct {
	ch        chan task
	h         workfn
	killClock atomic.Int64
}

type task struct {
	c    *conn
	data []byte
	id   uint64
}

var nWorker atomic.Int64
var _ = core.AddInfo("server.nWorker", &nWorker)

type workfn func(wb *WriteBuf, th *core.Thread, id uint64, rb []byte)

// NewWorkers creates a new goroutine pool
func NewWorkers(h workfn) *Workers {
	ch := make(chan task) // intentionally unbuffered
	ws := &Workers{ch: ch, h: h}
	go ws.killer()
	return ws
}

// Submit passes a task to a worker
func (ws *Workers) Submit(c *conn, id uint64, data []byte) {
	t := task{c: c, id: id, data: data}
	select {
	// try to use an existing worker goroutine
	case ws.ch <- t:
	default:
		// otherwise start a new one
		ws.killClock.Store(0) // prevent immediate kill
		go ws.worker(t)
	}
}

func (ws *Workers) worker(t task) {
	nWorker.Add(1)
	defer nWorker.Add(-1)
	// each worker has its own WriteBuf and Thread
	wb := newWriteBuf(nil, 0)
	th := core.NewThread(nil)
	defer func() {
		if e := recover(); e != nil {
			log.Println("worker panic:", e)
			dbg.PrintStack()
			if se, ok := e.(*core.SuExcept); ok {
				core.PrintStack(se.Callstack)
			} else {
				th.PrintStack()
			}
		}
	}()
	for {
		wb.conn = t.c
		wb.id = uint32(t.id)
		ws.h(wb, th, t.id, t.data) // do the task
		th.Invalidate()
		t = <-ws.ch                // blocking, wait for message
		if t.c == nil {
			return // got poison pill so terminate
		}
		th.Reset()
	}
}

func (ws *Workers) killer() {
	const interval = 2 * time.Second
	const createDelay = 10 // * interval
	const minWorkers = 3
	for {
		time.Sleep(interval)
		if ws.killClock.Add(1) > createDelay && nWorker.Load() > minWorkers {
			ws.ch <- task{} // send poison pill to a worker
		}
	}
}
