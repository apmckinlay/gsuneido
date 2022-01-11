// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"hash"
	"io"

	"hash/adler32"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

var _ = builtinRaw("Adler32(@args)",
	func(th *Thread, as *ArgSpec, args []Value) Value {
		sa := &suAdler32{hash: adler32.New()}
		iter := NewArgsIter(as, args)
		k, v := iter()
		if v == nil {
			return sa
		}
		for ; k == nil && v != nil; k, v = iter() {
			io.WriteString(sa.hash, ToStr(v))
		}
		return IntVal(int(int32(sa.hash.Sum32())))
	})

type suAdler32 struct {
	CantConvert
	hash hash.Hash32
}

var _ Value = (*suAdler32)(nil)

func (*suAdler32) Get(*Thread, Value) Value {
	panic("Adler32 does not support get")
}

func (*suAdler32) Put(*Thread, Value, Value) {
	panic("Adler32 does not support put")
}

func (*suAdler32) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic("Adler32 does not support update")
}

func (*suAdler32) RangeTo(int, int) Value {
	panic("Adler32 does not support range")
}

func (*suAdler32) RangeLen(int, int) Value {
	panic("Adler32 does not support range")
}

func (*suAdler32) Hash() uint32 {
	panic("Adler32 hash not implemented")
}

func (*suAdler32) Hash2() uint32 {
	panic("Adler32 hash not implemented")
}

func (*suAdler32) Compare(Value) int {
	panic("Adler32 compare not implemented")
}

func (*suAdler32) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call Adler32")
}

func (*suAdler32) String() string {
	return "aAdler32"
}

func (*suAdler32) Type() types.Type {
	return types.BuiltinClass
}

func (sa *suAdler32) Equal(other interface{}) bool {
	sa2, ok := other.(*suAdler32)
	return ok && sa == sa2
}

func (*suAdler32) SetConcurrent() {
	panic("Adler32 can not be shared between threads")
}

func (*suAdler32) Lookup(_ *Thread, method string) Callable {
	return adler32Methods[method]
}

var adler32Methods = Methods{
	"Update": method1("(string)", func(this, arg Value) Value {
		io.WriteString(this.(*suAdler32).hash, ToStr(arg))
		return this
	}),
	"Value": method0(func(this Value) Value {
		return IntVal(int(int32(this.(*suAdler32).hash.Sum32())))
	}),
}
