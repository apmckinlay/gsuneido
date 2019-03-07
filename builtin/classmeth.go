package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

// memberq is also used by InstanceMethods
var memberq = method1("(string)", func(this, arg Value) Value {
	return MemberQ(this.(Findable), arg)
})

func init() {
	ClassMethods = Methods{
		"*new*": methodRaw("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				return this.(*SuClass).New(t, as)
			}),
		"Base": method0(func(this Value) Value {
			b := this.(*SuClass).Parent()
			if b == nil {
				return False
			}
			return b
		}),
		"Members": method0(func(this Value) Value {
			return this.(*SuClass).Members()
		}),
		"Member?": memberq,
		"Size": method0(func(this Value) Value {
			return this.(*SuClass).Size()
		}),
		// TODO more methods
	}
}
