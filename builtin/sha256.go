// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"hash"
	"io"

	"crypto/sha256"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

type SuSha256 struct {
	CantConvert
	hash hash.Hash
}

var _ = builtinRaw("Sha256(@args)",
	func(th *Thread, as *ArgSpec, args []Value) Value {
		sa := &SuSha256{hash: sha256.New()}
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

var _ Value = (*SuSha256)(nil)

func (*SuSha256) Get(*Thread, Value) Value {
	panic("Sha256 does not support get")
}

func (*SuSha256) Put(*Thread, Value, Value) {
	panic("Sha256 does not support put")
}

func (*SuSha256) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic("Sha256 does not support update")
}

func (*SuSha256) RangeTo(int, int) Value {
	panic("Sha256 does not support range")
}

func (*SuSha256) RangeLen(int, int) Value {
	panic("Sha256 does not support range")
}

func (*SuSha256) Hash() uint32 {
	panic("Sha256 hash not implemented")
}

func (*SuSha256) Hash2() uint32 {
	panic("Sha256 hash not implemented")
}

func (*SuSha256) Compare(Value) int {
	panic("Sha256 compare not implemented")
}

func (*SuSha256) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call Sha256")
}

func (*SuSha256) String() string {
	return "aSha256"
}

func (*SuSha256) Type() types.Type {
	return types.BuiltinClass
}

func (sa *SuSha256) Equal(other interface{}) bool {
	sa2, ok := other.(*SuSha256)
	return ok && sa == sa2
}

func (*SuSha256) Lookup(_ *Thread, method string) Callable {
	return sha256Methods[method]
}

var sha256Methods = Methods{
	"Update": method1("(string)", func(this, arg Value) Value {
		io.WriteString(this.(*SuSha256).hash, ToStr(arg))
		return this
	}),
	"Value": method0(func(this Value) Value {
		return this.(*SuSha256).value()
	}),
}

func (sa *SuSha256) value() Value {
	var buf [sha256.Size]byte
	return SuStr(string(sa.hash.Sum(buf[0:0])))
}
