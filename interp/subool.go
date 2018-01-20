package interp

import "github.com/apmckinlay/gsuneido/util/dnum"

// SuBool is a boolean Value
type SuBool bool

var _ Value = True
var _ Packable = True

var (
	True  = SuBool(true)
	False = SuBool(false)
)

func (b SuBool) ToInt() int32 {
	if b == false {
		return 0
	}
	panic("can't convert true to number")
}

func (b SuBool) ToDnum() dnum.Dnum {
	if b == false {
		return dnum.Zero
	}
	panic("can't convert true to number")
}

func (b SuBool) ToStr() string {
	if b == true {
		return "true"
	} else {
		return "false"
	}
}

func (b SuBool) String() string {
	return b.ToStr()
}

func (b SuBool) Get(key Value) Value {
	panic("boolean does not support get")
}

func (b SuBool) Put(key Value, val Value) {
	panic("boolean does not support put")
}

func (b SuBool) Hash() uint32 {
	return uint32(b.ToInt())
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

func (b SuBool) PackSize() int {
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

func (_ SuBool) TypeName() string {
	return "Boolean"
}

func (_ SuBool) order() ord {
	return ordBool
}

func (b SuBool) cmp(other Value) int {
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
	} else {
		return True
	}
}

func (_ SuBool) Call(t *Thread, as ArgSpec) Value {
	panic("can't call Boolean")
}
