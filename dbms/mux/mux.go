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

type clientConn struct {
	conn
	next int32 // the next message id, must be accessed atomically
	lock sync.Mutex
	rchs map[int]respch // response channel per id, guarded by lock
}

type respch chan []byte

// ClientConn creates a new client connection.
// This should be one to one with the underlying connection.
func ClientConn(rw io.ReadWriteCloser) *clientConn {
	m := clientConn{conn: conn{rw: rw}, rchs: make(map[int]respch)}
	go m.conn.reader(m.client)
	return &m
}

// WriteBuffer creates a write buffer for each client thread
func (cc *clientConn) WriteBuffer() *writeBuffer {
	return newWriteBuffer(&cc.conn)
}

type serverConn struct {
	conn
}

type handler func(c *conn, id int, r []byte)

// ServerConn creates a new server connection.
// The supplied handler will be called with each received message.
func ServerConn(rw io.ReadWriteCloser, h handler) *serverConn {
	sc := serverConn{conn: conn{rw: rw}}
	go sc.conn.reader(func(id int, data []byte) {
		h(&sc.conn, id, data)
	})
	return &sc
}

// NewRequest returns a new message id,
// registering the supplied channel for the response.
func (cc *clientConn) NewRequest(rch chan []byte) int {
	id := int(atomic.AddInt32(&cc.next, 1))
	cc.lock.Lock()
	defer cc.lock.Unlock()
	cc.rchs[id] = rch
	return id
}

func (c *conn) Close() {
	c.rw.Close()
}

// write is called by writeBuffer to send part of a message.
// The id should come from NewRequest.
// final should be true for the last write of a message.
// If the sender can leave HeaderSize bytes of space at the start of data,
// then it can pass hdrSpace = true, and header & data can be written together.
func (c *conn) write(id int, data []byte, hdrSpace, final bool) {
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

func (*conn) putHdr(buf []byte, id, size int, final bool) {
	binary.BigEndian.PutUint32(buf, uint32(size))
	binary.BigEndian.PutUint32(buf[4:], uint32(id))
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
func (c *conn) reader(handler func(int, []byte)) {
	partial := make(map[int][]byte)
	hdr := make([]byte, HeaderSize)
	for {
		n, err := io.ReadFull(c.rw, hdr)
		if err != nil {
			c.err.Store(err)
			break
		}
		assert.This(n).Is(HeaderSize)
		size := int(binary.BigEndian.Uint32(hdr))
		id := int(binary.BigEndian.Uint32(hdr[4:]))
		buf := partial[id]
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

func (cc *clientConn) client(id int, r []byte) {
	cc.getrch(id) <- r
}

func (cc *clientConn) getrch(id int) respch {
	cc.lock.Lock()
	defer cc.lock.Unlock()
	ch, ok := cc.rchs[id]
	if !ok {
		cc.err.Store(errors.New("no chan for id"))
		cc.Close()
	}
	return ch
}
