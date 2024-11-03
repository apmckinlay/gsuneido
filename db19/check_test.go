// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestCheckStartStop(t *testing.T) {
	ck := NewCheck(&Database{})
	const ntrans = 20
	var trans [ntrans]*UpdateTran
	const ntimes = 5000
	for range ntimes {
		j := rand.Intn(ntrans)
		if trans[j] == nil {
			trans[j] = &UpdateTran{ct: ck.StartTran()}
		} else {
			if rand.Intn(2) == 1 {
				ck.Commit(trans[j])
			} else {
				ck.Abort(trans[j].ct, "")
			}
			trans[j] = nil
		}
	}
	for _, tn := range trans {
		if tn != nil {
			ck.Commit(tn)
		}
	}
	assert.T(t).This(ck.count()).Is(0)
}

func TestCheckLimit(t *testing.T) {
	ck := NewCheck(nil)
	for range MaxTrans {
		assert.T(t).That(ck.StartTran() != nil)
	}
	assert.T(t).That(ck.StartTran() == nil)
}

func TestCheckActions(t *testing.T) {
	checkerAbortT1 = true
	defer func() { checkerAbortT1 = false }()
	// writes
	script(t, "1o1 2o2 1c 2c")
	script(t, "1o4 1o5 2o6 2o7 1c 2c")
	script(t, "1o1 2o2 1c 2a")
	script(t, "1o1 2o2 1a 2c")
	script(t, "1o1 2o2 1a 2a")
	// conflict
	script(t, "1d1 2D1 1c")
	script(t, "1o1 1a 2o1 2c")
	script(t, "1d4 1d5 2d3 2D5")
	script(t, "1r55 1o5 2r55 2O8")
	script(t, "1r55 1o5 1c 2r55 2O8")
	script(t, "1r55 2r55 2o8 1O5")
	script(t, "1r57 1o9 2D6")
	// conflict with ended
	script(t, "1d1 1c 2D1")
	script(t, "2d1 2c 1D1 1C")
	script(t, "1r3 2d3 2c 2D3")

	// reads
	script(t, "1o4 1r68 2o3 2r77 2R35")
	script(t, "1o8 1r35 2O4")

	// don't check writes against committed reads
	script(t, "1r11 1c 2o1 2c")
	// but still check reads against committed writes
	script(t, "2o1 2c 1o9 1R11")

	// delayed read conflicts
	script(t, "1r24 2o3 1c")
	script(t, "1r24 2o3 1D3")
}

// script takes a string containing a space separated list of actions.
// Each action consists of:
//   - transaction number 1 or 2
//   - action type: (r)ead, (o)utput, (d)elete, (c)ommit, (a)bort
//   - read is followed by two characters specifying a key range
//   - write is followed by one character specifying a key
//
// If the type is capitalized (R, W, C, A) then the action is expected to fail
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
	ck := NewCheck(&Database{})
	ts := []*UpdateTran{{ct: ck.StartTran()}, {ct: ck.StartTran()}}
	for len(s) > 0 {
		t := ts[s[0]-'1']
		switch s[1] {
		case 'o':
			ok(ck.Output(t.ct, "mytable", []string{s[2:3]}))
			s = s[1:]
		case 'O':
			fail(ck.Output(t.ct, "mytable", []string{s[2:3]}))
			s = s[1:]
		case 'd':
			off := uint64(s[2] - '0')
			ok(ck.Delete(t.ct, "mytable", off, []string{s[2:3]}))
			s = s[1:]
		case 'D':
			off := uint64(s[2] - '0')
			fail(ck.Delete(t.ct, "mytable", off, []string{s[2:3]}))
			s = s[1:]
		case 'r':
			ok(ck.Read(t.ct, "mytable", 0, s[2:3], s[3:4]))
			s = s[2:]
		case 'R':
			fail(ck.Read(t.ct, "mytable", 0, s[2:3], s[3:4]))
			s = s[2:]
		case 'c':
			ok(ck.Commit(t))
		case 'C':
			fail(ck.Commit(t))
		case 'a':
			ok(ck.Abort(t.ct, ""))
		case 'A':
			fail(ck.Abort(t.ct, ""))
		}
		s = s[2:]
		for len(s) > 0 && s[0] == ' ' {
			s = s[1:]
		}
	}
}

func TestCheckMax(t *testing.T) {
	ck := NewCheck(nil)
	ct := ck.StartTran()
	randTable := str.UniqueRandom(3, 10)
	for range readMax {
		assert.True(ck.Read(ct, randTable(), 0, "bar", "foo"))
	}
	assert.False(ck.Read(ct, randTable(), 0, "bar", "foo"))
}

// func (t *CkTran) String() string {
// 	b := new(bytes.Buffer)
// 	fmt.Fprint(b, "T", t.start)
// 	if t.isEnded() {
// 		fmt.Fprint(b, "->", t.end)
// 	}
// 	fmt.Fprintln(b)
// 	for name, tbl := range t.tables {
// 		fmt.Fprintln(b, "    ", name)
// 		for i, writes := range tbl.writes {
// 			if writes != nil {
// 				fmt.Fprintln(b, "        writes", i, writes.String())
// 			}
// 		}
// 		for i, reads := range tbl.reads {
// 			if reads != nil {
// 				fmt.Fprintln(b, "        reads", i, reads.String())
// 			}
// 		}
// 	}
// 	return strings.TrimSpace(b.String())
// }

func BenchmarkCheck(b *testing.B) {
	const ntrans = 400
	const ntables = 7 // per transaction
	tables := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
		"k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w",
		"x", "y", "z"}
	const nacts = 5 // per table
	const nindexes = 5
	keys := make([]string, nindexes)
	makeKeys := func() []string {
		for i := range nindexes {
			keys[i] = strconv.Itoa(rand.Intn(9999))
		}
		return keys
	}
	const keyRange = 999999

	var ck *Check
	b.ResetTimer()
	for i := range b.N {
		if i%ntrans == 0 {
			ck = NewCheck(nil)
		}
		ct := ck.StartTran()
		for range ntables {
			table := tables[rand.Intn(len(tables))]
			for range nacts {
				switch rand.Intn(3) {
				case 0:
					index := rand.Intn(nindexes)
					n := rand.Intn(keyRange)
					from := strconv.Itoa(n)
					to := strconv.Itoa(n + rand.Intn(99))
					ck.Read(ct, table, index, from, to)
				case 1:
					ck.Output(ct, table, makeKeys())
				case 2:
					offset := uint64(rand.Intn(keyRange)) + 1 // not zero
					ck.Delete(ct, table, offset, makeKeys())
				}
			}
		}
	}
}
