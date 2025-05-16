// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"fmt"
	"runtime/metrics"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/dbms"
	qry "github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/regex"
)

var _ = exportMethods(&SuneidoObjectMethods, "suneido")

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

var _ = staticMethod(suneido_ParseQuery, "(query)")

func suneido_ParseQuery(th *Thread, args []Value) Value {
	dbms, ok := th.Dbms().(*dbms.DbmsLocal)
	if !ok {
		panic("Suneido.ParseQuery requires a local database")
	}
	t := dbms.Transaction(false)
	defer t.Complete()
	query := ToStr(args[0])
	q := qry.JustParse(t.(qry.QueryTran), query)
	return qry.NewSuQueryNode(q)
}

var _ = staticMethod(suneido_Regex, "(pattern)")

func suneido_Regex(arg Value) Value {
	return SuRegex{Pat: regex.Compile(ToStr(arg))}
}

var _ = staticMethod(suneido_GoMetric, "(name = false)")

func suneido_GoMetric(th *Thread, args []Value) Value {
	if args[0] == False {
		all := metrics.All()
		ob := SuObject{}
		for _, d := range all {
			ob.Add(SuStr(d.Name))
		}
		return &ob
	}
	sample := make([]metrics.Sample, 1)
	sample[0].Name = ToStr(args[0])
	metrics.Read(sample)
	switch sample[0].Value.Kind() {
	case metrics.KindUint64:
		return Int64Val(int64(sample[0].Value.Uint64()))
	case metrics.KindFloat64:
		return SuDnum{Dnum: dnum.FromFloat(sample[0].Value.Float64())}
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

var Nil []Value

func suneido_RuntimeError() Value {
	return Nil[123]
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

var _ = staticMethod(suneido_WarningsThrow, "(arg = true)")

func suneido_WarningsThrow(x Value) Value {
	switch x {
	case True:
		options.WarningsThrow.Store(options.AllWarningsThrow)
	case False:
		options.WarningsThrow.Store(options.NoWarningsThrow)
	default:
		options.WarningsThrow.Store(regex.Compile(ToStr(x)))
	}
	return nil
}

var _ = staticMethod(suneido_Info, "(name = false)")

func suneido_Info(x Value) Value {
	if x == False {
		return SuObjectOfStrs(InfoList())
	}
	return InfoStr(ToStr(x))
}

var _ = method(suneido_Members, "(all = false)")

func suneido_Members(this Value, all Value) Value {
	if !ToBool(all) {
		return NewSuSequence(IterMembers(ToContainer(this), true, true))
	}
	suneido := this.(*SuneidoObject)
	mems := make([]Value, 0, suneido.Size()+len(SuneidoObjectMethods))
	iter := IterMembers(suneido, true, true)
	for v := iter.Next(); v != nil; v = iter.Next() {
		mems = append(mems, v)
	}
	for k := range SuneidoObjectMethods {
		mems = append(mems, SuStr(k))
	}
	return SuObjectOf(mems...)
}

var _ = staticMethod(suneido_IndexUse, "()")

func suneido_IndexUse() Value {
	iu := qry.PullIdxUse()
	ob := &SuObject{}
	for k, v := range iu {
		ob.Set(SuStr(k), IntVal(v))
	}
	return ob
}

var _ = staticMethod(suneido_LibraryTags, "(@args)")

func suneido_LibraryTags(args Value) Value {
	ob := args.(*SuObject)
	tags := make([]string, 1+ob.ListSize())
	tags[0] = "" // untagged
	for i := range tags[1:] {
		tags[i+1] = "__" + ToStr(ob.ListGet(i))
	}
	options.LibraryTags = tags
	Global.UnloadAll() // same as Use/Unuse
	return nil
}

var _ = AddInfo("library.tags",
	func() string { return fmt.Sprintf("%#v", options.LibraryTags)[8:] })
