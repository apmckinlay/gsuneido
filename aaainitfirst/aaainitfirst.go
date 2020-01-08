// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package aaainitfirst is intended to be imported from main
// and initialized first so that any errors from other initialization
// is seen or logged.
// We want this to be initialized *first*
// therefore we need to be careful what gSuneido packages we import
// because anything we import will be initialized first.
// Also need to be careful that only gsuneido.go imports this package.
// If any other package import this their tests will create logs files.
package aaainitfirst
