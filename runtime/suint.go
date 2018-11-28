package runtime

// WARNING: reflect.DeepEqual (tests) does not work correctly with SuInt
// For pointers it compares what they point to
// Which for smi is always zero
// As a partial fix, smi's from -127 to 127 are set to themselves
// So as long as tests stick to small values they are ok

import (
	"math"
	"strconv"
	"unsafe"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/ints"
)

type smi byte

const smiRange = 1 << 16
const MinSuInt = math.MinInt16
const MaxSuInt = math.MaxInt16

var smispace [smiRange]smi // uninitialized BSS, no actual memory used
var smibase = uintptr(unsafe.Pointer(&smispace[0]))

func init() {
	// this is so that reflect.DeepEquals doesn't think small smi's are equal
	// (for tests)
	for i := -127; i < 127; i++ {
		*SuInt(i) = smi(i) // +1 to avoid zero
	}
}

// SuInt converts an int to *smi which implements Value
// will panic if out of int16 range
func SuInt(n int) *smi {
	offset := n - math.MinInt16
	return &smispace[offset] // will panic if out of range
}

// SmiToInt converts to int if its argument is a *smi
func SmiToInt(x interface{}) (int, bool) {
	if si, ok := x.(*smi); ok {
		return si.ToInt(), true
	}
	return 0, false
}

var _ Value = SuInt(0)
var _ Packable = SuInt(0)

func (si *smi) ToInt() int {
	p := unsafe.Pointer(si)
	offset := int(uintptr(p) - smibase)
	return offset + math.MinInt16
}

func (si *smi) ToDnum() dnum.Dnum {
	return dnum.FromInt(int64(si.ToInt()))
}

func (si *smi) ToStr() string {
	return strconv.Itoa(si.ToInt())
}

func (si *smi) String() string {
	return si.ToStr()
}

func (*smi) Get(Value) Value {
	panic("number does not support get")
}

func (*smi) Put(Value, Value) {
	panic("number does not support put")
}

func (*smi) RangeTo(int, int) Value {
	panic("number does not support range")
}

func (*smi) RangeLen(int, int) Value {
	panic("number does not support range")
}

func (si *smi) Hash() uint32 {
	return uint32(si.ToInt())
}

func (si *smi) Hash2() uint32 {
	return si.Hash()
}

func (si *smi) Equal(other interface{}) bool {
	if i2, ok := other.(*smi); ok {
		return si == i2
	} else if dn, ok := other.(SuDnum); ok {
		return 0 == dnum.Compare(si.ToDnum(), dn.Dnum)
	}
	return false
}

func (si *smi) PackSize() int {
	return PackSizeInt64(int64(si.ToInt()))
}

func (si *smi) Pack(buf []byte) []byte {
	return PackInt64(int64(si.ToInt()), buf)
}

func (*smi) TypeName() string {
	return "Number"
}

func (*smi) Order() Ord {
	return ordNum
}

func (si *smi) Compare(other Value) int {
	if cmp := ints.Compare(si.Order(), other.Order()); cmp != 0 {
		return cmp
	}
	if y, ok := other.(*smi); ok {
		return ints.Compare(si.ToInt(), y.ToInt())
	}
	return dnum.Compare(si.ToDnum(), other.ToDnum())
}

func (*smi) Call(*Thread, *ArgSpec) Value {
	panic("can't call Number")
}

// IntMethods is initialized by the builtin package
var IntMethods Methods

func (*smi) Lookup(method string) Value {
	if m := IntMethods[method]; m != nil {
		return m
	}
	return NumMethods[method]
}
