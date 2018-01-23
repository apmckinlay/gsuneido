package language

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/base"
	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/interp"
	"github.com/apmckinlay/gsuneido/interp/global"
	"github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

var _ = global.Add("Suneido", new(SuObject))
var _ = ptest.Add("execute", pt_execute)

func TestPtest(t *testing.T) {
	if !ptest.RunFile("execute.test") {
		t.Fail()
	}
}

func pt_execute(args []string) bool {
	//fmt.Println(args)
	src := "function () {\n" + args[0] + "\n}"
	fn := compile.Constant(src).(*SuFunc)
	th := interp.NewThread()
	expected := "**notfalse**"
	if len(args) > 1 {
		expected = args[1]
	}
	var ok bool
	var result Value
	if expected == "throws" {
		expected = "throws " + args[2]
		e := hamcrest.Catch(func() {
			result = th.Call(fn, interp.ArgSpec{})
		})
		if e == nil {
			ok = false
		} else {
			result = SuStr(e.(string))
			ok = strings.Contains(e.(string), args[2])
		}
	} else {
		result = th.Call(fn, interp.ArgSpec{})
		if expected == "**notfalse**" {
			ok = result != False
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
