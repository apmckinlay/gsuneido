// StrVal is a string Value

package value

import "strconv"

type StrVal string

func (sv StrVal) ToInt() int {
	i, _ := strconv.ParseInt(string(sv), 0, 32)
	return int(i)
}

func (sv StrVal) ToStr() string {
	return string(sv)
}

func (sv StrVal) Get(key Value) Value {
	return StrVal(string(sv)[key.ToInt()])
}

func (sv StrVal) Put(key Value, val Value) {
	panic("strings do not support put")
}

var _ Value = StrVal("") // confirm it implements Value
