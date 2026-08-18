// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestLiteralMembers(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class { x: 42  s: "hello"  b: false }`, "T")
	a.This(env.Members["x"]).Is(TNumber)
	a.This(env.Members["s"]).Is(TString)
	a.This(env.Members["b"]).Is(TFalse)
}

func TestLocalBoolLiteralsNotWidened(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			x = false
			return x
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TFalse)
}

func TestMemberBoolWideningCollapses(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		active: false
		Activate() { .active = true }
	}`, "T")
	a.This(env.Members["active"]).Is(TBoolean)
}

func TestInitMembers(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Init() {
			.count = 5
			.label = "hi"
		}
	}`, "T")
	a.This(env.Members["count"]).Is(TNumber)
	a.This(env.Members["label"]).Is(TString)
}

func TestGetCountResolvesViaInit(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Init() { .count = 5 }
		GetCount() { return .count }
	}`, "T")
	a.This(env.Returns["GetCount"]).Is(TNumber)
}

func TestCompoundAssignment(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Compound() {
			x = false
			x += 1
			return x
		}
	}`, "T")
	a.This(env.Returns["Compound"]).Is(TNumber)
}

func TestSymbol(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Sym() { return #foo }
	}`, "T")
	a.This(env.Returns["Sym"]).Is(TString)
}

func TestTernaryUnion(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Pick(c) { return c is 1 ? "yes" : 42 }
	}`, "T")
	ret := env.Returns["Pick"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TString) && u.Contains(TNumber))
}

func TestMemberUnionAcrossMethods(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Init()  { .v = 1 }
		Reset() { .v = "none" }
	}`, "T")
	u, ok := env.Members["v"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber) && u.Contains(TString))
}

func TestInlineParamSeedsScope(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		count: 0
		Bump(.count) {
			return count
		}
	}`, "T")
	a.This(env.Returns["Bump"]).Is(TNumber)
}

func TestUProducesDirtyUnion(t *testing.T) {
	a := assert.T(t)
	result := U(TFalse, TUnknown)
	u, ok := result.(Union)
	a.That(ok)
	a.That(u.IsDirty)
	a.That(u.Contains(TFalse))
}

func TestDirtyReturnUnionPreservesConcreteType(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Find(x) {
			if x is 0
				return false
			return .unknown_member
		}
	}`, "T")
	ret := env.Returns["Find"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.IsDirty)
	a.That(u.Contains(TFalse))
}

func TestDirtyUnionReceiverResolvesByKnownMember(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Append", Sig: "(ob) :object"},
		Reg{Receiver: "object", Name: "Copy", Sig: "() :object"},
		Reg{Receiver: "sequence", Name: "Copy", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Difference", Sig: "(other) :object"},
		Reg{Kind: "static", Class: "Database", Name: "Load", Sig: "(table :string, from :string = '', privateKey :string = '', passphrase :string = '') :number"},
		Reg{Kind: "static", Class: "Ftsearch", Name: "Load", Sig: "(data :string) :unknown"},
		Reg{Receiver: "object", Name: "Map", Sig: "(block) :object"},
		Reg{Receiver: "string", Name: "Map", Sig: "(block) :string"},
		Reg{Kind: "free", Name: "Object", Sig: "(@args) :object"},
		Reg{Kind: "free", Name: "Query1", Sig: "(@args) :false|object"},
		Reg{Receiver: "string", Name: "Replace", Sig: "(pattern, block = '', count = false) :string"},
		Reg{Receiver: "object", Name: "Replace", Sig: "(oldvalue, newvalue) :object"},
		Reg{Receiver: "date", Name: "Replace", Sig: "(@args) :date"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		querySelect() { return Query1('x') }
		getSelects() {
			if false is d = .querySelect()
				return #()
			return d.userselect_selects
		}
		defaultSel: #()
		Load() {
			.defaultSel = .getSelects().Copy()
			hisSel = .getSelects()
			hisSel = hisSel.Map({ it })
			combinedSel = .defaultSel.Copy().Append(hisSel.Difference(#()))
			return combinedSel
		}
		// isolates the Map step: a multi-overload builtin on a dirty union
		// must pick the Object overload, not fold in the String overload.
		MapOnly() {
			hisSel = .getSelects()
			return hisSel.Map({ it })
		}
		// generality guard: dispatch is annotation-table-driven, not keyed on
		// specific method names. Replace has three overloads (Object->object,
		// String->string, Date->date); on a dirty Object union it must still
		// collapse to the known member's overload (Object), never
		// Object | String | Date. add a future builtin and it behaves the same.
		ReplaceCase() {
			x = .getSelects()
			return x.Replace('a', 'b')
		}
	}`, "T")
	a.This(env.Returns["MapOnly"]).Is(TObject)
	a.This(env.Returns["Load"]).Is(TObject)
	a.This(env.Returns["ReplaceCase"]).Is(TObject)
}

func TestReturnThrowDoesNotLeakIntoReturnType(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Calc(op) {
			if op is "+"
				return 1
			return throw "bad op"
		}
	}`, "T")
	// Without the fix: return type is Number | String (string from the throw).
	a.This(env.Returns["Calc"]).Is(TNumber)
}

func TestReturnThrowOnlyPathIsVoid(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Always() {
			return throw "nope"
		}
	}`, "T")
	a.This(env.Returns["Always"]).Is(TVoid)
}

func TestInlineInitParamWidensBoolDefault(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(.flag = false) {}
	}`, "T")
	a.This(env.Members["flag"]).Is(U(TFalse, TUnknown))
}

func TestInlineInitParamWidensThenUnionsCleanly(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(.flag = false) {}
		SetIt() { .flag = true }
	}`, "T")
	a.This(env.Members["flag"]).Is(U(TBoolean, TUnknown))
}

func TestInlineInitParamNonBoolNotWidened(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(.count = 0) {}
	}`, "T")
	a.This(env.Members["count"]).Is(TNumber)
}

func TestForInLoopVarTypedUnknown(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Iter(ob) {
			for x in ob {
				return x
			}
		}
	}`, "T")
	a.This(env.Returns["Iter"]).Is(TUnknown)
}

// `for k, v in ob`: both loop variables resolve to TUnknown in scope.
func TestForInTwoVarsBothUnknown(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Pairs(ob) {
			for k, v in ob {
				return v
			}
		}
	}`, "T")
	a.This(env.Returns["Pairs"]).Is(TUnknown)
}

func TestSequencePrimitivePrintsAsSequence(t *testing.T) {
	a := assert.T(t)
	a.This(TSequence.String()).Is("sequence")
	u := Union{Types: []DynType{TSequence, TString}}
	// Either order is valid; just check both are present.
	a.That(u.Contains(TSequence))
	a.That(u.Contains(TString))
}

func TestNarrowStringPredicate(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { if String?(x) { return x } }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestNarrowNumberPredicate(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { if Number?(x) { return x } }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestNarrowDatePredicate(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { if Date?(x) { return x } }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TDate)
}

func TestNarrowBooleanPredicate(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { if Boolean?(x) { return x } }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TBoolean)
}

// Object? returns true for Objects (and Records, which are treated as Objects).
func TestNarrowObjectPredicate(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { if Object?(x) { return x } }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TObject)
}

// Record? narrows to TObject because records are not tracked separately.
func TestNarrowRecordPredicate(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { if Record?(x) { return x } }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TObject)
}

func TestNarrowClassPredicate(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { if Class?(x) { return x } }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TClass)
}

// Function? covers functions, methods, and blocks per the docs.
func TestNarrowFunctionPredicate(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { if Function?(x) { return x } }
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TFunction) && u.Contains(TBlock))
}

// `x is 5` narrows x to TNumber inside the body.
func TestNarrowIsNumberLiteral(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { if x is 5 { return x } }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
}

// `x is "hi"` narrows x to TString.
func TestNarrowIsStringLiteral(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { if x is "hi" { return x } }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

// `x is false` narrows x to TFalse.
func TestNarrowIsFalseLiteral(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { if x is false { return x } }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TFalse)
}

func TestNarrowIsntFalseRemovesFromUnion(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Maybe() {
			if .pick is true { return 42 }
			return false
		}
		pick: false
		Foo() {
			x = .Maybe()
			if x isnt false { return x }
			return false
		}
	}`, "T")
	// Foo's return is the narrowed-x branch (Number) plus the explicit false (False).
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber) && u.Contains(TFalse))
}

// `Type(x) is "Number"` narrows the same way Number? does.
func TestNarrowTypeStringEq(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { if Type(x) is "Number" { return x } }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
}

// `A and B` composes both refinements in the body.
func TestNarrowAndChainAppliesBoth(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(a, b) {
			if Number?(a) and String?(b) { return a + 1 }
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestNarrowNotPredicateNoPanic(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			if not String?(x) { return x }
		}
	}`, "T")
	// x is "not String" but otherwise unknown - return type stays Unknown.
	a.This(env.Returns["Foo"]).Is(TUnknown)
}

func TestNarrowElseBranchNegated(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			x = .Maybe()
			if Number?(x) { return x }
			else { return x }   // x had False removed, leaving... well, dirty
		}
		Maybe() { if .pick is true { return 42 }; return false }
		pick: false
	}`, "T")
	ret := env.Returns["Foo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TNumber) && u.Contains(TFalse))
}

func TestNarrowTernaryThenArmNarrowed(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { return Number?(x) ? x + 1 : false }
	}`, "T")
	// then arm produces TNumber, else arm produces TFalse.
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber) && u.Contains(TFalse))
}

func TestNarrowDoesNotEscapeIfBody(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			if Number?(x) { 1 }
			return x
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TUnknown)
}

func TestNarrowDoesNotRegressCompoundAssignment(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			x = false
			x += 1
			return x
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestNarrowMembersInPureBody(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: 0
		Reset() { .x = "reset" }
		Foo() {
			if Number?(.x) { return .x }
			return false
		}
	}`, "T")
	ret := env.Returns["Foo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TFalse))
	a.That(!u.Contains(TString)) // body-scan-driven narrowing dropped String
}

func TestNarrowMemberIsntFalseRemovesFalse(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		width: 0
		Reset() { .width = false }
		Area(h: number) {
			if .width isnt false
				return .width * h
			return 0
		}
	}`, "T")
	errCount := 0
	for _, d := range *env.Diagnostics {
		if d.Method == "Area" && d.Severity == SeverityError {
			errCount++
		}
	}
	a.This(errCount).Is(0)
}

func TestNarrowMemberPredicateCallNarrows(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: 0
		Off() { .x = false }
		Foo() {
			if Number?(.x) { return .x }
			return "fallback"
		}
	}`, "T")
	ret := env.Returns["Foo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TString))
	a.That(!u.Contains(TFalse)) // narrowing must drop False from the if-path
}

func TestNarrowMemberEarlyReturnFallthrough(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: 0
		Off() { .x = false }
		Foo() {
			if not Number?(.x)
				return "no"
			return .x
		}
	}`, "T")
	ret := env.Returns["Foo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TString))
	a.That(!u.Contains(TFalse)) // fall-through return must see narrowed .x
}

func TestNarrowMemberCallClears(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: 0
		Off() { .x = false }
		Touch() { .x = false }
		Foo() {
			if Number?(.x) {
				.Touch()
				return .x
			}
			return "no"
		}
	}`, "T")
	ret := env.Returns["Foo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TFalse))
}

func TestNarrowMemberAssignClears(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: 0
		Off() { .x = false }
		Foo() {
			if Number?(.x) {
				.x = false
				return .x
			}
			return "no"
		}
	}`, "T")
	ret := env.Returns["Foo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TFalse))
}

func TestNarrowMemberGuardPureBody(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: 0
		Touch() { .x = Object() }
		Foo() {
			if Object?(.x)
				return .x
			return "no"
		}
	}`, "T")
	ret := env.Returns["Foo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TObject))
	a.That(u.Contains(TString))
	a.That(!u.Contains(TNumber)) // .x narrowed to Object, no leak from class-field 0
}

func TestNarrowLocalSurvivesCallInBody(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Touch() { }
		Foo(x) : number {
			if Number?(x)
				{
				.Touch()
				return x + 1
				}
			return 0
		}
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}

func TestNarrowMemberGuardNested(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: 0
		Touch() { .x = Object() }
		Foo() {
			if Object?(.x)
				if .x isnt false
					return .x
			return "no"
		}
	}`, "T")
	ret := env.Returns["Foo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TObject))
	a.That(!u.Contains(TFalse))
}

func TestNarrowMemberGuardSurvivesCallInBody(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		members: false
		SetObj() { .members = Object() }
		Foo(x) {
			if Object?(.members)
				.members.AddUnique(x)
		}
	}`, "T")
	for _, d := range methodDiags(env, "Foo") {
		a.That(!strings.Contains(d.Msg, "not applicable"))
	}
}

func TestNarrowMemberElseIfGuard(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		valid_timer: false
		record_change_members: false
		SetObj() { .record_change_members = Object() }
		Foo(member) {
			if .valid_timer is false
				{
				.record_change_members = Object(member)
				.valid_timer = true
				}
			else if Object?(.record_change_members)
				.record_change_members.AddUnique(member)
			else
				.record_change_members = Object(member)
		}
	}`, "T")
	for _, d := range methodDiags(env, "Foo") {
		a.That(!strings.Contains(d.Msg, "not applicable"))
	}
}

func TestNarrowMemberReadBeforeCallKeepsNarrowing(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: 0
		Touch() { .x = Object() }
		Clobber() { .x = false }
		Foo() {
			if Object?(.x)
				{
				y = .x
				.Clobber()
				return y
				}
			return "no"
		}
	}`, "T")
	ret := env.Returns["Foo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TObject))
	a.That(u.Contains(TString))
	a.That(!u.Contains(TFalse))
}

func TestNarrowMemberUnrelatedAssign(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: 0
		y: 0
		Touch() { .x = Object() }
		Foo() {
			if Object?(.x)
				{
				.y = false
				return .x
				}
			return "no"
		}
	}`, "T")
	ret := env.Returns["Foo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TObject))
	a.That(u.Contains(TString))
	a.That(!u.Contains(TNumber))
}

func TestNarrowMemberIncClearsNarrowing(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: 0
		Touch() { .x = false }
		Foo() {
			if Number?(.x)
				{
				.x++
				return .x
				}
			return "no"
		}
	}`, "T")
	ret := env.Returns["Foo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TFalse))
}

func TestNarrowMemberCompoundAssign(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: 0
		Touch() { .x = false }
		Foo() {
			if Number?(.x)
				{
				.x += 1
				return .x
				}
			return "no"
		}
	}`, "T")
	ret := env.Returns["Foo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TFalse))
}

func TestNarrowOrTrueBranchSameTargetUnion(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			if Number?(x) or String?(x) { return x }
		}
	}`, "T")
	ret := env.Returns["Foo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TNumber) && u.Contains(TString))
}

func TestNarrowOrTrueDifferentTargetsNoNarrowing(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x, y) {
			if Number?(x) or String?(y) { return x }
		}
	}`, "T")
	// x stays Unknown - operand 2 (about y) might be the true one
	a.This(env.Returns["Foo"]).Is(TUnknown)
}

func TestNarrowOrTrueThreeWaySameTarget(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			if Number?(x) or String?(x) or Date?(x) { return x }
		}
	}`, "T")
	ret := env.Returns["Foo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TNumber) && u.Contains(TString) && u.Contains(TDate))
}

func TestNarrowAndChainSameTarget(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: number | string) : string {
			if not (Number?(x) and Number?(x)) { return x }
			return ""
		}
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}

func TestNarrowAndVsOrSameTargetReturns(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			if String?(x) and Number?(x) { return x }
		}
		Bar(x) {
			if String?(x) or Number?(x) { return x }
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
	ret := env.Returns["Bar"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TNumber) && u.Contains(TString))
}

func TestNarrowReassignKeepsDispatch(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Replace", Sig: "(pattern, block = '', count = false) :string"},
		Reg{Receiver: "object", Name: "Replace", Sig: "(oldvalue, newvalue) :object"},
		Reg{Receiver: "date", Name: "Replace", Sig: "(@args) :date"},
		Reg{Kind: "free", Name: "String?", Sig: "(value) :boolean"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(s) {
			if not String?(s)
				throw "oops"
			s = s.Replace('a', 'b')
			return s.Replace('c', 'd')
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestNarrowDefaultValueParam(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(autoClose = false) {
			if Number?(autoClose)
				return autoClose
			return false
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TFalse))
}

func TestNarrowMemberAssignmentInstallsNarrowing(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		prevConditions: false
		Reset() { .prevConditions = false }
		Foo() {
			.prevConditions = Object()
			return .prevConditions
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TObject)
}

func TestNarrowDispatchReassignChain(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Replace", Sig: "(pattern, block = '', count = false) :string"},
		Reg{Receiver: "object", Name: "Replace", Sig: "(oldvalue, newvalue) :object"},
		Reg{Receiver: "date", Name: "Replace", Sig: "(@args) :date"},
		Reg{Kind: "free", Name: "String?", Sig: "(value) :boolean"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(s) {
			if String?(s)
				{
				// there is only one Replace on recv string which return string
				// bug was s was stamped as date|object|string
				s = s.Replace('xyz', '')
				s = s.Replace("abc", "")
				return s.Replace('def', '')
				}
			return false
		}
	}`, "T")
	ret := env.Returns["Foo"]
	// then-branch returns String, else returns false -> Union{String, False}.
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TString))
	a.That(u.Contains(TFalse))
	a.That(!u.Contains(TDate))
	a.That(!u.Contains(TObject))
}

func TestForLoopJoinKeepsEntryTypeForDispatch(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Replace", Sig: "(pattern, block = '', count = false) :string"},
		Reg{Receiver: "object", Name: "Replace", Sig: "(oldvalue, newvalue) :object"},
		Reg{Receiver: "date", Name: "Replace", Sig: "(@args) :date"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(value) {
			s = ""
			for i in .loop
				{
				s = value
				}
			s = s.Replace("", "")
			return s
		}
	}`, "T")
	ret := env.Returns["Foo"]
	a.This(ret).Is(TString)
}

func TestForLoopJoinEmitsNoAmbiguityWarning(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(value) {
			s = ""
			for i in .loop
				{
				s = value
				}
			s = s.Replace("", "")
			return s
		}
	}`, "T")
	a.This(countDiagsWith(env, "ambiguous overloads")).Is(0)
}

func TestNarrowBooleanMemberIsntFalseGivesTrue(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Flag: false
		Coerce() { .Flag = 0 }
		Set() { .Flag = true }
		Echo() {
			if .Flag isnt false
				return .Flag
			return "no"
		}
	}`, "T")
	ret := env.Returns["Echo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TTrue))
	a.That(u.Contains(TString))
	a.That(!u.Contains(TBoolean)) // would indicate the split didn't happen
	a.That(!u.Contains(TFalse))   // False must be gone from the if-path
}

func TestNarrowBooleanIsTrueGivesTrue(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Echo(x: number | boolean) {
			if x is true
				return x
			return 0
		}
	}`, "T")
	u, ok := env.Returns["Echo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TTrue))
	a.That(!u.Contains(TBoolean))
	a.That(!u.Contains(TFalse))
}

// `if x is false` narrows toward TFalse.
func TestNarrowBooleanIsFalseGivesFalse(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Echo(x: number | boolean) {
			if x is false
				return x
			return 0
		}
	}`, "T")
	u, ok := env.Returns["Echo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TFalse))
	a.That(!u.Contains(TBoolean))
	a.That(!u.Contains(TTrue))
}

// `if x isnt true` narrows away TTrue, leaving the TFalse half.
func TestNarrowBooleanIsntTrueGivesFalse(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Echo(x: number | boolean) {
			if x isnt true
				return x
			return 0
		}
	}`, "T")
	u, ok := env.Returns["Echo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TFalse))
	a.That(!u.Contains(TBoolean))
	a.That(!u.Contains(TTrue))
}

func TestNarrowIsntFalseSplitsBoolean(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Echo(x: number | boolean) {
			if x isnt false
				return x
			return 0
		}
	}`, "T")
	u, ok := env.Returns["Echo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TTrue))
	a.That(!u.Contains(TBoolean))
	a.That(!u.Contains(TFalse))
}

func TestNarrowAwayConcreteMatch(t *testing.T) {
	a := assert.T(t)
	got := narrowAwaySet(TFalse, []DynType{TFalse})
	a.This(got).Is(TUnknown)
}

func TestNarrowAwayDirtyFalse(t *testing.T) {
	a := assert.T(t)
	existing := Union{Types: []DynType{TFalse}, IsDirty: true}
	got := narrowAwaySet(existing, []DynType{TFalse})
	a.This(got).Is(TUnknown)
}

func TestNarrowAwayBooleanRemoved(t *testing.T) {
	a := assert.T(t)
	got := narrowAwaySet(TBoolean, []DynType{TBoolean})
	a.This(got).Is(TUnknown)
}

func TestNarrowAwaySet_NonEmptyKeepReturnsKept(t *testing.T) {
	a := assert.T(t)
	existing := Union{Types: []DynType{TFalse, TNumber}}
	got := narrowAwaySet(existing, []DynType{TFalse})
	a.This(got).Is(TNumber)
}

func TestNarrowAwaySet_NoMatchReturnsExisting(t *testing.T) {
	a := assert.T(t)
	got := narrowAwaySet(TNumber, []DynType{TFalse})
	a.This(got).Is(TNumber)
}

func TestNarrowAwayValueTypedConcrete(t *testing.T) {
	a := assert.T(t)
	got := narrowAwaySet(TNumber, []DynType{TNumber})
	a.This(got).Is(TNumber)
}

func TestNarrowAwayValueTypedUnion(t *testing.T) {
	a := assert.T(t)
	existing := Union{Types: []DynType{TNumber}}
	got := narrowAwaySet(existing, []DynType{TNumber})
	a.This(got).Is(existing)
}

func TestNarrowAwayEndToEnd(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Set(x = false) { .x = x }
		Use(h: number) {
			if .x isnt false
				return .x + h
			return 0
		}
	}`, "T")
	errs, warns := 0, 0
	for _, d := range *env.Diagnostics {
		if d.Method != "Use" {
			continue
		}
		switch d.Severity {
		case SeverityError:
			errs++
		case SeverityWarning:
			warns++
		}
	}
	a.This(errs).Is(0) // no error on the guarded read
	a.That(warns >= 1) // at least one "cannot prove Number" warning
}

func TestThrowNarrowLoopCarried(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		InLoop(name) {
			names = #('Line Items')
			if false is idx = names.Find(name)
				throw 'nope'
			while names.Member?(idx)
				idx += 0.1
			return idx
		}
	}`, "T")
	errs, warns := 0, 0
	for _, d := range *env.Diagnostics {
		if d.Method != "InLoop" {
			continue
		}
		switch d.Severity {
		case SeverityError:
			errs++
		case SeverityWarning:
			warns++
		}
	}
	a.This(errs).Is(0) // no spurious False arm, so no hard operator error
	a.That(warns >= 1) // Find's `?` arm survives, a dirty-operand warning is fine
}

func TestLoopCarriedNarrowing_NoGuardStillErrors(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Find", Sig: "(value) :false|unknown"},
		Reg{Receiver: "string", Name: "Find", Sig: "(string :string, pos=0) :number"},
		Reg{Receiver: "object", Name: "Member?", Sig: "(member) :boolean"},
		Reg{Receiver: "class", Name: "Member?", Sig: "(member) :boolean"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		NoGuard(name) {
			names = #('Line Items')
			idx = names.Find(name)
			while names.Member?(idx)
				idx += 0.1
			return idx
		}
	}`, "T")
	a.That(errorCount(env, "NoGuard") >= 1)
}

func TestLoopCarriedBodyReassign(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Find", Sig: "(value) :false|unknown"},
		Reg{Receiver: "string", Name: "Find", Sig: "(string :string, pos=0) :number"},
		Reg{Receiver: "object", Name: "Member?", Sig: "(member) :boolean"},
		Reg{Receiver: "class", Name: "Member?", Sig: "(member) :boolean"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		ReassignFalseInLoop(name) {
			names = #('Line Items')
			if false is idx = names.Find(name)
				throw 'nope'
			while names.Member?(idx)
				{
				idx = names.Find(name)
				idx += 0.1
				}
			return idx
		}
	}`, "T")
	a.That(errorCount(env, "ReassignFalseInLoop") >= 1)
}

func TestLoopCarriedReturnGuard(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		ReturnGuard(name) {
			names = #('Line Items')
			if false is idx = names.Find(name)
				return false
			while names.Member?(idx)
				idx += 0.1
			return idx
		}
	}`, "T")
	a.This(errorCount(env, "ReturnGuard")).Is(0)
}

func TestLoopCarriedNarrowing_ForLoopDowngrades(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		ForLoopCarried(name) {
			names = #('Line Items')
			if false is idx = names.Find(name)
				throw 'nope'
			for (j = 0; j < 3; j++)
				idx += 0.1
			return idx
		}
	}`, "T")
	a.This(errorCount(env, "ForLoopCarried")).Is(0)
}

func TestNarrowEarlyReturnNot(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			if not Number?(x)
				return false
			return x
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber) && u.Contains(TFalse))
}

func TestNarrowEarlyReturnOr(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			other = 5
			if other is 0 or not Number?(x)
				return false
			return x
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber) && u.Contains(TFalse))
}

func TestNarrowEarlyReturnThenExitsElseFalls(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			if not Number?(x) {
				return false
			} else {
				1
			}
			return x
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber) && u.Contains(TFalse))
}

func TestNarrowEarlyReturnElseExitsThenFalls(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			if Number?(x) {
				1
			} else {
				return false
			}
			return x
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber) && u.Contains(TFalse))
}

func TestNarrowEarlyReturnInlineAssignIdentOnRhs(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Maybe() { if .pick is true { return 42 }; return false }
		pick: false
		Foo() {
			if false is x = .Maybe()
				throw "unreachable"
			return x
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestNarrowEarlyReturnInlineAssignIdentOnLhs(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Maybe() { if .pick is true { return 42 }; return false }
		pick: false
		Foo() {
			if ((x = .Maybe()) is false)
				throw "unreachable"
			return x
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestNarrowContinueGuardAssign(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Maybe() { if .pick is true { return 42 }; return false }
		pick: false
		Foo() {
			for x in #(1, 2, 3) {
				if false is idx = .Maybe()
					continue
				n = idx + 1
			}
		}
	}`, "T")
	a.That(!hasOperatorError(env, "Foo"))
}

func TestNarrowBreakGuardCStyleFor(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Maybe() { if .pick is true { return "x" }; return false }
		pick: false
		Foo() {
			for (i = 0; ; ++i)
				{
				if false isnt dir = .Maybe()
					break
				if i > 8
					throw "nope"
				}
			return dir.Trim()
		}
	}`, "T")
	bad := false
	for _, d := range methodDiags(env, "Foo") {
		if d.Severity == SeverityError && strings.Contains(d.Msg, "type false") {
			bad = true
		}
	}
	a.That(!bad)
}

func TestBranchMergeUnionsDifferingBranchTypes(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: 0
		Foo() {
			v = false
			if .x < 0
				v = false
			else if .x is 0
				v = "invalid"
			else if .x > 0
				v = 1
			return v
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TFalse))
	a.That(u.Contains(TString))
	a.That(u.Contains(TNumber))
}

func TestBranchMergeWithoutElseKeepsEntryType(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: 0
		Foo() {
			v = false
			if .x < 0
				v = -1
			else if .x is 0
				v = 0
			else if .x > 0
				v = 1
			return v
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TFalse))
	a.That(u.Contains(TNumber))
}

func TestBranchMergeExhaustiveChainCollapses(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: 0
		Foo() {
			v = false
			if .x < 0
				v = -1
			else if .x is 0
				v = 0
			else
				v = 1
			return v
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestBranchMergeThenOnly(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: 0
		Foo() {
			v = 0
			if .x is 0
				v = "yes"
			return v
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TString))
}

// Only the else-branch assigns (then-block is empty): same union shape.
func TestBranchMergeElseOnly(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: 0
		Foo() {
			v = 0
			if .x is 0 {
			} else
				v = "no"
			return v
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TString))
}

func TestNarrowAfterBranchMergedVariable(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: 0
		Foo() {
			v = false
			if .x < 0
				v = false
			else if .x is 0
				v = "invalid"
			else if .x > 0
				v = 1
			if Type(v) is "Number"
				return v
			return false
		}
	}`, "T")
	// then-branch returns v narrowed to Number; else returns false.
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TFalse))
}

func TestNarrowReassignmentInsideGuardInvalidates(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(y) {
			if Type(y) is "Number" {
				y = "now string"
				if Type(y) is "String"
					return y
			}
			return false
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TString))
	a.That(u.Contains(TFalse))
}

func TestNarrowNestedSameTypeIdempotent(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(z) {
			if Type(z) is "Number"
				if Type(z) is "Number"
					return z
			return false
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TFalse))
}

func TestNarrowParamDefaultReassign(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Mem: 42
		Foo(x = false) {
			if x is false
				x = .Mem
			if x isnt false
				return x
			return false
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TFalse))
}

func TestNarrowParamDefaultForIn(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Custom: 42
		FocusFirst(parentHwnd, custom = false) {
			if custom is false
				custom = .Custom
			if custom isnt false
				for field in custom.Members()
					return custom
			return false
		}
	}`, "T")
	u, ok := env.Returns["FocusFirst"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TFalse))
}

func TestNarrowUnrecognisedTypeStringIsNoOp(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			if Type(x) is "Broski"
				return x
			return false
		}
	}`, "T")
	// x stays Unknown so the then-branch return is dirty.
	ret := env.Returns["Foo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.IsDirty)
	a.That(u.Contains(TFalse))
}

// for c in "literal" - loop var is TString.
func TestForInStringLiteralElementType(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			for c in "abc"
				return c
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestForInStringConstructorElementType(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			for c in String(123)
				return c
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestForInStringTwoVarsBothUnknown(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			for k, v in "abc"
				return v
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TUnknown)
}

func TestForInStringLocalVarElementType(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			s = "abc"
			for c in s
				return c
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestTypeConstructorNumber(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class { Foo() { return Number("42") } }`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestTypeConstructorString(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class { Foo() { return String(123) } }`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestTypeConstructorObject(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class { Foo() { return Object(1, 2, 3) } }`, "T")
	a.This(env.Returns["Foo"]).Is(TObject)
}

func TestTypeConstructorRecord(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class { Foo() { return Record(a: 1) } }`, "T")
	a.This(env.Returns["Foo"]).Is(TObject)
}

func TestTypeConstructorDisplay(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class { Foo() { return Display(123) } }`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestTypeConstructorTypeReturnsString(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class { Foo(x) { return Type(x) } }`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestTypeConstructorDateReturnsDateOrFalse(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class { Foo(s) { return Date(s) } }`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TDate))
	a.That(u.Contains(TFalse))
}

func TestTypeConstructorDateNoArgsReturnsDate(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class { Foo() { return Date() } }`, "T")
	a.This(env.Returns["Foo"]).Is(TDate)
}

func TestDateNarrowingValidStringLiteral(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "number", Name: "Format", Sig: "(mask :string) :string"},
		Reg{Receiver: "date", Name: "Format", Sig: "(format) :string"},
	)
	a := assert.T(t)
	_, env := runPasses(`class { Foo() { return Date("9:45").Format("Hmm") } }`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
	a.This(errorCount(env, "Foo")).Is(0)
}

func TestDateNarrowingNamedFields(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "date", Name: "Year", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() { return Date(year: 1999, month: 1, day: 1).Year() }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
	a.This(errorCount(env, "Foo")).Is(0)
}

func TestDateNarrowingNonConstantStaysUnion(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "number", Name: "Format", Sig: "(mask :string) :string"},
		Reg{Receiver: "date", Name: "Format", Sig: "(format) :string"},
	)
	a := assert.T(t)
	_, env := runPasses(`class { Foo(s) { return Date(s).Format("Hmm") } }`, "T")
	a.That(errorCount(env, "Foo") >= 1)
}

func TestDateNarrowingInvalidStringIsFalse(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class { Foo() { return Date("not a date") } }`, "T")
	a.This(env.Returns["Foo"]).Is(TFalse)
}

func TestDateNarrowingInvalidFieldsIsFalse(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class { Foo() { return Date(year: 3500, month: 1, day: 1) } }`, "T")
	a.This(env.Returns["Foo"]).Is(TFalse)
}

func TestBracketRecordLiteralResolvesAsObject(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class { Foo() { return [a: 1, b: 2] } }`, "T")
	a.This(env.Returns["Foo"]).Is(TObject)
}

// [1, 2, 3] desugars to Call(Ident("Object"),...).
func TestBracketObjectLiteralResolvesAsObject(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class { Foo() { return [1, 2, 3] } }`, "T")
	a.This(env.Returns["Foo"]).Is(TObject)
}

func TestNarrowSwitchCondModePredicate(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			switch {
				case Number?(x): return x
			}
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestNarrowSwitchCondModeTypeString(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			switch {
				case Type(x) is "String": return x
			}
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestNarrowSwitchIdentLiteral(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			switch x {
				case "hello": return x
			}
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestNarrowSwitchMultiExpr(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			switch x {
				case 1, "hi": return x
			}
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber) && u.Contains(TString))
}

func TestNarrowSwitchDefaultNarrowsAwayCaseTypes(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			switch {
				case String?(x): return false
				default: return x + 1
			}
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TFalse) && u.Contains(TNumber))
}

func TestNarrowSwitchCondModeNegatedPredicate(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			switch {
				case not Number?(x): return Type(x)
			}
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestNarrowSwitchDefaultStrips(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			switch x {
				case 1: return "num"
				case "hi": return "str"
				default: return x
			}
		}
	}`, "T")
	a.That(env.Returns["Foo"] != TUnknown)
}

func TestNarrowSwitchAssign(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			switch {
				case Number?(x): x = "now a string"; return x
			}
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestNarrowSwitchSkipsGlobalScrutinee(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			switch X {
				case 5: return X
			}
		}
	}`, "T")
	// X is global, never narrowed, so the return type stays unresolved.
	a.This(env.Returns["Foo"]).Is(TUnknown)
}

// empty switch is valid (if useless) - must walk without crashing.
func TestNarrowSwitchEmptyDoesNotPanic(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			switch x {
			}
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TVoid)
}

func TestNarrowSwitchNested(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			switch {
				case Number?(x): {
					switch x {
						case 0: return x
						default: return x
					}
				}
			}
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestNarrowSwitchDoesNotNarrowMembers(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		foo: 0
		Bar() {
			switch {
				case Number?(.foo): return .foo
			}
		}
	}`, "T")
	a.This(env.Returns["Bar"]).Is(TNumber)
	a.This(env.Members["foo"]).Is(TNumber)
}

func TestSwitchBranchMergeUnionsArmAssignments(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			v = false
			switch x {
				case 1: v = "one"
				case 2: v = 2
			}
			return v
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TFalse))
	a.That(u.Contains(TString))
	a.That(u.Contains(TNumber))
}

func TestSwitchBranchMergeWithDefaultDropsEntry(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			v = false
			switch x {
				case 1: v = "one"
				case 2: v = 2
				default: v = #sym
			}
			return v
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TString))
	a.That(u.Contains(TNumber))
	a.That(!u.Contains(TFalse))
}

func TestSwitchMergeExhaustive(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			v = false
			switch x {
				case 1: v = 1
				case 2: v = 2
				default: v = 0
			}
			return v
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestSwitchBranchMergeOnlyDefaultAssigns(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			v = false
			switch x {
				case 1:
				case 2:
				default: v = "fallback"
			}
			return v
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TFalse))
	a.That(u.Contains(TString))
}

func TestSwitchMergeSingleArm(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			v = 0
			switch x {
				case 1: v = "one"
			}
			return v
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TString))
}

func TestSwitchBranchMergeCondModeArmAssignments(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			v = false
			switch {
				case Number?(x): v = x
				case String?(x): v = x
			}
			return v
		}
	}`, "T")
	ret := env.Returns["Foo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.IsDirty)
	a.That(u.Contains(TFalse))
}

func TestSwitchBranchMergeNestedSwitch(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x, y) {
			v = false
			switch x {
				case 1:
					switch y {
						case 10: v = "ten"
						case 20: v = 20
					}
				case 2:
					v = #other
			}
			return v
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TFalse))
	a.That(u.Contains(TString))
	a.That(u.Contains(TNumber))
}

func TestSwitchBranchMergeInsideIfElse(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(p, x) {
			v = false
			if p {
				switch x {
					case 1: v = "a"
					case 2: v = 2
				}
			} else {
				v = #else
			}
			return v
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TString))
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TFalse))
}

func TestSwitchBranchMergeIntroducesNewLocal(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			switch x {
				case 1: v = "one"
				case 2: v = 2
				default: v = false
			}
			return v
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TString))
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TFalse))
}

func TestSwitchBranchMergeCompoundAssignmentInArm(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			v = 0
			switch x {
				case 1: v += 1
				case 2: v = "two"
			}
			return v
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TString))
}

func TestSwitchMergeArmReturn(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			v = false
			switch x {
				case 1: v = "one"; return v
				case 2: v = 2
			}
			return v
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TString)) // arm-1 return
	a.That(u.Contains(TNumber)) // post-switch read after arm-2
	a.That(u.Contains(TFalse))  // post-switch read after no-match path
}

func TestSwitchBranchMergeMultiExprCaseBody(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			v = false
			switch x {
				case 1, 2, 3: v = "small"
				default: v = "other"
			}
			return v
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestSwitchBranchMergeEmptyKeepsEntryType(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			v = 42
			switch x {
			}
			return v
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestSwitchMemberAssignmentBackfills(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Init(x) {
			switch x {
				case 1: .v = "str"
				case 2: .v = 2
			}
		}
	}`, "T")
	u, ok := env.Members["v"].(Union)
	a.That(ok)
	a.That(u.Contains(TString))
	a.That(u.Contains(TNumber))
}

func TestSwitchPostSwitchNarrowingSeesMergedUnion(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			v = false
			switch x {
				case 1: v = "one"
				case 2: v = 2
			}
			if Type(v) is "Number"
				return v
			return false
		}
	}`, "T")
	// then-branch returns v narrowed to Number; else returns false.
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TFalse))
}

func TestSwitchBranchMergeFourArmsAllUnion(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			switch x {
				case 1: v = 1
				case 2: v = "two"
				case 3: v = false
				default: v = #four
			}
			return v
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TString)) // covers both "two" and #four
	a.That(u.Contains(TFalse))
}

func TestSwitchScrutineeInArmHasEntryType(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			x = 5
			switch x {
				case 1: return x
			}
			return x
		}
	}`, "T")
	// Both returns produce TNumber from the entry binding.
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestParseAnnotationEmptyIsUnknownNoError(t *testing.T) {
	a := assert.T(t)
	ty, err := ParseTypeAnnotation("")
	a.That(err == nil)
	a.This(ty).Is(TUnknown)
}

func TestParseAnnotationLowercasePrimitives(t *testing.T) {
	a := assert.T(t)
	cases := map[string]DynType{
		"boolean":  TBoolean,
		"number":   TNumber,
		"string":   TString,
		"date":     TDate,
		"object":   TObject,
		"record":   TObject,
		"class":    TClass,
		"function": TFunction,
		"block":    TBlock,
		"sequence": TSequence,
		"void":     TVoid,
		"unknown":  TUnknown,
		"false":    TFalse,
		"true":     TTrue,
	}
	for input, want := range cases {
		ty, err := ParseTypeAnnotation(input)
		a.That(err == nil)
		a.This(ty).Is(want)
	}
}

// a capitalized builtin name is still that builtin, not a class of the same name
func TestParseAnnotationUppercaseBuiltin(t *testing.T) {
	a := assert.T(t)
	ty, err := ParseTypeAnnotation("Number")
	a.That(err == nil)
	a.This(ty).Is(TNumber)
	ty, err = ParseTypeAnnotation("String")
	a.That(err == nil)
	a.This(ty).Is(TString)
	ty, err = ParseTypeAnnotation("Object")
	a.That(err == nil)
	a.This(ty).Is(TObject)
}

// Parser: pipe-separated alternatives fold into a Union.
func TestParseAnnotationUnion(t *testing.T) {
	a := assert.T(t)
	ty, err := ParseTypeAnnotation("string|number")
	a.That(err == nil)
	u, ok := ty.(Union)
	a.That(ok)
	a.That(u.Contains(TString) && u.Contains(TNumber))
}

func TestParseAnnotationCustomClassIsNominal(t *testing.T) {
	a := assert.T(t)
	ty, err := ParseTypeAnnotation("Foo")
	a.That(err == nil)
	a.This(ty).Is(Instance{Class: "Foo"})
}

// Parser: union of primitive and class name folds correctly.
func TestParseAnnotationPrimitivePlusClassUnion(t *testing.T) {
	a := assert.T(t)
	ty, err := ParseTypeAnnotation("number|Foo")
	a.That(err == nil)
	// TNumber | Instance{Foo}
	u, ok := ty.(Union)
	a.That(ok)
	a.That(u.Contains(TNumber) && u.Contains(Instance{Class: "Foo"}))
}

func TestParseAnnotationLowercaseUnknownIsError(t *testing.T) {
	a := assert.T(t)
	_, err := ParseTypeAnnotation("nubmer")
	a.That(err != nil)
}

func TestAnnotationUnionTypo(t *testing.T) {
	a := assert.T(t)
	ty, err := ParseTypeAnnotation("string|nubmer")
	a.That(err != nil)
	// String survives. Union with TUnknown produces a dirty union per U().
	u, ok := ty.(Union)
	a.That(ok)
	a.That(u.IsDirty)
	a.That(u.Contains(TString))
}

func TestParseAnnotationEmptyAlternativeIsError(t *testing.T) {
	a := assert.T(t)
	_, err := ParseTypeAnnotation("string||number")
	a.That(err != nil)
}

func TestAnnotationParamOverridesDefaultValueType(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: string = 0) { return x }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestAnnotationParamSeedsType(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: number) { return x }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestAnnotationParamUnionNarrowsInGuard(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: string | number) {
			if Number?(x) { return x }
			return x
		}
	}`, "T")
	// then-arm returns TNumber; fall-through returns the original union.
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TString))
}

func TestAnnotationReturnPinsPublishedReturn(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: number) : string {
			return x
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestAnnotationReturnMatchesInferred(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() : number {
			return 42
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestAnnotationReturnVisibleToInternalCallers(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Helper(x) : string {
			return x
		}
		Caller() {
			return .Helper(1)
		}
	}`, "T")
	a.This(env.Returns["Caller"]).Is(TString)
}

func TestAnnotationClassNameReturnIsNominal(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Make() : Foo {
			return new Foo()
		}
	}`, "T")
	a.This(env.Returns["Make"]).Is(Instance{Class: "Foo"})
}

func TestAnnotationReturnsMap(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Declared() : number { return 1 }
		Inferred() { return "hi" }
	}`, "T")
	_, declaredOk := env.AnnotatedReturns["Declared"]
	a.That(declaredOk)
	_, inferredOk := env.AnnotatedReturns["Inferred"]
	a.That(!inferredOk)
}

func TestAnnotationInlineInitSeedsMember(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Init(.count: number) {}
	}`, "T")
	a.This(env.Members["count"]).Is(TNumber)
}

func TestAnnotationInlineInitKeepsBool(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(.flag: false) {}
	}`, "T")
	a.This(env.Members["flag"]).Is(TFalse)
}

func TestAnnotationMismatchDoesNotCrash(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: number) : string {
			return x + 1
		}
	}`, "T")
	// Pinned to TString despite body returning TNumber.
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestAnnotationReturnUnionWidens(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() : number | string {
			return 1
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TString))
}

func TestAnnotationReturnViaSuper(t *testing.T) {
	a := assert.T(t)
	// Two classes via lineage requires using TypeInfer + a resolver.
	parentSrc := `class {
		Helper() : string { return 1 }
	}`
	childSrc := `Parent {
		Caller() {
			return super.Helper()
		}
	}`
	resolver := mapResolver{
		"Parent": parentSrc,
	}
	_, env := TypeInfer("Child", childSrc, resolver)
	a.This(env.Returns["Caller"]).Is(TString)
}

// `x is true` narrows to TTrue. Mirror of TestNarrowIsFalseLiteral.
func TestNarrowIsTrueLiteral(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { if x is true { return x } }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TTrue)
}

func TestNarrowIsntFalseElseNarrowsToFalse(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Maybe() { if .pick is true { return 42 }; return false }
		pick: false
		Foo() {
			x = .Maybe()
			if x isnt false { return 0 }
			else { return x }
		}
	}`, "T")
	// then-arm returns 0 (TNumber); else-arm returns x narrowed to TFalse.
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TFalse))
}

func TestNarrowNestedIfComposes(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			if Number?(x) {
				if x is 5 { return x }
			}
			return false
		}
	}`, "T")
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TFalse))
}

func TestNarrowTypeStringEqBoolean(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { if Type(x) is "Boolean" { return x } }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TBoolean)
}

// ---- Ternary narrowing regressions ----------------------------------------

// String?-guarded ternary: T branch narrows to TString, so return is TString.
func TestNarrowTernaryTBranchIsNarrowedIdent(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { return String?(x) ? x : "fallback" }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

// Negated guard: F branch narrows to TString, so return is TString.
func TestNarrowTernaryFBranchIsNarrowedIdent(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { return not String?(x) ? "fallback" : x }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestNarrowNestedTernaryAllBranchesString(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x, y) { return String?(x) ? x : (String?(y) ? y : "default") }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestInitMemberCopyFromOtherMember(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Init() { .a = 5; .b = .a }
	}`, "T")
	a.This(env.Members["a"]).Is(TNumber)
	a.This(env.Members["b"]).Is(TNumber)
}

func TestMemberAssignedOnlyInOneBranch(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		v: 0
		Update(flag) {
			if flag { .v = "hi" }
		}
	}`, "T")
	u, ok := env.Members["v"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TString))
}

func TestInlineInitBoolDefaultLocalAssign(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(.flag = false) { x = flag; return x }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(U(TFalse, TUnknown))
}

func TestMemberReadBeforeInit(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Read() { return .x }
		Init() { .x = 42 }
	}`, "T")
	a.This(env.Returns["Read"]).Is(TNumber)
}

func TestForInObjectLiteralElementUnknown(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Iter() { for x in [1,2,3] { return x } }
	}`, "T")
	a.This(env.Returns["Iter"]).Is(TUnknown)
}

func TestForInObjectLocalElementUnknown(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Iter() {
			ob = [1, 2, 3]
			for x in ob { return x }
		}
	}`, "T")
	a.This(env.Returns["Iter"]).Is(TUnknown)
}

func TestForInBodyNarrowing(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Iter(s: string) {
			for c in s {
				if Number?(c) { return c }
				return c
			}
		}
	}`, "T")
	ret := env.Returns["Iter"]
	if u, ok := ret.(Union); ok {
		a.That(u.Contains(TString))
	} else {
		a.This(ret).Is(TString)
	}
}

func TestForInAssignmentVisibleAfterLoop(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Iter(s: string) {
			for c in s { v = 1 }
			return v
		}
	}`, "T")
	a.This(env.Returns["Iter"]).Is(TNumber)
}

func TestNarrowSwitchScrutineeIdentToFalseLiteral(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			switch x {
				case false: return x
			}
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TFalse)
}

func TestSwitchOnlyDefault(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			v = false
			switch x {
				default: v = 1
			}
			return v
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestSubclassOverridesAnnotation(t *testing.T) {
	a := assert.T(t)
	parentSrc := `class {
		Foo() : number { return 1 }
	}`
	childSrc := `Parent {
		Foo() : string { return "x" }
		Caller() { return .Foo() }
	}`
	resolver := mapResolver{"Parent": parentSrc}
	_, env := TypeInfer("Child", childSrc, resolver)
	a.This(env.Returns["Caller"]).Is(TString)
}

func TestSuperCallMissingMethodIsUnknown(t *testing.T) {
	a := assert.T(t)
	parentSrc := `class { Other() { return 1 } }`
	childSrc := `Parent {
		Foo() { return super.NotThere() }
	}`
	resolver := mapResolver{"Parent": parentSrc}
	_, env := TypeInfer("Child", childSrc, resolver)
	a.This(env.Returns["Foo"]).Is(TUnknown)
}

func TestSuperCallNotFlagged(t *testing.T) {
	a := assert.T(t)
	parentSrc := `class {
		Helper() : string { return "ok" }
	}`
	childSrc := `Parent {
		Caller() {
			return super.Helper()
		}
	}`
	resolver := mapResolver{"Parent": parentSrc}
	_, env := TypeInfer("Child", childSrc, resolver)
	a.This(env.Returns["Caller"]).Is(TString)
	for _, d := range diagsByMethod(env)["Caller"] {
		a.That(!strings.Contains(d.Msg, "no built-in method"))
	}
}

func TestParseAnnotationThreeWayUnion(t *testing.T) {
	a := assert.T(t)
	ty, err := ParseTypeAnnotation("number|string|boolean")
	a.That(err == nil)
	u, ok := ty.(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TString))
	a.That(u.Contains(TBoolean))
}

func TestParseAnnotationDuplicateCollapses(t *testing.T) {
	a := assert.T(t)
	ty, err := ParseTypeAnnotation("number|number")
	a.That(err == nil)
	a.This(ty).Is(TNumber)
}

func TestAnnotationOverridesCallerArgType(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Helper(x: number) { return x }
		Caller() { return .Helper("oops") }
	}`, "T")
	a.This(env.Returns["Helper"]).Is(TNumber)
	a.This(env.Returns["Caller"]).Is(TNumber)
}

func TestBlockLiteralIsBlockType(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			b = {|x| x + 1}
			return b
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TBlock)
}

func TestMemberObjectLiteralIsObject(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		ob: [1, 2, 3]
	}`, "T")
	a.This(env.Members["ob"]).Is(TObject)
}

func TestBuiltinObjectMethods(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Add", Sig: "(@args) :object"},
		Reg{Receiver: "object", Name: "Assocs", Sig: "(list = true, named = true) :object"},
		Reg{Receiver: "object", Name: "BinarySearch", Sig: "(value, block = false) :number"},
		Reg{Receiver: "object", Name: "CompareAndSet", Sig: "(member, newValue, oldValue = nil) :boolean"},
		Reg{Receiver: "object", Name: "Copy", Sig: "() :object"},
		Reg{Receiver: "sequence", Name: "Copy", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Delete", Sig: "(@args) :object"},
		Reg{Receiver: "object", Name: "DeleteIf", Sig: "(block) :object"},
		Reg{Receiver: "object", Name: "Erase", Sig: "(@args) :object"},
		Reg{Receiver: "object", Name: "Eval", Sig: "(@args) :unknown"},
		Reg{Receiver: "string", Name: "Eval", Sig: "() :unknown"},
		Reg{Receiver: "class", Name: "Eval", Sig: "(@args) :unknown"},
		Reg{Receiver: "object", Name: "Eval2", Sig: "(@args) :object"},
		Reg{Receiver: "string", Name: "Eval2", Sig: "() :object"},
		Reg{Receiver: "class", Name: "Eval2", Sig: "(@args) :object"},
		Reg{Receiver: "object", Name: "Find", Sig: "(value) :false|unknown"},
		Reg{Receiver: "string", Name: "Find", Sig: "(string :string, pos=0) :number"},
		Reg{Receiver: "object", Name: "FindAll", Sig: "(value) :object"},
		Reg{Receiver: "object", Name: "FindAllIf", Sig: "(block) :object"},
		Reg{Receiver: "object", Name: "GetDefault", Sig: "(member, block) :unknown"},
		Reg{Receiver: "class", Name: "GetDefault", Sig: "(member, block) :unknown"},
		Reg{Receiver: "object", Name: "Has?", Sig: "(value) :boolean"},
		Reg{Receiver: "string", Name: "Has?", Sig: "(string :string) :boolean"},
		Reg{Receiver: "object", Name: "Iter", Sig: "() :object"},
		Reg{Receiver: "sequence", Name: "Iter", Sig: "() :object"},
		Reg{Receiver: "string", Name: "Iter", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Join", Sig: "(separator :string ='') :string"},
		Reg{Receiver: "sequence", Name: "Join", Sig: "(separator :string ='') :string"},
		Reg{Kind: "free", Name: "Max", Sig: "(@args) :unknown"},
		Reg{Receiver: "object", Name: "Max", Sig: "() :unknown"},
		Reg{Receiver: "object", Name: "Member?", Sig: "(member) :boolean"},
		Reg{Receiver: "class", Name: "Member?", Sig: "(member) :boolean"},
		Reg{Kind: "static", Class: "Database", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Date", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Ftsearch", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "LruCache", Name: "Members", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Members", Sig: "(list = true, named = true) :object"},
		Reg{Kind: "static", Class: "OpenPGP", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "PdfEncrypt", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Random", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Thread", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "TypeChecker", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Zlib", Name: "Members", Sig: "() :object"},
		Reg{Receiver: "class", Name: "Members", Sig: "(all = false) :object"},
		Reg{Kind: "free", Name: "Min", Sig: "(@args) :unknown"},
		Reg{Receiver: "object", Name: "Min", Sig: "() :unknown"},
		Reg{Receiver: "object", Name: "PopFirst", Sig: "() :unknown"},
		Reg{Receiver: "object", Name: "PopLast", Sig: "() :unknown"},
		Reg{Receiver: "class", Name: "Readonly?", Sig: "() :boolean"},
		Reg{Receiver: "object", Name: "Readonly?", Sig: "() :boolean"},
		Reg{Receiver: "object", Name: "Reverse!", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Set_default", Sig: "(value=nil) :object"},
		Reg{Receiver: "object", Name: "Set_readonly", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "object", Name: "Sort!", Sig: "(block = false) :object"},
		Reg{Receiver: "object", Name: "Unique!", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Values", Sig: "(list = true, named = true) :object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
	 	obj: #(1,2,3,4)
		Foo() {
			.add        = obj.Add(56, 78)
			.assocs     = obj.Assocs()
			.bs         = obj.BinarySearch(2)
			.cas        = obj.CompareAndSet("k", 1)
			.copy       = obj.Copy()
			.delete     = obj.Delete(0)
			.delIf      = obj.DeleteIf({|x| x > 2})
			.erase      = obj.Erase(0)
			.eval       = obj.Eval()
			.eval2      = obj.Eval2()
			.find       = obj.Find(1)
			.findAll    = obj.FindAll(1)
			.findIf     = obj.FindAllIf({|x| x > 2})
			.getDef     = obj.GetDefault(0, false)
			.has        = obj.Has?(1)
			.iter       = obj.Iter()
			.join       = obj.Join(",")
			.max        = obj.Max()
			.memberQ    = obj.Member?("k")
			.members    = obj.Members()
			.min        = obj.Min()
			.popF       = obj.PopFirst()
			.popL       = obj.PopLast()
			.ro         = obj.Readonly?()
			.reverse    = obj.Reverse!()
			.setDef     = obj.Set_default(0)
			.setRO      = obj.Set_readonly()
			.size       = obj.Size()
			.sort       = obj.Sort!()
			.uniq       = obj.Unique!()
			.values     = obj.Values()
		}
	}`, "T")

	// Concrete primitive returns
	a.This(env.Members["add"]).Is(TObject)
	a.This(env.Members["assocs"]).Is(TObject)
	a.This(env.Members["bs"]).Is(TNumber)
	a.This(env.Members["cas"]).Is(TBoolean)
	a.This(env.Members["copy"]).Is(TObject)
	a.This(env.Members["delete"]).Is(TObject)
	a.This(env.Members["delIf"]).Is(TObject)
	a.This(env.Members["erase"]).Is(TObject)
	a.This(env.Members["eval2"]).Is(TObject)
	a.This(env.Members["findAll"]).Is(TObject)
	a.This(env.Members["findIf"]).Is(TObject)
	a.This(env.Members["has"]).Is(TBoolean)
	a.This(env.Members["iter"]).Is(TObject)
	a.This(env.Members["join"]).Is(TString)
	a.This(env.Members["memberQ"]).Is(TBoolean)
	a.This(env.Members["members"]).Is(TObject)
	a.This(env.Members["ro"]).Is(TBoolean)
	a.This(env.Members["reverse"]).Is(TObject)
	a.This(env.Members["setDef"]).Is(TObject)
	a.This(env.Members["setRO"]).Is(TObject)
	a.This(env.Members["size"]).Is(TNumber)
	a.This(env.Members["sort"]).Is(TObject)
	a.This(env.Members["uniq"]).Is(TObject)
	a.This(env.Members["values"]).Is(TObject)

	u, ok := env.Members["find"].(Union)
	a.That(ok)
	a.That(u.IsDirty)
	a.That(u.Contains(TFalse))

	for _, name := range []string{"eval", "getDef", "max", "min", "popF", "popL"} {
		if _, present := env.Members[name]; present {
			t.Errorf("%s: expected absent from Members for TUnknown return, got %v",
				name, env.Members[name])
		}
	}
}

func TestBuiltinObjectMethodsStdlib(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "AddMany!", Sig: "(value, n) :object"},
		Reg{Receiver: "object", Name: "AddTo", Sig: "(ob) :unknown"},
		Reg{Receiver: "object", Name: "AddUnique", Sig: "(value) :object"},
		Reg{Receiver: "object", Name: "Any?", Sig: "(block) :boolean"},
		Reg{Receiver: "object", Name: "Append", Sig: "(ob) :object"},
		Reg{Receiver: "object", Name: "BinarySearch?", Sig: "(value) :boolean"},
		Reg{Receiver: "object", Name: "Concat", Sig: "(@iterables) :object"},
		Reg{Receiver: "string", Name: "Count", Sig: "(string :string) :number"},
		Reg{Kind: "static", Class: "Thread", Name: "Count", Sig: "() :number"},
		Reg{Receiver: "object", Name: "Count", Sig: "(value = #(0)) :number"},
		Reg{Receiver: "object", Name: "CountIf", Sig: "(block) :number"},
		Reg{Receiver: "object", Name: "DeepCopy", Sig: "(nesting = 0) :object"},
		Reg{Receiver: "object", Name: "DeepReplaceReferences!", Sig: "(oldName, newName, nesting = 0) :unknown"},
		Reg{Receiver: "object", Name: "Difference", Sig: "(other) :object"},
		Reg{Receiver: "object", Name: "Disjoint?", Sig: "(other) :boolean"},
		Reg{Receiver: "object", Name: "Drop", Sig: "(n) :object"},
		Reg{Receiver: "object", Name: "DuplicateValues", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Each", Sig: "(block) :object"},
		Reg{Receiver: "object", Name: "Each2", Sig: "(block) :object"},
		Reg{Receiver: "object", Name: "Empty?", Sig: "() :boolean"},
		Reg{Receiver: "object", Name: "EqualSet?", Sig: "(that) :boolean"},
		Reg{Receiver: "object", Name: "Every?", Sig: "(block) :boolean"},
		Reg{Receiver: "string", Name: "Extract", Sig: "(pattern, part=false) :false|string"},
		Reg{Receiver: "object", Name: "Extract", Sig: "(member, x = #extract_no_default) :unknown"},
		Reg{Receiver: "object", Name: "Filter", Sig: "(block) :object"},
		Reg{Receiver: "object", Name: "FindIf", Sig: "(block) :false|unknown"},
		Reg{Receiver: "object", Name: "FindLastIf", Sig: "(block) :false|unknown"},
		Reg{Receiver: "object", Name: "FindOne", Sig: "(block) :false|unknown"},
		Reg{Receiver: "object", Name: "First", Sig: "() :unknown"},
		Reg{Receiver: "object", Name: "FlatMap", Sig: "(block) :object"},
		Reg{Receiver: "object", Name: "Flatten", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Fold", Sig: "(val, block) :unknown"},
		Reg{Receiver: "object", Name: "GetDifferences", Sig: "(ob) :object"},
		Reg{Receiver: "object", Name: "GetInit", Sig: "(member, block) :unknown"},
		Reg{Receiver: "object", Name: "Grep", Sig: "(regex, block = false) :object"},
		Reg{Receiver: "object", Name: "HasIf?", Sig: "(block) :boolean"},
		Reg{Receiver: "object", Name: "HasNamed?", Sig: "() :boolean"},
		Reg{Receiver: "object", Name: "HasNonEmptyMember?", Sig: "(members) :boolean"},
		Reg{Receiver: "object", Name: "Instantiate", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Intersect", Sig: "(other) :object"},
		Reg{Receiver: "object", Name: "Intersects?", Sig: "(other) :boolean"},
		Reg{Receiver: "object", Name: "JoinCSV", Sig: "(fields = false) :string"},
		Reg{Receiver: "object", Name: "Last", Sig: "() :unknown"},
		Reg{Receiver: "object", Name: "ListToMembers", Sig: "() :object"},
		Reg{Receiver: "object", Name: "ListToNamed", Sig: "(@fields) :object"},
		Reg{Receiver: "object", Name: "Map", Sig: "(block) :object"},
		Reg{Receiver: "string", Name: "Map", Sig: "(block) :string"},
		Reg{Receiver: "object", Name: "Map!", Sig: "(block) :object"},
		Reg{Receiver: "object", Name: "Map2", Sig: "(block) :object"},
		Reg{Receiver: "object", Name: "MapMembers", Sig: "(block) :object"},
		Reg{Receiver: "object", Name: "MaxWith", Sig: "(block) :unknown"},
		Reg{Kind: "static", Class: "Database", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Date", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Ftsearch", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "LruCache", Name: "Members", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Members", Sig: "(list = true, named = true) :object"},
		Reg{Kind: "static", Class: "OpenPGP", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "PdfEncrypt", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Random", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Thread", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "TypeChecker", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Zlib", Name: "Members", Sig: "() :object"},
		Reg{Receiver: "class", Name: "Members", Sig: "(all = false) :object"},
		Reg{Receiver: "object", Name: "MembersIf", Sig: "(block) :object"},
		Reg{Receiver: "object", Name: "Merge", Sig: "(ob) :object"},
		Reg{Receiver: "object", Name: "MergeNew", Sig: "(ob) :object"},
		Reg{Receiver: "object", Name: "MergeUnion", Sig: "(other) :object"},
		Reg{Receiver: "object", Name: "MinWith", Sig: "(block) :unknown"},
		Reg{Receiver: "object", Name: "NotEmpty?", Sig: "() :boolean"},
		Reg{Receiver: "object", Name: "Nth", Sig: "(n) :unknown"},
		Reg{Receiver: "object", Name: "Project", Sig: "(@fields) :object"},
		Reg{Receiver: "object", Name: "ProjectValues", Sig: "(@fields) :object"},
		Reg{Receiver: "object", Name: "RandVal", Sig: "() :unknown"},
		Reg{Receiver: "object", Name: "Reduce", Sig: "(block) :unknown"},
		Reg{Receiver: "object", Name: "Remove", Sig: "(@values) :object"},
		Reg{Receiver: "object", Name: "Remove1", Sig: "(x) :unknown"},
		Reg{Receiver: "object", Name: "RemoveIf", Sig: "(block) :object"},
		Reg{Receiver: "string", Name: "Replace", Sig: "(pattern, block = '', count = false) :string"},
		Reg{Receiver: "object", Name: "Replace", Sig: "(oldvalue, newvalue) :object"},
		Reg{Receiver: "date", Name: "Replace", Sig: "(@args) :date"},
		Reg{Receiver: "object", Name: "SafeMembers", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Shuffle!", Sig: "() :object"},
		Reg{Receiver: "object", Name: "SortWith!", Sig: "(block) :object"},
		Reg{Receiver: "object", Name: "Sorted?", Sig: "(block = false) :boolean"},
		Reg{Receiver: "object", Name: "Subset?", Sig: "(other) :boolean"},
		Reg{Receiver: "object", Name: "Sum", Sig: "() :number"},
		Reg{Receiver: "object", Name: "SumWith", Sig: "(block) :number"},
		Reg{Receiver: "object", Name: "Swap", Sig: "(idx1, idx2) :object"},
		Reg{Receiver: "object", Name: "Take", Sig: "(n) :object"},
		Reg{Receiver: "object", Name: "Trim!", Sig: "(@values) :object"},
		Reg{Receiver: "object", Name: "Union", Sig: "(other) :object"},
		Reg{Receiver: "object", Name: "UniqueValues", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Val_or_func", Sig: "(member) :unknown"},
		Reg{Receiver: "object", Name: "Without", Sig: "(@values) :object"},
		Reg{Receiver: "object", Name: "WithoutFields", Sig: "(@fields) :object"},
		Reg{Receiver: "object", Name: "Zip", Sig: "() :object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			obj = #(1,2,3,4)
			other = #(2,3,5)
			// :boolean
			.empty       = obj.Empty?()
			.notEmpty    = obj.NotEmpty?()
			.hasNamed    = obj.HasNamed?()
			.intersectsQ = obj.Intersects?(other)
			.disjointQ   = obj.Disjoint?(other)
			.subsetQ     = obj.Subset?(other)
			.everyQ      = obj.Every?({|x| x > 0})
			.anyQ        = obj.Any?({|x| x > 0})
			.hasIfQ      = obj.HasIf?({|x| x > 0})
			.sortedQ     = obj.Sorted?()
			.hasNonEmpty = obj.HasNonEmptyMember?("a")
			.binSearchQ  = obj.BinarySearch?(2)
			.equalSetQ   = obj.EqualSet?(other)

			// :object - mappers
			.mapBang     = obj.Map!({|x| x})
			.mapM        = obj.Map({|x| x})
			.map2        = obj.Map2({|x| x})
			.mapMembers  = obj.MapMembers({|x| x})
			.flatMap     = obj.FlatMap({|x| x})

			// :object - set ops
			.intersect   = obj.Intersect(other)
			.difference  = obj.Difference(other)
			.unionO      = obj.Union(other)
			.mergeU      = obj.MergeUnion(other)
			.dupVals     = obj.DuplicateValues()
			.uniqueVals  = obj.UniqueValues()

			// :object - selection & projection
			.project     = obj.Project("a")
			.projectVals = obj.ProjectValues("a")
			.listToNamed = obj.ListToNamed("a")
			.listToMems  = obj.ListToMembers()
			.flatten     = obj.Flatten()
			.membersIf   = obj.MembersIf({|x| x > 0})
			.safeMembers = obj.SafeMembers()
			.removeMany  = obj.Remove(1, 2)
			.removeIf    = obj.RemoveIf({|x| x > 0})
			.without     = obj.Without(1)
			.withoutF    = obj.WithoutFields("a")
			.filter      = obj.Filter({|x| x > 0})
			.zip         = obj.Zip()
			.take        = obj.Take(2)
			.drop        = obj.Drop(2)
			.concat      = obj.Concat(other)
			.instantiate = obj.Instantiate()
			.grep        = obj.Grep("p")
			.deepCopy    = obj.DeepCopy()
			.trimBang    = obj.Trim!(1)

			// :object - mutators / merges
			.merge       = obj.Merge(other)
			.mergeNew    = obj.MergeNew(other)
			.replace     = obj.Replace(1, 2)
			.swap        = obj.Swap(0, 1)
			.addUnique   = obj.AddUnique(99)
			.append      = obj.Append(other)
			.shuffleBang = obj.Shuffle!()
			.addManyBang = obj.AddMany!(0, 3)
			.each        = obj.Each({|x| x})
			.each2       = obj.Each2({|x| x})
			.sortWithB   = obj.SortWith!({|a, b| a < b})
			.getDiffs    = obj.GetDifferences(other)

			// :number
			.sum         = obj.Sum()
			.sumWith     = obj.SumWith({|x| x})
			.count       = obj.Count()
			.countIf     = obj.CountIf({|x| x > 0})

			// :string
			.joinCSV     = obj.JoinCSV()

			// :false|unknown -> dirty Union with TFalse
			.findOne     = obj.FindOne({|x| x > 0})
			.findIf      = obj.FindIf({|x| x > 0})
			.findLastIf  = obj.FindLastIf({|x| x > 0})

			// :unknown (absent from env.Members)
			.first       = obj.First()
			.last        = obj.Last()
			.minWith     = obj.MinWith({|a, b| a < b})
			.maxWith     = obj.MaxWith({|a, b| a < b})
			.valOrFunc   = obj.Val_or_func("k")
			.fold        = obj.Fold(0, {|a, x| a + x})
			.reduce      = obj.Reduce({|a, x| a + x})
			.remove1     = obj.Remove1(1)
			.extract     = obj.Extract("k")
			.getInit     = obj.GetInit("k", {|| 1})
			.nth         = obj.Nth(0)
			.addTo       = obj.AddTo(other)
			.deepReplace = obj.DeepReplaceReferences!("a", "b")
			.randVal     = obj.RandVal()
		}
	}`, "T")

	for _, name := range []string{
		"empty", "notEmpty", "hasNamed",
		"intersectsQ", "disjointQ", "subsetQ",
		"everyQ", "anyQ", "hasIfQ", "sortedQ",
		"hasNonEmpty", "binSearchQ", "equalSetQ",
	} {
		a.This(env.Members[name]).Is(TBoolean)
	}

	for _, name := range []string{
		"mapBang", "mapM", "map2", "mapMembers", "flatMap",
		"intersect", "difference", "unionO", "mergeU",
		"dupVals", "uniqueVals",
		"project", "projectVals", "listToNamed", "listToMems",
		"flatten", "membersIf", "safeMembers",
		"removeMany", "removeIf", "without", "withoutF",
		"filter", "zip", "take", "drop", "concat",
		"instantiate", "grep", "deepCopy", "trimBang",
		"merge", "mergeNew", "replace", "swap", "addUnique",
		"append", "shuffleBang", "addManyBang",
		"each", "each2", "sortWithB", "getDiffs",
	} {
		a.This(env.Members[name]).Is(TObject)
	}

	for _, name := range []string{"sum", "sumWith", "count", "countIf"} {
		a.This(env.Members[name]).Is(TNumber)
	}

	a.This(env.Members["joinCSV"]).Is(TString)

	for _, name := range []string{"findOne", "findIf", "findLastIf"} {
		u, ok := env.Members[name].(Union)
		if !ok {
			t.Errorf("%s: expected Union for false|unknown, got %T", name, env.Members[name])
			continue
		}
		a.That(u.IsDirty)
		a.That(u.Contains(TFalse))
	}

	for _, name := range []string{
		"first", "last", "minWith", "maxWith", "valOrFunc",
		"fold", "reduce", "remove1", "extract", "getInit",
		"nth", "addTo", "deepReplace", "randVal",
	} {
		if _, present := env.Members[name]; present {
			t.Errorf("%s: expected absent from Members for TUnknown return, got %v",
				name, env.Members[name])
		}
	}
}

func TestBuiltinStringMethods(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Alpha?", Sig: "() :boolean"},
		Reg{Receiver: "string", Name: "AlphaNum?", Sig: "() :boolean"},
		Reg{Receiver: "string", Name: "Asc", Sig: "() :number"},
		Reg{Receiver: "string", Name: "Count", Sig: "(string :string) :number"},
		Reg{Kind: "static", Class: "Thread", Name: "Count", Sig: "() :number"},
		Reg{Receiver: "object", Name: "Count", Sig: "(value = #(0)) :number"},
		Reg{Receiver: "string", Name: "Detab", Sig: "() :string"},
		Reg{Receiver: "string", Name: "Entab", Sig: "() :string"},
		Reg{Receiver: "object", Name: "Eval2", Sig: "(@args) :object"},
		Reg{Receiver: "string", Name: "Eval2", Sig: "() :object"},
		Reg{Receiver: "class", Name: "Eval2", Sig: "(@args) :object"},
		Reg{Receiver: "string", Name: "Extract", Sig: "(pattern, part=false) :false|string"},
		Reg{Receiver: "object", Name: "Extract", Sig: "(member, x = #extract_no_default) :unknown"},
		Reg{Receiver: "object", Name: "Find", Sig: "(value) :false|unknown"},
		Reg{Receiver: "string", Name: "Find", Sig: "(string :string, pos=0) :number"},
		Reg{Receiver: "string", Name: "Find1of", Sig: "(chars :string, pos=0) :number"},
		Reg{Receiver: "string", Name: "FindLast", Sig: "(string :string, pos=false) :false|number"},
		Reg{Receiver: "string", Name: "FindLast1of", Sig: "(chars :string, pos=false) :false|number"},
		Reg{Receiver: "string", Name: "FromHex", Sig: "() :string"},
		Reg{Receiver: "string", Name: "FromUtf8", Sig: "() :string"},
		Reg{Receiver: "object", Name: "Has?", Sig: "(value) :boolean"},
		Reg{Receiver: "string", Name: "Has?", Sig: "(string :string) :boolean"},
		Reg{Receiver: "object", Name: "Iter", Sig: "() :object"},
		Reg{Receiver: "sequence", Name: "Iter", Sig: "() :object"},
		Reg{Receiver: "string", Name: "Iter", Sig: "() :object"},
		Reg{Receiver: "string", Name: "Lower", Sig: "() :string"},
		Reg{Receiver: "string", Name: "Lower?", Sig: "() :boolean"},
		Reg{Receiver: "string", Name: "MapN", Sig: "(n :number, block) :string"},
		Reg{Receiver: "string", Name: "Match", Sig: "(pattern, pos=false, prev :boolean =false) :false|object"},
		Reg{Receiver: "string", Name: "NthLine", Sig: "(n) :string"},
		Reg{Receiver: "string", Name: "Number?", Sig: "() :boolean"},
		Reg{Kind: "free", Name: "Number?", Sig: "(value) :boolean"},
		Reg{Receiver: "string", Name: "Numeric?", Sig: "() :boolean"},
		Reg{Receiver: "string", Name: "Prefix?", Sig: "(string :string, pos=0) :boolean"},
		Reg{Receiver: "string", Name: "Repeat", Sig: "(count) :string"},
		Reg{Receiver: "string", Name: "Replace", Sig: "(pattern, block = '', count = false) :string"},
		Reg{Receiver: "object", Name: "Replace", Sig: "(oldvalue, newvalue) :object"},
		Reg{Receiver: "date", Name: "Replace", Sig: "(@args) :date"},
		Reg{Receiver: "string", Name: "Reverse", Sig: "() :string"},
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "string", Name: "Split", Sig: "(separator :string = '') :object"},
		Reg{Receiver: "string", Name: "Suffix?", Sig: "(string :string) :boolean"},
		Reg{Receiver: "string", Name: "ToHex", Sig: "() :string"},
		Reg{Receiver: "string", Name: "ToUtf8", Sig: "() :string"},
		Reg{Receiver: "string", Name: "Tr", Sig: "(from :string, to :string ='') :string"},
		Reg{Receiver: "string", Name: "Unescape", Sig: "() :string"},
		Reg{Receiver: "object", Name: "Union", Sig: "(other) :object"},
		Reg{Receiver: "string", Name: "Upper", Sig: "() :string"},
		Reg{Receiver: "string", Name: "Upper?", Sig: "() :boolean"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			s = "hello world"
			.alpha       = s.Alpha?()
			.alphaNum    = s.AlphaNum?()
			.asc         = s.Asc()
			.count       = s.Count("l")
			.detab       = s.Detab()
			.entab       = s.Entab()
			.eval2       = s.Eval2()
			.extract     = s.Extract("l")
			.find        = s.Find("l")
			.find1of     = s.Find1of("lo")
			.findLast    = s.FindLast("l")
			.findLast1of = s.FindLast1of("lo")
			.fromUtf8    = s.FromUtf8()
			.has         = s.Has?("l")
			.iter        = s.Iter()
			.lower       = s.Lower()
			.lowerQ      = s.Lower?()
			.mapN        = s.MapN(1, {|x| x})
			.match       = s.Match("l")
			.nthLine     = s.NthLine(0)
			.numberQ     = s.Number?()
			.numericQ    = s.Numeric?()
			.prefixQ     = s.Prefix?("he")
			.repeat      = s.Repeat(2)
			.replace     = s.Replace("l", "r")
			.reverse     = s.Reverse()
			.size        = s.Size()
			.split       = s.Split(" ")
			.suffixQ     = s.Suffix?("d")
			.toHex       = s.ToHex()
			.fromHex     = s.FromHex()
			.toUtf8      = s.ToUtf8()
			.tr          = s.Tr("a", "b")
			.unescape    = s.Unescape()
			.upper       = s.Upper()
			.upperQ      = s.Upper?()
		}
	}`, "T")

	for _, name := range []string{
		"alpha", "alphaNum", "has", "lowerQ", "numberQ", "numericQ",
		"prefixQ", "suffixQ", "upperQ",
	} {
		a.This(env.Members[name]).Is(TBoolean)
	}

	for _, name := range []string{
		"asc", "count", "find", "find1of", "size",
	} {
		a.This(env.Members[name]).Is(TNumber)
	}

	for _, name := range []string{
		"detab", "entab", "fromUtf8", "lower", "mapN", "nthLine",
		"repeat", "replace", "reverse", "toHex", "fromHex", "toUtf8",
		"tr", "unescape", "upper",
	} {
		a.This(env.Members[name]).Is(TString)
	}

	a.This(env.Members["eval2"]).Is(TObject)
	a.This(env.Members["iter"]).Is(TObject)
	a.This(env.Members["split"]).Is(TObject)

	// Extract returns `:false|string` -> Union{TFalse, TString}.
	uExtract, ok := env.Members["extract"].(Union)
	a.That(ok)
	a.That(uExtract.Contains(TFalse))
	a.That(uExtract.Contains(TString))

	// FindLast / FindLast1of return `:false|number` -> Union{TFalse, TNumber}.
	for _, name := range []string{"findLast", "findLast1of"} {
		u, ok := env.Members[name].(Union)
		if !ok {
			t.Errorf("%s: expected Union for false|number, got %T", name, env.Members[name])
			continue
		}
		a.That(u.Contains(TFalse))
		a.That(u.Contains(TNumber))
	}

	// Match returns `:false|object` -> Union{TFalse, TObject}.
	uMatch, ok := env.Members["match"].(Union)
	a.That(ok)
	a.That(uMatch.Contains(TFalse))
	a.That(uMatch.Contains(TObject))
}

func TestBuiltinStringMethodsStdlib(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "AfterFirst", Sig: "(delimiter) :string"},
		Reg{Receiver: "string", Name: "AfterLast", Sig: "(delimiter) :string"},
		Reg{Receiver: "string", Name: "As", Sig: "(e) :string"},
		Reg{Receiver: "string", Name: "Base64Decode", Sig: "() :string"},
		Reg{Receiver: "string", Name: "Base64Encode", Sig: "() :string"},
		Reg{Receiver: "string", Name: "BeforeFirst", Sig: "(delimiter) :string"},
		Reg{Receiver: "string", Name: "BeforeLast", Sig: "(delimiter) :string"},
		Reg{Receiver: "string", Name: "Blank?", Sig: "() :boolean"},
		Reg{Receiver: "string", Name: "Capitalize", Sig: "() :string"},
		Reg{Receiver: "string", Name: "CapitalizeWords", Sig: "(lower = true) :string"},
		Reg{Receiver: "string", Name: "Capitalized?", Sig: "() :boolean"},
		Reg{Receiver: "string", Name: "Center", Sig: "(minSize, char = ' ') :string"},
		Reg{Receiver: "string", Name: "ChangeEol", Sig: "(eol) :string"},
		Reg{Receiver: "string", Name: "Divide", Sig: "(n = 1) :object"},
		Reg{Receiver: "string", Name: "DynamicName?", Sig: "() :boolean"},
		Reg{Receiver: "string", Name: "Ellipsis", Sig: "(maxLength, atEnd = false) :string"},
		Reg{Receiver: "string", Name: "Escape", Sig: "() :string"},
		Reg{Receiver: "string", Name: "ExtractAll", Sig: "(pattern) :false|object"},
		Reg{Receiver: "string", Name: "FindRx", Sig: "(rx) :number"},
		Reg{Receiver: "string", Name: "FindRxLast", Sig: "(rx) :false|number"},
		Reg{Receiver: "string", Name: "FirstLine", Sig: "() :string"},
		Reg{Receiver: "string", Name: "ForEach1of", Sig: "(chars, block) :unknown"},
		Reg{Receiver: "string", Name: "ForEachMatch", Sig: "(pat, block) :unknown"},
		Reg{Receiver: "string", Name: "GlobalName?", Sig: "() :boolean"},
		Reg{Receiver: "string", Name: "Has1of?", Sig: "(chars) :boolean"},
		Reg{Receiver: "string", Name: "Identifier?", Sig: "() :boolean"},
		Reg{Receiver: "string", Name: "In?", Sig: "(x) :boolean"},
		Reg{Receiver: "string", Name: "LeftFill", Sig: "(minSize, char = ' ') :string"},
		Reg{Receiver: "string", Name: "LeftTrim", Sig: "(chars = ' \t\r\n') :string"},
		Reg{Receiver: "string", Name: "LineAtPosition", Sig: "(pos) :string"},
		Reg{Receiver: "string", Name: "LineCount", Sig: "() :number"},
		Reg{Receiver: "string", Name: "LineFromPosition", Sig: "(pos) :number"},
		Reg{Receiver: "string", Name: "Lines", Sig: "() :object"},
		Reg{Receiver: "string", Name: "LocalName?", Sig: "() :boolean"},
		Reg{Receiver: "object", Name: "Map", Sig: "(block) :object"},
		Reg{Receiver: "string", Name: "Map", Sig: "(block) :string"},
		Reg{Kind: "static", Class: "Database", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Date", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Ftsearch", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "LruCache", Name: "Members", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Members", Sig: "(list = true, named = true) :object"},
		Reg{Kind: "static", Class: "OpenPGP", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "PdfEncrypt", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Random", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Thread", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "TypeChecker", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Zlib", Name: "Members", Sig: "() :object"},
		Reg{Receiver: "class", Name: "Members", Sig: "(all = false) :object"},
		Reg{Receiver: "string", Name: "RandChar", Sig: "() :string"},
		Reg{Receiver: "string", Name: "RemoveBlankLines", Sig: "() :string"},
		Reg{Receiver: "string", Name: "RemovePrefix", Sig: "(prefix) :string"},
		Reg{Receiver: "string", Name: "RemoveSuffix", Sig: "(suffix) :string"},
		Reg{Receiver: "string", Name: "ReplaceSubstr", Sig: "(i, n, s) :string"},
		Reg{Receiver: "string", Name: "RightFill", Sig: "(minSize, char = ' ') :string"},
		Reg{Receiver: "string", Name: "RightTrim", Sig: "(chars = ' \t\r\n') :string"},
		Reg{Receiver: "string", Name: "SafeEval", Sig: "() :unknown"},
		Reg{Receiver: "number", Name: "SafeEval", Sig: "() :number"},
		Reg{Receiver: "string", Name: "Shuffle", Sig: "() :string"},
		Reg{Receiver: "string", Name: "SplitCSV", Sig: "(fields = false, string_vals = false) :object"},
		Reg{Receiver: "string", Name: "SplitFixedLength", Sig: "(map) :object"},
		Reg{Receiver: "string", Name: "SplitOnFirst", Sig: "(delimiter = ' ') :object"},
		Reg{Receiver: "string", Name: "SplitOnLast", Sig: "(delimiter = ' ') :object"},
		Reg{Receiver: "string", Name: "StartPositionOfLine", Sig: "(line) :number"},
		Reg{Receiver: "string", Name: "StripInvalidChars", Sig: "() :string"},
		Reg{Receiver: "string", Name: "Trim", Sig: "(chars = ' \t\r\n') :string"},
		Reg{Receiver: "string", Name: "TruncateLeftFill", Sig: "(size, char = ' ') :string"},
		Reg{Receiver: "string", Name: "TruncateRightFill", Sig: "(size, char = ' ') :string"},
		Reg{Receiver: "string", Name: "UnCapitalize", Sig: "() :string"},
		Reg{Receiver: "string", Name: "UniqueChars", Sig: "() :string"},
		Reg{Receiver: "string", Name: "White?", Sig: "() :boolean"},
		Reg{Receiver: "string", Name: "WrapLines", Sig: "(width) :object"},
		Reg{Receiver: "string", Name: "Xor", Sig: "(key) :string"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			s = "hello world"
			// :boolean
			.capitalizedQ = s.Capitalized?()
			.whiteQ       = s.White?()
			.blankQ       = s.Blank?()
			.inQ          = s.In?("hello world list")
			.globalQ      = s.GlobalName?()
			.localQ       = s.LocalName?()
			.identQ       = s.Identifier?()
			.dynamicQ     = s.DynamicName?()
			.has1ofQ      = s.Has1of?("lo")

			// :number
			.lineCount    = s.LineCount()
			.lineFromPos  = s.LineFromPosition(0)
			.startPosLine = s.StartPositionOfLine(0)
			.findRx       = s.FindRx("[a-z]+")

			// :string
			.asAs         = s.As("e")
			.firstLine    = s.FirstLine()
			.removeBlank  = s.RemoveBlankLines()
			.removePre    = s.RemovePrefix("he")
			.removeSuf    = s.RemoveSuffix("ld")
			.changeEol    = s.ChangeEol("\n")
			.lineAtPos    = s.LineAtPosition(0)
			.capitalize   = s.Capitalize()
			.unCapitalize = s.UnCapitalize()
			.capWords     = s.CapitalizeWords()
			.trim         = s.Trim()
			.leftTrim     = s.LeftTrim()
			.rightTrim    = s.RightTrim()
			.replaceSub   = s.ReplaceSubstr(0, 1, "x")
			.leftFill     = s.LeftFill(10)
			.truncLeftFil = s.TruncateLeftFill(10)
			.rightFill    = s.RightFill(10)
			.truncRightFi = s.TruncateRightFill(10)
			.center       = s.Center(10)
			.beforeFirst  = s.BeforeFirst("l")
			.afterFirst   = s.AfterFirst("l")
			.beforeLast   = s.BeforeLast("l")
			.afterLast    = s.AfterLast("l")
			.base64Enc    = s.Base64Encode()
			.base64Dec    = s.Base64Decode()
			.xor          = s.Xor("k")
			.mapS         = s.Map({|x| x})
			.escape       = s.Escape()
			.ellipsis     = s.Ellipsis(5)
			.uniqueChars  = s.UniqueChars()
			.stripInvalid = s.StripInvalidChars()
			.randChar     = s.RandChar()
			.shuffle      = s.Shuffle()

			// :object
			.lines        = s.Lines()
			.splitCSV     = s.SplitCSV()
			.splitFixed   = s.SplitFixedLength(#())
			.splitOnFirst = s.SplitOnFirst()
			.splitOnLast  = s.SplitOnLast()
			.wrapLines    = s.WrapLines(80)
			.divide       = s.Divide()

			// :false|number
			.findRxLast   = s.FindRxLast("[a-z]+")

			// :false|object
			.extractAll   = s.ExtractAll("[a-z]+")

			// :unknown (absent from env.Members)
			.forEachMatch = s.ForEachMatch("[a-z]+", {|m| m})
			.forEach1of   = s.ForEach1of("lo", {|i| i})
			.safeEval     = s.SafeEval()
		}
	}`, "T")

	for _, name := range []string{
		"capitalizedQ", "whiteQ", "blankQ", "inQ",
		"globalQ", "localQ", "identQ", "dynamicQ", "has1ofQ",
	} {
		a.This(env.Members[name]).Is(TBoolean)
	}

	for _, name := range []string{
		"lineCount", "lineFromPos", "startPosLine", "findRx",
	} {
		a.This(env.Members[name]).Is(TNumber)
	}

	for _, name := range []string{
		"asAs", "firstLine", "removeBlank", "removePre", "removeSuf",
		"changeEol", "lineAtPos", "capitalize", "unCapitalize", "capWords",
		"trim", "leftTrim", "rightTrim", "replaceSub",
		"leftFill", "truncLeftFil", "rightFill", "truncRightFi", "center",
		"beforeFirst", "afterFirst", "beforeLast", "afterLast",
		"base64Enc", "base64Dec", "xor", "mapS", "escape", "ellipsis",
		"uniqueChars", "stripInvalid", "randChar", "shuffle",
	} {
		a.This(env.Members[name]).Is(TString)
	}

	for _, name := range []string{
		"lines", "splitCSV", "splitFixed", "splitOnFirst", "splitOnLast",
		"wrapLines", "divide",
	} {
		a.This(env.Members[name]).Is(TObject)
	}

	// FindRxLast returns :false|number.
	uFindRxLast, ok := env.Members["findRxLast"].(Union)
	a.That(ok)
	a.That(uFindRxLast.Contains(TFalse))
	a.That(uFindRxLast.Contains(TNumber))

	// ExtractAll returns :false|object.
	uExtractAll, ok := env.Members["extractAll"].(Union)
	a.That(ok)
	a.That(uExtractAll.Contains(TFalse))
	a.That(uExtractAll.Contains(TObject))

	for _, name := range []string{"forEachMatch", "forEach1of", "safeEval"} {
		if _, present := env.Members[name]; present {
			t.Errorf("%s: expected absent from Members for TUnknown return, got %v",
				name, env.Members[name])
		}
	}
}

func TestBuiltinNumberMethods(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "number", Name: "ACos", Sig: "() :number"},
		Reg{Receiver: "number", Name: "ASin", Sig: "() :number"},
		Reg{Receiver: "number", Name: "ATan", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Binary", Sig: "() :string"},
		Reg{Receiver: "number", Name: "Chr", Sig: "() :string"},
		Reg{Receiver: "number", Name: "Cos", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Exp", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Format", Sig: "(mask :string) :string"},
		Reg{Receiver: "date", Name: "Format", Sig: "(format) :string"},
		Reg{Receiver: "number", Name: "Frac", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Hex", Sig: "() :string"},
		Reg{Receiver: "number", Name: "Int", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Log", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Log10", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Log2", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Pow", Sig: "(number) :number"},
		Reg{Receiver: "number", Name: "Round", Sig: "(number) :number"},
		Reg{Receiver: "number", Name: "RoundDown", Sig: "(number) :number"},
		Reg{Receiver: "number", Name: "RoundUp", Sig: "(number) :number"},
		Reg{Receiver: "number", Name: "Sin", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Sqrt", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Tan", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			n = 5
			.acos       = n.ACos()
			.asin       = n.ASin()
			.atan       = n.ATan()
			.binary     = n.Binary()
			.chr        = n.Chr()
			.cos        = n.Cos()
			.exp        = n.Exp()
			.format     = n.Format("0.00")
			.frac       = n.Frac()
			.hex        = n.Hex()
			.int        = n.Int()
			.log        = n.Log()
			.log2       = n.Log2()
			.log10      = n.Log10()
			.pow        = n.Pow(2)
			.round      = n.Round(2)
			.roundUp    = n.RoundUp(2)
			.roundDown  = n.RoundDown(2)
			.sin        = n.Sin()
			.sqrt       = n.Sqrt()
			.tan        = n.Tan()
		}
	}`, "T")

	for _, name := range []string{
		"acos", "asin", "atan", "cos", "exp", "frac", "int",
		"log", "log2", "log10", "pow",
		"round", "roundUp", "roundDown", "sin", "sqrt", "tan",
	} {
		a.This(env.Members[name]).Is(TNumber)
	}

	for _, name := range []string{"binary", "chr", "format", "hex"} {
		a.This(env.Members[name]).Is(TString)
	}
}

func TestBuiltinNumberMethodsStdlib(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "number", Name: "Abs", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Ceiling", Sig: "() :number"},
		Reg{Receiver: "number", Name: "DecimalToPercent", Sig: "(round = 0) :number"},
		Reg{Receiver: "number", Name: "DollarFormat", Sig: "(mask) :string"},
		Reg{Receiver: "number", Name: "EnFrancais", Sig: "() :string"},
		Reg{Receiver: "number", Name: "EuroFormat", Sig: "(mask) :string"},
		Reg{Receiver: "number", Name: "Even?", Sig: "() :boolean"},
		Reg{Receiver: "number", Name: "Factorial", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Floor", Sig: "() :number"},
		Reg{Receiver: "number", Name: "FracDigits", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Gb", Sig: "() :number"},
		Reg{Receiver: "number", Name: "HoursInMinutes", Sig: "() :number"},
		Reg{Receiver: "number", Name: "HoursInSeconds", Sig: "() :number"},
		Reg{Receiver: "number", Name: "InchesInCanvasUnit", Sig: "() :number"},
		Reg{Receiver: "number", Name: "InchesInTwips", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Int?", Sig: "() :boolean"},
		Reg{Receiver: "number", Name: "IntDigits", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Kb", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Mb", Sig: "() :number"},
		Reg{Kind: "static", Class: "Database", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Date", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Ftsearch", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "LruCache", Name: "Members", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Members", Sig: "(list = true, named = true) :object"},
		Reg{Kind: "static", Class: "OpenPGP", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "PdfEncrypt", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Random", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Thread", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "TypeChecker", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Zlib", Name: "Members", Sig: "() :object"},
		Reg{Receiver: "class", Name: "Members", Sig: "(all = false) :object"},
		Reg{Receiver: "number", Name: "MinutesInMs", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Odd?", Sig: "() :boolean"},
		Reg{Receiver: "number", Name: "Of", Sig: "(block) :object"},
		Reg{Receiver: "number", Name: "OfStr", Sig: "(block) :string"},
		Reg{Receiver: "number", Name: "Pad", Sig: "(minSize, char = '0') :string"},
		Reg{Receiver: "number", Name: "PercentToDecimal", Sig: "() :number"},
		Reg{Receiver: "number", Name: "RoundToNearest", Sig: "(nearest = 1) :number"},
		Reg{Receiver: "number", Name: "RoundToPrecision", Sig: "(p) :number"},
		Reg{Receiver: "string", Name: "SafeEval", Sig: "() :unknown"},
		Reg{Receiver: "number", Name: "SafeEval", Sig: "() :number"},
		Reg{Receiver: "number", Name: "SecondsInHours", Sig: "() :number"},
		Reg{Receiver: "number", Name: "SecondsInMinutes", Sig: "() :number"},
		Reg{Receiver: "number", Name: "SecondsInMs", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Sign", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Times", Sig: "(block) :unknown"},
		Reg{Receiver: "number", Name: "ToRGB", Sig: "() :object"},
		Reg{Receiver: "number", Name: "ToWordSpanish", Sig: "() :string"},
		Reg{Receiver: "number", Name: "ToWords", Sig: "() :string"},
		Reg{Receiver: "number", Name: "ToWordsDutch", Sig: "() :string"},
		Reg{Receiver: "number", Name: "ToWordsFrench", Sig: "() :string"},
		Reg{Receiver: "number", Name: "ToWordsItalian", Sig: "() :string"},
		Reg{Receiver: "number", Name: "ToWordsSimple", Sig: "() :string"},
		Reg{Receiver: "number", Name: "TwipsInInch", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			n = 5
			// :number
			.ceiling       = n.Ceiling()
			.floor         = n.Floor()
			.abs           = n.Abs()
			.pctToDec      = n.PercentToDecimal()
			.decToPct      = n.DecimalToPercent()
			.factorial     = n.Factorial()
			.intDigits     = n.IntDigits()
			.fracDigits    = n.FracDigits()
			.sign          = n.Sign()
			.roundToPrec   = n.RoundToPrecision(2)
			.roundToNear   = n.RoundToNearest()
			.minutesInMs   = n.MinutesInMs()
			.secondsInMs   = n.SecondsInMs()
			.secondsInHrs  = n.SecondsInHours()
			.hoursInSecs   = n.HoursInSeconds()
			.hoursInMins   = n.HoursInMinutes()
			.secondsInMins = n.SecondsInMinutes()
			.inchesInTwips = n.InchesInTwips()
			.twipsInInch   = n.TwipsInInch()
			.inchesInCanv  = n.InchesInCanvasUnit()
			.kb            = n.Kb()
			.mb            = n.Mb()
			.gb            = n.Gb()
			.safeEval      = n.SafeEval()

			// :string
			.toWords       = n.ToWords()
			.enFrancais    = n.EnFrancais()
			.toWordsFr     = n.ToWordsFrench()
			.toWordsNl     = n.ToWordsDutch()
			.toWordSp      = n.ToWordSpanish()
			.toWordsIt     = n.ToWordsItalian()
			.toWordsSimple = n.ToWordsSimple()
			.pad           = n.Pad(3)
			.euroFormat    = n.EuroFormat("0.00")
			.dollarFormat  = n.DollarFormat("0.00")
			.ofStr         = n.OfStr({ "x" })

			// :boolean
			.intQ          = n.Int?()
			.evenQ         = n.Even?()
			.oddQ          = n.Odd?()

			// :object
			.toRGB         = n.ToRGB()
			.of            = n.Of({ 1 })

			// :unknown (absent from env.Members)
			.times         = n.Times({ 1 })
		}
	}`, "T")

	for _, name := range []string{
		"ceiling", "floor", "abs", "pctToDec", "decToPct", "factorial",
		"intDigits", "fracDigits", "sign", "roundToPrec", "roundToNear",
		"minutesInMs", "secondsInMs", "secondsInHrs", "hoursInSecs",
		"hoursInMins", "secondsInMins", "inchesInTwips", "twipsInInch",
		"inchesInCanv", "kb", "mb", "gb", "safeEval",
	} {
		a.This(env.Members[name]).Is(TNumber)
	}

	for _, name := range []string{
		"toWords", "enFrancais", "toWordsFr", "toWordsNl", "toWordSp",
		"toWordsIt", "toWordsSimple", "pad", "euroFormat", "dollarFormat",
		"ofStr",
	} {
		a.This(env.Members[name]).Is(TString)
	}

	for _, name := range []string{"intQ", "evenQ", "oddQ"} {
		a.This(env.Members[name]).Is(TBoolean)
	}

	for _, name := range []string{"toRGB", "of"} {
		a.This(env.Members[name]).Is(TObject)
	}

	if _, present := env.Members["times"]; present {
		t.Errorf("times: expected absent from Members for TUnknown return, got %v",
			env.Members["times"])
	}
}

func TestBuiltinDateMethods(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "date", Name: "Day", Sig: "() :number"},
		Reg{Receiver: "date", Name: "FormatEn", Sig: "(format :string) :string"},
		Reg{Receiver: "date", Name: "GetLocalGMTBias", Sig: "() :number"},
		Reg{Receiver: "date", Name: "Hour", Sig: "() :number"},
		Reg{Receiver: "date", Name: "Millisecond", Sig: "() :number"},
		Reg{Receiver: "date", Name: "MinusDays", Sig: "(date :date) :number"},
		Reg{Receiver: "date", Name: "MinusSeconds", Sig: "(date :date) :number"},
		Reg{Receiver: "date", Name: "Minute", Sig: "() :number"},
		Reg{Receiver: "date", Name: "Month", Sig: "() :number"},
		Reg{Receiver: "date", Name: "Plus", Sig: "(years=0, months=0, days=0, hours=0, minutes=0, seconds=0, milliseconds=0) :date"},
		Reg{Receiver: "date", Name: "Second", Sig: "() :number"},
		Reg{Receiver: "date", Name: "WeekDay", Sig: "(firstDay='Sun') :number"},
		Reg{Receiver: "date", Name: "Year", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(d: date) {
			.minusDays  = d.MinusDays(d)
			.minusSecs  = d.MinusSeconds(d)
			.formatEn   = d.FormatEn("yyyy-MM-dd")
			.gmtBias    = d.GetLocalGMTBias()
			.plus       = d.Plus(days: 1)
			.weekDay    = d.WeekDay()
			.year       = d.Year()
			.month      = d.Month()
			.day        = d.Day()
			.hour       = d.Hour()
			.minute     = d.Minute()
			.second     = d.Second()
			.ms         = d.Millisecond()
		}
	}`, "T")

	for _, name := range []string{
		"minusDays", "minusSecs", "gmtBias", "weekDay",
		"year", "month", "day", "hour", "minute", "second", "ms",
	} {
		a.This(env.Members[name]).Is(TNumber)
	}

	a.This(env.Members["formatEn"]).Is(TString)
	a.This(env.Members["plus"]).Is(TDate)
}

func TestBuiltinDateMethodsStdlib(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "date", Name: "DayOfYear", Sig: "() :number"},
		Reg{Receiver: "date", Name: "EndOfDay", Sig: "() :date"},
		Reg{Receiver: "date", Name: "EndOfMonth", Sig: "() :date"},
		Reg{Receiver: "date", Name: "EndOfMonthDay", Sig: "() :number"},
		Reg{Receiver: "number", Name: "Format", Sig: "(mask :string) :string"},
		Reg{Receiver: "date", Name: "Format", Sig: "(format) :string"},
		Reg{Receiver: "date", Name: "FromUnix", Sig: "(s) :date"},
		Reg{Receiver: "date", Name: "GMTime", Sig: "() :date"},
		Reg{Receiver: "date", Name: "GMTimeToLocal", Sig: "() :date"},
		Reg{Receiver: "date", Name: "InternetFormat", Sig: "() :string"},
		Reg{Receiver: "date", Name: "IsoWeekDay", Sig: "() :string"},
		Reg{Receiver: "date", Name: "LongDate", Sig: "() :string"},
		Reg{Receiver: "date", Name: "LongDateTime", Sig: "() :string"},
		Reg{Receiver: "date", Name: "Minus", Sig: "(@args) :date"},
		Reg{Receiver: "date", Name: "MinusHours", Sig: "(date) :number"},
		Reg{Receiver: "date", Name: "MinusMinutes", Sig: "(date) :number"},
		Reg{Receiver: "date", Name: "MinusMonths", Sig: "(date) :number"},
		Reg{Receiver: "date", Name: "NoTime", Sig: "() :date"},
		Reg{Receiver: "date", Name: "NoTime?", Sig: "() :boolean"},
		Reg{Receiver: "date", Name: "Quarter", Sig: "() :number"},
		Reg{Receiver: "string", Name: "Replace", Sig: "(pattern, block = '', count = false) :string"},
		Reg{Receiver: "object", Name: "Replace", Sig: "(oldvalue, newvalue) :object"},
		Reg{Receiver: "date", Name: "Replace", Sig: "(@args) :date"},
		Reg{Receiver: "date", Name: "ShortDate", Sig: "() :string"},
		Reg{Receiver: "date", Name: "ShortDateTime", Sig: "() :string"},
		Reg{Receiver: "date", Name: "ShortDateTimeSec", Sig: "() :string"},
		Reg{Receiver: "date", Name: "StartOfDay", Sig: "() :date"},
		Reg{Receiver: "date", Name: "StartOfYear", Sig: "() :date"},
		Reg{Receiver: "date", Name: "StdShortDate", Sig: "() :string"},
		Reg{Receiver: "date", Name: "StdShortDateTime", Sig: "() :string"},
		Reg{Receiver: "date", Name: "StdShortDateTimeSec", Sig: "() :string"},
		Reg{Receiver: "date", Name: "Time", Sig: "() :string"},
		Reg{Receiver: "date", Name: "UTC", Sig: "(gmtBias = false) :string"},
		Reg{Kind: "free", Name: "UnixTime", Sig: "() :number"},
		Reg{Receiver: "date", Name: "UnixTime", Sig: "() :number"},
		Reg{Receiver: "date", Name: "WeekNumber", Sig: "(firstday = 'sun') :number"},
		Reg{Receiver: "date", Name: "WeekStart", Sig: "(firstday = 'sun') :date"},
		Reg{Receiver: "date", Name: "isLeapYear?", Sig: "(year) :boolean"},
		Reg{Receiver: "date", Name: "translate", Sig: "(s) :string"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(d: date) {
			// :string
			.format         = d.Format("#####")
			.translate      = d.translate("Monday")
			.shortDate      = d.ShortDate()
			.longDate       = d.LongDate()
			.time           = d.Time()
			.shortDateTime  = d.ShortDateTime()
			.shortDTSec     = d.ShortDateTimeSec()
			.stdShortDate   = d.StdShortDate()
			.stdShortDT     = d.StdShortDateTime()
			.stdShortDTSec  = d.StdShortDateTimeSec()
			.longDateTime   = d.LongDateTime()
			.isoWeekDay     = d.IsoWeekDay()
			.utc            = d.UTC()
			.internetFormat = d.InternetFormat()

			// :number
			.dayOfYear      = d.DayOfYear()
			.minusMinutes   = d.MinusMinutes(d)
			.minusHours     = d.MinusHours(d)
			.minusMonths    = d.MinusMonths(d)
			.endOfMonthDay  = d.EndOfMonthDay()
			.weekNumber     = d.WeekNumber()
			.quarter        = d.Quarter()
			.unixTime       = d.UnixTime()

			// :boolean
			.noTimeQ        = d.NoTime?()
			.isLeapYearQ    = d.isLeapYear?(2024)

			// :date
			.noTime         = d.NoTime()
			.startOfDay     = d.StartOfDay()
			.endOfDay       = d.EndOfDay()
			.replace        = d.Replace(day: 1)
			.startOfYear    = d.StartOfYear()
			.endOfMonth     = d.EndOfMonth()
			.weekStart      = d.WeekStart()
			.gmtime         = d.GMTime()
			.gmtimeToLocal  = d.GMTimeToLocal()
			.minus          = d.Minus(days: 1)
			.fromUnix       = d.FromUnix(0)
		}
	}`, "T")

	for _, name := range []string{
		"format", "translate", "shortDate", "longDate", "time",
		"shortDateTime", "shortDTSec", "stdShortDate", "stdShortDT",
		"stdShortDTSec", "longDateTime", "isoWeekDay", "utc", "internetFormat",
	} {
		a.This(env.Members[name]).Is(TString)
	}

	for _, name := range []string{
		"dayOfYear", "minusMinutes", "minusHours", "minusMonths",
		"endOfMonthDay", "weekNumber", "quarter", "unixTime",
	} {
		a.This(env.Members[name]).Is(TNumber)
	}

	for _, name := range []string{"noTimeQ", "isLeapYearQ"} {
		a.This(env.Members[name]).Is(TBoolean)
	}

	for _, name := range []string{
		"noTime", "startOfDay", "endOfDay", "replace", "startOfYear",
		"endOfMonth", "weekStart", "gmtime", "gmtimeToLocal", "minus", "fromUnix",
	} {
		a.This(env.Members[name]).Is(TDate)
	}
}

func TestBuiltinDateStatics(t *testing.T) {
	withSigs(t,
		Reg{Kind: "static", Class: "Date", Name: "Begin", Sig: "() :date"},
		Reg{Kind: "static", Class: "Date", Name: "End", Sig: "() :date"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			.begin = Date.Begin()
			.end   = Date.End()
		}
	}`, "T")

	a.This(env.Members["begin"]).Is(TDate)
	a.This(env.Members["end"]).Is(TDate)

	for _, d := range methodDiags(env, "Foo") {
		if strings.Contains(d.Msg, "no built-in method") {
			t.Errorf("unexpected no-such-method diagnostic: %v", d)
		}
	}
}

func TestInfoWindowControlDefaultRefinement(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(autoClose = false)
			{
			if Number?(autoClose)
				Delay(autoClose.SecondsInMs(), .Destroy)
			}
		}`, "T")
	errCount := 0
	for _, d := range *env.Diagnostics {
		if d.Method == "New" && d.Severity == SeverityError {
			errCount++
		}
	}
	a.This(errCount).Is(0)
}

func TestEnhancedButtonComponentExample(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		svgImages: false
		setSvgImageColor(color, imageEl)
			{
			if color is 'black'
				return

			if .svgImages isnt false and .svgImages.Member?(color)
				imageEl.src = .svgImages[color]
			}
		clear()
			{
			.svgImages = Object()
			}
	}`, "T")
	for _, d := range methodDiags(env, "setSvgImageColor") {
		if d.Severity == SeverityError {
			t.Errorf("unexpected error: %v", d.Msg)
		}
	}
	a.This(env.Members["svgImages"]).Is(U(TFalse, TObject))
}

func TestLowercaseMemberFalseLiteralKeepsTFalse(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class { flag: false }`, "T")
	a.This(env.Members["flag"]).Is(TFalse)
}

func TestUppercaseMemberFalseLiteralStillWidens(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class { Flag: false }`, "T")
	a.This(env.Members["Flag"]).Is(U(TFalse, TUnknown))
}

func TestLowercaseMemberFalseUnion(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		data: false
		Reset() { .data = Object() }
	}`, "T")
	a.This(env.Members["data"]).Is(U(TFalse, TObject))
}

func TestUppercaseMemberFalseUnion(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Data: false
		Reset() { .Data = Object() }
	}`, "T")
	a.This(env.Members["Data"]).Is(markDirty(U(TFalse, TObject)))
}

func TestNarrowMemberAndChainSiblings(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		data: false
		Reset() { .data = Object() }
		Lookup(k) {
			if .data isnt false and .data.Member?(k)
				return .data[k]
			return false
		}
	}`, "T")
	for _, d := range methodDiags(env, "Lookup") {
		if d.Severity == SeverityError {
			t.Errorf("unexpected error: %v", d.Msg)
		}
	}
	a.That(env.Returns["Lookup"] != nil)
}

func TestNarrowMemberObjectAndChain(t *testing.T) {
	_, env := runPasses(`class {
		x: 0
		Reset() { .x = Object() }
		Lookup(k) {
			if Object?(.x) and .x.Member?(k)
				return .x[k]
			return false
		}
	}`, "T")
	for _, d := range methodDiags(env, "Lookup") {
		if d.Severity == SeverityError {
			t.Errorf("unexpected error: %v", d.Msg)
		}
	}
}

// ---- Assert-based member ground truths (AssertMemberPass) -----------------

func TestAssert_NewInlineParamsBecomeMembers(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(.x, .y) { Assert(Number?(x) and Number?(y)) }
		Sum() { return .x + .y }
	}`, "T")
	a.This(env.Members["x"]).Is(TNumber)
	a.This(env.Members["y"]).Is(TNumber)
	// .x + .y is well typed under the ground truth
	a.This(len(methodDiags(env, "Sum"))).Is(0)
	a.This(env.Returns["Sum"]).Is(TNumber)
}

func TestAssert_DirectMemberPredicate(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New() { Assert(String?(.name)) }
	}`, "T")
	a.This(env.Members["name"]).Is(TString)
}

// Type(x) is "..." form is understood the same as the Foo? predicate form.
func TestAssert_TypeStringForm(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(.x) { Assert(Type(x) is "Number") }
	}`, "T")
	a.This(env.Members["x"]).Is(TNumber)
}

func TestAssert_UsedAsWrongTypeFlags(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(.x) { Assert(Number?(x)) }
		Bad() { return .x $ "tail" }
	}`, "T")
	ds := methodDiags(env, "Bad")
	found := false
	for _, d := range ds {
		if d.Severity == SeverityError && strings.Contains(d.Msg, "operator") {
			found = true
		}
	}
	a.That(found)
}

func TestAssert_AssignmentConflictFlags(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(.x) { Assert(Number?(x)) }
		Reset() { .x = "nope" }
	}`, "T")
	ds := methodDiags(env, "Reset")
	found := false
	for _, d := range ds {
		if d.Severity == SeverityError &&
			strings.Contains(d.Msg, "ground truth") &&
			strings.Contains(d.Msg, "New") {
			found = true
		}
	}
	a.That(found)
}

// Assigning a compatible type to an asserted member is fine.
func TestAssert_CompatibleAssignmentOk(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(.x) { Assert(Number?(x)) }
		Reset() { .x = 0 }
	}`, "T")
	a.This(len(methodDiags(env, "Reset"))).Is(0)
	a.This(env.Members["x"]).Is(TNumber)
}

func TestAssert_ConflictingAssertsError(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New() { Assert(Number?(.x)) }
		Other() { Assert(String?(.x)) }
	}`, "T")
	found := false
	for _, d := range *env.Diagnostics {
		if d.Severity == SeverityError && strings.Contains(d.Msg, "conflicting Assert") {
			found = true
		}
	}
	a.That(found)
}

// Re-asserting the same type in two methods is consistent - no conflict.
func TestAssert_ConsistentAssertsNoConflict(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New() { Assert(Number?(.x)) }
		Check() { Assert(Number?(.x)) }
	}`, "T")
	a.This(env.Members["x"]).Is(TNumber)
	for _, d := range *env.Diagnostics {
		a.That(!strings.Contains(d.Msg, "conflicting Assert"))
	}
}

func TestAssert_PlainLocalIsNotMember(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(p) { Assert(Number?(p)) }
	}`, "T")
	_, ok := env.Members["p"]
	a.That(!ok)
}

func TestAssert_NamedArgFormIgnored(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(.x) { Assert(.x is: 5) }
	}`, "T")
	_, ok := env.Members["x"]
	a.That(!ok)
}

func TestAssert_AndCompositionPinsAll(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(.x, .y, .z) { Assert(Number?(x) and Number?(y) and String?(z)) }
	}`, "T")
	a.This(env.Members["x"]).Is(TNumber)
	a.This(env.Members["y"]).Is(TNumber)
	a.This(env.Members["z"]).Is(TString)
}

func TestAssert_OrCompositionPinsUnion(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(.x) { Assert(Number?(x) or String?(x)) }
	}`, "T")
	u, ok := env.Members["x"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.Contains(TString))
}

func TestAssert_UnionGroundTruthStillChecks(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(.x) { Assert(Number?(x) or String?(x)) }
		Bad() { return .x $ 1 }
	}`, "T")
	found := false
	for _, d := range methodDiags(env, "Bad") {
		if d.Severity == SeverityError && strings.Contains(d.Msg, "operator") {
			found = true
		}
	}
	a.That(found)
}

func TestAssert_PublicMemberNotPinned(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New() { Assert(Number?(.Width)) }
	}`, "T")
	_, ok := env.Members["Width"]
	a.That(!ok)
}

func TestAssert_PublicInlineInitNotPinned(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(.Width) { Assert(Number?(.Width)) }
	}`, "T")
	_, okUpper := env.Members["Width"]
	_, okLower := env.Members["width"]
	a.That(!okUpper)
	a.That(!okLower)
}

func TestAssert_PrivateMemberStillPinned(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New() { Assert(Number?(.count)) }
	}`, "T")
	a.This(env.Members["count"]).Is(TNumber)
}

// A class with no asserts is unaffected - the pass is a no-op.
func TestAssert_NoAssertsNoOp(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(.x = 0) {}
		Bump() { .x += 1 }
	}`, "T")
	a.This(env.Members["x"]).Is(TNumber)
	a.This(len(methodDiags(env, "Bump"))).Is(0)
}

// P1: New assigns the member directly and unconditionally.
func TestCtorDirectInitDemotesMember(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Count", Sig: "(string :string) :number"},
		Reg{Kind: "static", Class: "Thread", Name: "Count", Sig: "() :number"},
		Reg{Receiver: "object", Name: "Count", Sig: "(value = #(0)) :number"},
		Reg{Kind: "free", Name: "Object", Sig: "(@args) :object"},
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		data: false
		New() { .data = Object() }
		Count() { return .data.Size() }
	}`, "T")
	a.This(errorCount(env, "Count")).Is(0)
	a.This(env.Returns["Count"]).Is(TNumber)
}

func TestCtorHelperInitDemotesMember(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Count", Sig: "(string :string) :number"},
		Reg{Kind: "static", Class: "Thread", Name: "Count", Sig: "() :number"},
		Reg{Receiver: "object", Name: "Count", Sig: "(value = #(0)) :number"},
		Reg{Kind: "free", Name: "Object", Sig: "(@args) :object"},
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		data: false
		New() { .fill() }
		fill() { .data = Object() }
		Count() { return .data.Size() }
	}`, "T")
	a.This(errorCount(env, "Count")).Is(0)
	a.This(env.Returns["Count"]).Is(TNumber)
}

func TestCtorInitIfUnsetDemotesMember(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Copy", Sig: "() :object"},
		Reg{Receiver: "sequence", Name: "Copy", Sig: "() :object"},
		Reg{Kind: "free", Name: "Object", Sig: "(@args) :object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		columns: false
		New() { .SetColumns(Object()) }
		SetColumns(cols) {
			if .columns is cols
				return
			.columns = cols.Copy()
		}
		GetColumns() { return .columns.Copy() }
	}`, "T")
	a.This(errorCount(env, "GetColumns")).Is(0)
	a.This(env.Returns["GetColumns"]).Is(TObject)
}

func TestCtorDefaultIdiomArgDemotesMember(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Count", Sig: "(string :string) :number"},
		Reg{Kind: "static", Class: "Thread", Name: "Count", Sig: "() :number"},
		Reg{Receiver: "object", Name: "Count", Sig: "(value = #(0)) :number"},
		Reg{Kind: "free", Name: "Object", Sig: "(@args) :object"},
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		data: false
		New(data = false) { .Set(data isnt false ? data : Object()) }
		Set(data) {
			if .data is data
				return
			.data = Object()
		}
		Count() { return .data.Size() }
	}`, "T")
	a.This(errorCount(env, "Count")).Is(0)
	a.This(env.Returns["Count"]).Is(TNumber)
}

func TestCtorCompoundGuardDemotesMember(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Copy", Sig: "() :object"},
		Reg{Receiver: "sequence", Name: "Copy", Sig: "() :object"},
		Reg{Kind: "free", Name: "Object", Sig: "(@args) :object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		columns: false
		New() { .SetColumns(Object()) }
		SetColumns(cols, reset = false) {
			if .columns is cols and not reset
				return
			.columns = cols.Copy()
		}
		GetColumns() { return .columns.Copy() }
	}`, "T")
	a.This(errorCount(env, "GetColumns")).Is(0)
	a.This(env.Returns["GetColumns"]).Is(TObject)
}

func TestCtorConditionalInitNotDemoted(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Count", Sig: "(string :string) :number"},
		Reg{Kind: "static", Class: "Thread", Name: "Count", Sig: "() :number"},
		Reg{Receiver: "object", Name: "Count", Sig: "(value = #(0)) :number"},
		Reg{Kind: "free", Name: "Object", Sig: "(@args) :object"},
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		data: false
		New(flag) {
			if flag is 1
				.data = Object()
		}
		Count() { return .data.Size() }
	}`, "T")
	a.That(errorCount(env, "Count") >= 1)
}

func TestCtorReassignedFalseNotDemoted(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "record", Name: "Clear", Sig: "() :void"},
		Reg{Receiver: "string", Name: "Count", Sig: "(string :string) :number"},
		Reg{Kind: "static", Class: "Thread", Name: "Count", Sig: "() :number"},
		Reg{Receiver: "object", Name: "Count", Sig: "(value = #(0)) :number"},
		Reg{Kind: "free", Name: "Object", Sig: "(@args) :object"},
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		data: false
		New() { .data = Object() }
		Clear() { .data = false }
		Count() { return .data.Size() }
	}`, "T")
	a.That(errorCount(env, "Count") >= 1)
}

func TestCtorPublicMemberNotDemoted(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Data: false
		New() { .Data = Object() }
		Count() { return .Data.Size() }
	}`, "T")
	_, demoted := env.PostCtorMembers["Data"]
	a.That(!demoted)
}

func TestCtorNewKeepsSeedForUseBeforeInit(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Count", Sig: "(string :string) :number"},
		Reg{Kind: "static", Class: "Thread", Name: "Count", Sig: "() :number"},
		Reg{Receiver: "object", Name: "Count", Sig: "(value = #(0)) :number"},
		Reg{Kind: "free", Name: "Object", Sig: "(@args) :object"},
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		data: false
		New() {
			.data.Size()
			.data = Object()
		}
		Count() { return .data.Size() }
	}`, "T")
	// New reads .data before initializing it -> still errors (seed False|Object).
	a.That(errorCount(env, "New") >= 1)
	// Count runs post-construction -> demoted, clean.
	a.This(errorCount(env, "Count")).Is(0)
}

func TestCtorNonSentinelUnaffected(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		count: 0
		New() { .count = 5 }
		Get() { return .count }
	}`, "T")
	a.This(env.Returns["Get"]).Is(TNumber)
	a.This(errorCount(env, "Get")).Is(0)
}

func TestCtorRawMemberTypeRetainsSeed(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		data: false
		New() { .data = Object() }
		Count() { return .data.Size() }
	}`, "T")
	got := env.Members["data"]
	u, ok := got.(Union)
	a.That(ok)
	a.That(u.Contains(TFalse) && u.Contains(TObject))
}

func TestCtorPublishesSeedAndInstanceReturnViews(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: false
		New() { .x = "initialised" }
		Foo() { return .x }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
	a.That(isUnionOf(env.PreCtorReturns["Foo"], TFalse, TString))
}

func TestCtorSeedViewPropagatesThroughThisCall(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: false
		New() { .x = "initialised" }
		Foo() { return .x }
		Bar() { return .Foo() }
	}`, "T")
	a.This(env.Returns["Bar"]).Is(TString)
	a.That(isUnionOf(env.PreCtorReturns["Bar"], TFalse, TString))
}

func TestCtorNoDemotionLeavesSeedViewEmpty(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		count: 0
		New() { .count = 5 }
		Get() { return .count }
	}`, "T")
	a.This(len(env.PreCtorReturns)).Is(0)
}

func TestClassRefCallResolvesToSeedView(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithClassViews(`class {
		Get() { return Widget.Foo() }
	}`,
		map[string]map[string]DynType{"Widget": {"Foo": TString}},
		map[string]map[string]DynType{"Widget": {"Foo": U(TFalse, TString)}},
	)
	a.That(isUnionOf(env.Returns["Get"], TFalse, TString))
}

func TestInstanceCallResolvesToInstanceView(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithClassViews(`class {
		Get() { return (new Widget()).Foo() }
	}`,
		map[string]map[string]DynType{"Widget": {"Foo": TString}},
		map[string]map[string]DynType{"Widget": {"Foo": U(TFalse, TString)}},
	)
	a.This(env.Returns["Get"]).Is(TString)
}

func TestClassRefSeedViewViolatesStringContract(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithClassViews(`class {
		Get() : string { return Widget.Foo() }
	}`,
		map[string]map[string]DynType{"Widget": {"Foo": TString}},
		map[string]map[string]DynType{"Widget": {"Foo": U(TFalse, TString)}},
	)
	a.That(errorCount(env, "Get") >= 1)
}

func TestInstanceCallSatisfiesStringContract(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithClassViews(`class {
		Get() : string { return (new Widget()).Foo() }
	}`,
		map[string]map[string]DynType{"Widget": {"Foo": TString}},
		map[string]map[string]DynType{"Widget": {"Foo": U(TFalse, TString)}},
	)
	a.This(errorCount(env, "Get")).Is(0)
}

func TestClassRefInstanceFallback(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithClassViews(`class {
		Get() { return Widget.Value() }
	}`,
		map[string]map[string]DynType{"Widget": {"Value": TNumber}},
		map[string]map[string]DynType{},
	)
	a.This(env.Returns["Get"]).Is(TNumber)
}

func TestGuessDisagree(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Ensure(extra) {
			return extra.Map({ it }).Join(', ')
		}
	}`, "T")
	for _, d := range methodDiags(env, "Ensure") {
		a.That(d.Severity != SeverityError ||
			!strings.Contains(d.Msg, "at least one path"))
	}
}

func TestGuessReplace_NoDownstreamUnionError(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Clean(cmdline) {
			return cmdline.Replace('a', 'b').Trim()
		}
	}`, "T")
	for _, d := range methodDiags(env, "Clean") {
		a.That(d.Severity != SeverityError ||
			!strings.Contains(d.Msg, "at least one path"))
	}
}

func TestGuessDirtySentinel(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		shortest_key: false
		New(k) { .shortest_key = k }
		Key(record) {
			return .shortest_key.Split(',').Map({ it }).Join('x')
		}
	}`, "T")
	for _, d := range methodDiags(env, "Key") {
		a.That(d.Severity != SeverityError ||
			!strings.Contains(d.Msg, "at least one path"))
	}
}

func TestGuessDisagreeReturnIsDirty(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Map", Sig: "(block) :object"},
		Reg{Receiver: "string", Name: "Map", Sig: "(block) :string"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Pick(extra) { return extra.Map({ it }) }
	}`, "T")
	u, ok := env.Returns["Pick"].(Union)
	a.That(ok)
	a.That(u.IsDirty)
	a.That(u.Contains(TString) && u.Contains(TObject))
}

func TestGuessDisagree_StillWarns(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Join", Sig: "(separator :string ='') :string"},
		Reg{Receiver: "sequence", Name: "Join", Sig: "(separator :string ='') :string"},
		Reg{Receiver: "object", Name: "Map", Sig: "(block) :object"},
		Reg{Receiver: "string", Name: "Map", Sig: "(block) :string"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Ensure(extra) { return extra.Map({ it }).Join(', ') }
	}`, "T")
	found := false
	for _, d := range methodDiags(env, "Ensure") {
		if d.Severity == SeverityWarning && strings.Contains(d.Msg, "ambiguous overloads") {
			found = true
		}
	}
	a.That(found)
}

func TestGenuineUnion_StillErrors(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Join", Sig: "(separator :string ='') :string"},
		Reg{Receiver: "sequence", Name: "Join", Sig: "(separator :string ='') :string"},
		Reg{Kind: "free", Name: "Use", Sig: "(library :string) :boolean"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Use(c) {
			x = c is 1 ? "str" : #()
			x.Join(', ')
		}
	}`, "T")
	found := false
	for _, d := range methodDiags(env, "Use") {
		if d.Severity == SeverityError && strings.Contains(d.Msg, "at least one path") {
			found = true
		}
	}
	a.That(found)
}

func TestGuessSingleReturnStaysClean(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Sort!", Sig: "(block = false) :object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Pick(x) { return x.Sort!() }
	}`, "T")
	a.This(env.Returns["Pick"]).Is(TObject)
}

func TestGlobalStaticCall_NoBuiltinReceiverGuess(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		amazonDate() {
			return Date(AmazonAWS.GMTime().AfterFirst(', ').BeforeFirst(' GMT'))
		}
	}`, "AmazonV4Signing")
	for _, d := range methodDiags(env, "amazonDate") {
		a.That(d.Severity != SeverityError ||
			!strings.Contains(d.Msg, "not applicable to receiver of type Date"))
	}
}

func TestClassReceiverUniversalMethods(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		probe(rpt) {
			return Class?(rpt) and rpt.Member?('Export') and rpt.GetDefault('X', false)
		}
	}`, "T")
	for _, d := range methodDiags(env, "probe") {
		a.That(d.Severity != SeverityError ||
			!strings.Contains(d.Msg, "not applicable to receiver of type Class"))
	}
}

func TestSentinelMember_GuardedLazyInit_NoError(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		items: false
		AddItem(x) {
			if .items is false
				.items = Object()
			.items.Add(x)
		}
	}`, "T")
	a.That(!hasNotApplicableError(env, "AddItem"))
}

func TestSentinelAssignUseAcrossCall(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		items: false
		Build(other) {
			.items = Object()
			other.Size()
			.items.Add(1)
		}
	}`, "T")
	a.That(!hasNotApplicableError(env, "Build"))
}

func TestSentinelBoundedCall(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		items: false
		Reset() { .items = false }
		Build(other) {
			.items = Object()
			other.Size()
			.items.Add(1)
		}
	}`, "T")
	a.That(!hasNotApplicableError(env, "Build"))
}

func TestSentinelWritingCall(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Add", Sig: "(@args) :object"},
		Reg{Kind: "free", Name: "Object", Sig: "(@args) :object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		items: false
		reset() { .items = false }
		Build() {
			.items = Object()
			.reset()
			.items.Add(1)
		}
	}`, "T")
	a.That(hasNotApplicableError(env, "Build"))
}

func TestSentinelMember_OpaqueThisCallInvalidates(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Add", Sig: "(@args) :object"},
		Reg{Kind: "free", Name: "Object", Sig: "(@args) :object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		items: false
		Reset() { .items = false }
		Build() {
			.items = Object()
			.Send("Changed")
			.items.Add(1)
		}
	}`, "T")
	a.That(hasNotApplicableError(env, "Build"))
}

func TestSentinelMember_BlockArgWriteInvalidates(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Add", Sig: "(@args) :object"},
		Reg{Receiver: "object", Name: "Each", Sig: "(block) :object"},
		Reg{Kind: "free", Name: "Object", Sig: "(@args) :object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		items: false
		Reset() { .items = false }
		Build(ob) {
			.items = Object()
			ob.Each({ |x| .items = false })
			.items.Add(1)
		}
	}`, "T")
	a.That(hasNotApplicableError(env, "Build"))
}

func TestSentinelMember_UseBeforeInit_StillErrors(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Add", Sig: "(@args) :object"},
		Reg{Kind: "free", Name: "Object", Sig: "(@args) :object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		buf: false
		Go() {
			.buf.Add(1)
			.buf = Object()
		}
	}`, "T")
	a.That(hasNotApplicableError(env, "Go"))
}

func TestSentinelSetElsewhere(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Add", Sig: "(@args) :object"},
		Reg{Receiver: "object", Name: "Append", Sig: "(ob) :object"},
		Reg{Receiver: "object", Name: "Copy", Sig: "() :object"},
		Reg{Receiver: "sequence", Name: "Copy", Sig: "() :object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		columns: false
		setup(c) { .columns = c.Copy() }
		Append() { .columns.Add(1) }
	}`, "T")
	a.That(hasNotApplicableError(env, "Append"))
}

func TestSentinelPublicWidened(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Items: false
		AddItem(x) {
			if .Items is false
				.Items = Object()
			.Items.Add(x)
		}
	}`, "T")
	for _, d := range methodDiags(env, "AddItem") {
		a.That(!strings.Contains(d.Msg, "True"))
	}
}

func TestSentinelBoundedCallUnion(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		items: false
		alt() { .items = "hello" }
		Build(other) {
			.items = Object()
			other.Size()
			.items.Add(1)
		}
	}`, "T")
	a.That(!hasNotApplicableError(env, "Build"))
}

func TestBoolSentinel_ParamGuardedMethod_NoError(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(.text = false) {
			}
		Width() {
			return .text is false ? 0 : .text.Size()
		}
	}`, "T")
	a.That(!hasNotApplicableError(env, "Width"))
}

func TestBoolSentinel_PublicMemberGuarded_NoError(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Flag: false
		Check() {
			if .Flag isnt false
				return .Flag.Member?(1)
			return false
		}
	}`, "T")
	a.That(!hasNotApplicableError(env, "Check"))
}

func TestBoolSentinel_ArithmeticStillErrors(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(.pad = false) {
			}
		Calc() {
			return .pad + 1
		}
	}`, "T")
	a.That(hasOperatorError(env, "Calc"))
}

func TestNewConstructsInstance(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Make() { return new Point(1, 2) }
	}`, "T")
	a.This(env.Returns["Make"]).Is(Instance{Class: "Point"})
}

// the no-arg / no-parens form parses the same way and still constructs.
func TestNewNoArgsConstructsInstance(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Make() { return new Point }
	}`, "T")
	a.This(env.Returns["Make"]).Is(Instance{Class: "Point"})
}

func TestInstanceTypeString(t *testing.T) {
	a := assert.T(t)
	a.This(fmt.Sprint(Instance{Class: "Point"})).Is("Point")
}

func TestInstanceMethodResolvesViaRegistry(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithClasses(`class {
		Get() {
			c = new Counter()
			return c.Value()
		}
	}`, map[string]map[string]DynType{
		"Counter": {"Value": TNumber},
	})
	a.This(env.Returns["Get"]).Is(TNumber)
}

// directly off the constructor, no intermediate local.
func TestInstanceMethodOnFreshConstruction(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithClasses(`class {
		Get() { return (new Counter()).Value() }
	}`, map[string]map[string]DynType{
		"Counter": {"Value": TString},
	})
	a.This(env.Returns["Get"]).Is(TString)
}

func TestInstanceMethodUnknownWithoutRegistry(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Get() {
			c = new Counter()
			return c.Value()
		}
	}`, "T")
	a.This(env.Returns["Get"]).Is(TUnknown)
}

func TestStaticMemberReadResolves(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithStatics(`class {
		Foo() { return Ipsum.I }
	}`,
		map[string]map[string]DynType{"Ipsum": {}},
		map[string]map[string]DynType{"Ipsum": {"I": TString}},
		map[string]string{"Ipsum": ""})
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestStaticMemberReadIsExact(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithStatics(`class {
		Foo() { return Ref.Flag }
	}`,
		map[string]map[string]DynType{"Ref": {}},
		map[string]map[string]DynType{"Ref": {"Flag": TFalse}},
		map[string]string{"Ref": ""})
	a.This(env.Returns["Foo"]).Is(TFalse)
}

func TestStaticMemberSpecificGetter(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithStatics(`class {
		Foo() { return Ref.Num }
	}`,
		map[string]map[string]DynType{"Ref": {"Getter_Num": TNumber}},
		map[string]map[string]DynType{"Ref": {}},
		map[string]string{"Ref": ""})
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestStaticMemberGenericGetter(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithStatics(`class {
		Foo() { return Ref.Anything }
	}`,
		map[string]map[string]DynType{"Ref": {"Getter_": TString}},
		map[string]map[string]DynType{"Ref": {}},
		map[string]string{"Ref": ""})
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestStaticMemberDataBeatsGetter(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithStatics(`class {
		Foo() { return Ref.I }
	}`,
		map[string]map[string]DynType{"Ref": {"Getter_": TNumber}},
		map[string]map[string]DynType{"Ref": {"I": TString}},
		map[string]string{"Ref": ""})
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestStaticMemberInheritedFromBase(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithStatics(`class {
		Foo() { return Ref.I }
	}`,
		map[string]map[string]DynType{"Ref": {}, "Base": {}},
		map[string]map[string]DynType{"Ref": {}, "Base": {"I": TString}},
		map[string]string{"Ref": "Base", "Base": ""})
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestStaticMemberUnknownClassSilent(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() { return Whatever.X }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TUnknown)
	a.That(!hasMemberNotFound(env))
}

func TestStaticMemberNotFoundDiagnostic(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithStatics(`class {
		Foo() { return Ref.Nope }
	}`,
		map[string]map[string]DynType{"Ref": {}},
		map[string]map[string]DynType{"Ref": {"I": TString}},
		map[string]string{"Ref": ""})
	a.That(hasMemberNotFound(env))
	a.This(env.Returns["Foo"]).Is(TUnknown)
}

// a generic getter accepts any name, so "not found" is unprovable - silent.
func TestStaticMemberNotFoundSuppressedByGetter(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithStatics(`class {
		Foo() { return Ref.Nope }
	}`,
		map[string]map[string]DynType{"Ref": {"Getter_": TString}},
		map[string]map[string]DynType{"Ref": {}},
		map[string]string{"Ref": ""})
	a.That(!hasMemberNotFound(env))
}

func TestStaticMemberUnknownBase(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithStatics(`class {
		Foo() { return Ref.Nope }
	}`,
		map[string]map[string]DynType{"Ref": {}},
		map[string]map[string]DynType{"Ref": {}},
		map[string]string{"Ref": "MysteryBase"})
	a.That(!hasMemberNotFound(env))
}

func TestStaticMemberPrivatizedAccessResolves(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithStatics(`class {
		Foo() { return SvcCore.SvcCore_svcHooks }
	}`,
		map[string]map[string]DynType{"SvcCore": {}},
		map[string]map[string]DynType{"SvcCore": {"svcHooks": TObject}},
		map[string]string{"SvcCore": ""})
	a.This(env.Returns["Foo"]).Is(TObject)
	a.That(!hasMemberNotFound(env))
}

func TestStaticMemberPrivatizedAccessNotFlagged(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithStatics(`class {
		Foo() { return SvcCore.SvcCore_svcHooks }
	}`,
		map[string]map[string]DynType{"SvcCore": {}},
		map[string]map[string]DynType{"SvcCore": {}},
		map[string]string{"SvcCore": ""})
	a.That(!hasMemberNotFound(env))
}

func TestStaticMemberPrivatizedAccessAsCallArg(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithStatics(`class {
		Foo() { .SpyOn(SvcCore.SvcCore_svcHooks) }
	}`,
		map[string]map[string]DynType{"SvcCore": {}},
		map[string]map[string]DynType{"SvcCore": {}},
		map[string]string{"SvcCore": ""})
	a.That(!hasMemberNotFound(env))
}

func TestStaticMethodCallNotFlagged(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithStatics(`class {
		Foo() { return Ref.Compute() }
	}`,
		map[string]map[string]DynType{"Ref": {"Compute": TNumber}},
		map[string]map[string]DynType{"Ref": {}},
		map[string]string{"Ref": ""})
	a.That(!hasMemberNotFound(env))
	a.This(env.Returns["Foo"]).Is(TNumber)
}

// the annotation rule: a name is nominal because it is not a builtin, not
// because it is capitalized.
func TestAnnotationCaseInsensitive(t *testing.T) {
	a := assert.T(t)
	lower, err := ParseTypeAnnotation("string")
	a.That(err == nil)
	a.This(lower).Is(TString)
	upper, err := ParseTypeAnnotation("String")
	a.That(err == nil)
	a.This(upper).Is(TString)
	nominal, err := ParseTypeAnnotation("Foo")
	a.That(err == nil)
	a.This(nominal).Is(Instance{Class: "Foo"})
}

// Foo|String unions a nominal class type with the string primitive.
func TestAnnotationNominalUnion(t *testing.T) {
	a := assert.T(t)
	ty, err := ParseTypeAnnotation("Foo|String")
	a.That(err == nil)
	u, ok := ty.(Union)
	a.That(ok)
	a.That(u.Contains(Instance{Class: "Foo"}) && u.Contains(TString))
}

func TestInstanceUnknownMethodStaysUnknown(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithClasses(`class {
		Get() {
			c = new Counter()
			return c.Missing()
		}
	}`, map[string]map[string]DynType{
		"Counter": {"Value": TNumber},
	})
	a.This(env.Returns["Get"]).Is(TUnknown)
}

func TestComputedMethodCall_DirtyUnion(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Msg(args) {
			msg = args[0]
			if .Method?(msg)
				return this[msg](args)
			return 0
		}
	}`, "T")
	u, ok := env.Returns["Msg"].(Union)
	a.That(ok)
	a.That(u.IsDirty)           // the computed call can't be seen -> `?`
	a.That(u.Contains(TNumber)) // the `return 0` arm survives
}

func TestComputedMethodCall_SpreadArgs_DirtyUnion(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Msg(args) {
			msg = args[0]
			if .Method?(msg)
				return this[msg](@+1 args)
			return 0
		}
	}`, "T")
	u, ok := env.Returns["Msg"].(Union)
	a.That(ok)
	a.That(u.IsDirty)
	a.That(u.Contains(TNumber))
}

func TestComputedMethodCall_SoleReturnUnknown(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Call(target, msg) {
			return target[msg](msg)
		}
	}`, "T")
	a.This(env.Returns["Call"]).Is(TUnknown)
}

func TestLiteralMethodCall_ResolvesClean(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Bar(num) {
			return num - 1
		}
		Use() {
			return .Bar(1)
		}
	}`, "T")
	a.This(env.Returns["Bar"]).Is(TNumber)
	a.This(env.Returns["Use"]).Is(TNumber)
}

// omitted arg -> secs is its `false` default -> `if secs is false` true -> Number.
func TestCtxReturn_OmittedDefaultPrunesToNumber(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Pick(secs = false) {
			if secs is false
				return 1
			return "x"
		}
		Use() { return .Pick() }
	}`, "T")
	a.This(env.Returns["Use"]).Is(TNumber)
}

// non-false arg -> secs is Number -> `secs is false` false -> String. Like Timer(secs: 1).
func TestCtxReturn_ProvidedArgPrunesToString(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Pick(secs = false) {
			if secs is false
				return 1
			return "x"
		}
		Use() { return .Pick(secs: 1) }
	}`, "T")
	a.This(env.Returns["Use"]).Is(TString)
}

// call-site-insensitive sub-case: a local with a known value makes a guard dead.
func TestCtxReturn_LocalValuePrunesDeadArm(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			x = 1
			if x is false
				return "a"
			return 2
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestCtxReturn_MethodTypeStaysUnion(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Pick(secs = false) {
			if secs is false
				return 1
			return "x"
		}
	}`, "T")
	u, ok := env.Returns["Pick"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber) && u.Contains(TString))
}

func TestFuncRef_BareCallResolvesReturn(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithRefSrcs(`class {
		Foo() { return NameSplit("a b") }
	}`, map[string]string{
		"NameSplit": `function (s) { return "x" }`,
	})
	a.This(env.Returns["Foo"]).Is(TString)
}

func TestFuncRef_ArgDiscriminatedReturn(t *testing.T) {
	a := assert.T(t)
	refs := map[string]string{
		"Pick": `function (secs = false) {
			if secs is false
				return 1
			return "x"
		}`,
	}
	_, envOmit := runPassesWithRefSrcs(`class {
		Foo() { return Pick() }
	}`, refs)
	a.This(envOmit.Returns["Foo"]).Is(TNumber)

	_, envGiven := runPassesWithRefSrcs(`class {
		Foo() { return Pick(1) }
	}`, refs)
	a.This(envGiven.Returns["Foo"]).Is(TString)
}

const timerRefSrc = `class {
	CallClass(block, reps = 1, secs = false) {
		if secs is false
			return 1
		return "x"
	}
}`

// Timer() omits secs -> default false -> Number arm. The Timer-style case.
func TestCtxReturn_CrossClass_OmittedDefault(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithRefSrcs(`class {
		Foo() { return Timer() }
	}`, map[string]string{"Timer": timerRefSrc})
	a.This(env.Returns["Foo"]).Is(TNumber)
}

// Timer(secs: 1) -> secs is Number -> String arm.
func TestCtxReturn_CrossClass_ProvidedNamedArg(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithRefSrcs(`class {
		Foo() { return Timer(secs: 1) }
	}`, map[string]string{"Timer": timerRefSrc})
	a.This(env.Returns["Foo"]).Is(TString)
}

// secs is param index 2; positional Timer(0, 1, 2) binds secs = 2 -> String.
func TestCtxReturn_CrossClass_PositionalArg(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithRefSrcs(`class {
		Foo() { return Timer(0, 1, 2) }
	}`, map[string]string{"Timer": timerRefSrc})
	a.This(env.Returns["Foo"]).Is(TString)
}

// an unknown arg can't decide the guard -> fall back to the full union.
func TestCtxReturnCrossClassUnion(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithRefSrcs(`class {
		Foo(s) { return Timer(secs: s) }
	}`, map[string]string{"Timer": timerRefSrc})
	u, ok := env.Returns["Foo"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber) && u.Contains(TString))
}

// the Memoize/Contributions shape: CallClass comes from the base class, so a
// bare Contributions(...) is its return, not an instance of Contributions
const memoizeRefSrc = `class {
	CallClass(@args) {
		result = Suneido.GetInit('m', { LruCache(.Func) }).Get(@args)
		if Object?(result) and not result.Readonly?()
			result = result.Copy()
		return result
	}
}`

func TestCtxReturn_CallClassInheritedFromBase(t *testing.T) {
	withSigs(t,
		// memoizeRefSrc is a file-scope const, so these are its builtins
		Reg{Receiver: "object", Name: "Copy", Sig: "() :object"},
		Reg{Receiver: "object", Name: "GetInit", Sig: "(member, block) :unknown"},
		Reg{Kind: "free", Name: "Object", Sig: "(@args) :object"},
		Reg{Kind: "free", Name: "Object?", Sig: "(value) :boolean"},
		Reg{Receiver: "object", Name: "Readonly?", Sig: "() :boolean"},
	)
	a := assert.T(t)
	_, env := runPassesWithRefList(`class {
		Foo() { return Contributions("X") }
	}`, []RefSource{
		{Name: "Memoize", Src: memoizeRefSrc},
		{Name: "Contributions", Src: `Memoize { Func(name) { return Object() } }`},
	})
	a.This(env.Returns["Foo"]).Isnt(Instance{Class: "Contributions"})
}

// a base whose CallClass tells us nothing is the Singleton shape - it builds
// `new this` behind an opaque cache - so the instance answer is kept
func TestCtxReturn_UninformativeInheritedCallClass(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithRefList(`class {
		Foo() { return Plugins() }
	}`, []RefSource{
		{Name: "Singleton", Src: `class {
			CallClass() { return Suneido.GetInit(.name(), { new this }) }
			name() { return Name(this) }
		}`},
		{Name: "Plugins", Src: `Singleton { New() {} Get(x) { return x } }`},
	})
	a.This(env.Returns["Foo"]).Is(Instance{Class: "Plugins"})
}

// without the base among the references the chain is unproven, so the instance
// answer is kept but flagged as a guess: checks warn rather than error
func TestCtxReturn_CallClassUnknownBaseGuesses(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithRefList(`class {
		Foo() {
			for fn in Contributions("X")
				return fn
			return 0
		}
	}`, []RefSource{
		{Name: "Contributions", Src: `Memoize { Func(name) { return Object() } }`},
	})
	a.This(errorCount(env, "Foo")).Is(0)
	a.That(hasDiag(env, "Foo", SeverityWarning, "iteration on"))
}

// no class in the chain defines CallClass, so calling it does create an instance
func TestCtxReturn_NoCallClassInChainIsInstance(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithRefList(`class {
		Foo() { return Derived() }
	}`, []RefSource{
		{Name: "Base", Src: `class { Val() { return 1 } }`},
		{Name: "Derived", Src: `Base { Other() { return 2 } }`},
	})
	a.This(env.Returns["Foo"]).Is(Instance{Class: "Derived"})
}

func TestClassReceiverInheritsObjectMethods(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		F() { return (class { X: 123 }).Val_or_func("X") }
	}`, "T")
	a.This(len(methodDiags(env, "F"))).Is(0)
}

func TestCond_IfStringErrors(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: string) { if x { } }
	}`, "T")
	ds := methodDiags(env, "Foo")
	errs := condDiags(ds, SeverityError)
	a.This(len(ds)).Is(1)
	a.This(len(errs)).Is(1)
	for _, d := range errs {
		a.That(strings.Contains(d.Msg, "boolean"))
		a.That(strings.Contains(d.Msg, "string"))
	}
}

func TestCond_IfNumberErrors(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: number) { if x { } }
	}`, "T")
	ds := methodDiags(env, "Foo")
	a.This(len(ds)).Is(1)
	a.This(len(condDiags(ds, SeverityError))).Is(1)
}

func TestCond_WhileStringErrors(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: string) { while x { } }
	}`, "T")
	ds := methodDiags(env, "Foo")
	a.This(len(ds)).Is(1)
	a.This(len(condDiags(ds, SeverityError))).Is(1)
}

func TestCond_DoWhileStringErrors(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: string) { do { } while x }
	}`, "T")
	ds := methodDiags(env, "Foo")
	a.This(len(ds)).Is(1)
	a.This(len(condDiags(ds, SeverityError))).Is(1)
}

func TestCond_ForClassicStringErrors(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: string) { for (i = 0; x; ++i) { } }
	}`, "T")
	ds := methodDiags(env, "Foo")
	a.This(len(ds)).Is(1)
	a.This(len(condDiags(ds, SeverityError))).Is(1)
}

func TestCond_TernaryStringErrors(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: string) { return x ? 1 : 2 }
	}`, "T")
	ds := methodDiags(env, "Foo")
	a.This(len(ds)).Is(1)
	a.This(len(condDiags(ds, SeverityError))).Is(1)
}

func TestCond_IfUnionAllNonBooleanErrors(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: number | string) { if x { } }
	}`, "T")
	ds := methodDiags(env, "Foo")
	errs := condDiags(ds, SeverityError)
	a.This(len(ds)).Is(1)
	a.This(len(errs)).Is(1)
	for _, d := range errs {
		a.That(strings.Contains(d.Msg, "boolean"))
	}
}

func TestCond_IfSentinelMisuseErrors(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		pick: false
		Maybe() { if .pick is true { return 42 }  return false }
		Foo() {
			x = .Maybe()
			if x { }
		}
	}`, "T")
	ds := methodDiags(env, "Foo")
	a.This(len(ds)).Is(1)
	a.This(len(condDiags(ds, SeverityError))).Is(1)
}

func TestCond_NotOperandStringErrors(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: string) { return not x }
	}`, "T")
	ds := methodDiags(env, "Foo")
	a.This(len(ds)).Is(1)
	a.This(len(condDiags(ds, SeverityError))).Is(1)
}

func TestCond_AndOperandsStringErrors(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: string, y: string) { return x and y }
	}`, "T")
	ds := methodDiags(env, "Foo")
	errs := condDiags(ds, SeverityError)
	a.This(len(errs)).Is(2)
	pos := map[int]bool{}
	for _, d := range errs {
		pos[d.Pos] = true
	}
	a.This(len(pos)).Is(2) // distinct positions
}

func TestCond_OrOperandsStringErrors(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: string, y: string) { return x or y }
	}`, "T")
	ds := methodDiags(env, "Foo")
	a.This(len(condDiags(ds, SeverityError))).Is(2)
}

// Pure Unknown: nothing concrete to disprove -> guessed boolean -> silent.
func TestCond_IfUnknownTreatedAsBoolean(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { if x { } }
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}

func TestCond_IfDirtyBooleanSentinelOk(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Flag: false
		Foo() { if .Flag { } }
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}

func TestCond_IfBooleanOk(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: boolean) { if x { } }
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}

func TestCond_IfComparisonOk(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: number) { if x > 0 { } }
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}

func TestCond_IfPredicateOk(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: number | string) { if Number?(x) { } }
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}

func TestCond_IfIsntFalseGuardOk(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		pick: false
		Maybe() { if .pick is true { return 42 }  return false }
		Foo() {
			x = .Maybe()
			if x isnt false { }
		}
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}

func TestCond_WhileBooleanOk(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: boolean) { while x { } }
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}

func TestCond_AndBooleanOperandsOk(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: boolean, y: boolean) { return x and y }
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}

func TestCallsite_Matched_NoDiag(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Use(o: object) {
			o.Has?(1)
		}
	}`, "T")
	a.This(len(methodDiags(env, "Use"))).Is(0)
}

func TestCallsite_ReceiverMismatch_Error(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Sort!", Sig: "(block = false) :object"},
		Reg{Kind: "free", Name: "Use", Sig: "(library :string) :boolean"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Use(s: string) {
			s.Sort!()
		}
	}`, "T")
	ds := methodDiags(env, "Use")
	a.This(len(ds) >= 1).Is(true)
	found := false
	for _, d := range ds {
		if d.Severity == SeverityError && strings.Contains(d.Msg, "not applicable") {
			found = true
		}
	}
	a.That(found)
}

func TestCallsite_GuessSingle_Warn(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Sort!", Sig: "(block = false) :object"},
		Reg{Kind: "free", Name: "Use", Sig: "(library :string) :boolean"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Use(x) {
			x.Sort!()
		}
	}`, "T")
	ds := methodDiags(env, "Use")
	found := false
	for _, d := range ds {
		if d.Severity == SeverityWarning && strings.Contains(d.Msg, "assuming") {
			found = true
		}
	}
	a.That(found)
}

func TestCallsite_NoSuchMethod_Warn(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Use(x) {
			x.TotallyMadeUpMethod(1)
		}
	}`, "T")
	ds := methodDiags(env, "Use")
	found := false
	for _, d := range ds {
		if d.Severity == SeverityWarning && strings.Contains(d.Msg, "no built-in method") {
			found = true
		}
	}
	a.That(found)
}

func TestCallsite_UnionMixed_Error(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Sort!", Sig: "(block = false) :object"},
		Reg{Kind: "free", Name: "Use", Sig: "(library :string) :boolean"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Use(x: string|object) {
			x.Sort!()
		}
	}`, "T")
	ds := methodDiags(env, "Use")
	found := false
	for _, d := range ds {
		if d.Severity == SeverityError && strings.Contains(d.Msg, "at least one path") {
			found = true
		}
	}
	a.That(found)
}

func TestCallsite_UnionOfInstances_NoFalseMixed(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Use(svcTable: SvcLibrary|SvcBook) {
			svcTable.Method?(#foo)
		}
	}`, "T")
	ds := methodDiags(env, "Use")
	for _, d := range ds {
		if d.Severity == SeverityError && strings.Contains(d.Msg, "at least one path") {
			t.Fatalf("false positive union-mixed error on instance union: %s", d.Msg)
		}
	}
	a.This(len(ds)).Is(0)
}

func TestCallsiteUnionFoldsReturn(t *testing.T) {
	a := assert.T(t)
	_, env := runPassesWithClasses(`class {
		Use(x: Foo|Bar) {
			return x.Val()
		}
	}`, map[string]map[string]DynType{
		"Foo": {"Val": TNumber},
		"Bar": {"Val": TNumber},
	})
	a.This(env.Returns["Use"]).Is(TNumber)
}

func TestCallsiteMixedUnionSilent(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Use(x: Foo|string) {
			x.Upper()
		}
	}`, "T")
	ds := methodDiags(env, "Use")
	for _, d := range ds {
		if d.Severity == SeverityError && strings.Contains(d.Msg, "at least one path") {
			t.Fatalf("false positive union-mixed error on instance|primitive union: %s", d.Msg)
		}
	}
	a.This(len(ds)).Is(0)
}

func TestCallsite_ArgCheck_OkPath(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Use(o: object, s: string) {
			o.Has?(s)
		}
	}`, "T")
	a.This(len(methodDiags(env, "Use"))).Is(0)
}

func TestCallsite_DefaultThenReplaceIdiom(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		opts: ()
		Use(opts = false) {
			if opts is false
				opts = .opts
			opts.Members()
		}
	}`, "T")
	a.This(len(methodDiags(env, "Use"))).Is(0)
}

func TestCallsite_DefaultThenReplaceUnionEntry(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		pick: false
		opts: ()
		Maybe() { if .pick is true { return #() }  return false }
		Use() {
			opts = .Maybe()
			if opts is false
				opts = .opts
			opts.Members()
		}
	}`, "T")
	a.This(len(methodDiags(env, "Use"))).Is(0)
}

func TestCallsite_DefaultThenReplaceUnrelatedCond(t *testing.T) {
	withSigs(t,
		Reg{Kind: "static", Class: "Database", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Date", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Ftsearch", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "LruCache", Name: "Members", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Members", Sig: "(list = true, named = true) :object"},
		Reg{Kind: "static", Class: "OpenPGP", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "PdfEncrypt", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Random", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Thread", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "TypeChecker", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Zlib", Name: "Members", Sig: "() :object"},
		Reg{Receiver: "class", Name: "Members", Sig: "(all = false) :object"},
		Reg{Kind: "free", Name: "Use", Sig: "(library :string) :boolean"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		opts: ()
		Use(x: number, opts = false) {
			if x is 5
				opts = .opts
			opts.Members()
		}
	}`, "T")
	ds := methodDiags(env, "Use")
	found := false
	for _, d := range ds {
		if d.Severity == SeverityError && strings.Contains(d.Msg, "at least one path") {
			found = true
		}
	}
	a.That(found)
}

func TestCallsiteDefaultReplaceNested(t *testing.T) {
	withSigs(t,
		Reg{Kind: "static", Class: "Database", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Date", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Ftsearch", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "LruCache", Name: "Members", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Members", Sig: "(list = true, named = true) :object"},
		Reg{Kind: "static", Class: "OpenPGP", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "PdfEncrypt", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Random", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Thread", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "TypeChecker", Name: "Members", Sig: "() :object"},
		Reg{Kind: "static", Class: "Zlib", Name: "Members", Sig: "() :object"},
		Reg{Receiver: "class", Name: "Members", Sig: "(all = false) :object"},
		Reg{Kind: "free", Name: "Use", Sig: "(library :string) :boolean"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		opts: ()
		Use(c: boolean, opts = false) {
			if opts is false
				if c
					opts = .opts
			opts.Members()
		}
	}`, "T")
	ds := methodDiags(env, "Use")
	found := false
	for _, d := range ds {
		if d.Severity == SeverityError && strings.Contains(d.Msg, "at least one path") {
			found = true
		}
	}
	a.That(found)
}

func TestIfFunctionsBecomeClassContainers(t *testing.T) {
	a := assert.T(t)
	psrc := ParseClass(`function () {}`)
	cobj := NewClassObject("function", psrc)

	a.This(cobj.Name == "function")
	a.This(cobj.Base == "")
	a.This(len(cobj.Members) == 0)
	a.This(len(cobj.Methods) == 1)
}

func TestCompoundAssignLhsCoercionIsHardError(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() { n = 5; n $= "x"; return n }
	}`, "T")
	a.That(hasDiag(env, "Foo", SeverityError, `operator "$="`))

	_, env = runPasses(`class {
		Bar() { s = ""; s += 1; return s }
	}`, "T")
	a.That(hasDiag(env, "Bar", SeverityError, `operator "+="`))
}

// member LHS goes through the Mem stamp (pre-type) - keep it erroring.
func TestCompoundAssignMemberCoercion(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		n: 5
		Foo() { .n $= "x" }
	}`, "T")
	a.That(hasDiag(env, "Foo", SeverityError, `operator "$="`))
}

func TestCompoundAssignWellTypedSilent(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() { n = 5; n += 1; s = "a"; s $= "b"; return n }
	}`, "T")
	a.This(countDiagsWith(env, `"+="`)).Is(0)
	a.This(countDiagsWith(env, `"$="`)).Is(0)
}

func TestIncDecAdvancesScopeToNumber(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			s = ""
			s++
			return s $ "x"
		}
	}`, "T")
	// the ++ itself coerces (String operand) AND the later $ reads a Number
	a.That(hasDiag(env, "Foo", SeverityError, `operator "++"`))
	a.That(hasDiag(env, "Foo", SeverityError, `operator "$"`))
}

func TestCompoundAssignResultTypeFlows(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() { s = "a"; s $= "b"; return s }
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
}

// ---- assert assignment check runs post-narrowing (the Rect.ss Set shape) --

func TestAssertCheck_GuardedAssignNoFalsePositive(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(.x) { Assert(Number?(x)) }
		Set(x = false) {
			if x isnt false
				.x = x
		}
	}`, "Rect")
	a.That(!errOn(env, "Set", "ground truth"))
}

func TestAssertCheck_ProvableViolationStillErrors(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(.x) { Assert(Number?(x)) }
		Clear() { .x = false }
	}`, "Rect")
	a.That(errOn(env, "Clear", "ground truth"))
}

// ---- method-local asserts narrow, New asserts pin -------------------------

func TestLocalAssert_NarrowsMemberForRestOfMethod(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		focused: false
		Blur() { .focused = false }
		Move(row) { .focused = row }
		vk_home(shift) {
			Assert(Number?(.focused))
			return -.focused
		}
	}`, "T")
	a.That(!errOn(env, "vk_home", `operator "-"`))
	a.This(env.Returns["vk_home"]).Is(TNumber)
}

func TestLocalAssert_DoesNotPinClassWide(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		focused: false
		Blur() { .focused = false }
		vk_home(shift) {
			Assert(Number?(.focused))
			return -.focused
		}
	}`, "T")
	a.That(!errOn(env, "Blur", "ground truth"))
}

func TestLocalAssert_LocalRefinedAfterOnly(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			Assert(Number?(x))
			return x + 1
		}
	}`, "T")
	a.That(!errOn(env, "Foo", `operator "+"`))
	a.This(env.Returns["Foo"]).Is(TNumber)
}

// ---- write-set call invalidation + helper postconditions, composed --------

func TestHelperPostcondition_SelectableShape(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		focused: false
		Blur() { .focused = false }
		GetNumRows() { return 10 }
		GetNumVisibleRows() { return 5 }
		RepaintRow(r) { }
		selectable?()
			{
			if .GetNumRows() < 1
				{
				.focused = false
				return false
				}
			if .focused is false
				.focused = 0
			if .GetNumVisibleRows() < 1
				return false
			return true
			}
		Select(amt)
			{
			if not .selectable?()
				return
			newFocus = .focused + amt
			.RepaintRow(.focused)
			.focused = newFocus
			return newFocus
			}
	}`, "T")
	a.That(!errOn(env, "Select", `operator "+"`))
}

// the then-arm form: `if .helper() { use members }`.
func TestHelperPostcondition_ThenArm(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		width: false
		Clear() { .width = false }
		Resize(w) { .width = w * 1 }
		ready?()
			{
			if .width is false
				return false
			return true
			}
		Area(h)
			{
			if .ready?()
				return .width * h
			return 0
			}
	}`, "T")
	a.That(!errOn(env, "Area", `operator "*"`))
}

func TestHelperPostcondition_NoFalseLicense(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		width: false
		Clear() { .width = false }
		odd?()
			{
			if .Count() % 2 is 1
				return true
			return false
			}
		Count() { return 3 }
		Area(h)
			{
			if .odd?()
				return .width * h
			return 0
			}
	}`, "T")
	a.That(errOn(env, "Area", `operator "*"`))
}

// public helpers are overridable cross-class, so no postconditions for them.
func TestHelperPostcondition_PublicHelperExcluded(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		width: false
		Clear() { .width = false }
		Ready?()
			{
			if .width is false
				return false
			return true
			}
		Area(h)
			{
			if .Ready?()
				return .width * h
			return 0
			}
	}`, "T")
	a.That(errOn(env, "Area", `operator "*"`))
}

func TestNewThisResolvesOwnClass(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		CallClass(src) {
			parser = new this(src)
			return parser.Expression()
		}
		Expression() { return 42 }
	}`, "Tdop")
	a.This(env.Returns["CallClass"]).Is(TNumber)
}

// loop join: zero-iteration soundness ------------------------------------

func TestLoopJoinZeroIterationWhile(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(c) {
			x = 0
			while c is true
				{ x = "s" }
			return x
		}
	}`, "T")
	ret := env.Returns["Foo"]
	a.That(unionContains(ret, TNumber))
	a.That(unionContains(ret, TString))
}

func TestLoopJoinZeroIterationForIn(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(ob) {
			x = 0
			for v in ob
				{ x = "s" }
			return x
		}
	}`, "T")
	ret := env.Returns["Foo"]
	a.That(unionContains(ret, TNumber))
	a.That(unionContains(ret, TString))
}

// loop fixpoint: loop-carried feedback ------------------------------------

func TestLoopFixpointFeedbackTypes(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(c) {
			s = ""
			while c is true
				{ s = s.Size() }
			return s
		}
	}`, "T")
	ret := env.Returns["Foo"]
	a.That(unionContains(ret, TString))
	a.That(unionContains(ret, TNumber))
	a.That(hasDiag(env, "Foo", SeverityError, "not applicable on at least one path"))
}

func TestUnionMixedDiagNamesReceiverAndBadArms(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(c) {
			s = ""
			while c is true
				{ s = s.Size() }
			return s
		}
	}`, "T")
	// receiver expression is echoed back, up front after the method name...
	a.That(hasDiag(env, "Foo", SeverityError, "on `s`"))
	// ...and the failing arm (Number has no Size) is named explicitly.
	a.That(hasDiag(env, "Foo", SeverityError, "no overload for number"))
	// the leading phrase other tests match on is preserved.
	a.That(hasDiag(env, "Foo", SeverityError, "not applicable on at least one path"))
}

// while-guard narrowing ----------------------------------------------------

func TestWhileGuardNarrowsBody(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(n = false) {
			while n isnt false
				{ return n + 1 }
			return 0
		}
	}`, "T")
	a.That(!hasDiag(env, "Foo", SeverityError, `operator "+"`))
}

func TestWhilePredicateGuardNarrowsBody(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(n) {
			total = 0
			while Number?(n)
				{
				total += n
				n = .next()
				}
			return total
		}
	}`, "T")
	a.That(!hasDiag(env, "Foo", SeverityError, `operator "+="`))
	a.That(!hasDiag(env, "Foo", SeverityWarning, `operator "+="`))
}

func TestWhileGuardNarrowsMember(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(.width = false) { }
		Total(h)
			{
			total = 0
			while .width isnt false
				{
				total = .width * h
				.width = false
				}
			return total
			}
	}`, "T")
	a.That(!hasDiag(env, "Total", SeverityError, `operator "*"`))
}

func TestLoopCarriedAssignInvalidates(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x, c) {
			if Number?(x)
				{
				while c is true
					{
					y = x + 1
					x = "s"
					}
				return x
				}
			return 0
		}
	}`, "T")
	a.That(hasDiag(env, "Foo", SeverityError, `operator "+"`) ||
		hasDiag(env, "Foo", SeverityWarning, `operator "+"`))
}

// for-loop guard ------------------------------------------------------------

func TestForCondGuardNarrowsBody(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(n = false) {
			for (i = 0; n isnt false; ++i)
				{ return n + i }
			return 0
		}
	}`, "T")
	a.That(!hasDiag(env, "Foo", SeverityError, `operator "+"`))
}

// trinary narrow-to-unknown through ExprPos ---------------------------------

func TestTrinaryNarrowAwayExprPos(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x = false) {
			y = x is false ? 0 : x
			return y + 1
		}
	}`, "T")
	a.That(!hasDiag(env, "Foo", SeverityError, `operator "+"`))
}

// catch variable -------------------------------------------------------------

func TestCatchVarIsString(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			try
				return 1
			catch (e)
				return e $ "!"
			}
	}`, "T")
	ret := env.Returns["Foo"]
	a.That(unionContains(ret, TNumber))
	a.That(unionContains(ret, TString))
	a.That(!hasDiag(env, "Foo", SeverityWarning, `operator "$"`))
	a.That(!hasDiag(env, "Foo", SeverityError, `operator "$"`))
}

func TestTryCatchBranchScopes(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			x = 0
			try
				x = "s"
			catch (unused)
				x = x
			return x
		}
	}`, "T")
	ret := env.Returns["Foo"]
	a.That(unionContains(ret, TNumber))
	a.That(unionContains(ret, TString))
}

func TestReturnContract_ParamPropagatesToLocal(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: string) {
			y = x
			return y
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TString)
	a.This(len(*env.Diagnostics)).Is(0)
}

// OK: number subset of number - silent.
func TestReturnContract_OkExact(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Add(a: number, b: number) : number {
			return a + b
		}
	}`, "T")
	a.This(env.Returns["Add"]).Is(TNumber)
	a.This(len(*env.Diagnostics)).Is(0)
}

// OK at both return sites: member subset of declared union.
func TestReturnContract_OkUnionMembers(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		EitherOk(flag: boolean) : number | string {
			if flag { return 1 }
			return "zero"
		}
	}`, "T")
	a.This(len(*env.Diagnostics)).Is(0)
}

// ERROR: provable mismatch. number not subset of string.
func TestReturnContract_ErrorLiteralMismatch(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		BadLiteral() : string {
			return 42
		}
	}`, "T")
	ds := diagsByMethod(env)["BadLiteral"]
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityError)
	// Public return is still the declared contract.
	a.This(env.Returns["BadLiteral"]).Is(TString)
}

func TestReturnContract_ErrorThroughOperator(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		BadAdd(x: number) : string {
			return x + 1
		}
	}`, "T")
	ds := diagsByMethod(env)["BadAdd"]
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityError)
}

func TestReturnContractPerSite(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Mixed(flag: boolean) : string {
			if flag { return 42 }
			return "hi"
		}
	}`, "T")
	ds := diagsByMethod(env)["Mixed"]
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityError)
	pos := ds[0].Pos
	a.That(pos > 0)
}

func TestReturnContract_WarningOnDirtyUnion(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		DirtyUnion(flag: boolean) : string | number {
			z = flag ? "hi" : GetMystery()
			return z
		}
	}`, "T")
	ds := diagsByMethod(env)["DirtyUnion"]
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityWarning)
}

func TestReturnContract_WarningOnPureUnknown(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Ambiguous() : string | number {
			return GetMystery()
		}
	}`, "T")
	ds := diagsByMethod(env)["Ambiguous"]
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityWarning)
}

func TestReturnContract_ErrorAndWarningCoexist(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		LeakAndAmbiguous(flag: boolean) : string | number {
			if flag { return true }
			return GetMystery()
		}
	}`, "T")
	ds := diagsByMethod(env)["LeakAndAmbiguous"]
	a.This(len(ds)).Is(2)
	a.This(countSeverity(ds, SeverityError)).Is(1)
	a.This(countSeverity(ds, SeverityWarning)).Is(1)
}

func TestReturnContract_NarrowingSilencesWarning(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Round(x: number | string) : number {
			if Number?(x) { return x }
			return 0
		}
	}`, "T")
	ds := diagsByMethod(env)["Round"]
	a.This(len(ds)).Is(0)
}

func TestReturnContract_NarrowerBodyNoDiagnostic(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() : number | string {
			return 1
		}
	}`, "T")
	a.This(len(*env.Diagnostics)).Is(0)
}

func TestReturnContract_NoAnnotationNoDiagnostic(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Untyped() {
			return GetMystery()
		}
		Mixed(flag: boolean) {
			if flag { return 1 }
			return "two"
		}
	}`, "T")
	a.This(len(*env.Diagnostics)).Is(0)
}

func TestReturnContract_ReturnThrowSkipped(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Bail() : number {
			return throw "boom"
		}
	}`, "T")
	ds := diagsByMethod(env)["Bail"]
	a.This(len(ds)).Is(0)
}

func TestReturnContract_BooleanLeakIsError(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() : string | number {
			return true
		}
	}`, "T")
	ds := diagsByMethod(env)["Foo"]
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityError)
}

func TestReturnContractGroundTruth(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Identity(x: number) : string {
			return x
		}
	}`, "T")
	ds := diagsByMethod(env)["Identity"]
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityError)
}

func TestReturnContract_ThrowInBranchIgnored(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(flag) : number {
			if flag { return throw "nope" }
			return 1
		}
	}`, "T")
	a.This(env.Returns["Foo"]).Is(TNumber)
	a.This(len(*env.Diagnostics)).Is(0)
}

func TestReturnContract_ThrowSkippedUnknownWarns(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(flag) : number {
			if flag { return throw "bad" }
			return GetMystery()
		}
	}`, "T")
	ds := diagsByMethod(env)["Foo"]
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityWarning)
}

func TestGuessBooleanContextWarns(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Match", Sig: "(pattern, pos=false, prev :boolean =false) :false|object"},
		Reg{Kind: "static", Class: "Database", Name: "Token", Sig: "() :string"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		New() { .Scan = TdopScanner() }
		Statement(alt_end = false) {
			expr = TdopStmtExpr()
			if expr.Match("IDENTIFIER") and .Scan.Token().Match("COMMA")
				return 1
			mustMatch = not .Scan.Token().Match(alt_end)
			return mustMatch
		}
	}`, "T")
	sawCapped := false
	for _, d := range *env.Diagnostics {
		if strings.Contains(d.Msg, "expects boolean") {
			a.That(d.Severity != SeverityError)
			if strings.Contains(d.Msg, "type guessed from builtin overloads") {
				sawCapped = true
			}
		}
	}
	a.That(sawCapped)
}

// same cap for operator operands (arithmetic on a guessed type).
func TestGuessAssumedTypeCapsOperatorAtWarning(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			return x.Find(5) + 1
		}
	}`, "T")
	for _, d := range *env.Diagnostics {
		if strings.Contains(d.Msg, `operator "+"`) {
			a.That(d.Severity != SeverityError)
		}
	}
}

func TestGuessMarkClearedWhenReceiverResolves(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			s = "abc"
			return s.Size() $ "!"
		}
	}`, "T")
	found := false
	for _, d := range *env.Diagnostics {
		if d.Severity == SeverityError && strings.Contains(d.Msg, `operator "$"`) {
			found = true
		}
	}
	a.That(found)
}

func TestThisCallFallbackScansAllOverloads(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Find", Sig: "(value) :false|unknown"},
		Reg{Receiver: "string", Name: "Find", Sig: "(string :string, pos=0) :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() { return .Find(5) }
	}`, "T")
	ret := env.Returns["Foo"]
	u, ok := ret.(Union)
	a.That(ok)
	a.That(u.Contains(TNumber))
	a.That(u.IsDirty)
}

func TestExtensionClassThisIsString(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Find", Sig: "(value) :false|unknown"},
		Reg{Receiver: "string", Name: "Find", Sig: "(string :string, pos=0) :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Trimmed(c) {
			i = .Find(c)
			return i + 1
		}
	}`, "Strings")
	a.This(env.Returns["Trimmed"]).Is(TNumber)
	for _, d := range *env.Diagnostics {
		if strings.Contains(d.Msg, `operator "+"`) {
			t.Fatalf("unexpected diagnostic: %s", d.Msg)
		}
	}
}

// and a Dates.ss shape: this-dispatch must reach the Date overloads.
func TestExtensionClassThisIsDate(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "date", Name: "Plus", Sig: "(years=0, months=0, days=0, hours=0, minutes=0, seconds=0, milliseconds=0) :date"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Tomorrow() { return .Plus(days: 1) }
	}`, "Dates")
	a.This(env.Returns["Tomorrow"]).Is(TDate)
}

func TestNonExtensionClassThisUnchanged(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() { return .Size() }
	}`, "MyHelper")
	// Size agrees across overloads (:number) - multi-agree guess, clean.
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestOp_NumericOk(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Add(a: number, b: number) : number { return a + b }
	}`, "T")
	a.This(len(methodDiags(env, "Add"))).Is(0)
}

func TestOp_ConcatOk(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Cat(a: string, b: string) : string { return a $ b }
	}`, "T")
	a.This(len(methodDiags(env, "Cat"))).Is(0)
}

// ---- Errors --------------------------------------------------------------

func TestOp_ErrorStringInArithmetic(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		BadAdd(s: string, n: number) : number { return s + n }
	}`, "T")
	ds := methodDiags(env, "BadAdd")
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityError)
	a.That(strings.Contains(ds[0].Msg, "number"))
	a.That(strings.Contains(ds[0].Msg, "string"))
}

func TestOp_ErrorNumberInConcat(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		BadCat(s: string, n: number) : string { return s $ n }
	}`, "T")
	ds := methodDiags(env, "BadCat")
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityError)
}

func TestOp_PerOperandErrorsAreIndependent(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		TwoBad(s: string, t: string) : number { return s * t }
	}`, "T")
	ds := methodDiags(env, "TwoBad")
	a.This(len(ds)).Is(2)
	a.This(ds[0].Severity).Is(SeverityError)
	a.This(ds[1].Severity).Is(SeverityError)
	a.That(ds[0].Pos != ds[1].Pos)
}

func TestOp_UnaryNegationOnBoolean(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		BadUnary(flag: boolean) : number { return -flag }
	}`, "T")
	ds := methodDiags(env, "BadUnary")
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityError)
}

func TestOp_WarningOnUnknownOperand(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Amb(n) : number { return n + 1 }
	}`, "T")
	ds := methodDiags(env, "Amb")
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityWarning)
}

// ---- Narrowing silences the check ----------------------------------------

func TestOp_NarrowingSilencesWarning(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Round(x: number | string) : number {
			if Number?(x) { return x + 1 }
			return 0
		}
	}`, "T")
	a.This(len(methodDiags(env, "Round"))).Is(0)
}

func TestOp_AndChainProgressiveNarrowing(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Round(x: number | string) : number {
			if Number?(x) and x + 1 > 0
				return x + 1
			return 0
		}
	}`, "T")
	a.This(len(methodDiags(env, "Round"))).Is(0)
}

func TestOp_OrChainProgressiveNarrowing(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Test(x: number | string) : boolean {
			return not Number?(x) or x + 1 > 0
		}
	}`, "T")
	a.This(len(methodDiags(env, "Test"))).Is(0)
}

func TestOp_AndChainThreeOperandsAccumulate(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Test(x: number | string | boolean) : number {
			if not String?(x) and not Boolean?(x) and x + 1 > 0
				return x + 1
			return 0
		}
	}`, "T")
	a.This(len(methodDiags(env, "Test"))).Is(0)
}

func TestOp_AndChainInlineAssignmentNarrowing(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		pick: false
		Maybe() {
			if .pick is true { return 42 }
			return false
		}
		Foo(n: number) {
			if (n > 0 and (false isnt pos = .Maybe()) and pos + 1 > 0)
				return pos
			return 0
		}
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}

func TestOp_OrChainTrueBodyStringStillErrors(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Test(x: number | string) : number {
			if Number?(x) or String?(x)
				return x + 1
			return 0
		}
	}`, "T")
	ds := methodDiags(env, "Test")
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityError)
}

func TestOp_NarrowedTernaryThroughTempLocal(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		pick: false
		Maybe() { if .pick is true { return 42 }  return false }
		Foo() : number {
			np = .Maybe()
			np2 = Number?(np) ? np : 0
			return 1 + np2
		}
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}

func TestOp_NarrowedTernaryInlineStillWorks(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		pick: false
		Maybe() { if .pick is true { return 42 }  return false }
		Foo() : number {
			np = .Maybe()
			return 1 + (Number?(np) ? np : 0)
		}
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}

func TestOp_NarrowedIsntFalseThroughTempLocal(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		pick: false
		Maybe() { if .pick is true { return 42 }  return false }
		Foo() : number {
			np = .Maybe()
			np2 = np isnt false ? np : 0
			return 1 + np2
		}
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}

// ---- No-else if-merge narrows the post-if scope --------------------------

func TestOp_DefaultThenReplaceArithmetic(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(n = false) : number {
			if n is false
				n = 0
			return n + 1
		}
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}

func TestOpDefaultReplaceUnknown(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(n) : number {
			if n is false
				n = 0
			return n + 1
		}
	}`, "T")
	ds := methodDiags(env, "Foo")
	found := false
	for _, d := range ds {
		if d.Severity == SeverityWarning &&
			strings.Contains(d.Msg, `operator "+"`) &&
			strings.Contains(d.Msg, "cannot prove") {
			found = true
		}
	}
	a.That(found)
}

// ---- Operators outside scope are NOT checked -----------------------------

func TestOp_LooseOperatorsAreSilent(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Cmp(s: string, n: number) : boolean { return s is n }
	}`, "T")
	a.This(len(methodDiags(env, "Cmp"))).Is(0)
}

func TestOp_CoexistsWithReturnContract(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Combo(s: string) : string { return s + 1 }
	}`, "T")
	ds := methodDiags(env, "Combo")
	a.This(len(ds)).Is(2)
	hasOpErr := false
	hasReturnErr := false
	for _, d := range ds {
		switch {
		case strings.Contains(d.Msg, `operator`):
			hasOpErr = true
		case strings.Contains(d.Msg, "return type"):
			hasReturnErr = true
		}
	}
	a.That(hasOpErr)
	a.That(hasReturnErr)
}

func TestOp_BitwiseRequiresNumber(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		BadBit(s: string, n: number) : number { return s & n }
	}`, "T")
	ds := methodDiags(env, "BadBit")
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityError)
}

func TestOp_ShiftRequiresNumber(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		BadShift(s: string, n: number) : number { return s << n }
	}`, "T")
	ds := methodDiags(env, "BadShift")
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityError)
}

func TestOp_ModuloRequiresNumber(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Bad(s: string, n: number) : number { return s % n }
	}`, "T")
	ds := methodDiags(env, "Bad")
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityError)
}

func TestOp_ConcatWithNumberOperand(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Bad(s: string, n: number) : string { return s $ n }
	}`, "T")
	ds := methodDiags(env, "Bad")
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityError)
}

// ---- ExprPos wrapping (regression) ---------------------------------------

func TestOp_OperatorInsideTrinaryIsChecked(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x: string, y) {
			z = x + Number?(y) ? y : 0
		}
	}`, "T")
	ds := methodDiags(env, "Foo")
	hasErr := false
	for _, d := range ds {
		if d.Severity == SeverityError &&
			strings.Contains(d.Msg, "operator") &&
			strings.Contains(d.Msg, "string") {
			hasErr = true
		}
	}
	a.That(hasErr)
}

func TestOp_CompoundAssignDocsExistingQuirk(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Compound(s: string) { s += "x" }
	}`, "T")
	ds := methodDiags(env, "Compound")
	a.This(len(ds)).Is(2)
	for _, d := range ds {
		a.This(d.Severity).Is(SeverityError)
		a.That(strings.Contains(d.Msg, `"+="`))
	}
}

func TestThisCallResolvesParentMethod(t *testing.T) {
	a := assert.T(t)
	parentSrc := `class {
        OurNewMethod() { return 10 }
    }`
	childSrc := `Parent {
        Foo() { return .OurNewMethod() }
    }`
	resolver := mapResolver{"Parent": parentSrc}
	_, env := TypeInfer("Child", childSrc, resolver)
	a.This(env.Returns["Foo"]).Is(TNumber)
}

func TestThisCallCaseSensitive(t *testing.T) {
	a := assert.T(t)
	parentSrc := `class {
        OurNewMethod() { return "from parent" }
    }`
	childSrc := `Parent {
        ourNewMethod() { return 4 }
        Foo() { return .OurNewMethod() }
    }`
	resolver := mapResolver{"Parent": parentSrc}
	_, env := TypeInfer("Child", childSrc, resolver)
	a.This(env.Returns["Foo"]).Is(TString) // not TNumber
}

func TestOp_CompareCrossTypeWarns(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Bad(n: number, s: string) : boolean { return n < s }
	}`, "T")
	ds := methodDiags(env, "Bad")
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityWarning)
	a.That(strings.Contains(ds[0].Msg, "cross-type"))
	a.That(strings.Contains(ds[0].Msg, "StrictCompare"))
}

func TestOp_CompareSameTypeSilent(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Cmp(a: number, b: number) : boolean { return a < b }
	}`, "T")
	a.This(len(methodDiags(env, "Cmp"))).Is(0)
}

func TestOp_CompareUnionOperandSilent(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Maybe() { if .pick is true { return 42 }; return false }
		pick: false
		Foo() : boolean { x = .Maybe(); return x < 10 }
	}`, "T")
	a.This(len(methodDiags(env, "Foo"))).Is(0)
}

func TestOp_CompareUnknownSilent(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Cmp(a, b) : boolean { return a < b }
	}`, "T")
	a.This(len(methodDiags(env, "Cmp"))).Is(0)
}

func TestOp_NarrowedReturnPropagatesIntoCallSite(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		val: false
		SetVal(x) { .val = x }
		Fetch() {
			v = .val
			return String?(v) ? v : "fallback"
		}
		Caller() { return 'prefix_' $ .Fetch() }
	}`, "T")
	// Fetch() must resolve to String after narrowing.
	a.This(env.Returns["Fetch"]).Is(TString)
	// Caller() uses $ (requires String operands): must produce no diagnostics.
	a.This(len(methodDiags(env, "Caller"))).Is(0)
}

func TestOp_CompareAllOperatorsWarn(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		F1(n: number, s: string) : boolean { return n <  s }
		F2(n: number, s: string) : boolean { return n <= s }
		F3(n: number, s: string) : boolean { return n >  s }
		F4(n: number, s: string) : boolean { return n >= s }
	}`, "T")
	for _, m := range []string{"F1", "F2", "F3", "F4"} {
		ds := methodDiags(env, m)
		a.This(len(ds)).Is(1)
		a.This(ds[0].Severity).Is(SeverityWarning)
	}
}

func TestFreeFunctionSigNotMethodCandidate(t *testing.T) {
	withSigs(t,
		Reg{Kind: "free", Name: "Query1", Sig: "(@args) :false|object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Get(t = false) {
			DoWithTran(t)
				{|t|
				rec = t.Query1('q')
				}
			return rec
		}
	}`, "T")
	a.This(len(methodDiags(env, "Get"))).Is(0)

	_, env = runPasses(`class {
		Foo(x) { return x.Query1('q') }
	}`, "T")
	ds := methodDiags(env, "Foo")
	a.This(len(ds)).Is(1)
	a.This(ds[0].Severity).Is(SeverityWarning)
	a.That(strings.Contains(ds[0].Msg, `no built-in method "Query1"`))

	// the free-call form still resolves through the free signature.
	_, env = runPasses(`class {
		Bar() { return Query1('q') }
	}`, "T")
	a.That(unionContains(env.Returns["Bar"], TObject))
	a.That(unionContains(env.Returns["Bar"], TFalse))
}

// ---- branch joins must propagate kills, not just additions ----

// A call inside a branch that writes the member has to invalidate the
// refinement for the statements *after* the if, even though the call ran
// against a cloned scope.
func TestNarrowMemberKillInsideIfEscapesBranch(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		f1: #20200101
		M1() { .f1 = "x" }
		M0() {
			.f1 = #20200101
			if (true) { .M1() }
			return .f1
		}
	}`, "T")
	a.That(unionContains(env.Returns["M0"], TString))
}

// same kill, reached through a nested if so the top-level-assignment merge
// cannot see it either
func TestNarrowMemberNestedReassign(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		f1: #20200101
		M0(p) {
			.f1 = #20200101
			if (true) { if (p) { .f1 = "x" } }
			return .f1
		}
	}`, "T")
	a.That(unionContains(env.Returns["M0"], TString))
}

// wrapping a statement in if(true) must not change the answer - the
// metamorphic property TestPropMetaIfTrueWrap checks
func TestNarrowIfTrueWrapAgreesWithFlatBody(t *testing.T) {
	a := assert.T(t)
	const flat = `class {
		f0: #20200101
		f1: #20200101
		M0(p0) {
			.f1 = (false ? .f0 : .f0)
			v0 = (.M0("a") < #20200101)
			return .f1
		}
	}`
	const wrapped = `class {
		f0: #20200101
		f1: #20200101
		M0(p0) {
			.f1 = (false ? .f0 : .f0)
			if (true) { v0 = (.M0("a") < #20200101) }
			return .f1
		}
	}`
	_, flatEnv := runPasses(flat, "T")
	_, wrappedEnv := runPasses(wrapped, "T")
	a.This(wrappedEnv.Returns["M0"]).Is(flatEnv.Returns["M0"])
}

// a refinement no branch touches survives the join untouched
func TestNarrowMemberSurvivesUnrelatedBranch(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		f1: #20200101
		M0(p) {
			.f1 = #20200101
			if (p) { x = 1 }
			return .f1
		}
	}`, "T")
	a.This(env.Returns["M0"]).Is(TDate)
}

// ---- parentheses on an assignment RHS must not dirty the member ----

func TestMemberNotDirtiedByParenthesizedRhs(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		f0: #20200101
		f1: #20200101
		M0() { .f1 = (false ? .f0 : .f0) }
	}`, "T")
	a.This(env.Members["f1"]).Is(TDate)
}

// after a guard that empties a member's known type, the refinement's two
// consumers deliberately disagree. a read applies it unconditionally - even
// the widening to TUnknown - which is what keeps the provably-dead path
// silent, at the cost of a dirty return arm. an assignment only accepts a
// refinement through isNarrower, so the local keeps the precise seed type
// and publishes clean. see narrowEqAssign vs the *ast.Mem case in narrowWalk;
// neither side should be "fixed" to match the other - widening the local
// spreads dirt into published var types, and skipping the widening at reads
// reintroduces the guarded-param false positives.
func TestContradictionGuardReadVsAssign(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		f0: false
		Direct() {
			if false is .f0
				return false
			return .f0
		}
		ViaLocal() {
			if false is .f0
				return false
			x = .f0
			return x
		}
	}`, "T")
	u, ok := env.Returns["Direct"].(Union)
	a.That(ok)
	a.That(u.IsDirty)
	a.That(u.Contains(TFalse))
	a.This(env.Returns["ViaLocal"]).Is(TFalse)
}
