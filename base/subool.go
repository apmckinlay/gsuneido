package base

import "github.com/apmckinlay/gsuneido/util/dnum"

// SuBool is a boolean Value
type SuBool bool

var _ Value = True
var _ Packable = True

var (
	True  = SuBool(true)
	False = SuBool(false)
)

func (b SuBool) ToInt() int {
	if b == false {
		return 0
	}
	panic("can't convert true to integer")
}

func (b SuBool) ToDnum() dnum.Dnum {
	if b == false {
		return dnum.Zero
	}
	panic("can't convert true to number")
}

func (b SuBool) String() string {
	if b == true {
		return "true"
	}
	return "false"
}

func (b SuBool) ToStr() string {
	return b.String()
}

func (SuBool) Get(Value) Value {
	panic("boolean does not support get")
}

func (SuBool) Put(Value, Value) {
	panic("boolean does not support put")
}

func (SuBool) RangeTo(int, int) Value {
	panic("boolean does not support range")
}

func (SuBool) RangeLen(int, int) Value {
	panic("boolean does not support range")
}

func (b SuBool) Hash() uint32 {
	if b == false {
		return 0x11111111
	}
	return 0x22222222
}

func (b SuBool) hash2() uint32 {
	return b.Hash()
}

func (b SuBool) Equals(other interface{}) bool {
	if b2, ok := other.(SuBool); ok {
		return b == b2
	}
	return false
}

func (SuBool) PackSize() int {
	return 1
}

func (b SuBool) Pack(buf []byte) []byte {
	if b == true {
		buf = append(buf, packTrue)
	} else {
		buf = append(buf, packFalse)
	}
	return buf
}

func (SuBool) TypeName() string {
	return "Boolean"
}

func (SuBool) Order() ord {
	return ordBool
}

func (b SuBool) Cmp(other Value) int {
	if b == other {
		return 0
	} else if b {
		return 1
	} else {
		return -1
	}
}

func (b SuBool) Not() SuBool {
	if b == True {
		return False
	}
	return True
}
