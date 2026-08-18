// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package typechecker

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/typechecker/annotations"
	"github.com/apmckinlay/gsuneido/typechecker/internal/engine"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// Reg is a signature registration, as TypeChecker.RegisterSignatures takes.
type Reg = annotations.Registration

// withSigs registers sigs for this test only, through the same entry point
// stdlib uses. The signature table is process-wide, so it is restored on
// cleanup and tests using this must not run in parallel.
func withSigs(t *testing.T, sigs ...Reg) {
	t.Helper()
	prev := engine.Annotations()
	if _, err := RegisterSignatures(sigs); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { engine.SetAnnotations(prev) })
}

// runRequest runs method over refs/args and fails the test on an error.
func runRequest(t *testing.T, method string, refs, args []SourceEntry) Result {
	t.Helper()
	res, err := Process(Request{Method: method, References: refs, Arguments: args})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	return res
}

func inferOne(t *testing.T, refs []SourceEntry, name, src string) TypeInfo {
	t.Helper()
	res := runRequest(t, "TypeInfer", refs, []SourceEntry{{name, src}})
	return res.Results[0].(TypeInfo)
}

// classDiags returns every diagnostic (warning or error) reported against class.
func classDiags(res Result, class string) []ResultDiagnostic {
	var out []ResultDiagnostic
	for _, d := range res.Diagnostics.Warnings {
		if d.Class == class {
			out = append(out, d)
		}
	}
	for _, d := range res.Diagnostics.Errors {
		if d.Class == class {
			out = append(out, d)
		}
	}
	return out
}

func TestParseConfidenceFilter(t *testing.T) {
	keeps := func(raw string, c float64) bool {
		f, err := parseConfidenceFilter(map[string]string{"confidence": raw})
		if err != nil {
			t.Fatalf("parse %q: %v", raw, err)
		}
		return f.passes(c)
	}
	// absent / empty / "all" -> no filtering (everything passes)
	for _, m := range []map[string]string{nil, {"confidence": ""}, {"confidence": "all"}} {
		f, err := parseConfidenceFilter(m)
		assert.T(t).True(err == nil)
		assert.T(t).True(f.passes(0.0))
	}
	assert.T(t).True(keeps(">=0.70", 0.90))
	assert.T(t).True(keeps(">=0.70", 0.70))
	assert.T(t).False(keeps(">=0.70", 0.40))
	assert.T(t).False(keeps(">0.70", 0.70))
	assert.T(t).True(keeps("<0.70", 0.40))
	assert.T(t).False(keeps("<0.70", 0.70))
	assert.T(t).True(keeps("<=0.40", 0.40))
	assert.T(t).True(keeps("==0.90", 0.90))
	assert.T(t).False(keeps("==0.90", 0.70))
	// bare number is a >= minimum
	assert.T(t).True(keeps("0.70", 0.90))
	assert.T(t).False(keeps("0.70", 0.40))

	if _, err := parseConfidenceFilter(map[string]string{"confidence": ">=high"}); err == nil {
		t.Errorf("want error for malformed threshold")
	}
}

func TestProcessNoPhantomDiags(t *testing.T) {
	withSigs(t) // no builtin signatures needed
	a := assert.T(t)
	flagged := `class {
		Bad(plus, minus)
			{
			return plus - minus
			}
		}`
	clean := `class {
		Ok()
			{
			}
		}`
	res := runRequest(t, "TypeInfer", nil, []SourceEntry{
		{"Flagged", flagged},
		{"Clean", clean},
	})

	for _, d := range classDiags(res, "Clean") {
		t.Errorf("phantom diagnostic on Clean: %+v", d)
	}

	// Sanity check: the real diagnostics on Flagged still arrive.
	sawFlagged := false
	for _, d := range res.Diagnostics.Warnings {
		if d.Class == "Flagged" && d.Method == "Bad" {
			sawFlagged = true
			break
		}
	}
	a.That(sawFlagged)
}

func TestProcess_ReferencesNotReported(t *testing.T) {
	withSigs(t) // no builtin signatures needed
	a := assert.T(t)
	// `plus - minus` on untyped params trips an operator diagnostic
	ref := `class { Bad(plus, minus) { return plus - minus } }`
	arg := `class { Ok() { } }`
	res := runRequest(t, "TypeInfer",
		[]SourceEntry{{"RefBad", ref}}, []SourceEntry{{"Main", arg}})

	// results are parallel to arguments only - the reference is not one
	a.This(len(res.Results)).Is(1)

	for _, d := range classDiags(res, "RefBad") {
		t.Errorf("reference diagnostic leaked: %+v", d)
	}
}

func TestProcessWithReference(t *testing.T) {
	withSigs(t) // no builtin signatures needed
	a := assert.T(t)
	ref := `class { Value() { return 5 } }`
	arg := `class {
		Get()
			{
			c = new Counter()
			return c.Value()
			}
		}`
	got := inferOne(t, []SourceEntry{{"Counter", ref}}, "Main", arg)
	a.This(got.Methods["Get"]["$return"]).Is("number")
}

func TestProcess_StaticMemberReadTypesLocal(t *testing.T) {
	withSigs(t) // no builtin signatures needed
	a := assert.T(t)
	ipsum := `class {
		I: "ip-sum"
		CallClass() { return Date() }
	}`
	hello := `class {
		Foo()
			{
			b = Ipsum.I
			cls = Ipsum()
			}
	}`
	foo := inferOne(t, []SourceEntry{{"Ipsum", ipsum}}, "Hello", hello).Methods["Foo"]
	a.This(foo["b"]).Is("string")
	a.This(foo["cls"]).Is("date")
}

func TestProcessClassRefSeedView(t *testing.T) {
	withSigs(t) // no builtin signatures needed
	a := assert.T(t)
	ref := `class {
		x: false
		New() { .x = "init" }
		Foo() { return .x }
	}`
	arg := `class {
		OnClass() { return Widget.Foo() }
		OnInstance() { return (new Widget()).Foo() }
	}`
	got := inferOne(t, []SourceEntry{{"Widget", ref}}, "Main", arg)
	a.This(got.Methods["OnInstance"]["$return"]).Is("string")
	// fold order isn't fixed, so accept either rendering of the seed union.
	onClass := got.Methods["OnClass"]["$return"]
	a.That(onClass == "false|string" || onClass == "string|false")
}

func TestProcessNoReference(t *testing.T) {
	withSigs(t) // no builtin signatures needed
	a := assert.T(t)
	arg := `class {
		Get()
			{
			c = new Counter()
			return c.Value()
			}
		}`
	got := inferOne(t, nil, "Main", arg)
	// the method is still reported (so consumers see it exists)...
	_, hasMethod := got.Methods["Get"]
	a.That(hasMethod)
	_, hasReturn := got.Methods["Get"]["$return"]
	a.That(!hasReturn)
}

func TestProcess_ReportsParamsLocalsAndMembers(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Add", Sig: "(@args) :object"},
	)
	a := assert.T(t)
	src := `class {
		limit: 10
		Add(x = 0)
			{
			sum = x + .limit
			return sum
			}
		}`
	got := inferOne(t, nil, "Adder", src)

	add := got.Methods["Add"]
	a.This(add["x"]).Is("number")       // parameter (from its default)
	a.This(add["sum"]).Is("number")     // local variable
	a.This(add["$return"]).Is("number") // method return
	a.This(got.Members["limit"]).Is("number")
}

func TestTypeAnnotateNativeSyntax(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Add", Sig: "(@args) :object"},
	)
	a := assert.T(t)
	src := `class {
	limit: 10
	Add(x = 0)
		{
		sum = x + .limit
		return sum
		}
	}`
	res := runRequest(t, "TypeAnnotate", nil, []SourceEntry{{"Adder", src}})
	out := res.Results[0].(string)
	t.Logf("annotated:\n%s", out)

	// parameter -> native colon syntax on the signature, before the default
	a.That(strings.Contains(out, "Add(x :number = 0)"))
	// return type -> native colon syntax between ) and {
	a.That(strings.Contains(out, ") :number"))
	// local variable -> inline comment, not native syntax
	a.That(strings.Contains(out, "sum /* number */"))
	// member -> inline :: comment, not native (would collide with `limit: 10`)
	a.That(strings.Contains(out, "limit: 10 /* :: number */"))
}

func TestTypeAnnotateRespellsExistingAnnotations(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Add", Sig: "(@args) :object"},
		Reg{Kind: "free", Name: "Object", Sig: "(@args) :object"},
	)
	a := assert.T(t)
	src := `class {
	Add(x: number = 0) : number
		{
		return x
		}
	Pick(s: String | False)
		{
		return s
		}
	Make() : Point
		{
		return Object()
		}
	}`
	res := runRequest(t, "TypeAnnotate", nil, []SourceEntry{{"Adder", src}})
	out := res.Results[0].(string)
	t.Logf("annotated:\n%s", out)

	// a declared annotation is respelled in place, never duplicated
	a.That(strings.Contains(out, "Add(x :number = 0) :number"))
	a.That(!strings.Contains(out, ":number :number"))
	// builtins lowercase, union arms joined by a bare |, both declared...
	a.That(strings.Contains(out, "Pick(s :string|false)"))
	a.That(!strings.Contains(out, "String | False"))
	// ...and inferred
	a.That(strings.Contains(out, ") :false|string"))
	// a class name keeps its capitalization
	a.That(strings.Contains(out, "Make() :Point"))
}

// an annotation split over lines is left as written rather than reflowed
func TestTypeAnnotateLeavesMultilineAnnotation(t *testing.T) {
	withSigs(t) // no builtin signatures needed
	a := assert.T(t)
	src := `class {
	Foo(x: number)
		: String
		{
		return "a"
		}
	}`
	res := runRequest(t, "TypeAnnotate", nil, []SourceEntry{{"Foo", src}})
	out := res.Results[0].(string)
	t.Logf("annotated:\n%s", out)

	a.That(strings.Contains(out, "Foo(x :number)\n\t\t: String\n"))
	a.That(!strings.Contains(out, ":string"))
}

func TestTypeAnnotateIsIdempotent(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Add", Sig: "(@args) :object"},
	)
	a := assert.T(t)
	src := `class {
	limit: 10
	Add(x = 0)
		{
		sum = x + .limit
		return sum
		}
	}`
	once := runRequest(t, "TypeAnnotate", nil,
		[]SourceEntry{{"Adder", src}}).Results[0].(string)
	twice := runRequest(t, "TypeAnnotate", nil,
		[]SourceEntry{{"Adder", once}}).Results[0].(string)
	a.This(twice).Is(once)
}
