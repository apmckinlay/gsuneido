package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

// methods common to SuClass and SuInstance

func init() {
	BaseMethods = Methods{
		"Base": method0(func(this Value) Value {
			return find(this, func(v Value) Value { return v })
		}),
		"Base?": method1("(class)", func(this, c Value) Value {
			return find(this, func(v Value) Value {
				if v == c {
					return True
				}
				return nil
			})
		}),
		//TODO Eval, Eval2
		//TODO GetDefault
		"Member?": method1("(string)", func(this, arg Value) Value {
			m := IfStr(arg)
			result := this.(Findable).Finder(func(v Value, mb *MemBase) Value {
				if _, ok := mb.Data[m]; ok {
					return True
				}
				return nil
			})
			if result == nil {
				result = False
			}
			return result
		}),
		"Members": method0(func(this Value) Value {
			//TODO all:
			return this.(Findable).Finder(func(_ Value, mb *MemBase) Value {
				return mb.Members()
			})
		}),
		//TODO Method?
		//TODO MethodClass
		"Size": method0(func(this Value) Value {
			return this.(Findable).Finder(func(_ Value, mb *MemBase) Value {
				return IntVal(mb.Size())
			})
		}),
	}
}

func find(x Value, fn func(Value) Value) Value {
	first := true
	result := x.(Findable).Finder(func(v Value, _ *MemBase) Value {
		if first {
			first = false
			return nil
		}
		return fn(v)
	})
	if result == nil {
		result = False
	}
	return result
}
