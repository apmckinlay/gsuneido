// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"strings"
	"testing"
)

func TestVoidCallOperand(t *testing.T) {
	src := `class {
	Void() { return }
	TestVoid() { return 123 + .Void() }
	}`
	_, env := runPasses(src, "T")
	found := false
	for _, d := range diagList(env) {
		if strings.Contains(d.Msg, "got unknown") {
			t.Errorf("void return degraded to unknown: %s", d.Msg)
		}
		if d.Severity == SeverityError && strings.Contains(d.Msg, `"no return value"`) {
			found = true
		}
	}
	if !found {
		t.Error("no error reported for void call used as + operand")
	}
}

func TestVoidCallReceiver(t *testing.T) {
	src := `class {
	Void() { return }
	TestChain() { return .Void().Size() }
	}`
	cls, ok := safeParse(src, "T")
	if !ok {
		t.Fatal("source did not parse")
	}
	env := NewTypeEnv()
	sigs := AnnotationSet{"Size": {
		{Receiver: TString, Returns: TNumber},
		{Receiver: TObject, Returns: TNumber},
	}}
	RunPipeline(cls, env, &PassCtx{Annotations: sigs}, nil)
	found := false
	for _, d := range diagList(env) {
		if d.Severity == SeverityError && strings.Contains(d.Msg, "void") {
			found = true
		}
	}
	if !found {
		t.Error("no error for method call on void receiver")
	}
}

func TestVoidAllThrowNotVoid(t *testing.T) {
	src := `class {
	Fail() { throw "should not be called" }
	Use() { return 1 + .Fail() }
	}`
	_, env := runPasses(src, "T")
	if got := env.Returns["Fail"]; got == TVoid {
		t.Error("all-throw method published Void")
	}
	for _, d := range diagList(env) {
		if strings.Contains(d.Msg, `"no return value"`) {
			t.Errorf("false positive on all-throw callee: %s", d.Msg)
		}
	}
}

func TestVoidCallPassthrough(t *testing.T) {
	src := `class {
	Void() { return }
	Passthru() { return .Void() }
	Discard() { .Void(); return 5 }
	}`
	_, env := runPasses(src, "T")
	for _, d := range diagList(env) {
		if strings.Contains(d.Msg, `"no return value"`) {
			t.Errorf("false positive on legal pass-through: %s", d.Msg)
		}
	}
}
