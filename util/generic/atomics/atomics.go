// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package atomics

import "sync/atomic"

type String struct {
	v atomic.Value
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

type Value[T any] struct {
	v atomic.Value
}

func (a *Value[T]) Store(x T) {
	a.v.Store(x)
}

func (a *Value[T]) Load() T {
	x := a.v.Load()
	if x == nil {
		var z T
		return z
	}
	return x.(T)
}
