// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"hash"
	"io"

	"crypto/sha1"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

type suSha1 struct {
	CantConvert
	hash hash.Hash
}

var _ = builtinRaw("Sha1(@args)",
	func(th *Thread, as *ArgSpec, args []Value) Value {
		sa := &suSha1{hash: sha1.New()}
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

var _ Value = (*suSha1)(nil)

func (*suSha1) Get(*Thread, Value) Value {
	panic("Sha1 does not support get")
}

func (*suSha1) Put(*Thread, Value, Value) {
	panic("Sha1 does not support put")
}

func (*suSha1) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic("Sha1 does not support update")
}

func (*suSha1) RangeTo(int, int) Value {
	panic("Sha1 does not support range")
}

func (*suSha1) RangeLen(int, int) Value {
	panic("Sha1 does not support range")
}

func (*suSha1) Hash() uint32 {
	panic("Sha1 hash not implemented")
}

func (*suSha1) Hash2() uint32 {
	panic("Sha1 hash not implemented")
}

func (*suSha1) Compare(Value) int {
	panic("Sha1 compare not implemented")
}

func (*suSha1) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call Sha1")
}

func (*suSha1) String() string {
	return "aSha1"
}

func (*suSha1) Type() types.Type {
	return types.BuiltinClass
}

func (sa *suSha1) Equal(other interface{}) bool {
	sa2, ok := other.(*suSha1)
	return ok && sa == sa2
}

func (*suSha1) SetConcurrent() {
	panic("Sha1 can not be shared between threads")
}

func (*suSha1) Lookup(_ *Thread, method string) Callable {
	return sha1Methods[method]
}

var sha1Methods = Methods{
	"Update": method1("(string)", func(this, arg Value) Value {
		io.WriteString(this.(*suSha1).hash, ToStr(arg))
		return this
	}),
	"Value": method0(func(this Value) Value {
		return this.(*suSha1).value()
	}),
}

func (sa *suSha1) value() Value {
	var buf [sha1.Size]byte
	return SuStr(string(sa.hash.Sum(buf[0:0])))
}
