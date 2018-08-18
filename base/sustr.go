package base

import (
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hash"
)

// SuStr is a string Value
type SuStr string

var _ Value = SuStr("")
var _ Packable = SuStr("")

func (ss SuStr) ToInt() int {
	if string(ss) == "" {
		return 0
	}
	panic("can't convert String to number")
}

func (ss SuStr) ToDnum() dnum.Dnum {
	if string(ss) == "" {
		return dnum.Zero
	}
	panic("can't convert String to number")
}

func (ss SuStr) ToStr() string {
	return string(ss)
}

// String returns a human readable string with quotes and escaping
// TODO: handle escaping
func (ss SuStr) String() string {
	return "'" + string(ss) + "'"
}

func (ss SuStr) Get(key Value) Value {
	return SuStr(string(ss)[Index(key)])
}

func (SuStr) Put(Value, Value) {
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
	}
	return 1 + len(ss)
}

func (ss SuStr) Pack(buf []byte) []byte {
	if len(ss) == 0 {
		return buf
	}
	buf = append(buf, packString)
	buf = append(buf, string(ss)...)
	return buf
}

func UnpackSuStr(buf []byte) Value {
	return SuStr(string(buf))
}

func (SuStr) TypeName() string {
	return "String"
}

func (SuStr) Order() ord {
	return ordStr
}

func (ss SuStr) Cmp(other Value) int {
	// COULD optimize this to not convert Concat to string
	s1 := ss.String()
	s2 := other.String()
	switch {
	case s1 < s2:
		return -1
	case s1 > s2:
		return +1
	default:
		return 0
	}
}
