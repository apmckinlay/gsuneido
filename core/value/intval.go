package value

import "strconv"

// IntVal is an integer Value
type IntVal int32

func (iv IntVal) ToInt() int {
	return int(iv)
}

func (iv IntVal) ToStr() string {
	return strconv.Itoa(int(iv))
}

func (iv IntVal) String() string {
	return iv.ToStr()
}

func (iv IntVal) Get(key Value) Value {
	panic("number does not support get")
}

func (iv IntVal) Put(key Value, val Value) {
	panic("number does not support put")
}

func (iv IntVal) Hash() uint32 {
	return uint32(iv)
}

func (iv IntVal) Hash2() uint32 {
	return iv.Hash()
}

func (iv IntVal) Equals(other interface{}) bool {
	return iv == other.(IntVal)
}

var _ Value = IntVal(0) // confirm it implements Value
