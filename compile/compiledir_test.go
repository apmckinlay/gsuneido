// +build slow

package compile

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

//	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestCompileDir(t *testing.T) {
	filepath.Walk("../../stdlib", walk)
}

func walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	if strings.Contains(path, "/Win32/") {
		return nil
	}
	fmt.Println(path)
	data, _ := iou.ReadFile(path)
	text := string(data)
	Constant(text)
	// e := Catch(func() { Constant(text) })
	// if e != nil {
	// 	fmt.Println("\t", e)
	// 	return fmt.Errorf("%v", e)
	// }
	// if e == nil {
	// 	fmt.Println(path)
	// }
	return nil
}
