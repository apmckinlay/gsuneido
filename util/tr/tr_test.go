package tr

import (
	"testing"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

func Test_makset(t *testing.T) {
	test := func(s, expected string) {
		Assert(t).That(string(Set("")), Equals(""))
	}
	test("foo", "foo")
	test("-foo", "-foo")
	test("foo-", "foo-")
	test("-0-9-", "-0123456789-")
}

func TestPtest(t *testing.T) {
	if !ptest.RunFile("tr.test") {
		t.Fail()
	}
}
