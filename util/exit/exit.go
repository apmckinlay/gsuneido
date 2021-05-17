// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package exit

import (
	"log"
	"os"
)

var exitfns []func()

func Add(fn func()) {
	exitfns = append(exitfns, fn)
}

func Exit(code int) {
	for _, fn := range exitfns {
		func() {
			defer func() {
				if e := recover(); e != nil {
					log.Println("ERROR during Exit:", e)
				}
			}()
			fn()
		}()
	}
	os.Exit(code)
}
