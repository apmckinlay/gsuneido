package runtime

type SuExcept struct {
	SuStr
	Callstack *SuObject
}

func NewSuExcept(t *Thread, s SuStr) *SuExcept {
	return &SuExcept{SuStr: s, Callstack: CallStack(t)}
}

// SuValue interface ------------------------------------------------

func (*SuExcept) TypeName() string {
	return "Except"
}

// SuExceptMethods is initialized by the builtin package
var SuExceptMethods Methods

func (*SuExcept) Lookup(method string) Value {
	if m := SuExceptMethods[method]; m != nil {
		return m
	}
	return StringMethods[method]
}

// callstack --------------------------------------------------------

// NOTE: it might be more efficient
// to capture the call stack in an internal format (not an SuObject)
// and only build the SuObject if required

// callstack captures the call stack
func CallStack(t *Thread) *SuObject {
	cs := &SuObject{}
	for i := t.fp - 1; i >= 0; i-- {
		fr := t.frames[i]
		call := &SuObject{}
		call.Put(SuStr("fn"), fr.fn)
		locals := &SuObject{}
		for i, v := range fr.locals {
			if v != nil {
				locals.Put(SuStr(fr.fn.Names[i]), v)
			}
		}
		call.Put(SuStr("locals"), locals)
		cs.Add(call)
	}
	return cs
}
