// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mux

import "github.com/apmckinlay/gsuneido/util/hacks"

const bufSize = 4 * 1024

// writeBuffer is used to combine small writes
type writeBuffer struct {
	*conn
	buf []byte
}

func newWriteBuffer(c *conn) *writeBuffer {
	return &writeBuffer{conn: c, buf: make([]byte, HeaderSize, bufSize)}
}

// space returns the amount of space remaining in the buffer
func (b *writeBuffer) space() int {
	return bufSize - len(b.buf)
}

// Write writes part of a message. If it is small it will be buffered.
// final should be true for the last write of a message.
func (b *writeBuffer) Write(id int, data []byte, final bool) {
	if len(data) > b.space() {
		b.flush(id, false)
	}
	if len(data) >= bufSize {
		b.conn.write(id, data, false, final)
	} else {
		b.buf = append(b.buf, data...)
		if final {
			b.flush(id, true)
		}
	}
}

// WriteString is like Write, but for a string.
func (b *writeBuffer) WriteString(id int, data string, final bool) {
	if len(data) > b.space() {
		b.flush(id, false)
	}
	if len(data) >= bufSize {
		// it would be safer/better to use []byte(s)
		// but strings are used for large data so we want to avoid copying
		b.conn.write(id, hacks.Stobs(data), false, final)
	} else {
		b.buf = append(b.buf, data...)
		if final {
			b.flush(id, true)
		}
	}
}

// WriteByte is like Write, but for a byte.
func (b *writeBuffer) WriteByte(id int, data byte, final bool) {
	if b.space() == 0 {
		b.flush(id, false)
	}
	b.buf = append(b.buf, data)
	if final {
		b.flush(id, true)
	}
}

func (b *writeBuffer) flush(id int, final bool) {
	b.conn.write(id, b.buf, true, final)
	b.buf = b.buf[:HeaderSize]
}
