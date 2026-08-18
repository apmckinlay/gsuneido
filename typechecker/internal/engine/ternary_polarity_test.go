// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"strings"
	"testing"
)

// `x is false ? A : x` - the else arm runs only when x isnt false, so the
// member (or local) it feeds can never hold the False arm. the receiver
// checks must not claim a False path at 0.90 on provably-correct code.
// ConstructorExecPass already applied this rule (evalTrinary); the generic
// trinary stampers must agree with it.
func TestTernaryFalseGuardPolarity(t *testing.T) {
	cases := []struct{ label, src string }{
		{"ctor member", `class {
	New(children = false)
		{
		.children = children is false ? Object() : children
		}
	AddChild(c)
		{
		.children.Add(c)
		}
	}`},
		{"non-ctor member", `class {
	Set(children = false)
		{
		.children = children is false ? Object() : children
		}
	AddChild(c)
		{
		.children.Add(c)
		}
	}`},
		{"isnt flip", `class {
	Set(children = false)
		{
		.children = children isnt false ? children : Object()
		}
	AddChild(c)
		{
		.children.Add(c)
		}
	}`},
	}
	for _, c := range cases {
		_, env := runPasses(c.src, "T")
		for _, d := range diagList(env) {
			if d.Severity == SeverityError &&
				strings.Contains(d.Msg, "not applicable") {
				t.Errorf("%s: false-arm survived the guard: %s", c.label, d.Msg)
			}
		}
	}
}
