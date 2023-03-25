// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mux

import (
	"fmt"
	"math"

	"github.com/apmckinlay/gsuneido/dbms/commands"
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/trace"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

const bufSize = 4 * 1024

// WriteBuf is used to combine small writes
type WriteBuf struct {
	*conn
	buf []byte
	id  uint32
}

func newWriteBuf(c *conn, id uint32) *WriteBuf {
	return &WriteBuf{conn: c, id: id, buf: make([]byte, HeaderSize, bufSize)}
}

func (wb *WriteBuf) Id() uint32 {
	return wb.id
}

// space returns the amount of space remaining in the buffer
func (wb *WriteBuf) space() int {
	return bufSize - len(wb.buf)
}

// Write writes part of a message. If it is small it will be buffered.
// final should be true for the last write of a message.
func (wb *WriteBuf) Write(data []byte) *WriteBuf {
	if len(data) > wb.space() {
		wb.flush(false)
	}
	if len(data) >= bufSize {
		wb.conn.write(wb.id, data, false, false)
	} else {
		wb.buf = append(wb.buf, data...)
	}
	return wb
}

// WriteString is like Write, but for a string.
func (wb *WriteBuf) WriteString(data string) *WriteBuf {
	if len(data) > wb.space() {
		wb.flush(false)
	}
	if len(data) >= bufSize {
		// it would be safer/better to use []byte(s)
		// but strings are used for large data so we want to avoid copying
		wb.conn.write(wb.id, hacks.Stobs(data), false, false)
	} else {
		wb.buf = append(wb.buf, data...)
	}
	return wb
}

// Write1 is like Write, but for a single byte.
func (wb *WriteBuf) Write1(data byte) *WriteBuf {
	if wb.space() == 0 {
		wb.flush(false)
	}
	wb.buf = append(wb.buf, data)
	return wb
}

func (wb *WriteBuf) EndMsg() {
	wb.flush(true)
}

func (wb *WriteBuf) flush(final bool) {
	wb.conn.write(wb.id, wb.buf, true, final)
	wb.buf = wb.buf[:HeaderSize]
}

//-------------------------------------------------------------------

type ReadWrite struct {
	ReadBuf
	WriteBuf
}

const traceLimit = 100

// PutCmd writes a command byte
func (wb *WriteBuf) PutCmd(cmd commands.Command) *WriteBuf {
	trace.ClientServer.Println(">", cmd)
	wb.ResetWrite()
	wb.Write1(byte(cmd))
	return wb
}

// PutBool writes a boolean
func (wb *WriteBuf) PutBool(b bool) *WriteBuf {
	trace.ClientServer.Println("    ->", b)
	if b {
		wb.Write1(1)
	} else {
		wb.Write1(0)
	}
	return wb
}

// PutByte writes a byte
func (wb *WriteBuf) PutByte(b byte) *WriteBuf {
	trace.ClientServer.Println("    ->", b)
	wb.Write1(b)
	return wb
}

// PutStr writes a size prefixed string
func (wb *WriteBuf) PutStr(s string) *WriteBuf {
	if trace.ClientServer.On() {
		if len(s) < traceLimit {
			trace.ClientServer.Println("    ->", s)
		} else {
			trace.ClientServer.Println("    -> string", len(s))
		}
	}
	return wb.PutStr_(s)
}

// PutStr_ writes a size prefixed string without trace
func (wb *WriteBuf) PutStr_(s string) *WriteBuf {
	limit(int64(len(s)))
	wb.putInt(len(s))
	wb.WriteString(s)
	return wb
}

// PutBuf writes a string without a size prefix
func (wb *WriteBuf) PutBuf(s string) *WriteBuf {
	limit(int64(len(s)))
	wb.WriteString(s)
	return wb
}

// PutStrs writes a list of strings
func (wb *WriteBuf) PutStrs(strs []string) *WriteBuf {
	wb.PutInt(len(strs))
	for _, s := range strs {
		wb.PutStr(s)
	}
	return wb
}

// PutRec writes a size prefixed Record
func (wb *WriteBuf) PutRec(r Record) *WriteBuf {
	if trace.ClientServer.On() {
		if len(r) < traceLimit {
			trace.ClientServer.Println("    ->", r)
		} else {
			trace.ClientServer.Println("    -> record", len(r))
		}
	}
	return wb.PutStr_(string(r))
}

// PutInt writes a zig zag encoded varint
func (wb *WriteBuf) PutInt(i int) *WriteBuf {
	trace.ClientServer.Println("    ->", i)
	return wb.PutInt64(int64(i))
}

func (wb *WriteBuf) putInt(i int) *WriteBuf {
	return wb.PutInt64(int64(i))
}

// PutInt64 writes a zig zag encoded varint
func (wb *WriteBuf) PutInt64(i int64) *WriteBuf {
	i = (i << 1) ^ (i >> 63) // zig zag encoding
	n := uint64(i)
	for n > 0x7f {
		wb.Write1(byte(n | 0x80))
		n >>= 7
	}
	wb.Write1(byte(n))
	return wb
}

// PutInts writes a size prefixed list of ints
func (wb *WriteBuf) PutInts(ints []int) *WriteBuf {
	wb.PutInt(len(ints))
	for _, n := range ints {
		wb.PutInt(n)
	}
	return wb
}

// PutResult writes (true, true, PutVal) for non-nil or (true, false) for nil.
func (wb *WriteBuf) PutResult(v Value) *WriteBuf {
	wb.PutBool(true) // no error
	if v == nil {
		return wb.PutBool(false)
	}
	return wb.PutBool(true).PutVal(v)
}

// PutVal writes a size prefixed Pack'ed value
func (wb *WriteBuf) PutVal(val Value) *WriteBuf {
	packed := Pack(val.(Packable))
	if trace.ClientServer.On() {
		if len(packed) < traceLimit {
			trace.ClientServer.Println("    ->", val)
		} else {
			trace.ClientServer.Println("    ->", val.Type())
		}
	}
	return wb.PutStr_(packed)
}

func (wb *WriteBuf) ResetWrite() {
	wb.buf = wb.buf[:HeaderSize] // discard content, keep capacity
}

const maxio = 1024 * 1024 // 1 mb

// limit checks if the size is negative or greater than maxio
func limit(n int64) int {
	assert.That(n >= 0)
	if n > maxio {
		err := fmt.Sprintf("client server io too large (%d > %d)", n, maxio)
		if options.Action == "server" {
			panic(err)
		}
		Fatal(err)
	}
	return int(n)
}

//-------------------------------------------------------------------

type ReadBuf struct {
	buf []byte
}

func (rb *ReadBuf) Remaining() int {
	return len(rb.buf)
}

func (rb *ReadBuf) SetBuf(buf []byte) {
	rb.buf = buf
}

func (rb *ReadBuf) GetByte() byte {
	b := rb.buf[0]
	rb.buf = rb.buf[1:]
	return b
}

func (rb *ReadBuf) GetCmd() commands.Command {
	b := rb.GetByte()
	icmd := commands.Command(b)
	trace.ClientServer.Println("<", icmd)
	return icmd
}

// GetBool reads a boolean
func (rb *ReadBuf) GetBool() bool {
	b := rb.GetByte()
	if b != 0 && b != 1 {
		Fatal("invalid boolean value from server", b, string(b))
	}
	trace.ClientServer.Println("    <-", b == 1)
	return b == 1
}

func (rb *ReadBuf) GetChar() byte {
	b := rb.GetByte()
	trace.ClientServer.Println("    <-", string(b))
	return b
}

// GetInt reads a zig zag encoded varint
func (rb *ReadBuf) GetInt() int {
	n := rb.GetInt64()
	assert.That(int64(math.MinInt) <= n && n <= int64(math.MaxInt))
	trace.ClientServer.Println("    <-", n)
	return int(n)
}

// GetInt64 reads a zig zag encoded varint
func (rb *ReadBuf) GetInt64() int64 {
	shift := uint(0)
	n := uint64(0)
	for {
		b := rb.GetByte()
		n |= uint64(b&0x7f) << shift
		shift += 7
		if 0 == (b & 0x80) {
			break
		}
	}
	tmp := ((int64(n<<63) >> 63) ^ int64(n)) >> 1
	tmp = tmp ^ int64(n&(1<<63))
	return tmp
}

// GetN reads n bytes and returns them in a string
func (rb *ReadBuf) GetN(n int) string {
	s := string(rb.buf[:n]) // ??? hacks.BStoS ???
	rb.buf = rb.buf[n:]
	return s
}

// GetSize returns GetInt, checking the size against the maxio limit
func (rb *ReadBuf) GetSize() int {
	return limit(rb.GetInt64())
}

// GetRec reads a size prefixed string
func (rb *ReadBuf) GetRec() Record {
	n := rb.GetSize()
	rec := Record(rb.GetN(n))
	if trace.ClientServer.On() {
		if len(rec) < traceLimit {
			trace.ClientServer.Println("    <-", rec)
		} else {
			trace.ClientServer.Println("    <- record", len(rec))
		}
	}
	return rec
}

// GetStr reads a size prefixed string
func (rb *ReadBuf) GetStr() string {
	n := rb.GetSize()
	s := rb.GetN(n)
	trace.ClientServer.Println("    <-", s)
	return s
}

// GetStr_ reads a size prefixed string without tracing
func (rb *ReadBuf) GetStr_() string {
	n := rb.GetSize()
	s := rb.GetN(n)
	return s
}

func (rb *ReadBuf) GetStrs() []string {
	n := rb.GetInt()
	list := make([]string, 0, n)
	for ; n > 0; n-- {
		list = append(list, rb.GetStr())
	}
	return list
}

// GetVal reads a packed value
func (rb *ReadBuf) GetVal() Value {
	packed := rb.GetStr_()
	val := Unpack(packed)
	if trace.ClientServer.On() {
		if len(packed) < traceLimit {
			trace.ClientServer.Println("    ->", val)
		} else {
			trace.ClientServer.Println("    ->", val.Type())
		}
	}
	return val
}

// ValueResult reads an optional packed value
func (rb *ReadBuf) ValueResult() Value {
	if rb.GetBool() {
		return rb.GetVal()
	}
	return nil
}
