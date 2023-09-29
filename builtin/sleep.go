// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"time"

	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(Sleep, "(ms)")

func Sleep(arg Value) Value {
	ms := ToInt(arg)
	time.Sleep(time.Duration(ms) * time.Millisecond)
	return nil
}
