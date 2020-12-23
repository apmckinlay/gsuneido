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

type SuMd5 struct {
	CantConvert
	hash hash.Hash
}

var _ = builtinRaw("Md5(@args)",
	func(th *Thread, as *ArgSpec, args []Value) Value {
		sa := &SuMd5{hash: md5.New()}
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

var _ Value = (*SuMd5)(nil)

func (*SuMd5) Get(*Thread, Value) Value {
	panic("Md5 does not support get")
}

func (*SuMd5) Put(*Thread, Value, Value) {
	panic("Md5 does not support put")
}

func (*SuMd5) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic("Md5 does not support update")
}

func (*SuMd5) RangeTo(int, int) Value {
	panic("Md5 does not support range")
}

func (*SuMd5) RangeLen(int, int) Value {
	panic("Md5 does not support range")
}

func (*SuMd5) Hash() uint32 {
	panic("Md5 hash not implemented")
}

func (*SuMd5) Hash2() uint32 {
	panic("Md5 hash not implemented")
}

func (*SuMd5) Compare(Value) int {
	panic("Md5 compare not implemented")
}

func (*SuMd5) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call Md5")
}

func (*SuMd5) String() string {
	return "aMd5"
}

func (*SuMd5) Type() types.Type {
	return types.BuiltinClass
}

func (sa *SuMd5) Equal(other interface{}) bool {
	sa2, ok := other.(*SuMd5)
	return ok && sa == sa2
}

func (*SuMd5) Lookup(_ *Thread, method string) Callable {
	return md5Methods[method]
}

var md5Methods = Methods{
	"Update": method1("(string)", func(this, arg Value) Value {
		io.WriteString(this.(*SuMd5).hash, ToStr(arg))
		return this
	}),
	"Value": method0(func(this Value) Value {
		return this.(*SuMd5).value()
	}),
}

func (sa *SuMd5) value() Value {
	var buf [md5.Size]byte
	return SuStr(string(sa.hash.Sum(buf[0:0])))
}
