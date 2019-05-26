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
		sa := &SuAdler32{hash: adler32.New()}
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

type SuAdler32 struct {
	CantConvert
	hash hash.Hash32
}

var _ Value = (*SuAdler32)(nil)

func (*SuAdler32) Get(*Thread, Value) Value {
	panic("Adler32 does not support get")
}

func (*SuAdler32) Put(*Thread, Value, Value) {
	panic("Adler32 does not support put")
}

func (*SuAdler32) RangeTo(int, int) Value {
	panic("Adler32 does not support range")
}

func (*SuAdler32) RangeLen(int, int) Value {
	panic("Adler32 does not support range")
}

func (*SuAdler32) Hash() uint32 {
	panic("Adler32 hash not implemented")
}

func (*SuAdler32) Hash2() uint32 {
	panic("Adler32 hash not implemented")
}

func (*SuAdler32) Compare(Value) int {
	panic("Adler32 compare not implemented")
}

func (*SuAdler32) Call(*Thread, *ArgSpec) Value {
	panic("can't call Adler32")
}

func (*SuAdler32) String() string {
	return "aAdler32"
}

func (*SuAdler32) Type() types.Type {
	return types.Adler32
}

func (sa *SuAdler32) Equal(other interface{}) bool {
	if sa2, ok := other.(*SuAdler32); ok {
		return sa == sa2
	}
	return false
}

func (*SuAdler32) Lookup(_ *Thread, method string) Callable {
	return adler32Methods[method]
}

var adler32Methods = Methods{
	"Update": method1("(string)", func(this, arg Value) Value {
		io.WriteString(this.(*SuAdler32).hash, ToStr(arg))
		return this
	}),
	"Value": method0(func(this Value) Value {
		return IntVal(int(int32(this.(*SuAdler32).hash.Sum32())))
	}),
}
