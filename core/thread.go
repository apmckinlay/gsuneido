// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/generic/atomics"
	"github.com/apmckinlay/gsuneido/util/generic/cache"
	"github.com/apmckinlay/gsuneido/util/regex"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/tr"
)

// MainThread is injected by gsuneido.go to use for debugging
var MainThread *Thread

// See interp.go and args.go for the rest of the Thread methods

// maxStack is the size of the value stack, fixed size for performance
const maxStack = 1024

// maxFrames is the size of the frame stack, fixed size for performance
const maxFrames = 256

type Thread struct {

	// stack is the Value stack for arguments and expressions.
	// The end of the slice is the top of the stack.
	stack [maxStack]Value

	// Session is the name of the database session for clients and standalone.
	// Server tracks client session names separately.
	// Needs atomic because we access MainThread from other threads.
	session atomics.String

	// dbms is the database (client or local) for this Thread
	dbms IDbms

	// Suneido is a per-thread SuneidoObject that overrides the global one.
	// Needs atomic because sequence.go wrapIter may access from other threads.
	Suneido atomic.Pointer[SuneidoObject]

	// subThreadOf is used for sessions
	subThreadOf *Thread

	// blockReturnFrame is the parent frame of the block that is returning
	blockReturnFrame *Frame

	// RxCache is per thread so no locking is required
	rxCache *cache.Cache[string, regex.Pattern]
	// TrCache is per thread so no locking is required
	trCache *cache.Cache[string, tr.Set]

	// Name is the name of the thread (default is Thread-#)
	Name string

	Nonce string

	profile profile

	// rules is a stack of the currently running rules, used by SuRecord
	rules activeRules

	// frames are the Frame's making up the call stack.
	frames [maxFrames]Frame

	// Quote is used by Display to request specific quotes
	Quote int

	// OpCount counts op codes in interp, for polling
	OpCount int

	// sp is the stack pointer, top is stack[sp-1]
	sp int

	// spMax is the "high water" mark for sp
	spMax int

	// fp is the frame pointer, top is frames[fp-1]
	fp int

	// fpMax is the "high water" mark for fp
	fpMax int

	// Num is a unique number assigned to the thread
	Num int32

	// UIThread is only set for the main UI thread.
	// It controls whether interp checks for UI requests from other threads.
	UIThread bool

	// InHandler is used to detect nested handler calls
	InHandler bool

	// ReturnThrow is set by op.ReturnThrow and used by op.Call*Discard
	// and by some built-in functions.
	ReturnThrow bool

	// Sviews are the session view definitions for this thread
	sv *Sviews

	Rand *rand.Rand
}

var threadNum atomic.Int32

// NewThread creates a new thread.
// It is primarily used for user initiated threads.
// Internal threads can just use a zero Thread.
func NewThread(parent *Thread) *Thread {
	th := setup(&Thread{})
	if parent != nil {
		if suneido := parent.Suneido.Load(); suneido != nil {
			suneido.SetConcurrent()
			th.Suneido.Store(suneido)
		}
		th.sv = parent.sv
	}
	return th
}

func setup(th *Thread) *Thread {
	th.Num = threadNum.Add(1)
	th.Name = "Thread-" + strconv.Itoa(int(th.Num))
	mts := ""
	if MainThread != nil {
		mts = MainThread.Session()
	}
	th.session.Store(str.Opt(mts, ":") + th.Name)
	return th
}

// Reset is equivalent to calling NewThread(nil)
func (th *Thread) Reset() {
	*th = Thread{} // zero it
	setup(th)
}

func (th *Thread) Session() string {
	if th.subThreadOf != nil {
		th = th.subThreadOf
	}
	return th.session.Load()
}

func (th *Thread) SetSession(s string) {
	if th.subThreadOf != nil {
		th = th.subThreadOf
	}
	th.session.Store(s)
}

func (th *Thread) SetSviews(sv *Sviews) {
	th.sv = sv
}

// Push pushes a value onto the value stack
func (th *Thread) Push(x Value) {
	if th.sp >= maxStack {
		panic("value stack overflow")
	}
	th.stack[th.sp] = x
	th.sp++
}

// Pop pops a value off the value stack
func (th *Thread) Pop() Value {
	th.sp--
	return th.stack[th.sp]
}

// Top returns the top of the value stack (without modifying the stack)
func (th *Thread) Top() Value {
	return th.stack[th.sp-1]
}

// Swap exchanges the top two values on the stack
func (th *Thread) Swap() {
	th.stack[th.sp-1], th.stack[th.sp-2] = th.stack[th.sp-2], th.stack[th.sp-1]
}

// GetState and RestoreState are used by callbacks_windows.go

type state struct {
	fp int
	sp int
}

func (th *Thread) GetState() state {
	return state{fp: th.fp, sp: th.sp}
}

func (th *Thread) RestoreState(st any) {
	s := st.(state)
	th.fp = s.fp
	th.sp = s.sp
}

// Callstack captures the call stack
func (th *Thread) Callstack() *SuObject {
	// NOTE: it might be more efficient
	// to capture the call stack in an internal format
	// and only build the SuObject if required
	cs := &SuObject{}
	for i := th.fp - 1; i >= 0; i-- {
		fr := th.frames[i]
		call := &SuObject{}
		call.Set(SuStr("fn"), fr.fn)
		call.Set(SuStr("srcpos"), IntVal(fr.fn.CodeToSrcPos(fr.ip-1)))
		call.Set(SuStr("locals"), th.locals(i))
		cs.Add(call)
		if cs.Size() > 50 {
			break
		}
	}
	return cs
}

func (th *Thread) Locals(i int) *SuObject {
	return th.locals(th.fp - 1 - i)
}

func (th *Thread) locals(i int) *SuObject {
	fr := th.frames[i]
	locals := &SuObject{}
	if fr.this != nil {
		locals.Set(SuStr("this"), fr.this)
	}
	for i, v := range fr.locals.v {
		if v != nil && fr.fn != nil && i < len(fr.fn.Names) {
			locals.Set(SuStr(fr.fn.Names[i]), v)
		}
	}
	return locals
}

// PrintStack outputs the thread's call stack to stderr
func (th *Thread) PrintStack() {
	th.printStack(os.Stderr, 20)
}

// TraceStack outputs the thread's call stack to trace
func (th *Thread) TraceStack() {
	th.printStack(trace.Writer, 6)
}

func (th *Thread) printStack(w io.Writer, levels int) {
	limit := max(th.fp-levels, 0)
	for i := th.fp - 1; i >= limit; i-- {
		frame := th.frames[i]
		fmt.Fprintln(w, frame.fn)
	}
}

func PrintStack(cs *SuObject) {
	if cs == nil {
		return
	}
	for i := 0; i < cs.ListSize(); i++ {
		frame := cs.ListGet(i)
		fn := frame.Get(nil, SuStr("fn"))
		fmt.Fprintln(os.Stderr, fn)
	}
}

func (th *Thread) TraceCaller() {
	if i := th.fp - 1; i >= 0 {
		trace.Println(th.frames[i].fn)
	}
}

// SetDbms is used to set up the main thread initially
func (th *Thread) SetDbms(dbms IDbms) {
	th.dbms = dbms
}

// GetDbms requires dependency injection
var GetDbms func() IDbms

var DbmsAuth = false

func (th *Thread) Dbms() IDbms {
	if th.dbms == nil {
		th.dbms = GetDbms()
		if s := th.session.Load(); s != "" {
			// session id was set before connecting
			th.dbms.SessionId(th, s)
		}
	}
	return th.dbms.Unwrap()
}

// Close closes the thread's dbms connection (if it has one)
func (th *Thread) Close() {
	if th.dbms != nil && options.Action == "client" {
		th.dbms.Close()
		th.dbms = nil
	}
}

// SubThread is a NewThread with the same dbms as this thread.
// This is used for the RunOnGoSide and SuneidoAPP threads.
// We want a new thread for isolation e.g. for exceptions or dynamic variables
// but we don't need the overhead of another dbms connection.
// WARNING: This should only be used where it is guaranteed
// that the Threads will NOT be used concurrently.
func (th *Thread) SubThread() *Thread {
	t2 := NewThread(th)
	t2.dbms = th.dbms
	t2.subThreadOf = th
	return t2
}

func (th *Thread) Cat(x, y Value) Value {
	return OpCat(th, x, y)
}

func (th *Thread) SessionId(id string) string {
	if th.dbms == nil {
		// don't create a connection just to get/set the session id
		if id != "" {
			th.SetSession(id)
		}
		return th.Session()
	}
	return th.dbms.SessionId(th, id)
}

func (th *Thread) RunWithMainSuneido(fn func() Value) Value {
	defer func(orig *SuneidoObject) {
		th.Suneido.Store(orig)
	}(th.Suneido.Load())
	th.Suneido.Store(nil)
	return fn()
}

func (th *Thread) Regex(x Value) regex.Pattern {
	if sr, ok := x.(SuRegex); ok {
		return sr.Pat
	}
	if th.rxCache == nil {
		th.rxCache = cache.New(regex.Compile)
	}
	return th.rxCache.Get(ToStr(x))
}

func (th *Thread) TrSet(x Value) tr.Set {
	if th.trCache == nil {
		th.trCache = cache.New(tr.New)
	}
	return th.trCache.Get(ToStr(x))
}

func (th *Thread) Sviews() *Sviews {
	if th == nil {
		return nil
	}
	return th.sv
}

//-------------------------------------------------------------------

var tsCount int
var tsLimit int
var tsLast SuDate
var tsLock sync.Mutex

// Timestamp is a Thread method for convenience, so it has access to the dbms.
// It is not "per thread".
// This is the "client" version of Timestamp.
func (th *Thread) Timestamp() PackableValue {
	tsLock.Lock()
	defer tsLock.Unlock()
	if tsCount++; tsCount < tsLimit {
		// fast path
		if tsLimit == TsInitialBatch {
			tsLast = tsLast.AddMs(1)
			return tsLast
		}
		return SuTimestamp{SuDate: tsLast, extra: uint8(tsCount)}
	}
	// fetch a new timestamp, slow path
	if tsLimit == 0 {
		go tsExpire()
	}
	tsLast = th.Dbms().Timestamp()
	tsCount = 0
	if tsLast.Millisecond() < TsThreshold {
		tsLimit = TsInitialBatch
	} else {
		tsLimit = 256
	}
	return tsLast
}

func tsExpire() {
	for {
		time.Sleep(1 * time.Second)
		// clear tsLast to force fetching a new timestamp
		tsLock.Lock()
		tsCount = tsLimit + 1
		tsLock.Unlock()
	}
}
