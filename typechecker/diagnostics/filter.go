// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package diagnostics

func FilterDiagnostics(diags []Diagnostic, cfg Config) []Diagnostic {
	out := diags[:0]
	for _, d := range diags {
		keep := true
		switch d.Flag {
		case FlagStrictStringConcat:
			keep = applyLevel(&d, cfg.StrictStringConcat)
		case FlagStrictCrossTypeCompares:
			keep = applyLevel(&d, cfg.StrictCrossTypeCompares)
		}
		if keep {
			out = append(out, d)
		}
	}
	return out
}

// returns false to drop, true to keep
func applyLevel(d *Diagnostic, l Level) bool {
	switch l {
	case LevelOff:
		return false
	case LevelWarn:
		d.Severity = SeverityWarning
	case LevelError:
		d.Severity = SeverityError
	}
	return true
}
