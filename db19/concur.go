// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"sort"
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
	mergeChan := make(chan merge, chanBuffers)
	allDone := make(chan void)
	go merger(db, mergeChan, persistInterval, allDone)
	db.ck = StartCheckCo(mergeChan, allDone)
}

func merger(db *Database, mergeChan chan merge,
	persistInterval time.Duration, allDone chan void) {

	em := startMergeWorkers()
	merges := &mergeList{}
	ticker := time.NewTicker(persistInterval)
	prevState := db.GetState()
loop:
	for {
		select {
		case m := <-mergeChan: // receive mergeMsg's from commit
			if m.tn == 0 { // zero value means channel closed
				break loop
			}
			merges.reset()
			merges.add(m)
			merges.drain(mergeChan)
			db.Merge(em.merge, merges)
			db.GetState().meta.CheckTnMerged(m.tn)
		case <-ticker.C:
			// fmt.Println("Persist")
			state := db.GetState()
			if state != prevState {
				db.Persist(false)
				prevState = state
			}
		}
	}
	close(em.jobChan)
	db.Persist(true) // flatten on shutdown (required by quick check)
	close(allDone)
}

// mergeSingle is a single threaded merge for tran_test
func mergeSingle(state *DbState, merges *mergeList) []meta.MergeUpdate {
	var results []meta.MergeUpdate
	for i, table := range merges.tables {
		result := state.meta.Merge(table, merges.tns[i:i+1])
		results = append(results, result)
	}
	return results
}

const nMergeWorkers = 8 // ???

type execMulti struct {
	jobChan    chan job
	resultChan chan meta.MergeUpdate
}

func startMergeWorkers() *execMulti {
	em := &execMulti{
		jobChan:    make(chan job, 1),
		resultChan: make(chan meta.MergeUpdate, 1),
	}
	for i := 0; i < nMergeWorkers; i++ {
		go em.worker()
	}
	return em
}

func (em *execMulti) merge(state *DbState, merges *mergeList) []meta.MergeUpdate {
	//TODO if only one table, just merge it in this thread
	// and avoid overhead of channels and worker
	sort.Stable(merges)
	i := 0
	j := 0
	n := merges.Len()
	for ; j < n && merges.tables[j] == merges.tables[i]; j++ {
	}
	sent := 0
	for i < n {
		select {
		case em.jobChan <- job{meta: state.meta,
			table: merges.tables[i], tns: merges.tns[i:j]}:
			sent++
			for i = j; j < n && merges.tables[j] == merges.tables[i]; j++ {
			}
		case result := <-em.resultChan:
			merges.results = append(merges.results, result)
		}
	}
	for len(merges.results) < sent {
		result := <-em.resultChan
		merges.results = append(merges.results, result)
	}
	return merges.results
}

type job struct {
	meta  *meta.Meta
	table string
	tns   []int
}

func (em *execMulti) worker() {
	for j := range em.jobChan {
		em.resultChan <- j.meta.Merge(j.table, j.tns)
	}
}

// mergeList uses parallel arrays for tables and tns to allow slices of tns
type mergeList struct {
	tables  []string
	tns     []int //TODO just need count
	results []meta.MergeUpdate
}

func (ml *mergeList) add(m merge) {
	for _, table := range m.tables {
		ml.tns = append(ml.tns, m.tn)
		ml.tables = append(ml.tables, table)
	}
}

func (ml *mergeList) drain(mergeChan chan merge) {
	for {
		select {
		case m := <-mergeChan:
			if m.tn == 0 { // zero value means channel closed
				return
			}
			ml.add(m)
		default:
			return
		}
	}
}

func (ml *mergeList) reset() {
	ml.tables = ml.tables[:0]
	ml.tns = ml.tns[:0]
	ml.results = ml.results[:0]
}

// sort interface

func (ml *mergeList) Len() int {
	return len(ml.tns)
}

func (ml *mergeList) Less(i, j int) bool {
	return ml.tables[i] < ml.tables[j]
}

func (ml *mergeList) Swap(i, j int) {
	ml.tables[i], ml.tables[j] = ml.tables[j], ml.tables[i]
	ml.tns[i], ml.tns[j] = ml.tns[j], ml.tns[i]
}
