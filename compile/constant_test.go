package compile

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/dnum"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/ptest"
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
	test("0xfffff", SuDnum{dnum.FromInt(0xfffff)})
	test("0xffffffff", SuDnum{dnum.FromInt(-1)})
	test("0377", SuInt(255))
	test("'hi wo'", SuStr("hi wo"))
	test("/* comment */ true", True)

	Assert(t).That(Constant("#20140425").String(), Equals("#20140425"))
	Assert(t).That(Constant("function () {}").TypeName(), Equals("Function"))
}

var _ = ptest.Add("compile", pt_compile)

func TestPtest(t *testing.T) {
	if !ptest.RunFile("constant.test") {
		t.Fail()
	}
}

func pt_compile(args []string, _ []bool) bool {
	expectedType := args[1]
	expected := args[2]
	var actual Value
	ok := true
	e := Catch(func() {
		actual = Constant(args[0])
		if actual.TypeName() != expectedType {
			ok = false
		}
		if Show(actual) != expected {
			ok = false
		}
		if !ok {
			fmt.Println("\tgot:", "<"+actual.TypeName()+">", Show(actual))
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
