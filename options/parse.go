// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package options

import (
	"math/bits"
	"strconv"
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
		case match(&args, "-help"), match(&args, "-h"), match(&args, "-?"):
			setAction("help")
		case match(&args, "-check"):
			setAction("check")
		case match(&args, "-fullcheck"):
			setAction("check")
			FullCheck = true
		case match(&args, "-client"), match(&args, "-c"):
			setAction("client")
			Arg = "127.0.0.1"
			args = optionalArg(args, &Arg)
		case match(&args, "-compact"):
			setAction("compact")
		case match(&args, "-dump"), match(&args, "-d"):
			setAction("dump")
			args = optionalArg(args, &Arg)
		case match(&args, "-load"), match(&args, "-l"):
			setAction("load")
			args = optionalArg(args, &Arg)
		case match(&args, "-passphrase"), match(&args, "-pp"):
			if Action != "load" {
				error("passphrase only valid with -load")
			}
			args = optionalArg(args, &Passphrase)
			if Passphrase == "" {
				error("passphrase required")
			}
		case match(&args, "-port"), match(&args, "-p"):
			args = optionalArg(args, &Port)
			if Port == "" {
				error("port number required")
			} else if _, ok := atoui(Port); !ok {
				error("invalid port number")
			}
		case match(&args, "-repair"):
			setAction("repair")
		case match(&args, "-server"), match(&args, "-s"):
			setAction("server")
		case match(&args, "-unattended"), match(&args, "-u"):
			//TEMP for backward compatibility
		case match(&args, "-version"), match(&args, "-v"):
			Action = "version"
		case match(&args, "-ignoreversion"), match(&args, "-iv"):
			IgnoreVersion = true
		case match(&args, "-timeout"), match(&args, "-to"):
			to := ""
			args = optionalArg(args, &to)
			if to == "" {
				error("timeout value required")
			} else if n, ok := atoui(to); ok {
				TimeoutMinutes = n
			} else {
				error("invalid timeout value")
			}
		case match(&args, "-web"), match(&args, "-w"):
			WebServer = true
			args = optEqualArg(args, &WebPort)
			if WebPort != "" {
				if _, ok := atoui(WebPort); !ok {
					error("invalid web port number")
				}
			}
		case match(&args, "-printstates"):
			setAction("printstates")
		case match(&args, "-checkstates"):
			setAction("checkstates")
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
		error("port should only be specified with -server or -client, not " +
			Action)
	}
	if Port == "" && (Action == "client" || Action == "server") {
		Port = "3147"
	}
	CmdLine = remainder(args)
}

func atoui(s string) (int, bool) {
	n, err := strconv.ParseUint(s, 10, bits.UintSize)
	return int(n), err == nil
}

func match(pargs *[]string, s string) bool {
	arg := (*pargs)[0]
	if arg == s {
		*pargs = (*pargs)[1:]
		return true
	}
	if strings.HasPrefix(arg, s+"=") {
		(*pargs)[0] = strings.TrimPrefix(arg, s) // leave the '='
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

func optionalArg(args []string, parg *string) []string {
	if len(args) > 0 {
		if strings.HasPrefix(args[0], "=") {
			*parg = args[0][1:]
			args = args[1:]
		} else if !strings.HasPrefix(args[0], "-") { //DEPRECATED
			*parg = args[0]
			args = args[1:]
		}
	}
	return args
}

func optEqualArg(args []string, parg *string) []string {
	if len(args) > 0 && strings.HasPrefix(args[0], "=") {
		*parg = args[0][1:]
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
