package value

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
	if b == true {
		return 1
	} else {
		return 0
	}
}

func (b SuBool) ToDnum() dnum.Dnum {
	if b == true {
		return dnum.One
	} else {
		return dnum.Zero
	}
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
		buf = append(buf, PACK_TRUE)
	} else {
		buf = append(buf, PACK_FALSE)
	}
	return buf
}

func (_ SuBool) TypeName() string {
	return "Boolean"
}

func (_ SuBool) order() Order {
	return OrdBool
}

func (b SuBool) cmp(other Value) int {
	return int(b.ToInt() - other.(SuBool).ToInt())
}
