// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/options"
)

var httpServer *http.Server

func startHttpStatus() {
	http.HandleFunc("/", httpStatus)
	port, _ := strconv.Atoi(options.Port)
	addr := ":" + strconv.Itoa(port+1)
	go func() {
		httpServer = &http.Server{Addr: addr}
		err := httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Println("Server Monitor:", err)
		}
	}()
}
func httpStatus(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w,
		`<html>
		<head>
		<title>Suneido Server Monitor</title>
		<meta http-equiv="refresh" content="5" />
		</head>
		<body>
		<h1>Suneido Server Monitor</h1>
		`+body()+`
		</body>
		</html>`)
}

func body() string {
	extra := ""
	switch options.DbStatus.Load() {
	case "starting":
		extra = "<h2 style=\"color: blue;\">Starting ...</h2>"
	case "corrupted":
		extra = "<h2 style=\"color: red;\">Database damage detected - " +
			"operating in read-only mode</h2>"
	case "checking":
		return `<h2 style="color: red;">Checking database ...<h2>`
	case "repairing":
		return `<h2 style="color: red;">Repairing database ...<h2>`
	}
	return extra + `
		<p>Built: ` + builtin.Built() + `</p>
		<p>Heap: ` + mb(builtin.HeapSys()) + `</p>
		<p>Database: ` + mb(dbmsLocal.Size()) + `
		` + threads() + `
		` + trans() + `
		` + dbms.Conns()
}

func mb(n uint64) string {
	return strconv.FormatUint(((n+512*1024)/(1024*1024)), 10) + "mb"
}

func threads() string {
	list := builtin.ThreadList()
	sort.Strings(list)
	var sb strings.Builder
	fmt.Fprintf(&sb, "<p>Threads: (%d) ", len(list))
	sep := ""
	for _, s := range list {
		sb.WriteString(sep)
		sb.WriteString(s)
		sep = ", "
	}
	sb.WriteString("<p>\n")
	return sb.String()
}

func trans() string {
	list := dbmsLocal.Transactions()
	n := list.Size()
	var sb strings.Builder
	fmt.Fprintf(&sb, "<p>Transactions: (%d) ", n)
	sep := ""
	for i := 0; i < n; i++ {
		sb.WriteString(sep)
		sb.WriteString(list.ListGet(i).String())
		sep = ", "
	}
	sb.WriteString("<p>\n")
	return sb.String()
}
