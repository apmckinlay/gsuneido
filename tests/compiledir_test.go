// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build interactive

package tests

// compiles a directory tree of files
// e.g. as exported by LibToFiles
// for testing the compiler
// run with: go test -tags slow -run TestCompileDir

import (
	"fmt"
	iou "io/ioutil"
	"os"
	filepath "path/filepath"
	"strings"
	"testing"

	_ "github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile"
)

func TestCompileDir(*testing.T) {
	filepath.Walk("../../stdlib", walk)
	fmt.Println("TOTAL SIZE", totalSize)
}

var totalSize = 0

func walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	if strings.Contains(path, "/Win32/") ||
		strings.HasSuffix(path, ".css") ||
		strings.HasSuffix(path, ".js") ||
		strings.HasSuffix(path, ".DS_Store") {
		return nil
	}
	// fmt.Println(path)
	data, _ := iou.ReadFile(path)
	text := string(data)
	if strings.Contains(text, "struct") ||
		strings.Contains(text, "dll") ||
		strings.Contains(text, "callback") {
		return nil
	}
	totalSize += len(text)
	// Constant(text)
	_, results := compile.Checked(nil, text)
	if len(results) > 0 {
		fmt.Println(path)
		for _, s := range results {
			fmt.Println("   ", s)
		}
		//return fmt.Errorf("errors")
	}

	// e := assert.Catch(func() { Constant(text) })
	// if e != nil {
	// 	fmt.Println("\t", e)
	// 	return fmt.Errorf("%v", e)
	// }
	// if e == nil {
	// 	fmt.Println(path)
	// }
	return nil
}
