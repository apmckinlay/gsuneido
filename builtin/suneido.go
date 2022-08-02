// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func init() {
	SuneidoObjectMethods = Methods{
		"Compile": method("(source, errob = false)",
			func(t *Thread, _ Value, args []Value) Value {
				src := ToStr(args[0])
				if args[1] == False {
					return compile.Constant(src)
				}
				ob := ToContainer(args[1])
				val, checks := compile.Checked(t, src)
				for _, w := range checks {
					ob.Add(SuStr(w))
				}
				return val
			}),
		"Parse": method("(source)",
			func(t *Thread, _ Value, args []Value) Value {
				src := ToStr(args[0])
				p := compile.AstParser(src)
				ast := p.Const()
				if p.Token != tokens.Eof {
					p.Error("did not parse all input")
				}
				return ast
			}),
		// simulate various kinds of errors for testing
		"Crash!": method0(
			func(Value) Value {
				// force a crash, mostly to test output capture
				go func() { panic("Crash!") }()
				return nil
			}),
		"BoundsFail": method0(
			func(Value) Value {
				return []Value{}[1]
			}),
		"AssertFail": method0(
			func(Value) Value {
				assert.That(false)
				return nil
			}),
		"ShouldNotReachHere": method0(
			func(Value) Value {
				assert.ShouldNotReachHere()
				return nil
			}),
	}
}
