package language

import (
	"fmt"
	"testing"

	_ "github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

var _ = ptest.Add("method", pt_method)

func TestMethodPtest(*testing.T) {
	if !ptest.RunFile("number.test") {
		//t.Fail()
	}
}

func pt_method(args []string, str []bool) bool {
	ob := toValue(args, str, 0)
	method := args[1]
	expected := toValue(args, str, len(args)-1)
	f := ob.Lookup(method)
	if f == nil {
		fmt.Print("\tmethod not found: ", method)
		return false
	}
	th := NewThread()
	var result Value
	switch len(args) - 3 {
	case 0:
		result = f.Call1(th, ob)
	case 1:
		result = f.Call2(th, ob, toValue(args, str, 2))
	}
	ok := result.Equal(expected)
	if !ok {
		fmt.Printf("\tgot: %v %#v", result, result)
	}
	return ok
}

func toValue(args []string, str []bool, i int) Value {
	if str[i] {
		return SuStr(args[i])
	}
	return compile.Constant(args[i])
}
