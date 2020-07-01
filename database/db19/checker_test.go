// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"math/rand"
	"testing"

	"github.com/apmckinlay/gsuneido/util/verify"
)

func TestCheckerStartStop(*testing.T) {
	ck := NewCheck()
	const ntrans = 20
	var trans [ntrans]int
	const ntimes = 5000
	for i := 0; i < ntimes; i++ {
		j := rand.Intn(ntrans)
		if trans[j] == 0 {
			trans[j] = ck.StartTran()
		} else {
			if rand.Intn(2) == 1 {
				ck.Commit(trans[j])
			} else {
				ck.Abort(trans[j])
			}
			trans[j] = 0
		}
	}
	// fmt.Println("#", len(ck.trans))
	for _, tn := range trans {
		if tn != 0 {
			ck.Commit(tn)
		}
	}
	verify.That(len(ck.trans) == 0)
}

func TestCheckerDelete(t *testing.T) {
	script(t, "1d1 2d2 1c 2c")
	script(t, "1d1 2d2 1c 2a")
	script(t, "1d1 2d2 1a 2c")
	script(t, "1d1 2d2 1a 2a")
	// conflict
	script(t, "1d1 2d1 1a 2a")
	script(t, "1d1 2d1 1a 2c")
	script(t, "1d1 2d1 1c 2C")
	script(t, "1d1 2d1 1c 2A")
	// conflict with ended
	script(t, "1d1 1c 2D1 2C")
	script(t, "2d1 2c 1D1 1C")
}

func script(t *testing.T, s string) {
	ok := func(result bool) {
		if result != true {
			t.Log("incorrect at:", s)
			t.FailNow()
		}
	}
	fail := func(result bool) {
		if result != false {
			t.Log("incorrect at:", s)
			t.FailNow()
		}
	}
	ck := NewCheck()
	ts := []int{ck.StartTran(), ck.StartTran()}
	for len(s) > 0 {
		t := ts[s[0]-'1']
		switch s[1] {
		case 'd':
			ok(ck.Delete(t, "mytable", uint64(s[2]-'0')))
			s = s[1:]
		case 'D':
			fail(ck.Delete(t, "mytable", uint64(s[2]-'0')))
			s = s[1:]
		case 'c':
			ok(ck.Commit(t))
		case 'C':
			fail(ck.Commit(t))
		case 'a':
			ok(ck.Abort(t))
		case 'A':
			fail(ck.Abort(t))
		}
		s = s[2:]
		for len(s) > 0 && s[0] == ' ' {
			s = s[1:]
		}
	}
}
