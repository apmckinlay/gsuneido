package value

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hash"
)

// StrVal is a string Value
type StrVal string

var _ Value = StrVal("") // confirm it implements Value

func (sv StrVal) ToInt() int32 {
	i, _ := strconv.ParseInt(string(sv), 0, 32)
	return int32(i)
}

func (sv StrVal) ToDnum() dnum.Dnum {
	dn, err := dnum.Parse(string(sv))
	if err != nil {
		panic("can't convert this string to a number")
	}
	return dn
}

func (sv StrVal) ToStr() string {
	return string(sv)
}

func (sv StrVal) String() string {
	return "'" + string(sv) + "'"
}

func (sv StrVal) Get(key Value) Value {
	return StrVal(string(sv)[key.ToInt()])
}

func (sv StrVal) Put(key Value, val Value) {
	panic("strings do not support put")
}

func (sv StrVal) Hash() uint32 {
	return hash.HashString(string(sv))
}

func (sv StrVal) hash2() uint32 {
	return sv.Hash()
}

func (sv StrVal) Equals(other interface{}) bool {
	if s2, ok := other.(StrVal); ok {
		return sv == s2
	}
	if cv, ok := other.(CatVal); ok && cv.n == len(sv) {
		for i := 0; i < cv.n; i++ {
			if cv.b.a[i] != string(sv)[i] {
				return false
			}
			return true
		}
	}
	return false
}

func (sv StrVal) PackSize() int {
	if len(sv) == 0 {
		return 0
	} else {
		return 1 + len(sv)
	}
}

func (sv StrVal) Pack(buf []byte) []byte {
	n := len(sv)
	if n == 0 {
		return buf
	}
	i := len(buf)
	buf = buf[:i+1+n]
	buf[i] = STRING
	copy(buf[i+1:], string(sv))
	return buf
}

func UnpackStrVal(buf []byte) Value {
	return StrVal(string(buf))
}
