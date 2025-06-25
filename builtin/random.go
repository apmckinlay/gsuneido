// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	crypto "crypto/rand"
	"math/rand/v2"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

type suRandomGlobal struct {
	SuBuiltin
}

func init() {
	Global.Builtin("Random", &suRandomGlobal{
		SuBuiltin{Fn: Random,
			BuiltinParams: BuiltinParams{ParamSpec: params("(limit = false)")}}})
}

func Random(th *Thread, args []Value) Value {
	initRand(th)
	if args[0] == False {
		return Int64Val(th.Rand.Int64())
	}
	limit := IfInt(args[0])
	return IntVal(th.Rand.IntN(limit))
}

func initRand(th *Thread) {
	if th.Rand == nil {
		th.Rand = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	}
}

var randomMethods = methods("rnd")

var _ = staticMethod(rnd_Seed, "(seed)")

func rnd_Seed(th *Thread, args []Value) Value {
	seed := uint64(IfInt(args[0]))
	// using the same value twice is not ideal, but it's all we have
	th.Rand = rand.New(rand.NewPCG(seed, seed))
	return nil
}

var _ = staticMethod(rnd_Bytes, "(nbytes)")

func rnd_Bytes(arg Value) Value {
	n := ToInt(arg)
	if n < 0 || 128 < n {
		panic("Random.Bytes: allowed range is 0 to 128")
	}
	buf := make([]byte, n)
	crypto.Read(buf)
	return SuStr(hacks.BStoS(buf))
}

var _ = staticMethod(rnd_Members, "()")

func rnd_Members() Value {
	return rnd_members
}

var rnd_members = methodList(randomMethods)

func (r *suRandomGlobal) Lookup(th *Thread, method string) Value {
	if f, ok := randomMethods[method]; ok {
		return f
	}
	return r.SuBuiltin.Lookup(th, method) // for Params
}
