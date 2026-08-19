// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"testing"

	"pgregory.net/rapid"

	"github.com/apmckinlay/gsuneido/typechecker/internal/oracle"
)

func TestPropUArmsFit(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		a, b := oracle.GenDynType(rt), oracle.GenDynType(rt)
		u := U(a, b)
		if !fits(a, u) {
			rt.Fatalf("operand %v does not fit U(a,b)=%v", a, u)
		}
		if !fits(b, u) {
			rt.Fatalf("operand %v does not fit U(a,b)=%v", b, u)
		}
	})
}
