// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"
	"strings"
	"testing"

	"pgregory.net/rapid"

	"github.com/apmckinlay/gsuneido/typechecker/internal/synth"
)

// an error needs clean evidence: at least one non-? arm in Got.
// evidence that is nothing but a guess (pure ? / unknown) caps at warning.
func TestPropErrorEvidenceHasCleanArm(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		src := genSource(rt, synth.Config{})
		_, env := runGenerated(rt, src)
		for _, d := range diagList(env) {
			if d.Severity != SeverityError || len(d.Got) == 0 {
				continue
			}
			clean := false
			for _, ty := range d.Got {
				if members, _ := decomposeForCheck(ty); len(members) > 0 {
					clean = true
				}
			}
			if !clean {
				rt.Fatalf("error built on purely dirty evidence %v: %+v\nsrc:\n%s", d.Got, d, src)
			}
		}
	})
}

var guessMethodNames = []string{"Size", "Substr", "AddDays", "Format", "Sort!"}

var guessArgLits = []string{"1", "42", `"a"`, "true", "#20200101"}

// builds a class whose methods call builtin names on unknown-typed params,
// so CallsiteResolutionPass takes the guessing path.
func genGuessSource(rt *rapid.T) string {
	var b strings.Builder
	b.WriteString("class {\n")
	nmeth := rapid.IntRange(1, 3).Draw(rt, "nmeth")
	for i := range nmeth {
		fmt.Fprintf(&b, "M%d(p0, p1) {\n", i)
		nstmts := rapid.IntRange(1, 4).Draw(rt, "nstmts")
		for range nstmts {
			call := genGuessCall(rt)
			switch rapid.IntRange(0, 6).Draw(rt, "stmtkind") {
			case 0:
				fmt.Fprintf(&b, "v = %s\nreturn v\n", call)
			case 1:
				op := rapid.SampledFrom([]string{"+", "$", "and"}).Draw(rt, "op")
				lit := rapid.SampledFrom(guessArgLits).Draw(rt, "oplit")
				fmt.Fprintf(&b, "return (%s %s %s)\n", call, op, lit)
			case 2:
				fmt.Fprintf(&b, "if (%s)\n{ return 1 }\n", call)
			case 3:
				name := rapid.SampledFrom(guessMethodNames).Draw(rt, "chain")
				fmt.Fprintf(&b, "return %s.%s()\n", call, name)
			case 4:
				op := rapid.SampledFrom([]string{"+", "$", "and"}).Draw(rt, "vop")
				lit := rapid.SampledFrom(guessArgLits).Draw(rt, "vlit")
				fmt.Fprintf(&b, "v = %s\nreturn (v %s %s)\n", call, op, lit)
			case 5:
				fmt.Fprintf(&b, "v = %s\nw = v\nif (w)\n{ return 1 }\n", call)
			default:
				fmt.Fprintf(&b, "%s\n", call)
			}
		}
		b.WriteString("return 0\n}\n")
	}
	b.WriteString("}")
	return b.String()
}

func genGuessCall(rt *rapid.T) string {
	recv := rapid.SampledFrom([]string{"p0", "p1"}).Draw(rt, "recv")
	name := rapid.SampledFrom(guessMethodNames).Draw(rt, "callee")
	nargs := rapid.IntRange(0, 5).Draw(rt, "nargs")
	args := make([]string, nargs)
	for i := range args {
		args[i] = rapid.SampledFrom(guessArgLits).Draw(rt, "arg")
	}
	return fmt.Sprintf("%s.%s(%s)", recv, name, strings.Join(args, ", "))
}

// checks never error on a guess: not at the call node itself, and not
// where the guess flows through locals or chained receivers. every receiver
// this generator produces is guess-derived, so a "not applicable" error here
// means provenance leaked.
func TestPropGuessesNeverError(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		src := genGuessSource(rt)
		_, env := runGenerated(rt, src)
		guessPos := make(map[int]bool, len(env.GuessedCalls))
		for n := range env.GuessedCalls {
			guessPos[nodePos(n)] = true
		}
		for _, d := range diagList(env) {
			if d.Severity != SeverityError {
				continue
			}
			if guessPos[d.Pos] {
				rt.Fatalf("error at guessed call pos %d: %+v\nsrc:\n%s", d.Pos, d, src)
			}
			if strings.Contains(d.Msg, "not applicable") {
				rt.Fatalf("receiver error on guess-derived receiver: %+v\nsrc:\n%s", d, src)
			}
		}
	})
}

// a guess caps findings at warning through chained receivers, tainted locals, and aliases
func TestGuessProvenanceLeaks(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "number", Name: "Round", Sig: "(number) :number"},
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
	)
	cases := []struct{ label, src, want string }{
		{"chained receiver", `class { M(p) { return p.Size().Size() } }`, "not applicable"},
		{"var operand", "class { M(p) { v = p.Size()\nreturn v $ \"x\" } }", `operator "$" expects`},
		{"alias condition", "class { M(p) { v = p.Size()\nw = v\nif (w) { return 1 }\nreturn 0 } }", "condition expects boolean"},
		{"var receiver arity", "class { M(p) { v = p.Size()\nreturn v.Round(1, 2, 3) } }", "too many arguments"},
	}
	for _, c := range cases {
		cls, ok := safeParse(c.src, "T")
		if !ok {
			t.Fatalf("%s: unparseable", c.label)
		}
		env, panicMsg := safeRun(cls)
		if panicMsg != "" {
			t.Fatalf("%s: %s", c.label, panicMsg)
		}
		found := false
		for _, d := range diagList(env) {
			if d.Severity == SeverityError {
				t.Errorf("%s: error on a guess: %+v", c.label, d)
			}
			if d.Severity == SeverityWarning && strings.Contains(d.Msg, c.want) {
				found = true
			}
		}
		if !found {
			t.Errorf("%s: expected a downgraded warning containing %q", c.label, c.want)
		}
	}
}
