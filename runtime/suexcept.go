// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import "github.com/apmckinlay/gsuneido/runtime/types"

type SuExcept struct {
	SuStr
	Callstack *SuObject
}

func NewSuExcept(t *Thread, s SuStr) *SuExcept {
	return &SuExcept{SuStr: s, Callstack: t.Callstack()}
}

// SuValue interface ------------------------------------------------

func (*SuExcept) Type() types.Type {
	return types.Except
}

// SuExceptMethods is initialized by the builtin package
var SuExceptMethods Methods

func (*SuExcept) Lookup(t *Thread, method string) Callable {
	if m := SuExceptMethods[method]; m != nil {
		return m
	}
	return Lookup(t, StringMethods, gnStrings, method)
}
