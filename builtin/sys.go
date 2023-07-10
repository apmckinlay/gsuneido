// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"errors"
	"io/fs"
	"log"
	"net"
	"os"
	"runtime"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var sysMem = systemMemory() // cache

var _ = builtin(SystemMemory, "()")

func SystemMemory() Value {
	return Int64Val(int64(sysMem))
}

var _ = builtin(MemoryArena, "()")

func MemoryArena() Value {
	return Int64Val(int64(HeapSys()))
}

func HeapSys() uint64 {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms.HeapSys
}

var _ = builtin(GetCurrentDirectory, "()")

func GetCurrentDirectory() Value {
	dir, err := os.Getwd()
	if err != nil {
		panic("GetCurrentDirectory: " + err.Error())
	}
	return SuStr(dir)
}

// NOTE: temp file is NOT deleted automatically on exit
// (same as cSuneido, but different from jSuneido)
var _ = builtin(GetTempFileName, "(path, prefix)")

func GetTempFileName(path, prefix Value) Value {
	f, err := os.CreateTemp(ToStr(path), ToStr(prefix)+"*.tmp")
	if err != nil {
		panic("GetTempFileName: " + err.Error())
	}
	filename := f.Name()
	f.Close()
	filename = strings.Replace(filename, `\`, `/`, -1)
	return SuStr(filename)
}

var _ = builtin(CreateDirectory, "(dirname)")

func CreateDirectory(arg Value) Value {
	err := os.Mkdir(ToStr(arg), 0755)
	return SuBool(err == nil)
}

var _ = builtin(DeleteFileApi, "(filename)")

func DeleteFileApi(arg Value) Value {
	err := os.Remove(ToStr(arg))
	return SuBool(err == nil)
}

var _ = builtin(FileExistsQ, "(filename)")

func FileExistsQ(arg Value) Value {
	filename := ToStr(arg)
	_, err := os.Stat(filename)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Println("INFO: FileExists?", filename, err)
	}
	return SuBool(err == nil)
}

var _ = builtin(DirExistsQ, "(filename)")

func DirExistsQ(arg Value) Value {
	filename := ToStr(arg)
	info, err := os.Stat(filename)
	if err == nil {
		return SuBool(info.Mode().IsDir())
	}
	if !errors.Is(err, fs.ErrNotExist) {
		log.Println("INFO: DirExists?", filename, err)
	}
	return False
}

var _ = builtin(MoveFile, "(from, to)")

func MoveFile(from, to Value) Value {
	err := os.Rename(ToStr(from), ToStr(to))
	if err != nil {
		panic("MoveFile: " + err.Error())
	}
	return True
}

var _ = builtin(DeleteDir, "(dir)")

func DeleteDir(dir Value) Value {
	dirname := ToStr(dir)
	info, err := os.Stat(dirname)
	if errors.Is(err, os.ErrNotExist) {
		return False
	}
	if err != nil {
		panic("DeleteDir: " + err.Error())
	}
	if !info.Mode().IsDir() {
		return False
	}
	err = os.RemoveAll(dirname)
	if err != nil {
		panic("DeleteDir: " + err.Error())
	}
	return True
}

var _ = builtin(GetMacAddresses, "()")

func GetMacAddresses() Value {
	ob := &SuObject{}
	if intfcs, err := net.Interfaces(); err == nil {
		for _, intfc := range intfcs {
			if s := string(intfc.HardwareAddr); s != "" {
				ob.Add(SuStr(s))
			}
		}
	}
	return ob
}

var _ = builtin(GetTempPath, "()")

func GetTempPath() Value {
	s := os.TempDir()
	s = strings.ReplaceAll(s, `\`, `/`)
	if !strings.HasSuffix(s, "/") {
		s += "/"
	}
	return SuStr(s)
}

var _ = builtin(OSName, "()")

func OSName() Value {
	os := runtime.GOOS
	os = strings.Replace(os, "darwin", "macos", 1)
	return SuStr(os)
}
