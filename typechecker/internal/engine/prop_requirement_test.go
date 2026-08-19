// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"pgregory.net/rapid"

	"github.com/apmckinlay/gsuneido/typechecker/internal/synth"
)

// requirement chains correct by construction: W1 -> W2 -> ... -> WD -> Scanner,
// optionally one level also demanding Number (a conflict). the generator knows
// the ground truth: a bad driver literal yields exactly one finding iff no
// conflict sits on the path, and findings are identical under every reference
// order (the restart-on-conflict fixpoint's confluence claim).
func TestPropRequirementChainOracle(t *testing.T) {
	withSigs(t,
		Reg{Kind: "free", Name: "Scanner", Sig: "(string :string) :unknown"},
	)
	rapid.Check(t, func(rt *rapid.T) {
		depth := rapid.IntRange(1, 4).Draw(rt, "depth")
		conflictAt := rapid.IntRange(0, depth).Draw(rt, "conflictAt") // 0 = none
		bad := rapid.Bool().Draw(rt, "bad")

		srcs := map[string]string{}
		names := []string{}
		for k := 1; k <= depth; k++ {
			name := fmt.Sprintf("W%d", k)
			next := fmt.Sprintf("W%d(p)", k+1)
			if k == depth {
				next = "Scanner(p)"
			}
			extra := ""
			if k == conflictAt {
				extra = "NumOnly(v: number) { return v }\n\t"
				next = ".NumOnly(p)\n\t\treturn " + next
			}
			srcs[name] = "class {\n\t" + extra + "CallClass(p) { return " + next + " }\n\t}"
			names = append(names, name)
		}
		lit := `"s"`
		if bad {
			lit = "0"
		}
		driver := "class {\n\tUse() { return W1(" + lit + ") }\n\t}"

		run := func(order []string) []string {
			refs := make([]RefSource, len(order))
			for i, n := range order {
				refs[i] = RefSource{Name: n, Src: srcs[n]}
			}
			regs := BuildReferenceRegistry(refs)
			_, env := runPassesWith(driver, "T", func(e *TypeEnv) { regs.Seed(e) })
			var out []string
			for _, d := range diagList(env) {
				if strings.Contains(d.Msg, "is passed on to") {
					out = append(out, fmt.Sprintf("%d:%s:%d:%s", d.Severity, d.Method, d.Pos, d.Msg))
				}
			}
			sort.Strings(out)
			return out
		}

		canonical := run(names)

		shuffled := append([]string{}, names...)
		for i := len(shuffled) - 1; i > 0; i-- {
			j := rapid.IntRange(0, i).Draw(rt, "shuf")
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		}
		alt := run(shuffled)

		if !strSlicesEqual(canonical, alt) {
			rt.Fatalf("ref order changed findings\n canon=%v\n shuf=%v\n order=%v",
				canonical, alt, shuffled)
		}
		want := 0
		if bad && conflictAt == 0 {
			want = 1
		}
		if len(canonical) != want {
			rt.Fatalf("want %d findings (bad=%v conflictAt=%d depth=%d), got %v",
				want, bad, conflictAt, depth, canonical)
		}
	})
}

// graft a requirement chain into an error-free wellTyped synth program: a
// satisfied requirement must add no error anywhere, a violated one must add
// exactly the one finding at the grafted call site - base program shapes
// (members, narrowing, self-calls) must neither mask it nor amplify it.
func TestPropRequirementWellTypedAntifault(t *testing.T) {
	withSigs(t,
		Reg{Kind: "free", Name: "Scanner", Sig: "(string :string) :unknown"},
	)
	rapid.Check(t, func(rt *rapid.T) {
		p := synth.GenProgram(rt, synth.Config{WellTyped: true})
		src := synth.Render(&p)
		clsA, okA := safeParse(src, "T")
		if !okA {
			rt.Skip("unparseable")
		}
		envA, pa := safeRun(clsA)
		if pa != "" {
			rt.Fatalf("panic: %s\nsrc:\n%s", pa, src)
		}
		if _, bad := anySeverityError(envA); bad {
			rt.Skip("base wellTyped program already has an error")
		}

		depth := rapid.IntRange(1, 3).Draw(rt, "gdepth")
		bad := rapid.Bool().Draw(rt, "gbad")
		lit := `"s"`
		if bad {
			lit = "0"
		}
		var g strings.Builder
		for k := 1; k <= depth; k++ {
			if k == depth {
				fmt.Fprintf(&g, "\tZ%d(q) { Scanner(q)\n\t\treturn q }\n", k)
			} else {
				fmt.Fprintf(&g, "\tZ%d(q) { return .Z%d(q) }\n", k, k+1)
			}
		}
		fmt.Fprintf(&g, "\tZu() { return .Z1(%s) }\n", lit)
		grafted := strings.TrimSuffix(src, "\t}\n") + g.String() + "\t}\n"
		clsB, okB := safeParse(grafted, "T")
		if !okB {
			rt.Fatalf("graft did not parse:\n%s", grafted)
		}
		envB, pb := safeRun(clsB)
		if pb != "" {
			rt.Fatalf("panic (grafted): %s\nsrc:\n%s", pb, grafted)
		}
		var errs []Diagnostic
		for _, d := range diagList(envB) {
			if d.Severity == SeverityError {
				errs = append(errs, d)
			}
		}
		if !bad {
			if len(errs) != 0 {
				rt.Fatalf("satisfied requirement errored: %+v\nsrc:\n%s", errs, grafted)
			}
			return
		}
		if len(errs) != 1 || errs[0].Method != "Zu" ||
			!strings.Contains(errs[0].Msg, "is passed on to Scanner which requires string") {
			rt.Fatalf("want exactly the Zu requirement error, got %+v\nsrc:\n%s", errs, grafted)
		}
	})
}

// guess interleaving. X family: the only demand on p sits behind a
// guessed sig (r.Prefix?/Tr on an unknown receiver bind guessed String
// params) - no requirement may form, so the mistyped Xu call stays silent.
// Y family: a guessed receiver use of the same param must not suppress the
// genuine annotation-grounded chain below it.
func TestPropRequirementGuessInterleave(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Prefix?", Sig: "(string :string, pos=0) :boolean"},
		Reg{Kind: "free", Name: "Scanner", Sig: "(string :string) :unknown"},
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "string", Name: "Tr", Sig: "(from :string, to :string ='') :string"},
	)
	rapid.Check(t, func(rt *rapid.T) {
		p := synth.GenProgram(rt, synth.Config{WellTyped: true})
		src := synth.Render(&p)
		clsA, okA := safeParse(src, "T")
		if !okA {
			rt.Skip("unparseable")
		}
		envA, pa := safeRun(clsA)
		if pa != "" {
			rt.Fatalf("panic: %s\nsrc:\n%s", pa, src)
		}
		if _, bad := anySeverityError(envA); bad {
			rt.Skip("base wellTyped program already has an error")
		}

		guessed := rapid.SampledFrom([]string{"Prefix?", "Tr"}).Draw(rt, "guessed")
		depth := rapid.IntRange(1, 2).Draw(rt, "ydepth")
		var g strings.Builder
		fmt.Fprintf(&g, "\tX1(r, p) { return r.%s(p) }\n", guessed)
		g.WriteString("\tXu() { return .X1(\"s\", 0) }\n")
		for k := 1; k <= depth; k++ {
			body := fmt.Sprintf("return .Y%d(r)", k+1)
			if k == depth {
				body = "Scanner(r)\n\t\treturn r"
			}
			if k == 1 {
				body = "r.Size()\n\t\t" + body
			}
			fmt.Fprintf(&g, "\tY%d(r) { %s }\n", k, body)
		}
		g.WriteString("\tYu() { return .Y1(0) }\n")
		grafted := strings.TrimSuffix(src, "\t}\n") + g.String() + "\t}\n"
		clsB, okB := safeParse(grafted, "T")
		if !okB {
			rt.Fatalf("graft did not parse:\n%s", grafted)
		}
		envB, pb := safeRun(clsB)
		if pb != "" {
			rt.Fatalf("panic (grafted): %s\nsrc:\n%s", pb, grafted)
		}
		var errs []Diagnostic
		for _, d := range diagList(envB) {
			if d.Severity == SeverityError {
				errs = append(errs, d)
			}
			if strings.Contains(d.Msg, "is passed on to") &&
				(d.Method == "Xu" || d.Method == "X1") {
				rt.Fatalf("requirement formed behind a guessed sig: %+v\nsrc:\n%s", d, grafted)
			}
		}
		if len(errs) != 1 || errs[0].Method != "Yu" ||
			!strings.Contains(errs[0].Msg, "is passed on to Scanner which requires string") {
			rt.Fatalf("guessed receiver use suppressed the genuine chain: %+v\nsrc:\n%s", errs, grafted)
		}
	})
}
