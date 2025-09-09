// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

/*
slot1 consists of:
- 1 byte number of entries
- array of:
	= 1 byte size
	- 5 byte database offset
- followed by the contiguous key data (sorted)
*/

type slot1 []byte

// Len returns the number of entries in the slot
func (s slot1) Len() int {
	return int(s[0])
}

// Search does a linear search for the key and 
// returns the index where the key is found, or where it should be inserted.
func (s slot1) Search(key string) int {
	n := s.Len()
	pos := 1
	keyDataStart := 1 + n*6 
	keyDataPos := keyDataStart
	for i := range n {
		keySize := int(s[pos])
		pos += 6
		entryKey := s[keyDataPos : keyDataPos+keySize]
		keyDataPos += keySize
		if string(entryKey) >= key {
			return i 
		}
	}
	return n 
}
