// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import "testing"

func TestWidenBool(t *testing.T) {
	got := widenBool(TFalse)
	u, ok := got.(Union)
	if !ok || !u.IsDirty || len(u.Types) != 1 || u.Types[0] != TFalse {
		t.Errorf("widenBool(TFalse) = %v, want dirty False", got)
	}
	if got := widenBool(TTrue); got != TBoolean {
		t.Errorf("widenBool(TTrue) = %v, want Boolean", got)
	}
	for _, ty := range []DynType{TNumber, TString, TBoolean, TObject, TDate, TUnknown} {
		if got := widenBool(ty); got != ty {
			t.Errorf("widenBool(%v) = %v, want unchanged", ty, got)
		}
	}
}

func TestSeedWidensOnlyPublicMembers(t *testing.T) {
	src := `class { Pub: false
priv: false
Flag: true
quiet: true }`
	cls, ok := safeParse(src, "T")
	if !ok {
		t.Fatal("class did not parse")
	}
	env := NewTypeEnv()
	LocalInference(cls, env, &PassCtx{})

	if got := env.Members["Pub"]; func() bool {
		u, ok := got.(Union)
		return !ok || !u.IsDirty || len(u.Types) != 1 || u.Types[0] != TFalse
	}() {
		t.Errorf("public false seed = %v, want dirty False", env.Members["Pub"])
	}
	if got := env.Members["priv"]; got != TFalse {
		t.Errorf("private false seed = %v, want False", got)
	}
	if got := env.Members["Flag"]; got != TBoolean {
		t.Errorf("public true seed = %v, want Boolean", got)
	}
	if got := env.Members["quiet"]; got != TTrue {
		t.Errorf("private true seed = %v, want True", got)
	}
}
