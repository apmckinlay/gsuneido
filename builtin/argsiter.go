package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

type ArgsIter func() (Value, Value)

// ArgsIter returns an iterator function
// that can be called to return successive name,value pairs
// with name = nil for unnamed values
// It returns nil,nil when there are no more values
func NewArgsIter(as *ArgSpec, args []Value) ArgsIter {
	if as.Each != 0 {
		iter := ToObject(args[0]).Iter()
		if as.Each == EACH1 {
			iter() // skip first
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
