// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"fmt"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/bloom"
)

type suBloom struct {
	ValueBase[suBloom]
	data *bloom.Bloom
}

var _ = builtin(Bloom, "(n, p)")

func Bloom(n, p Value) Value {
	m, k := bloom.Calc(ToInt(n), ToDnum(p).ToFloat())
	if m/8 > StringLimit {
		panic(fmt.Sprint("Bloom size ", m/8, " > max ", StringLimit))
	}
	return suBloom{data: bloom.New(m, k)}
}

var _ Value = (*suBloom)(nil)

func (b suBloom) Equal(other any) bool {
	return b == other
}

func (suBloom) Lookup(_ *Thread, method string) Value {
	return bloomMethods[method]
}

var bloomMethods = methods("bloom")

var _ = method(bloom_Add, "(value)")

func bloom_Add(this, arg Value) Value {
	this.(suBloom).data.Add(arg.Hash())
	return this
}

var _ = method(bloom_Test, "(value)")

func bloom_Test(this, arg Value) Value {
	return SuBool(this.(suBloom).data.Test(arg.Hash()))
}

var _ = method(bloom_Size, "()")

func bloom_Size(this Value) Value {
	return IntVal(this.(suBloom).data.Size())
}
