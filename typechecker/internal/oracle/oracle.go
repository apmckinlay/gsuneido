// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package oracle

import (
	"fmt"

	"pgregory.net/rapid"

	"github.com/apmckinlay/gsuneido/typechecker/typealgebra"
)

var atomTypes = []typealgebra.DynType{
	typealgebra.TVoid, typealgebra.TBoolean, typealgebra.TFalse, typealgebra.TTrue,
	typealgebra.TNumber, typealgebra.TString, typealgebra.TDate,
	typealgebra.TFunction, typealgebra.TBlock, typealgebra.TClass, typealgebra.TObject,
	typealgebra.TSequence,
	typealgebra.Instance{Class: "Ca"}, typealgebra.Instance{Class: "Cb"},
	typealgebra.Instance{Class: "Cc"},
	typealgebra.TUnknown,
}

func GenDynType(t *rapid.T) typealgebra.DynType {
	var rec func(d int) typealgebra.DynType
	rec = func(d int) typealgebra.DynType {
		if d <= 0 {
			return rapid.SampledFrom(atomTypes).Draw(t, "atom")
		}
		switch rapid.IntRange(0, 3).Draw(t, "op") {
		case 0:
			return rapid.SampledFrom(atomTypes).Draw(t, "atom")
		case 1:
			return typealgebra.U(rec(d-1), rec(d-1))
		case 2:
			return typealgebra.MarkDirty(rec(d - 1))
		default:
			return typealgebra.U(rec(d-1), typealgebra.U(rec(d-1), rec(d-1)))
		}
	}
	return rec(rapid.IntRange(0, 3).Draw(t, "depth"))
}

func CanonArms(t typealgebra.DynType) (arms map[string]bool, dirty bool) {
	arms = map[string]bool{}
	var add func(x typealgebra.DynType)
	add = func(x typealgebra.DynType) {
		switch v := x.(type) {
		case typealgebra.Primitive:
			if v == typealgebra.TUnknown {
				dirty = true
				return
			}
			arms[v.String()] = true
		case typealgebra.Instance:
			arms["#"+v.Class] = true
		case typealgebra.Union:
			if v.IsDirty {
				dirty = true
			}
			for _, a := range v.Types {
				add(a)
			}
		default:
			arms[fmt.Sprintf("%T:%v", x, x)] = true
		}
	}
	add(t)
	// mirror Fold's boolean lattice: any 2+ of {True,False,Boolean} -> Boolean
	n := 0
	for _, b := range []string{"true", "false", "boolean"} {
		if arms[b] {
			n++
		}
	}
	if n >= 2 {
		delete(arms, "true")
		delete(arms, "false")
		arms["boolean"] = true
	}
	return arms, dirty
}

func SemEq(a, b typealgebra.DynType) bool {
	aa, ad := CanonArms(a)
	bb, bd := CanonArms(b)
	if ad != bd || len(aa) != len(bb) {
		return false
	}
	for k := range aa {
		if !bb[k] {
			return false
		}
	}
	return true
}

func DirtyOf(t typealgebra.DynType) bool { _, d := CanonArms(t); return d }

func ArmCount(t typealgebra.DynType) int { a, _ := CanonArms(t); return len(a) }

// FoldAny normalizes a possibly-unfolded value (bare atoms pass through).
func FoldAny(t typealgebra.DynType) typealgebra.DynType {
	if u, ok := t.(typealgebra.Union); ok {
		return u.Fold()
	}
	return t
}

func AssertNormalForm(rt *rapid.T, f typealgebra.DynType) {
	u, ok := f.(typealgebra.Union)
	if !ok {
		return // a bare Primitive/Instance is trivially normal
	}
	boolCount := 0
	seen := map[string]bool{}
	for _, arm := range u.Types {
		if arm == typealgebra.TUnknown {
			rt.Fatalf("TUnknown arm in %v", f)
		}
		if _, nested := arm.(typealgebra.Union); nested {
			rt.Fatalf("nested Union arm in %v", f)
		}
		k := arm.String()
		if seen[k] {
			rt.Fatalf("duplicate arm %q in %v", k, f)
		}
		seen[k] = true
		if arm == typealgebra.TTrue || arm == typealgebra.TFalse || arm == typealgebra.TBoolean {
			boolCount++
		}
	}
	if boolCount > 1 {
		rt.Fatalf(">1 boolean arm in %v", f)
	}
	if u.IsDirty && len(u.Types) < 1 {
		rt.Fatalf("dirty union with 0 arms: %v", f)
	}
	if !u.IsDirty && len(u.Types) < 2 {
		rt.Fatalf("clean union with <2 arms: %v", f)
	}
}
