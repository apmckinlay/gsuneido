// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mux

import (
	"bytes"
	"net"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/race"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestMux(t *testing.T) {
	p1, p2 := net.Pipe()
	client := NewClientConn(p1)
	var n atomic.Int32
	workers := NewWorkers(func(wb *WriteBuf, _ *core.Thread, id uint64, data []byte) {
		n.Add(1)
		wb.Write(bytes.ToUpper(data)).EndMsg()
	})
	msc := NewServerConn(p2)
	go msc.Run(workers.Submit)
	nmsgs := 1000
	nthreads := 11
	if testing.Short() || race.Enabled {
		nmsgs = 100
		nthreads = 5
	}
	var wg sync.WaitGroup
	clientThread := func() {
		session := client.NewClientSession()
		for range nmsgs {
			a := str.Random(1, 100)
			b := str.Random(1, 2*bufSize)
			session.WriteString(a)
			session.WriteString(b)
			session.EndMsg()
			data := session.read()
			assert.This(string(data)).Is(str.ToUpper(a + b))
		}
	}
	for range nthreads {
		wg.Go(clientThread)
	}
	wg.Wait()
	assert.T(t).This(n.Load()).Is(nmsgs * nthreads)
}
