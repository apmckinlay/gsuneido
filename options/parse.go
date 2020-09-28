// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package options

import (
	"os"
	"strings"
)

// Parse processes the command line options
// returning the remaining arguments
func Parse(args []string) {
loop:
	for len(args) > 0 && (args[0] == "" || args[0][0] == '-') {
		arg := args[0]
		args = args[1:]
		switch arg {
		case "":
			// ignore
		case "-c", "-client":
			setAction("client")
			Arg = "127.0.0.1"
			args = optionalArg(args)
		case "-check":
			setAction("check")
		case "-repair":
			setAction("repair")
		case "-compact":
			setAction("compact")
		case "-d", "-dump":
			setAction("dump")
			args = optionalArg(args)
		case "-l", "-load":
			setAction("load")
			args = optionalArg(args)
		case "-r", "-repl":
			setAction("repl")
		case "-p", "-port":
			if len(args) > 0 && args[0][0] != '-' {
				Port = args[0]
				args = args[1:]
			} else {
				error(arg + " must be followed by port number")
			}
		case "-s", "-server":
			setAction("server")
		case "-u", "-unattended":
			Unattended = true
		case "-v", "-version":
			Action = "version"
		case "--":
			break loop
		default:
			Action = "help"
		}
		if Action == "error" {
			return
		}
	}
	if Port != "" && Action != "client" && Action != "server" {
		error("port should only be specifed with -server or -client, not " +
			Action)
	}
	if Port == "" && (Action == "client" || Action == "server") {
		Port = "3147"
	}
	CmdLine = remainder(args)
	if Action == "client" {
		temp := os.TempDir() + "/"
		Errlog = temp + "suneido" + Port + ".err"
		Outlog = temp + "suneido" + Port + ".out"
	}
}

func setAction(action string) {
	if Action == "" {
		Action = action
	} else {
		error("only one action is allowed, can't have both " + Action +
			" and " + action)
	}
}

func optionalArg(args []string) []string {
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		Arg = args[0]
		args = args[1:]
	}
	return args
}

func error(err string) {
	Action = "error"
	Error = err
}

func remainder(args []string) string {
	var sb strings.Builder
	sep := ""
	for _, arg := range args {
		sb.WriteString(sep)
		sep = " "
		sb.WriteString(EscapeArg(arg))
	}
	return sb.String()
}
