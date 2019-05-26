package builtin

import (
	"hash"
	"io"

	"crypto/sha1"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

var _ = builtinRaw("Sha1(@args)",
	func(th *Thread, as *ArgSpec, args []Value) Value {
		sa := &SuSha1{hash: sha1.New()}
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

type SuSha1 struct {
	CantConvert
	hash hash.Hash
}

var _ Value = (*SuSha1)(nil)

func (*SuSha1) Get(*Thread, Value) Value {
	panic("Sha1 does not support get")
}

func (*SuSha1) Put(*Thread, Value, Value) {
	panic("Sha1 does not support put")
}

func (*SuSha1) RangeTo(int, int) Value {
	panic("Sha1 does not support range")
}

func (*SuSha1) RangeLen(int, int) Value {
	panic("Sha1 does not support range")
}

func (*SuSha1) Hash() uint32 {
	panic("Sha1 hash not implemented")
}

func (*SuSha1) Hash2() uint32 {
	panic("Sha1 hash not implemented")
}

func (*SuSha1) Compare(Value) int {
	panic("Sha1 compare not implemented")
}

func (*SuSha1) Call(*Thread, *ArgSpec) Value {
	panic("can't call Sha1")
}

func (*SuSha1) String() string {
	return "aSha1"
}

func (*SuSha1) Type() types.Type {
	return types.Sha1
}

func (sa *SuSha1) Equal(other interface{}) bool {
	if sa2, ok := other.(*SuSha1); ok {
		return sa == sa2
	}
	return false
}

func (*SuSha1) Lookup(_ *Thread, method string) Callable {
	return sha1Methods[method]
}

var sha1Methods = Methods{
	"Update": method1("(string)", func(this, arg Value) Value {
		io.WriteString(this.(*SuSha1).hash, ToStr(arg))
		return this
	}),
	"Value": method0(func(this Value) Value {
		return this.(*SuSha1).value()
	}),
}

func (sa *SuSha1) value() Value {
	var buf [sha1.Size]byte
	return SuStr(string(sa.hash.Sum(buf[0:0])))
}
