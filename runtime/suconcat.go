package runtime

import (
	"bytes"
	"strings"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/ints"
)

// SuConcat is a Value used to optimize string concatenation
// NOTE: Not thread safe
type SuConcat struct {
	b *shared
	n int
}

type shared struct {
	a []byte
	// MAYBE have a string to cache?
}

var _ Value = SuConcat{}
var _ Packable = SuConcat{}

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
	return SuConcat{bb, c.n + len(s)}
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
	return SuConcat{bb, c.n + cv2.Len()}
}

// Value interface --------------------------------------------------

// ToInt converts an SuConcat to an integer (Value interface)
func (c SuConcat) ToInt() int {
	if c.n == 0 {
		return 0
	}
	panic("can't convert String to integer")
}

// ToDnum converts an SuConcat to a Dnum (Value interface)
func (c SuConcat) ToDnum() dnum.Dnum {
	if c.n == 0 {
		return dnum.Zero
	}
	panic("can't convert String to number")
}

// ToStr converts an SuConcat to a string (Value interface)
func (c SuConcat) ToStr() string {
	return string(c.b.a[:c.n])
}

// String returns a quoted string (Value interface)
// TODO: handle escaping
func (c SuConcat) String() string {
	return "'" + c.ToStr() + "'"
}

// Get returns the character at a given index (Value interface)
func (c SuConcat) Get(key Value) Value {
	return SuStr(string(c.b.a[:c.n][key.ToInt()]))
}

// Put is not applicable to SuConcat (Value interface)
func (SuConcat) Put(Value, Value) {
	panic("strings do not support put")
}

func (c SuConcat) RangeTo(i int, j int) Value {
	// TODO prep indexes
	return SuStr(c.ToStr()[i:j])
}

func (c SuConcat) RangeLen(i int, n int) Value {
	// TODO prep indexes
	return SuStr(c.ToStr()[i : i+n])
}

// Hash returns a hash value for an SuConcat (Value interface)
func (c SuConcat) Hash() uint32 {
	return hash.HashBytes(c.b.a[:c.n])
}

// Hash2 is used to hash nested values (Value interface)
func (c SuConcat) Hash2() uint32 {
	return c.Hash()
}

// Equals returns true if other is an equal SuConcat or SuStr (Value interface)
func (c SuConcat) Equal(other interface{}) bool {
	// check string first assuming more common than concat
	if s2, ok := other.(SuStr); ok {
		// according to benchmark, this doesn't allocate
		return c.n == len(s2) && c.ToStr() == string(s2)
	}
	if c2, ok := other.(SuConcat); ok {
		return c.n == c2.n && bytes.Equal(c.b.a[:c.n], c2.b.a[:c.n])
	}
	return false
}

// TypeName returns the name of this type (Value interface)
func (SuConcat) TypeName() string {
	return "String"
}

// Order returns the ordering of SuDnum (Value interface)
func (SuConcat) Order() Ord {
	return ordStr
}

// Compare compares an SuDnum to another Value (Value interface)
func (c SuConcat) Compare(other Value) int {
	if cmp := ints.Compare(c.Order(), other.Order()); cmp != 0 {
		return cmp
	}
	return strings.Compare(c.ToStr(), other.ToStr())
}

// Packable interface -----------------------------------------------

func (c SuConcat) PackSize() int {
	if c.n == 0 {
		return 0
	}
	return 1 + c.n
}

func (c SuConcat) Pack(buf []byte) []byte {
	buf = append(buf, packString)
	buf = append(buf, c.b.a[:c.n]...)
	return buf
}
