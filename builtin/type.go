package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin1("Type(value)",
	func(arg Value) Value {
		return SuStr(arg.TypeName())
	})

var _ = builtin1("Boolean?(value)",
	func(arg Value) Value {
		return SuBool(arg.TypeName() == "Boolean")
	})

var _ = builtin1("Number?(value)",
	func(arg Value) Value {
		return SuBool(arg.TypeName() == "Number")
	})

var _ = builtin1("String?(value)",
	func(arg Value) Value {
		return SuBool(arg.TypeName() == "String")
	})

var _ = builtin1("Date?(value)",
	func(arg Value) Value {
		return SuBool(arg.TypeName() == "Date")
	})

var _ = builtin1("Object?(value)",
	func(arg Value) Value {
		return SuBool(arg.TypeName() == "Object")
	})

var _ = builtin1("Record?(value)",
	func(arg Value) Value {
		return SuBool(arg.TypeName() == "Record")
	})

var _ = builtin1("Class?(value)",
	func(arg Value) Value {
		return SuBool(arg.TypeName() == "Class")
	})

var _ = builtin1("Instance?(value)",
	func(arg Value) Value {
		return SuBool(arg.TypeName() == "Instance")
	})

var _ = builtin1("Function?(value)",
	func(arg Value) Value {
		switch arg.TypeName() {
		case "Function", "Method", "BuiltinFunction":
			return True
		}
		return False
	})
