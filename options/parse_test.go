package options

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestParse(t *testing.T) {
	parse := func(args ...string) func(string) {
		Repl, Client, Port, Version, Help = false, "", "", false, false
		args = Parse(args)
		if len(args) == 0 {
			args = nil
		}
		s := ""
		if Repl {
			s += " repl"
		}
		if Client != "" {
			s += " " + Client
		}
		if Port != "" {
			s += " " + Port
		}
		if Version {
			s += " version"
		}
		if Help {
			s += " help"
		}
		for _,a := range args {
			s += " | " + a
		}
		if s != "" {
			s = s[1:]
		}
		return func(expected string) {
			t.Helper()
			Assert(t).That(s, Equals(expected))
		}
	}
	parse()("")
	parse("-r")("repl")
	parse("-repl")("repl")
	parse("-c")("127.0.0.1")
	parse("-client")("127.0.0.1")
	parse("-c", "--")("127.0.0.1")
	parse("-c", "1.2.3.4")("1.2.3.4")
	parse("-p", "1234")("1234")
	parse("-c", "-p", "1234")("127.0.0.1 1234")
	parse("-c", "localhost", "-p", "1234")("localhost 1234")
	parse("-c", "--", "foo", "bar")("127.0.0.1 | foo | bar")
}
