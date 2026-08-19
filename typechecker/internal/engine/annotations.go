// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/compile/ast"

	"github.com/apmckinlay/gsuneido/typechecker/annotations"
	"github.com/apmckinlay/gsuneido/typechecker/typealgebra"
)

type (
	AnnotationSet = annotations.Set
	Signature     = annotations.Signature
	Param         = annotations.Param
)

var builtinAnnotations atomic.Pointer[AnnotationSet]

// SetAnnotations atomically publishes a replacement table, built by
// annotations.LoadRegistrations. Runs in flight keep the snapshot they
// took at NewPassCtx; new runs see the new table.
func SetAnnotations(set AnnotationSet) {
	builtinAnnotations.Store(&set)
}

// shared, not copied - passes only read it, which
// TestPropSharedAnnotationsNotMutated enforces. Tables are immutable once
// published; SetAnnotations replaces the pointer atomically. Empty until
// signatures are registered - in the exe, builtin pushes its own
// declarations at startup, so only tests ever see it empty.
func Annotations() AnnotationSet {
	if p := builtinAnnotations.Load(); p != nil {
		return *p
	}
	return AnnotationSet{}
}

// SignaturesRegistered reports whether a table has been published. Registering
// an empty table is a legitimate answer ("this exe has no builtins to declare"),
// so callers must not infer "never registered" from an empty Annotations().
func SignaturesRegistered() bool {
	return builtinAnnotations.Load() != nil
}

func ParseTypeAnnotation(s string) (DynType, error) {
	return typealgebra.ParseAnnotation(s)
}

func signatureFromAst(fn *ast.Function) (Signature, error) {
	return annotations.SignatureFromAst(fn)
}
