// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"hash"
	"io"

	"crypto/md5"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

type suMd5 struct {
	CantConvert
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

func (*suMd5) Get(*Thread, Value) Value {
	panic("Md5 does not support get")
}

func (*suMd5) Put(*Thread, Value, Value) {
	panic("Md5 does not support put")
}

func (*suMd5) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic("Md5 does not support update")
}

func (*suMd5) RangeTo(int, int) Value {
	panic("Md5 does not support range")
}

func (*suMd5) RangeLen(int, int) Value {
	panic("Md5 does not support range")
}

func (*suMd5) Hash() uint32 {
	panic("Md5 hash not implemented")
}

func (*suMd5) Hash2() uint32 {
	panic("Md5 hash not implemented")
}

func (*suMd5) Compare(Value) int {
	panic("Md5 compare not implemented")
}

func (*suMd5) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call Md5")
}

func (*suMd5) String() string {
	return "aMd5"
}

func (*suMd5) Type() types.Type {
	return types.BuiltinClass
}

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
