package main

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/ptest"
)

func TestPtest(t *testing.T) {
	if !ptest.RunFile("execute.test") {
		t.Fail()
	}
}
