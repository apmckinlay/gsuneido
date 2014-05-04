package value

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/util/dnum"
)

// IntVal is an integer Value
type IntVal int32

var _ Value = IntVal(0) // confirm it implements Value

func (iv IntVal) ToInt() int32 {
	return int32(iv)
}

func (iv IntVal) ToDnum() dnum.Dnum {
	return dnum.FromInt64(int64(iv))
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

func (iv IntVal) hash2() uint32 {
	return iv.Hash()
}

func (iv IntVal) Equals(other interface{}) bool {
	if i2, ok := other.(IntVal); ok {
		return iv == i2
	}
	return false
}

func (iv IntVal) Sign() int {
	if iv < 0 {
		return -1
	} else if iv > 0 {
		return 1
	} else {
		return 0
	}
}

func (iv IntVal) Exp() int {
	return 0
}

func (iv IntVal) Coef() uint64 {
	if int32(iv) >= 0 {
		return uint64(iv)
	} else {
		return uint64(-int64(iv))
	}
}

func (iv IntVal) PackSize() int {
	return packSizeNum(iv)
}

func (iv IntVal) Pack(buf []byte) []byte {
	return packNum(iv, buf)
}
