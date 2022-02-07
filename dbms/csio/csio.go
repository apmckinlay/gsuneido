// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package csio

import (
	"bufio"
	"io"
	"math"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/trace"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hacks"

	"github.com/apmckinlay/gsuneido/dbms/commands"
)

const traceLimit = 100

// ReadWrite handles encode/decode for the Suneido client/server protocol.
// It uses bufio for buffering.
type ReadWrite struct {
	r   *bufio.Reader
	w   *bufio.Writer
	err func(string)
}

const maxio = 1024 * 1024 // 1 mb

// NewReadWrite returns a new ReadWrite
func NewReadWrite(rw io.ReadWriter, err func(string)) *ReadWrite {
	return &ReadWrite{r: bufio.NewReader(rw), w: bufio.NewWriter(rw), err: err}
}

// PutCmd writes a command byte
func (rw *ReadWrite) PutCmd(cmd commands.Command) *ReadWrite {
	trace.ClientServer.Println(">", cmd)
	rw.w.WriteByte(byte(cmd))
	return rw
}

// PutBool writes a boolean
func (rw *ReadWrite) PutBool(b bool) *ReadWrite {
	trace.ClientServer.Println("    ->", b)
	if b {
		rw.w.WriteByte(1)
	} else {
		rw.w.WriteByte(0)
	}
	return rw
}

// PutByte writes a byte
func (rw *ReadWrite) PutByte(b byte) *ReadWrite {
	trace.ClientServer.Println("    ->", b)
	rw.w.WriteByte(b)
	return rw
}

// PutStr writes a size prefixed string
func (rw *ReadWrite) PutStr(s string) *ReadWrite {
	if trace.ClientServer.On() {
		if len(s) < traceLimit {
			trace.ClientServer.Println("    ->", s)
		} else {
			trace.ClientServer.Println("    -> string", len(s))
		}
	}
	return rw.PutStr_(s)
}

// PutStr_ writes a size prefixed string without trace
func (rw *ReadWrite) PutStr_(s string) *ReadWrite {
	limit(int64(len(s)))
	rw.putInt(len(s))
	rw.w.WriteString(s)
	return rw
}

// PutBuf writes a string without a size prefix
func (rw *ReadWrite) PutBuf(s string) *ReadWrite {
	limit(int64(len(s)))
	rw.w.WriteString(s)
	return rw
}

// PutStrs writes a list of strings
func (rw *ReadWrite) PutStrs(strs []string) *ReadWrite {
	rw.PutInt(len(strs))
	for _, s := range strs {
		rw.PutStr(s)
	}
	return rw
}

// PutRec writes a size prefixed Record
func (rw *ReadWrite) PutRec(r Record) *ReadWrite {
	if trace.ClientServer.On() {
		if len(r) < traceLimit {
			trace.ClientServer.Println("    ->", r)
		} else {
			trace.ClientServer.Println("    -> record", len(r))
		}
	}
	return rw.PutStr_(string(r))
}

// PutInt writes a zig zag encoded varint
func (rw *ReadWrite) PutInt(i int) *ReadWrite {
	trace.ClientServer.Println("    ->", i)
	return rw.PutInt64(int64(i))
}

func (rw *ReadWrite) putInt(i int) *ReadWrite {
	return rw.PutInt64(int64(i))
}

// PutInt64 writes a zig zag encoded varint
func (rw *ReadWrite) PutInt64(i int64) *ReadWrite {
	i = (i << 1) ^ (i >> 63) // zig zag encoding
	n := uint64(i)
	for n > 0x7f {
		rw.w.WriteByte(byte(n | 0x80))
		n >>= 7
	}
	rw.w.WriteByte(byte(n))
	return rw
}

// PutInts writes a size prefixed list of ints
func (rw *ReadWrite) PutInts(ints []int) *ReadWrite {
	rw.PutInt(len(ints))
	for _, n := range ints {
		rw.PutInt(n)
	}
	return rw
}

// PutResult writes (true, true, PutVal) for non-nil or (true, false) for nil.
func (rw *ReadWrite) PutResult(v Value) *ReadWrite {
	rw.PutBool(true) // no error
	if v == nil {
		return rw.PutBool(false)
	}
	return rw.PutBool(true).PutVal(v)
}

// PutVal writes a size prefixed Pack'ed value
func (rw *ReadWrite) PutVal(val Value) *ReadWrite {
	packed := Pack(val.(Packable))
	if trace.ClientServer.On() {
		if len(packed) < traceLimit {
			trace.ClientServer.Println("    ->", val)
		} else {
			trace.ClientServer.Println("    ->", val.Type())
		}
	}
	return rw.PutStr_(packed)
}

func (rw *ReadWrite) ResetWrite(w io.Writer) {
	rw.w.Reset(w)
}

//-------------------------------------------------------------------

func (rw *ReadWrite) GetCmd() commands.Command {
	icmd := commands.Command(rw.GetByte_())
	trace.ClientServer.Println("<", icmd)
	return icmd
}

// GetBool reads a boolean
func (rw *ReadWrite) GetBool() bool {
	b := rw.GetByte_()
	if b != 0 && b != 1 {
		Fatal("invalid boolean value from server", b)
	}
	trace.ClientServer.Println("    <-", b == 1)
	return b == 1
}

func (rw *ReadWrite) GetChar() byte {
	b := rw.GetByte_()
	trace.ClientServer.Println("    <-", string(b))
	return b
}

func (rw *ReadWrite) GetByte_() byte {
	b, err := rw.r.ReadByte()
	rw.ck(err)
	return b
}

func (rw *ReadWrite) ck(err error) {
	if err != nil {
		rw.err(err.Error())
	}
}

// GetInt reads a zig zag encoded varint
func (rw *ReadWrite) GetInt() int {
	n := rw.GetInt64()
	assert.That(int64(math.MinInt) <= n && n <= int64(math.MaxInt))
	trace.ClientServer.Println("    <-", n)
	return int(n)
}

// GetInt64 reads a zig zag encoded varint
func (rw *ReadWrite) GetInt64() int64 {
	shift := uint(0)
	n := uint64(0)
	for {
		b := rw.GetByte_()
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
func (rw *ReadWrite) GetN(n int) string {
	buf := make([]byte, n)
	_, err := io.ReadFull(rw.r, buf)
	rw.ck(err)
	return hacks.BStoS(buf) // safe since buf doesn't escape
}

// GetSize returns GetInt, checking the size against the maxio limit
func (rw *ReadWrite) GetSize() int {
	return limit(rw.GetInt64())
}

// GetRec reads a size prefixed string
func (rw *ReadWrite) GetRec() Record {
	n := rw.GetSize()
	r := Record(rw.GetN(n))
	if trace.ClientServer.On() {
		if len(r) < traceLimit {
			trace.ClientServer.Println("    <-", r)
		} else {
			trace.ClientServer.Println("    <- record", len(r))
		}
	}
	return r
}

// GetStr reads a size prefixed string
func (rw *ReadWrite) GetStr() string {
	n := rw.GetSize()
	s := rw.GetN(n)
	trace.ClientServer.Println("    <-", s)
	return s
}

// GetStr_ reads a size prefixed string without tracing
func (rw *ReadWrite) GetStr_() string {
	n := rw.GetSize()
	s := rw.GetN(n)
	return s
}

func (rw *ReadWrite) GetStrs() []string {
	n := rw.GetInt()
	list := make([]string, 0, n)
	for ; n > 0; n-- {
		list = append(list, rw.GetStr())
	}
	return list
}

// GetVal reads a packed value
func (rw *ReadWrite) GetVal() Value {
	packed := rw.GetStr_()
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
func (rw *ReadWrite) ValueResult() Value {
	if rw.GetBool() {
		return rw.GetVal()
	}
	return nil
}

// Flush flushes the Writer
func (rw *ReadWrite) Flush() {
	rw.w.Flush()
}

// limit checks if the size is negative or greater than maxio
func limit(n int64) int {
	if n < 0 || maxio < n {
		Fatal("bad io size:", n)
	}
	return int(n)
}

// Request is used by DbmsClient.
// It does Flush and GetBool for the result.
// If the result is false, it does GetStr for the error and panics with it.
func (rw *ReadWrite) Request() {
	rw.ck(rw.w.Flush())
	if !rw.GetBool() {
		err := rw.GetStr()
		trace.ClientServer.Println(err)
		panic(err + " (from server)")
	}
}
