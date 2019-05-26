package runtime

// Container is an interface to objects and records
type Container interface {
	Value
	Add(val Value)
	ListSize() int
	ListGet(i int) Value
	NamedSize() int
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

func containerEqual(x Container, y Container, inProgress pairs) bool {
	if x == y { // pointer comparison
		return true // same object
	}
	if x.ListSize() != y.ListSize() || x.NamedSize() != y.NamedSize() {
		return false
	}
	if inProgress.contains(x, y) {
		return true
	}
	inProgress.push(x, y) // no need to pop due to pass by value
	for i := 0; i < x.ListSize(); i++ {
		if !deepEqual(x.ListGet(i), y.ListGet(i), inProgress) {
			return false
		}
	}
	if x.NamedSize() > 0 {
		iter := x.Iter2(false, true)
		for {
			k, v := iter()
			if k == nil {
				break
			}
			yk := y.GetIfPresent(nil, k) //TODO need thread for record
			if yk == nil || !deepEqual(v.(Value), yk.(Value), inProgress) {
				return false
			}
		}
	}
	return true
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
