// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"errors"
	"log"
	"net"
	"os"
	"runtime"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("MemoryArena()", func() Value {
	return Int64Val(int64(HeapSys()))
})

func HeapSys() uint64 {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms.HeapSys
}

var _ = builtin0("GetCurrentDirectory()",
	func() Value {
		dir, err := os.Getwd()
		if err != nil {
			panic("GetCurrentDirectory: " + err.Error())
		}
		return SuStr(dir)
	})

// NOTE: temp file is NOT deleted automatically on exit
// (same as cSuneido, but different from jSuneido)
var _ = builtin2("GetTempFileName(path, prefix)",
	func(path, prefix Value) Value {
		f, err := os.CreateTemp(ToStr(path), ToStr(prefix)+"*.tmp")
		if err != nil {
			panic("GetTempFileName: " + err.Error())
		}
		filename := f.Name()
		f.Close()
		filename = strings.Replace(filename, `\`, `/`, -1)
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
		filename := ToStr(arg)
		_, err := os.Stat(filename)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Println("INFO: FileExists?", filename, err)
		}
		return SuBool(err == nil)
	})

var _ = builtin1("DirExists?(filename)",
	func(arg Value) Value {
		info, err := os.Stat(ToStr(arg))
		if err == nil {
			return SuBool(info.Mode().IsDir())
		}
		if os.IsNotExist(err) {
			return False
		}
		panic("DirExists?: " + err.Error())
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
	})

var _ = builtin0("GetMacAddresses()",
	func() Value {
		ob := &SuObject{}
		if intfcs, err := net.Interfaces(); err == nil {
			for _, intfc := range intfcs {
				if s := string(intfc.HardwareAddr); s != "" {
					ob.Add(SuStr(s))
				}
			}
		}
		return ob
	})
