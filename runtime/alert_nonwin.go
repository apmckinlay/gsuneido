// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !windows portable

package runtime

import (
	"fmt"
	"log"
	"os"
)

func Alert(args ...interface{}) {
	fmt.Println(args...)
}

func Fatal(args ...interface{}) {
	s := fmt.Sprintln(args...)
	log.Print("FATAL: ", s)
	os.Exit(1)
}
