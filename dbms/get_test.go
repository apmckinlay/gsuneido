// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"testing"
)

func TestExists(t *testing.T) {
	rowToRecord(existsRow, existsHdr)
}
