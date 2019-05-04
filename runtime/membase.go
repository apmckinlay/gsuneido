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
	Finder(t *Thread, fn func(v Value, mb *MemBase) Value) Value
}

func (mb *MemBase) AddMembersTo(ob *SuObject) {
	for m := range mb.Data {
		ob.Add(SuStr(m))
	}
}

func (mb *MemBase) Size() int {
	return len(mb.Data)
}

func (mb *MemBase) Copy() MemBase {
	copy := make(map[string]Value, len(mb.Data))
	for k, v := range mb.Data {
		copy[k] = v
	}
	return MemBase{Data: copy}
}
