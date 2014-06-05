package main

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/interp"
	"github.com/apmckinlay/gsuneido/util/ptest"
	"github.com/apmckinlay/gsuneido/value"
)

var _ = ptest.Add("execute", pt_execute)

func pt_execute(args []string) bool {
	src := "function () {\n" + args[0] + "\n}"
	//fmt.Println(src)
	fn := compile.Constant(src).(*value.SuFunc)
	th := interp.NewThread()
	result := th.Call(fn, interp.SimpleArgSpecs[0])
	expected := "true"
	if len(args) > 1 {
		expected = args[1]
	}
	if result.String() != expected {
		fmt.Println("got: " + result.String())
		fmt.Println("expected: " + expected)
	}
	return result.String() == expected
}
