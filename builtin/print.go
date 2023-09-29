// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"io"
	"os"

	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(PrintStdout, "(string)")

func PrintStdout(_ *Thread, args []Value) Value {
	s := ToStr(args[0])
	io.WriteString(os.Stdout, s)
	return nil
}
