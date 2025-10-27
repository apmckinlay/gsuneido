// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tests

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestTimestamp(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	db19.StartTimestamps()
	var lock sync.Mutex
	var prev Value = SuDate{}
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			var th Thread
			th.SetDbms(&dbms.DbmsLocal{})
			var ts Value
			for range 100 {
				n := rand.Intn(50)
				for range n {
					lock.Lock()
					ts = th.Timestamp()
					assert.That(ts.Compare(prev) > 0)
					assert.That(prev.Compare(ts) < 0)
					assert.False(ts.Equal(prev))
					assert.False(prev.Equal(ts))
					prev = ts
					lock.Unlock()
				}
				fmt.Println(ts)
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)))
			}
		})
	}
	wg.Wait()
}
