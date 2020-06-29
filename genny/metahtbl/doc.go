// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package metahtbl implements a simple hash map
// with linear probing open addressing.
// The flat structure means we can duplicate efficiently for copy-on-write.
// It is specifically written for meta Info and Schema, not general purpose.
package metahtbl

