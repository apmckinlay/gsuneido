// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestExecTool(t *testing.T) {
	assert := assert.T(t)
	{
		result, err := execTool("1 + 2 \n")
		assert.That(err == nil)
		assert.That(len(result.Warnings) == 0)
		assert.This(result.Results).Is("[3]")
	}
	{
		result, err := execTool("return")
		assert.That(err == nil)
		assert.That(len(result.Warnings) == 0)
		assert.This(result.Results).Is("[]")
	}
	{
		result, err := execTool("return 123")
		assert.That(err == nil)
		assert.That(len(result.Warnings) == 0)
		assert.This(result.Results).Is("[123]")
	}
	{
		result, err := execTool("return 1, 'string'")
		assert.That(err == nil)
		assert.That(len(result.Warnings) == 0)
		assert.This(result.Results).Is(`[1, "string"]`)
	}
	{
		result, err := execTool("x = 1; y = 2")
		assert.That(err == nil)
		assert.This(result.Results).Is("[2]")
		assert.That(strings.Contains(result.Warnings[0], "initialized but not used: x"))
		assert.That(strings.Contains(result.Warnings[0], "@line:1"))
		assert.That(strings.Contains(result.Warnings[1], "initialized but not used: y"))
		assert.That(strings.Contains(result.Warnings[1], "@line:1"))
	}
	{
		_, err := execTool("throw 'exception'")
		assert.That(strings.Contains(err.Error(), `execute error: "exception"`))
		assert.That(strings.Contains(err.Error(), "\n"))
	}
	{
		_, err := execTool("x")
		assert.That(strings.Contains(err.Error(), "execute error: uninitialized variable: x"))
		assert.That(strings.Contains(err.Error(), "\n"))
	}
	{
		_, err := execTool("if true") // syntax error - missing block
		assert.That(err != nil)
		assert.That(strings.Contains(err.Error(), "execute error: syntax error"))
		assert.That(strings.Contains(err.Error(), "@line:2")) // error at closing }
	}
	{
		result, err := execTool("x = 1\ny = 2")
		assert.That(err == nil)
		assert.That(strings.Contains(result.Warnings[0], "@line:1"))
		assert.That(strings.Contains(result.Warnings[1], "@line:2"))
	}
	// Set Print to a function that dispatches to Suneido.Print,
	// matching the pattern of stdlib Print.ss.
	oldPrint := core.Global.GetIfPresent("Print")
	defer core.Global.TestDef("Print", oldPrint)
	core.Global.TestDef("Print", &core.SuBuiltinRaw{
		Fn: func(th *core.Thread, as *core.ArgSpec, args []core.Value) core.Value {
			suneido := core.Global.GetName(th, "Suneido")
			if suneido == nil {
				return nil
			}
			printFn := suneido.Get(th, core.SuStr("Print"))
			if printFn == nil {
				return nil
			}
			th.Call(printFn, args[0])
			return nil
		},
		BuiltinParams: core.BuiltinParams{ParamSpec: core.ParamSpec1},
	})
	{
		result, err := execTool("Print('hello')")
		assert.That(err == nil)
		assert.This(result.Print).Is("hello")
	}
	{
		result, err := execTool("Print('a')\nPrint('b')")
		assert.That(err == nil)
		assert.This(result.Print).Is("ab")
	}
}

func TestCheckTool(t *testing.T) {
	assert := assert.T(t)
	{
		result, err := checkTool("1 + 2 \n")
		assert.That(err == nil)
		assert.That(len(result.Warnings) == 0)
	}
	{
		result, err := checkTool("x = 1; y = 2")
		assert.That(err == nil)
		assert.That(strings.Contains(result.Warnings[0], "initialized but not used: x"))
		assert.That(strings.Contains(result.Warnings[0], "@line:1"))
		assert.That(strings.Contains(result.Warnings[1], "initialized but not used: y"))
		assert.That(strings.Contains(result.Warnings[1], "@line:1"))
	}
	{
		_, err := checkTool("throw 'exception'")
		assert.That(err == nil) // checkTool should not throw errors for exceptions in code
	}
	{
		_, err := checkTool("x")
		assert.That(err == nil) // checkTool should not throw errors for uninitialized variables
	}
	{
		_, err := checkTool("if true") // syntax error - missing block
		assert.That(err != nil)
		assert.That(strings.Contains(err.Error(), "check error: syntax error"))
		assert.That(strings.Contains(err.Error(), "@line:2")) // error at closing }
	}
	{
		result, err := checkTool("x = 1\ny = 2")
		assert.That(err == nil)
		assert.That(strings.Contains(result.Warnings[0], "@line:1"))
		assert.That(strings.Contains(result.Warnings[1], "@line:2"))
	}
}
