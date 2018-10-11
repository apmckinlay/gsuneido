package compile

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/base"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestConstant(t *testing.T) {
	test := func(src string, expected Value) {
		Assert(t).That(Constant(src), Equals(expected))
	}
	test("true", True)
	test("false", False)
	test("0", SuInt(0))
	test("-123", SuInt(-123))
	test("+456", SuInt(456))
	test("0xff", SuInt(255))
	test("0377", SuInt(255))
	test("'hi wo'", SuStr("hi wo"))

	Assert(t).That(Constant("#20140425").String(), Equals("#20140425"))
	Assert(t).That(Constant("function () {}").TypeName(), Equals("Function"))
}

func TestConstantObject(t *testing.T) {
	test := func(src string, expected string) {
		//fmt.Println(">>>", src)
		Assert(t).That(Constant(src).String(), Equals(expected))
	}
	test("()", "#()")
	test("{}", "#()")
	test("[]", "#()")
	test("#()", "#()")
	test("#{}", "#()")
	test("#[]", "#()")
	test("#(123)", "#(123)")
	test("#(12, 34)", "#(12, 34)")
	test("#(a:)", "#(a:)")
	test("#(a: 123)", "#(a: 123)")
	test("#(1, 2, a: 3)", "#(1, 2, a: 3)")
	test("#(1 2 a: 3)", "#(1, 2, a: 3)")
	test("#(-1: -1, 'foo bar': foobar)", "#(-1: -1, 'foo bar': 'foobar')")
	test("#(-1: -1, #20140513: 'May 13')", "#(-1: -1, #20140513: 'May 13')")
}
