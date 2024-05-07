// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package system

import (
	rand "math/rand/v2"
	"time"
)

func Retry(fn func() error) (e error) {
	for i := 8; i <= 128; i *= 2 {
		e = fn()
		if e == nil {
			return nil
		}
		time.Sleep(time.Millisecond * time.Duration(i+rand.IntN(i)))
	}
	return e // the last error
}
