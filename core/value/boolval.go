package value

// BoolVal is a boolean Value
type BoolVal bool

var (
	True  = BoolVal(true)
	False = BoolVal(false)
)

func (bv BoolVal) ToInt() int {
	if bv == true {
		return 1
	} else {
		return 0
	}
}

func (bv BoolVal) ToStr() string {
	if bv == true {
		return "true"
	} else {
		return "false"
	}
}

func (bv BoolVal) String() string {
	return bv.ToStr()
}

func (bv BoolVal) Get(key Value) Value {
	panic("boolean does not support get")
}

func (bv BoolVal) Put(key Value, val Value) {
	panic("boolean does not support put")
}

func (bv BoolVal) Hash() uint32 {
	return uint32(bv.ToInt())
}

func (bv BoolVal) Hash2() uint32 {
	return bv.Hash()
}

func (bv BoolVal) Equals(other interface{}) bool {
	if b2, ok := other.(BoolVal); ok {
		return bv == b2
	}
	return false
}

var _ Value = BoolVal(true) // confirm it implements Value
