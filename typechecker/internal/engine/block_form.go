// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
)

// A builtin that takes a block does one thing when given one and another when
// not - Dir returns its entries, or calls the block and returns nothing - so
// its annotation has to union both. The parser turns a trailing block into an
// argument named "block" (`f(){ }` parses as `Call(f block:Block())`), so the
// callsite says which half applies and the union can be dropped.
func blockFormReturn(sig *Signature, args []ast.Arg) DynType {
	u, ok := sig.Returns.(Union)
	if !ok || !u.Contains(TVoid) {
		return sig.Returns
	}
	slot := blockParamIndex(sig)
	if slot < 0 {
		return sig.Returns
	}
	supplied, certain := blockArgSupplied(slot, args)
	if !certain {
		return sig.Returns
	}
	if supplied {
		return TVoid
	}
	return unionWithoutVoid(u)
}

func blockParamIndex(sig *Signature) int {
	for i := range sig.Params {
		if sig.Params[i].Name == "block" {
			return i
		}
	}
	return -1
}

// certain is false for an args-spread, where the block may or may not be in
// there and the union has to stand
func blockArgSupplied(slot int, args []ast.Arg) (supplied, certain bool) {
	for i := range args {
		arg := &args[i]
		if isAtArg(arg) {
			return false, false
		}
		if arg.Name != nil {
			if name, ok := argName(arg); ok && name == "block" {
				return true, true
			}
			continue // some other named arg, keep looking
		}
		if i == slot { // positional args bind by index, as matchParam does
			supplied = true
		}
	}
	return supplied, true
}

func unionWithoutVoid(u Union) DynType {
	out := make([]DynType, 0, len(u.Types))
	for _, t := range u.Types {
		if t != TVoid {
			out = append(out, t)
		}
	}
	return Union{Types: out, IsDirty: u.IsDirty}.Fold()
}
