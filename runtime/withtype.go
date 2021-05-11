// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/runtime/types"
)

// WithType is used by the repl
func WithType(x Value) string {
	if x == nil {
		return "nil"
	}
	var s string
	if ss, ok := x.ToStr(); ok && !needQuote(ss) {
		s = fmt.Sprint(ss)
	} else {
		s = fmt.Sprint(x)
	}
	if x.Type() != types.Boolean {
		if _, ok := x.(SuStr); !ok {
			t := fmt.Sprintf("%T", x)
			if strings.HasPrefix(t, "runtime.") {
				t = t[8:]
			}
			if strings.HasPrefix(t, "*runtime.") {
				t = "*" + t[9:]
			}
			s += " <" + t + ">"
		}
	}
	return s
}

func needQuote(s string) bool {
	for _, c := range []byte(s) {
		if (c < ' ' || '~' < c) && c != '\t' && c != '\r' && c != '\n' {
			return true
		}
	}
	return false
}
