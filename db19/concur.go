// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"time"

	"github.com/apmckinlay/gsuneido/db19/meta"
)

type merge struct {
	tn     int
	tables []string
}

type void = struct{}

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
	// meta.Executor = execSingle // inject
	meta.Executor = startMergeWorkers().exec // inject
	mergeChan := make(chan merge, chanBuffers)
	allDone := make(chan void)
	go merger(db, mergeChan, persistInterval, allDone)
	db.ck = StartCheckCo(mergeChan, allDone)
}

func merger(db *Database, mergeChan chan merge,
	persistInterval time.Duration, allDone chan void) {

	ticker := time.NewTicker(persistInterval)
	prevState := db.GetState()
loop:
	for {
		select {
		case m := <-mergeChan: // receive mergeMsg's from commit
			if m.tn == 0 { // zero value means channel closed
				break loop
			}
			db.Merge(m.tn, m.tables)
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

func execSingle(tn int, tables []string,
	fn func(tn int, table string) meta.Update) []meta.Update {
	results := make([]meta.Update, len(tables))
	for i, table := range tables {
		results[i] = fn(tn, table)
	}
	return results
}

const nMergeWorkers = 8 // ???

type execMulti struct {
	work    chan job
	results chan meta.Update
}

func startMergeWorkers() *execMulti {
	work := make(chan job, 1)
	results := make(chan meta.Update, 1)
	for i := 0; i < nMergeWorkers; i++ {
		go mergeWorker(work, results)
	}
	return &execMulti{work: work, results: results}
}

func (em *execMulti) exec(tn int, tables []string,
	fn func(tn int, table string) meta.Update) []meta.Update {
	results := make([]meta.Update, len(tables))
	i := 0
	j := 0
	for i < len(tables) {
		select {
		case em.work <- job{fn: fn, tn: tn, table: tables[i]}:
			i++
		case results[j] = <-em.results:
			j++
		}
	}
	for j < len(tables) {
		results[j] = <-em.results
		j++
	}
	return results
}

type job struct {
	fn    func(tn int, table string) meta.Update
	tn    int
	table string
}

func mergeWorker(work chan job, results chan meta.Update) {
	for w := range work {
		results <- w.fn(w.tn, w.table)
	}
}
