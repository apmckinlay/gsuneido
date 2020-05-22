// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build interactive

package model

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

const nIndex = 5

type Index struct {
	lock   sync.Mutex
	worker chan Task
}

type Task struct {
	id    int
	idx   int
	index *Index
}

func (task Task) String() string {
	return fmt.Sprint("Task ", task.id, " for index ", task.idx)
}

var indexes [nIndex]Index

const maxPoolSize = 4

var pool = make(chan chan Task, maxPoolSize)

func TestIndexWorkers(*testing.T) {
	source := make(chan Task)
	go manager(source)
	for i := 0; i < 10; i++ {
		idx := rand.Intn(nIndex)
		source <- Task{id: i, idx: idx, index: &indexes[idx]}
		randSleep(5)
	}
	randSleep(100)
}

func manager(source chan Task) {
	wid := 0
	for task := range source {
		fmt.Println("manager received", task)
		index := &indexes[task.idx]
		index.lock.Lock()
		if index.worker != nil {
			index.worker <- task
		} else {
			select {
			case workerChan := <-pool:
				fmt.Println("got worker from pool")
				index.worker = workerChan
				index.worker <- task
			default:
				fmt.Println("create worker")
				index.worker = make(chan Task, 8)
				go worker(wid, task, index.worker)
				wid++
			}
		}
		index.lock.Unlock()
	}
}

func worker(wid int, task Task, tasks chan Task) {
	for {
		idx := task.idx
		fmt.Println("\tworker", wid, "executing", task)
		execute(task)
		index := task.index
		index.lock.Lock()
		select {
		case task = <-tasks: // non-blocking
			fmt.Println("\tworker", wid, "RECEIVED", task)
			index.lock.Unlock()
			if task.idx != idx {
				panic("wrong index")
			}
		default:
			select {
			case pool <- tasks: // return worker to pool
				fmt.Println("\tworker", wid, "no more tasks - returned to pool")
				index.worker = nil
				index.lock.Unlock()
				task = <-tasks // block until re-assigned
			default:
				fmt.Println("\tworker", wid, "no more tasks - terminating")
				index.lock.Unlock()
				return // pool full, so terminate
			}
		}
	}
}

func execute(Task) {
	randSleep(10)
}
