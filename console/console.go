package console

import (
	"github.com/apmckinlay/gsuneido/options"
	"log"
	"os"
)

func LogFileAlso() {
	log.SetOutput(logWriter{})
}

type logWriter struct{}

var errCount = 0
const errLimit = 100
const tooMany = "too many errors, exiting"

// Write outputs to Stderr and error.log
func (lw logWriter) Write(p []byte) (n int, err error) {
	os.Stderr.Write(p)
	f, err := os.OpenFile(options.Errlog, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	if errCount++; errCount > errLimit {
		os.Stderr.WriteString(tooMany)
		f.WriteString(tooMany)
		f.Close()
		os.Exit(-1)
	}
	return f.Write(p)
}
