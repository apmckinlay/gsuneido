// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"hash"
	"io"

	"crypto/md5"

	. "github.com/apmckinlay/gsuneido/runtime"
)

type suMd5 struct {
	ValueBase[*suMd5]
	hash hash.Hash
}

var _ = builtinRaw("Md5(@args)",
	func(th *Thread, as *ArgSpec, args []Value) Value {
		sa := &suMd5{hash: md5.New()}
		iter := NewArgsIter(as, args)
		k, v := iter()
		if v == nil {
			return sa
		}
		for ; k == nil && v != nil; k, v = iter() {
			io.WriteString(sa.hash, ToStr(v))
		}
		return sa.value()
	})

var _ Value = (*suMd5)(nil)

func (sa *suMd5) Equal(other interface{}) bool {
	sa2, ok := other.(*suMd5)
	return ok && sa == sa2
}

func (*suMd5) Lookup(_ *Thread, method string) Callable {
	return md5Methods[method]
}

var md5Methods = Methods{
	"Update": method1("(string)", func(this, arg Value) Value {
		io.WriteString(this.(*suMd5).hash, ToStr(arg))
		return this
	}),
	"Value": method0(func(this Value) Value {
		return this.(*suMd5).value()
	}),
}

func (sa *suMd5) value() Value {
	var buf [md5.Size]byte
	return SuStr(string(sa.hash.Sum(buf[0:0])))
}
