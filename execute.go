package main

import (
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/interp"
	"github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/ptest"
	"github.com/apmckinlay/gsuneido/value"
)

var _ = ptest.Add("execute", pt_execute)

func pt_execute(args []string) bool {
	//fmt.Println(args)
	src := "function () {\n" + args[0] + "\n}"
	fn := compile.Constant(src).(*value.SuFunc)
	th := interp.NewThread()
	expected := "**notfalse**"
	if len(args) > 1 {
		expected = args[1]
	}
	var ok bool
	var result value.Value
	if expected == "throws" {
		expected = "throws " + args[2]
		e := hamcrest.Catch(func() {
			result = th.Call(fn, interp.SimpleArgSpecs[0])
		})
		if e == nil {
			ok = false
		} else {
			result = value.SuStr(e.(string))
			ok = strings.Contains(e.(string), args[2])
		}
	} else {
		result = th.Call(fn, interp.SimpleArgSpecs[0])
		if expected == "**notfalse**" {
			ok = result != value.False
		} else {
			expectedValue := compile.Constant(expected)
			ok = result.Equals(expectedValue)
		}
	}
	if !ok {
		fmt.Println("got: " + result.String())
		fmt.Println("expected: " + expected)
	}
	return ok
}
