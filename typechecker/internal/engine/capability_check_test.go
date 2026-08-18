// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// `Query1(.q).name` throws exactly like `Query1(.q).Size()` does: gsuneido
// evaluates a member read as Get first, and SuBool.Get panics outright
// instead of returning nil, so it never reaches the method-lookup fallback.
func TestMemberReadOnFalseArmErrors(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Match", Sig: "(pattern, pos=false, prev :boolean =false) :false|object"},
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		q: ""
		E() {
			x = "s".Match(.q).Size()
			return x
		}
		N() {
			y = "s".Match(.q).name
			return y
		}
	}`, "T")
	a.That(hasDiag(env, "E", SeverityError, "not applicable on at least one path"))
	a.That(hasDiag(env, "N", SeverityError, "not readable on at least one path"))
	a.That(hasDiag(env, "N", SeverityError, "no members on false"))
	a.That(hasDiag(env, "N", SeverityError, "on `'s'.Match(this.q)`"))
}

// the read and the call are the same defect, so they carry the same weight
func TestMemberReadAndCallAgreeOnSeverity(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		q: ""
		E() { return "s".Match(.q).Size() }
		N() { return "s".Match(.q).name }
	}`, "T")
	a.This(errorCount(env, "E")).Is(errorCount(env, "N"))
}

func TestMemberReadOnPlainFalseErrors(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			x = false
			return x.name
		}
	}`, "T")
	a.That(hasDiag(env, "Foo", SeverityError, "not readable on receiver of type false"))
}

// a boolean rejects every member name, method names included - SuBool.Get
// panics before Lookup can bind one
func TestMemberReadOfMethodNameOnFalseErrors(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			x = false
			return x.Size
		}
	}`, "T")
	a.That(hasDiag(env, "Foo", SeverityError, "not readable on receiver of type false"))
}

// this-member reads chain: the receiver of the outer read is .q's type
func TestMemberReadThroughThisMemberUnion(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Match", Sig: "(pattern, pos=false, prev :boolean =false) :false|object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		New() { .q = "s".Match("x") }
		Foo() { return .q.name }
	}`, "T")
	a.That(hasDiag(env, "Foo", SeverityError, "not readable on at least one path"))
}

func TestMemberReadCleanReceiversSilent(t *testing.T) {
	cases := []struct {
		label, src string
	}{
		{"object receiver", `class { Foo() { return Object().name } }`},
		{"unknown member on object is unknowable", `class { Foo() { return Object().nosuch } }`},
		{"this member", `class { q: ""  Foo() { return .q } }`},
		{"class static", `class { Foo() { return Suneido.name } }`},
		{"string receiver", `class { Foo() { return "s".name } }`},
		{"unknown receiver", `class { Foo(p) { return p.name } }`},
		{"loop variable", `class { Foo(ob) { for x in ob { return x.name } return "" } }`},
		{"narrowed away", `class {
			Foo() {
				x = "s".Match("y")
				if x is false
					{ return "" }
				return x.name
			}
		}`},
		// join across branches then narrow: the arm the guard leaves is clean
		{"reassigned then guarded", `class {
			Foo(c) {
				x = false
				if c is true
					{ x = Object() }
				if x is false
					{ return "" }
				return x.name
			}
		}`},
	}
	for _, c := range cases {
		_, env := runPasses(c.src, "T")
		if hasDiag(env, "Foo", SeverityError, "not readable") {
			t.Errorf("%s: unexpected member-read error", c.label)
		}
	}
}

// a receiver whose type rests on a guess cannot prove anything
func TestMemberReadOnAssumedReceiverWarns(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Match", Sig: "(pattern, pos=false, prev :boolean =false) :false|object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(p) {
			v = p.Match("x")
			return v.name
		}
	}`, "T")
	a.That(!hasDiag(env, "Foo", SeverityError, "not readable"))
	a.That(hasDiag(env, "Foo", SeverityWarning, "not readable"))
}

// `rec = Query1(q)` then `rec.name = x` is the same defect on the write side:
// SuBool.Put panics just like SuBool.Get does
func TestMemberWriteOnFalseArmErrors(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Match", Sig: "(pattern, pos=false, prev :boolean =false) :false|object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			rec = "s".Match("x")
			rec.name = 1
		}
	}`, "T")
	a.That(hasDiag(env, "Foo", SeverityError, "not writable on at least one path"))
}

// a compound assignment reads before it writes, so it stays a read
func TestMemberUpdateOnFalseArmReportsRead(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Match", Sig: "(pattern, pos=false, prev :boolean =false) :false|object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			rec = "s".Match("x")
			rec.n += 1
		}
	}`, "T")
	a.That(hasDiag(env, "Foo", SeverityError, "not readable on at least one path"))
}

// the call check already owns `.Size()`; the read check must not double up
func TestMemberReadDoesNotDuplicateCallDiagnostic(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Match", Sig: "(pattern, pos=false, prev :boolean =false) :false|object"},
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() { return "s".Match("x").Size() }
	}`, "T")
	a.This(countDiagsWith(env, "not readable")).Is(0)
	a.This(countDiagsWith(env, "not applicable")).Is(1)
}

// ---- the other positions a boolean arm reaches ------------------------

// `ob[0]` is a Mem whose member is not an identifier, so it lands in the same
// walk as `ob.name` - and SuBool.Get panics for both
func TestSubscriptOnFalseArmErrors(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Match", Sig: "(pattern, pos=false, prev :boolean =false) :false|object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			rec = "s".Match("x")
			return rec[0]
		}
	}`, "T")
	a.That(hasDiag(env, "Foo", SeverityError, "subscript on `rec` not readable on at least one path"))
	a.That(hasDiag(env, "Foo", SeverityError, "no members on false"))
}

func TestRangeOnFalseArmErrors(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Match", Sig: "(pattern, pos=false, prev :boolean =false) :false|object"},
	)
	cases := []struct{ label, src string }{
		{"range-to", `class { Foo() { rec = "s".Match("x")  return rec[1..2] } }`},
		{"range-len", `class { Foo() { rec = "s".Match("x")  return rec[1::2] } }`},
	}
	for _, c := range cases {
		_, env := runPasses(c.src, "T")
		if !hasDiag(env, "Foo", SeverityError, "range on `rec` not applicable on at least one path") {
			t.Errorf("%s: expected a range error", c.label)
		}
		if !hasDiag(env, "Foo", SeverityError, "false is not rangeable") {
			t.Errorf("%s: expected the rangeable note", c.label)
		}
	}
}

func TestCallOfFalseArmErrors(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			f = false
			return f()
		}
	}`, "T")
	a.That(hasDiag(env, "Foo", SeverityError, "call on `f` not applicable on receiver of type false"))
}

func TestIterationOverFalseArmErrors(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Match", Sig: "(pattern, pos=false, prev :boolean =false) :false|object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			rec = "s".Match("x")
			for f in rec
				{ return f }
			return ""
		}
	}`, "T")
	a.That(hasDiag(env, "Foo", SeverityError, "iteration on `rec` not applicable on at least one path"))
	a.That(hasDiag(env, "Foo", SeverityError, "false is not iterable"))
}

// each of these is legal Suneido that a naive "not an Object" rule would flag
func TestBooleanUseCleanPositionsSilent(t *testing.T) {
	cases := []struct{ label, src string }{
		// `"Foo"(ob)` is the string-as-method idiom - strings are callable
		{"string callee", `class { Foo(ob) { f = "Trim"  return f(ob) } }`},
		{"global callee", `class { Foo() { return Object() } }`},
		{"function local", `class { Foo() { f = function () { return 1 }  return f() } }`},
		// SuStr.Iter exists, so strings iterate
		{"iterate a string", `class { Foo() { for c in "abc" { return c } return "" } }`},
		{"iterate an object", `class { Foo() { for x in Object(1, 2) { return x } return 0 } }`},
		{"for-range is not iteration", `class { Foo() { for i in 1..10 { return i } return 0 } }`},
		{"subscript an object", `class { Foo() { return Object(1)[0] } }`},
		{"range a string", `class { Foo() { return "abc"[1..2] } }`},
		{"unknown receiver", `class { Foo(p) { return p[0] } }`},
		{"iterate unknown", `class { Foo(p) { for x in p { return x } return 0 } }`},
		{"narrowed away", `class {
			Foo() {
				rec = "s".Match("y")
				if rec is false
					{ return "" }
				for f in rec
					{ return f }
				return ""
			}
		}`},
	}
	for _, c := range cases {
		_, env := runPasses(c.src, "T")
		for _, d := range methodDiags(env, "Foo") {
			if d.Severity == SeverityError {
				t.Errorf("%s: unexpected error: %s", c.label, d.Msg)
			}
		}
	}
}

// ---- capability sets beyond booleans ----------------------------------
// every row below was verified by executing the position in gsuneido:
// iterate/range are supported only by String, Object, Record and Sequence;
// call additionally by Function, Class, Block and Instance; member writes
// only by Instance, Object, Record and Sequence.

func capSrc(bind, use string) string {
	return "class {\n\tFoo(c) {\n\t\tx = " + bind + "\n\t\t" + use + "\n\t}\n}"
}

// a lone receiver names the type; only the union form carries a note clause
func expectCap(t *testing.T, bind, use, subject, verb, typeName string) {
	t.Helper()
	_, env := runPasses(capSrc(bind, use), "T")
	want := subject + " on `x` not " + verb + " on receiver of type " + typeName
	if !hasDiag(env, "Foo", SeverityError, want) {
		t.Errorf("%s: expected %q, got %v", bind, want, methodDiags(env, "Foo"))
	}
}

const iterUse = "for y in x\n\t\t\t{ return y }\n\t\treturn 0"

func TestIterationOverNonIterableArms(t *testing.T) {
	cases := []struct{ bind, typeName string }{
		{"5", "number"},
		{"#20240102", "date"},
		{"function () { return 1 }", "function"},
		{"{|y| y }", "block"},
		{"class { }", "class"},
		{"new this()", "T"},
	}
	for _, c := range cases {
		expectCap(t, c.bind, iterUse, "iteration", "applicable", c.typeName)
	}
}

// `for x in Contributions(...)` - Contributions inherits CallClass from
// Memoize, so the call is not an instance and the loop is legal
func TestIterationOverInheritedCallClassResult(t *testing.T) {
	_, env := runPassesWithRefList(`class {
		Foo() {
			for fn in Contributions("AllowDelete")
				return fn
			return 0
		}
	}`, []RefSource{
		{Name: "Memoize", Src: `class {
			CallClass(@args) {
				result = Suneido.GetInit('m', { LruCache(.Func) }).Get(@args)
				if Object?(result) and not result.Readonly?()
					result = result.Copy()
				return result
			}
		}`},
		{Name: "Contributions", Src: `Memoize { Func(name) { return Object() } }`},
	})
	for _, d := range methodDiags(env, "Foo") {
		if d.Severity == SeverityError {
			t.Errorf("unexpected error: %s", d.Msg)
		}
	}
}

// the shape that prompted this: a union where only one arm can iterate
func TestIterationOverUnionWithNonIterableArm(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(c) {
			x = 5
			if c is true
				{ x = Object() }
			for y in x
				{ return y }
			return 0
		}
	}`, "T")
	a.That(hasDiag(env, "Foo", SeverityError,
		"iteration on `x` not applicable on at least one path of union receiver number|object; number is not iterable"))
}

func TestRangeOverNonRangeableArms(t *testing.T) {
	cases := []struct{ bind, typeName string }{
		{"5", "number"},
		{"#20240102", "date"},
		{"function () { return 1 }", "function"},
		{"{|y| y }", "block"},
		{"class { }", "class"},
		{"new this()", "T"},
	}
	for _, c := range cases {
		expectCap(t, c.bind, "return x[1..2]", "range", "applicable", c.typeName)
		expectCap(t, c.bind, "return x[1::2]", "range", "applicable", c.typeName)
	}
}

func TestCallOfNonCallableArms(t *testing.T) {
	withSigs(t,
		Reg{Kind: "free", Name: "Seq", Sig: "(from=false, to=false, by=1) :sequence"},
	)
	cases := []struct{ bind, typeName string }{
		{"5", "number"},
		{"#20240102", "date"},
		{"Seq(0, 3)", "sequence"},
	}
	for _, c := range cases {
		expectCap(t, c.bind, "return x()", "call", "applicable", c.typeName)
	}
}

func TestMemberWriteOnNonWritableArms(t *testing.T) {
	cases := []struct{ bind, typeName string }{
		{"5", "number"},
		{`"s"`, "string"},
		{"#20240102", "date"},
		{"function () { return 1 }", "function"},
		{"{|y| y }", "block"},
		{"class { }", "class"},
	}
	for _, c := range cases {
		expectCap(t, c.bind, "x.foo = 1\n\t\treturn x", `member "foo"`, "writable", c.typeName)
	}
	// subscript writes go through the same Put that member writes do
	expectCap(t, "5", "x[0] = 1\n\t\treturn x", "subscript", "writable", "number")
}

// only the union form carries the note clause
func TestCapabilityUnionNotes(t *testing.T) {
	cases := []struct{ label, use, want string }{
		{"iterate", iterUse, "number is not iterable"},
		{"range", "return x[1..2]", "number is not rangeable"},
		{"call", "return x()", "number is not callable"},
		{"write", "x.foo = 1\n\t\treturn x", "number has no writable members"},
		{"read", "return x.foo", "no members on false"},
	}
	for _, c := range cases {
		bind := "5\n\t\tif c is true\n\t\t\t{ x = Object() }"
		if c.label == "read" {
			bind = "false\n\t\tif c is true\n\t\t\t{ x = Object() }"
		}
		_, env := runPasses(capSrc(bind, c.use), "T")
		if !hasDiag(env, "Foo", SeverityError, c.want) {
			t.Errorf("%s: expected %q, got %v", c.label, c.want, methodDiags(env, "Foo"))
		}
	}
}

// the supported cells of the matrix must stay silent
func TestCapabilitySupportedPositionsSilent(t *testing.T) {
	cases := []struct{ label, bind, use string }{
		{"iterate string", `"abc"`, "for y in x { return y } return 0"},
		{"iterate object", "Object(1)", "for y in x { return y } return 0"},
		{"iterate sequence", "Seq(0, 3)", "for y in x { return y } return 0"},
		{"range string", `"abc"`, "return x[1..2]"},
		{"range object", "Object(1)", "return x[1..2]"},
		{"range sequence", "Seq(0, 3)", "return x[1..2]"},
		{"call function", "function () { return 1 }", "return x()"},
		{"call block", "{|y| y }", "return x(1)"},
		{"call class", "class { }", "return x()"},
		// Object(1,2)(k) succeeds, and Instance is callable via a Call method
		{"call object", "Object(1)", "return x(0)"},
		{"call instance", "new this()", "return x()"},
		// strings are callable with a `this` argument
		{"call string", `"Trim"`, "return x(Object())"},
		{"write object", "Object(1)", "x.foo = 1  return x"},
		{"write sequence", "Seq(0, 3)", "x.foo = 1  return x"},
		// instances accept member writes, and Getter_ can serve any read
		{"write instance", "new this()", "x.foo = 1  return x"},
		{"read instance", "new this()", "return x.foo"},
		{"read number method", "5", "return x.Round"},
		{"read date method", "#20240102", "return x.Year"},
	}
	for _, c := range cases {
		_, env := runPasses(capSrc(c.bind, c.use), "T")
		for _, d := range methodDiags(env, "Foo") {
			if d.Severity == SeverityError {
				t.Errorf("%s: unexpected error: %s", c.label, d.Msg)
			}
		}
	}
}
