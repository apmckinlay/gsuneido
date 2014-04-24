// IntVal is an integer Value

package value

import "strconv"

type IntVal int

func (iv IntVal) ToInt() int {
	return int(iv)
}

func (iv IntVal) ToStr() string {
	return strconv.Itoa(int(iv))
}

func (iv IntVal) Get(key Value) Value {
	panic("number does not support get")
}

func (iv IntVal) Put(key Value, val Value) {
	panic("number does not support put")
}

var _ Value = IntVal(0) // confirm it implements Value
