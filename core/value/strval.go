package value

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hash"
)

// StrVal is a string Value
type StrVal string

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

func (sv StrVal) Hash2() uint32 {
	return sv.Hash()
}

func (sv StrVal) Equals(other interface{}) bool {
	// TODO CatVal
	if s2, ok := other.(StrVal); ok {
		return sv == s2
	}
	return false
}

var _ Value = StrVal("") // confirm it implements Value
