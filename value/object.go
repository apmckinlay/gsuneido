package value

import (
	"bytes"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hmap"
)
import "unicode"

// Object is a Suneido object
// i.e. a container with both list and hash members
type Object struct {
	list     []Value
	hash     *hmap.Hmap
	readonly bool
}

var _ Value = &Object{} // confirm it implements Value

func (ob *Object) Get(key Value) Value {
	if iv, ok := key.(IntVal); ok {
		i := int(iv)
		if 0 <= i && i <= ob.ListSize() {
			return ob.list[i]
		}
	}
	if ob.hash == nil {
		return nil
	} else {
		return ob.hash.Get(key).(Value)
	}
}

func (ob *Object) Put(key Value, val Value) {
	if iv, ok := key.(IntVal); ok {
		i := int(iv)
		if i == ob.ListSize() {
			ob.Add(val)
			return
		} else if 0 <= i && i < ob.ListSize() {
			ob.list[i] = val
			return
		}
	}
	ob.ensureHash()
	ob.hash.Put(key, val)
}

func (ob *Object) ToInt() int32 {
	panic("cannot convert object to integer")
}

func (ob *Object) ToDnum() dnum.Dnum {
	panic("cannot convert object to number")
}

func (ob *Object) ToStr() string {
	panic("cannot convert object to string")
}

func (ob *Object) ListSize() int {
	return len(ob.list)
}

func (ob *Object) HashSize() int {
	if ob.hash == nil {
		return 0
	} else {
		return ob.hash.Size()
	}
}

// Size returns the number of values in the object
func (ob *Object) Size() int {
	return ob.ListSize() + ob.HashSize()
}

// Add appends a value to the list portion
func (ob *Object) Add(val Value) {
	ob.mustBeMutable()
	ob.list = append(ob.list, val)
	ob.migrate()
}

func (ob *Object) mustBeMutable() {
	if ob.readonly {
		panic("cannot modify readonly object")
	}
}

func (ob *Object) ensureHash() {
	if ob.hash == nil {
		ob.hash = hmap.NewHmap(0)
	}
}

func (ob *Object) migrate() {
	if ob.hash == nil {
		return
	}
	for {
		x := ob.hash.Del(IntVal(ob.ListSize()))
		if x == nil {
			break
		}
		ob.list = append(ob.list, x.(Value))
	}
}

func (ob *Object) String() string {
	var buf bytes.Buffer
	buf.WriteString("#(")
	for _, v := range ob.list {
		buf.WriteString(v.String())
		buf.WriteString(", ")
	}
	if ob.hash != nil {
		iter := ob.hash.Iter()
		for {
			k, v := iter.Next()
			if k == nil {
				break
			}
			if ks, ok := k.(StrVal); ok && isIdentifier(string(ks)) {
				buf.WriteString(string(ks))
			} else {
				buf.WriteString(k.(Value).String())
			}
			buf.WriteString(": ")
			buf.WriteString(v.(Value).String())
			buf.WriteString(", ")
		}
	}
	if buf.Len() > 2 {
		// remove trailing ", "
		buf.Truncate(buf.Len() - 2)
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

func (ob *Object) Hash() uint32 {
	hash := ob.hash2()
	if ob.ListSize() > 0 {
		hash = 31*hash + ob.list[0].Hash()
	}
	if 0 < ob.HashSize() && ob.HashSize() <= 4 {
		iter := ob.hash.Iter()
		for {
			k, v := iter.Next()
			if k == nil {
				break
			}
			hash = 31*hash + k.(Value).hash2()
			hash = 31*hash + v.(Value).hash2()
		}
	}
	return hash
}

// hash2 is shallow so prevents infinite recursion
func (ob *Object) hash2() uint32 {
	hash := uint32(17)
	hash = 31*hash + uint32(ob.HashSize())
	hash = 31*hash + uint32(ob.ListSize())
	return hash
}

func (ob *Object) Equals(other interface{}) bool {
	ob2, ok := other.(*Object)
	if !ok {
		return false
	}
	return equals2(ob, ob2, newpairs())
}

func equals2(x *Object, y *Object, inProgress pairs) bool {
	if x == y { // pointer comparison
		return true // same object
	}
	if x.ListSize() != y.ListSize() || x.HashSize() != y.HashSize() {
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
	if x.HashSize() > 0 {
		iter := x.hash.Iter()
		for {
			k, v := iter.Next()
			if k == nil {
				break
			}
			if !equals3(v.(Value), y.hash.Get(k).(Value), inProgress) {
				return false
			}
		}
	}
	return true
}

func equals3(x Value, y Value, inProgress pairs) bool {
	xo, xok := x.(*Object)
	if !xok {
		return x.Equals(y)
	}
	yo, yok := y.(*Object)
	if !yok {
		return false
	}
	return equals2(xo, yo, inProgress)

}
