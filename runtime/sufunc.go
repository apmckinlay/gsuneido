package runtime

import (
	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/str"
)

// SuFunc is a compiled Suneido function, method, or block
type SuFunc struct {
	ParamSpec

	// Nlocals is the number of parameters and local variables
	Nlocals uint8

	// Code is the actual byte code
	Code []byte //TODO change to string

	// ArgSpecs used by calls in the code
	ArgSpecs []*ArgSpec

	// ClassName is used to privatize dot params
	ClassName string

	// Id is a unique identifier for a function defining blocks
	Id uint32
	// OuterId is the Id of the outer SuFunc
	// It is used by interp to handle block return
	OuterId uint32

	// SrcPos contains pairs of source and code position deltas
	SrcPos string
	// SrcBase is the starting point for the SrcPos source deltas
	SrcBase int
}

// Value interface (mostly handled by ParamSpec) --------------------

var _ Value = (*SuFunc)(nil)

func (f *SuFunc) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&f.ParamSpec, as)
	for i, flag := range f.Flags {
		if flag&DotParam == DotParam {
			name := f.Names[i]
			if flag&PubParam == PubParam {
				name = str.Capitalize(name)
			} else { // privatize
				name = f.ClassName + "_" + name
			}
			t.this.Put(SuStr(name), args[i])
		}
	}
	return t.Call(f)
}

func (f *SuFunc) Type() types.Type {
	if f.OuterId != 0 {
		return types.Block
	}
	return types.Function
}

// SuFuncMethods is initialized by the builtin package
var SuFuncMethods Methods

func (*SuFunc) Lookup(method string) Value {
	if m, ok := ParamsMethods[method]; ok {
		return m
	}
	return SuFuncMethods[method]
}

func (f *SuFunc) String() string {
	s := ""
	if f.Name != "" {
		s = f.Name + " "
	}
	s += "/* "
	if f.ClassName != "" {
		s += "method"
	} else {
		s += "function"
	}
	s += " */"
	return s
}

func (f *SuFunc) CodeToSrcPos(ip int) int {
	ip-- // because interp will have already incremented
	sp := f.SrcBase
	cp := 0
	for i := 0; i < len(f.SrcPos); i += 2 {
		prev := sp
		sp += int(f.SrcPos[i])
		cp += int(f.SrcPos[i+1])
		if cp > ip {
			return prev
		}
	}
	return sp // ???
}
