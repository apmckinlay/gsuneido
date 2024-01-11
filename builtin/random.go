// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"math/rand"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
)

type suRandomGlobal struct {
	SuBuiltin
}

func init() {
	Global.Builtin("Random", &suRandomGlobal{
		SuBuiltin{Fn: Random,
			BuiltinParams: BuiltinParams{ParamSpec: params("(limit)")}}})
}

func Random(th *Thread, args []Value) Value {
	initRand(th)
	limit := IfInt(args[0])
	return IntVal(th.Rand.Intn(limit))
}

func initRand(th *Thread) {
	if th.Rand == nil {
		th.Rand = rand.New(rand.NewSource(time.Now().UnixNano() * rand.Int63()))
	}
}

var randomMethods = methods()

var _ = staticMethod(rnd_Seed, "(seed)")

func rnd_Seed(th *Thread, args []Value) Value {
	initRand(th)
	th.Rand.Seed(int64(IfInt(args[0])))
	return nil
}

func (d *suRandomGlobal) Lookup(th *Thread, method string) Callable {
	if f, ok := randomMethods[method]; ok {
		return f
	}
	return d.SuBuiltin.Lookup(th, method) // for Params
}
