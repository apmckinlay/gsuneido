package value

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/util/dnum"
)

// SuInt is an integer Value
type SuInt int32

var _ Value = SuInt(0) // confirm it implements Value

func (iv SuInt) ToInt() int32 {
	return int32(iv)
}

func (iv SuInt) ToDnum() dnum.Dnum {
	return dnum.FromInt64(int64(iv))
}

func (iv SuInt) ToStr() string {
	return strconv.Itoa(int(iv))
}

func (iv SuInt) String() string {
	return iv.ToStr()
}

func (iv SuInt) Get(key Value) Value {
	panic("number does not support get")
}

func (iv SuInt) Put(key Value, val Value) {
	panic("number does not support put")
}

func (iv SuInt) Hash() uint32 {
	return uint32(iv)
}

func (iv SuInt) hash2() uint32 {
	return iv.Hash()
}

func (iv SuInt) Equals(other interface{}) bool {
	if i2, ok := other.(SuInt); ok {
		return iv == i2
	}
	return false
}

func (iv SuInt) Sign() int {
	if iv < 0 {
		return -1
	} else if iv > 0 {
		return 1
	} else {
		return 0
	}
}

func (iv SuInt) Exp() int {
	return 0
}

func (iv SuInt) Coef() uint64 {
	if int32(iv) >= 0 {
		return uint64(iv)
	} else {
		return uint64(-int64(iv))
	}
}

func (iv SuInt) PackSize() int {
	return packSizeNum(iv)
}

func (iv SuInt) Pack(buf []byte) []byte {
	return packNum(iv, buf)
}
