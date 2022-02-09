// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package trace

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type what int

// var cur = ClientServer
var cur = what(0)

func Set(w int) int {
	prev := cur
	cur = what(w)
	return int(prev)
}

const (
	Functions what = 1 << iota
	Statements
	Opcodes
	Records
	Libraries
	SlowQuery
	Query
	Symbol
	AllIndex
	Table
	Select
	TempIndex
	QueryOpt

	Console
	LogFile

	ClientServer
	Exceptions
	Globals

	JoinOpt
	Dbms
)

func (w what) String() string {
	return map[what]string{
		Functions:    "FUNC ",
		Statements:   "STMT ",
		Opcodes:      "OP ",
		Records:      "REC ",
		Libraries:    "LIB ",
		SlowQuery:    "SLOWQUERY ",
		Query:        "QUERY ",
		Symbol:       "SYM ",
		AllIndex:     "ALLINDEX ",
		Table:        "TABLE ",
		Select:       "SELECT ",
		TempIndex:    "TEMPINDEX ",
		QueryOpt:     "QUERYOPT ",
		Console:      "CONSOLE ",
		LogFile:      "LOGFILE ",
		ClientServer: "CS ",
		Exceptions:   "EXCEPT ",
		Globals:      "GLOBAL ",
		JoinOpt:      "JOINOPT ",
		Dbms:         "DBMS ",
	}[w]
}

func (w what) Println(first interface{}, rest ...interface{}) {
	// kept short in hopes it will be inlined
	if cur&w != 0 {
		w.println(first, rest)
	}
}

func (w what) println(first interface{}, rest []interface{}) {
	s := w.String() + fmt.Sprint(first) + " " + fmt.Sprintln(rest...)
	Print(s)
}

func Print(s string) {
	if cur&LogFile != 0 || cur&(LogFile|Console) == 0 {
		logPrintln(s)
	}
	if cur&Console != 0 || cur&(LogFile|Console) == 0 {
		consolePrintln(s)
	}
}

func (w what) On() bool {
	return cur&w != 0
}

var traceLog *os.File
var traceLogOnce sync.Once

func logPrintln(s string) {
	traceLogOnce.Do(func() {
		var err error
		traceLog, err = os.OpenFile("trace.log",
			os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Println("ERROR", err)
		}
	})
	if traceLog != nil {
		traceLog.WriteString(s)
	}
}
