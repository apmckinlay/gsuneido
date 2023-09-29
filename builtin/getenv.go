// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"os"

	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(Getenv, "(string)")

func Getenv(arg Value) Value {
	return SuStr(os.Getenv(ToStr(arg)))
}
