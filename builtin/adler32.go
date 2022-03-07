// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"hash"
	"io"

	"hash/adler32"

	. "github.com/apmckinlay/gsuneido/runtime"
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
	ValueBase[*suAdler32]
	hash hash.Hash32
}

var _ Value = (*suAdler32)(nil)

func (sa *suAdler32) Equal(other interface{}) bool {
	sa2, ok := other.(*suAdler32)
	return ok && sa == sa2
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
