// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/core"
)

// ```suneido
// d = Date('#20240101')
//
//	^^^^^^^^^^^^^^^^^ literal always parses -> Date (no False arm)
//
// e = Date('junk')
//
//	^^^^^^^^^^^^ literal never parses -> False
//
// ```
func DateNarrowingPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	for _, fn := range cls.SortedMethods {
		for _, stmt := range fn.Body {
			walkDateCalls(stmt, env)
		}
	}
	return false
}

type dateProof int

const (
	dateUnprovable  dateProof = iota // non-constant arg, or runtime panic - we know nothing
	dateProvenValid                  // the runtime builds a real date         -> Date
	dateProvenFalse                  // the runtime constructor returns false   -> False
)

func walkDateCalls(n ast.Node, env TypeEnv) {
	if n == nil {
		return
	}
	if call, ok := n.(*ast.Call); ok {
		if id, ok := call.Fn.(*ast.Ident); ok && id.Name == "Date" {
			switch proveDateCall(call) {
			case dateProvenValid:
				env.MarkValidDate(call)
			case dateProvenFalse:
				env.MarkFalseDate(call)
			}
		}
	}
	n.Children(func(c ast.Node) ast.Node {
		walkDateCalls(c, env)
		return c
	})
}

var dateArgSlot = map[string]int{
	"string": 0, "pattern": 1,
	"year": 2, "month": 3, "day": 4,
	"hour": 5, "minute": 6, "second": 7, "millisecond": 8,
}

func proveDateCall(call *ast.Call) dateProof {
	var slot [9]core.Value
	var has [9]bool
	pos := 0
	for i := range call.Args {
		a := &call.Args[i]
		c, ok := unwrapConstant(a.E)
		if !ok {
			return dateUnprovable
		}
		var idx int
		if a.Name != nil {
			name, ok := a.Name.(core.SuStr)
			if !ok {
				return dateUnprovable
			}
			idx, ok = dateArgSlot[string(name)]
			if !ok {
				return dateUnprovable
			}
		} else {
			idx = pos
			pos++
		}
		if idx >= len(slot) {
			return dateUnprovable
		}
		slot[idx] = c.Val
		has[idx] = true
	}
	return evalConstDate(&slot, has)
}

func evalConstDate(slot *[9]core.Value, has [9]bool) (proof dateProof) {
	defer func() {
		if recover() != nil {
			proof = dateUnprovable
		}
	}()
	fields := false
	for i := 2; i <= 8; i++ {
		if has[i] {
			fields = true
		}
	}
	if has[0] {
		if fields {
			return dateUnprovable // Date(string, year:..) panics at runtime - undecidable
		}
		return evalDateString(slot, has)
	}
	if fields {
		return evalDateFields(slot, has)
	}
	return dateProvenValid // Date() -> Now(), always valid
}

func evalDateString(slot *[9]core.Value, has [9]bool) dateProof {
	v := slot[0]
	if isDateValue(v) {
		return dateProvenValid
	}
	s := core.AsStr(v)
	if strings.HasPrefix(s, "#") || isTimestampLiteral(s) {
		return dateProofOf(core.DateFromLiteral(s) != core.NilDate)
	}
	order := "yMd"
	if has[1] {
		order = core.AsStr(slot[1])
	}
	return dateProofOf(core.ParseDate(s, order) != core.NilDate)
}

func evalDateFields(slot *[9]core.Value, has [9]bool) dateProof {
	year, month, day := 1700, 1, 1
	hour, minute, second, ms := 0, 0, 0, 0
	if has[2] {
		year = core.ToInt(slot[2])
	}
	if has[3] {
		month = core.ToInt(slot[3])
	}
	if has[4] {
		day = core.ToInt(slot[4])
	}
	if has[5] {
		hour = core.ToInt(slot[5])
	}
	if has[6] {
		minute = core.ToInt(slot[6])
	}
	if has[7] {
		second = core.ToInt(slot[7])
	}
	if has[8] {
		ms = core.ToInt(slot[8])
	}
	return dateProofOf(core.NormalizeDate(year, month, day, hour, minute, second, ms) != core.NilDate)
}

func dateProofOf(valid bool) dateProof {
	if valid {
		return dateProvenValid
	}
	return dateProvenFalse
}

func isDateValue(v core.Value) bool {
	switch v.(type) {
	case core.SuDate, core.SuTimestamp:
		return true
	}
	return false
}

func isTimestampLiteral(s string) bool {
	if len(s) != 21 || s[8] != '.' {
		return false
	}
	for i := 0; i < len(s); i++ {
		if i == 8 {
			continue
		}
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
