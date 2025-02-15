// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"fmt"

	. "github.com/apmckinlay/gsuneido/core"
)

// sunAPP is called by goside.go <- cside.c <- sunapp.cpp
// cside.c ensures it is always called on the main UI thread.
func sunAPP(url string) (result string) {
	state := MainThread.GetState()
	defer func() {
		if err := recover(); err != nil {
			result = fmt.Sprint("SuneidoApp("+url+")", err)
		}
		MainThread.RestoreState(state)
	}()
	f := Global.GetName(MainThread, "SuneidoAPP")
	return ToStr(MainThread.Call(f, SuStr(url)))
}
