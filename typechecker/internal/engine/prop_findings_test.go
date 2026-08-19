// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"
	"strings"
	"testing"
)

func TestDuplicateSiblingDiagnostics(t *testing.T) {
	src := `class { M0(p0) { .M1((p0.Number?()), (p0.Number?())) } M1(a, b) { } }`
	env := mustRun(t, src)
	counts := map[string]int{}
	for _, d := range diagList(env) {
		counts[diagKey(d)]++
	}
	maxDup := 0
	for _, c := range counts {
		if c > maxDup {
			maxDup = c
		}
	}
	if maxDup < 2 {
		t.Fatalf("expected duplicate sibling diagnostics (current behavior); got max multiplicity %d - "+
			"if the checker now de-duplicates, update TestPropPipelineWellFormed to assert uniqueness", maxDup)
	}
}

// found by TestPropWellTypedPublishedTypes: a member assigned a union with an
// unresolved-callsite arm on fixpoint iteration 1 kept the transient dirty
// flag forever, because member unions accumulate monotonically. the fix makes
// collectThisAssignments union only proven arms and MemberDirtyPass re-dirty
// any RHS whose converged type still decomposes dirty.
func TestMemberAssignTransientDirtNotPublished(t *testing.T) {
	cases := []struct {
		label, src string
		wantDirty  bool
	}{
		{"direct call", "class { f0: 1\nM0()\n{\nreturn 1\n}\nM2()\n{\n.f0 = .M0()\nreturn 0\n} }", false},
		{"ternary-wrapped call", "class { f0: 1\nM0()\n{\nreturn 1\n}\nM2()\n{\n.f0 = (false ? 1 : .M0())\nreturn 0\n} }", false},
		{"genuinely unresolvable arm", "class { f0: 1\nM2(p0)\n{\n.f0 = (Number?(p0) ? 1 : Zzz.Foo())\nreturn 0\n} }", true},
	}
	for _, c := range cases {
		_, env := runPasses(c.src, "T")
		got := env.Members["f0"]
		_, dirty := decomposeForCheck(got)
		if dirty != c.wantDirty {
			t.Errorf("%s: Members[f0]=%v dirty=%v, want dirty=%v", c.label, got, dirty, c.wantDirty)
		}
	}
}

// found by TestPropDirtModel: a local assigned from a ternary whose arm reads
// a guarded param kept its stale pre-narrowing dirty flow stamp, because
// isNarrower had no clean-over-dirty rule for the primitive-vs-dirty-union
// case, so the assign refinement never rescued the read.
func TestAssignRefinementRescuesDirtyFlow(t *testing.T) {
	cases := []struct{ label, src string }{
		{"direct", "class { M0(p0)\n{\nif (Boolean?(p0))\n{\nv3 = p0\nreturn v3\n}\nreturn false\n} }"},
		{"ternary arm", "class { f0: 1\nM0(p0)\n{\nif (Boolean?(p0))\n{\nif (p0)\n{\nv3 = ((.f0 < 1) ? false : p0)\nreturn v3\n}\n}\nreturn false\n} }"},
	}
	for _, c := range cases {
		_, env := runPasses(c.src, "T")
		got := env.Returns["M0"]
		if _, dirty := decomposeForCheck(got); dirty {
			t.Errorf("%s: Returns[M0]=%v is dirty; the guarded assignment should have cleaned it", c.label, got)
		}
	}
}

func TestFixpointReportsRecursiveReturnError(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Number?", Sig: "() :boolean"},
		Reg{Kind: "free", Name: "Number?", Sig: "(value) :boolean"},
	)
	src := `class
	{
	M0(p0, p1 = true)
	{
		v0 = (1 ? .M0(p0, 1, p0) : (p0.Number?()))
		while (v0)
		{
			if (1)
			{
				v2 = ((v0.Number?()) + 1)
				v1 = (1 + 1)
				v1 = 1
			}
			else
			{
				v1 = (not v0)
				.M0(v0, v0, .M0(1))
				return v1
			}
		}
	}
	}
`
	env := mustRun(t, src)
	found := false
	for _, d := range diagList(env) {
		if strings.Contains(d.Msg, "not applicable to receiver of type boolean") {
			found = true
		}
	}
	if !found {
		t.Fatalf("the recursive-return error is missing (got %d diagnostics) - "+
			"fixpoint detection may have regressed", len(diagList(env)))
	}
}

func TestPipelineDeterministicOnRecursiveMembers(t *testing.T) {
	src := `class
	{
	f0: #(1, 2, 3)
	f1: #(1, 2, 3)
	M0(p0)
	{
		if (Object?(p0))
		{
			return .f0
			.f1 = p0
		}
		if (Boolean?(p0))
		{
			.f1 = .f0
			v2 = "x0"
		}
		return .M1()
		return .f0
	}
	M1()
	{
		return .f1
		while (true)
		{
			.f0 = .M1()
			return .f1
			.f1 = (true ? .M1() : #(1, 2, 3))
		}
		return .M0((3 - 85))
	}
	M2()
	{
		return (.f1 isnt #(1, 2, 3))
		.f1 = .M0((13 * 1))
		return .M2()
	}
	}
`
	first := ""
	for i := range 100 {
		env := mustRun(t, src)
		got := fmt.Sprintf("f0=%v|f1=%v|M0=%v|M1=%v|M2=%v",
			env.Members["f0"], env.Members["f1"],
			env.Returns["M0"], env.Returns["M1"], env.Returns["M2"])
		if i == 0 {
			first = got
		} else if got != first {
			t.Fatalf("nondeterministic pipeline output:\n run 0: %s\n run %d: %s", first, i, got)
		}
	}
}
