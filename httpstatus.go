// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"runtime/metrics"
	"sort"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/options"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var httpServer *http.Server

func startHttpStatus() {
	if httpServer != nil {
		return // already started
	}
	http.HandleFunc("/", httpStatus)
	http.HandleFunc("/metrics/", httpMetrics)
	http.HandleFunc("/info/", httpInfo)
	port := "3148"
	if options.WebPort != "" {
		port = options.WebPort
	} else if options.Port != "" {
		p, _ := strconv.Atoi(options.Port)
		port = strconv.Itoa(p + 1)
	}
	addr := ":" + port
	go func() {
		httpServer = &http.Server{Addr: addr}
		err := httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Println("ERROR: Server Monitor:", err)
		}
	}()
}

func httpStatus(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w,
		`<html>
		<head>
		<title>Suneido Monitor</title>
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
	s := extra + `
		<p>Built: ` + options.BuiltStr() + `</p>
		<p>Heap: ` + heap() + `</p>` +
		threads()
	if dbmsLocal != nil {
		s += `<p>Database: ` + mb(dbmsLocal.Size()) + `
		` + trans() + `
		` + dbms.Conns()
	}
	return s + `<p><a href="info/">Suneido Info</a> &nbsp;&nbsp;
			<a href="metrics/">Go metrics</a> &nbsp;&nbsp;
			<a href="debug/pprof/">Go pprof</a>	</p>`
}

func mb(n uint64) string {
	return strconv.FormatUint(((n+512*1024)/(1024*1024)), 10) + "mb"
}

func heap() string {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return mb(ms.HeapSys) + " (" + mb(ms.HeapInuse) + " in use)"
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
	for i := range n {
		sb.WriteString(sep)
		sb.WriteString(list.ListGet(i).String())
		sep = ", "
	}
	sb.WriteString("<p>\n")
	return sb.String()
}

func httpMetrics(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if req.URL.Path == "/metrics" || req.URL.Path == "/metrics/" {
		io.WriteString(w,
			`<html>
			<head><title>Go metrics</title></head>
			<body>
			<h1>Go metrics</h1>
			`)
		for _, d := range metrics.All() {
			if d.Kind == metrics.KindUint64 || d.Kind == metrics.KindFloat64 {
				fmt.Fprintf(w, `<a href="%s">%s</a><br />`+"\n",
					strings.TrimPrefix(d.Name, "/"), d.Name)
				fmt.Fprintln(w, "<p>"+d.Description+"</p>")
			}
		}
		io.WriteString(w, `</body></html>`)
	} else {
		sample := make([]metrics.Sample, 1)
		sample[0].Name = req.URL.Path[8:] // remove /metrics
		metrics.Read(sample)
		var x any
		switch sample[0].Value.Kind() {
		case metrics.KindUint64:
			x = printer.Sprintf("%d", sample[0].Value.Uint64())
		case metrics.KindFloat64:
			x = sample[0].Value.Float64()
		default:
			x = "unsupported"
		}
		fmt.Fprint(w, req.URL.Path, " = ", x)
	}
}

var printer = message.NewPrinter(language.English)

func httpInfo(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if req.URL.Path == "/info" || req.URL.Path == "/info/" {
		io.WriteString(w,
			`<html>
			<head><title>Suneido Info</title></head>
			<body>
			<h1>Suneido Info</h1>
			<pre>`)
		for _, name := range core.InfoList() {
			fmt.Fprintf(w, `<a href="%s">%s</a>`+"\n\n", name, name)
		}
		io.WriteString(w, `</pre></body></html>`)
	} else {
		s := core.InfoStr(req.URL.Path[6:]) // skip /info/
		fmt.Fprint(w, req.URL.Path, " = ", s)
	}
}
