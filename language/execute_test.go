package language

import (
	"fmt"
	"strconv"
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
var _ = ptest.Add("lang_rangeto", pt_lang_rangeto)
var _ = ptest.Add("lang_rangelen", pt_lang_rangelen)

func TestPtest(t *testing.T) {
	if !ptest.RunFile("execute.test") || !ptest.RunFile("lang.test") {
		//t.Fail()
	}
}

func pt_execute(args []string) bool {
	//fmt.Println(args)
	src := "function () {\n" + args[0] + "\n}"
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
			fn := compile.Constant(src).(*SuFunc)
			result = th.Call(fn, nil)
		})
		if e == nil {
			ok = false
		} else {
			result = SuStr(e.(string))
			ok = strings.Contains(e.(string), args[2])
		}
	} else {
		fn := compile.Constant(src).(*SuFunc)
		result = th.Call(fn, nil)
		if expected == "**notfalse**" {
			ok = result != False
		} else {
			expectedValue := compile.Constant(expected)
			ok = result.Equal(expectedValue)
		}
	}
	if !ok {
		fmt.Println("\tgot: " + result.String())
		fmt.Println("\texpected: " + expected)
	}
	return ok
}

func pt_lang_rangeto(args []string) bool {
	s := args[0]
	from, _ := strconv.Atoi(args[1])
	to, _ := strconv.Atoi(args[2])
	expected := SuStr(args[3])
	actual := SuStr(s).RangeTo(from, to)
	if !actual.Equal(expected) {
		fmt.Println("\tgot:", actual)
		fmt.Println("\texpected:", expected)
		return false
	}
	return true
}

func pt_lang_rangelen(args []string) bool {
	s := args[0]
	from, _ := strconv.Atoi(args[1])
	n := 9999
	if len(args) == 4 {
		n, _ = strconv.Atoi(args[2])
	}
	expected := args[len(args)-1]
	actual := SuStr(s).RangeLen(from, n)
	if !actual.Equal(SuStr(expected)) {
		fmt.Println("\tgot:", actual)
		fmt.Println("\texpected:", expected)
		return false
	}

	list := strToList(s)
	expectedList := strToList(expected)
	actualList := list.RangeLen(from, n)
	if !actualList.Equal(expectedList) {
		fmt.Println("\tgot:", actualList)
		fmt.Println("\texpected:", expectedList)
		return false
	}

	return true
}

func strToList(s string) *SuObject {
	ob := SuObject{}
	for _, c := range s {
		ob.Add(SuStr(string(c)))
	}
	return &ob
}
