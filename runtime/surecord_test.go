package runtime

import (
	"testing"

	"github.com/apmckinlay/gsuneido/runtime/types"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestSuRecord(t *testing.T) {
	r := new(SuRecord)
	Assert(t).That(r.Type(), Equals(types.Record))
	Assert(t).That(r.String(), Equals("[]"))
	r.Set(SuStr("a"), SuInt(123))
	Assert(t).That(r.String(), Equals("[a: 123]"))
}
