// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

type arityParam struct {
	name       string
	hasDefault bool
}

type aritySig struct {
	params []arityParam
	atRest bool // lone @args variadic
}

type arityCall struct {
	positional int
	named      []string
	spread     bool // @spread present
}

// modelArity is the reference verdict as a sorted key multiset.
func modelArity(sig aritySig, call arityCall) []string {
	if sig.atRest || call.spread {
		return nil
	}
	if call.positional > len(sig.params) {
		return []string{"tooMany|E"}
	}
	named := map[string]bool{}
	for _, n := range call.named {
		named[n] = true
	}
	var keys []string
	for i, p := range sig.params {
		if i < call.positional && named[p.name] {
			keys = append(keys, "double:"+p.name+"|W")
		}
	}
	var missing []string
	for i, p := range sig.params {
		if !p.hasDefault && i >= call.positional && !named[p.name] {
			missing = append(missing, p.name)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		keys = append(keys, "missing:"+strings.Join(missing, ",")+"|E")
	}
	sort.Strings(keys)
	return keys
}

func sevLetter(s Severity) string {
	if s == SeverityError {
		return "E"
	}
	return "W"
}

func actualArity(env TypeEnv, method string) []string {
	var keys []string
	for _, d := range diagList(env) {
		if d.Method != method {
			continue
		}
		switch {
		case strings.Contains(d.Msg, "too many arguments"):
			keys = append(keys, "tooMany|"+sevLetter(d.Severity))
		case strings.Contains(d.Msg, "passed both positionally and by name"):
			keys = append(keys, "double:"+doubleBindParam(d.Msg)+"|"+sevLetter(d.Severity))
		case strings.Contains(d.Msg, "missing argument"):
			list := d.Msg[strings.LastIndex(d.Msg, ": ")+2:]
			parts := strings.Split(list, ", ")
			sort.Strings(parts)
			keys = append(keys, "missing:"+strings.Join(parts, ",")+"|"+sevLetter(d.Severity))
		}
	}
	sort.Strings(keys)
	return keys
}

func doubleBindParam(msg string) string {
	const pre = `argument "`
	_, after, ok := strings.Cut(msg, pre)
	if !ok {
		return "?"
	}
	rest := after
	if before, _, ok := strings.Cut(rest, `"`); ok {
		return before
	}
	return "?"
}

// ---- generators ----

func genAritySig(t *rapid.T) aritySig {
	if rapid.IntRange(0, 4).Draw(t, "atrest") == 0 {
		return aritySig{atRest: true}
	}
	np := rapid.IntRange(0, 3).Draw(t, "arnp")
	params := make([]arityParam, np)
	sawDefault := false
	for i := range np {
		hasDef := sawDefault || rapid.Bool().Draw(t, "ardef")
		if hasDef {
			sawDefault = true // Suneido wants defaults trailing
		}
		params[i] = arityParam{name: fmt.Sprintf("p%d", i), hasDefault: hasDef}
	}
	return aritySig{params: params}
}

func genArityCall(t *rapid.T, sig aritySig) arityCall {
	if rapid.IntRange(0, 5).Draw(t, "arspread") == 0 {
		return arityCall{spread: true}
	}
	positional := rapid.IntRange(0, len(sig.params)+2).Draw(t, "arpos")
	var named []string
	for _, p := range sig.params {
		if rapid.Bool().Draw(t, "arnamed") {
			named = append(named, p.name)
		}
	}
	return arityCall{positional: positional, named: named}
}

// ---- renderer ----

func renderAritySig(sig aritySig) string {
	if sig.atRest {
		return "@args"
	}
	parts := make([]string, len(sig.params))
	for i, p := range sig.params {
		if p.hasDefault {
			parts[i] = p.name + " = 1"
		} else {
			parts[i] = p.name
		}
	}
	return strings.Join(parts, ", ")
}

func renderArityClass(sig aritySig, call arityCall) string {
	var b strings.Builder
	b.WriteString("class\n\t{\n")
	b.WriteString("\tM0()\n\t\t{\n")
	if call.spread {
		b.WriteString("\t\tv0 = #(1, 2, 3)\n\t\t.M1(@v0)\n")
	} else {
		var args []string
		for i := 0; i < call.positional; i++ {
			args = append(args, fmt.Sprintf("%d", i+1))
		}
		for _, n := range call.named {
			args = append(args, n+": 1")
		}
		fmt.Fprintf(&b, "\t\t.M1(%s)\n", strings.Join(args, ", "))
	}
	b.WriteString("\t\t}\n")
	fmt.Fprintf(&b, "\tM1(%s)\n\t\t{\n\t\t}\n", renderAritySig(sig))
	b.WriteString("\t}\n")
	return b.String()
}

// ---- property ----

func TestPropArityModel(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		sig := genAritySig(rt)
		call := genArityCall(rt, sig)
		src := renderArityClass(sig, call)
		cls, ok := safeParse(src, "T")
		if !ok {
			rt.Fatalf("arity source did not parse:\n%s", src)
		}
		env, panicMsg := safeRun(cls)
		if panicMsg != "" {
			rt.Fatalf("pipeline panicked: %s\nsrc:\n%s", panicMsg, src)
		}
		want := modelArity(sig, call)
		got := actualArity(env, "M0")
		if !strSlicesEqual(want, got) {
			rt.Fatalf("arity verdict mismatch\n want=%v\n got =%v\nsrc:\n%s", want, got, src)
		}
	})
}

func TestArityKnownCases(t *testing.T) {
	cases := []struct {
		sig  aritySig
		call arityCall
		want []string
	}{
		{aritySig{params: []arityParam{{"p0", false}, {"p1", false}}}, arityCall{positional: 1}, []string{"missing:p1|E"}},
		{aritySig{params: []arityParam{{"p0", false}, {"p1", false}}}, arityCall{positional: 0}, []string{"missing:p0,p1|E"}},
		{aritySig{params: []arityParam{{"p0", false}, {"p1", false}}}, arityCall{positional: 3}, []string{"tooMany|E"}},
		{aritySig{params: []arityParam{{"p0", false}}}, arityCall{positional: 1, named: []string{"p0"}}, []string{"double:p0|W"}},
		{aritySig{params: []arityParam{{"p0", false}, {"p1", true}}}, arityCall{positional: 1}, nil},
	}
	for _, c := range cases {
		if got := modelArity(c.sig, c.call); !strSlicesEqual(got, c.want) {
			t.Fatalf("model sig=%v call=%v: want %v got %v", c.sig, c.call, c.want, got)
		}
		src := renderArityClass(c.sig, c.call)
		_, env := runPasses(src, "T")
		if got := actualArity(env, "M0"); !strSlicesEqual(got, c.want) {
			t.Fatalf("actual: want %v got %v\nsrc:\n%s", c.want, got, src)
		}
	}
}

func TestArityUnresolvedAndSpreadNoDiag(t *testing.T) {
	for _, src := range []string{
		`class { M0() { .Munknown(1, 2, 3) } }`,                        // unknown own method
		`class { M0() { SomeGlobal.Foo(1, 2, 3) } }`,                   // unknown global receiver
		`class { M0() { v0 = #(1, 2)  .M1(@v0) } M1(p0, p1, p2) { } }`, // @spread
		`class { M0() { .M1(1, 2, 3, 4, 5) } M1(@args) { } }`,          // variadic callee
	} {
		_, env := runPasses(src, "T")
		if got := actualArity(env, "M0"); len(got) > 0 {
			t.Fatalf("expected no arity diagnostics for %q, got %v", src, got)
		}
	}
}
