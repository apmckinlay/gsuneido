// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"sync/atomic"

	"golang.org/x/exp/maps"
)

var infos = map[string]any{}

// AddInfo registers a reference to a value.
// It is NOT thread-safe and should only be called during initialization.
// ref's should be thread-safe.
// The return type is to allow var _ = AddInfo("name", ref)
func AddInfo(name string, ref any) struct{} {
	infos[name] = ref
	return struct{}{}
}

func InfoStr(name string) Value {
	x, ok := infos[name]
	if !ok {
		return nil
	}
	switch x := x.(type) {
	case *int:
		return IntVal(*x)
	case *atomic.Int64:
		return Int64Val(x.Load())
	case *atomic.Int32:
		return IntVal(int(x.Load()))
	case func() int:
		return IntVal(x())
	case func() string:
		return SuStr(x())
	case func() Value:
		return x()
	}
	panic("unknown info type")
}

func InfoList() []string {
	return maps.Keys(infos)
}
