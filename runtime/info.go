// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"strconv"
	"sync/atomic"

	"golang.org/x/exp/maps"
)

var infos = map[string]any{}

// AddInfo registers a reference to a value.
// It is NOT thread-safe and should only be called during initialization.
// The return type is to allow var _ = AddInfo("name", ref)
func AddInfo(name string, ref any) struct{} {
	infos[name] = ref
	return struct{}{}
}

func InfoStr(name string) string {
	x, ok := infos[name]
	if !ok {
		return ""
	}
	switch x := x.(type) {
	case *atomic.Int64:
		return strconv.FormatInt(x.Load(), 10)
	}
	panic("unknown info type")
}

func InfoList() *SuObject {
	return SuObjectOfStrs(maps.Keys(infos))
}
