// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package exit

import (
	"log"
	"os"
	"sync"
	"time"
)

var exitfns []func()

var hanger sync.Mutex

// Add registers a function to be called on exit.
func Add(fn func()) {
	exitfns = append(exitfns, fn)
}

// Exit calls RunFuncs and then os.Exit
// It also starts a failsafe timer which will exit in 10 seconds regardless.
func Exit(code int) {
	RunFuncs()
	os.Exit(code)
}

// RunFuncs runs the Add'ed exit functions.
// Only the first caller will run them, any other callers will block.
// The functions are run in the reverse order that they were Add'ed.
func RunFuncs() {
	// failsafe in case exit funcs don't return
	go func() {
		time.Sleep(10 * time.Second)
		log.Fatalln("FATAL exit timeout")
	}()

	hanger.Lock() // never unlocked

	for i := len(exitfns) - 1; i >= 0; i-- {
		func() {
			defer func() {
				if e := recover(); e != nil {
					log.Println("ERROR during Exit:", e)
				}
			}()
			exitfns[i]()
		}()
	}
}

// Wait should only be called after Exit or RunFuncs. It blocks until exit.
func Wait() {
	hanger.Lock() // should be locked
	log.Fatalln("FATAL exit.Wait: shouldn't reach here")
}
