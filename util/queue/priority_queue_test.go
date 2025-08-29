// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package queue

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestPriorityQueueBasic(t *testing.T) {
	pq := NewPriorityQueue()

	// Test single element
	pq.Put(1, 100, "first")
	value := pq.Get()
	assert.T(t).This(value).Is("first")
}

func TestPriorityQueuePriorityOrder(t *testing.T) {
	pq := NewPriorityQueue()

	// Add elements from different transactions with different priorities
	pq.Put(3, 1, "high")
	pq.Put(1, 2, "low")
	pq.Put(2, 3, "medium")

	// Should get highest priority first
	value := pq.Get()
	assert.T(t).This(value).Is("high")

	// Then medium priority
	value = pq.Get()
	assert.T(t).This(value).Is("medium")

	// Finally low priority
	value = pq.Get()
	assert.T(t).This(value).Is("low")
}

func TestPriorityQueueSameTransaction(t *testing.T) {
	pq := NewPriorityQueue()

	// Add multiple elements from same transaction
	pq.Put(1, 100, "first")
	pq.Put(3, 100, "third")
	pq.Put(2, 100, "second")

	// Should get FIFO order (oldest first) for same transaction
	value := pq.Get()
	assert.T(t).This(value).Is("first")

	value = pq.Get()
	assert.T(t).This(value).Is("third")

	value = pq.Get()
	assert.T(t).This(value).Is("second")
}

func TestPriorityQueueMixedTransactions(t *testing.T) {
	pq := NewPriorityQueue()

	// Transaction 1: priority 1, then 3
	pq.Put(1, 1, "t1-low")
	pq.Put(3, 1, "t1-high")

	// Transaction 2: priority 2
	pq.Put(2, 2, "t2-medium")

	// Should get highest priority among oldest per transaction
	// t1 oldest = priority 1, t2 oldest = priority 2
	// So t2 wins with priority 2
	value := pq.Get()
	assert.T(t).This(value).Is("t2-medium")

	// Now t1 oldest is still priority 1
	value = pq.Get()
	assert.T(t).This(value).Is("t1-low")

	// Finally the remaining t1 element
	value = pq.Get()
	assert.T(t).This(value).Is("t1-high")
}

func TestPriorityQueueBlocking(t *testing.T) {
	pq := NewPriorityQueue()

	// Test blocking when empty
	done := make(chan bool)
	go func() {
		time.Sleep(50 * time.Millisecond)
		pq.Put(1, 1, "delayed")
		done <- true
	}()

	start := time.Now()
	value := pq.Get()
	elapsed := time.Since(start)

	assert.T(t).This(value).Is("delayed")
	assert.T(t).This(elapsed > 40*time.Millisecond).Is(true)
	<-done

	// Test blocking when full
	for i := 0; i < bufSize; i++ {
		pq.Put(i, i, i)
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		pq.Get() // make space
		done <- true
	}()

	start = time.Now()
	pq.Put(99, 99, "blocked")
	elapsed = time.Since(start)

	assert.T(t).This(elapsed > 40*time.Millisecond).Is(true)
	<-done
}

func TestPriorityQueueConcurrent(t *testing.T) {
	const nproducers = 4
	const nitemsperproducer = 20
	const totalItems = nproducers * nitemsperproducer

	pq := NewPriorityQueue()
	var produced atomic.Int32
	var consumed atomic.Int32
	var wg sync.WaitGroup

	// Start single consumer
	wg.Go(func() {
		for consumed.Load() < totalItems {
			pq.Get()
			consumed.Add(1)
		}
	})

	// Start multiple producers
	for i := range nproducers {
		wg.Go(func() {
			for j := range nitemsperproducer {
				pq.Put(j%3, i, j)
				produced.Add(1)
			}
		})
	}

	wg.Wait()
	assert.T(t).This(int(produced.Load())).Is(totalItems)
	assert.T(t).This(int(consumed.Load())).Is(totalItems)
}

func TestPriorityQueueEdgeCases(t *testing.T) {
	pq := NewPriorityQueue()

	// Test with zero priority
	pq.Put(0, 1, "zero")
	value := pq.Get()
	assert.T(t).This(value).Is("zero")

	// Test with negative priority
	pq.Put(-1, 2, "negative")
	pq.Put(1, 3, "positive")

	// Positive should win
	value = pq.Get()
	assert.T(t).This(value).Is("positive")

	value = pq.Get()
	assert.T(t).This(value).Is("negative")

	// Test with nil value
	pq.Put(1, 1, nil)
	value = pq.Get()
	assert.T(t).This(value).Is(nil)
}

// Benchmarks comparing PriorityQueue vs simple Go channel

func BenchmarkPriorityQueue(b *testing.B) {
	pq := NewPriorityQueue()
	val := &struct{}{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pq.Put(i%3, i%5, val)
		pq.Get()
	}
}

func BenchmarkGoChannel(b *testing.B) {
	ch := make(chan int, bufSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch <- i
		<-ch
	}
}

func BenchmarkPriorityQueueConcurrent(b *testing.B) {
	pq := NewPriorityQueue()
	val := &struct{}{}

	// Pre-fill queue to avoid initial blocking
	for i := 0; i < bufSize/2; i++ {
		pq.Put(i, i, val)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Do both operations to maintain queue size
			pq.Put(1, 1, val)
			pq.Get()
		}
	})
}

func BenchmarkGoChannelConcurrent(b *testing.B) {
	ch := make(chan int, bufSize)

	// Pre-fill channel to avoid initial blocking
	for i := 0; i < bufSize/2; i++ {
		ch <- i
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Do both operations to maintain channel size
			ch <- 1
			<-ch
		}
	})
}

func TestPriorityQueue2(t *testing.T) {
	pq := NewPriorityQueue()
	val := &struct{}{}
	var wg sync.WaitGroup
	const nthreads = 32
	for range nthreads {
		wg.Add(1)
		go func() {
			for range 100_000 {
				pq.Put(1, 1, val)
			}
			wg.Done()
			fmt.Println("done")
		}()
	}
	wg.Add(1)
	go func() {
		for range 100_000 * nthreads {
			pq.Get()
		}
		wg.Done()
	}()
	wg.Wait()
}
