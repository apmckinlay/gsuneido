// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/google/uuid"
)

var _ = builtin(UuidString, "()")

func UuidString() Value {
	return SuStr(uuid.New().String())
}
