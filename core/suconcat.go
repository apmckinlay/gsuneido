// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"bytes"
	"cmp"
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/pack"
)

// SuConcat is a Value used to optimize string concatenation
// WARNING: zero value is not valid, use NewSuConcat
type SuConcat struct {
	ValueBase[SuConcat]
	buf *scbuf
	n   int
}

// scbuf is the potentially shared byte slice buffer
// i.e. multiple SuConcat's may point to the same scbuf
type scbuf struct {
	bs []byte
	// concurrent is set by SetConcurrent
	// once it is set, shared must be immutable
	concurrent bool
}

// NewSuConcat returns an empty SuConcat
func NewSuConcat() SuConcat {
	return SuConcat{buf: &scbuf{}}
}

func (c SuConcat) SetConcurrent() {
	c.buf.concurrent = true
}

func (c SuConcat) IsConcurrent() Value {
	return SuBool(c.buf.concurrent)
}

// Len returns the length of an SuConcat
func (c SuConcat) Len() int {
	return c.n
}

const StringLimit = 32_000_000 // ???

func CheckStringSize(op string, n int) {
	if n > StringLimit {
		panic(fmt.Sprint("ERROR ", op + ": string > ", StringLimit))
	}
}

// Add appends a string to an SuConcat
func (c SuConcat) Add(s string) SuConcat {
	CheckStringSize("concatenate", c.n + len(s))
	buf := c.buf
	if buf.concurrent || // shared between threads
		len(buf.bs) != c.n { // another SuConcat has appended their own stuff
		// copy to our own new buffer
		a := append(make([]byte, 0, c.n+len(s)), buf.bs[:c.n]...)
		buf = &scbuf{bs: a}
	}
	buf.bs = append(buf.bs, s...)
	return SuConcat{buf: buf, n: c.n + len(s)}
}

func (c SuConcat) Iter() Iter {
	return &stringIter{s: c.toStr()}
}

// Value interface --------------------------------------------------

var _ Value = (*SuConcat)(nil)

func (c SuConcat) ToInt() (int, bool) {
	return 0, c.n == 0
}

func (c SuConcat) ToDnum() (dnum.Dnum, bool) {
	return dnum.Zero, c.n == 0
}

func (c SuConcat) AsStr() (string, bool) {
	return c.toStr(), true
}

func (c SuConcat) ToStr() (string, bool) {
	return c.toStr(), true
}

func (c SuConcat) toStr() string {
	// use the same trick as strings.Builder to avoid allocation
	s := hacks.BStoS(c.buf.bs)
	return s[:c.n]
}

// String returns a quoted string
func (c SuConcat) String() string {
	return escapeStr(c.toStr(), 0)
}

// Get returns the character at a given index
func (c SuConcat) Get(_ *Thread, key Value) Value {
	return strGet(c.toStr(), key)
}

func (c SuConcat) RangeTo(from int, to int) Value {
	from = prepFrom(from, c.n)
	to = prepTo(from, to, c.n)
	return SuStr(c.toStr()[from:to])
}

func (c SuConcat) RangeLen(from int, n int) Value {
	from = prepFrom(from, c.n)
	n = prepLen(n, c.n-from)
	return SuStr(c.toStr()[from : from+n])
}

func (c SuConcat) Hash() uint64 {
	return hash.Bytes(c.buf.bs[:c.n])
}

func (c SuConcat) Hash2() uint64 {
	return c.Hash()
}

func (c SuConcat) Equal(other any) bool {
	// check string first assuming more common than concat
	if s2, ok := other.(SuStr); ok {
		// according to benchmark, this doesn't allocate
		return c.n == len(s2) && c.toStr() == string(s2)
	}
	if c2, ok := other.(SuConcat); ok {
		return c.n == c2.n && bytes.Equal(c.buf.bs[:c.n], c2.buf.bs[:c.n])
	}
	return false
}

func (SuConcat) Type() types.Type {
	return types.String
}

func (c SuConcat) Compare(other Value) int {
	if cmp := cmp.Compare(ordStr, Order(other)); cmp != 0 {
		return cmp * 2
	}
	// now know other is a string so AsStr won't panic
	return strings.Compare(c.toStr(), AsStr(other))
}

func (c SuConcat) Call(th *Thread, this Value, as *ArgSpec) Value {
	ss := SuStr(c.toStr())
	return ss.Call(th, this, as)
}

func (SuConcat) Lookup(th *Thread, method string) Callable {
	return Lookup(th, StringMethods, gnStrings, method)
}

// Packable interface -----------------------------------------------

var _ Packable = SuConcat{}

func (c SuConcat) PackSize(*uint64) int {
	if c.n == 0 {
		return 0
	}
	return 1 + c.n
}

func (c SuConcat) PackSize2(*uint64, packStack) int {
	return c.PackSize(nil)
}

func (c SuConcat) Pack(_ *uint64, buf *pack.Encoder) {
	if c.n > 0 {
		buf.Put1(PackString).Put(c.buf.bs[:c.n])
	}
}
