package builtin

import (
	"time"

	. "github.com/apmckinlay/gsuneido/runtime"
)

type SuThreadGlobal struct {
	SuBuiltin1
}

func init() {
	name, ps := paramSplit("Thread(block)")
	Global.Builtin(name, &SuThreadGlobal{
		SuBuiltin1{threadCallClass, BuiltinParams{ParamSpec: *ps}}})
}

func threadCallClass(arg Value) Value {
	t2 := NewThread()
	go func() {
		defer t2.Close()
		t2.CallWithArgs(arg)
	}()
	return nil
}

var threadMethods = Methods{
	"Name": method("(name=false)", func(t *Thread, _ Value, args ...Value) Value {
		if args[0] != False {
			t.Name = ToStr(args[0])
		}
		return SuStr(t.Name)
	}),
}

func (d *SuThreadGlobal) Lookup(t *Thread, method string) Callable {
	if f, ok := threadMethods[method]; ok {
		return f
	}
	return d.Lookup(t, method) // for Params
}

func (d *SuThreadGlobal) String() string {
	return "Thread /* builtin class */"
}

var _ = builtin2("Scheduled(ms, block)",
	func(arg, block Value) Value {
		ms := time.Duration(ToInt(arg)) * time.Millisecond
		t2 := NewThread()
		go func() {
			defer t2.Close()
			time.Sleep(ms)
			t2.CallWithArgs(block)
		}()
		return nil
	})
