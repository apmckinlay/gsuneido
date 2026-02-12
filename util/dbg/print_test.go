// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbg

func ExamplePrint_string() {
	PrintString("hello")
	PrintSlice([]byte("hello"))
	PrintString("01234567890123456789012345678901234567890123456789")
	// Output:
	// 68656c6c6f
	//  h e l l o
	// 68656c6c6f
	//  h e l l o
	// 3031323334353637383930313233343536373839303132333435363738393031323334 ...
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 ...
}
