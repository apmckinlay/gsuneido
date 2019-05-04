package language

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/runtime"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

func TestNaming(t *testing.T) {
	test := func(src, expected string) {
		t.Helper()
		c := compile.Constant("function () {\n" + src + "\n}").(*SuFunc)
		result := NewThread().Call(c)
		Assert(t).That(result, Equals(SuStr(expected)))
	}
	test(`foo = function(){}; Name(foo)`, "foo")
	test(`foo = class{}; Name(foo)`, "foo")
	test(`foo = bar = class{}; Name(bar)`, "bar")
	test(`Def('Tmp', 'class { F(){} }'); Name(Tmp.F)`, "Tmp.F")
	test(`Def('Tmp', 'class { Inner: class { F(){} } }');
		Name(Tmp.Inner.F)`, "Tmp.Inner.F")
	test(`Def('Tmp', 'function(){ myclass = class { F(){} } }');
		Name(Tmp().F)`, "Tmp myclass.F")
	test(`Def('Tmp', 'function() { Object(class{}) }'); Name(Tmp()[0])`,
		"Tmp ?")
	test(`Def('Tmp', 'class { A() { class { B(){} } } }'); Name(Tmp.A().B)`,
		"Tmp.A ?.B")
}

var _ = Global.Add("Suneido", new(SuObject))
var _ = ptest.Add("execute", pt_execute)
var _ = ptest.Add("lang_rangeto", pt_lang_rangeto)
var _ = ptest.Add("lang_rangelen", pt_lang_rangelen)
var _ = ptest.Add("compare", pt_compare)
var _ = ptest.Add("compare_packed", pt_compare_packed)

func TestBuiltinString(t *testing.T) {
	f := Global.GetName(nil, "Type")
	Assert(t).That(f.String(), Equals("Type /* builtin function */"))
	f = Global.GetName(nil, "Object")
	Assert(t).That(f.String(), Equals("Object /* builtin function */"))
}

// func TestTmp(t *testing.T) {
// 	args := []string{`f = function (b){ c={}; b() }; f({ return 123 }); 456`, `123`}
// 	strq := []bool{}
// 	if ! pt_execute(args, strq) {
// 		t.Fail()
// 	}
// }

func TestPtestExecute(t *testing.T) {
	if !ptest.RunFile("execute.test") {
		t.Fail()
	}
}

func TestPtestStrings(t *testing.T) {
	if !ptest.RunFile("strings.test") {
		t.Fail()
	}
}

func TestPtestObjects(t *testing.T) {
	if !ptest.RunFile("objects.test") {
		t.Fail()
	}
}

func TestPtestRecords(t *testing.T) {
	if !ptest.RunFile("records.test") {
		t.Fail()
	}
}

func TestPtestLang(t *testing.T) {
	if !ptest.RunFile("lang.test") {
		t.Fail()
	}
}

func TestPtestClass(t *testing.T) {
	if !ptest.RunFile("class.test") {
		t.Fail()
	}
}

func init() {
	Def()
}

func pt_execute(args []string, _ []bool) bool {
	src := "function () {\n" + args[0] + "\n}"
	th := NewThread()
	expected := "**notfalse**"
	if len(args) > 1 {
		expected = args[1]
	}
	var success bool
	var actual Value
	if expected == "throws" {
		expected = "throws " + args[2]
		e := Catch(func() {
			fn := compile.Constant(src).(*SuFunc)
			actual = th.Call(fn)
		})
		if e == nil {
			success = false
		} else if es, ok := e.(string); ok {
			actual = SuStr(es)
			success = strings.Contains(es, args[2])
		} else if ss, ok := e.(SuStr); ok {
			actual = ss
			success = strings.Contains(string(ss), args[2])
		} else if se, ok := e.(*SuExcept); ok {
			actual = se.SuStr
			success = strings.Contains(string(se.SuStr), args[2])
		} else {
			actual = SuStr(fmt.Sprintf("%#v", e))
			success = false
		}
	} else {
		fn := compile.Constant(src).(*SuFunc)
		actual = th.Call(fn)
		if expected == "**notfalse**" {
			success = actual != False
		} else {
			expectedValue := compile.Constant(expected)
			success = actual.Equal(expectedValue)
			expected = WithType(expectedValue)
		}
	}
	if !success {
		fmt.Printf("\tgot: %s  expected: %s\n", WithType(actual), expected)
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
		fmt.Printf("\tgot: %v  expected: %v\n", actual, expected)
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
		fmt.Printf("\tgot: %v  expected: %v\n", actual, expected)
		return false
	}

	list := strToList(s)
	expectedList := strToList(expected)
	actualList := list.RangeLen(from, n)
	if !actualList.Equal(expectedList) {
		fmt.Printf("\tgot: %v  expected: %v\n", actualList, expectedList)
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
		x := constant(args[i])
		if x.Compare(x) != 0 {
			return false
		}
		for j := i + 1; j < n; j++ {
			y := constant(args[j])
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
		x := constant(args[i])
		xp := Pack(x.(Packable))
		x2 := Unpack(xp)
		if !x.Equal(x2) {
			fmt.Println("pack/unpack, got:", x2, " expected:", x)
		}
		for j := i + 1; j < n; j++ {
			y := constant(args[j])
			yp := Pack(y.(Packable))
			if strings.Compare(xp, yp) >= 0 || strings.Compare(yp, xp) <= 0 {
				fmt.Println(x, "should be less than", y)
				return false
			}
		}
	}
	return true
}

func constant(s string) Value {
	if s == "inf" {
		return Inf
	} else if s == "-inf" {
		return NegInf
	}
	return compile.Constant(s)
}

// compare to BenchmarkJit in interp_test.go
func BenchmarkInterp(b *testing.B) {
	src := `function (x,y) { x + y }`
	if !Global.Exists("ADD") {
		Global.Add("ADD", compile.Constant(src).(*SuFunc))
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
	f := Global.GetName(nil, "Type")
	as := ArgSpec1
	th := NewThread()
	th.Push(SuInt(123))
	for i := 0; i < b.N; i++ {
		f.Call(th, as)
	}
}
