package builtin

import (
	"fmt"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtinRaw("Print(@args)",
	func(t *Thread, as *ArgSpec, args ...Value) Value {
		iter := NewArgsIter(as, args)
		sep := ""
		for {
			name, value := iter()
			if value == nil {
				break
			}
			fmt.Print(sep)
			if name != nil {
				print(t, name)
				fmt.Print(": ")
			}
			print(t, value)
			sep = " "
		}
		fmt.Println()
		return nil
	})

func print(t *Thread, v Value) {
	if s,ok := v.ToStr(); ok {
		fmt.Print(s)
	} else {
		fmt.Print(display(t, v))
	}
}
