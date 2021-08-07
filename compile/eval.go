// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/regex"
)

// EvalString executes string containing Suneido code
// i.e. string.Eval()
func EvalString(t *Thread, s string) Value {
	s = strings.Trim(s, " \t\r\n")
	if isGlobal(s) {
		// optimize if just a global name
		return Global.GetName(t, s)
	}
	s = "function () {\n" + s + "\n}"
	fn := NamedConstant("", "eval", s).(*SuFunc)
	return t.Invoke(fn, nil)
}

var rxGlobal = regex.Compile(`\A[A-Z][_a-zA-Z0-9]*?[!?]?\Z`)

func isGlobal(s string) bool {
	return rxGlobal.Matches(s)
}
