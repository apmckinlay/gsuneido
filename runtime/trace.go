// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"fmt"
	"os"
	"sync"

	"github.com/apmckinlay/gsuneido/aaainitfirst"
	"github.com/apmckinlay/gsuneido/options"
)

var traceLog *os.File
var traceLogOnce sync.Once
var traceConOnce sync.Once

func Trace(args ...interface{}) {
	s := fmt.Sprint(args...)
	if options.Trace&options.TraceLogFile != 0 {
		traceLogOnce.Do(func() {
			traceLog, _ = os.OpenFile("trace.log",
				os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		})
		if traceLog != nil {
			traceLog.WriteString(s)
		}
	}
	if options.Trace&options.TraceConsole != 0 {
		traceConOnce.Do(aaainitfirst.OutputToConsole)
		if aaainitfirst.ConsoleAttached() {
			os.Stdout.WriteString(s)
		}
	}
}
