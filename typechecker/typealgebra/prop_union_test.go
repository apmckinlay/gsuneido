// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package typealgebra_test

import (
	"testing"

	"github.com/apmckinlay/gsuneido/core"
	"pgregory.net/rapid"

	"github.com/apmckinlay/gsuneido/typechecker/internal/oracle"
	. "github.com/apmckinlay/gsuneido/typechecker/typealgebra"
)

// thin local names over the shared oracle so the properties read naturally.
func genDynType(t *rapid.T) DynType { return oracle.GenDynType(t) }
func semEq(a, b DynType) bool       { return oracle.SemEq(a, b) }
func dirtyOf(t DynType) bool        { return oracle.DirtyOf(t) }
func armCount(t DynType) int        { return oracle.ArmCount(t) }
func foldAny(t DynType) DynType     { return oracle.FoldAny(t) }

func canonArms(t DynType) (map[string]bool, bool) { return oracle.CanonArms(t) }

func TestPropFoldIdempotent(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		x := genDynType(rt)
		if !semEq(foldAny(x), foldAny(foldAny(x))) {
			rt.Fatalf("fold not idempotent: %v -> %v", foldAny(x), foldAny(foldAny(x)))
		}
	})
}

func TestPropNormalForm(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		oracle.AssertNormalForm(rt, genDynType(rt))
	})
}

func TestPropUCommutative(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		a, b := genDynType(rt), genDynType(rt)
		if !semEq(U(a, b), U(b, a)) {
			rt.Fatalf("U not commutative: U(%v,%v)=%v vs U(%v,%v)=%v", a, b, U(a, b), b, a, U(b, a))
		}
	})
}

func TestPropUAssociative(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		a, b, c := genDynType(rt), genDynType(rt), genDynType(rt)
		if !semEq(U(U(a, b), c), U(a, U(b, c))) {
			rt.Fatalf("U not associative for %v %v %v: %v vs %v",
				a, b, c, U(U(a, b), c), U(a, U(b, c)))
		}
	})
}

func TestPropUIdempotent(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		a := genDynType(rt)
		if !semEq(U(a, a), a) {
			rt.Fatalf("U not idempotent: U(%v,%v)=%v", a, a, U(a, a))
		}
	})
}

func TestPropUAbsorptive(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		a, b := genDynType(rt), genDynType(rt)
		if !semEq(U(a, U(a, b)), U(a, b)) {
			rt.Fatalf("U not absorptive: U(%v,U(%v,%v))=%v vs U(%v,%v)=%v",
				a, a, b, U(a, U(a, b)), a, b, U(a, b))
		}
	})
}

func TestPropUDirtyIff(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		a, b := genDynType(rt), genDynType(rt)
		want := dirtyOf(a) || dirtyOf(b)
		if got := dirtyOf(U(a, b)); got != want {
			rt.Fatalf("U(%v,%v)=%v dirty=%v want %v", a, b, U(a, b), got, want)
		}
	})
}

func TestPropUArmCountBound(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		a, b := genDynType(rt), genDynType(rt)
		if got, lim := armCount(U(a, b)), armCount(a)+armCount(b); got > lim {
			rt.Fatalf("|arms U(%v,%v)|=%d exceeds %d", a, b, got, lim)
		}
	})
}

func TestPropMarkDirty(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		x := genDynType(rt)
		md := MarkDirty(x)
		if !dirtyOf(md) {
			rt.Fatalf("MarkDirty(%v)=%v not dirty", x, md)
		}
		if !semEq(md, MarkDirty(md)) {
			rt.Fatalf("MarkDirty not idempotent: %v vs %v", md, MarkDirty(md))
		}
		xa, _ := canonArms(x)
		ma, _ := canonArms(md)
		if len(xa) != len(ma) {
			rt.Fatalf("MarkDirty changed arms: %v -> %v", x, md)
		}
		for k := range xa {
			if !ma[k] {
				rt.Fatalf("MarkDirty dropped arm %q: %v -> %v", k, x, md)
			}
		}
	})
}

func TestPropAllUnknownFoldsToUnknown(t *testing.T) {
	if U(TUnknown, TUnknown) != TUnknown {
		t.Fatalf("U(TUnknown,TUnknown)=%v want TUnknown", U(TUnknown, TUnknown))
	}
	if MarkDirty(TUnknown) != TUnknown {
		t.Fatalf("MarkDirty(TUnknown)=%v want TUnknown", MarkDirty(TUnknown))
	}
}

func TestPropStringCanonical(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		a, b := genDynType(rt), genDynType(rt)
		s1, s2 := U(a, b).String(), U(b, a).String()
		if s1 != s2 {
			rt.Fatalf("non-canonical string for semEq values: %q vs %q", s1, s2)
		}
	})
}

// TestPropDynTypeOfSuValue - total, and maps bool literals to TTrue/TFalse.
func TestPropDynTypeOfSuValue(t *testing.T) {
	cases := []struct {
		val  core.Value
		want DynType
	}{
		{core.True, TTrue},
		{core.False, TFalse},
		{core.SuStr("hello"), TString},
		{core.IntVal(7), TNumber},
	}
	for _, c := range cases {
		got := DynTypeOfSuValue(c.val)
		if got == nil {
			t.Fatalf("DynTypeOfSuValue(%v) returned nil", c.val)
		}
		if got != c.want {
			t.Fatalf("DynTypeOfSuValue(%v)=%v want %v", c.val, got, c.want)
		}
	}
}
