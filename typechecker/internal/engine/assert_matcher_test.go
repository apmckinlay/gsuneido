// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"strings"
	"testing"
)

func assertDiags(t *testing.T, src string) []string {
	t.Helper()
	_, env := runPasses(src, "T")
	var errs []string
	for _, d := range diagList(env) {
		if d.Severity == SeverityError {
			errs = append(errs, d.Msg)
		}
	}
	return errs
}

// matcher-style Assert proves its claim on fall-through, same as the
// old-style predicate form walkBlock already handles.
func TestAssertMatcherNarrows(t *testing.T) {
	cases := []struct{ label, src, wantErr string }{
		{"isString", `class {
	F(x) {
		Assert(x isString:)
		return x + 1
		}
	}`, `operator "+" expects number, got string`},
		{"isString with msg", `class {
	F(x) {
		Assert(x isString:, msg: "oops")
		return x + 1
		}
	}`, `operator "+" expects number, got string`},
		{"isNumber", `class {
	F(x) {
		Assert(x isNumber:)
		return x $ "a"
		}
	}`, `operator "$" expects string, got number`},
		{"is literal", `class {
	F(x) {
		Assert(x is: 5)
		return x $ "a"
		}
	}`, `operator "$" expects string, got number`},
		{"isType string", `class {
	F(x) {
		Assert(x isType: "Number")
		return x $ "a"
		}
	}`, `operator "$" expects string, got number`},
		{"member subject", `class {
	F() {
		Assert(.x isNumber:)
		return .x $ "a"
		}
	}`, `operator "$" expects string, got number`},
	}
	for _, c := range cases {
		errs := assertDiags(t, c.src)
		found := false
		for _, e := range errs {
			if strings.Contains(e, c.wantErr) {
				found = true
			}
		}
		if !found {
			t.Errorf("%s: matcher fact not applied, errors=%v", c.label, errs)
		}
	}
}

// matchers that prove nothing about the subject's type stay silent, and the
// msg arg is never mistaken for a matcher.
func TestAssertMatcherConservative(t *testing.T) {
	cases := []string{
		`class {
	F(x) {
		Assert(x has: "a")
		return x + 1
		}
	}`,
		`class {
	F(x) {
		Assert(x greaterThan: 3)
		return x + 1
		}
	}`,
	}
	for _, src := range cases {
		errs := assertDiags(t, src)
		for _, e := range errs {
			if strings.Contains(e, "expects number") {
				t.Errorf("non-type matcher narrowed the subject: %v", e)
			}
		}
	}
}
