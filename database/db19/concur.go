// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"time"
)

type void struct{}

type concur struct {
	commitChan  chan *UpdateTran
	mergeChan   chan int
	persistChan chan void
	allDone     chan void
}

func StartConcur(persistInterval time.Duration) *concur {
	c := concur{
		commitChan:  make(chan *UpdateTran),
		mergeChan:   make(chan int),
		persistChan: make(chan void),
		allDone:     make(chan void),
	}
	go c.committer()
	go c.merger(persistInterval)
	go c.persister()
	return &c
}

func (c *concur) Stop() {
	close(c.commitChan)
	<-c.allDone
}

func (c *concur) committer() {
	for tran := range c.commitChan {
		c.mergeChan <- tran.Commit()
	}
	close(c.mergeChan)
}

func (c *concur) merger(persistInterval time.Duration) {
	ticker := time.NewTicker(persistInterval)
loop:
	for {
		select {
		case tn := <-c.mergeChan: // receive mergeMsg's from commit
			if tn == 0 { // zero value means channel closed
				break loop
			}
			Merge(tn)
		case <-ticker.C:
			// send ticks from here so we get back pressure
			fmt.Println("Persist")
			c.persistChan <- void{}
		}
	}
	c.persistChan <- void{}
	close(c.persistChan)
}

func (c *concur) persister() {
	for range c.persistChan {
		Persist()
	}
	close(c.allDone)
}
