// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestParseFunction(t *testing.T) {
	fn, ok := parse("", "function (a, b) { return a + b }").(*ast.Function)
	assert.T(t).That(ok)
	assert.T(t).This(len(fn.Params)).Is(2)
}

func TestParseClass(t *testing.T) {
	val := parse("", "Foo { New() { } }")
	assert.T(t).This(core.StrOrEmpty(val.Get(nil, core.StrVal("type")))).Is("Class")
}

func TestParseLiteral(t *testing.T) {
	assert.T(t).This(parse("", "123").String()).Is("123")
}

func TestParseError(t *testing.T) {
	saved := os.Stderr
	r, w, err := os.Pipe()
	assert.T(t).That(err == nil)
	os.Stderr = w
	t.Cleanup(func() { os.Stderr = saved })

	val := parse("Bad", "function (")
	assert.T(t).That(val == nil)

	w.Close()
	out, err := io.ReadAll(r)
	assert.T(t).That(err == nil)
	assert.T(t).That(strings.Contains(string(out), `"Bad"`))
}
