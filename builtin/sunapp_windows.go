// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"fmt"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var sunappThread *Thread

func sunAPP(url string) (result string) {
	if sunappThread == nil {
		sunappThread = UIThread.SubThread()
	}
	defer func() {
		if err := recover(); err != nil {
			result = fmt.Sprint("SuneidoApp("+url+")", err)
		}
	}()
	f := Global.GetName(sunappThread, "SuneidoAPP")
	return ToStr(sunappThread.Call(f, SuStr(url)))
}
