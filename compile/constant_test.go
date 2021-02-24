// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

func TestConstant(t *testing.T) {
	test := func(src string, expected Value) {
		t.Helper()
		assert.T(t).This(Constant(src)).Is(expected)
	}
	test("true", True)
	test("false", False)
	test("0", SuInt(0))
	test("-123", SuInt(-123))
	test("+456", SuInt(456))
	test("0xff", SuInt(255))
	test("0xfffff", SuDnum{Dnum: dnum.FromInt(0xfffff)})
	test("0xffffffff", SuDnum{Dnum: dnum.FromInt(-1)})
	test("0377", SuInt(377))
	test("'hi wo'", SuStr("hi wo"))
	test("#foo", SuStr("foo"))
	test("/* comment */ true", True)

	assert.T(t).This(Constant("#20140425").String()).Is("#20140425")
	assert.T(t).This(Constant("function () {}").Type()).Is(types.Function)
}

// ptest ------------------------------------------------------------

var _ = ptest.Add("compile", pt_compile)

func TestPtestConstant(t *testing.T) {
	if !ptest.RunFile("constant.test") {
		t.Fail()
	}
}

func pt_compile(args []string, _ []bool) bool {
	expectedType := args[1]
	expected := args[2]
	var actual Value
	ok := true
	e := assert.Catch(func() {
		actual = Constant(args[0])
		if actual.Type().String() != expectedType {
			ok = false
		}
		if Show(actual) != expected {
			ok = false
		}
		if !ok {
			fmt.Println("\tgot:", "<"+actual.Type().String()+">", Show(actual))
		}
	})
	if e != nil {
		if _, str := e.(string); !str {
			fmt.Println(e)
			ok = false
		} else if expectedType != "throws" ||
			!strings.Contains(e.(string), expected) {
			fmt.Println("\tgot:", e)
			ok = false
		}
	}
	return ok
}
