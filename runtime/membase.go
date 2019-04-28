package runtime

// MemBase is the shared base for SuClass and SuInstance
type MemBase struct {
	Data map[string]Value
	CantConvert
}

func NewMemBase() MemBase {
	return MemBase{Data: map[string]Value{}}
}

type Findable interface {
	// Finder applies fn to ob and all its parents
	// stopping if fn returns something other than nil, and returning that value.
	// Implemented by SuClass and SuInstance
	Finder(fn func(v Value, mb *MemBase) Value) Value
}

func (ob *MemBase) Members() *SuObject { // TODO sequence
	mems := new(SuObject)
	for m := range ob.Data {
		mems.Add(SuStr(m))
	}
	return mems
}

func (ob *MemBase) Size() int {
	return len(ob.Data)
}

func (ob *MemBase) Copy() MemBase {
	copy := make(map[string]Value, len(ob.Data))
	for k, v := range ob.Data {
		copy[k] = v
	}
	return MemBase{Data: copy}
}
