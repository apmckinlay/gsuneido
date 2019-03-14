package runtime

type SuExcept struct {
	SuStr
	Callstack *SuObject
}

func NewSuExcept(t *Thread, s SuStr) *SuExcept {
	return &SuExcept{SuStr: s, Callstack: t.CallStack()}
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
