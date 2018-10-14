package base

import (
	"strings"
	"unicode"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hmap"
	"github.com/apmckinlay/gsuneido/util/ints"
)

// SuObject is a Suneido object
//
// i.e. a container with both list and named members
//
// NOTE: Not thread safe
type SuObject struct {
	list     []Value
	named    hmap.Hmap
	readonly bool
}

var _ Value = (*SuObject)(nil)

//TODO var _ Packable = &SuObject{}

// Get returns the value associated with a key, or nil if not found
func (ob *SuObject) Get(key Value) Value {
	if i := index(key); 0 <= i && i < ob.ListSize() {
		return ob.list[i]
	}
	x := ob.named.Get(key)
	if x == nil {
		return nil
	}
	return x.(Value)
}

func index(key Value) int {
	if i, ok := SmiToInt(key); ok {
		return i
	}
	if dn, ok := key.(SuDnum); ok {
		if i, ok := dn.Dnum.ToInt(); ok {
			return i
		}
	}
	return -1 // invalid list index
}

// ListGet returns a value from the list, panics if index out of range
func (ob *SuObject) ListGet(i int) Value {
	return ob.list[i]
}

// Put adds or updates the given key and value
// The value will be added to the list if the key is the "next"
func (ob *SuObject) Put(key Value, val Value) {
	i := index(key)
	if i == ob.ListSize() {
		ob.Add(val)
		return
	} else if 0 <= i && i < ob.ListSize() {
		ob.list[i] = val
		return
	}
	ob.named.Put(key, val)
}

func (ob *SuObject) RangeTo(from int, to int) Value {
	size := ob.Size()
	from = prepFrom(from, size)
	to = prepTo(from, to, size)
	return ob.rangeTo(from, to)
}

func (ob *SuObject) RangeLen(from int, n int) Value {
	size := ob.Size()
	from = prepFrom(from, size)
	n = prepLen(n, size-from)
	return ob.rangeTo(from, from+n)
}

func (ob *SuObject) rangeTo(i int, j int) *SuObject {
	list := make([]Value, j-i)
	copy(list, ob.list[i:j])
	return &SuObject{list: list}
}

func (*SuObject) ToInt() int {
	panic("cannot convert object to integer")
}

func (*SuObject) ToDnum() dnum.Dnum {
	panic("cannot convert object to number")
}

func (*SuObject) ToStr() string {
	panic("cannot convert object to string")
}

func (ob *SuObject) ListSize() int {
	return len(ob.list)
}

func (ob *SuObject) NamedSize() int {
	return ob.named.Size()
}

// Size returns the number of values in the object
func (ob *SuObject) Size() int {
	return ob.ListSize() + ob.NamedSize()
}

// Add appends a value to the list portion
func (ob *SuObject) Add(val Value) {
	ob.mustBeMutable()
	ob.list = append(ob.list, val)
	ob.migrate()
}

func (ob *SuObject) mustBeMutable() {
	if ob.readonly {
		panic("cannot modify readonly object")
	}
}

func (ob *SuObject) migrate() {
	for {
		x := ob.named.Del(NumFromInt(ob.ListSize()))
		if x == nil {
			break
		}
		ob.list = append(ob.list, x.(Value))
	}
}

func (ob *SuObject) String() string {
	var buf strings.Builder
	sep := ""
	buf.WriteString("#(")
	for _, v := range ob.list {
		buf.WriteString(sep)
		sep = ", "
		buf.WriteString(v.String())
	}
	iter := ob.named.Iter()
	for {
		k, v := iter()
		if k == nil {
			break
		}
		buf.WriteString(sep)
		sep = ", "
		if ks, ok := k.(SuStr); ok && isIdentifier(string(ks)) {
			buf.WriteString(string(ks))
		} else {
			buf.WriteString(k.(Value).String())
		}
		buf.WriteString(":")
		if v != True {
			buf.WriteString(" ")
			buf.WriteString(v.(Value).String())
		}
	}
	buf.WriteString(")")
	return buf.String()
}

func isIdentifier(s string) bool {
	last := len(s) - 1
	if last < 0 {
		return false
	}
	for i, c := range s {
		if !(c == '_' || unicode.IsLetter(c) ||
			(i > 0 && '0' <= c && c <= '9') ||
			(i == last && (c == '?' || c == '!'))) {
			return false
		}
	}
	return true
}

func (ob *SuObject) Hash() uint32 {
	hash := ob.Hash2()
	if ob.ListSize() > 0 {
		hash = 31*hash + ob.list[0].Hash()
	}
	if 0 < ob.NamedSize() && ob.NamedSize() <= 4 {
		iter := ob.named.Iter()
		for {
			k, v := iter()
			if k == nil {
				break
			}
			hash = 31*hash + k.(Value).Hash2()
			hash = 31*hash + v.(Value).Hash2()
		}
	}
	return hash
}

// Hash2 is shallow so prevents infinite recursion
func (ob *SuObject) Hash2() uint32 {
	hash := uint32(17)
	hash = 31*hash + uint32(ob.NamedSize())
	hash = 31*hash + uint32(ob.ListSize())
	return hash
}

func (ob *SuObject) Equal(other interface{}) bool {
	ob2, ok := other.(*SuObject)
	if !ok {
		return false
	}
	return equals2(ob, ob2, newpairs())
}

func equals2(x *SuObject, y *SuObject, inProgress pairs) bool {
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
		if !equals3(x.list[i], y.list[i], inProgress) {
			return false
		}
	}
	if x.NamedSize() > 0 {
		iter := x.named.Iter()
		for {
			k, v := iter()
			if k == nil {
				break
			}
			yk := y.named.Get(k)
			if yk == nil || !equals3(v.(Value), yk.(Value), inProgress) {
				return false
			}
		}
	}
	return true
}

func equals3(x Value, y Value, inProgress pairs) bool {
	xo, xok := x.(*SuObject)
	if !xok {
		return x.Equal(y)
	}
	yo, yok := y.(*SuObject)
	if !yok {
		return false
	}
	return equals2(xo, yo, inProgress)
}

func (SuObject) TypeName() string {
	return "Object"
}

func (SuObject) Order() Ord {
	return ordObject
}

func (ob *SuObject) Compare(other Value) int {
	if cmp := ints.Compare(ob.Order(), other.Order()); cmp != 0 {
		return cmp
	}
	return cmp2(ob, other.(*SuObject), newpairs())
}

func cmp2(x *SuObject, y *SuObject, inProgress pairs) int {
	if x == y { // pointer comparison
		return 0
	}
	if inProgress.contains(x, y) {
		return 0
	}
	inProgress.push(x, y) // no need to pop due to pass by value
	for i := 0; i < x.Size() && i < y.Size(); i++ {
		if cmp := cmp3(x.ListGet(i), y.ListGet(i), inProgress); cmp != 0 {
			return cmp
		}
	}
	return ints.Compare(x.Size(), y.Size())
}

func cmp3(x Value, y Value, inProgress pairs) int {
	xo, xok := x.(*SuObject)
	yo, yok := y.(*SuObject)
	if !xok || !yok {
		return x.Compare(y)
	}
	return cmp2(xo, yo, inProgress)
}

// Slice returns a copy of the object, with the first n list elements removed
func (ob *SuObject) Slice(n int) *SuObject {
	newNamed := ob.named.Copy()

	if n > len(ob.list) {
		return &SuObject{named: *newNamed, readonly: false}
	}
	list := make([]Value, len(ob.list)-n)
	copy(list, ob.list[n:])
	return &SuObject{list: list, named: *newNamed, readonly: false}
}
