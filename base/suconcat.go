package base

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hash"
)

// SuConcat is used to optimize string concatenation
//
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

func NewSuConcat() SuConcat {
	return SuConcat{b: &shared{}}
}

func (c SuConcat) Len() int {
	return c.n
}

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

func (c SuConcat) ToInt() int32 {
	i, _ := strconv.ParseInt(c.ToStr(), 0, 32)
	return int32(i)
}

func (c SuConcat) ToDnum() dnum.Dnum {
	dn, err := dnum.Parse(c.ToStr())
	if err != nil {
		panic("can't convert this string to a number")
	}
	return dn
}

func (c SuConcat) ToStr() string {
	return string(c.b.a[:c.n])
}

func (c SuConcat) String() string {
	return "'" + c.ToStr() + "'"
}

func (c SuConcat) Get(key Value) Value {
	return SuStr(string(c.b.a[:c.n][key.ToInt()]))
}

func (c SuConcat) Put(key Value, val Value) {
	panic("strings do not support put")
}

func (c SuConcat) Hash() uint32 {
	return hash.HashBytes(c.b.a[:c.n])
}

func (c SuConcat) hash2() uint32 {
	return c.Hash()
}

func (c SuConcat) Equals(other interface{}) bool {
	if c2, ok := other.(SuConcat); ok {
		return c == c2
	}
	if s2, ok := other.(SuStr); ok && c.n == len(s2) {
		for i := 0; i < c.n; i++ {
			if c.b.a[i] != string(s2)[i] {
				return false
			}
			return true
		}
	}
	return false
}

func (c SuConcat) PackSize() int {
	if c.n == 0 {
		return 0
	} else {
		return 1 + c.n
	}
}

func (c SuConcat) Pack(buf []byte) []byte {
	buf = append(buf, packString)
	buf = append(buf, c.b.a[:c.n]...)
	return buf
}

func (_ SuConcat) TypeName() string {
	return "String"
}

func (_ SuConcat) Order() ord {
	return ordStr
}

func (c SuConcat) Cmp(other Value) int {
	// COULD optimize this to not convert Concat to string
	s1 := c.ToStr()
	s2 := other.ToStr()
	switch {
	case s1 < s2:
		return -1
	case s1 > s2:
		return +1
	default:
		return 0
	}

}

func (_ SuConcat) Call(c CallContext) Value {
	panic("concat call not implemented") //TODO
}
