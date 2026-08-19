// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/core"
)

func TestBlockFormReturn(t *testing.T) {
	withSigs(t,
		Reg{Kind: "free", Name: "Dir", Sig: "(path :string ='*', files :boolean =false, details :boolean =false, block=false) :object|void"},
	)
	ret := func(src string) DynType {
		t.Helper()
		_, env := runPasses(src, "T")
		return env.LookupReturn("Foo")
	}
	tests := []struct {
		name string
		src  string
		want DynType
	}{
		{"no args", `class { Foo() { return Dir() } }`, TObject},
		{"no block", `class { Foo() { return Dir("x") } }`, TObject},
		{"other positionals only",
			`class { Foo() { return Dir("x", true, true) } }`, TObject},
		{"other named only",
			`class { Foo() { return Dir("x", details: true) } }`, TObject},
		{"trailing block", `class { Foo() { return Dir("x") { } } }`, TVoid},
		{"named block", `class { Foo() { return Dir("x", block: { }) } }`, TVoid},
		{"positional block",
			`class { Foo() { return Dir("x", false, false, { }) } }`, TVoid},
		{"named block after named",
			`class { Foo() { return Dir("x", details: true, block: { }) } }`, TVoid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ret(tt.src); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlockFormReturnGuards(t *testing.T) {
	objOrVoid := Union{Types: []DynType{TObject, TVoid}}.Fold()
	objOrStr := Union{Types: []DynType{TObject, TString}}.Fold()
	withBlock := []Param{{Name: "path"}, {Name: "block"}}
	noBlock := []Param{{Name: "path"}}
	blockArg := []ast.Arg{{Name: core.SuStr("block")}}
	spread := []ast.Arg{{Name: core.SuStr("@")}}

	tests := []struct {
		name string
		sig  Signature
		args []ast.Arg
		want DynType
	}{
		{"non-union return left alone",
			Signature{Params: withBlock, Returns: TObject}, nil, TObject},
		{"union without void left alone",
			Signature{Params: withBlock, Returns: objOrStr}, nil, objOrStr},
		{"void union but no block param left alone",
			Signature{Params: noBlock, Returns: objOrVoid}, nil, objOrVoid},
		{"unknown return left alone",
			Signature{Params: withBlock, Returns: TUnknown}, nil, TUnknown},
		{"no block arg drops void",
			Signature{Params: withBlock, Returns: objOrVoid}, nil, TObject},
		{"block arg yields void",
			Signature{Params: withBlock, Returns: objOrVoid}, blockArg, TVoid},
		{"args spread keeps both halves",
			Signature{Params: withBlock, Returns: objOrVoid}, spread, objOrVoid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := blockFormReturn(&tt.sig, tt.args)
			if got.String() != tt.want.String() {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlockFormLeavesOtherBlockBuiltinsAlone(t *testing.T) {
	ret := func(src string) DynType {
		t.Helper()
		_, env := runPasses(src, "T")
		return env.LookupReturn("Foo")
	}
	for _, call := range []string{`RunPiped("x")`, `File("x")`} {
		without := ret(`class { Foo() { return ` + call + ` } }`)
		with := ret(`class { Foo() { return ` + call + ` { } } }`)
		if without.String() != with.String() {
			t.Errorf("%s changed with a block: %v vs %v", call, without, with)
		}
	}
}

func TestBlockFormReturnUndecidable(t *testing.T) {
	withSigs(t,
		Reg{Kind: "free", Name: "Dir", Sig: "(path :string ='*', files :boolean =false, details :boolean =false, block=false) :object|void"},
	)
	_, env := runPasses(`class { Foo(@args) { return Dir(@args) } }`, "T")
	got := env.LookupReturn("Foo")
	if !isUnionOf(got, TObject, TVoid) {
		t.Errorf("got %v, want object|void", got)
	}
}
