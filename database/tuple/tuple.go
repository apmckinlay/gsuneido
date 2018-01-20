/*
Package tuple implements Suneido's tuples - a sequence of packed values.

Tuple is an interface that is implemented by two main types, TupleM and TupleB.
*/
package tuple

import (
	"bytes"

	v "github.com/apmckinlay/gsuneido/interp"
	"github.com/apmckinlay/gsuneido/util/ints"
)

/*
Tuple is a sequence of packed values.
*/
type Tuple interface {
	// Size returns the number of values
	Size() int
	// GetRaw returns the byte slice containing the i'th value
	GetRaw(i int) []byte
	// Get returns the i'th value (unpacked)
	Get(i int) v.Value
	// Compare returns -1 for less than, 0 for equal, +1 for greater than
	Compare(other Tuple) int
}

// these are standalone functions so they can work on any kind of Tuple
// but they are only exposed as methods that delegate back

func compare(x, y Tuple) int {
	xn := x.Size()
	yn := y.Size()
	for i := 0; i < xn && i < yn; i++ {
		cmp := bytes.Compare(x.GetRaw(i), y.GetRaw(i))
		if cmp != 0 {
			return cmp
		}
	}
	return ints.Compare(xn, yn)
}

func get(t Tuple, i int) v.Value {
	return v.Unpack(t.GetRaw(i))
}
