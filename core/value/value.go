package value

import "strconv"
import "fmt"

// Value is used to reference a Suneido value
type Value interface {
	ToStr() string
	ToInt() int
	Get(key Value) Value
	Put(key Value, val Value)
}

// IntValue is an integer Value
type IntValue struct {
	i int
	NoGetPut
}

func (iv IntValue) ToInt() int {
	return iv.i
}

func (iv IntValue) ToStr() string {
	return strconv.Itoa(iv.i)
}

// StrValue is a string Value
type StrValue string

func (sv StrValue) ToInt() int {
	i, _ := strconv.ParseInt(string(sv), 0, 32)
	return int(i)
}

func (sv StrValue) ToStr() string {
	return string(sv)
}

func (sv StrValue) Get(key Value) Value {
	switch key := key.(type) {
	case IntValue:
		if 0 <= key.i && key.i < len(string(sv)) {
			return StrValue(string(sv)[key.i])
		} else {
			return StrValue("")
		}
	default:
		panic("string subscripts must be integers")
	}

}

func (nc StrValue) Put(key Value, val Value) {
	panic("strings do not support put")
}

// Object is a Suneido object
// i.e. a container with both list and named members
type Object struct {
	list  []Value
	named map[Value]Value
	NoConv
}

func (ob Object) Get(key Value) Value {
	return ob.named[key]
}

func (ob Object) Put(key Value, val Value) {
	ob.named[key] = val
}

// NoConv is embedded in a Value struct that does not support ToStr or ToInt
// It defines these methods to call panic.
type NoConv struct {
}

func (nc NoConv) ToInt() int {
	panic(fmt.Sprintf("can't convert to integer"))
}

func (nc NoConv) ToStr() string {
	panic("can't convert to string")
}

// NoGetPut is embedded in a Value that does not support Get or Put
// It defines these methods to call panic.
type NoGetPut struct {
}

func (nc NoGetPut) Get(key Value) Value {
	panic("doesn't support get")
}

func (nc NoGetPut) Put(key Value, val Value) {
	panic("doesn't support put")
}
