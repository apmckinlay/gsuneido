package runtime

import (
	"bytes"
	"strings"

	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/pack"
)

// SuConcat is a Value used to optimize string concatenation
// NOTE: Not thread safe
type SuConcat struct {
	b *shared
	n int
	CantConvert
}

type shared struct {
	a []byte
	// MAYBE have a string to cache?
}

// NewSuConcat returns an empty SuConcat
func NewSuConcat() SuConcat {
	return SuConcat{b: &shared{}}
}

// Len returns the length of an SuConcat
func (c SuConcat) Len() int {
	return c.n
}

// Add appends a string to an SuConcat
func (c SuConcat) Add(s string) SuConcat {
	bb := c.b
	if len(bb.a) != c.n {
		// another reference has appended their own stuff so make our own buf
		a := append(make([]byte, 0, c.n+len(s)), bb.a[:c.n]...)
		bb = &shared{a}
	}
	bb.a = append(bb.a, s...)
	return SuConcat{b: bb, n: c.n + len(s)}
}

// AddSuConcat appends an SuConcat to an SuConcat
func (c SuConcat) AddSuConcat(cv2 SuConcat) SuConcat {
	// avoid converting cv2 to string
	bb := c.b
	if len(bb.a) != c.n {
		// another reference has appended their own stuff so make our own buf
		a := append(make([]byte, 0, c.n+cv2.Len()), bb.a[:c.n]...)
		bb = &shared{a}
	}
	bb.a = append(bb.a, cv2.b.a...)
	return SuConcat{b: bb, n: c.n + cv2.Len()}
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
	return string(c.b.a[:c.n])
}

// String returns a quoted string
func (c SuConcat) String() string {
	// TODO: handle escaping
	return "'" + c.toStr() + "'"
}

// Get returns the character at a given index
func (c SuConcat) Get(_ *Thread, key Value) Value {
	return strGet(c.toStr(), key)
}

func (SuConcat) Put(*Thread, Value, Value) {
	panic("strings do not support put")
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

func (c SuConcat) Hash() uint32 {
	return hash.HashBytes(c.b.a[:c.n])
}

func (c SuConcat) Hash2() uint32 {
	return c.Hash()
}

func (c SuConcat) Equal(other interface{}) bool {
	// check string first assuming more common than concat
	if s2, ok := other.(SuStr); ok {
		// according to benchmark, this doesn't allocate
		return c.n == len(s2) && c.toStr() == string(s2)
	}
	if c2, ok := other.(SuConcat); ok {
		return c.n == c2.n && bytes.Equal(c.b.a[:c.n], c2.b.a[:c.n])
	}
	return false
}

func (SuConcat) Type() types.Type {
	return types.String
}

func (c SuConcat) Compare(other Value) int {
	if cmp := ints.Compare(ordStr, Order(other)); cmp != 0 {
		return cmp
	}
	// now know other is a string so AsStr won't panic
	return strings.Compare(c.toStr(), AsStr(other))
}

func (c SuConcat) Call(t *Thread, this Value, as *ArgSpec) Value {
	ss := SuStr(c.toStr())
	return ss.Call(t, this, as)
}

func (SuConcat) Lookup(_ *Thread, method string) Callable {
	return StringMethods[method]
}

// Packable interface -----------------------------------------------

var _ Packable = SuConcat{}

func (c SuConcat) PackSize(int) int {
	if c.n == 0 {
		return 0
	}
	return 1 + c.n
}

func (c SuConcat) Pack(buf *pack.Encoder) {
	if c.n > 0 {
		buf.Put1(PackString).Put(c.b.a[:c.n])
	}
}
