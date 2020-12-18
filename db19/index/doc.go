// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

/*
Package index implements gSuneido database indexes.

Indexes store pairs of string keys and uint64 offsets.
Keys must be unique within a particular index.

Suneido style iterators (i.e. Next/Prev/Rewind) have a current key value
and behave as if Next/Prev get the first value greater than or less than
the current value, regardless of modifications to the data.
e.g. the current or next value could be deleted,
or a new next value inserted.
However, their implementations use a current position to be more efficient.
Although the iterators expose a Cur method, this is for internal purposes.
The current value is only really valid immediately after Next or Prev.
At the Suneido level Next and Prev return the new current value.
*/
package index
