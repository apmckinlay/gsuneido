package ptest

import "testing"

func TestPtest(t *testing.T) {
	if !RunFile("ptest.test") {
		t.Fail()
	}
}

func init() {
	Add("ptest", ptest)
}

func ptest(args []string, _ []bool) bool {
	return args[0] == args[1]
}
