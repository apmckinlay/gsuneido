// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package mux handles multiple concurrent requests & responses
// over a single connection.
package mux

import (
	"encoding/binary"
	"errors"
	"io"
	"sync"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/runtime/trace"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/bytes"
)

const HeaderSize = 4 + 4 + 1 /* size + id + final */

type conn struct {
	rw    io.ReadWriteCloser // the underlying connection
	wlock sync.Mutex         // used by write to keep header and data together
	hdr   [HeaderSize]byte   // used by write, guarded by wlock
	err   atomic.Value
}

func (c *conn) Close() {
	c.rw.Close()
}

type clientConn struct {
	conn
	nextSession uint32 // the next session id, must be accessed atomically
	lock        sync.Mutex
	rchs        map[uint32]respch // response channel per id, guarded by lock
}

type respch chan []byte

// NewClientConn creates a new client connection.
// This should be one to one with the underlying connection.
func NewClientConn(rw io.ReadWriteCloser) *clientConn {
	m := clientConn{conn: conn{rw: rw}, rchs: make(map[uint32]respch)}
	go m.conn.reader(m.client)
	return &m
}

type serverConn struct {
	conn
	id uint32
}

type handler func(c *conn, id uint64, r []byte)

var nextServerConn uint32

// NewServerConn creates a new server connection.
// The supplied handler will be called with each received message.
func NewServerConn(rw io.ReadWriteCloser, h handler) *serverConn {
	sid := atomic.AddUint32(&nextServerConn, 1)
	sc := serverConn{conn: conn{rw: rw}, id: sid}
	go sc.conn.reader(func(cid uint32, data []byte) {
		h(&sc.conn, uint64(sid)<<32|uint64(cid), data)
	})
	return &sc
}

type ClientSession struct {
	ReadWrite
	cc  *clientConn
	rch respch
}

// NewClientSession returns a new ClientSession
func (cc *clientConn) NewClientSession() *ClientSession {
	cid := atomic.AddUint32(&cc.nextSession, 1)
	rch := make(respch, 1)
	wb := newWriteBuf(&cc.conn, cid)
	cc.lock.Lock()
	defer cc.lock.Unlock()
	cc.rchs[cid] = rch
	return &ClientSession{cc: cc, rch: rch,
		ReadWrite: ReadWrite{writeBuf: *wb}}
}

func (cs *ClientSession) read() []byte {
	return <-cs.rch
}

// Request is used by DbmsClient.
// It does Flush and GetBool for the result.
// If the result is false, it does GetStr for the error and panics with it.
func (cs *ClientSession) Request() {
	cs.EndMsg()
	cs.readBuf.buf = cs.read()
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
		c.err.Store(err)
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
			c.err.Store(err)
			break
		}
		assert.This(n).Is(HeaderSize)
		size := int(binary.BigEndian.Uint32(hdr))
		id := binary.BigEndian.Uint32(hdr[4:])
		buf := partial[id] // nil (empty buf) if not found
		i := len(buf)
		if i+size > maxSize {
			c.err.Store(errors.New("message size greater than max"))
			break
		}
		buf = bytes.Grow(buf, size)
		n, err = io.ReadFull(c.rw, buf[i:])
		if err != nil {
			c.err.Store(err)
			break
		}
		assert.This(n).Is(size)
		if hdr[8] == 0 {
			partial[id] = buf
		} else if hdr[8] == 1 {
			delete(partial, id)
			handler(id, buf) // process message
		} else {
			c.err.Store(errors.New("bad final byte"))
			break
		}
	}
	c.Close()
}

func (cc *clientConn) client(id uint32, data []byte) {
	// need to send id for client to pipeline messages
	cc.getrch(id) <- data
}

func (cc *clientConn) getrch(id uint32) respch {
	cc.lock.Lock()
	defer cc.lock.Unlock()
	ch, ok := cc.rchs[id]
	if !ok {
		cc.err.Store(errors.New("no chan for id"))
		cc.Close()
	}
	return ch
}
