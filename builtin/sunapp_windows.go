package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var sunappThread *Thread

func sunAPP(url string) string {
	if sunappThread == nil {
		sunappThread = UIThread.SubThread()
	}
	f := Global.GetName(sunappThread, "SuneidoAPP")
	result := sunappThread.Call(f, SuStr(url))
	return ToStr(result)
}
