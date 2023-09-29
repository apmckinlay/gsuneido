// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package atomics

import goatomic "sync/atomic"

type String struct {
	v goatomic.Value
}

func (as *String) Store(s string) {
	as.v.Store(s)
}

func (as *String) Load() string {
	x := as.v.Load()
	if x == nil {
		return ""
	}
	return x.(string)
}
