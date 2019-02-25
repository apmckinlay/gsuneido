package runtime

type SuExcept struct {
	SuStr
	Callstack *SuObject
}

func NewSuExcept(t *Thread, s SuStr) *SuExcept {
	// capture call stack
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
	return &SuExcept{SuStr: s, Callstack: cs}
}

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
