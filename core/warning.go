// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"fmt"
	"log"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/options"
)

var warningCount atomic.Int64

const maxWarnings = 100

func Warning(args ...any) {
	s := fmt.Sprintln(args...)
	if options.WarningsThrow.Load().Matches(s) {
		panic(s)
	}
	n := warningCount.Add(1)
	if n < maxWarnings {
		log.Print("WARNING: ", s)
	}
	if n == maxWarnings {
		log.Println("WARNING: too many warnings - stopping logging")
	}
}
