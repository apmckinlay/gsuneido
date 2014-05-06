package value

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hash"
)

// SuConcat is used to optimize string concatenation
type SuConcat struct {
	b *shared
	n int
}

type shared struct {
	a []byte
	// MAYBE have a string
}

func NewSuConcat() SuConcat {
	return SuConcat{b: &shared{}}
}

func (cv SuConcat) Len() int {
	return cv.n
}

func (cv SuConcat) Add(s string) SuConcat {
	bb := cv.b
	if len(bb.a) != cv.n {
		// another reference has appended their own stuff so make our own buf
		a := append(make([]byte, 0, cv.n+len(s)), bb.a[:cv.n]...)
		bb = &shared{a}
	}
	bb.a = append(bb.a, s...)
	return SuConcat{bb, cv.n + len(s)}
}

func (cv SuConcat) AddSuConcat(cv2 SuConcat) SuConcat {
	// avoid converting cv2 to string
	bb := cv.b
	if len(bb.a) != cv.n {
		// another reference has appended their own stuff so make our own buf
		a := append(make([]byte, 0, cv.n+cv2.Len()), bb.a[:cv.n]...)
		bb = &shared{a}
	}
	bb.a = append(bb.a, cv2.b.a...)
	return SuConcat{bb, cv.n + cv2.Len()}
}

func (cv SuConcat) ToInt() int32 {
	i, _ := strconv.ParseInt(cv.ToStr(), 0, 32)
	return int32(i)
}

func (cv SuConcat) ToDnum() dnum.Dnum {
	dn, err := dnum.Parse(cv.ToStr())
	if err != nil {
		panic("can't convert this string to a number")
	}
	return dn
}

func (cv SuConcat) ToStr() string {
	return string(cv.b.a[:cv.n])
}

func (cv SuConcat) String() string {
	return "'" + cv.ToStr() + "'"
}

func (cv SuConcat) Get(key Value) Value {
	return SuStr(string(cv.b.a[:cv.n][key.ToInt()]))
}

func (cv SuConcat) Put(key Value, val Value) {
	panic("strings do not support put")
}

func (cv SuConcat) Hash() uint32 {
	return hash.HashBytes(cv.b.a[:cv.n])
}

func (cv SuConcat) hash2() uint32 {
	return cv.Hash()
}

func (cv SuConcat) Equals(other interface{}) bool {
	if c2, ok := other.(SuConcat); ok {
		return cv == c2
	}
	if s2, ok := other.(SuStr); ok && cv.n == len(s2) {
		for i := 0; i < cv.n; i++ {
			if cv.b.a[i] != string(s2)[i] {
				return false
			}
			return true
		}
	}
	return false
}

func (cv SuConcat) PackSize() int {
	if cv.n == 0 {
		return 0
	} else {
		return 1 + cv.n
	}
}

func (cv SuConcat) Pack(buf []byte) []byte {
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

var _ Value = SuConcat{} // confirm it implements Value
