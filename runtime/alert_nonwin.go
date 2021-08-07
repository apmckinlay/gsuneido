// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !windows portable com

package runtime

import (
	"fmt"
	"log"
	"github.com/apmckinlay/gsuneido/util/exit"
)

func Alert(args ...interface{}) {
	fmt.Println(args...)
}

func Fatal(args ...interface{}) {
	s := fmt.Sprintln(args...)
	log.Print("FATAL: ", s)
	exit.Exit(1)
}
