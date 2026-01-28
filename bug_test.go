package main

import (
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
)

func TestBug(t *testing.T) {
	Libload = libload // dependency injection
	MainThread = &mainThread

	openDbms()
	defer db.CloseKeepMapped()

	src := `
		Init.Repl()
		retryException = 'some exception'
		count = 0
		block = { 
			try  
				{
		Print(count, "++")
				if count++ < 2
					throw 'case failed: testing'
		Print('Should see me!!')
				}
			catch (err, 'case failed:')
				{
				Suneido.X = err
				throw retryException
				}
			}
		Retry(block, maxRetries: 3, minDelayMs: 100, :retryException)
		`
	compile.EvalString(&mainThread, src)
}
