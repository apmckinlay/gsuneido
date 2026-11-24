// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestUninit(t *testing.T) {
	test := func(body string, want Value) {
		t.Helper()
		th := &Thread{}
		assert.T(t).This(compile.EvalString(th, body)).Is(want)
	}
	test("return Uninit?('UninitBuiltinMissing')", True)
	test("return Uninit?('Uninit?')", False)
}
