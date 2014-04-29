package value

import (
	"bytes"
	"reflect"

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

func (ob *Object) ToInt() int {
	panic("cannot convert object to integer")
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
	hash := ob.Hash2()
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
			hash = 31*hash + k.(Value).Hash2()
			hash = 31*hash + v.(Value).Hash2()
		}
	}
	return hash
}

func (ob *Object) Hash2() uint32 {
	hash := uint32(17)
	hash = 31*hash + uint32(ob.HashSize())
	hash = 31*hash + uint32(ob.ListSize())
	return hash
}

func (ob *Object) Equals(other interface{}) bool {
	// TODO probably want to implement this myself in terms of Equals
	return reflect.DeepEqual(ob, other)
}
