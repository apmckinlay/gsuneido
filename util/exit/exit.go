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

func Add(fn func()) {
	exitfns = append(exitfns, fn)
}

func Exit(code int) {
	// First call gets in, any later ones just block here until exit
	hanger.Lock() // never unlocked

	// failsafe in case this goroutine doesn't get to exit
	go func() {
		time.Sleep(5 * time.Second)
		log.Fatalln("exit failsafe")
	}()
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
	os.Exit(code)
}

func Wait() {
	hanger.Lock()
	log.Fatalln("exit.Wait: shouldn't reach here")
}
