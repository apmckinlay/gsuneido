// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package mux handles multiple concurrent requests & responses
// over a single connection.
package mux

import (
	"encoding/binary"
	"io"
	"sync"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/trace"
	"github.com/apmckinlay/gsuneido/util/assert"
	myatomic "github.com/apmckinlay/gsuneido/util/generic/atomic"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
)

const HeaderSize = 4 + 4 + 1 /* size + id + final */

type conn struct {
	rw    io.ReadWriteCloser // the underlying connection
	wlock sync.Mutex         // used by write to keep header and data together
	hdr   [HeaderSize]byte   // used by write, guarded by wlock
	err   myatomic.String
}

func (c *conn) Close() {
	// if err := c.err.Load(); err != "" && err != "EOF" {
	// 	log.Println("mux:", err)
	// }
	c.rw.Close()
}

type ClientConn struct {
	conn
	nextSession atomic.Uint32 // the next session id
	lock        sync.Mutex
	rchs        map[uint32]respch // response channel per id, guarded by lock
}

type respch chan []byte

// NewClientConn creates a new client connection.
// This should be one to one with the underlying connection.
func NewClientConn(rw io.ReadWriteCloser) *ClientConn {
	m := ClientConn{conn: conn{rw: rw}, rchs: make(map[uint32]respch)}
	go m.conn.reader(m.client)
	return &m
}

type ServerConn struct {
	conn
	id uint32
}

type handler func(c *conn, id uint64, r []byte)

var nextServerConn atomic.Uint32

// NewServerConn creates a new server connection.
// The supplied handler will be called with each received message.
// For dbms server the handler is Workers.Submit
func NewServerConn(rw io.ReadWriteCloser) *ServerConn {
	connId := nextServerConn.Add(1)
	return &ServerConn{conn: conn{rw: rw}, id: connId}
}

func (sc *ServerConn) Id() uint32 {
	return sc.id
}

func (sc *ServerConn) Run(h handler) {
	sc.conn.reader(func(sessionId uint32, data []byte) {
		h(&sc.conn, uint64(sc.id)<<32|uint64(sessionId), data)
	})
}

type ClientSession struct {
	ReadWrite
	cc  *ClientConn
	rch respch
}

// NewClientSession returns a new ClientSession
func (cc *ClientConn) NewClientSession() *ClientSession {
	sessionId := uint32(cc.nextSession.Add(1))
	rch := make(respch, 1)
	wb := newWriteBuf(&cc.conn, sessionId)
	cc.lock.Lock()
	defer cc.lock.Unlock()
	cc.rchs[sessionId] = rch
	return &ClientSession{cc: cc, rch: rch, ReadWrite: ReadWrite{WriteBuf: *wb}}
}

func (cs *ClientSession) read() []byte {
	return <-cs.rch
}

// Request is used by DbmsClient.
// It does Flush and GetBool for the result.
// If the result is false, it does GetStr for the error and panics with it.
func (cs *ClientSession) Request() {
	cs.EndMsg()
	cs.ReadBuf.buf = cs.read()
	if !cs.GetBool() {
		err := cs.GetStr()
		trace.ClientServer.Println(err)
		panic(err + " (from server)")
	}
}

// write is called by writeBuffer to send part of a message.
// final should be true for the last write of a message.
// If the sender can leave HeaderSize bytes of space at the start of data,
// then it can pass hdrSpace = true, and header & data can be written together.
func (c *conn) write(id uint32, data []byte, hdrSpace, final bool) {
	c.wlock.Lock()
	defer c.wlock.Unlock()
	var err error
	if hdrSpace {
		c.putHdr(data, id, len(data)-HeaderSize, final)
		_, err = c.rw.Write(data)
	} else {
		c.putHdr(c.hdr[:], id, len(data), final)
		_, err = c.rw.Write(c.hdr[:])
		if err == nil {
			_, err = c.rw.Write(data)
		}
	}
	if err != nil {
		c.err.Store(err.Error())
		c.Close()
		return
	}
}

func (*conn) putHdr(buf []byte, id uint32, size int, final bool) {
	binary.BigEndian.PutUint32(buf, uint32(size))
	binary.BigEndian.PutUint32(buf[4:], id)
	if final {
		buf[8] = 1
	} else {
		buf[8] = 0
	}
}

const maxSize = 1024 * 1024 // 1 mb

// reader reads from the connection
// and calls handler when it has a complete message.
// client and server have different handlers.
// Any errors close the connection.
func (c *conn) reader(handler func(uint32, []byte)) {
	partial := make(map[uint32][]byte)
	hdr := make([]byte, HeaderSize)
	for {
		n, err := io.ReadFull(c.rw, hdr)
		if err != nil {
			c.err.Store(err.Error())
			break
		}
		assert.This(n).Is(HeaderSize)
		size := int(binary.BigEndian.Uint32(hdr))
		sessionId := binary.BigEndian.Uint32(hdr[4:])
		buf := partial[sessionId] // nil (empty buf) if not found
		i := len(buf)
		if i+size > maxSize {
			c.err.Store("message size greater than max")
			break
		}
		buf = slc.Grow(buf, size)
		n, err = io.ReadFull(c.rw, buf[i:])
		if err != nil {
			c.err.Store(err.Error())
			break
		}
		assert.This(n).Is(size)
		if hdr[8] == 0 {
			partial[sessionId] = buf
		} else if hdr[8] == 1 {
			delete(partial, sessionId)
			assert.That(buf != nil)
			handler(sessionId, buf) // process message
		} else {
			c.err.Store("bad final byte")
			break
		}
	}
	handler(0, nil) // notify handler
	c.Close()
}

func (cc *ClientConn) client(id uint32, data []byte) {
	// need to send id for client to pipeline messages
	if data == nil {
		runtime.Fatal("lost connection")
	}
	cc.getrch(id) <- data
}

func (cc *ClientConn) getrch(id uint32) respch {
	cc.lock.Lock()
	defer cc.lock.Unlock()
	ch, ok := cc.rchs[id]
	if !ok {
		cc.err.Store("no chan for id")
		cc.Close()
	}
	return ch
}
