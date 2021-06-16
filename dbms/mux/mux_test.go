// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mux

import (
	"bytes"
	"net"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/race"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestMux(t *testing.T) {
	p1, p2 := net.Pipe()
	client := ClientConn(p1)
	n := int32(0)
	p := NewPool(func(w *writeBuffer, id int, data []byte) {
		atomic.AddInt32(&n, 1)
		w.Write(id, bytes.ToUpper(data), true)
	})
	server := ServerConn(p2, p.Submit)
	var nmsgs = 1000
	var nthreads = 11
	if testing.Short() || race.Enabled {
		nmsgs = 100
		nthreads = 5
	}
	var wg sync.WaitGroup
	thread := func() {
		cw := client.WriteBuffer() // each client thread has its own buffer
		ch := make(chan []byte, 1) // and its own result channel
		for i := 0; i < nmsgs; i++ {
			a := str.Random(1, 100)
			b := str.Random(1, 2*bufSize)
			id := client.NewRequest(ch)
			cw.WriteString(id, a, false)
			cw.WriteString(id, b, true)
			x := <-ch
			assert.This(string(x)).Is(str.ToUpper(a + b))
		}
		wg.Done()
	}
	for i := 0; i < nthreads; i++ {
		wg.Add(1)
		go thread()
	}
	wg.Wait()
	assert.T(t).This(atomic.LoadInt32(&n)).Is(nmsgs * nthreads)
	client.Close()
	server.Close()
}
