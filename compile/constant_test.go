// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dnum"
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
	test("-0x1", SuInt(-1))
	test("0xfffff", SuDnum{Dnum: dnum.FromInt(0xfffff)})
	test("-0xfffff", SuDnum{Dnum: dnum.FromInt(-0xfffff)})
	test("0377", SuInt(377)) // Suneido does not support octal literals
	test("'hi wo'", SuStr("hi wo"))
	test("#foo", SuStr("foo"))
	test("/* comment */ true", True)

	assert.T(t).This(Constant("#20140425").String()).Is("#20140425")
	assert.T(t).This(Constant("function () {}").Type()).Is(types.Function)
}

func TestConstantBool(t *testing.T) {
	test := func(src, expectedType, expected string) {
		t.Helper()
		actual := Constant(src)
		if actual.Type().String() != expectedType {
			t.Errorf("%q: expected type %q, got %q", src, expectedType, actual.Type().String())
		}
		if Show(actual) != expected {
			t.Errorf("%q: expected %q, got %q", src, expected, Show(actual))
		}
	}
	test("true", "Boolean", "true")
	test("false", "Boolean", "false")
}

func TestConstantNumber(t *testing.T) {
	test := func(src, expectedType, expected string) {
		t.Helper()
		actual := Constant(src)
		if actual.Type().String() != expectedType {
			t.Errorf("%q: expected type %q, got %q", src, expectedType, actual.Type().String())
		}
		if Show(actual) != expected {
			t.Errorf("%q: expected %q, got %q", src, expected, Show(actual))
		}
	}
	test("0", "Number", "0")
	test("123", "Number", "123")
	test("-123", "Number", "-123")
	test("+123", "Number", "123")
	test("123.", "Number", "123")
	test("0xff", "Number", "255")
	test(".001", "Number", ".001")
	test("0.001", "Number", ".001")
	test("1e-3", "Number", ".001")
	test(".1e-2", "Number", ".001")
	test("100000000000000000000", "Number", "1e20")
	test(".00000000000000000001", "Number", "1e-20")
	test("0.00000000000000000001", "Number", "1e-20")
}

func TestConstantNumberErrors(t *testing.T) {
	throws := func(src, expectedErr string) {
		t.Helper()
		err := assert.Catch(func() {
			Constant(src)
		})
		if err == nil {
			t.Errorf("%q: expected error containing %q, got nil", src, expectedErr)
		} else if !strings.Contains(err.(string), expectedErr) {
			t.Errorf("%q: expected error containing %q, got %q", src, expectedErr, err)
		}
	}
	throws("+true", "syntax error")
	throws("-false", "syntax error")
}

func TestConstantDateSymbol(t *testing.T) {
	test := func(src, expectedType, expected string) {
		t.Helper()
		actual := Constant(src)
		if actual.Type().String() != expectedType {
			t.Errorf("%q: expected type %q, got %q", src, expectedType, actual.Type().String())
		}
		if Show(actual) != expected {
			t.Errorf("%q: expected %q, got %q", src, expected, Show(actual))
		}
	}
	test("#20181023", "Date", "#20181023")
	test("#xy", "String", `"xy"`)
	test("#Foo", "String", `"Foo"`)
	test("#_foo", "String", `"_foo"`)
	test("#_X", "String", `"_X"`)
}

func TestConstantObject(t *testing.T) {
	test := func(src, expectedType, expected string) {
		t.Helper()
		actual := Constant(src)
		if actual.Type().String() != expectedType {
			t.Errorf("%q: expected type %q, got %q", src, expectedType, actual.Type().String())
		}
		if Show(actual) != expected {
			t.Errorf("%q: expected %q, got %q", src, expected, Show(actual))
		}
	}
	test("()", "Object", "#()")
	test("{}", "Record", "[]")
	test("[]", "Record", "[]")
	test("#()", "Object", "#()")
	test("#{}", "Record", "[]")
	test("#[]", "Record", "[]")
	test("#(123)", "Object", "#(123)")
	test("#(+123)", "Object", "#(123)")
	test("#(-123)", "Object", "#(-123)")
	test("#(12, 34)", "Object", "#(12, 34)")
	test("#(ab)", "Object", `#("ab")`)
	test("#(a:)", "Object", "#(a:)")
	test("#(a: 123)", "Object", "#(a: 123)")
	test("#(1, 2, a: 3, b: 4)", "Object", "#(1, 2, a: 3, b: 4)")
	test("#(-1: -2)", "Object", "#(-1: -2)")
	test(`#("foo bar": foobar)`, "Object", `#("foo bar": "foobar")`)
	test("#(20181023)", "Object", "#(20181023)")
	test("#(20181023:)", "Object", "#(20181023:)")
	test("#(20181023: false)", "Object", "#(20181023: false)")
	test("#(#20181023)", "Object", "#(#20181023)")
	test("#(#20181023:)", "Object", "#(#20181023:)")
	test("#(#20181023: 123)", "Object", "#(#20181023: 123)")
	test("#(is)", "Object", `#("is")`)
	test("#(is: 123)", "Object", "#(is: 123)")
	test("#(true: 123)", "Object", "#(true: 123)")
	test(`#("true": 123)`, "Object", `#("true": 123)`)
	test("#(class: 123)", "Object", "#(class: 123)")
	test("#(function: 123)", "Object", "#(function: 123)")
	test("#(default: 123)", "Object", "#(default: 123)")
	test("#(function(){})", "Object", "#(/* function */)")
	test("#(class{})", "Object", "#(/* class */)")
	test("#(flag:)", "Object", "#(flag:)")
	test("#(flag: \n )", "Object", "#(flag:)")
	test("[]", "Record", "[]")
	test("[b: 1, c: 2]", "Record", "[b: 1, c: 2]")
	test("[1, 2]", "Object", "#(1, 2)")
}

func TestConstantObjectErrors(t *testing.T) {
	throws := func(src, expectedErr string) {
		t.Helper()
		err := assert.Catch(func() {
			Constant(src)
		})
		if err == nil {
			t.Errorf("%q: expected error containing %q, got nil", src, expectedErr)
		} else if !strings.Contains(err.(string), expectedErr) {
			t.Errorf("%q: expected error containing %q, got %q", src, expectedErr, err)
		}
	}
	throws("#(a: 1, a: 2)", "duplicate member name")
}

func TestConstantClass(t *testing.T) {
	test := func(src, expectedType, expected string) {
		t.Helper()
		actual := Constant(src)
		if actual.Type().String() != expectedType {
			t.Errorf("%q: expected type %q, got %q", src, expectedType, actual.Type().String())
		}
		if Show(actual) != expected {
			t.Errorf("%q: expected %q, got %q", src, expected, Show(actual))
		}
	}
	test("class { }", "Class", "class{}")
	test("Base { }", "Class", "Base{}")
	test("class Base { }", "Class", "Base{}")
	test("class : Base { }", "Class", "Base{}")
	test("class { B: 2 \n A: 1 }", "Class", "class{A: 1; B: 2}")
	test("Base { Foo(){} \n Bar(){} }", "Class", "Base{Bar(); Foo()}")
	test("class { N: ; Foo: function(){} }", "Class", "class{Foo(); N: true}")
}

func TestConstantClassErrors(t *testing.T) {
	throws := func(src, expectedErr string) {
		t.Helper()
		err := assert.Catch(func() {
			Constant(src)
		})
		if err == nil {
			t.Errorf("%q: expected error containing %q, got nil", src, expectedErr)
		} else if !strings.Contains(err.(string), expectedErr) {
			t.Errorf("%q: expected error containing %q, got %q", src, expectedErr, err)
		}
	}
	throws("class { 123 }", "class members must be named")
	throws("class { 123: 456 }", "class member names must be strings")
}

func TestConstantFunction(t *testing.T) {
	test := func(src, expectedType, expected string) {
		t.Helper()
		actual := Constant(src)
		if actual.Type().String() != expectedType {
			t.Errorf("%q: expected type %q, got %q", src, expectedType, actual.Type().String())
		}
		if Show(actual) != expected {
			t.Errorf("%q: expected %q, got %q", src, expected, Show(actual))
		}
	}
	test("function () {}", "Function", "function()")
	test("function (a,b,c) {}", "Function", "function(a,b,c)")
	test("function (@args) {}", "Function", "function(@args)")
	test("function (a,b=1,c=2) {}", "Function", "function(a,b=1,c=2)")
	test("function (a,_b,_c=false) {}", "Function", "function(a,_b,_c=false)")
}

func TestConstantFunctionErrors(t *testing.T) {
	throws := func(src, expectedErr string) {
		t.Helper()
		err := assert.Catch(func() {
			Constant(src)
		})
		if err == nil {
			t.Errorf("%q: expected error containing %q, got nil", src, expectedErr)
		} else if !strings.Contains(err.(string), expectedErr) {
			t.Errorf("%q: expected error containing %q, got %q", src, expectedErr, err)
		}
	}
	throws("function (a,b=1,c) {}", "default parameters must come last")
	throws("function (.a) {}", "dot parameters only allowed in class methods")
	throws("function (a, b, a) {}", "duplicate function parameter")
	throws("function(){ .New() }", "can't explicitly call New method")
}
