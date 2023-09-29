// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"fmt"
	"log"

	"github.com/apmckinlay/gsuneido/util/exit"
)

func Fatal(args ...any) {
	s := fmt.Sprintln(args...)
	log.Print("FATAL: ", s)
	// if args[0] != "lost connection" && args[0] != "Can't connect." {
	// 	dbg.PrintStack()
	// }
	Fatal2(s)
	exit.Exit(1)
}
