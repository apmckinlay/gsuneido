// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package diagnostics

import (
	"github.com/apmckinlay/gsuneido/typechecker/typealgebra"
)

func ScoreConfidence(d *Diagnostic) float64 {
	if d.Confidence != 0 {
		return d.Confidence
	}
	if d.Flag != FlagNone {
		return 0.20
	}
	switch d.Severity {
	case SeverityError:
		switch gotKind(d.Got) {
		case gotDirty:
			return 0.40 // provable-looking but built on a guess (`| ?`)
		case gotUnionArm:
			return 0.70 // proven for one arm of a clean union; narrowing may legitimately exclude it
		default:
			return 0.90 // fully concrete, or a proven error with no operand type (e.g. no-such-method)
		}
	case SeverityWarning:
		if gotKind(d.Got) == gotDirty {
			return 0.25
		}
		return 0.35
	}
	return 0.35
}

type gotClass int

const (
	gotConcrete gotClass = iota // single concrete type, or no operand-type info
	gotUnionArm                 // clean union with >1 arm - the bad type is one of several
	gotDirty                    // a dirty type (`| ?`): the checker is reporting a guess
)

func gotKind(got []typealgebra.DynType) gotClass {
	kind := gotConcrete
	for _, t := range got {
		switch {
		case isDirty(t):
			return gotDirty // dominates - report immediately
		case isCleanMultiUnion(t):
			kind = gotUnionArm
		}
	}
	return kind
}

func isDirty(t typealgebra.DynType) bool {
	if t == typealgebra.TUnknown {
		return true
	}
	u, ok := t.(typealgebra.Union)
	return ok && u.IsDirty
}

func isCleanMultiUnion(t typealgebra.DynType) bool {
	u, ok := t.(typealgebra.Union)
	return ok && !u.IsDirty && len(u.Types) > 1
}
