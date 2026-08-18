// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestAritySelfCall(t *testing.T) {
	a := assert.T(t)
	msgs := arityMsgs(`class {
		foo(x) {}
		bar() { .foo() }
		baz() { .foo(1, "hello") }
	}`)
	a.This(len(msgs)).Is(2)
	j := strings.Join(msgs, " | ")
	a.That(strings.Contains(j, "missing argument to foo: x"))
	a.That(strings.Contains(j, "too many arguments to foo: 2 given, takes at most 1"))
}

func TestArityCorrectCallSilent(t *testing.T) {
	assert.T(t).This(len(arityMsgs(`class {
		foo(x) {}
		bar() { .foo(42) }
	}`))).Is(0)
}

func TestArityDefaultsAndOmitMiddle(t *testing.T) {
	assert.T(t).This(len(arityMsgs(`class {
		box(w, h = 10, color = "white") { return 0 }
		a() { .box(5) }
		b() { .box(5, color: "red") }
		c() { .box(w: 1, h: 2, color: "x") }
		d() { .box(color: "x", w: 1) }
	}`))).Is(0)
}

func TestArityMissingRequiredWithDefaultsPresent(t *testing.T) {
	a := assert.T(t)
	j := joined(`class {
		box(w, h = 10) { return 0 }
		a() { .box() }
	}`)
	a.That(strings.Contains(j, "missing argument to box: w"))
}

// @args callee accepts anything - never checked.
func TestArityVariadicCalleeSilent(t *testing.T) {
	assert.T(t).This(len(arityMsgs(`class {
		collect(@args) { return args }
		a() { .collect() }
		b() { .collect(1, 2, 3, x: 9) }
	}`))).Is(0)
}

// @ / @+1 spread at the call site makes the positional count unknown - bail.
func TestAritySpreadArgSilent(t *testing.T) {
	assert.T(t).This(len(arityMsgs(`class {
		foo(x) {}
		a(ob) { .foo(@ob) }
		b(ob) { .foo(@+1 ob) }
	}`))).Is(0)
}

func TestArityDynamicParamOptional(t *testing.T) {
	assert.T(t).This(len(arityMsgs(`class {
		themed(_color) { return color }
		a() { .themed() }
	}`))).Is(0)
}

func TestArityExtraNamedIgnored(t *testing.T) {
	assert.T(t).This(len(arityMsgs(`class {
		foo(x) {}
		a() { .foo(1, bogus: 9) }
	}`))).Is(0)
}

func TestArityExtraPositionalOverflows(t *testing.T) {
	a := assert.T(t)
	j := joined(`class {
		foo(x) {}
		a() { .foo(1, 2, bogus: 9) }
	}`)
	a.That(strings.Contains(j, "too many arguments to foo: 2 given, takes at most 1"))
}

func TestArityTrailingBlock(t *testing.T) {
	assert.T(t).This(len(arityMsgs(`class {
		applyBlock(block) { return block() }
		a() { .applyBlock() { 42 } }
	}`))).Is(0)
}

func TestArityConstructionNew(t *testing.T) {
	a := assert.T(t)
	a.That(strings.Contains(joined(`class {
		New(x) {}
		a() { return new this() }
	}`), "missing argument to new this: x"))
	a.That(strings.Contains(joined(`class {
		New(x) {}
		a() { return new this(1, 2) }
	}`), "too many arguments to new this"))
	a.This(len(arityMsgs(`class {
		New(x) {}
		a() { return new this(1) }
	}`))).Is(0)
}

func TestArityDefaultCtorZeroArgs(t *testing.T) {
	a := assert.T(t)
	j := joined(`class {
		foo() {}
		a() { return new this(1, 2, 3) }
	}`)
	a.That(strings.Contains(j, "too many arguments to new this: 3 given, takes at most 0"))
}

func TestArityMemberInitPublicNamed(t *testing.T) {
	assert.T(t).This(len(arityMsgs(`class {
		New(.Size, .Color = "white") {}
		a() { return new this(size: 5) }
	}`))).Is(0)
}

func TestArityConstructionWithBaseSilent(t *testing.T) {
	assert.T(t).This(len(arityMsgs(`class : SomeBase {
		foo() {}
		a() { return new this(1, 2, 3) }
	}`))).Is(0)
}

func TestArityNewBypassesCallClass(t *testing.T) {
	assert.T(t).This(len(arityMsgs(`class {
		CallClass(when) { return new this(when, false) }
		New(at, skipWeekends) {}
		go() { return new this(1, 2) }
	}`))).Is(0)
}

func TestArityBareCallUsesCallClass(t *testing.T) {
	a := assert.T(t)
	j := joined(`class {
		CallClass(when) { return 0 }
		New(at, skipWeekends) {}
		go() { return T(1, 2) }
	}`)
	a.That(strings.Contains(j, "too many arguments to T: 2 given, takes at most 1"))
}

func TestArityAbstractStubSkipped(t *testing.T) {
	assert.T(t).This(len(arityMsgs(`Base {
		Output() { throw 'MUST IMPLEMENT IN DERIVED CLASS' }
		BeforeTest(name) { .Output(name) }
	}`))).Is(0)
}

func TestArityReferenceStaticCall(t *testing.T) {
	a := assert.T(t)
	classes, sigs := refOf("Json",
		`class { Encode(value) { return "" } Decode(s) { return s } }`)

	bad := arityMsgsRef(`class { go() { return Json.Encode(1, 2) } }`, classes, sigs)
	a.That(strings.Contains(strings.Join(bad, " | "),
		"too many arguments to Encode: 2 given, takes at most 1"))

	missing := arityMsgsRef(`class { go() { return Json.Encode() } }`, classes, sigs)
	a.That(strings.Contains(strings.Join(missing, " | "),
		"missing argument to Encode: value"))

	ok := arityMsgsRef(`class { go() { return Json.Encode(1) } }`, classes, sigs)
	a.This(len(ok)).Is(0)
}

func TestArityReferenceInstanceCall(t *testing.T) {
	a := assert.T(t)
	classes, sigs := refOf("Parser",
		`class { New(src) {} Statement(depth) { return depth } }`)

	bad := arityMsgsRef(`class {
		go() { p = new Parser("x"); return p.Statement(1, 2) }
	}`, classes, sigs)
	a.That(strings.Contains(strings.Join(bad, " | "),
		"too many arguments to Statement: 2 given, takes at most 1"))

	ok := arityMsgsRef(`class {
		go() { p = new Parser("x"); return p.Statement(1) }
	}`, classes, sigs)
	a.This(len(ok)).Is(0)
}

func TestArityReferenceConstruction(t *testing.T) {
	a := assert.T(t)
	classes, sigs := refOf("Parser",
		`class { CallClass(a, b) { return new this(a) } New(src) {} }`)

	bad := arityMsgsRef(`class { go() { return new Parser(1, 2) } }`, classes, sigs)
	a.That(strings.Contains(strings.Join(bad, " | "),
		"too many arguments to new Parser: 2 given, takes at most 1"))

	// bare Parser(...) goes through CallClass(a, b) instead
	bare := arityMsgsRef(`class { go() { return Parser(1, 2, 3) } }`, classes, sigs)
	a.That(strings.Contains(strings.Join(bare, " | "),
		"too many arguments to Parser: 3 given, takes at most 2"))
}

// bare Derived(...) binds the base's CallClass, the way SuClass.Call does
func TestArityBareCallUsesInheritedCallClass(t *testing.T) {
	a := assert.T(t)
	refs := []RefSource{
		{Name: "Base", Src: `class { CallClass(a) { return a } }`},
		{Name: "Derived", Src: `Base { Other() { return 1 } }`},
	}
	bad := arityMsgsRefList(`class { go() { return Derived(1, 2, 3) } }`, refs)
	a.That(strings.Contains(strings.Join(bad, " | "),
		"too many arguments to Derived: 3 given, takes at most 1"))

	a.This(len(arityMsgsRefList(`class { go() { return Derived(1) } }`, refs))).Is(0)
}

// an own CallClass overrides the inherited one
func TestArityOwnCallClassBeatsInherited(t *testing.T) {
	a := assert.T(t)
	bad := arityMsgsRefList(`class { go() { return Derived(1, 2) } }`, []RefSource{
		{Name: "Base", Src: `class { CallClass(a, b) { return a } }`},
		{Name: "Derived", Src: `Base { CallClass(a) { return a } }`},
	})
	a.That(strings.Contains(strings.Join(bad, " | "),
		"too many arguments to Derived: 2 given, takes at most 1"))
}

// no CallClass in the chain: the call constructs, so New binds - also inherited
func TestArityBareCallUsesInheritedNew(t *testing.T) {
	a := assert.T(t)
	bad := arityMsgsRefList(`class { go() { return Derived(1, 2) } }`, []RefSource{
		{Name: "Base", Src: `class { New(a) {} }`},
		{Name: "Derived", Src: `Base { Other() { return 1 } }`},
	})
	a.That(strings.Contains(strings.Join(bad, " | "),
		"too many arguments to Derived: 2 given, takes at most 1"))
}

// the base is not among the references, so the chain is unproven: no check
func TestArityUnknownBaseSkipsCall(t *testing.T) {
	a := assert.T(t)
	a.This(len(arityMsgsRefList(`class { go() { return Derived(1, 2, 3) } }`,
		[]RefSource{{Name: "Derived", Src: `Missing { Other() { return 1 } }`}}))).Is(0)
}

func TestArityInheritedMethod(t *testing.T) {
	a := assert.T(t)
	classes, sigs := refOf("Base", `class { foo(x) { return x } }`)

	bad := arityMsgsRef(`Base { go() { return .foo(1, 2, 3) } }`, classes, sigs)
	a.That(strings.Contains(strings.Join(bad, " | "),
		"too many arguments to foo: 3 given, takes at most 1"))

	a.This(len(arityMsgsRef(`Base { go() { return .foo(1) } }`, classes, sigs))).Is(0)

	// an own override shadows the inherited method: own foo(a, b) needs 2
	over := arityMsgsRef(`Base {
		foo(a, b) { return a + b }
		go() { return .foo(1) }
	}`, classes, sigs)
	a.That(strings.Contains(strings.Join(over, " | "), "missing argument to foo: b"))
}

func TestAritySuperCall(t *testing.T) {
	a := assert.T(t)
	classes, sigs := refOf("Base", `class { New(x) {} bar(y) { return y } }`)

	badNew := arityMsgsRef(`Base { New() { super(1, 2, 3) } }`, classes, sigs)
	a.That(strings.Contains(strings.Join(badNew, " | "),
		"too many arguments to super: 3 given, takes at most 1"))

	badBar := arityMsgsRef(`Base { go() { return super.bar(1, 2) } }`, classes, sigs)
	a.That(strings.Contains(strings.Join(badBar, " | "),
		"too many arguments to bar: 2 given, takes at most 1"))

	a.This(len(arityMsgsRef(`Base { New() { super(1) } }`, classes, sigs))).Is(0)
}

func TestArityInheritedUnknownBaseSilent(t *testing.T) {
	assert.T(t).This(len(arityMsgsRef(
		`Unseen { go() { return .foo(1, 2, 3) } }`, nil, nil))).Is(0)
}

func TestArityUnknownGlobalSilent(t *testing.T) {
	assert.T(t).This(len(arityMsgsRef(
		`class { go() { return Mystery.Whatever(1, 2, 3) } }`, nil, nil))).Is(0)
}

func TestAritySuneidoSemantics(t *testing.T) {
	a := assert.T(t)
	classes, sigs := refOf("Arity", `class
		{
		New() { }
		Input(x, y = false) :false|number
			{
			if y isnt false
				return x + Number(y)
			return false
			}
		}`)
	silent := func(call string) {
		t.Helper()
		msgs := arityMsgsRef(`class { go() { return `+call+` } }`, classes, sigs)
		a.This(len(msgs)).Is(0)
	}
	silent(`Arity().Input(1, 2)`)
	silent(`Arity().Input(1, y: 2)`)
	silent(`Arity().Input(x: 1, y: 2)`)

	silent(`Arity().Input(1, 2, y: 4)`)
	warns := overrideWarns(`class { go() { return Arity().Input(1, 2, y: 4) } }`, classes, sigs)
	a.This(len(warns)).Is(1)
	a.That(strings.Contains(warns[0], `argument "y"`))

	bad := arityMsgsRef(`class { go() { return Arity().Input(1, 2, 3) } }`, classes, sigs)
	a.That(strings.Contains(strings.Join(bad, " | "),
		"too many arguments to Input: 3 given, takes at most 2"))
}

func TestArityDoubleBindWarning(t *testing.T) {
	a := assert.T(t)
	warns := overrideWarns(`class {
		add(x, y) { return x + y }
		go() { return .add(1, 2, y: 4) }
	}`, nil, nil)
	a.This(len(warns)).Is(1)
	a.That(strings.Contains(warns[0], `argument "y"`))
	a.That(strings.Contains(warns[0], "this may be confusing"))

	a.This(len(overrideWarns(`class {
		add(x, y) { return x + y }
		go() { return .add(1, y: 2) }
	}`, nil, nil))).Is(0)

	a.This(len(overrideWarns(`class {
		add(x, y) { return x + y }
		go() { return .add(1, 2) }
	}`, nil, nil))).Is(0)

	a.This(len(overrideWarns(`class {
		add(x, y) { return x + y }
		go() { return .add(1, 2, 3, y: 4) }
	}`, nil, nil))).Is(0)
}

func TestArityDoubleBindConfidence(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		add(x, y) { return x + y }
		go() { return .add(1, 2, y: 4) }
	}`, "T")
	found := false
	for _, d := range *env.Diagnostics {
		if strings.Contains(d.Msg, "passed both positionally and by name") {
			found = true
			a.This(ScoreConfidence(&d)).Is(0.70)
		}
	}
	a.That(found)
}

func TestArityUnnamedAfterNamedIsParseError(t *testing.T) {
	a := assert.T(t)
	defer func() {
		r := recover()
		a.That(r != nil)
		a.That(strings.Contains(fmt.Sprint(r), "named"))
	}()
	ParseClass(`class { go() { return Arity().Input(x: 1, 2) } }`)
	t.Fatal("expected a parse error for an un-named argument after a named one")
}
