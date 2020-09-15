// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package meta handles the database metadata.
// The metadata is split into two parts - info and schema.
// Info is the fast changing stats and indexes,
// that change with every update.
// Schema is the slow changing columns and indexes
// that are changed infrequently.
package meta
