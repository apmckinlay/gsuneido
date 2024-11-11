// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

// MemBase is the shared base for SuClass and SuInstance
type MemBase struct {
	Data map[string]Value
	MayLock
}

func NewMemBase() MemBase {
	return MemBase{Data: map[string]Value{}}
}

type Findable interface {
	// Finder applies fn to a MemBase and all its parents
	// stopping if fn returns something other than nil, and returning that value.
	// Implemented by SuClass and SuInstance
	Finder(th *Thread, fn func(v Value, mb *MemBase) Value) Value
}

func (mb *MemBase) AddMembersTo(ob *SuObject) {
	if mb.Lock() {
		defer mb.lock.Unlock()
	}
	for m := range mb.Data {
		ob.Add(SuStr(m))
	}
}

func (mb *MemBase) Size() int {
	if mb.Lock() {
		defer mb.lock.Unlock()
	}
	return len(mb.Data)
}

func (mb *MemBase) Copy() MemBase {
	if mb.Lock() {
		defer mb.lock.Unlock()
	}
	copy := make(map[string]Value, len(mb.Data))
	for k, v := range mb.Data {
		copy[k] = v
	}
	return MemBase{Data: copy}
}

func (mb *MemBase) Has(m string) bool {
	if mb.Lock() {
		defer mb.lock.Unlock()
	}
	_, ok := mb.Data[m]
	return ok
}

var _ iter2able = (*MemBase)(nil)

func (mb *MemBase) Iter2(bool, bool) func() (Value, Value) {
	if mb.Lock() {
		defer mb.lock.Unlock()
	}
	// can't use iter.Pull2 because it requires calling stop
	data := make([]Value, 0, 2*len(mb.Data)) // snapshot
	for k, v := range mb.Data {
		data = append(data, SuStr(k), v)
	}
	i := 0
	return func() (Value, Value) {
		if i >= len(data) {
			return nil, nil
		}
		i += 2
		return data[i-2], data[i-1]
	}
}
