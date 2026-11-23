// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"log"
	"math"
	"slices"
	"strings"

	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/util/ascii"

	"github.com/apmckinlay/gsuneido/compile/lexer"
	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
)

var nextId uint16
var nodeMap = map[string]entry{}
var nodeVec []string
var edges [][]uint16
var nTables, nRoots, nLibNames int

type entry struct {
	id   uint16
	bits uint8
}

const (
	hasRule    = 1 << iota // Rule_<name>
	hasField               // Field_<name>
	hasControl             // <Name>Control
	hasFormat              // <Name>Format
)

const alias = math.MaxInt16

var variants = []struct {
	bit            uint8
	prefix, suffix string
}{
	{hasRule, "Rule_", ""},
	{hasField, "Field_", ""},
	{hasControl, "", "Control"},
	{hasFormat, "", "Format"},
}

// buildGraph loads the dependency graph from suneido.db
func buildGraph() *db19.Database {
	db, err := db19.OpenDb("../../suneido.db", stor.Read, false)
	if err != nil {
		log.Fatal(err)
	}
	// first collect all the name nodes
	tableNames(db)
	nTables = int(nextId)
	rootNames(db)
	nRoots = int(nextId) - nTables
	libNames(db)
	nLibNames = int(nextId) - nTables - nRoots
	edges = make([][]uint16, nextId)
	setupVariants()

	// next collect all the reference edges
	stamps = make([]uint32, nextId)
	prev = make([]uint16, nextId)
	scanLibs(db)
	return db
}

func addNode(name string) {
	if _, ok := nodeMap[name]; ok {
		return
	}
	id := nextId
	nextId++
	nodeMap[name] = entry{id: id}
	nodeVec = append(nodeVec, name)
}

var excludeTables = map[string]struct{}{"Accountinglib": {}, "apmlib": {},
	"autoupdate_lib": {}, "autoupdate_test_betacadjavalib": {},
	"autoupdate_test_betacadlib": {}, "autoupdate_test_betausjavalib": {},
	"autoupdate_test_betauslib": {}, "autoupdate_test_curcadlib": {},
	"autoupdate_test_curuslib": {}, "autoupdate_test_lib": {},
	"axislib": {}, "axoncurrentcadlib": {}, "axoncurrentuslib": {},
	"axonlib": {}, "axonmedialib": {}, "axonservicelib": {},
	"betacadlib": {}, "betacadlinuxlib": {}, "betauslib": {},
	"betauslinuxlib": {}, "caronlib": {}, "carslib": {},
	"configlib": {}, "Contrib": {}, "demobookoptions": {},
	"demobookoptionsjslib": {}, "demobookoptionswebgui": {},
	"emailfwdbetalib": {}, "emailfwdlib": {}, "etalib": {},
	"etaprlib": {}, "ewlib": {}, "ftweblib": {}, "groupslib": {},
	"invenlib": {}, "joblib": {}, "Kenlib": {}, "Kenlibwebgui": {},
	"leadslib": {}, "mylib": {}, "pcmiler": {}, "polib": {},
	"prcadlib": {}, "prereleaselib": {}, "prlib": {}, "pruslib": {},
	"salesdemo": {}, "stdlib": {}, "sujslib": {}, "sujswebgui": {},
	"suneidoholdinglib": {}, "Test_lib": {}, "ticketlib": {},
	"towlib": {}, "vclib": {}, "weblib": {}, "wolib": {}}

// tableNames adds every table name in suneido.db to the node map,
// assigning each a consecutive id. Tables in excludeTables are skipped.
func tableNames(db *db19.Database) {
	for sc := range db.GetState().Meta.Tables() {
		if _, ok := excludeTables[sc.Table]; !ok {
			addNode(sc.Table)
		}
	}
}

// rootNames adds the text of each ETA entry that is a single capitalized
// identifier.
func rootNames(db *db19.Database) {
	rt := db.NewReadTran(core.AllPerms)
	defer rt.Abort()
	sc := db.GetState().Meta.GetRoSchema("ETA")
	textFld := slices.Index(sc.Columns, "text")
	scanTable(rt, "ETA", func(rec core.Record) {
		if name, ok := singleIdent(rec.GetStr(textFld)); ok {
			addNode(name)
		}
	})
}

// singleIdent returns text if it is a single capitalized identifier.
func singleIdent(text string) (string, bool) {
	lxr := lexer.NewLexer(text)
	item := lxr.NextSkip()
	if item.Token != tok.Identifier || !ascii.IsUpper(item.Text[0]) {
		return "", false
	}
	if lxr.NextSkip().Token != tok.Eof {
		return "", false
	}
	return item.Text, true
}

var libTables = []string{
	"carslib", "polib", "wolib", "invenlib", "etaprlib", "prcadlib", "prlib",
	"joblib", "ticketlib", "pcmiler", "towlib", "etalib", "axislib",
	"Accountinglib", "axonlib", "stdlib",
}

// libNames adds the name of each library item (group is -1) from the
// standard library tables.
func libNames(db *db19.Database) {
	rt := db.NewReadTran(core.AllPerms)
	defer rt.Abort()
	scanLibraries(rt, func(name, _ string) {
		addNode(name)
	})
}

// scanLibs sets edges[id] to the node ids referenced by each library item's text.
func scanLibs(db *db19.Database) {
	rt := db.NewReadTran(core.AllPerms)
	defer rt.Abort()
	scanLibraries(rt, func(name, text string) {
		if e, ok := nodeMap[name]; ok {
			addEdges(e.id, scanIdentifiers(name, text))
		}
	})
}

// addEdges adds ids to edges[id], skipping duplicates. A name can be defined
// in more than one library table, so the references are unioned rather than
// letting one definition overwrite another.
func addEdges(id uint16, ids []uint16) {
	for _, ref := range ids {
		if !slices.Contains(edges[id], ref) {
			edges[id] = append(edges[id], ref)
		}
	}
}

func setupVariants() {
	add := func(name string, bit uint8) {
		if e, ok := nodeMap[name]; ok {
			// if basename exists, set the bit
			e.bits |= bit
			nodeMap[name] = e
		} else {
			// if basename doesn't exist, create an "alias" entry for it
			nodeMap[name] = entry{id: alias, bits: bit}
		}
	}
	prefix := func(prefix, name string, bit uint8) {
		if strings.HasPrefix(name, prefix) &&
			len(name) > len(prefix) &&
			ascii.IsLower(name[len(prefix)]) {
			add(name[len(prefix):], bit)
		}
	}
	suffix := func(name, suffix string, bit uint8) {
		if strings.HasSuffix(name, suffix) &&
			len(name) > len(suffix) &&
			ascii.IsUpper(name[0]) {
			add(name[:len(name)-len(suffix)], bit)
		}
	}
	for i := nTables + nRoots; i < nTables+nRoots+nLibNames; i++ {
		name := nodeVec[i]
		prefix("Rule_", name, hasRule)
		prefix("Field_", name, hasField)
		suffix(name, "Control", hasControl)
		suffix(name, "Format", hasFormat)
	}
	// a table references its trigger, so the trigger is reachable from the table
	// BUT currently triggers run on the server with full permissions
	// for id := range nTables {
	// 	if e, ok := nodeMap["Trigger_"+nodeVec[id]]; ok && e.id != alias {
	// 		addEdges(uint16(id), []uint16{e.id})
	// 	}
	// }
}

// scanLibraries calls fn for each library definition (group is -1)
// in the standard library tables.
func scanLibraries(rt *db19.ReadTran, fn func(name, text string)) {
	for _, table := range libTables {
		if !rt.HasTable(table) {
			continue
		}
		sc := rt.GetSchema(table)
		nameFld := slices.Index(sc.Columns, "name")
		groupFld := slices.Index(sc.Columns, "group")
		textFld := slices.Index(sc.Columns, "text")
		scanTable(rt, table, func(rec core.Record) {
			if core.ToInt(rec.GetVal(groupFld)) == -1 {
				name := rec.GetStr(nameFld)
				if skipName(name) {
					return
				}
				fn(name, rec.GetStr(textFld))
			}
		})
	}
}

// skipName returns true for library items that are not Suneido source
// or that only contain tests.
func skipName(name string) bool {
	return strings.HasSuffix(name, ".js") ||
		strings.HasSuffix(name, ".css") ||
		strings.HasSuffix(name, "Test") ||
		strings.HasSuffix(name, "Tests") ||
		strings.Contains(name, "__webgui") ||
		strings.HasPrefix(name, "Update_20")
}

// scanTable calls fn for every record in table.
func scanTable(rt *db19.ReadTran, table string, fn func(core.Record)) {
	iter := rt.IndexIter(table, 0)
	for iter.Next(rt); !iter.Eof(); iter.Next(rt) {
		fn(iter.GetRecord(rt))
	}
}

//-------------------------------------------------------------------

var excludeNames = map[string]struct{}{"Params": {}, "Working": {},
	"Reporter": {}, "InfoControl": {}, "QueryFormat": {}}

// excluded returns true if name should not be treated as a reference.
func excluded(name string) bool {
	_, ok := excludeNames[name]
	return ok
}

var stamps []uint32
var stamp uint32

// scanIdentifiers returns the unique ids of the globals and table names in
// src that are in nodeMap. String constants may contain queries, so their
// identifiers are scanned for table names.
func scanIdentifiers(name, src string) []uint16 {
	stamp++
	return scanSrc(name, src, nil)
}

func scanSrc(name, src string, ids []uint16) []uint16 {
	walkRefs(parse(name, src), func(typ, value string) {
		switch typ {
		case "global":
			if excluded(value) {
				return
			}
			if e, ok := nodeMap[value]; ok && e.id != alias && markSeen(e.id) {
				ids = append(ids, e.id)
			}
		case "string":
			ids = scanString(value, ids)
		default:
			panic("unexpected type: " + typ)
		}
	})
	return ids
}

// scanString returns ids plus the ids of the table names and capitalized
// variants in src's identifiers.
func scanString(src string, ids []uint16) []uint16 {
	lxr := lexer.NewLexer(src)
	for {
		item := lxr.NextSkip()
		if item.Token == tok.Eof {
			return ids
		}
		if item.Token != tok.Identifier || excluded(item.Text) {
			continue
		}
		e, ok := nodeMap[item.Text]
		if !ok {
			continue
		}
		if int(e.id) < nTables && markSeen(e.id) {
			ids = append(ids, e.id)
		}
		if ascii.IsUpper(item.Text[0]) {
			for _, v := range variants {
				if e.bits&v.bit == 0 {
					continue
				}
				name := v.prefix + item.Text + v.suffix
				if excluded(name) {
					continue
				}
				if ve, ok := nodeMap[name]; ok &&
					ve.id != alias && markSeen(ve.id) {
					ids = append(ids, ve.id)
				}
			}
		}
	}
}

// markSeen returns true if id has not been seen with this stamp.
func markSeen(id uint16) bool {
	if stamps[id] == stamp {
		return false
	}
	stamps[id] = stamp
	return true
}
