// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"strings"
	"testing"
)

// a trailing expression statement is the implicit return (codegen exprStmt
// lastStmt keeps its value): the method must not publish Void, and calling a
// method on its result must not error as a Void receiver.
func TestImplicitReturn_TrailingExpr(t *testing.T) {
	src := `class {
	Lit() { 5 }
	Wrap(q) { Transaction(read:) {|t| t.QueryMax(q) } }
	Use() { m = .Wrap(1)
		return m.Extract("[0-9]+")
		}
	Tail() { if true { 5 } }
	}`
	_, env := runPasses(src, "T")
	if got := env.Returns["Lit"]; got != TNumber {
		t.Errorf("Lit returns %v, want Number", got)
	}
	if got := env.Returns["Wrap"]; got == TVoid {
		t.Error("Wrap published Void; trailing expression is the implicit return")
	}
	if got := env.Returns["Tail"]; got != TVoid {
		t.Errorf("Tail returns %v, want Void (control-flow tail returns nothing)", got)
	}
	for _, d := range diagList(env) {
		if d.Severity == SeverityError && strings.Contains(d.Msg, "receiver of type void") {
			t.Fatalf("Void-receiver false positive survived: %s", d.Msg)
		}
	}
}
