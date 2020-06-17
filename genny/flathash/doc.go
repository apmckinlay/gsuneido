// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package flathash implements a simple hash map
// with linear probing open addressing.
// The flat structure means we can duplicate efficiently for copy-on-write.
// The key type must support ==.
// The zero value of keys is used to identify noKey slots.
// The zero value of values is returned for failed searched.
// Instantiations must define hash and keyOf to complete the generated code.
package flathash

