// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package system

import "os"

func RenameBak(from string, to string) error {
	err := Retry(func() error { return os.Remove(to + ".bak") })
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	err = Retry(func() error { return os.Rename(to, to+".bak") })
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	err = Retry(func() error { return os.Rename(from, to) })
	if err != nil {
		return err
	}
	return nil
}
