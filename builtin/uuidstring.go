// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"uuid"

	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(UuidString, "() :string")

func UuidString() Value {
	return SuStr(uuid.New().String())
}
