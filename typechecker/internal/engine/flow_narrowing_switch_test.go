// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// reaching default proves x isnt 1, which says nothing about x not being a
// Number, so the Number arm has to survive - narrowing it away turned `x + 1`
// into an error on code that runs
func TestSwitchDefaultKeepsValueTypedArms(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: number | string) {
			switch (x) {
			case 1:
				return 0
			default:
				return x + 1
			}
		}
	}`, "T")
	a.That(hasDiag(env, "Foo", SeverityError, "got number|string"))
	a.That(!hasDiag(env, "Foo", SeverityError, "got string"))
}

// a boolean case is subtractable - false is the only value of its type - so
// reaching default really does leave only the Number arm
func TestSwitchDefaultRemovesBooleanArm(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: number | false) {
			switch (x) {
			case false:
				return 0
			default:
				return x + 1
			}
		}
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}

// `switch (true)` is the same cond-mode dispatch as bare `switch` - the
// parens must not push it into scrutinee mode, where every case body would
// see the unnarrowed type
func TestSwitchParenTrueCondMode(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x = false) {
			switch (true) {
			case Number?(x):
				return x + 1
			}
			return 0
		}
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}
