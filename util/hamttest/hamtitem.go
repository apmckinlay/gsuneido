// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hamttest

//go:generate genny -in ../../genny/hamt/hamt.go -out testhamt.go -pkg hamttest gen "Item=Foo KeyType=int"

type Foo struct {
	key  int
	data string
}

func (foo *Foo) Key() int {
	return foo.key
}

func FooHash(key int) uint32 {
	return uint32(key) & 0xffff // reduce bits to force overflows
}
