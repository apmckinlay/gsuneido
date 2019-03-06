package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	SequenceMethods = Methods{
		"Join": method1("(separator='')", func(this, arg Value) Value {
			seq := this.(*SuSequence)
			separator := ToStr(arg)
			sep := ""
			iter := seq.Iter()
			var buf strings.Builder
			for {
				val := iter.Next()
				if val == nil {
					break
				}
				buf.WriteString(sep)
				sep = separator
				if s, ok := val.IfStr(); ok {
					buf.WriteString(s)
				} else {
					buf.WriteString(val.String())
				}
			}
			return SuStr(buf.String())
		}),
	}
}
