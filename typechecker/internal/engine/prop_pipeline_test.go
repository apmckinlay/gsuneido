// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"
	"testing"

	"pgregory.net/rapid"

	"github.com/apmckinlay/gsuneido/typechecker/internal/synth"
)

// safeParse recovers ParseClass's designed panic. ok=false means unparseable.
func safeParse(src, name string) (cls *ClassObject, ok bool) {
	defer func() {
		if recover() != nil {
			cls, ok = nil, false
		}
	}()
	return NewClassObject(name, ParseClass(src)), true
}

// safeRun runs the pipeline, capturing any panic as a non-empty message.
func safeRun(cls *ClassObject) (env TypeEnv, panicMsg string) {
	defer func() {
		if r := recover(); r != nil {
			panicMsg = fmt.Sprint(r)
		}
	}()
	env = NewTypeEnv()
	RunPipeline(cls, env, NewPassCtx(), nil)
	return env, ""
}

func genSource(t *rapid.T, cfg synth.Config) string {
	p := synth.GenProgram(t, cfg)
	return synth.Render(&p)
}

func classMethodSet(cls *ClassObject) map[string]bool {
	s := make(map[string]bool, len(cls.Methods))
	for name := range cls.Methods {
		s[name] = true
	}
	return s
}

func runGenerated(rt *rapid.T, src string) (*ClassObject, TypeEnv) {
	cls, ok := safeParse(src, "T")
	if !ok {
		rt.Skip("unparseable render (parser constant-folder rejected)")
	}
	env, panicMsg := safeRun(cls)
	if panicMsg != "" {
		rt.Fatalf("pipeline panicked: %s\nsrc:\n%s", panicMsg, src)
	}
	return cls, env
}

func assertEnvSemEq(rt *rapid.T, label, src string, a, b TypeEnv) {
	for _, m := range []struct {
		Name string
		x, y map[string]DynType
	}{
		{"Members", a.Members, b.Members},
		{"PostCtorMembers", a.PostCtorMembers, b.PostCtorMembers},
		{"Returns", a.Returns, b.Returns},
		{"PreCtorReturns", a.PreCtorReturns, b.PreCtorReturns},
	} {
		if !typeMapsSemEq(m.x, m.y) {
			rt.Fatalf("%s: %s differ\n a=%v\n b=%v\nsrc:\n%s", label, m.Name, m.x, m.y, src)
		}
	}
	ka := sortedDiagKeys(diagList(a), diagKey)
	kb := sortedDiagKeys(diagList(b), diagKey)
	if !strSlicesEqual(ka, kb) {
		rt.Fatalf("%s: Diagnostics differ\n a=%v\n b=%v\nsrc:\n%s", label, ka, kb, src)
	}
}

func TestPropPipelineWellFormed(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		src := genSource(rt, synth.Config{})
		cls, env := runGenerated(rt, src)

		for _, published := range []map[string]DynType{
			env.Members, env.PostCtorMembers, env.Returns, env.PreCtorReturns,
		} {
			for _, ty := range published {
				assertNormalForm(rt, ty)
			}
		}

		methods := classMethodSet(cls)
		for _, d := range diagList(env) {
			if d.Pos < 0 || d.Pos > len(src) {
				rt.Fatalf("diagnostic Pos %d out of [0,%d]: %+v\nsrc:\n%s", d.Pos, len(src), d, src)
			}
			if d.Msg == "" {
				rt.Fatalf("empty diagnostic Msg: %+v\nsrc:\n%s", d, src)
			}
			if d.Method != "" && !methods[d.Method] {
				rt.Fatalf("diagnostic method %q not on class\nsrc:\n%s", d.Method, src)
			}
		}

		// view coherence
		for k, pv := range env.PostCtorMembers {
			mv, ok := env.Members[k]
			if !ok {
				rt.Fatalf("PostCtorMembers key %q not in Members\nsrc:\n%s", k, src)
			}
			if !fits(pv, mv) {
				rt.Fatalf("PostCtorMembers[%q]=%v does not fit Members[%q]=%v\nsrc:\n%s", k, pv, k, mv, src)
			}
		}
		if !typeMapKeysSubset(env.PreCtorReturns, env.Returns) {
			rt.Fatalf("PreCtorReturns keys not subset of Returns\n pre=%v\n ret=%v\nsrc:\n%s",
				env.PreCtorReturns, env.Returns, src)
		}
	})
}

// TestPropPipelineDeterministic - two fresh parse+run cycles agree semantically.
func TestPropPipelineDeterministic(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		src := genSource(rt, synth.Config{})
		_, e1 := runGenerated(rt, src)
		_, e2 := runGenerated(rt, src)
		assertEnvSemEq(rt, "determinism", src, e1, e2)
	})
}

func TestPropWellTypedDeterministic(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		src := genSource(rt, synth.Config{WellTyped: true})
		_, e1 := runGenerated(rt, src)
		_, e2 := runGenerated(rt, src)
		assertEnvSemEq(rt, "wellTyped-determinism", src, e1, e2)
	})
}

func TestPropPipelineDeterminismAnyProgram(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		src := genSource(rt, synth.Config{})
		_, e1 := runGenerated(rt, src)
		_, e2 := runGenerated(rt, src)
		assertEnvSemEq(rt, "anyProgram-determinism", src, e1, e2)
	})
}

func TestPipelineLoneFunctionNoPanic(t *testing.T) {
	for _, src := range []string{
		`function (x) { return x + 1 }`,
		`function () { return 1 }`,
		`function (x, y = 2) { return x $ y }`,
	} {
		cls, ok := safeParse(src, "F")
		if !ok {
			t.Fatalf("lone function did not parse: %q", src)
		}
		if _, panicMsg := safeRun(cls); panicMsg != "" {
			t.Fatalf("pipeline panicked on %q: %s", src, panicMsg)
		}
	}
}

func TestPipelineObjectLiteralNoPanic(t *testing.T) {
	for _, src := range []string{`#()`, `#(1, 2, 3)`, `#(a: 1, b: 2)`} {
		cls, ok := safeParse(src, "O")
		if !ok {
			continue
		}
		if _, panicMsg := safeRun(cls); panicMsg != "" {
			t.Fatalf("pipeline panicked on object literal %q: %s", src, panicMsg)
		}
	}
}

func FuzzPipeline(f *testing.F) {
	for _, s := range []string{
		`class { }`,
		`class { F() { return 1 } }`,
		`class { a: 1  F(x) { return (x + a) } }`,
		`class { f0: false  M0(p0) { if (p0) { return 1 } return "x" } }`,
		`function (x) { return x }`,
		`#(1, 2, 3)`,
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, src string) {
		cls, ok := safeParse(src, "Fuzz")
		if !ok {
			return // ParseClass's designed panic, recovered
		}
		if _, panicMsg := safeRun(cls); panicMsg != "" {
			t.Fatalf("pipeline panicked on %q: %s", src, panicMsg)
		}
	})
}
