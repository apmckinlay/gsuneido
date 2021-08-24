// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/apmckinlay/gsuneido/db19/meta"
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
// Concurrency is separate so we can test functionality
// without any goroutines or channels.
//
// To stop we close the checker channel, and then each following stage
// closes its output channel.
// Finally the merger closes the allDone channel
// so we know the shutdown has finished.
func StartConcur(db *Database, persistInterval time.Duration) {
	mergeChan := make(chan interface{}, chanBuffers)
	resultChan := make(chan error)
	allDone := make(chan void)
	go merger(db, mergeChan, resultChan, persistInterval, allDone)
	db.ck = StartCheckCo(db, mergeChan, resultChan, allDone)
}

type mergeT struct {
	db         *Database
	mergeChan  chan interface{}
	merges     *mergeList
	em         *execMulti
	resultChan chan error
}

func merger(db *Database, mergeChan chan interface{}, resultChan chan error,
	persistInterval time.Duration, allDone chan void) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			log.Fatalln("FATAL ERROR in merger:", e)
		}
	}()
	em := startMergeWorkers()
	ep := startExecPersistMulti()
	// ep := &execPersistSingle{}
	merges := &mergeList{}
	mt := mergeT{db: db, mergeChan: mergeChan, merges: merges, em: em,
		resultChan: resultChan}
	ticker := time.NewTicker(persistInterval)
	prevState := db.GetState()
loop:
	for {
		select {
		case m := <-mergeChan: // receive mergeMsg's from commit
			if m == nil { // channel closed
				break loop
			}
			mt.dispatch(m)
		case <-ticker.C:
			if db.GetState() != prevState {
				prevState = db.Persist(ep, false)
			}
		}
	}
	close(em.jobChan)
	if db.GetState() != prevState {
		db.Persist(ep, false)
	}
	close(allDone)
}

func (mt *mergeT) dispatch(m interface{}) {
	for {
		switch m2 := m.(type) {
		case []string: // merge
			mt.merges.reset()
			mt.merges.add(m2)
			m = mt.merges.drain(mt.mergeChan)
			mt.db.Merge(mt.em.merge, mt.merges)
			// db.Merge(mergeSingle, merges)
			if m == nil {
				return
			}
		case func() error: // run
			mt.resultChan <- run(m2)
			return
		}
	}
}

func run(fn func() error) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	return fn()
}

// mergeSingle is a single threaded merge for tran_test
func mergeSingle(state *DbState, merges *mergeList) []meta.MergeUpdate {
	var results []meta.MergeUpdate
	for _, tn := range merges.tn {
		result := state.Meta.Merge(tn.table, tn.nmerge)
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
	// if only one table, just merge it in this thread
	// and avoid overhead of channels and worker
	if len(merges.tn) == 1 {
		m := merges.tn[0]
		result := state.Meta.Merge(m.table, m.nmerge)
		return append(merges.results, result)
	}
	for i := 0; i < len(merges.tn); {
		select {
		case em.jobChan <- job{meta: state.Meta,
			table: merges.tn[i].table, nmerge: merges.tn[i].nmerge}:
			i++
		case result := <-em.resultChan:
			merges.results = append(merges.results, result)
		}
	}
	for len(merges.results) < len(merges.tn) {
		result := <-em.resultChan
		merges.results = append(merges.results, result)
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
	tn      []tableCount
	results []meta.MergeUpdate
}

type tableCount struct {
	table  string
	nmerge int
}

func (ml *mergeList) add(m []string) {
outer:
	for _, table := range m {
		for i := range ml.tn {
			if ml.tn[i].table == table {
				ml.tn[i].nmerge++
				continue outer
			}
		}
		ml.tn = append(ml.tn, tableCount{table: table, nmerge: 1})
	}
}

func (ml *mergeList) drain(mergeChan chan interface{}) interface{} {
	for {
		select {
		case m := <-mergeChan:
			if m == nil { // channel closed
				return nil
			}
			if m2, ok := m.([]string); ok {
				ml.add(m2)
			} else {
				return m // create or drop that needs to be handled
			}
		default: // channel empty
			return nil
		}
	}
}

func (ml *mergeList) reset() {
	ml.tn = ml.tn[:0]
	ml.results = ml.results[:0]
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
	count      int
	results    []meta.PersistUpdate
	workChan   chan func() meta.PersistUpdate
	resultChan chan meta.PersistUpdate
}

const nPersistWorkers = 8 // ???

func startExecPersistMulti() *execPersistMulti {
	workChan := make(chan func() meta.PersistUpdate, 1)
	resultChan := make(chan meta.PersistUpdate, 1)
	for i := 0; i < nPersistWorkers; i++ {
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
