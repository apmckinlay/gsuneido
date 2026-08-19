// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"testing"

	"pgregory.net/rapid"

	"github.com/apmckinlay/gsuneido/typechecker/internal/synth"
)

// demotion, checked against the generator's ground truth: an optional member
// (seeded `: false`, assigned only in New) must demote to its clean payload
// exactly when New assigns it on every path, and must never demote otherwise.

func TestPropConstructorDemotion(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		p := synth.GenProgram(rt, synth.Config{WellTyped: true})
		hasNew := false
		for _, m := range p.Methods {
			if m.Name == "New" {
				hasNew = true
			}
		}
		if !hasNew {
			rt.Skip("no constructor generated")
		}
		src := synth.Render(&p)
		cls, ok := safeParse(src, "T")
		if !ok {
			rt.Skip("unparseable render")
		}
		env, panicMsg := safeRun(cls)
		if panicMsg != "" {
			rt.Fatalf("pipeline panicked: %s\nsrc:\n%s", panicMsg, src)
		}
		if d, bad := anySeverityError(env); bad {
			rt.Skip("base wellTyped program errors (covered by TestPropWellTypedNoError): " + d.Msg)
		}
		for _, mem := range p.Members {
			switch mem.Ctor {
			case synth.CtorAlways:
				got, ok := env.PostCtorMembers[mem.Name]
				if !ok {
					rt.Fatalf("always-assigned %q not demoted (no PostCtorMembers entry)\nsrc:\n%s",
						mem.Name, src)
				}
				if isDirty(got) || containsFalse(got) {
					rt.Fatalf("demoted %q not clean non-false: %v\nsrc:\n%s", mem.Name, got, src)
				}
				if !armsFit(got, gtypeDynType(mem.Payload)) {
					rt.Fatalf("demoted %q=%v outside payload %v\nsrc:\n%s",
						mem.Name, got, gtypeDynType(mem.Payload), src)
				}
				// the seed view must keep the False arm
				if !containsFalse(env.Members[mem.Name]) {
					rt.Fatalf("seed view of %q lost its False arm: %v\nsrc:\n%s",
						mem.Name, env.Members[mem.Name], src)
				}
			case synth.CtorMaybe, synth.CtorNever:
				if got, ok := env.PostCtorMembers[mem.Name]; ok {
					rt.Fatalf("unsound demotion of conditionally-assigned %q to %v\nsrc:\n%s",
						mem.Name, got, src)
				}
			}
		}
	})
}

// pins the demotion contract on hand-checkable sources, including the traps
// the generator is built around (early return, false put back by a method).
func TestDemotionKnownCases(t *testing.T) {
	cases := []struct {
		label, src string
		wantDemote bool
	}{
		{"always assigned", `class { x: false
New()
{
.x = 5
}
Get()
{
return .x
} }`, true},
		{"conditional assign (implicit else keeps seed)", `class { x: false
New(a)
{
if (a)
{
.x = 5
}
} }`, false},
		{"never assigned", `class { x: false
New()
{
v0 = 1
} }`, false},
		{"method puts false back (C2)", `class { x: false
New()
{
.x = 5
}
Reset()
{
.x = false
} }`, false},
		{"early return before the assignment (C1)", `class { x: false
New(a)
{
if (a)
{
return
}
.x = 5
} }`, false},
		{"provably-true condition is effectively unconditional", `class { x: false
New()
{
if ((false isnt true))
{
.x = 5
}
} }`, true},
		{"assignment from a dirty member (might be false, C2)", `class { x: false
f: "a"
New()
{
.x = .f
}
M0(p0)
{
if (String?(p0))
{
.f = p0
}
return 0
} }`, false},
	}
	for _, c := range cases {
		cls, ok := safeParse(c.src, "T")
		if !ok {
			t.Fatalf("%s: did not parse:\n%s", c.label, c.src)
		}
		env, panicMsg := safeRun(cls)
		if panicMsg != "" {
			t.Fatalf("%s: pipeline panicked: %s", c.label, panicMsg)
		}
		got, demoted := env.PostCtorMembers["x"]
		if demoted != c.wantDemote {
			t.Errorf("%s: demoted=%v want %v (PostCtorMembers[x]=%v, Members[x]=%v)",
				c.label, demoted, c.wantDemote, got, env.Members["x"])
		}
		if c.wantDemote && demoted {
			if isDirty(got) || containsFalse(got) || !armsFit(got, TNumber) {
				t.Errorf("%s: demoted type %v is not clean Number", c.label, got)
			}
		}
	}
	// demotion restamps reads: Get returns the demoted type, not False|Number
	_, env := runPasses(`class { x: false
New()
{
.x = 5
}
Get()
{
return .x
} }`, "T")
	if got := env.Returns["Get"]; containsFalse(got) {
		t.Errorf("Returns[Get]=%v still carries the False arm after demotion", got)
	}
}
