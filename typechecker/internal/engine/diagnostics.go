// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"github.com/apmckinlay/gsuneido/typechecker/diagnostics"
)

// aliases so the passes don't each import the diagnostics package
type (
	Diagnostic = diagnostics.Diagnostic
	Severity   = diagnostics.Severity
	Flag       = diagnostics.Flag
	Config     = diagnostics.Config
)

const (
	SeverityWarning = diagnostics.SeverityWarning
	SeverityError   = diagnostics.SeverityError

	FlagNone                    = diagnostics.FlagNone
	FlagStrictStringConcat      = diagnostics.FlagStrictStringConcat
	FlagStrictCrossTypeCompares = diagnostics.FlagStrictCrossTypeCompares
)

func ScoreConfidence(d *Diagnostic) float64 { return diagnostics.ScoreConfidence(d) }
