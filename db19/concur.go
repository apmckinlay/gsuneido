// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"log"
	"time"

	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/util/dbg"
	"github.com/apmckinlay/gsuneido/util/exit"
)

type void = struct{}

const chanBuffers = 4 // ???

// StartConcur starts the database pipeline -
// starting the goroutines and connecting them with channels.
//
// checker -> merger
//
// persist is called by merger every persistInterval
//
// To stop we close the checker channel, and then each following stage
// closes its output channel.
// Finally the merger closes the allDone channel
// so we know the shutdown has finished.
func StartConcur(db *Database, persistInterval time.Duration) {
	mergeChan := make(chan todo, chanBuffers)
	allDone := make(chan void)
	go merger(db, db.GetState(), mergeChan, persistInterval, allDone)
	db.ck = StartCheckCo(db, mergeChan, allDone)
}

func merger(db *Database, prevState *DbState, mergeChan chan todo,
	persistInterval time.Duration, allDone chan void) {
	defer func() {
		if e := recover(); e != nil {
			dbg.PrintStack()
			log.Fatalln("FATAL: in merger:", e)
		}
	}()
	em := startMergeWorkers()
	ep := startExecPersistMulti()
	// ep := &execPersistSingle{}
	merges := &mergeList{}
	ticker := time.NewTicker(persistInterval)
loop:
	for {
		select {
		case td, ok := <-mergeChan: // receive todo's from checkco
			if !ok { // channel closed
				break loop
			}
			for {
				if td.ret != nil {
					if td.fn != nil {
						td.ret <- td.run()
					} else {
						// persist
						if db.GetState() != prevState {
							prevState = db.persist(ep, true)
						}
						td.ret <- prevState
					}
					break
				}
				merges.start(td)
				td, ok = merges.drain(mergeChan)
				db.Merge(em.merge, merges)
				// db.Merge(mergeSingle, merges)
				if !ok {
					break
				}
			}
		case <-ticker.C:
			if db.GetState() != prevState {
				prevState = db.persist(ep, true)
			}
		}
	}
	close(em.jobChan)
	if db.GetState() != prevState ||
		prevState.Off != db.Store.Size()-uint64(stateLen) {
		exit.Progress("persist starting")
		db.persist(ep, true)
		exit.Progress("persist finished")
	}
	close(allDone)
}

type todo struct {
	meta   *meta.Meta
	fn     func()
	ret    chan any
	tables []string
}

func (td *todo) run() (err any) {
	defer func() {
		if e := recover(); e != nil {
			err = e
		}
	}()
	td.fn()
	return nil
}

// mergeSingle is a single threaded merge for tran_test
func mergeSingle(m *meta.Meta, merges *mergeList) []meta.MergeUpdate {
	var results []meta.MergeUpdate
	for _, tn := range merges.tn {
		result := m.Merge(tn.table, tn.nmerge)
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
	for range nMergeWorkers {
		go em.worker()
	}
	return em
}

func (em *execMulti) merge(met *meta.Meta, merges *mergeList) []meta.MergeUpdate {
	// if only one table, just merge it in this thread
	// and avoid overhead of channels and worker
	if len(merges.tn) == 1 {
		mrg := merges.tn[0]
		result := met.Merge(mrg.table, mrg.nmerge)
		return append(merges.results, result)
	}
	nresults := 0
	for i := 0; i < len(merges.tn); {
		select {
		case em.jobChan <- job{meta: met,
			table: merges.tn[i].table, nmerge: merges.tn[i].nmerge}:
			i++
		case result := <-em.resultChan:
			merges.results = append(merges.results, result)
			nresults++
		}
	}
	for nresults < len(merges.tn) {
		result := <-em.resultChan
		merges.results = append(merges.results, result)
		nresults++
	}

	return merges.results
}

type job struct {
	meta   *meta.Meta
	table  string
	nmerge int
}

func (em *execMulti) worker() {
	for j := range em.jobChan {
		em.resultChan <- j.meta.Merge(j.table, j.nmerge)
	}
}

type mergeList struct {
	meta    *meta.Meta
	tn      []tableCount
	results []meta.MergeUpdate
}

type tableCount struct {
	table  string
	nmerge int
}

func (ml *mergeList) start(m todo) {
	ml.meta = m.meta
	ml.tn = ml.tn[:0]
	ml.results = ml.results[:0]
	ml.add(m.tables)
}

func (ml *mergeList) add(tables []string) {
outer:
	for _, table := range tables {
		for i := range ml.tn {
			if ml.tn[i].table == table {
				ml.tn[i].nmerge++
				continue outer
			}
		}
		ml.tn = append(ml.tn, tableCount{table: table, nmerge: 1})
	}
}

// drain returns the next message that can't be added to the mergeList
// and must be processed separately
func (ml *mergeList) drain(mergeChan chan todo) (_ todo, _ bool) {
	for {
		select {
		case td, ok := <-mergeChan:
			if !ok { // channel closed
				return
			}
			if td.ret == nil && ml.meta.SameSchemaAs(td.meta) {
				ml.add(td.tables)
			} else {
				return td, true // not added to merge (sync or persist)
			}
		default: // channel empty
			return
		}
	}
}

// ------------------------------------------------------------------

type execPersist interface {
	Submit(fn func() meta.PersistUpdate)
	Results() []meta.PersistUpdate
}

type execPersistSingle struct {
	results []meta.PersistUpdate
}

func (ep *execPersistSingle) Submit(fn func() meta.PersistUpdate) {
	result := fn()
	ep.results = append(ep.results, result)

}

func (ep *execPersistSingle) Results() []meta.PersistUpdate {
	results := ep.results
	ep.results = ep.results[:0]
	return results
}

//-------------------------------------------------------------------

type execPersistMulti struct {
	workChan   chan func() meta.PersistUpdate
	resultChan chan meta.PersistUpdate
	results    []meta.PersistUpdate
	count      int
}

const nPersistWorkers = 8 // ???

func startExecPersistMulti() *execPersistMulti {
	workChan := make(chan func() meta.PersistUpdate, 1)
	resultChan := make(chan meta.PersistUpdate, 1)
	for range nPersistWorkers {
		go persistWorker(workChan, resultChan)
	}
	return &execPersistMulti{workChan: workChan, resultChan: resultChan}
}

func (ep *execPersistMulti) Submit(fn func() meta.PersistUpdate) {
	for {
		select {
		case ep.workChan <- fn:
			ep.count++
			return
		case result := <-ep.resultChan:
			ep.count--
			ep.results = append(ep.results, result)
		}
	}
}

func (ep *execPersistMulti) Results() []meta.PersistUpdate {
	for ; ep.count > 0; ep.count-- {
		result := <-ep.resultChan
		ep.results = append(ep.results, result)
	}
	results := ep.results
	ep.results = ep.results[:0]
	return results
}

func persistWorker(
	workChan chan func() meta.PersistUpdate,
	resultChan chan meta.PersistUpdate) {
	for fn := range workChan {
		resultChan <- fn()
	}
}
