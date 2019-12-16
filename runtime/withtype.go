// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"fmt"
	"strings"

	cm "github.com/apmckinlay/gsuneido/util/cmatch"

	"github.com/apmckinlay/gsuneido/runtime/types"
)

var binary = cm.InRange(' ', '~').Or(cm.AnyOf("\r\n\t")).Negate()

func WithType(x Value) string {
	if x == nil {
		return "nil"
	}
	var s string
	if ss, ok := x.ToStr(); ok && binary.IndexIn(ss) != -1 {
		s = fmt.Sprintf("%q", ss)
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
