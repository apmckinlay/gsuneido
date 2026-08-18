// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package annotations

import "testing"

func TestParamBaseName(t *testing.T) {
	cases := []struct{ in, want string }{
		{"x", "x"},
		{"foo", "foo"},
		{"_x", "x"},
		{".x", "x"},
		{".X", "x"},
		{".Name", "name"},
		{"._x", "x"},
		{"._Name", "name"},
		{"", ""},
	}
	for _, c := range cases {
		if got := paramBaseName(c.in); got != c.want {
			t.Errorf("paramBaseName(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestParamIsDynamic(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"_x", true},
		{"._x", true},
		{"x", false},
		{".x", false},
		{".X", false},
		{"", false},
	}
	for _, c := range cases {
		if got := paramIsDynamic(c.in); got != c.want {
			t.Errorf("paramIsDynamic(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}
