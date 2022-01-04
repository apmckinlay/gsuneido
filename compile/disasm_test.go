// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestDisasm(t *testing.T) {
	test := func(src, expected string) {
		t.Helper()
		ast := parseFunction(src)
		fn := codegen("", "", ast, nil).(*SuFunc)
		s := DisasmMixed(fn, src)
		assert.T(t).This(s).Like(expected)
	}

	test(`function () {
		return function() { f() }
		}`,
		`16: return function() {
				0: Value /* function */
        >    36: f() }
		>            0: Load f
		>            2: CallFuncNilOk ()`)

	test(`function () {
		return function() { f(); g() }
		}`,
		`16: return function() {
                    0: Value /* function */
        >    36: f();
        >             0: Load f
        >             2: CallFuncDiscard ()
        >    41: g() }
        >             4: Load g
        >             6: CallFuncNilOk ()`)

	test(`function () {
		fn = function(){}
		fn()
		}`,
		`16: fn = function(){}
		36: fn()
				0: Value fn /* function */
				2: CallFuncNilOk ()`)

	test(`function(){}`,
		``)

	test(`function(){ f() }`,
		`12: f()
				0: Load f
				2: CallFuncNilOk ()`)

	test(`function(){ f(); g() }`,
		`12: f();
				0: Load f
				2: CallFuncDiscard ()
		17: g()
				4: Load g
				6: CallFuncNilOk ()`)

	test(`function (x, y) {
		a = x
		b = y
		return a + b
		}`,
		`20: a = x
			    0: Load x
			    2: Store a
			    4: Pop
	    28: b = y
			    5: Load y
			    7: Store b
			    9: Pop
	    36: return a + b
			   10: Load a
			   12: Load b
			   14: Add`)

	test(`function (x) {
		if x
			{
			return 0
			}
		else
			return 1
		}`,
		`17: if x
				0: Load x
				2: JumpFalse 10
		25: {
		30: return 0
			}
			else
				5: Zero
				6: Return
				7: Jump 12
		54: return 1
				10: One
				11: Return`)

	test(`function (x) {
		try
			return x
		catch(e)
			return e
		}`,
		`17: try
			    0: Try 10 ""
	    24: return x
			    4: Load x
			    6: Return
			    7: Catch 16
	    35: catch(e)
			   10: Store e
			   12: Pop
	    47: return e
			   13: Load e
			   15: Return`)
}
