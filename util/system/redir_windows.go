// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package system

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/sys/windows"
)

// Redirect redirects stderr and stdout to a file
// (unless already redirected).
// This is to capture Go errors when running as a service or GUI program.
func Redirect(filename string, getId func() string) error {
	if redirected() {
		return nil
	}
	log.SetFlags(0) // don't need timestamp since we're adding it here
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Fprintln(f, e)
			}
		}()
		rdr := bufio.NewReader(r)
		// copy from pipe to file
		for {
			line, err := rdr.ReadBytes('\n')
			if err == io.EOF {
				break
			}
			io.WriteString(f, time.Now().Format("2006/01/02 15:04:05 ")+
				str.Opt(getId(), " - "))
			f.Write(line)
		}
		f.Close()
	}()

	wh := windows.Handle(w.Fd())
	err = windows.SetStdHandle(windows.STD_ERROR_HANDLE, wh)
	if err != nil {
		return err
	}
	err = windows.SetStdHandle(windows.STD_OUTPUT_HANDLE, wh)
	if err != nil {
		return err
	}
	// redo initialization
	os.Stdout = w
	os.Stderr = w
	log.SetOutput(w)
	syscall.Stdout = syscall.Handle(wh)
	syscall.Stderr = syscall.Handle(wh)
	return nil
}

func redirected() bool {
	handle, _ := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if handle != 0 {
		dwFileType, _ := windows.GetFileType(handle)
		return dwFileType == windows.FILE_TYPE_DISK ||
			dwFileType == windows.FILE_TYPE_PIPE
	}
	return false
}
