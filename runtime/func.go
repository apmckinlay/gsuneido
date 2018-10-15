package runtime

import (
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/str"
)

// Func is the base for callables
type Func struct {
	ParamSpec
}

// Value interface (except TypeName)

func (f *Func) ToInt() int {
	panic("cannot convert function to integer")
}

func (f *Func) ToDnum() dnum.Dnum {
	panic("cannot convert function to number")
}

func (f *Func) ToStr() string {
	panic("cannot convert function to string")
}

func (f *Func) String() string {
	var buf strings.Builder
	buf.WriteString("function (")
	sep := ""
	v := 0 // index into Values
	for i := 0; i < f.Nparams; i++ {
		buf.WriteString(sep)
		buf.WriteString(flagsToName(f.Strings[i], f.Flags[i]))
		if i >= f.Nparams-f.Ndefaults {
			buf.WriteString("=")
			buf.WriteString(fmt.Sprint(f.Values[v]))
			v++
		}
		sep = ", "
	}
	buf.WriteString(")")
	return buf.String()
}

func flagsToName(p string, flags Flag) string {
	if flags == AtParam {
		p = "@" + p
	}
	if flags&PubParam == PubParam {
		p = str.Capitalize(p)
	}
	if flags&DynParam == DynParam {
		p = "_" + p
	}
	if flags&DotParam == DotParam {
		p = "." + p
	}
	return p
}

func (f *Func) Get(Value) Value {
	panic("function does not support get")
}

func (f *Func) Put(Value, Value) {
	panic("function does not support put")
}

func (Func) RangeTo(int, int) Value {
	panic("function does not support range")
}

func (Func) RangeLen(int, int) Value {
	panic("function does not support range")
}

func (f *Func) Hash() uint32 {
	panic("function hash not implemented")
}

func (f *Func) Hash2() uint32 {
	panic("function hash not implemented")
}

func (f *Func) Equal(other interface{}) bool {
	if f2, ok := other.(*Func); ok {
		return f == f2
	}
	return false
}

func (*Func) Order() Ord {
	return OrdOther
}

func (f *Func) Compare(Value) int {
	panic("function compare not implemented")
}
