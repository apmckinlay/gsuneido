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
	finder(fn func(*MemBase) Value) Value
}

func (ob *MemBase) Members() *SuObject { // TODO sequence
	mems := new(SuObject)
	for m := range ob.Data {
		mems.Add(SuStr(m))
	}
	return mems
}

func (ob *MemBase) Size() Value {
	return IntVal(len(ob.Data))
}

func MemberQ(ob Findable, mem Value) Value {
	if ss, ok := mem.(SuStr); ok {
		m := string(ss)
		result := ob.finder(func(ob *MemBase) Value {
			if _, ok := ob.Data[m]; ok {
				return True
			}
			return nil
		})
		if result == True {
			return True
		}
	}
	return False
}
