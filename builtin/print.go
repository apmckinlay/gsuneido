package builtin

import (
	"fmt"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = rawbuiltin("Print(@args)",
	func(t *Thread, as *ArgSpec, args ...Value) Value {
		iter := ArgsIter(as, args)
		sep := ""
		for {
			name, value := iter()
			if value == nil {
				break
			}
			fmt.Print(sep)
			if name != nil {
				print(name)
				fmt.Print(": ")
			}
			print(value)
			sep = " "
		}
		fmt.Println()
		return nil
	})

func print(v Value) {
	if ss, ok := v.(SuStr); ok {
		fmt.Print(string(ss))
	} else {
		fmt.Print(v)
	}
}

// ArgsIter returns an iterator function
// that can be called to return successive name,value pairs
// with name = nil for unnamed values
// It returns nil,nil when there are no more values
func ArgsIter(as *ArgSpec, args []Value) func() (Value, Value) {
	if as.Each != 0 {
		iter := args[0].(*SuObject).Iter()
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
