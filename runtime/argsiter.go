// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

type ArgsIter func() (Value, Value)

// NewArgsIter returns an iterator function
// that can be called to return successive name,value pairs
// with name = nil for unnamed values
// It returns nil,nil when there are no more values
func NewArgsIter(as *ArgSpec, args []Value) ArgsIter {
	if as.Each != 0 {
		iter := ToContainer(args[0]).ArgsIter()
		for n := as.Each; n > 1; n-- {
			iter() // skip
		}
		return iter
	}
	next := 0
	return func() (Value, Value) {
		i := next
		if i >= len(args) {
			return nil, nil
		}
		next++
		unnamed := as.Unnamed()
		if i < unnamed {
			return nil, args[i]
		}
		return as.Names[as.Spec[i-unnamed]], args[i]
	}
}
