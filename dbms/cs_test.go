// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms/mux"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestClientServer(*testing.T) {
	// trace.Set(int(trace.ClientServer))
	options.BuiltDate = "Dec 29 2020 12:34"
	db := db19.CreateDb(stor.HeapStor(8192))
	dbmsLocal := NewDbmsLocal(db)
	p1, p2 := net.Pipe()
	workers = mux.NewWorkers(doRequest)
	go newServerConn(dbmsLocal, p1)
	errmsg := checkHello(p2)
	assert.This(errmsg).Is("")
	p2.Write(hello())
	c := NewDbmsClient(p2)
	ses := c.NewSession()
	ses.Get(nil, "tables sort table", core.Next)

	ses2 := c.NewSession()
	ses2.Get(nil, "tables sort table", core.Prev)
	ses2.Close()

	time.Sleep(25 * time.Millisecond)
}

var A atomic.Bool
var M sync.Mutex

func BenchmarkOne(b *testing.B) {
	for b.Loop() {
		func() {
			if A.Load() {
				M.Lock()
				defer M.Unlock()
			}
		}()
	}
}

func BenchmarkTwo(b *testing.B) {
	for b.Loop() {
		func() {
			M.Lock()
			defer M.Unlock()
		}()
	}
}
