// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package exit

import (
	"log"
	"os"
	"sync"
	"time"
)

type exitfn struct {
	desc string
	fn   func()
}

var exitfns []exitfn

var hanger sync.Mutex

// Add registers a function to be called on exit.
func Add(desc string, fn func()) {
	exitfns = append(exitfns, exitfn{desc: desc, fn: fn})
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
	i := len(exitfns) - 1
	var t time.Time
	ds := make([]time.Duration, len(exitfns))
	// failsafe in case exit funcs don't return
	go func() {
		time.Sleep(10 * time.Second)
		for j := len(exitfns) - 1; j > i; j-- {
			log.Println("Exit:", exitfns[j].desc, ds[j])
		}
		if i >= 0 {
			log.Println("Exit:", exitfns[i].desc, "didn't finish")
		}
		log.Fatalln("FATAL exit timeout")
	}()

	hanger.Lock() // never unlocked
	for ; i >= 0; i-- {
		func() {
			defer func() {
				if e := recover(); e != nil {
					log.Println("ERROR during Exit" + exitfns[i].desc, ":", e)
				}
			}()
			t = time.Now()
			exitfns[i].fn()
			ds[i] = time.Since(t)
		}()
	}
}

// Wait should only be called after Exit or RunFuncs. It blocks until exit.
func Wait() {
	hanger.Lock() // should be locked
	log.Fatalln("FATAL exit.Wait: shouldn't reach here")
}
