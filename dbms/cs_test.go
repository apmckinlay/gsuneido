// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !gui

package dbms

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms/mux"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestClientServer(*testing.T) {
	// trace.Set(int(trace.ClientServer))
	Global.TestDef("Auth",
		compile.Constant("function (@args) { return true }"))
	defer Global.UnloadAll()
	options.BuiltDate = "Dec 29 2020 12:34"
	db := db19.CreateDb(stor.HeapStor(8192))
	dbmsLocal := NewDbmsLocal(db)
	p1, p2 := net.Pipe()
	workers = mux.NewWorkers(doRequest)
	go newServerConn(dbmsLocal, p1, serverTLSConfig())
	// Exchange hello over plain connection
	errmsg := checkHello(p2)
	assert.This(errmsg).Is("")
	p2.Write(hello())
	// Upgrade client side to TLS
	tlsConn := tls.Client(p2, clientTLSConfig())
	if err := tlsConn.Handshake(); err != nil {
		panic(err)
	}
	c := NewDbmsClient(tlsConn)
	ses := c.NewSession()
	ses.Auth(&Thread{}, &SuObject{})
	args := SuObjectOf(SuStr("tables sort table"))
	ses.Get(nil, args, Next)

	ses2 := c.NewSession()
	ses2.Get(nil, args, Prev)
	ses2.Close()

	time.Sleep(25 * time.Millisecond)
}

func TestServerEvalWhitelist(t *testing.T) {
	assert := assert.T(t)
	// a Go Auth that sets up a ServerEval whitelist
	// since the Perm builtin is not available in dbms package tests
	Global.TestDef("Auth", &SuBuiltinRaw{
		Fn: func(th *Thread, as *ArgSpec, args []Value) Value {
			th.Perms().AddServerEval("F1")
			return True
		},
		ParamSpec: ParamSpecAt})
	Global.TestDef("F1", compile.Constant("function (@args) { return 123 }"))
	Global.TestDef("F2", compile.Constant("function (@args) { return 456 }"))
	defer Global.UnloadAll()
	options.BuiltDate = "Dec 29 2020 12:34"
	db := db19.CreateDb(stor.HeapStor(8192))
	dbmsLocal := NewDbmsLocal(db)
	p1, p2 := net.Pipe()
	workers = mux.NewWorkers(doRequest)
	go newServerConn(dbmsLocal, p1, serverTLSConfig())
	// Exchange hello over plain connection
	errmsg := checkHello(p2)
	assert.This(errmsg).Is("")
	p2.Write(hello())
	// Upgrade client side to TLS
	tlsConn := tls.Client(p2, clientTLSConfig())
	if err := tlsConn.Handshake(); err != nil {
		panic(err)
	}
	c := NewDbmsClient(tlsConn)
	ses := c.NewSession()
	assert.True(ses.Auth(&Thread{}, &SuObject{}))

	// allowed by the whitelist
	assert.This(ToInt(ses.Exec(nil, SuObjectOf(SuStr("F1"))))).Is(123)
	// not allowed by the whitelist, even though the function exists
	assert.This(func() { ses.Exec(nil, SuObjectOf(SuStr("F2"))) }).
		Panics("ServerEval: not permitted: F2")

	time.Sleep(25 * time.Millisecond)
}

func TestTLSPinning(t *testing.T) {
	assert := assert.T(t)
	// with the expected certs on both ends, the connection is established
	assert.That(tryHandshake(serverTLSConfig(), clientTLSConfig()))
	// a client without a client certificate (i.e. built before mutual
	// pinning was added) is rejected by the server
	assert.That(!tryHandshake(serverTLSConfig(), oldClientConfig()))
	// a client with a certificate the server does not trust is rejected
	serverPair, err := tls.X509KeyPair(ServerCert, ServerKey)
	assert.That(err == nil)
	assert.That(!tryHandshake(serverTLSConfig(), &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{serverPair},
	}))
	// a server with a certificate the client does not trust is rejected
	clientPair, err := tls.X509KeyPair(ClientCert, ClientKey)
	assert.That(err == nil)
	assert.That(!tryHandshake(&tls.Config{
		Certificates: []tls.Certificate{clientPair},
	}, clientTLSConfig()))
}

// oldClientConfig is like clientTLSConfig but without a client certificate
func oldClientConfig() *tls.Config {
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(ServerCert)
	if !ok {
		panic("failed to append embedded cert to pool")
	}
	return &tls.Config{
		RootCAs:    caCertPool,
		ServerName: "localhost",
	}
}

// tryHandshake attempts a TLS handshake over TCP loopback.
// It returns true only if both ends accept the connection.
// Both ends must be checked because, with TLS 1.3, a client can finish
// its own handshake before the server has validated its certificate.
// It uses TCP rather than net.Pipe so a rejected handshake can't deadlock
// on writing the alert.
func tryHandshake(serverConfig, clientConfig *tls.Config) bool {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.That(err == nil)
	defer l.Close()
	srvErr := make(chan error, 1)
	go func() {
		conn, err := l.Accept()
		if err != nil {
			srvErr <- err
			return
		}
		sconn := tls.Server(conn, serverConfig)
		srvErr <- sconn.Handshake()
		sconn.Close()
	}()
	conn, err := net.Dial("tcp", l.Addr().String())
	assert.That(err == nil)
	cconn := tls.Client(conn, clientConfig)
	err = cconn.Handshake()
	serr := <-srvErr
	cconn.Close()
	return err == nil && serr == nil
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
