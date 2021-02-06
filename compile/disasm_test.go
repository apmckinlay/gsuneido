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
		ast := parseFunction(src)
		fn := codegen("", "", ast)
		s := DisasmMixed(fn, src)
		assert.T(t).This(s).Like(expected)
	}

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
