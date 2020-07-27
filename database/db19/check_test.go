// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestCheckStartStop(t *testing.T) {
	ck := NewCheck()
	const ntrans = 20
	var trans [ntrans]*UpdateTran
	const ntimes = 5000
	for i := 0; i < ntimes; i++ {
		j := rand.Intn(ntrans)
		if trans[j] == nil {
			trans[j] = &UpdateTran{ct: ck.StartTran()}
		} else {
			if rand.Intn(2) == 1 {
				ck.Commit(trans[j])
			} else {
				ck.Abort(trans[j].ct)
			}
			trans[j] = nil
		}
	}
	for _, tn := range trans {
		if tn != nil {
			ck.Commit(tn)
		}
	}
	Assert(t).That(len(ck.trans), Equals(0))
}

func TestCheckLimit(t *testing.T) {
	ck := NewCheck()
	for i := 0; i < maxTrans; i++ {
		Assert(t).True(ck.StartTran() != nil)
	}
	Assert(t).True(ck.StartTran() == nil)
}

func TestCheckActions(t *testing.T) {
	checkerAbortT1 = true
	defer func() { checkerAbortT1 = false }()
	// writes
	script(t, "1w1 2w2 1c 2c")
	script(t, "1w4 1w5 2w6 2w7 1c 2c")
	script(t, "1w1 2w2 1c 2a")
	script(t, "1w1 2w2 1a 2c")
	script(t, "1w1 2w2 1a 2a")
	// conflict
	script(t, "1w1 2W1 1c")
	script(t, "1w1 1a 2w1 2c")
	script(t, "1w4 1w5 2w3 2W5")
	// conflict with ended
	script(t, "1w1 1c 2W1")
	script(t, "2w1 2c 1W1 1C")

	// reads
	script(t, "1w4 1r68 2r77 2R35")
	script(t, "1r35 2W4")
}

func script(t *testing.T, s string) {
	t.Helper()
	ok := func(result bool) {
		t.Helper()
		if result != true {
			t.Log("incorrect at:", s)
			t.FailNow()
		}
	}
	fail := func(result bool) {
		t.Helper()
		if result != false {
			t.Log("incorrect at:", s)
			t.FailNow()
		}
	}
	ck := NewCheck()
	ts := []*UpdateTran{{ct: ck.StartTran()}, {ct: ck.StartTran()}}
	for len(s) > 0 {
		t := ts[s[0]-'1']
		switch s[1] {
		case 'w':
			ok(ck.Write(t.ct, "mytable", []string{"", s[2:3]}))
			s = s[1:]
		case 'W':
			fail(ck.Write(t.ct, "mytable", []string{"", s[2:3]}))
			s = s[1:]
		case 'r':
			ok(ck.Read(t.ct, "mytable", 1, s[2:3], s[3:4]))
			s = s[2:]
		case 'R':
			fail(ck.Read(t.ct, "mytable", 1, s[2:3], s[3:4]))
			s = s[2:]
		case 'c':
			ok(ck.Commit(t))
		case 'C':
			fail(ck.Commit(t))
		case 'a':
			ok(ck.Abort(t.ct))
		case 'A':
			fail(ck.Abort(t.ct))
		}
		s = s[2:]
		for len(s) > 0 && s[0] == ' ' {
			s = s[1:]
		}
	}
}

func (t *CkTran) String() string {
	b := new(bytes.Buffer)
	fmt.Fprint(b, "T", t.start)
	if t.isEnded() {
		fmt.Fprint(b, "->", t.end)
	}
	fmt.Fprintln(b)
	for name, tbl := range t.tables {
		fmt.Fprintln(b, "    ", name)
		for i, writes := range tbl.writes {
			if writes != nil {
				fmt.Fprintln(b, "        writes", i, writes.String())
			}
		}
		for i, reads := range tbl.reads {
			if reads != nil {
				fmt.Fprintln(b, "        reads", i, reads.String())
			}
		}
	}
	return strings.TrimSpace(b.String())
}
