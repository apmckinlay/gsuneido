// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"strings"
	"testing"
)

// conflict-cascade confluence: B.y takes String from Scanner in round 1 but a
// Number demand chains in later (B -> D -> E -> annotated). whichever round A
// reads B.y in, the final findings must be identical for every ref order.
func TestRequirement_ConflictCascadeConfluent(t *testing.T) {
	srcs := map[string]string{
		"AA": `class {
	CallClass(x) { return BB(x) }
	}`,
		"BB": `class {
	CallClass(y) { Scanner(y)
		return DD(y)
		}
	}`,
		"DD": `class {
	CallClass(w) { return EE(w) }
	}`,
		"EE": `class {
	NumOnly(v: number) { return v }
	CallClass(v) { return .NumOnly(v) }
	}`,
	}
	driver := `class {
	Use() { return AA(0) }
	}`
	findings := func(order []string) []string {
		refs := make([]RefSource, len(order))
		for i, n := range order {
			refs[i] = RefSource{Name: n, Src: srcs[n]}
		}
		regs := BuildReferenceRegistry(refs)
		_, env := runPassesWith(driver, "T", func(e *TypeEnv) { regs.Seed(e) })
		var out []string
		for _, d := range diagList(env) {
			if strings.Contains(d.Msg, "is passed on to") {
				out = append(out, d.Msg)
			}
		}
		return out
	}
	a := findings([]string{"BB", "DD", "EE", "AA"})
	b := findings([]string{"AA", "EE", "DD", "BB"})
	if len(a) != len(b) {
		t.Fatalf("ref order changed findings:\n order1 (%d): %v\n order2 (%d): %v", len(a), a, len(b), b)
	}
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("ref order changed finding %d:\n %s\n %s", i, a[i], b[i])
		}
	}
	// B.y ends CONFLICT (String from Scanner, Number chained via DD -> EE), so
	// it is genuinely polymorphic and nothing upstream may keep a requirement
	// derived from its transient String phase
	if len(a) != 0 {
		t.Fatalf("requirement survived its retracted justification: %v", a)
	}
}
