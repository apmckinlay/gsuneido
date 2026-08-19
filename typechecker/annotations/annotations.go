// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package annotations

import (
	"github.com/apmckinlay/gsuneido/typechecker/typealgebra"
)

type Set map[string][]Signature

type Param struct {
	Name       string
	Typ        typealgebra.DynType
	HasDefault bool
	Inferred   bool
	Why        string
}

type Signature struct {
	Receiver typealgebra.DynType // nil = unspecified
	Params   []Param
	AtParam  bool
	Returns  typealgebra.DynType
}

func receiverKey(r typealgebra.DynType) string {
	if r == nil {
		return ""
	}
	return r.String()
}
