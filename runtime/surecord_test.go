package runtime

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestSuRecord(t *testing.T) {
	r := new(SuRecord)
	Assert(t).That(r.TypeName(), Equals("Record"))
	Assert(t).That(r.String(), Equals("[]"))
	r.Add(Zero)
	r.Put(SuStr("a"), SuInt(123))
	Assert(t).That(r.String(), Equals("[0, a: 123]"))
}
