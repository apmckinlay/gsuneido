package runtime

// Container is an interface to objects and records
type Container interface {
	Value
	Add(val Value)
	ListSize() int
	ListGet(i int) Value
	NamedSize() int
	NamedGet(k Value) Value
	Copy() Container
	Slice(n int) Container
	Clear()
	Insert(at int, val Value)
	Erase(t *Thread, key Value) bool
	Delete(t *Thread, key Value) bool
	GetIfPresent(t *Thread, key Value) Value
	IsReadOnly() bool
	SetReadOnly()
	ArgsIter() func() (Value, Value)
	Iter2(list bool, named bool) func() (Value, Value)
	HasKey(key Value) bool
	ToObject() *SuObject
	ToRecord(t *Thread, hdr *Header) Record
}

// ContainerFind returns the key of the first occurrence of the value and true
// or False,false if not found. The order of named members is not defined.
func ContainerFind(ob Container, val Value) (Value, bool) {
	iter := ob.Iter2(true, true)
	for k, v := iter(); v != nil; k, v = iter() {
		if v.Equal(val) {
			return k, true
		}
	}
	return False, false
}

// iterators

func IterValues(ob Container, list, named bool) Iter {
	return &obIter{ob: ob, list: list, named: named,
		iter:   ob.Iter2(list, named),
		result: func(k, v Value) Value { return v }}
}

func IterMembers(ob Container, list, named bool) Iter {
	return &obIter{ob: ob, list: list, named: named,
		iter:   ob.Iter2(list, named),
		result: func(k, v Value) Value { return k }}
}

func IterAssocs(ob Container, list, named bool) Iter {
	return &obIter{ob: ob, list: list, named: named,
		iter:   ob.Iter2(list, named),
		result: func(k, v Value) Value { return NewSuObject(k, v) }}
}

type obIter struct {
	ob          Container
	list, named bool
	iter        func() (Value, Value)
	result      func(Value, Value) Value
}

func (it *obIter) Next() Value {
	k, v := it.iter()
	if v == nil {
		return nil
	}
	return it.result(k, v)
}
func (it *obIter) Dup() Iter {
	oi := *it
	oi.iter = it.ob.Iter2(it.list, it.named)
	return &oi
}
func (it *obIter) Infinite() bool {
	return false
}
