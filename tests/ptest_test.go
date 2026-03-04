// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tests

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

func init() {
	builtin.DefDef()
}

func TestPtestLang(t *testing.T) {
	if !ptest.RunFile("lang.test") {
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

func TestPtestClass(t *testing.T) {
	if !ptest.RunFile("class.test") {
		t.Fail()
	}
}

func TestPtestExecute(t *testing.T) {
	if !ptest.RunFile("execute.test") {
		t.Fail()
	}
}

func TestMethodPtest(t *testing.T) {
	if !ptest.RunFile("number.test") {
		t.Fail()
	}
}

// ------------------------------------------------------------------

var _ = ptest.Add("execute", ptExecute)

func ptExecute(args []string, _ []bool) bool {
	src := "function () {\n" + args[0] + "\n}"
	var th Thread
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

var _ = ptest.Add("lang_rangeto", ptLangRangeto)

func ptLangRangeto(args []string, _ []bool) bool {
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

var _ = ptest.Add("lang_rangelen", ptLangRangelen)

func ptLangRangelen(args []string, _ []bool) bool {
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

var _ = ptest.Add("compare", ptCompare)

func ptCompare(args []string, _ []bool) bool {
	n := len(args)
	for i := range n {
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

var _ = ptest.Add("compare_packed", ptComparePacked)

func ptComparePacked(args []string, _ []bool) bool {
	n := len(args)
	for i := range n {
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
	switch s {
	case "inf":
		return Inf
	case "-inf":
		return NegInf
	}
	return compile.Constant(s)
}

var _ = ptest.Add("method", ptMethod)

func ptMethod(args []string, str []bool) bool {
	ob := toValue(args, str, 0)
	method := args[1]
	expected := toValue(args, str, len(args)-1)
	f := ob.Lookup(nil, method)
	if f == nil {
		fmt.Print("\tmethod not found: ", method)
		return false
	}
	th := &Thread{}
	for i := 2; i < len(args)-1; i++ {
		th.Push(toValue(args, str, i))
	}
	nargs := len(args) - 3
	result := f.Call(th, ob, &StdArgSpecs[nargs])
	ok := result.Equal(expected)
	if !ok {
		fmt.Printf("\tgot: %s  expected: %s\n",
			WithType(result), WithType(expected))
	}
	return ok
}

func toValue(args []string, str []bool, i int) Value {
	if str[i] {
		return SuStr(args[i])
	}
	return compile.Constant(args[i])
}
