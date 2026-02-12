// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/core"
)

var sandboxEnabled atomic.Bool
var sandboxRoot string

func EnableSandbox() {
	root, err := os.Getwd()
	if err != nil {
		core.Fatal(err)
	}
	enableSandbox(root)
}

func enableSandbox(root string) {
	sandboxRoot = filepath.Clean(root)
	sandboxEnabled.Store(true)
}

func resetSandbox() {
	sandboxEnabled.Store(false)
	sandboxRoot = ""
}

func sandboxed() bool {
	return sandboxEnabled.Load()
}

func guardSandbox(op string) {
	if sandboxEnabled.Load() {
		panic("sandbox: " + op + " disabled")
	}
}

func sandboxPath(op, name string) (string, error) {
	if !sandboxEnabled.Load() {
		return name, nil
	}
	if filepath.IsAbs(name) {
		return "", fmt.Errorf("%s: absolute paths disabled in sandbox", op)
	}
	if filepath.VolumeName(name) != "" {
		return "", fmt.Errorf("%s: volume paths disabled in sandbox", op)
	}
	if strings.HasPrefix(name, string(filepath.Separator)) {
		return "", fmt.Errorf("%s: root-relative paths disabled in sandbox", op)
	}
	if hasDotDot(name) {
		return "", fmt.Errorf("%s: parent paths disabled in sandbox", op)
	}
	root := sandboxRoot
	if root == "" {
		return "", fmt.Errorf("%s: sandbox root not set", op)
	}
	cleaned := filepath.Clean(name)
	full := filepath.Join(root, cleaned)
	rel, err := filepath.Rel(root, full)
	if err != nil {
		return "", fmt.Errorf("%s: invalid path", op)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("%s: path outside sandbox", op)
	}
	return full, nil
}

func hasDotDot(name string) bool {
	path := filepath.ToSlash(name)
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if part == ".." {
			return true
		}
	}
	return false
}
