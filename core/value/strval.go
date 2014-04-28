// StrVal is a string Value

package value

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/util/hash"
)

type StrVal string

func (sv StrVal) ToInt() int {
	i, _ := strconv.ParseInt(string(sv), 0, 32)
	return int(i)
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
	return hash.Hash(string(sv))
}

func (sv StrVal) Hash2() uint32 {
	return sv.Hash()
}

func (sv StrVal) Equals(other interface{}) bool {
	return sv == other.(StrVal)
}

var _ Value = StrVal("") // confirm it implements Value
