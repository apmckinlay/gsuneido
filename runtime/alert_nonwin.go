// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !win32

package runtime

import (
	"fmt"
	"log"
	"os"
)

func Alert(args ...interface{}) {
	s := fmt.Sprintln(args...)
	log.Print("Alert: ", s)
}

func Fatal(args ...interface{}) {
	s := fmt.Sprintln(args...)
	log.Print("FATAL: ", s)
	os.Exit(1)
}
