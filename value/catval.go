package value

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hash"
)

// CatVal is used to optimize string concatenation
type CatVal struct {
	b *shared
	n int
}

type shared struct {
	a []byte
	// MAYBE have a string
}

func NewCatVal() CatVal {
	return CatVal{b: &shared{}}
}

func (cv CatVal) Len() int {
	return cv.n
}

func (cv CatVal) Add(s string) CatVal {
	bb := cv.b
	if len(bb.a) != cv.n {
		// another reference has appended their own stuff so make our own buf
		a := append(make([]byte, 0, cv.n+len(s)), bb.a[:cv.n]...)
		bb = &shared{a}
	}
	bb.a = append(bb.a, s...)
	return CatVal{bb, cv.n + len(s)}
}

func (cv CatVal) AddCatVal(cv2 CatVal) CatVal {
	// avoid converting cv2 to string
	bb := cv.b
	if len(bb.a) != cv.n {
		// another reference has appended their own stuff so make our own buf
		a := append(make([]byte, 0, cv.n+cv2.Len()), bb.a[:cv.n]...)
		bb = &shared{a}
	}
	bb.a = append(bb.a, cv2.b.a...)
	return CatVal{bb, cv.n + cv2.Len()}
}

func (cv CatVal) ToInt() int32 {
	i, _ := strconv.ParseInt(cv.ToStr(), 0, 32)
	return int32(i)
}

func (cv CatVal) ToDnum() dnum.Dnum {
	dn, err := dnum.Parse(cv.ToStr())
	if err != nil {
		panic("can't convert this string to a number")
	}
	return dn
}

func (cv CatVal) ToStr() string {
	return string(cv.b.a[:cv.n])
}

func (cv CatVal) String() string {
	return "'" + cv.ToStr() + "'"
}

func (cv CatVal) Get(key Value) Value {
	return StrVal(string(cv.b.a[:cv.n][key.ToInt()]))
}

func (cv CatVal) Put(key Value, val Value) {
	panic("strings do not support put")
}

func (cv CatVal) Hash() uint32 {
	return hash.HashBytes(cv.b.a[:cv.n])
}

func (cv CatVal) hash2() uint32 {
	return cv.Hash()
}

func (cv CatVal) Equals(other interface{}) bool {
	if c2, ok := other.(CatVal); ok {
		return cv == c2
	}
	if s2, ok := other.(StrVal); ok && cv.n == len(s2) {
		for i := 0; i < cv.n; i++ {
			if cv.b.a[i] != string(s2)[i] {
				return false
			}
			return true
		}
	}
	return false
}

func (cv CatVal) PackSize() int {
	if cv.n == 0 {
		return 0
	} else {
		return 1 + cv.n
	}
}

func (cv CatVal) Pack(buf []byte) []byte {
	n := cv.n
	if n == 0 {
		return buf
	}
	i := len(buf)
	buf = buf[:i+1+n]
	buf[i] = STRING
	copy(buf[i+1:], cv.b.a[:cv.n])
	return buf
}

var _ Value = CatVal{} // confirm it implements Value
