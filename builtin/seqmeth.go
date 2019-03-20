package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

// for SuSequence

func init() {
	SequenceMethods = Methods{
		"Copy": method0(func(this Value) Value {
			return this.(*SuSequence).Copy()
		}),
		"Infinite?": method0(func(this Value) Value {
			return SuBool(this.(*SuSequence).Infinite())
		}),
		"Instantiated?": method0(func(this Value) Value {
			return SuBool(this.(*SuSequence).Instantiated())
		}),
		"Iter": method0(func(this Value) Value {
				iter := this.(*SuSequence).Iter()
				if wi, ok := iter.(*wrapIter); ok {
					return wi.iter
				}
				return SuIter{Iter: iter}
			}),
		"Join": method1("(separator='')", func(this, arg Value) Value {
			seq := this.(*SuSequence)
			separator := IfStr(arg)
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
