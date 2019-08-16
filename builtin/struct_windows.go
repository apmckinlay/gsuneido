package builtin

import (
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
)

type SuStructGlobal struct {
	SuBuiltin
	size int
}

func init() {
	Global.Builtin("INITCOMMONCONTROLSEX",
		&SuStructGlobal{size: int(unsafe.Sizeof(INITCOMMONCONTROLSEX{}))})
	//TODO add other structs
}

func (*SuStructGlobal) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call struct")
}

var structSize = method0(func(this Value) Value {
	return IntVal(this.(*SuStructGlobal).size)
})

func (d *SuStructGlobal) Lookup(_ *Thread, method string) Callable {
	if method == "Size" {
		return structSize
	}
	return nil
}

func (d *SuStructGlobal) String() string {
	return "/* builtin struct */"
}
