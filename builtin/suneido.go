// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"runtime/metrics"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util2/regex"
)

var _ = exportMethods(&SuneidoObjectMethods)

var _ = staticMethod(suneido_Compile, "(source, errob = false)")

func suneido_Compile(th *Thread, args []Value) Value {
	src := ToStr(args[0])
	if args[1] == False {
		return compile.Constant(src)
	}
	ob := ToContainer(args[1])
	val, checks := compile.Checked(th, src)
	for _, w := range checks {
		ob.Add(SuStr(w))
	}
	return val
}

var _ = staticMethod(suneido_Parse, "(source)")

func suneido_Parse(th *Thread, args []Value) Value {
	src := ToStr(args[0])
	p := compile.AstParser(src)
	ast := p.Const()
	if p.Token != tokens.Eof {
		p.Error("did not parse all input")
	}
	return ast
}

var _ = staticMethod(suneido_Regex, "(pattern)")

func suneido_Regex(arg Value) Value {
	return SuRegex{Pat: regex.Compile(ToStr(arg))}
}

var _ = staticMethod(suneido_GoMetric, "(name)")

func suneido_GoMetric(th *Thread, args []Value) Value {
	sample := make([]metrics.Sample, 1)
	sample[0].Name = ToStr(args[0])
	metrics.Read(sample)
	switch sample[0].Value.Kind() {
	case metrics.KindUint64:
		return Int64Val(int64(sample[0].Value.Uint64()))
	case metrics.KindFloat64:
		return SuDnum{Dnum: dnum.FromFloat(float64(sample[0].Value.Float64()))}
	default:
		return False
	}
}

// force various kinds of errors for testing

var _ = staticMethod(suneido_CrashX, "()")

func suneido_CrashX() Value {
	go func() { panic("Crash!") }()
	return nil
}

var _ = staticMethod(suneido_AssertFail, "()")

func suneido_AssertFail() Value {
	assert.Msg("Suneido.AssertFail").That(false)
	return nil
}

var _ = staticMethod(suneido_ShouldNotReachHere, "()")

func suneido_ShouldNotReachHere() Value {
	panic(assert.ShouldNotReachHere())
}

var _ = staticMethod(suneido_RuntimeError, "()")

func suneido_RuntimeError() Value {
	var x []Value
	return x[123]
}

var _ = staticMethod(suneido_StrictCompare, "(bool)")

func suneido_StrictCompare(x Value) Value {
	options.StrictCompare = ToBool(x)
	return nil
}

var _ = staticMethod(suneido_StrictCompareDb, "(bool)")

func suneido_StrictCompareDb(x Value) Value {
	options.StrictCompareDb = ToBool(x)
	return nil
}
