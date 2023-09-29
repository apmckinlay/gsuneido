// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/options"
)

var _ = exportMethods(&ClassMethods)

func init() {
	ClassMethods["*new*"] =
		&SuBuiltinMethodRaw{Fn: class_new, ParamSpec: params("(@args)")}
}

func class_new(th *Thread, as *ArgSpec, this Value, args []Value) Value {
	return this.(*SuClass).New(th, as)
}

var _ = method(class_BaseQ, "(class)")

func class_BaseQ(th *Thread, this Value, args []Value) Value {
	return nilToFalse(this.(*SuClass).Finder(th,
		func(v Value, _ *MemBase) Value {
			if v == args[0] {
				return True
			}
			return nil
		}))
}

var _ = method(class_MethodQ, "(string)")

func class_MethodQ(th *Thread, this Value, args []Value) Value {
	m := ToStr(args[0])
	return nilToFalse(this.(Findable).Finder(th, func(c Value, mb *MemBase) Value {
		if _, ok := c.(*SuClass); ok {
			if x, ok := mb.Data[m]; ok && x.Type() == types.Function {
				return True
			}
		}
		return nil
	}))
}

var _ = method(class_MethodClass, "(string)")

func class_MethodClass(th *Thread, this Value, args []Value) Value {
	m := ToStr(args[0])
	return nilToFalse(this.(Findable).Finder(th, func(c Value, mb *MemBase) Value {
		if _, ok := c.(*SuClass); ok {
			if x, ok := mb.Data[m]; ok && x.Type() == types.Function {
				return c
			}
		}
		return nil
	}))
}

var _ = method(class_ReadonlyQ, "()")

func class_ReadonlyQ(this Value) Value {
	return True
}

var _ = method(class_StartCoverage, "(count = false)")

func class_StartCoverage(this, a Value) Value {
	if !options.Coverage.Load() {
		panic("coverage not enabled")
	}
	c := this.(*SuClass)
	c.StartCoverage(ToBool(a))
	return nil
}

var _ = method(class_StopCoverage, "()")

func class_StopCoverage(this Value) Value {
	c := this.(*SuClass)
	return c.StopCoverage()
}
