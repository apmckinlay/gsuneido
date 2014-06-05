package tr

import (
	"testing"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

func Test_makset(t *testing.T) {
	Assert(t).That(makset(""), Equals(""))
	Assert(t).That(makset("foo"), Equals("foo"))
	Assert(t).That(makset("-foo"), Equals("-foo"))
	Assert(t).That(makset("foo-"), Equals("foo-"))
	Assert(t).That(makset("-0-9-"), Equals("-0123456789-"))
}

func TestPtest(t *testing.T) {
	if !ptest.RunFile("tr.test") {
		t.Fail()
	}
}
