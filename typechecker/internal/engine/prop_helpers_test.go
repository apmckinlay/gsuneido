// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"
	"sort"

	"pgregory.net/rapid"

	"github.com/apmckinlay/gsuneido/typechecker/internal/oracle"
)

func semEq(a, b DynType) bool { return oracle.SemEq(a, b) }

func assertNormalForm(rt *rapid.T, f DynType) { oracle.AssertNormalForm(rt, f) }

func fits(a, b DynType) bool {
	bMembers, bDirty := decomposeForCheck(b)
	if bDirty {
		return true
	}
	aMembers, aDirty := decomposeForCheck(a)
	if aDirty {
		return false
	}
	for _, m := range aMembers {
		ok := false
		for _, bm := range bMembers {
			if memberFits(m, bm) {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

// --- diagnostic / env comparison helpers (Property B onward) ---

// diagKey is the full identity of a diagnostic (severity, method, pos, msg).
func diagKey(d Diagnostic) string {
	return fmt.Sprintf("%d\x00%s\x00%d\x00%s", d.Severity, d.Method, d.Pos, d.Msg)
}

// diagKeyNoPos drops Pos - for reorder/comment transforms where positions shift.
func diagKeyNoPos(d Diagnostic) string {
	return fmt.Sprintf("%d\x00%s\x00%s", d.Severity, d.Method, d.Msg)
}

func diagList(env TypeEnv) []Diagnostic {
	if env.Diagnostics == nil {
		return nil
	}
	return *env.Diagnostics
}

// sortedDiagKeys renders a diagnostic slice as a sorted key multiset.
func sortedDiagKeys(ds []Diagnostic, key func(Diagnostic) string) []string {
	ks := make([]string, 0, len(ds))
	for _, d := range ds {
		ks = append(ks, key(d))
	}
	sort.Strings(ks)
	return ks
}

func strSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// typeMapsSemEq compares two name->DynType maps by keys and semantic equality.
func typeMapsSemEq(a, b map[string]DynType) bool {
	if len(a) != len(b) {
		return false
	}
	for k, va := range a {
		vb, ok := b[k]
		if !ok || !semEq(va, vb) {
			return false
		}
	}
	return true
}

func typeMapKeysSubset(sub, super map[string]DynType) bool {
	for k := range sub {
		if _, ok := super[k]; !ok {
			return false
		}
	}
	return true
}
