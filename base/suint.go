package base

import (
	"math"
	"strconv"
	"unsafe"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/ints"
)

type smi byte

const smiRange = 1 << 16

var space [smiRange]smi
var base = uintptr(unsafe.Pointer(&space[0]))

// SuInt converts an int to *smi which implements Value
// will panic if out of int16 range
func SuInt(n int) *smi {
	offset := int(n) - math.MinInt16
	return &space[offset] // will panic if out of range
}

// Su2Int converts to int if possible
func Su2Int(x interface{}) (int, bool) {
	if si, ok := x.(*smi); ok {
		return int(si.ToInt()), true
	}
	return 0, false
}

var _ Value = SuInt(0)
var _ Packable = SuInt(0)

func (si *smi) ToInt() int32 {
	p := unsafe.Pointer(si)
	offset := int32(uintptr(p) - base)
	return offset + math.MinInt16
}

func (si *smi) ToDnum() dnum.Dnum {
	return dnum.FromInt64(int64(si.ToInt()))
}

func (si *smi) ToStr() string {
	return strconv.Itoa(int(si.ToInt()))
}

func (si *smi) String() string {
	return si.ToStr()
}

func (si *smi) Get(key Value) Value {
	panic("number does not support get")
}

func (si *smi) Put(key Value, val Value) {
	panic("number does not support put")
}

func (si *smi) Hash() uint32 {
	return uint32(si.ToInt())
}

func (si *smi) hash2() uint32 {
	return si.Hash()
}

func (si *smi) Equals(other interface{}) bool {
	if i2, ok := other.(*smi); ok {
		return si == i2
	} else if dn, ok := other.(SuDnum); ok {
		return 0 == dnum.Cmp(si.ToDnum(), dn.Dnum)
	}
	return false
}

func (si *smi) Sign() int {
	i := si.ToInt()
	if i < 0 {
		return -1
	} else if i > 0 {
		return 1
	} else {
		return 0
	}
}

func (si *smi) Exp() int {
	return 0
}

func (si *smi) Coef() uint64 {
	i := si.ToInt()
	if i >= 0 {
		return uint64(i)
	}
	return uint64(-int64(i))
}

func (si *smi) PackSize() int {
	return packSizeNum(si)
}

func (si *smi) Pack(buf []byte) []byte {
	return packNum(si, buf)
}

func (*smi) TypeName() string {
	return "Number"
}

func (*smi) Order() ord {
	return ordNum
}

func (x *smi) Cmp(other Value) int {
	if y, ok := other.(*smi); ok {
		return ints.Compare(int(x.ToInt()), int(y.ToInt()))
	}
	return dnum.Cmp(x.ToDnum(), other.ToDnum())
}
