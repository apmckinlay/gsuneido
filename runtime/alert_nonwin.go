// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !windows || portable || com
// +build !windows portable com

package runtime

import (
	"fmt"
	"github.com/apmckinlay/gsuneido/util/exit"
	"log"
)

func Alert(args ...interface{}) {
	fmt.Println(args...)
}

func Fatal(args ...interface{}) {
	s := fmt.Sprintln(args...)
	log.Print("FATAL: ", s)
	exit.Exit(1)
}
