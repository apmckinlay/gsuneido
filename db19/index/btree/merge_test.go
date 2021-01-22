// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/db19/index/testdata"
	"github.com/apmckinlay/gsuneido/db19/stor"
)

func TestMergeAndSave(*testing.T) {
	nMerges := 2000
	opsPerMerge := 1000
	if testing.Short() {
		nMerges = 200
		opsPerMerge = 200
	}
	d := testdata.New()
	GetLeafKey = d.GetLeafKey
	defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
	MaxNodeSize = 64
	st := stor.HeapStor(8192)
	st.Alloc(1) // avoid offset 0
	bt := CreateBtree(st, nil)

	for i := 0; i < nMerges; i++ {
		_ = t && trace("---")
		x := &ixbuf.T{}
		for j := 0; j < opsPerMerge; j++ {
			k := rand.Intn(4)
			switch {
			case k == 0 || k == 1 || d.Len() == 0:
				x.Insert(d.Gen())
			case k == 2:
				_, key, _ := d.Rand()
				off := d.NextOff()
				x.Update(key, off)
				d.Update(key, off)
			case k == 3:
				i, key, off := d.Rand()
				x.Delete(key, off)
				d.Delete(i)
			}
		}
		bt = bt.MergeAndSave(x.Iter())
	}
	bt.Check(nil)
	d.Check(bt)
	d.CheckIter(bt.Iterator())
}

//-------------------------------------------------------------------

func (st *state) print() {
	fmt.Println("state:", st.bt.treeLevels)
	for _, m := range st.path {
		fmt.Println("   ", &m)
		fmt.Println("       ", m.node.knowns())
	}
}

func (m *merge) String() string {
	limit := m.limit
	if limit == "" {
		limit = `""`
	}
	mod := ""
	if m.modified {
		mod = " modified"
	}
	return fmt.Sprint("off ", m.off, " pos ", m.pos, " limit ", limit, mod)
}
