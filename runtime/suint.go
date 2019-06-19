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

	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/pack"
)

type smi byte

const MinSuInt = math.MinInt16
const MaxSuInt = math.MaxInt16

var smispace [1 << 16]smi // uninitialized BSS, no actual memory used
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

// SuIntToInt converts to int if its argument is *smi
func SuIntToInt(x interface{}) (int, bool) {
	if si, ok := x.(*smi); ok {
		return si.toInt(), true
	}
	return 0, false
}

func (si *smi) toInt() int {
	p := unsafe.Pointer(si)
	offset := int(uintptr(p) - smibase)
	return offset + math.MinInt16
}

// Value interface --------------------------------------------------

var _ Value = (*smi)(nil)

func (si *smi) ToInt() (int, bool) {
	return si.toInt(), true
}

func (si *smi) IfInt() (int, bool) {
	return si.toInt(), true
}

func (si *smi) ToDnum() (dnum.Dnum, bool) {
	return dnum.FromInt(int64(ToInt(si))), true
}

func (*smi) ToContainer() (Container, bool) {
	return nil, false
}

func (si *smi) AsStr() (string, bool) {
	return si.String(), true
}

func (si *smi) ToStr() (string, bool) {
	return "", false
}

func (si *smi) String() string {
	return strconv.Itoa(si.toInt())
}

func (*smi) Get(*Thread, Value) Value {
	panic("number does not support get")
}

func (*smi) Put(*Thread, Value, Value) {
	panic("number does not support put")
}

func (*smi) RangeTo(int, int) Value {
	panic("number does not support range")
}

func (*smi) RangeLen(int, int) Value {
	panic("number does not support range")
}

func (si *smi) Hash() uint32 {
	return uint32(si.toInt())
}

func (si *smi) Hash2() uint32 {
	return si.Hash()
}

func (si *smi) Equal(other interface{}) bool {
	if i2, ok := other.(*smi); ok {
		return si == i2
	} else if dn, ok := other.(SuDnum); ok {
		dn2, _ := si.ToDnum()
		return 0 == dnum.Compare(dn2, dn.Dnum)
	}
	return false
}

func (*smi) Type() types.Type {
	return types.Number
}

func (si *smi) Compare(other Value) int {
	if cmp := ints.Compare(ordNum, Order(other)); cmp != 0 {
		return cmp
	}
	if y, ok := other.(*smi); ok {
		return ints.Compare(si.toInt(), y.toInt())
	}
	dn, _ := si.ToDnum()
	return dnum.Compare(dn, ToDnum(other))
}

func (*smi) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call Number")
}

// IntMethods is initialized by the builtin package
var IntMethods Methods

var anSuDnum = SuDnum{}

func (*smi) Lookup(t *Thread, method string) Callable {
	if m := IntMethods[method]; m != nil {
		return m
	}
	return anSuDnum.Lookup(t, method)
}

// Packable interface -----------------------------------------------

// TODO: avoid conversion to Dnum

var _ Packable = SuInt(0)

func (si *smi) PackSize(nest int) int {
	return SuDnum{Dnum: dnum.FromInt(int64(si.toInt()))}.PackSize(nest)
}

func (si *smi) Pack(buf *pack.Encoder) {
	SuDnum{Dnum: dnum.FromInt(int64(si.toInt()))}.Pack(buf)
}
