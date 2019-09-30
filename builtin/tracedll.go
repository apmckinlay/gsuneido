package builtin

// import (
// 	"fmt"
// 	"sync"

// 	"golang.org/x/sys/windows"
// )

// type TraceDLL struct {
// 	name string
// 	dll  *windows.LazyDLL
// }

// type TraceProc struct {
// 	name  string
// 	proc  *windows.LazyProc
// 	once  sync.Once
// 	calls bool
// }

// func NewTraceDLL(name string) TraceDLL {
// 	return TraceDLL{name: name, dll: windows.NewLazyDLL(name)}
// }

// func (td TraceDLL) MustFindProc(name string) *TraceProc {
// 	return &TraceProc{name: td.name + ":" + name, proc: td.dll.NewProc(name)}
// }

// func (td TraceDLL) NewTraceProc(name string) *TraceProc {
// 	return &TraceProc{name: td.name + ":" + name, proc: td.dll.NewProc(name), calls: true}
// }

// // WARNING: may not be safe because uintptr may be garbage collected
// func (tp *TraceProc) Call(args ...uintptr) (uintptr, uintptr, error) {
// 	tp.once.Do(func() { fmt.Println(tp.name) })
// 	if tp.calls {
// 		fmt.Printf("%s %#v\n", tp.name, args)
// 	}
// 	r1, r2, err := tp.proc.Call(args...)
// 	if tp.calls {
// 		fmt.Println(tp.name, "=>", r1)
// 	}
// 	return r1, r2, err
// }

// func (tp *TraceProc) Addr() uintptr {
// 	return tp.proc
// }
