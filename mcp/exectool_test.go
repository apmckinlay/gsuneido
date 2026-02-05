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
		result, err := exectool("1 + 2 \n")
		assert.That(err == nil)
		assert.That(strings.Contains(result, "warnings: []"))
		assert.That(strings.Contains(result, "results: [3]"))
	}
	{
		result, err := exectool("return")
		assert.That(err == nil)
		assert.That(strings.Contains(result, "warnings: []"))
		assert.That(strings.Contains(result, "results: []"))
	}
	{
		result, err := exectool("return 123")
		assert.That(err == nil)
		assert.That(strings.Contains(result, "warnings: []"))
		assert.That(strings.Contains(result, "results: [123]"))
	}
	{
		result, err := exectool("return 1, 'string'")
		assert.That(err == nil)
		assert.That(strings.Contains(result, "warnings: []"))
		assert.That(strings.Contains(result, `results: [1, "string"]`))
	}
	{
		result, err := exectool("x = 1; y = 2")
		assert.That(err == nil)
		assert.That(strings.Contains(result, "results: [2]"))
		assert.That(strings.Contains(result, "WARNING: initialized but not used: x @14"))
		assert.That(strings.Contains(result, "WARNING: initialized but not used: y @21"))
	}
	{
		_, err := exectool("throw 'exception'")
		assert.This(err.Error()).Is(`execute error: "exception"`)
	}
	{
		_, err := exectool("x")
		assert.This(err.Error()).Is("execute error: uninitialized variable: x")
	}
}
