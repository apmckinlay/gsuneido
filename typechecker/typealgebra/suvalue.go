// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package typealgebra

import (
	"strings"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"
)

func opString(val core.Value) string {
	op := val.Get(nil, core.SuStr("op"))
	if op == nil {
		return ""
	}
	return strings.Trim(op.String(), "\"")
}

func DynTypeOfSuValue(val core.Value) DynType {
	if b, ok := val.(core.SuBool); ok {
		switch b.String() {
		case "false":
			return TFalse
		case "true":
			return TTrue
		default:
			return TBoolean
		}
	}

	ty := val.Type()
	switch ty {
	case types.Number:
		return TNumber
	case types.String:
		return TString
	case types.Date:
		return TDate
	case types.Boolean:
		return TBoolean
	case types.AstNode:
		if t := val.Get(nil, core.SuStr("type")); t != nil {
			tystr := strings.Trim(t.String(), "\"")
			switch tystr {
			case "Object", "Record":
				return TObject
			case "Class":
				return TClass
			case "Function":
				return TFunction
			case "Nary", "Binary", "Unary":
				return ResultTypeOfOp(opString(val))
			}
		}
		// unhandled AstNode kind - no static type for it
		return TUnknown
	default:
		// unhandled SuValue type - no static type for it
		return TUnknown
	}
}
