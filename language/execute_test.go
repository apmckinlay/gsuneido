package language

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/base"
	"github.com/apmckinlay/gsuneido/interp/global"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

var _ = global.Add("Suneido", new(SuObject))

func TestPtest(t *testing.T) {
	if !ptest.RunFile("execute.test") {
		t.Fail()
	}
}
