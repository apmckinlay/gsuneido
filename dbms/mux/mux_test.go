// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mux

import (
	"bytes"
	"net"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/race"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestMux(t *testing.T) {
	p1, p2 := net.Pipe()
	client := NewClientConn(p1)
	var n atomic.Int32
	ws := NewWorkers(func(wb *WriteBuf, _ *runtime.Thread, id uint64, data []byte) {
		n.Add(1)
		wb.Write(bytes.ToUpper(data)).EndMsg()
	})
	NewServerConn(p2, ws.Submit) // use pool to execute requests
	var nmsgs = 1000
	var nthreads = 11
	if testing.Short() || race.Enabled {
		nmsgs = 100
		nthreads = 5
	}
	var wg sync.WaitGroup
	clientThread := func() {
		session := client.NewClientSession()
		for i := 0; i < nmsgs; i++ {
			a := str.Random(1, 100)
			b := str.Random(1, 2*bufSize)
			session.WriteString(a)
			session.WriteString(b)
			session.EndMsg()
			data := session.read()
			assert.This(string(data)).Is(str.ToUpper(a + b))
		}
		wg.Done()
	}
	for i := 0; i < nthreads; i++ {
		wg.Add(1)
		go clientThread()
	}
	wg.Wait()
	assert.T(t).This(n.Load()).Is(nmsgs * nthreads)
}
