// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"time"
)

type void struct{}

const chanBuffers = 4 // ???

// StartConcur starts the database pipeline -
// starting the goroutines and connecting them with channels.
//
// checker -> merger
//
// persist is done by merger every persistInterval
//
// Concurrency is separate so we can test functionality
// without any goroutines or channels.
//
// To stop we close the checker channel, and then each following stage
// closes its output channel.
// Finally the merger closes the allDone channel
// so we know the shutdown has finished.
func StartConcur(db *Database, persistInterval time.Duration) {
	mergeChan := make(chan int, chanBuffers)
	allDone := make(chan void)
	go merger(db, mergeChan, persistInterval, allDone)
	db.ck = StartCheckCo(mergeChan, allDone)
}

func merger(db *Database, mergeChan chan int,
	persistInterval time.Duration, allDone chan void) {
	ticker := time.NewTicker(persistInterval)
	prevState := db.GetState()
loop:
	for {
		select {
		case tn := <-mergeChan: // receive mergeMsg's from commit
			if tn == 0 { // zero value means channel closed
				break loop
			}
			db.Merge(tn)
		case <-ticker.C:
			// fmt.Println("Persist")
			state := db.GetState()
			if state != prevState {
				db.Persist(false)
				prevState = state
			}
		}
	}
	db.Persist(true) // flatten on shutdown (required by quick check)
	close(allDone)
}
