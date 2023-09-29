// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import "fmt"

var printBuiltin = &SuBuiltinRaw{printBuiltinFn, BuiltinParams{ParamSpec: ParamSpecAt}}

func printBuiltinFn(th *Thread, as *ArgSpec, args []Value) Value {
	iter := NewArgsIter(as, args)
	sep := ""
	for {
		name, value := iter()
		if value == nil {
			break
		}
		fmt.Print(sep)
		if name != nil {
			print(th, name)
			fmt.Print(": ")
		}
		print(th, value)
		sep = " "
	}
	fmt.Println()
	return nil
}

func print(th *Thread, v Value) {
	if s, ok := v.ToStr(); ok {
		fmt.Print(s)
	} else {
		fmt.Print(Display(th, v))
	}
}

type Displayable interface {
	Display(th *Thread) string
}

func Display(th *Thread, val Value) string {
	if d, ok := val.(Displayable); ok {
		return d.Display(th)
	}
	if d, ok := val.(ToStringable); ok {
		return d.ToString(th)
	}
	return val.String()
}
