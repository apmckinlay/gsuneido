// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"fmt"
	"log"
	"os"
)

// Exit is injected by gSuneido as either exit.Exit or stopService
var Exit func(code int)

func Fatal(args ...any) {
	s := fmt.Sprintln(args...)
	log.Print("FATAL: ", s)
	// if args[0] != "lost connection" && args[0] != "Can't connect." {
	// 	dbg.PrintStack()
	// }
	Fatal2(s)
	if Exit == nil { // e.g. during initialization
		os.Exit(1)
	}
	Exit(1)
}
