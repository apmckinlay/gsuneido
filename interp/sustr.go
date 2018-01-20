package interp

import (
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hash"
)

// SuStr is a string Value
type SuStr string

var _ Value = SuStr("")
var _ Packable = SuStr("")

func (ss SuStr) ToInt() int32 {
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

func (ss SuStr) String() string {
	return "'" + string(ss) + "'"
}

func (ss SuStr) Get(key Value) Value {
	return SuStr(string(ss)[index(key)])
}

func index(v Value) int32 {
	if i, ok := v.(SuInt); ok {
		return int32(i)
	}
	if d, ok := v.(SuDnum); ok {
		if d.IsInt() {
			if i, err := d.Int32(); err == nil {
				return i
			}
		}
	}
	panic("string subscripts must be integers")
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

func (_ SuStr) TypeName() string {
	return "String"
}

func (_ SuStr) order() ord {
	return ordStr
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

func (_ SuStr) Call(t *Thread, as ArgSpec) Value {
	panic("String call not implemented")
}
