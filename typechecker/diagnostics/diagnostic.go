// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package diagnostics

import (
	"github.com/apmckinlay/gsuneido/typechecker/typealgebra"
)

type Severity int

const (
	SeverityWarning Severity = iota
	SeverityError
)

type Diagnostic struct {
	Severity   Severity
	Method     string
	Pos        int
	Msg        string
	Flag       Flag
	Got        []typealgebra.DynType
	Confidence float64
}
