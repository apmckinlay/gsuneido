// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package typechecker

import (
	"strings"
	"testing"
)

func TestDirAnnotationDropsVoid(t *testing.T) {
	withSigs(t,
		Reg{Kind: "free", Name: "Dir", Sig: "(path :string ='*', files :boolean =false, details :boolean =false, block=false) :object|void"},
	)
	src := `class { F() { d = Dir("asd") } }`
	res := runRequest(t, "TypeAnnotate", nil,
		[]SourceEntry{{Name: "T", Src: src}})
	got, ok := res.Results[0].(string)
	if !ok {
		t.Fatalf("want annotated source, got %T", res.Results[0])
	}
	if strings.Contains(got, "object|void") {
		t.Errorf("void still leaking into a blockless Dir call:\n%s", got)
	}
	if !strings.Contains(got, "/* object */") {
		t.Errorf("want d annotated as object:\n%s", got)
	}
}
