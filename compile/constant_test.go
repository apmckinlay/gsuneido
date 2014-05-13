package compile

import (
	"testing"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	v "github.com/apmckinlay/gsuneido/value"
)

func TestConstant(t *testing.T) {
	test := func(src string, expected v.Value) {
		Assert(t).That(Constant(src), Equals(expected))
	}
	test("true", v.True)
	test("false", v.False)
	test("0", v.SuInt(0))
	test("-123", v.SuInt(-123))
	test("+456", v.SuInt(456))
	test("0xff", v.SuInt(255))
	test("0377", v.SuInt(255))
	test("'hi wo'", v.SuStr("hi wo"))

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
	test("#(1, 2, a: 3, b: 4)", "#(1, 2, a: 3, b: 4)")
	test("#(1 2 a: 3 b: 4)", "#(1, 2, a: 3, b: 4)")
	test("#(-1: -1, +2: +2, 'foo bar': foobar, #20140513: 'May 13')",
		"#(-1: -1, 2: 2, 'foo bar': 'foobar', #20140513: 'May 13')")
}
