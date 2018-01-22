package base

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/ints"
)

// SuInt is an integer Value
type SuInt int32

var _ Value = SuInt(0)
var _ Packable = SuInt(0)

func (si SuInt) ToInt() int32 {
	return int32(si)
}

func (si SuInt) ToDnum() dnum.Dnum {
	return dnum.FromInt64(int64(si))
}

func (si SuInt) ToStr() string {
	return strconv.Itoa(int(si))
}

func (si SuInt) String() string {
	return si.ToStr()
}

func (si SuInt) Get(key Value) Value {
	panic("number does not support get")
}

func (si SuInt) Put(key Value, val Value) {
	panic("number does not support put")
}

func (si SuInt) Hash() uint32 {
	return uint32(si)
}

func (si SuInt) hash2() uint32 {
	return si.Hash()
}

func (si SuInt) Equals(other interface{}) bool {
	if i2, ok := other.(SuInt); ok {
		return si == i2
	} else if dn, ok := other.(SuDnum); ok {
		return 0 == dnum.Cmp(si.ToDnum(), dn.Dnum)
	}
	return false
}

func (si SuInt) Sign() int {
	if si < 0 {
		return -1
	} else if si > 0 {
		return 1
	} else {
		return 0
	}
}

func (si SuInt) Exp() int {
	return 0
}

func (si SuInt) Coef() uint64 {
	if int32(si) >= 0 {
		return uint64(si)
	} else {
		return uint64(-int64(si))
	}
}

func (si SuInt) PackSize() int {
	return packSizeNum(si)
}

func (si SuInt) Pack(buf []byte) []byte {
	return packNum(si, buf)
}

func (_ SuInt) TypeName() string {
	return "Number"
}

func (_ SuInt) Order() ord {
	return ordNum
}

func (x SuInt) Cmp(other Value) int {
	if y, ok := other.(SuInt); ok {
		return ints.Compare(int(x.ToInt()), int(y.ToInt()))
	} else {
		return dnum.Cmp(x.ToDnum(), other.ToDnum())
	}
}
