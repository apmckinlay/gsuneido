// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	crypto "crypto/rand"
	"math/rand"
	"time"

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
		return Int64Val(th.Rand.Int63n(1_0000_0000_0000_0000)) // dnum range
	}
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

var _ = builtin(RandomBytes, "(nbytes)")

func RandomBytes(arg Value) Value {
	n := ToInt(arg)
	if n < 0 || 128 < n {
		panic("RandomBytes: allowed range is 0 to 128")
	}
	buf := make([]byte, n)
	crypto.Read(buf)
	return SuStr(hacks.BStoS(buf))
}
