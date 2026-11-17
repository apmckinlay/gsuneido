// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSandboxFilePath(t *testing.T) {
	DisableSandbox()
	path, err := sandboxPath("test", "file.txt")
	assert.T(t).That(err == nil)
	assert.T(t).This(path).Is("file.txt")

	root := t.TempDir()
	enableSandbox(root)
	defer DisableSandbox()

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
	DisableSandbox()
	guardSandbox("System")

	root := t.TempDir()
	enableSandbox(root)
	defer DisableSandbox()

	assert.T(t).This(func() { guardSandbox("System") }).Panics("sandbox: System disabled")
}

func TestSandboxPath(t *testing.T) {
	DisableSandbox()
	root := t.TempDir()
	enableSandbox(root)
	defer DisableSandbox()

	_, err := sandboxPath("CreateDir", string(filepath.Separator)+"abs")
	assert.T(t).That(err != nil)
	assert.T(t).That(strings.Contains(err.Error(), "paths disabled in sandbox"))
}

func TestSandboxSymlinkEscape(t *testing.T) {
	DisableSandbox()
	root := t.TempDir()
	outside := t.TempDir()
	assert.T(t).This(os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("x"), 0600)).Is(nil)
	assert.T(t).This(os.Mkdir(filepath.Join(root, "inside"), 0700)).Is(nil)
	if os.Symlink(outside, filepath.Join(root, "dirlink")) != nil {
		t.Skip("symlinks not supported")
	}
	assert.T(t).This(os.Symlink(filepath.Join(outside, "secret.txt"), filepath.Join(root, "filelink"))).Is(nil)
	assert.T(t).This(os.Symlink(filepath.Join(outside, "missing"), filepath.Join(root, "dangling"))).Is(nil)
	assert.T(t).This(os.Symlink("inside", filepath.Join(root, "okdir"))).Is(nil)
	enableSandbox(root)
	defer DisableSandbox()

	for _, name := range []string{"dirlink/secret.txt", "dirlink/new.txt", "filelink", "dangling"} {
		_, err := sandboxPath("test", name)
		assert.T(t).That(err != nil)
		assert.T(t).This(err.Error()).Is("test: path outside sandbox")
	}
	// missing paths and links that stay inside the root are allowed
	for _, name := range []string{"missing/new.txt", "inside/new.txt", "okdir/new.txt"} {
		_, err := sandboxPath("test", name)
		assert.T(t).That(err == nil)
	}
}
