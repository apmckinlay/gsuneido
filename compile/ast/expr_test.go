// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/core"

	"github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestIsColumn(t *testing.T) {
	cols := []string{"a", "b", "c"}
	assert.T(t).True(IsColumn(&Ident{Name: "a"}, cols))
	assert.T(t).True(IsColumn(&Ident{Name: "c"}, cols))
	assert.T(t).True(IsColumn(&Ident{Name: "b_lower!"}, cols))
	assert.T(t).False(IsColumn(&Ident{Name: "x_lower!"}, cols))
	assert.T(t).False(IsColumn(&Ident{Name: "b_upper!"}, cols))
}

func TestExpr_CanEvalRaw(t *testing.T) {
	cols := []string{"a", "b", "c"}
	e := &Binary{
		Lhs: &Ident{Name: "b_lower!"},
		Tok: tokens.Is,
		Rhs: &Constant{Val: SuStr("foobar")},
	}
	assert.T(t).True(e.CanEvalRaw(cols))
	rec := new(RecordBuilder).Add(IntVal(0)).Add(SuStr("FooBar")).Build()
	ctx := &Context{
		Hdr: SimpleHeader(cols),
		Row: []DbRec{{Record: rec}},
	}
	assert.T(t).This(e.Eval(ctx)).Is(True)
}
