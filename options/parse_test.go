// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package options

import (
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestParse(t *testing.T) {
	test := func(argstr, expected string) {
		t.Helper()
		args := strings.Fields(argstr)
		Action, Arg, Port, CmdLine, Error = "", "", "", "", ""
		TimeoutMinutes = 0
		WebServer, WebPort = false, ""
		Parse(args)
		s := Action
		if Arg != "" {
			s += " " + Arg
		}
		if Port != "3147" && Port != "" {
			s += " port " + Port
		}
		if TimeoutMinutes != 0 {
			s += " timeout=" + strconv.Itoa(TimeoutMinutes)
		}
		if WebServer {
			s += " web"
			if WebPort != "" {
				s += "=" + WebPort
			}
		}
		if CmdLine != "" {
			s += " | " + CmdLine
		}
		s = strings.TrimPrefix(s, " ")
		if Action == "error" {
			s = "error " + Error
		}
		if strings.HasPrefix(expected, "error ") {
			assert.T(t).Msg(s).
				That(strings.HasPrefix(s, expected))
		} else {
			assert.T(t).This(s).Is(expected)
		}
	}
	test("", "")
	test("-c", "client 127.0.0.1")
	test("-client", "client 127.0.0.1")
	test("-c --", "client 127.0.0.1")
	test("-c 1.2.3.4", "client 1.2.3.4")
	test("-c=1.2.3.4", "client 1.2.3.4")
	test("-client 1.2.3.4", "client 1.2.3.4")
	test("-c -p 1234", "client 127.0.0.1 port 1234")
	test("-c -p=1234", "client 127.0.0.1 port 1234")
	test("-s -port=1234", "server port 1234")
	test("-c localhost -p 1234", "client localhost port 1234")
	test("-c1.2.3.4", "error invalid command line argument")
	test("-c -- foo bar", "client 127.0.0.1 | foo bar")
	test("-client -- foo bar", "client 127.0.0.1 | foo bar")
	test("-client=1.2.3.4 foo bar", "client 1.2.3.4 | foo bar")
	test("-client1.2.3.4", "error invalid command line argument")
	test("-c 1.2.3.4 foo bar", "client 1.2.3.4 | foo bar")
	test("-client 1.2.3.4 foo bar", "client 1.2.3.4 | foo bar")

	test("-help", "help")
	test("-h", "help")
	test("-?", "help")

	test("-check", "check")
	test("-compact", "compact")

	test("-load -client", "error only one action is allowed")
	test("-load", "load")
	test("-load stdlib", "load stdlib")

	test("-p", "error port number required")
	test("-port", "error port number required")
	test("-c -p1234", "error invalid command line argument")
	test("-c -port1234", "error invalid command line argument")
	test("-s -port=1.2.3.4", "error invalid port number")
	test("-check -port=1234",
		"error port should only be specified with -server or -client")

	test("-dump", "dump")
	test("-dump stdlib", "dump stdlib")

	test("-server", "server")
	test("-repair", "repair")

	test("-to=44", "timeout=44")
	test("-to", "error timeout value required")
	test("-to=1.2", "error invalid timeout value")

	test("-v", "version")
	test("-version", "version")

	test("-w", "web")
	test("-web", "web")
	test("-w=1234", "web=1234")
	test("-web=1234", "web=1234")
	test("-w foo", "web | foo")
	test("-web=1.2.3.4", "error invalid web port number")

	test("-xyz", "error invalid command line argument")

}

func TestEscapeArg(t *testing.T) {
	test := func(s, expected string) {
		t.Helper()
		assert.T(t).This(EscapeArg(s)).Is(expected)
	}
	test(`foo`, `foo`)
	test(`foo bar`, `"foo bar"`)
	test(`ab"c`, `ab\"c`)
	test(`\`, `\`)
	test(`a\\\b`, `a\\\b`)
	test(`a\"b`, `a\\\"b`)
}
