package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var memberq = method1("(string)", func(this, arg Value) Value {
	return MemberQ(this.(Findable), arg)
})

func init() {
	ClassMethods = Methods{
		"*new*": rawmethod("(@args)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				return this.(*SuClass).New(t, as)
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
