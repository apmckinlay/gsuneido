// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"hash"
	"io"

	"crypto/sha256"

	. "github.com/apmckinlay/gsuneido/runtime"
)

type suSha256 struct {
	ValueBase[*suSha256]
	hash hash.Hash
}

var _ = builtinRaw("Sha256(@args)",
	func(th *Thread, as *ArgSpec, args []Value) Value {
		sa := &suSha256{hash: sha256.New()}
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

var _ Value = (*suSha256)(nil)

func (sa *suSha256) Equal(other interface{}) bool {
	sa2, ok := other.(*suSha256)
	return ok && sa == sa2
}

func (*suSha256) Lookup(_ *Thread, method string) Callable {
	return sha256Methods[method]
}

var sha256Methods = Methods{
	"Update": method1("(string)", func(this, arg Value) Value {
		io.WriteString(this.(*suSha256).hash, ToStr(arg))
		return this
	}),
	"Value": method0(func(this Value) Value {
		return this.(*suSha256).value()
	}),
}

func (sa *suSha256) value() Value {
	var buf [sha256.Size]byte
	return SuStr(string(sa.hash.Sum(buf[0:0])))
}
