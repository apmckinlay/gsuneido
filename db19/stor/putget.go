// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

// Put methods return the writer so they can be chained.

type Writer struct {
	buf []byte
}

// NewWriter returns a new Writer on a byte slice
func NewWriter(buf []byte) *Writer {
	return &Writer{buf[:0]}
}

// Put1 writes an unsigned byte value
func (w *Writer) Put1(n int) *Writer {
	if n < 0 || 1<<8 <= n {
		panic("stor.Writer.Put1 value outside range")
	}
	w.buf = append(w.buf,
		byte(n))
	return w
}

// Put2 writes an unsigned two byte value
func (w *Writer) Put2(n int) *Writer {
	if n < 0 || 1<<16 <= n {
		panic("stor.Writer.Put2 value outside range")
	}
	w.buf = append(w.buf,
		byte(n),
		byte(n>>8))
	return w
}

// Put3 writes an unsigned three byte value
func (w *Writer) Put3(n int) *Writer {
	if n < 0 || 1<<24 <= n {
		panic("stor.Writer.Put3 value outside range")
	}
	w.buf = append(w.buf,
		byte(n),
		byte(n>>8),
		byte(n>>16))
	return w
}

// Put4 writes an unsigned four byte value
func (w *Writer) Put4(n int) *Writer {
	if n < 0 || 1<<32 <= n {
		panic("stor.Writer.Put4 value outside range")
	}
	w.buf = append(w.buf,
		byte(n),
		byte(n>>8),
		byte(n>>16),
		byte(n>>24))
	return w
}

// Put5 writes an unsigned five byte value
func (w *Writer) Put5(n uint64) *Writer {
	if n < 0 || 1<<40 <= n {
		panic("stor.Writer.Put5 value outside range")
	}
	w.buf = append(w.buf,
		byte(n),
		byte(n>>8),
		byte(n>>16),
		byte(n>>24),
		byte(n>>32))
	return w
}

// PutStr writes a string with a maximum length of 64k
func (w *Writer) PutStr(s string) *Writer {
	w.Put2(len(s))
	w.buf = append(w.buf, s...)
	return w
}

// LenStr returns the space requirecd by PutStr
func LenStr(s string) int {
	return 2 + len(s)
}

// PutStrs writes a slice of up to 64k strings
func (w *Writer) PutStrs(ss []string) *Writer {
	w.Put2(len(ss))
	for _, col := range ss {
		w.PutStr(col)
	}
	return w
}

// LenStrs returns the space requirecd by PutStrs
func LenStrs(ss []string) int {
	n := 2
	for _, s := range ss {
		n += LenStr(s)
	}
	return n
}

// Len returns the current position within this writer
func (w *Writer) Len() int {
	return len(w.buf)
}

//-------------------------------------------------------------------

type Reader struct {
	buf []byte
}

func (stor *Stor) Reader(off uint64) *Reader {
	return NewReader(stor.Data(off))
}

// Reader returns a Reader based on the offset returned by Writer.Close
func NewReader(buf []byte) *Reader {
	return &Reader{buf}
}

// Get1 reads an unsigned byte value
func (r *Reader) Get1() int {
	n := int(r.buf[0])
	r.buf = r.buf[1:]
	return n
}

// Get2 reads an unsigned two byte value
func (r *Reader) Get2() int {
	n := int(r.buf[0]) + int(r.buf[1])<<8
	r.buf = r.buf[2:]
	return n
}

// Get3 reads an unsigned three byte value
func (r *Reader) Get3() int {
	n := int(r.buf[0]) + int(r.buf[1])<<8 + int(r.buf[2])<<16
	r.buf = r.buf[3:]
	return n
}

// Get4 reads an unsigned four byte value
func (r *Reader) Get4() int {
	n := int(r.buf[0]) + int(r.buf[1])<<8 + int(r.buf[2])<<16 + int(r.buf[3])<<24
	r.buf = r.buf[4:]
	return n
}

// Get5 reads an unsigned five byte value
func (r *Reader) Get5() uint64 {
	n := uint64(r.buf[0]) + uint64(r.buf[1])<<8 + uint64(r.buf[2])<<16 +
		uint64(r.buf[3])<<24 + uint64(r.buf[4])<<32
	r.buf = r.buf[5:]
	return n
}

// GetStr reads a string
func (r *Reader) GetStr() string {
	n := r.Get2()
	s := string(r.buf[:n])
	r.buf = r.buf[n:]
	return s
}

// GetStrs reads a slice of strings
func (r *Reader) GetStrs() []string {
	n := r.Get2()
	if n == 0 {
		return nil
	}
	ss := make([]string, n)
	for i := 0; i < n; i++ {
		ss[i] = r.GetStr()
	}
	return ss
}

// Remaining returns the number of unread bytes left
func (r *Reader) Remaining() int {
	return len(r.buf)
}
