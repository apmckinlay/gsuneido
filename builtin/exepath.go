// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"os"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin(ExePath, "()")

func ExePath() Value {
	path, err := os.Executable()
	if err != nil {
		panic("ExePath " + err.Error())
	}
	path = strings.ReplaceAll(path, "\\", "/")
	return SuStr(path)
}
