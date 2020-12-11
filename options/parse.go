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
		if args[0] == "" {
			args = args[1:]
			continue
		}
		switch {
		case match(&args, "-client"), match(&args, "-c"):
			setAction("client")
			Arg = "127.0.0.1"
			args = optionalArg(args)
		case match(&args, "-check"):
			setAction("check")
		case match(&args, "-repair"):
			setAction("repair")
		case match(&args, "-compact"):
			setAction("compact")
		case match(&args, "-dump"), match(&args, "-d"):
			setAction("dump")
			args = optionalArg(args)
		case match(&args, "-load"), match(&args, "-l"):
			setAction("load")
			args = optionalArg(args)
		case match(&args, "-repl"), match(&args, "-r"):
			if Mode == "gui" {
				error("-repl not support for gui mode")
			} else {
				setAction("repl")
			}
		case match(&args, "-port"), match(&args, "-p"):
			if len(args) > 0 && args[0][0] != '-' {
				Port = args[0]
				args = args[1:]
			} else {
				error("port number required")
			}
		case match(&args, "-server"), match(&args, "-s"):
			setAction("server")
		case match(&args, "-unattended"), match(&args, "-u"):
			Unattended = true
		case match(&args, "-version"), match(&args, "-v"):
			Action = "version"
		case match(&args, "-ignoreversion"), match(&args, "-iv"):
			// for compatibility with cSuneido
		case match(&args, "-norelaunch"), match(&args, "-nr"):
			NoRelaunch = true
		case match(&args, "--"):
			break loop
		default:
			error("invalid command line argument: " + args[0])
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
	}
}

func match(pargs *[]string, s string) bool {
	arg := (*pargs)[0]
	if arg == s {
		*pargs = (*pargs)[1:]
		return true
	} else if strings.HasPrefix(arg, s) {
		(*pargs)[0] = strings.TrimSpace(arg[len(s):])
		return true
	}
	return false
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
