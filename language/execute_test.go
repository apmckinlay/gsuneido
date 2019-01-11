package language

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"

	_ "github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/runtime"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

func TestCall(*testing.T) {
	f := compile.Constant(`function() {
		function () { 123 }(a: 1)
	}`)
	th := NewThread()
	f.Call(th, ArgSpec0)
}

var _ = AddGlobal("Suneido", new(SuObject))
var _ = ptest.Add("execute", pt_execute)
var _ = ptest.Add("lang_rangeto", pt_lang_rangeto)
var _ = ptest.Add("lang_rangelen", pt_lang_rangelen)
var _ = ptest.Add("compare", pt_compare)
var _ = ptest.Add("compare_packed", pt_compare_packed)

func TestBuiltinString(t *testing.T) {
	f := GetGlobal(GlobalNum("Type"))
	Assert(t).That(f.String(), Equals("Type /* builtin function */"))
	f = GetGlobal(GlobalNum("Object"))
	Assert(t).That(f.String(), Equals("Object /* builtin function */"))
}

func TestPtestExecute(t *testing.T) {
	if !ptest.RunFile("execute.test") {
		t.Fail()
	}
}

func TestPtestLang(t *testing.T) {
	if !ptest.RunFile("lang.test") {
		t.Fail()
	}
}

func TestPtestClassImpl(t *testing.T) {
	if !ptest.RunFile("classimpl.test") {
		t.Fail()
	}
}

func init() {
	def := func(nameVal, val Value) Value {
		name := string(nameVal.(SuStr))
		if ss, ok := val.(SuStr); ok {
			val = compile.NamedConstant(name, string(ss))
		}
		TestGlobal(name, val)
		return nil
	}
	AddGlobal("Def", &Builtin2{def, BuiltinParams{ParamSpec: ParamSpec2}})
}

func pt_execute(args []string, _ []bool) bool {
	//fmt.Println(args)
	if strings.Contains(args[0], "Seq(") || strings.Contains(args[0], ".Eval(") {
		fmt.Println("skipped", args) // TODO
		return true
	}

	src := "function () {\n" + args[0] + "\n}"
	th := NewThread()
	expected := "**notfalse**"
	if len(args) > 1 {
		expected = args[1]
	}
	var success bool
	var result interface{}
	if expected == "throws" {
		expected = "throws " + args[2]
		e := Catch(func() {
			fn := compile.Constant(src).(*SuFunc)
			result = th.Call(fn)
		})
		if e == nil {
			success = false
		} else if es, ok := e.(string); ok {
			result = SuStr(es)
			success = strings.Contains(es, args[2])
		} else {
			result = SuStr(fmt.Sprint(e))
			success = false
		}
	} else {
		fn := compile.Constant(src).(*SuFunc)
		actual := th.Call(fn)
		result = actual
		if expected == "**notfalse**" {
			success = actual != False
		} else {
			expectedValue := compile.Constant(expected)
			success = actual.Equal(expectedValue)
		}
	}
	if !success {
		fmt.Println("\tgot:", result)
		fmt.Println("\texpected: " + expected)
	}
	return success
}

func pt_lang_rangeto(args []string, _ []bool) bool {
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

func pt_lang_rangelen(args []string, _ []bool) bool {
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

func pt_compare(args []string, _ []bool) bool {
	n := len(args)
	for i := 0; i < n; i++ {
		x := compile.Constant(args[i])
		if x.Compare(x) != 0 {
			return false
		}
		for j := i + 1; j < n; j++ {
			y := compile.Constant(args[j])
			if x.Compare(y) >= 0 || y.Compare(x) <= 0 {
				fmt.Println(x, "should be less than", y)
				return false
			}
		}
	}
	return true
}

func pt_compare_packed(args []string, _ []bool) bool {
	n := len(args)
	for i := 0; i < n; i++ {
		x := compile.Constant(args[i])
		xp := Pack(x.(Packable))
		for j := i + 1; j < n; j++ {
			y := compile.Constant(args[j])
			yp := Pack(y.(Packable))
			if bytes.Compare(xp, yp) >= 0 || bytes.Compare(yp, xp) <= 0 {
				fmt.Println(x, "should be less than", y)
				return false
			}
		}
	}
	return true
}

// compare to BenchmarkJit in interp_test.go
func BenchmarkInterp(b *testing.B) {
	src := `function (x,y) { x + y }`
	if !GlobalExists("ADD") {
		AddGlobal("ADD", compile.Constant(src).(*SuFunc))
	}
	src = `function () {
		sum = 0
		for (i = 0; i < 100; ++i)
			sum = ADD(sum, i)
		return sum
	}`
	fn := compile.Constant(src).(*SuFunc)
	th := &Thread{}
	for n := 0; n < b.N; n++ {
		th.Reset()
		result := th.Call(fn)
		if !result.Equal(SuInt(4950)) {
			panic("wrong result " + result.String())
		}
	}
}

func BenchmarkCall(b *testing.B) {
	f := GetGlobal(GlobalNum("Type"))
	as := ArgSpec1
	th := NewThread()
	th.Push(SuInt(123))
	for i := 0; i < b.N; i++ {
		f.Call(th, as)
	}
}
