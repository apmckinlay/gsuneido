// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"time"
)

type void struct{}

const chanBuffers = 2 // ???

// StartConcur starts the database pipeline -
// starting the goroutines and connecting them with channels.
//
// checker -> merger -> persister
//
// persister is triggered by merger every persistInterval
//
// Concurrency is separate so we can test functionality
// without any goroutines or channels.
//
// To stop we close the checker channel, and then each following stage
// closes its output channel.
// Finally the persister closes the allDone channel
// so we know the shutdown has finished.
func StartConcur(persistInterval time.Duration) *CheckCo {
	mergeChan := make(chan int, chanBuffers)
	persistChan := make(chan void, chanBuffers)
	allDone := make(chan void)
	go merger(mergeChan, persistChan, persistInterval)
	go persister(persistChan, allDone)
	return StartCheckCo(mergeChan, allDone)
}

func merger(mergeChan chan int, persistChan chan void,
	persistInterval time.Duration) {
	ticker := time.NewTicker(persistInterval)
loop:
	for {
		select {
		case tn := <-mergeChan: // receive mergeMsg's from commit
			if tn == 0 { // zero value means channel closed
				break loop
			}
			Merge(tn)
		case <-ticker.C:
			// send ticks from here so we get back pressure
			// fmt.Println("Persist")
			persistChan <- void{}
		}
	}
	persistChan <- void{}
	close(persistChan)
}

func persister(persistChan chan void, allDone chan void) {
	for range persistChan {
		Persist()
	}
	close(allDone)
}
