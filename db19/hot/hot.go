// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package hot tracks stats for frequently used (hot) columns
//
//  1. "busy" tracks columns that are frequently use in WHERE clauses.
//     This is appended to a text file either every two hours or on exit.
//  2. "stats" tracks data cardinality (HLL), common values (SS), and quantiles (KLL)
//     This data is collected by compact for the "busy" columns.
//
// "tally" as in BusyTally and StatsTally refers to the data collection.
package hot
