// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable && !com

package builtin

import (
	"log"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(exit, "(code = 0)")

func exit(arg Value) Value {
	if arg == True {
		Exit(0) // immediate exit
	}
	go func() {
		time.Sleep(15 * time.Second) // longer than exit.Exit timeout
		log.Fatalln("FATAL: wingui Exit() timeout")
	}()
	postQuit(uintptr(IfInt(arg)))
	return nil
}
