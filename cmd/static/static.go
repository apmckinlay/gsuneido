// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		usage()
		return
	}
	if len(args) == 1 && args[0] == "repl" {
		db := buildGraph()
		defer db.Close()
		printCounts()
		repl()
		return
	}
	if len(args) == 1 && args[0] == "between" {
		db := buildGraph()
		defer db.Close()
		printCounts()
		run("roots", "between", -1, false)
		return
	}
	if len(args) == 1 {
		db := buildGraph()
		defer db.Close()
		printCounts()
		printEdges(args[0])
		return
	}
	from, to, minCount, sortByCount, ok := parseCommand(args)
	if !ok {
		usage()
		os.Exit(1)
	}
	db := buildGraph()
	// can't close the database because strings are in the memory map
	defer db.Close()
	printCounts()
	run(from, to, minCount, sortByCount)
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage:
  static
  static {name}
  static between
  static {from} to {to}
  static {from} to counts^{min}
  static repl

{from} is "roots", "tables", or a node name
{to} is "tables", "counts", "between", or a node name
{name} prints its incoming and outgoing edges
between prints the top 100 nodes on the most shortest paths to tables
counts is sorted by name
counts^ is sorted by count
counts^{min} also only shows entries with a count of at least {min}`)
}

// parseCommand parses "{from} to {to}" or "{from} to counts^[{min}]".
// minCount is -1 when no filter was given. sortByCount is set for counts^.
func parseCommand(args []string) (from, to string, minCount int, sortByCount, ok bool) {
	if len(args) != 3 || args[1] != "to" {
		return "", "", 0, false, false
	}
	from, to = args[0], args[2]
	if rest, found := strings.CutPrefix(to, "counts^"); found {
		if rest == "" {
			return from, "counts", -1, true, true
		}
		n, err := strconv.Atoi(rest)
		if err != nil {
			return "", "", 0, false, false
		}
		return from, "counts", n, true, true
	}
	return from, to, -1, false, true
}

// splitFields splits on whitespace, allowing double-quoted node names.
func splitFields(line string) []string {
	var fields []string
	var field strings.Builder
	quoted := false
	inField := false
	for _, r := range line {
		switch {
		case r == '"':
			quoted = !quoted
			inField = true
		case !quoted && (r == ' ' || r == '\t'):
			if inField {
				fields = append(fields, field.String())
				field.Reset()
				inField = false
			}
		default:
			field.WriteRune(r)
			inField = true
		}
	}
	if inField {
		fields = append(fields, field.String())
	}
	return fields
}

// repl reads command lines until a line with just "q".
func repl() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "q" {
			return
		}
		if line == "" {
			continue
		}
		from, to, minCount, sortByCount, ok := parseCommand(splitFields(line))
		if !ok {
			fmt.Fprintf(os.Stderr, "invalid command: %s\n", line)
			usage()
			continue
		}
		run(from, to, minCount, sortByCount)
		fmt.Fprintln(os.Stderr)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
