// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import "github.com/apmckinlay/gsuneido/util/verify"

// Writer allows writing data to a stor with an unknown length
// by breaking it up into blocks.
// (If the length of the data is known, just Alloc)
// Block size is a tradeoff between per-block overhead and last block wastage.
// Put methods return the writer so they can be chained.
// Note: Individual Put calls will NOT straddle blocks.
// This is so the Put's and Get's only have to check the block space once.
// WARNING: This means your Put's and Get's should correspond.
// e.g. if you write with Put3, you should read with Get3.
type Writer struct {
	stor      *Stor
	blockSize int
	blocks    []uint64
	buf       []byte
}

// Writer returns a new Writer for the Stor
func (stor *Stor) Writer(blockSize int) *Writer { //TODO make this a method on Stor
	return &Writer{stor: stor, blockSize: blockSize}
}

// Put1 writes an unsigned byte value
func (w *Writer) Put1(n int) *Writer {
	if n < 0 || 1<<8 <= n {
		panic("stor.Writer.Put1 value outside range")
	}
	if len(w.buf) < 1 {
		w.allocNext()
	}
	w.buf[0] = byte(n)
	w.buf = w.buf[1:]
	return w
}

// Put2 writes an unsigned two byte value
func (w *Writer) Put2(n int) *Writer {
	if n < 0 || 1<<16 <= n {
		panic("stor.Writer.Put2 value outside range")
	}
	if len(w.buf) < 2 {
		w.allocNext()
	}
	w.buf[0] = byte(n)
	w.buf[1] = byte(n >> 8)
	w.buf = w.buf[2:]
	return w
}

// Put3 writes an unsigned three byte value
func (w *Writer) Put3(n int) *Writer {
	if n < 0 || 1<<24 <= n {
		panic("stor.Writer.Put3 value outside range")
	}
	if len(w.buf) < 3 {
		w.allocNext()
	}
	w.buf[0] = byte(n)
	w.buf[1] = byte(n >> 8)
	w.buf[2] = byte(n >> 16)
	w.buf = w.buf[3:]
	return w
}

// Put4 writes an unsigned four byte value
func (w *Writer) Put4(n int) *Writer {
	if n < 0 || 1<<32 <= n {
		panic("stor.Writer.Put4 value outside range")
	}
	if len(w.buf) < 4 {
		w.allocNext()
	}
	w.buf[0] = byte(n)
	w.buf[1] = byte(n >> 8)
	w.buf[2] = byte(n >> 16)
	w.buf[3] = byte(n >> 24)
	w.buf = w.buf[4:]
	return w
}

// Put5 writes an unsigned five byte value
func (w *Writer) Put5(n uint64) *Writer {
	if n < 0 || 1<<40 <= n {
		panic("stor.Writer.Put5 value outside range")
	}
	if len(w.buf) < 5 {
		w.allocNext()
	}
	w.buf[0] = byte(n)
	w.buf[1] = byte(n >> 8)
	w.buf[2] = byte(n >> 16)
	w.buf[3] = byte(n >> 24)
	w.buf[4] = byte(n >> 32)
	w.buf = w.buf[5:]
	return w
}

// PutStr writes a string with a maximum length of 64k
func (w *Writer) PutStr(s string) *Writer {
	n := len(s)
	w.Put2(n)
	if len(w.buf) < n {
		w.allocNext()
	}
	copy(w.buf, s)
	w.buf = w.buf[len(s):]
	return w
}

// PutInts writes a slice of <256 int's, each <64k
func (w *Writer) PutInts(ints []int) *Writer {
	w.Put1(len(ints))
	for _, n := range ints {
		w.Put2(n)
	}
	return w
}

func (w *Writer) allocNext() {
	off, buf := w.stor.Alloc(w.blockSize)
	w.buf = buf
	w.blocks = append(w.blocks, off)
}

// Offset returns the current position within this writer
func (w *Writer) Pos() int {
	return w.blockSize*len(w.blocks) - len(w.buf)
}

// Close writes out the list of blocks and returns the offset of the list.
func (w *Writer) Close() uint64 {
	needed := 5 + 5*len(w.blocks)
	if len(w.buf) < needed {
		w.allocNext()
	}
	verify.That(len(w.buf) >= needed)
	off := w.blocks[len(w.blocks)-1] + uint64(w.blockSize-len(w.buf))
	w.Put3(w.blockSize)
	w.Put2(len(w.blocks))
	for _, o := range w.blocks {
		w.Put5(o)
	}
	return off
}

//-------------------------------------------------------------------

type Reader struct {
	stor      *Stor
	blockSize int
	blocks    []uint64
	buf       []byte
	bi        int
}

// Reader returns a Reader based on the offset returned by Writer.Close
func (stor *Stor) Reader(off uint64) *Reader {
	r := &Reader{stor: stor}
	r.buf = stor.Data(off)
	r.blockSize = r.Get3()
	nblocks := r.Get2()
	r.blocks = make([]uint64, nblocks)
	for i := 0; i < nblocks; i++ {
		r.blocks[i] = r.Get5()
	}
	r.bi = -1
	r.nextBlock()
	return r
}

func (r *Reader) nextBlock() {
	r.bi++
	data := r.stor.Data(r.blocks[r.bi])
	r.buf = data[:r.blockSize]
}

// Pos sets the position within the written data
// normally with a value from Writer.Pos
func (r *Reader) Pos(pos int) *Reader {
	r.bi = pos / r.blockSize
	ib := pos % r.blockSize
	data := r.stor.Data(r.blocks[r.bi])
	r.buf = data[ib:r.blockSize]
	return r
}

// Get1 reads an unsigned byte value
func (r *Reader) Get1() int {
	if len(r.buf) < 1 {
		r.nextBlock()
	}
	n := int(r.buf[0])
	r.buf = r.buf[1:]
	return n
}

// Get2 reads an unsigned two byte value
func (r *Reader) Get2() int {
	if len(r.buf) < 2 {
		r.nextBlock()
	}
	n := int(r.buf[0]) + int(r.buf[1])<<8
	r.buf = r.buf[2:]
	return n
}

// Get3 reads an unsigned three byte value
func (r *Reader) Get3() int {
	if len(r.buf) < 3 {
		r.nextBlock()
	}
	n := int(r.buf[0]) + int(r.buf[1])<<8 + int(r.buf[2])<<16
	r.buf = r.buf[3:]
	return n
}

// Get4 reads an unsigned four byte value
func (r *Reader) Get4() int {
	if len(r.buf) < 4 {
		r.nextBlock()
	}
	n := int(r.buf[0]) + int(r.buf[1])<<8 + int(r.buf[2])<<16 + int(r.buf[3])<<24
	r.buf = r.buf[4:]
	return n
}

// Get5 reads an unsigned five byte value
func (r *Reader) Get5() uint64 {
	if len(r.buf) < 5 {
		r.nextBlock()
	}
	n := uint64(r.buf[0]) + uint64(r.buf[1])<<8 + uint64(r.buf[2])<<16 +
		uint64(r.buf[3])<<24 + uint64(r.buf[4])<<32
	r.buf = r.buf[5:]
	return n
}

// GetStr reads a string
func (r *Reader) GetStr() string {
	n := r.Get2()
	if len(r.buf) < n {
		r.nextBlock()
	}
	s := string(r.buf[:n])
	r.buf = r.buf[n:]
	return s
}

// GetStr reads a slice of int's
func (r *Reader) GetInts() []int {
	n := r.Get1()
	ints := make([]int, n)
	for i := 0; i < n; i++ {
		ints[i] = r.Get2()
	}
	return ints
}

func (r *Reader) Stor() *Stor {
	return r.stor
}
