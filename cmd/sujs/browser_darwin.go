// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build darwin

package main

import (
	"os"
	"path/filepath"
)

// macOS browsers installed as .app bundles are not in PATH.
// Check known /Applications locations directly.
var macAppBrowsers = []string{
	// Chromium / Chrome
	"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
	"/Applications/Chromium.app/Contents/MacOS/Chromium",
	"/Applications/Ungoogled Chromium.app/Contents/MacOS/Chromium",
	// Edge
	"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
	// Brave
	"/Applications/Brave Browser.app/Contents/MacOS/Brave Browser",
	// Vivaldi
	"/Applications/Vivaldi.app/Contents/MacOS/Vivaldi",
	// Thorium
	"/Applications/Thorium.app/Contents/MacOS/Thorium",
	// Opera
	"/Applications/Opera.app/Contents/MacOS/Opera",
	// Arc
	"/Applications/Arc.app/Contents/MacOS/Arc",
	// Helium
	"/Applications/Helium.app/Contents/MacOS/Helium",
}

// findBrowserApp checks /Applications .app bundles before falling back
// to the PATH-based search used on other platforms.
func findBrowserApp() string {
	for _, p := range macAppBrowsers {
		if _, err := os.Stat(p); err == nil {
			return filepath.Clean(p)
		}
	}
	return ""
}
