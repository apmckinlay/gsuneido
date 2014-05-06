package value

import "github.com/apmckinlay/gsuneido/util/dnum"

// SuBool is a boolean Value
type SuBool bool

var (
	True  = SuBool(true)
	False = SuBool(false)
)

func (bv SuBool) ToInt() int32 {
	if bv == true {
		return 1
	} else {
		return 0
	}
}

func (bv SuBool) ToDnum() dnum.Dnum {
	if bv == true {
		return dnum.One
	} else {
		return dnum.Zero
	}
}

func (bv SuBool) ToStr() string {
	if bv == true {
		return "true"
	} else {
		return "false"
	}
}

func (bv SuBool) String() string {
	return bv.ToStr()
}

func (bv SuBool) Get(key Value) Value {
	panic("boolean does not support get")
}

func (bv SuBool) Put(key Value, val Value) {
	panic("boolean does not support put")
}

func (bv SuBool) Hash() uint32 {
	return uint32(bv.ToInt())
}

func (bv SuBool) hash2() uint32 {
	return bv.Hash()
}

func (bv SuBool) Equals(other interface{}) bool {
	if b2, ok := other.(SuBool); ok {
		return bv == b2
	}
	return false
}

func (bv SuBool) PackSize() int {
	return 1
}

func (bv SuBool) Pack(buf []byte) []byte {
	i := len(buf)
	buf = buf[:i+1]
	if bv == true {
		buf[i] = TRUE
	} else {
		buf[i] = FALSE
	}
	return buf
}

var _ Value = SuBool(true) // confirm it implements Value
