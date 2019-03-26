package runtime

import (
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/runtime/types"
)

func WithType(x Value) string {
	s := fmt.Sprint(x)
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
