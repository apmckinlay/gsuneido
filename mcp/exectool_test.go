// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mcp

import (
	"strings"
	"testing"

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
		assert.That(strings.Contains(result.Warnings[1], "initialized but not used: y"))
	}
	{
		_, err := execTool("throw 'exception'")
		assert.This(err.Error()).Is(`execute error: "exception"`)
	}
	{
		_, err := execTool("x")
		assert.This(err.Error()).Is("execute error: uninitialized variable: x")
	}
}
