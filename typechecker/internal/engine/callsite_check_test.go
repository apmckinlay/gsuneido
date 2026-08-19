// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"strings"
	"testing"
)

// annotated params on user-defined methods are hard contracts, enforced per callsite
func TestUserMethodArgTypeCheck(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
	)
	cases := []struct {
		label, src string
		wantSev    int // -1 = expect no argument diagnostic at all
		want       string
	}{
		{"positional mismatch",
			"class { F(x :number) { return x }\nG() { return .F(\"s\") } }",
			int(SeverityError), `argument "x": type string not assignable to declared number`},
		{"named mismatch",
			"class { F(x :number) { return x }\nG() { return .F(x: \"s\") } }",
			int(SeverityError), `argument "x": type string not assignable to declared number`},
		{"clean fit is silent",
			"class { F(x :number) { return x }\nG() { return .F(1) } }",
			-1, ""},
		{"guess-derived arg caps at warning",
			"class { F(x :string) { return x }\nG(p) { v = p.Size()\nreturn .F(v) } }",
			int(SeverityWarning), "argument \"x\": type number not assignable to declared string (type guessed from builtin overloads)"},
		{"named overrides positional - overridden value not checked",
			"class { F(x :number) { return x }\nG() { return .F(\"s\", x: 1) } }",
			-1, ""},
		{"named overrides positional - named value still checked",
			"class { F(x :number) { return x }\nG() { return .F(1, x: \"s\") } }",
			int(SeverityError), `argument "x": type string not assignable to declared number`},
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
			if !strings.Contains(d.Msg, "assignable to declared") {
				continue
			}
			if c.wantSev == -1 {
				t.Errorf("%s: unexpected argument diagnostic: %+v", c.label, d)
				continue
			}
			if int(d.Severity) == c.wantSev && strings.Contains(d.Msg, c.want) {
				found = true
			} else {
				t.Errorf("%s: wrong argument diagnostic: %+v", c.label, d)
			}
		}
		if c.wantSev != -1 && !found {
			t.Errorf("%s: expected sev=%d diagnostic containing %q", c.label, c.wantSev, c.want)
		}
	}
}

// massage silently drops a named arg that matches no param
func TestDiscardedNamedArgWarning(t *testing.T) {
	cases := []struct {
		label, src string
		want       bool
	}{
		{"misspelled named arg",
			"class { F(x) { return x }\nG() { return .F(1, bogus: 2) } }", true},
		{"matching named arg",
			"class { F(x) { return x }\nG() { return .F(x: 1) } }", false},
		{"block exempt",
			"class { F(x) { return x }\nG() { return .F(1, block: 2) } }", false},
		{"guessed sig exempt",
			"class { G(p) { return p.Size(bogus: 1) } }", false},
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
			if strings.Contains(d.Msg, "discarded at runtime") {
				found = true
				if d.Severity != SeverityWarning {
					t.Errorf("%s: discarded-named must be a warning: %+v", c.label, d)
				}
			}
		}
		if found != c.want {
			t.Errorf("%s: discarded-named warning present=%v, want %v", c.label, found, c.want)
		}
	}
}
