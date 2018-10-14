package builtin

import (
	. "github.com/apmckinlay/gsuneido/base"
	. "github.com/apmckinlay/gsuneido/interp"
	"github.com/apmckinlay/gsuneido/interp/global"
)

var _ = global.Add("Object", Builtin(suObject))

func suObject(as *ArgSpec, args ...Value) Value {
	if as.Unnamed >= EACH {
		ob := args[0].(*SuObject)
		return ob.Slice(int(as.Unnamed - EACH))
	}
	ob := SuObject{}
	i := 0
	for ; i < int(as.Unnamed); i++ {
		ob.Add(args[i])
	}
	for _,n := range as.Spec {
		ob.Put(SuStr(as.Names[n]), args[i])
		i++
	}
	return &ob
}
