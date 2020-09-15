// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

// Merge combines an fbtree with an mbtree to produce a new fbtree.
// It does not modify the original fbtree or mbtree.
func Merge(fb *fbtree, mb *mbtree) *fbtree {
	return fb.Update(func(mfb *fbtree) {
		mb.ForEach(func(key string, off uint64) {
			if (off & tombstone) == 0 {
				mfb.Insert(key, off)
			} else {
				mfb.Delete(key, off&^tombstone)
			}
		})
	})
}
