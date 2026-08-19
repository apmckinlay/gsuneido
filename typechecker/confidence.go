// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package typechecker

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/typechecker/diagnostics"
)

type rankedDiag struct {
	severity   diagnostics.Severity
	confidence float64
	entry      ResultDiagnostic
}

type confidenceFilter struct {
	active bool
	keep   func(c float64) bool
}

func (f confidenceFilter) passes(c float64) bool {
	return !f.active || f.keep(c)
}

func parseConfidenceFilter(raw map[string]string) (confidenceFilter, error) {
	s, ok := raw["confidence"]
	if !ok {
		return confidenceFilter{}, nil
	}
	s = strings.TrimSpace(s)
	if s == "" || strings.EqualFold(s, "all") {
		return confidenceFilter{}, nil
	}
	op := ">=" // a bare number is a minimum threshold
	for _, cand := range []string{">=", "<=", "==", "=", ">", "<"} {
		if strings.HasPrefix(s, cand) {
			op = cand
			s = strings.TrimSpace(s[len(cand):])
			break
		}
	}
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return confidenceFilter{}, fmt.Errorf("config.confidence: invalid threshold %q (want e.g. \">=0.70\")", raw["confidence"])
	}
	var keep func(float64) bool
	switch op {
	case ">=":
		keep = func(c float64) bool { return c >= n }
	case ">":
		keep = func(c float64) bool { return c > n }
	case "<=":
		keep = func(c float64) bool { return c <= n }
	case "<":
		keep = func(c float64) bool { return c < n }
	case "==", "=":
		keep = func(c float64) bool { return c == n }
	}
	return confidenceFilter{active: true, keep: keep}, nil
}

func rankDiagnostics(collected []rankedDiag, filter confidenceFilter) DiagnosticSet {
	sort.SliceStable(collected, func(i, j int) bool {
		return collected[i].confidence > collected[j].confidence
	})
	diags := DiagnosticSet{
		Errors:   []ResultDiagnostic{},
		Warnings: []ResultDiagnostic{},
	}
	for _, rd := range collected {
		if !filter.passes(rd.confidence) {
			continue
		}
		e := rd.entry
		e.Msg = fmt.Sprintf("[%.2f] %s", rd.confidence, e.Msg)
		switch rd.severity {
		case diagnostics.SeverityError:
			diags.Errors = append(diags.Errors, e)
		case diagnostics.SeverityWarning:
			diags.Warnings = append(diags.Warnings, e)
		}
	}
	return diags
}
