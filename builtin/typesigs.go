// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/typechecker"
	"github.com/apmckinlay/gsuneido/typechecker/annotations"
)

func init() {
	if err := typechecker.SetBuiltinSignatures(builtinRegistrations()); err != nil {
		Fatal("builtin type signatures:", err)
	}
}

var sigReceiver = map[string]string{
	"string": "string",
	"num":    "number",
	"int":    "number",
	"ob":     "object",
	"record": "record",
	"date":   "date",
	"class":  "class",
	"base":   "class",
	"func":   "function",
	"seq":    "sequence",
}

var sigGeneral = map[string]string{"record": "ob", "int": "num"}

var sigClass = map[string]string{
	"dateStatic":  "Date",
	"db":          "Database",
	"ftsearch":    "Ftsearch",
	"lruStatic":   "LruCache",
	"opgp":        "OpenPGP",
	"pcm":         "Pcmsrv64",
	"pe":          "PdfEncrypt",
	"rnd":         "Random",
	"sqs":         "Query",
	"suneido":     "Suneido",
	"thread":      "Thread",
	"typechecker": "TypeChecker",
	"zlib":        "Zlib",
}

func builtinRegistrations() []annotations.Registration {
	covered := map[string]map[string]bool{}
	for spec, gen := range sigGeneral {
		names := map[string]bool{}
		for _, ts := range builtinTypeSignatures {
			if ts.Kind == "method" && ts.Prefix == gen {
				names[ts.Name] = true
			}
		}
		covered[spec] = names
	}
	regs := make([]annotations.Registration, 0, len(builtinTypeSignatures))
	for _, ts := range builtinTypeSignatures {
		r := annotations.Registration{Kind: ts.Kind, Name: ts.Name, Sig: ts.Sig}
		switch ts.Kind {
		case "method":
			recv, ok := sigReceiver[ts.Prefix]
			if !ok || covered[ts.Prefix][ts.Name] {
				continue
			}
			r.Receiver = recv
		case "static":
			class, ok := sigClass[ts.Prefix]
			if !ok {
				Fatal("no class name for static method prefix:", ts.Prefix)
			}
			r.Class = class
		}
		regs = append(regs, r)
	}
	return regs
}
