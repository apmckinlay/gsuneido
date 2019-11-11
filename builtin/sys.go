package builtin

import (
	"io/ioutil"
	"os"
	"runtime"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("MemoryArena()", func() Value {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return Int64Val(int64(ms.HeapSys))
})

var _ = builtin0("GetCurrentDirectory()",
	func() Value {
		dir, err := os.Getwd()
		if err != nil {
			panic("GetCurrentDirectory: " + err.Error())
		}
		return SuStr(dir)
	})

var _ = builtin0("GetTempPath()",
	func() Value {
		return SuStr(os.TempDir())
	})

// NOTE: temp file is NOT deleted automatically on exit
// (same as cSuneido, but different from jSuneido)
var _ = builtin2("GetTempFileName(path, prefix)",
	func(path, prefix Value) Value {
		f, err := ioutil.TempFile(ToStr(path), ToStr(prefix))
		if err != nil {
			panic("GetTempFileName: " + err.Error())
		}
		filename := f.Name()
		f.Close()
		return SuStr(filename)
	})

var _ = builtin1("CreateDirectory(dirname)",
	func(arg Value) Value {
		err := os.Mkdir(ToStr(arg), 0755)
		return SuBool(err == nil)
	})

var _ = builtin1("DeleteFileApi(filename)",
	func(arg Value) Value {
		err := os.Remove(ToStr(arg))
		return SuBool(err == nil)
	})

var _ = builtin1("FileExists?(filename)",
	func(arg Value) Value {
		_, err := os.Stat(ToStr(arg))
		return SuBool(err == nil)
	})

var _ = builtin1("DirExists?(filename)",
	func(arg Value) Value {
		info, err := os.Stat(ToStr(arg))
		if err == nil {
			return SuBool((info.Mode() & os.ModeDir) == os.ModeDir)
		}
		if os.IsNotExist(err) {
			return False
		}
		panic("DirExists: " + err.Error())
	})

var _ = builtin2("MoveFile(from, to)",
	func(from, to Value) Value {
		err := os.Rename(ToStr(from), ToStr(to))
		if err != nil {
			panic("MoveFile: " + err.Error())
		}
		return True
	})

var _ = builtin1("DeleteDir(dir)",
	func(dir Value) Value {
		err := os.RemoveAll(ToStr(dir))
		if err != nil {
			panic("DeleteDir: " + err.Error())
		}
		return True
	})
