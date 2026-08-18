// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/apmckinlay/gsuneido/typechecker/annotations"
)

// Reg is a signature registration, as TypeChecker.RegisterSignatures takes.
type Reg = annotations.Registration

// withSigs registers sigs for this test only, so each test states the builtin
// signatures it depends on instead of inheriting a shared table. The signature
// table is process-wide, so it is restored on cleanup and tests using this
// must not run in parallel.
func withSigs(t *testing.T, sigs ...Reg) {
	t.Helper()
	set, err := annotations.LoadRegistrations(sigs)
	if err != nil {
		t.Fatal(err)
	}
	prev := Annotations()
	SetAnnotations(set)
	t.Cleanup(func() { SetAnnotations(prev) })
}
