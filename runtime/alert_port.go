// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !windows || portable || com

package runtime

import (
	"fmt"
	"log"

	"github.com/apmckinlay/gsuneido/util/dbg"
	"github.com/apmckinlay/gsuneido/util/exit"
)

func Alert(args ...any) {
	fmt.Println(args...)
}

func AlertCancel(args ...any) bool {
	Alert(args...)
	return true
}

func Fatal(args ...any) {
	s := fmt.Sprintln(args...)
	log.Print("FATAL: ", s)
	dbg.PrintStack()
	exit.Exit(1)
}
