package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/google/uuid"
)

var _ = builtin0("UuidString()", func() Value {
	return SuStr(uuid.New().String())
})
