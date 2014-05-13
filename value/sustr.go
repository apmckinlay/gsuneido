package value

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hash"
)

// SuStr is a string Value
type SuStr string

var _ Value = SuStr("") // confirm it implements Value

func (ss SuStr) ToInt() int32 {
	i, _ := strconv.ParseInt(string(ss), 0, 32)
	return int32(i)
}

func (ss SuStr) ToDnum() dnum.Dnum {
	dn, err := dnum.Parse(string(ss))
	if err != nil {
		panic("can't convert this string to a number")
	}
	return dn
}

func (ss SuStr) ToStr() string {
	return string(ss)
}

func (ss SuStr) String() string {
	return "'" + string(ss) + "'"
}

func (ss SuStr) Get(key Value) Value {
	return SuStr(string(ss)[key.ToInt()])
}

func (ss SuStr) Put(key Value, val Value) {
	panic("strings do not support put")
}

func (ss SuStr) Hash() uint32 {
	return hash.HashString(string(ss))
}

func (ss SuStr) hash2() uint32 {
	return ss.Hash()
}

func (ss SuStr) Equals(other interface{}) bool {
	if s2, ok := other.(SuStr); ok {
		return ss == s2
	}
	if cv, ok := other.(SuConcat); ok && cv.n == len(ss) {
		for i := 0; i < cv.n; i++ {
			if cv.b.a[i] != string(ss)[i] {
				return false
			}
			return true
		}
	}
	return false
}

func (ss SuStr) PackSize() int {
	if len(ss) == 0 {
		return 0
	} else {
		return 1 + len(ss)
	}
}

func (ss SuStr) Pack(buf []byte) []byte {
	n := len(ss)
	if n == 0 {
		return buf
	}
	i := len(buf)
	buf = buf[:i+1+n]
	buf[i] = PACK_STRING
	copy(buf[i+1:], string(ss))
	return buf
}

func UnpackSuStr(buf []byte) Value {
	return SuStr(string(buf))
}

func (_ SuStr) TypeName() string {
	return "String"
}

func (_ SuStr) order() Order {
	return OrdStr
}

func (c SuStr) cmp(other Value) int {
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
