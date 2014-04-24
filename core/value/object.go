// Object is a Suneido object
// i.e. a container with both list and named members

package value

type Object struct {
	list  []Value
	named map[Value]Value
}

func (ob Object) Get(key Value) Value {
	return ob.named[key]
}

func (ob Object) Put(key Value, val Value) {
	ob.named[key] = val
}

func (ob Object) ToInt() int {
	panic("cannot convert object to integer")
}

func (ob Object) ToStr() string {
	panic("cannot convert object to string")
}

var _ Value = Object{} // confirm it implements Value
