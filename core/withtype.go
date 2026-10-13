// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/core/types"
)

// WithType is used by the repl
func WithType(x Value) string {
	if x == nil {
		return "nil"
	}
	s := x.String()
	if x.Type() != types.Boolean {
		t := fmt.Sprintf("%T", x)
		t = strings.TrimPrefix(t, "runtime.")
		if strings.HasPrefix(t, "*runtime.") {
			t = "*" + t[9:]
		}
		s += " <" + t + ">"
	}
	return s
}
