// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"iter"
	"log"
	"os"
	"path/filepath"
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
)

const maxDir = 10000

var _ = builtin(dir, "(path='*', files=false, details=false, block=false)")

func dir(th *Thread, args []Value) Value {
	path := ToStr(args[0])
	justfiles := ToBool(args[1])
	details := ToBool(args[2])
	block := args[3]
	if block == False {
		ob := &SuObject{}
		for entry := range dirEntries(path, justfiles, details) {
			if ob.Size() >= maxDir {
				logPanic("ERROR: Dir: too many files")
			}
			ob.Add(entry)
		}
		return ob
	}
	// block form
	dirEntries(path, justfiles, details)(
		func(entry Value) (ret bool) {
			defer func() {
				if e := recover(); e != nil && e != BlockContinue {
					panic(e)
				}
				ret = true
			}()
			th.Call(block, entry)
			return true
		})
	return nil
}

func dirEntries(dir string, justfiles, details bool) iter.Seq[Value] {
	return func(yield func(Value) bool) {
		dir, pat := filepath.Split(dir)
		if dir == "" {
			dir = "."
		}
		if strings.HasSuffix(pat, "*.*") {
			pat = pat[:len(pat)-2] // switch *.* to *
		}
		match := newMatcher(pat)
		f, err := os.Open(dir)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				logPanic("ERROR: Dir:", err)
			}
			return
		}
		defer func() {
			f.Close()
			if e := recover(); e != nil && e != BlockBreak {
				panic(e)
			}
		}()
		for {
			list, err := f.ReadDir(100)
			if err != nil {
				if err != io.EOF {
					logPanic("ERROR: Dir:", err)
				}
				break
			}
			for _, ent := range list {
				name := ent.Name()
				if match(name) && (!justfiles || !ent.IsDir()) {
					suffix := ""
					if ent.IsDir() {
						suffix = "/"
					}
					var entry Value = SuStr(ent.Name() + suffix)
					if details {
						info, err := ent.Info()
						if err != nil {
							logPanic("ERROR: Dir:", err)
						}
						ob := &SuObject{}
						ob.Set(SuStr("name"), entry)
						ob.Set(SuStr("size"), Int64Val(info.Size()))
						ob.Set(SuStr("date"), FromGoTime(info.ModTime()))
						entry = ob
					}
					if !yield(entry) {
						return
					}
				}
			}
		}
	}
}

type matcher func(name string) bool

func newMatcher(pat string) matcher {
	if strings.Count(pat, "*") > 1 {
		panic("Dir only handles one '*'")
	}
	if strings.Contains(pat, "?") {
		panic("Dir does not handle '?'")
	}
	before, after, found := strings.Cut(pat, "*")
	if !found {
		return func(name string) bool {
			return pat == name
		}
	}
	return func(name string) bool {
		return strings.HasPrefix(name, before) && strings.HasSuffix(name, after)
	}
}

func logPanic(args ...any) {
	// log should include ERROR, panic should not
	s := fmt.Sprintln(args...)
	log.Println(s)
	panic(strings.TrimPrefix(s, "ERROR: "))
}
