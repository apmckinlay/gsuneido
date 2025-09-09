// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

/*
slot2 consists of:
- 1 byte number of entries
- array of:
	= 2 byte field offset
	- 5 byte database offset
- 2 byte node size (end of data)
- followed by the contiguous key data (sorted)
*/

type slot2 []byte

// Len returns the number of entries
func (s slot2) Len() int {
	return int(s[0])
}

// getKey returns the i'th field
func (s slot2) getKey(i int) []byte {
	base := 1 + i*7
	fieldPos := uint16(s[base])<<8 | uint16(s[base+1])
	endPos := uint16(s[base+7])<<8 | uint16(s[base+8])
	return s[fieldPos:endPos]
}

// Search does a binary search for the key and
// returns the index where the key is found, or where it should be inserted.
func (s slot2) Search(key string) int {
	low, high := 0, s.Len()-1
	for low <= high {
		mid := (low + high) / 2
		midKey := s.getKey(mid)
		if string(midKey) < key {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return low
}
