// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"sync"
	"time"
)

var tsOnce sync.Once
var tsLock sync.Mutex
var ts time.Time

func Timestamp() time.Time {
	tsOnce.Do(func() {
		ts = time.Now()
		go ticker()
	})
	tsLock.Lock()
	defer tsLock.Unlock()
	ts = ts.Add(1 * time.Millisecond)
	return ts
}

func ticker() {
	for {
		time.Sleep(1 * time.Second)
		t := time.Now()
		tsLock.Lock()
		if t.After(ts) {
			ts = t
		}
		tsLock.Unlock()
	}
}
