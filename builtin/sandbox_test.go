// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSandboxFilePath(t *testing.T) {
	resetSandbox()
	path, err := sandboxPath("test", "file.txt")
	assert.T(t).That(err == nil)
	assert.T(t).This(path).Is("file.txt")

	root := t.TempDir()
	enableSandbox(root)
	defer resetSandbox()

	path, err = sandboxPath("test", "sub/dir/file.txt")
	assert.T(t).That(err == nil)
	assert.T(t).This(path).Is(filepath.Join(root, "sub/dir/file.txt"))

	_, err = sandboxPath("test", string(filepath.Separator)+"abs")
	assert.T(t).That(err != nil)
	assert.T(t).That(strings.Contains(err.Error(), "paths disabled in sandbox"))

	_, err = sandboxPath("test", "../up")
	assert.T(t).This(err.Error()).Is("test: parent paths disabled in sandbox")

	_, err = sandboxPath("test", "sub/../file.txt")
	assert.T(t).This(err.Error()).Is("test: parent paths disabled in sandbox")
}

func TestGuardSandbox(t *testing.T) {
	resetSandbox()
	guardSandbox("System")

	root := t.TempDir()
	enableSandbox(root)
	defer resetSandbox()

	assert.T(t).This(func() { guardSandbox("System") }).Panics("sandbox: System disabled")
}

func TestSandboxPath(t *testing.T) {
	resetSandbox()
	root := t.TempDir()
	enableSandbox(root)
	defer resetSandbox()

	_, err := sandboxPath("CreateDir", string(filepath.Separator)+"abs")
	assert.T(t).That(err != nil)
	assert.T(t).That(strings.Contains(err.Error(), "paths disabled in sandbox"))
}
