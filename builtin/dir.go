// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"errors"
	"io"
	"io/fs"
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
		forEachDir(path, justfiles, details, func(entry Value) {
			if ob.Size() >= maxDir {
				panic("Dir: too many files")
			}
			ob.Add(entry)
		})
		return ob
	}
	// block form
	forEachDir(path, justfiles, details, func(entry Value) {
		th.Call(block, entry)
	})
	return nil
}

func forEachDir(dir string, justfiles, details bool, fn func(entry Value)) {
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
		// should panic, but cSuneido doesn't
		if !errors.Is(err, fs.ErrNotExist) &&
			!strings.Contains(err.Error(), "syntax is incorrect") { // Windows
			log.Println("ERROR: Dir:", err)
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
		list, err := f.Readdir(100)
		if err != nil {
			if err != io.EOF {
				panic(err.Error())
			}
			break
		}
		for _, info := range list {
			name := info.Name()
			if match(name) && (!justfiles || !info.IsDir()) {
				suffix := ""
				if info.IsDir() {
					suffix = "/"
				}
				var entry Value = SuStr(info.Name() + suffix)
				if details {
					ob := &SuObject{}
					ob.Set(SuStr("name"), entry)
					ob.Set(SuStr("size"), Int64Val(info.Size()))
					ob.Set(SuStr("date"), FromGoTime(info.ModTime()))
					entry = ob
				}
				func() {
					defer func() {
						if e := recover(); e != nil && e != BlockContinue {
							panic(e)
						}
					}()
					fn(entry)
				}()
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
