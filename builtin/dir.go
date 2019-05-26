package builtin

import (
	"os"
	"path"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

const maxDir = 10000

var _ = builtin("Dir(path='*', files=false, details=false, block=false)",
	func(t *Thread, args []Value) Value {
		path := strings.ReplaceAll(ToStr(args[0]), "\\", "/")
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
			t.CallWithArgs(block, entry)
		})
		return nil
	})

var name = SuStr("name")

func forEachDir(dir string, justfiles, details bool, fn func(entry Value)) {
	dir, pat := path.Split(dir)
	if dir == "" {
		dir = "."
	}
	if pat == "" {
		pat = "*"
	} else if strings.HasSuffix(pat, "*.*") {
		pat = pat[:len(pat)-3] + "*"
	}
	f, err := os.Open(dir)
	defer f.Close()
	if err != nil {
		panic("Dir: " + err.Error())
	}
	list, err := f.Readdir(100)
	for _, info := range list {
		name := info.Name()
		match, _ := path.Match(pat, name)
		if match && (!justfiles || !info.IsDir()) {
			suffix := ""
			if info.IsDir() {
				suffix = "/"
			}
			var entry Value = SuStr(info.Name() + suffix)
			if details {
				ob := &SuObject{}
				ob.Set(SuStr("name"), entry)
				ob.Set(SuStr("size"), Int64Val(info.Size()))
				ob.Set(SuStr("date"), FromTime(info.ModTime()))
				entry = ob
			}
			fn(entry)
		}
	}
}
