// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !gui

package dbms

import (
	"crypto/tls"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
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
	cert, err := tls.X509KeyPair(ServerCert, ServerKey)
	if err != nil {
		panic(err)
	}
	serverConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	go newServerConn(dbmsLocal, p1, serverConfig)
	// Exchange hello over plain connection
	errmsg := checkHello(p2)
	assert.This(errmsg).Is("")
	p2.Write(hello())
	// Upgrade client side to TLS
	clientConfig := &tls.Config{
		InsecureSkipVerify: true, // For testing with self-signed cert
	}
	tlsConn := tls.Client(p2, clientConfig)
	if err := tlsConn.Handshake(); err != nil {
		panic(err)
	}
	c := NewDbmsClient(tlsConn)
	ses := c.NewSession()
	args := SuObjectOf(SuStr("tables sort table"))
	ses.Get(nil, args, Next)

	ses2 := c.NewSession()
	ses2.Get(nil, args, Prev)
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
