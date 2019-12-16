// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/go-ole/go-ole"
)

func init() {
	connection := &ole.Connection{}
	err := connection.Initialize()
	if err != nil {
		panic("can't initialize COM/OLE")
	}
}

type suComObject struct {
	CantConvert
	iunk  *ole.IUnknown
	idisp *ole.IDispatch
}

var _ = builtin1("COMobject(progid)",
	func(arg Value) Value {
		if n, ok := arg.ToInt(); ok {
			iunk := (*ole.IUnknown)(unsafe.Pointer(uintptr(n)))
			idisp, err := iunk.QueryInterface(ole.IID_IDispatch)
			if err == nil && idisp != nil {
				return &suComObject{idisp: idisp}
			}
			return &suComObject{iunk: iunk}
		}
		panic("COMobject: only numeric implemented")
	})

var suComObjectMethods = Methods{
	"Dispatch?": method0(func(this Value) Value {
		return SuBool(this.(*suComObject).idisp != nil)
	}),
	"Release": method0(func(this Value) Value {
		this.(*suComObject).idisp.Release()
		return nil
	}),
}

var _ Value = (*suComObject)(nil)

func (sco *suComObject) Get(_ *Thread, mem Value) Value {
	name := ToStr(mem)
	v, err := sco.idisp.GetProperty(name)
	if err != nil {
		panic("COMobject get " + name + ": " + err.Error())
	}
	return variantToSu(v)
}

func variantToSu(v *ole.VARIANT) Value {
	val := v.Value()
	if val == nil {
		return nil
	}
	switch val := val.(type) {
	case bool:
		return SuBool(val)
	case string:
		return SuStr(val)
	case int16:
		return IntVal(int(val))
	case int32:
		return IntVal(int(val))
	case int:
		return IntVal(int(val))
	case int64:
		return IntVal(int(val))
	case uint16:
		return IntVal(int(val))
	case uint32:
		return IntVal(int(val))
	case uint:
		return IntVal(int(val))
	case uint64:
		return IntVal(int(val))
	}
	panic("COMobject: can't convert to Suneido value")
}

func (sco *suComObject) Put(_ *Thread, mem Value, val Value) {
	name := ToStr(mem)
	_, err := sco.idisp.PutProperty(name, suToGo(val))
	if err != nil {
		panic("COMobject put " + name + ": " + err.Error())
	}
}

func (*suComObject) RangeTo(int, int) Value {
	panic("COMobject does not support range")
}

func (*suComObject) RangeLen(int, int) Value {
	panic("COMobject does not support range")
}

func (*suComObject) Hash() uint32 {
	panic("COMobject hash not implemented")
}

func (*suComObject) Hash2() uint32 {
	panic("COMobject hash not implemented")
}

func (*suComObject) Equal(interface{}) bool {
	panic("COMobject equals not implemented")
}

func (*suComObject) Compare(Value) int {
	panic("COMobject compare not implemented")
}

func (*suComObject) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call a COMobject instance")
}

func (*suComObject) String() string {
	return "COMobject"
}

func (*suComObject) Type() types.Type {
	return types.BuiltinClass
}

func (sco *suComObject) Lookup(_ *Thread, method string) Callable {
	if f, ok := suComObjectMethods[method]; ok {
		return f
	}
	return &SuBuiltinMethodRaw{
		Fn: func(t *Thread, as *ArgSpec, this Value, args []Value) Value {
			return sco.call(method, as, args)
		}}
}

func (sco *suComObject) call(method string, as *ArgSpec, args []Value) Value {
	v, err := sco.idisp.CallMethod(method, comArgs(as, args)...)
	if err != nil {
		panic("COMobject call " + method + ": " + err.Error())
	}
	return variantToSu(v)
}

type any = interface{}

func comArgs(as *ArgSpec, args []Value) []any {
	ca := make([]any, 0, 8)
	ai := NewArgsIter(as, args)
	for {
		name, val := ai()
		if name != nil || val == nil {
			break
		}
		y := suToGo(val)
		ca = append(ca, y)
	}
	return ca
}

func suToGo(x Value) any {
	switch x := x.(type) {
	case SuBool:
		return x == True
	case SuStr:
		return string(x)
	}
	if n, ok := x.ToInt(); ok {
		return n
	}
	panic("COMobject: can't convert from Suneido value")
}
