package csio

import (
	"bufio"
	"io"

	"github.com/apmckinlay/gsuneido/database/clientserver/commands"
)

//TODO err checking

// ReadWrite handles encode/decode for the Suneido client/server protocol.
// It uses bufio for buffering.
type ReadWrite struct {
	r *bufio.Reader
	w *bufio.Writer
}

const maxio = 1024 * 1024 // 1 mb

func NewReadWrite(rw io.ReadWriter) *ReadWrite {
	return &ReadWrite{r: bufio.NewReader(rw), w: bufio.NewWriter(rw)}
}

func (rw *ReadWrite) PutCmd(cmd commands.Command) *ReadWrite {
	rw.w.Write([]byte{byte(cmd)})
	return rw
}

func (rw *ReadWrite) PutStr(s string) *ReadWrite {
	rw.PutInt(int64(len(s)))
	rw.w.Write([]byte(s))
	return rw
}

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

func (rw *ReadWrite) Get(n int) []byte {
	buf := make([]byte, n)
	rw.r.Read(buf)
	return buf
}

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

func (rw *ReadWrite) GetSize() int {
	return limit(rw.GetInt())
}

func (rw *ReadWrite) GetStr() string {
	n := rw.GetInt()
	s := string(rw.Get(limit(n)))
	return s
}

func (rw *ReadWrite) Flush() {
	rw.w.Flush()
}

func limit(n int64) int {
	if n < 0 || maxio < n {
		panic("bad io size")
	}
	return int(n)
}

func (rw *ReadWrite) Request() {
	rw.w.Flush()
	if !rw.GetBool() {
		err := rw.GetStr()
		panic(err + " (from server)")
	}
}
