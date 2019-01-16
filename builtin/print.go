package builtin

import (
	"fmt"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = rawbuiltin("Print(@args)",
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
