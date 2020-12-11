// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hamttest

import "github.com/apmckinlay/gsuneido/db19/stor"

//go:generate genny -in ../../genny/hamt/hamt.go -out testhamt.go -pkg hamttest gen "Item=*Foo KeyType=int"

type Foo struct {
	key  int
	data string
}

func FooKey(foo *Foo) int {
	return foo.key
}

func FooHash(key int) uint32 {
	return uint32(key) & 0xffff // reduce bits to force overflows
}

func (f *Foo) storSize() int {
	return 0
}

func (f *Foo) Write(*stor.Writer) {
}

func ReadFoo(*stor.Stor, *stor.Reader) *Foo {
	return nil
}
