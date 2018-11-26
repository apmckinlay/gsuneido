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
	th := NewThread()
	th.Push(ob) // "this"
	for i := 2; i < len(args)-1; i++ {
		th.Push(toValue(args, str, i))
	}
	f := ob.Lookup(method)
	if f == nil {
		fmt.Print("\tmethod not found: ", method)
		return false
	}
	result := f.Call(th, &ArgSpec{Unnamed: byte(len(args)-3)})
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
