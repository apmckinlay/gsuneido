package csio

import (
	"bufio"
	"io"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"

	"github.com/apmckinlay/gsuneido/database/dbms/commands"
)

//TODO err checking

// ReadWrite handles encode/decode for the Suneido client/server protocol.
// It uses bufio for buffering.
type ReadWrite struct {
	r *bufio.Reader
	w *bufio.Writer
}

const maxio = 1024 * 1024 // 1 mb

// NewReadWrite returns a new ReadWrite
func NewReadWrite(rw io.ReadWriter) *ReadWrite {
	return &ReadWrite{r: bufio.NewReader(rw), w: bufio.NewWriter(rw)}
}

// PutCmd writes a command byte
func (rw *ReadWrite) PutCmd(cmd commands.Command) *ReadWrite {
	rw.w.WriteByte(byte(cmd))
	return rw
}

// PutStr writes a size prefixed string
func (rw *ReadWrite) PutStr(s string) *ReadWrite {
	limit(int64(len(s)))
	rw.PutInt(int64(len(s)))
	rw.w.WriteString(s)
	return rw
}

// PutInt writes a zig zag encoded varint
func (rw *ReadWrite) PutInt(i int64) *ReadWrite {
	i = (i << 1) ^ (i >> 63) // zig zag encoding
	n := uint64(i)
	for n > 0x7f {
		rw.w.WriteByte(byte(n | 0x80))
		n >>= 7
	}
	rw.w.WriteByte(byte(n))
	return rw
}

// GetBool reads a boolean
func (rw *ReadWrite) GetBool() bool {
	b, _ := rw.r.ReadByte()
	switch b {
	case 0:
		return false
	case 1:
		return true
	default:
		panic("bad boolean")
	}
}

// Get reads n bytes and returns it in a newly allocated buffer
func (rw *ReadWrite) Get(n int) []byte {
	buf := make([]byte, n)
	io.ReadFull(rw.r, buf)
	return buf
}

// GetInt reads a zig zag encoded varint
func (rw *ReadWrite) GetInt() int64 {
	shift := uint(0)
	n := uint64(0)
	for {
		b, _ := rw.r.ReadByte()
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

// GetSize returns GetInt, checking the size against the maxio limit
func (rw *ReadWrite) GetSize() int {
	return limit(rw.GetInt())
}

// GetStr reads a size prefixed string
func (rw *ReadWrite) GetStr() string {
	n := rw.GetSize()
	buf := rw.Get(n)
	return *(*string)(unsafe.Pointer(&buf)) // safe since buf doesn't escape
}

func (rw *ReadWrite) GetVal() Value {
	return Unpack(rw.GetStr())
}

// Flush flushes the Writer
func (rw *ReadWrite) Flush() {
	rw.w.Flush()
}

// limit panics if the size is negative or greater than maxio
func limit(n int64) int {
	if 0 <= n && n < maxio {
		return int(n)
	}
	panic("bad io size")
}

// Request does Flush and GetBool for the result.
// If the result is false, it does GetStr for the error and panics with it.
func (rw *ReadWrite) Request() {
	rw.w.Flush()
	if !rw.GetBool() {
		err := rw.GetStr()
		panic(err + " (from server)")
	}
}
