// Governed by the MIT license found in the LICENSE file.

package nrc

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/exit"
	"github.com/apmckinlay/gsuneido/util/generic/atomics"
	"github.com/apmckinlay/gsuneido/util/generic/hamt"
)

// ever nevict batches, we evict batches older than ibatch - nkeep
const (
	nkeep  = 5000 // ???
	nevict = 512  // ???
)

var ibatch atomic.Uint32

type Hash [20]byte // sha1

func (h Hash) String() string {
	return fmt.Sprintf("%x", h[:4]) // for debugging, show prefix
}

// rmap is the read-only cache read by transactions
var rmap atomics.Value[hamt.Hamt[Hash, *item]]

type nrwrite struct {
	h Hash
	n int
}

type Batch struct {
	ws []nrwrite
}

// Add adds a write to the batch
func (c *Batch) Add(h Hash, n int) {
	c.ws = append(c.ws, nrwrite{h: h, n: n})
}

var hits, misses, adds, dups int

// Get does a lookup in rmap.
// If found, it re-adds it to the batch to track that it was recently used.
func (w *Batch) Get(h Hash) (int, bool) {
	if it, ok := rmap.Load().Get(h); ok {
		hits++
		atomic.StoreUint32(&it.batch, ibatch.Load())
		return int(it.n), true
	}
	misses++
	return 0, false
}

func init() {
	exit.Add("nrcache", func() {
		fmt.Println("nrcache", "adds", adds, "dups:", dups, "batches", ibatch.Load(), "hits:", hits, "misses:", misses)
	})
}

type Intfc interface {
	Add(h Hash, n int)
	Get(h Hash) (int, bool)
}

var _ Intfc = (*Batch)(nil)

var once sync.Once

// Save sends the batch to the writer goroutine
func (b Batch) Save() {
	once.Do(func() {
		c = make(chan Batch, 8) // ???
		go writer()
	})
	if len(b.ws) > 0 {
		c <- b
	}
}

// c is a channel of batches connecting Save to writer
var c chan Batch

// writer reads batches from c and creates a new version of rmap.
// Periodically, rmap is updated.
func writer() {
	mu := rmap.Load().Mutable()
	tick := time.Tick(5 * time.Second) // ???
	for {
		select {
		case batch := <-c:
			b := ibatch.Load()
			for _, w := range batch.ws {
				adds++
				if _, ok := mu.Get(w.h); ok {
					dups++
				}
				mu.Put(&item{key: w.h, n: uint32(w.n), batch: b})
			}
			if b > nkeep && b%nevict == 0 {
				evict(&mu)
				r := mu.Freeze()
				rmap.Store(r)
				mu = r.Mutable()
			}
			ibatch.Add(1)
		case <-tick:
			r := mu.Freeze()
			rmap.Store(r)
			mu = r.Mutable()
	}
	}
}

var size int

func evict(mu *hamt.Hamt[Hash, *item]) {
	b := ibatch.Load()
	if b <= nkeep || b%nevict != 0 {
		return
	}
	size = 0
	cutoff := b - nkeep
	ev := 0
	rmap.Load().ForEach(func(it *item) {
		if atomic.LoadUint32(&it.batch) < cutoff {
			mu.Delete(it.key)
			ev++
		}
		size++
	})
	fmt.Println("---", "batch", ibatch.Load(), "cutoff", cutoff, "evicted", ev, "size", size)
	fmt.Println("---", "adds", adds, "dups:", dups, "hits:", hits, "misses:", misses)
}

//-------------------------------------------------------------------

type item struct {
	n     uint32
	batch uint32 // should be accessed atomically
	key   Hash
}

func (it *item) Key() Hash {
	return it.key
}

func (it *item) Hash(k Hash) uint64 {
	return *(*uint64)(unsafe.Pointer(&k))
}

func (*item) Cksum() uint32 {
	return 0
}

func (*item) StorSize() int {
	return 0
}

func (*item) IsTomb() bool {
	return false
}

func (*item) LastMod() int {
	return 0
}

func (*item) SetLastMod(mod int) {
}

func (*item) Write(w *stor.Writer) {
}
