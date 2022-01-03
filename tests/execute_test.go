// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package language

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

func TestNaming(t *testing.T) {
	th := &Thread{}
	test := func(src, expected string) {
		t.Helper()
		c := compile.Constant("function () {\n" + src + "\n}").(*SuFunc)
		result := th.Invoke(c, nil)
		assert.T(t).This(result).Is(SuStr(expected))
	}
	test(`foo = function(){}; Name(foo)`, "foo")
	test(`foo = class{}; Name(foo)`, "foo")
	test(`foo = bar = class{}; Name(bar)`, "bar")
	test(`Def('Tmp', 'function(){}'); Name(Tmp)`, "Tmp")
	test(`Def('Tmp', 'function(){ return function(){} }'); Name(Tmp())`, "Tmp")
	test(`Def('Tmp', 'function(){ return {} }'); Name(Tmp())`, "Tmp")
	test(`Def('Tmp', 'function(){ fn = function(){} }'); Name(Tmp())`, "Tmp fn")
	test(`Def('Tmp', 'function(){ b = {} }'); Name(Tmp())`, "Tmp b")
	test(`Def('Tmp', 'class { F(){} }'); Name(Tmp.F)`, "Tmp.F")
	test(`Def('Tmp', 'class { Inner: class { F(){} } }');
		Name(Tmp.Inner.F)`, "Tmp.Inner.F")
	test(`Def('Tmp', 'function(){ myclass = class { F(){} } }');
		Name(Tmp().F)`, "Tmp myclass.F")
	test(`Def('Tmp', 'function() { Object(class{}) }'); Name(Tmp()[0])`,
		"Tmp")
	test(`Def('Tmp', 'class { A() { class { B(){} } } }'); Name(Tmp.A().B)`,
		"Tmp.A.B")
}

func BenchmarkCat(b *testing.B) {
	c := compile.Constant(
		`function ()
			{
			s = ''
			for (i = 0; i < 1000; ++i)
				s $= "abc"
			}`).(*SuFunc)
	th := &Thread{}
	for i := 0; i < b.N; i++ {
		th.Reset()
		th.Invoke(c, nil)
	}
}

func BenchmarkJoin(b *testing.B) {
	c := compile.Constant(
		`function ()
			{
			ob = Object()
			for (i = 0; i < 1000; ++i)
				ob.Add("abc")
			ob.Join()
			}`).(*SuFunc)
	th := &Thread{}
	for i := 0; i < b.N; i++ {
		th.Reset()
		th.Invoke(c, nil)
	}
}

func BenchmarkBase(b *testing.B) {
	c := compile.Constant(
		`function ()
			{
			for (i = 0; i < 1000; ++i)
				;
			}`).(*SuFunc)
	th := &Thread{}
	for i := 0; i < b.N; i++ {
		th.Reset()
		th.Invoke(c, nil)
	}
}

var _ = Global.Add("Suneido", new(SuObject))
var _ = ptest.Add("execute", pt_execute)
var _ = ptest.Add("lang_rangeto", pt_lang_rangeto)
var _ = ptest.Add("lang_rangelen", pt_lang_rangelen)
var _ = ptest.Add("compare", pt_compare)
var _ = ptest.Add("compare_packed", pt_compare_packed)

func TestBuiltinString(t *testing.T) {
	f := Global.GetName(nil, "Type")
	assert.T(t).This(f.String()).Is("Type /* builtin function */")
	f = Global.GetName(nil, "Object")
	assert.T(t).This(f.String()).Is("Object /* builtin function */")
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
	builtin.Def()
}

func pt_execute(args []string, _ []bool) bool {
	src := "function () {\n" + args[0] + "\n}"
	th := &Thread{}
	expected := "**notfalse**"
	if len(args) > 1 {
		expected = args[1]
	}
	var success bool
	var actual Value
	if expected == "throws" {
		expected = "throws " + args[2]
		e := assert.Catch(func() {
			fn := compile.Constant(src).(*SuFunc)
			actual = th.Invoke(fn, nil)
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
		actual = th.Invoke(fn, nil)
		if actual == nil {
			success = expected == "nil"
		} else if expected == "**notfalse**" {
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
		result := th.Invoke(fn, nil)
		if !result.Equal(SuInt(4950)) {
			panic("wrong result " + result.String())
		}
	}
}

func BenchmarkCall(b *testing.B) {
	f := Global.GetName(nil, "Type")
	as := &ArgSpec1
	th := &Thread{}
	th.Push(SuInt(123))
	for i := 0; i < b.N; i++ {
		f.Call(th, nil, as)
	}
}

func TestCoverage(t *testing.T) {
	atomic.StoreInt64(&options.Coverage, 1)
	fn := compile.Constant(`function()
		{
		x = 0
		for (i = 0; i < 10; ++i)
			x += i
		return x
		}`).(*SuFunc)
	fn.StartCoverage(true)
	th := &Thread{}
	th.Invoke(fn, nil)
	cover := fn.StopCoverage()
	assert.T(t).This(cover).
		Is(compile.Constant("#(17: 1, 25: 1, 53: 10, 62: 1)").(*SuObject))
}
