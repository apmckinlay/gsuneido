// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"fmt"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var sunappThread *Thread

// sunAPP is called by goside.go <- interact <- cside.c <- sunapp.cpp
func sunAPP(url string) (result string) {
	if sunappThread == nil {
		sunappThread = UIThread.SubThread()
	}
	for i := 0; i < 1; i++ { // workaround for 1.14 bug
		defer func() {
			if err := recover(); err != nil {
				result = fmt.Sprint("SuneidoApp("+url+")", err)
			}
		}()
	}
	f := Global.GetName(sunappThread, "SuneidoAPP")
	return ToStr(sunappThread.Call(f, SuStr(url)))
}
