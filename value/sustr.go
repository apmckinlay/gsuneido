package value

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hash"
)

// SuStr is a string Value
type SuStr string

var _ Value = SuStr("") // confirm it implements Value

func (sv SuStr) ToInt() int32 {
	i, _ := strconv.ParseInt(string(sv), 0, 32)
	return int32(i)
}

func (sv SuStr) ToDnum() dnum.Dnum {
	dn, err := dnum.Parse(string(sv))
	if err != nil {
		panic("can't convert this string to a number")
	}
	return dn
}

func (sv SuStr) ToStr() string {
	return string(sv)
}

func (sv SuStr) String() string {
	return "'" + string(sv) + "'"
}

func (sv SuStr) Get(key Value) Value {
	return SuStr(string(sv)[key.ToInt()])
}

func (sv SuStr) Put(key Value, val Value) {
	panic("strings do not support put")
}

func (sv SuStr) Hash() uint32 {
	return hash.HashString(string(sv))
}

func (sv SuStr) hash2() uint32 {
	return sv.Hash()
}

func (sv SuStr) Equals(other interface{}) bool {
	if s2, ok := other.(SuStr); ok {
		return sv == s2
	}
	if cv, ok := other.(SuConcat); ok && cv.n == len(sv) {
		for i := 0; i < cv.n; i++ {
			if cv.b.a[i] != string(sv)[i] {
				return false
			}
			return true
		}
	}
	return false
}

func (sv SuStr) PackSize() int {
	if len(sv) == 0 {
		return 0
	} else {
		return 1 + len(sv)
	}
}

func (sv SuStr) Pack(buf []byte) []byte {
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

func UnpackSuStr(buf []byte) Value {
	return SuStr(string(buf))
}
