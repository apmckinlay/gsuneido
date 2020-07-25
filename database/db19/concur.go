// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"time"
)

type void struct{}

func StartConcur(persistInterval time.Duration) *CheckCo {
	commitChan := make(chan *UpdateTran)
	mergeChan := make(chan int)
	persistChan := make(chan void)
	allDone := make(chan void)
	go committer(commitChan, mergeChan)
	go merger(mergeChan, persistChan, persistInterval)
	go persister(persistChan, allDone)
	return StartCheckCo(commitChan, allDone)
}

func committer(commitChan chan *UpdateTran, mergeChan chan int) {
	for tran := range commitChan {
		tn := tran.commit()
		mergeChan <- tn
	}
	close(mergeChan)
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
