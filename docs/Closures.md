# Closures in gSuneido

A closure in gSuneido is a block (anonymous function) that captures variables from its enclosing scope. This document explains how closures are implemented across the codebase.

A "block" is a syntactic construct `{|x,y| x + y }`
- it returns the value of its last statement
- `return` from a block returns from the enclosing function (via a panic)
- `break` or `continue` from a block panic "block:break" or "block:continue"

A block may be compiled several different ways:
- as a standalone function
- as a closure to implement block return
- as a closure to handle shared (captured) variables

## Overview

Closures in gSuneido work by capturing the local variables of the enclosing function at the point the closure is created. When the closure is later called, it has access to these captured variables.

## Compilation (`compile/codegen.go`)

### Determining Closure vs Function

Before code generation, the [`Blocks`](compile/ast/blocks.go:38) function in the AST package determines whether each block should be compiled as a standalone function or as a closure. A block becomes a closure if it shares any variables with its enclosing scope.

```go
func Blocks(f *Function) {
    // first traverse the ast and collect outer variables
    // and a list of blocks, their params & variables, and their parent if nested.
    var b bloks
    vars := make(strset)
    b.params(f.Params, vars)
    for _, stmt := range f.Body {
        b.statement(stmt, vars)
    }
    // then check for variable sharing
    for _, x := range b.bloks {
        x.block.CompileAsFunction = true
    }
    for i, x := range b.bloks {
        _, this := x.vars["this"]
        _, super := x.vars["super"]
        if this || super || x.hasRet || shares(x.vars, vars) ||
            (x.parent != nil && shares(x.vars, x.parent.params)) {
            closure(x)
        }
        for j := i + 1; j < len(b.bloks); j++ {
            y := b.bloks[j]
            if shares(x.vars, y.vars) {
                closure(x)
                closure(y)
            }
        }
    }
}
```

A block is compiled as a closure (not a function) if:
- It uses `this` or `super` (needs access to enclosing object)
- It contains a `return` statement (needs block return semantics to return from enclosing function)
- It shares variables with the outer function scope
- It shares variables with its parent block's parameters
- It shares variables with a sibling block

The `closure` function marks the block and all its parent blocks as closures:

```go
func closure(x *blok) {
    for x != nil {
        x.block.CompileAsFunction = false
        x = x.parent
    }
}
```

### Code Generation for Closures

The [`codegenClosureBlock`](compile/codegen.go:121) function generates bytecode for closure blocks:

```go
func codegenClosureBlock(ast *ast.Function, outercg *cgen) (*SuFunc, []string) {
    base := len(outercg.Names)
    cg := cgen{outerFn: outercg.outerFn, base: outercg.base, isBlock: true,
        cover: outercg.cover}
    cg.Names = outercg.Names
    cg.Lib = outercg.Lib
    cg.Name = outercg.Name

    f := cg.codegen(ast)

    // hide parameters from outer function
    outerNames := f.Names
    f.Names = make([]string, len(outerNames))
    f.Offset = uint8(base)
    copy(f.Names, outerNames)
    for i := range int(f.Nparams) {
        outerNames[base+i] += "|" + strconv.Itoa(base+i)
    }
    return f, outerNames
}
```

Key aspects:

- The closure shares the outer function's `Names` table, allowing access to outer local variables
- The `Offset` field marks where the closure's own parameters begin in the combined locals array
- Parameters from the outer function are "hidden" by adding a suffix (e.g., `x|0`)

### Block Compilation

In the [`block`](compile/codegen.go:1079) method, closures are compiled with the `op.Closure` opcode:

```go
func (cg *cgen) block(b *ast.Block) {
    f := &b.Function
    var fn *SuFunc
    if b.CompileAsFunction {
        fn = codegen2(cg.Lib, b.Name, f, cg.outerFn, cg.prevDef)
        cg.emitValue(fn)
    } else {
        // closure
        fn, cg.Names = codegenClosureBlock(f, cg)
        i := cg.value(fn)
        cg.emitUint8(op.Closure, i)
    }
    fn.IsBlock = true
}
```

## Runtime Representation (`core/suclosure.go`)

### SuClosure Structure

The [`SuClosure`](core/suclosure.go:14) struct represents a closure instance:

```go
type SuClosure struct {
    this Value
    // parent is the Frame of the outer function that created this closure.
    // It is used by interp to handle block returns.
    parent *Frame
    locals []Value // if concurrent, then read-only
    *SuFunc
    concurrent bool
}
```

Key fields:

- `this`: The `this` value from the enclosing scope (for method closures)
- `parent`: Reference to the Frame where the closure was created (for block returns)
- `locals`: The captured local variables from the enclosing scope
- `SuFunc`: The compiled function/block being executed
- `concurrent`: Whether the closure is being used in a concurrent context

### Calling a Closure

The [`Call`](core/suclosure.go:36) method sets up a new frame to execute the closure:

```go
func (b *SuClosure) Call(th *Thread, this Value, as *ArgSpec) Value {
    bf := b.SuFunc

    v := b.locals
    if b.concurrent {
        // make a mutable copy of the locals for the frame
        v = slc.Clone(b.locals)
    }

    // normally done by SuFunc Call
    args := th.Args(&b.ParamSpec, as)

    // copy args
    for i := range int(b.Nparams) {
        v[int(bf.Offset)+i] = args[i]
    }

    if this == nil {
        this = b.this
    }
    if th.fp >= len(th.frames) {
        panic("function call overflow")
    }
    fr := &th.frames[th.fp]
    fr.fn = bf
    fr.this = this
    fr.blockParent = b.parent
    fr.locals = locals{v: v, onHeap: true}
    return th.run()
}
```

The closure call:

1. Uses the captured locals as the base for the new frame's locals
2. For concurrent closures, clones the locals to allow independent mutation
3. Copies the closure's arguments into the appropriate offset position
4. Sets `blockParent` to enable block return semantics
5. Marks locals as `onHeap: true` since they're now heap-allocated

## Interpretation (`core/interp.go`)

### The Closure Opcode

The [`op.Closure`](core/interp.go:609) opcode creates a closure at runtime:

```go
case op.Closure:
    fr.locals.moveToHeap()
    fn := fr.fn.Values[fetchUint8()].(*SuFunc)
    parent := fr
    if fr.blockParent != nil {
        parent = fr.blockParent
    }
    block := &SuClosure{SuFunc: fn, locals: fr.locals.v, this: fr.this,
        parent: parent}
    th.Push(block)
```

The closure creation process:

1. **Move locals to heap**: Calls [`moveToHeap`](core/frame.go:38) to copy stack-based locals to the heap
2. **Capture parent frame**: Uses the current frame (or its blockParent) for block returns
3. **Create closure object**: Packages the function, captured locals, and `this` value

### Moving Locals to Heap

The [`moveToHeap`](core/frame.go:38) method in the `locals` struct:

```go
func (ls *locals) moveToHeap() {
    if ls.onHeap {
        return
    }
    // not concurrent at this point
    oldlocals := ls.v
    ls.v = slc.Clone(oldlocals)
    ls.onHeap = true
}
```

This ensures captured variables outlive the function that created them.

## Frame Structure (`core/frame.go`)

### Frame Definition

The [`Frame`](core/frame.go:9) struct holds execution context:

```go
type Frame struct {
    this Value           // instance if running a method
    fn *SuFunc           // the Function being executed
    blockParent *Frame   // used for block returns
    locals locals        // local variables (on stack or heap)
    ip int               // instruction pointer
    catchJump int
    catchSp int
}
```

### Locals Structure

```go
type locals struct {
    v []Value
    // onHeap is true when locals have been moved from the stack to the heap
    onHeap bool
}
```

The `onHeap` flag tracks whether locals are stack-allocated (normal function or heap-allocated (closure).

## Thread Context (`core/thread.go`)

The [`Thread`](core/thread.go:39) struct manages execution:

```go
type Thread struct {
    thread1
    thread2
}

type thread1 struct {
    stack [maxStack]Value    // value stack for arguments and expressions
    frames [maxFrames]Frame  // the call stack
    sp int                   // stack pointer
    fp int                   // frame pointer
    // ...
}
```

The thread maintains:

- A value stack for expressions and arguments
- A frame stack for function/block calls

## Block Returns

Closures support "block returns" - returning from the enclosing function from within a block. This is handled via:

1. **Panic mechanism**: [`BlockReturn`](core/interp.go:19) is a special panic used to transfer control
2. **blockParent chain**: Each closure captures its parent frame
3. **BlockReturn opcode**: In [`op.BlockReturn`](core/interp.go:626):

```go
case op.BlockReturn:
    th.blockReturnFrame = fr.blockParent
    panic(BlockReturn)
```

The panic is caught by the [`run`](core/interp.go:45) method which propagates it to the appropriate frame.

## Break and Continue from Blocks

Blocks also support `break` and `continue` statements. This allows loops to be terminated or continued from within a block passed to a higher-order function like `Each`.

### Implementation

The [`BlockBreak`](core/interp.go:17) and [`BlockContinue`](core/interp.go:18) are special exceptions similar to `BlockReturn`:

```go
var BlockBreak = BuiltinSuExcept("block:break")
var BlockContinue = BuiltinSuExcept("block:continue")
```

These are emitted by the compiler when `break` or `continue` is used within a block context:

```go
// from compile/codegen.go
func (cg *cgen) breakStmt(labels *Labels) {
    if labels != nil {
        labels.brk = cg.emitJump(op.Jump, labels.brk)
    } else if cg.isBlock {
        cg.emit(op.BlockBreak)  // throw "block:break"
    } else {
        panic("break can only be used within a loop")
    }
}

func (cg *cgen) continueStmt(labels *Labels) {
    if labels != nil {
        // ...
    } else if cg.isBlock {
        cg.emit(op.BlockContinue)  // throw "block:continue"
    } else {
        panic("continue can only be used within a loop")
    }
}
```

In the interpreter, these opcodes simply panic with the special exception:

```go
case op.BlockBreak:
    panic(BlockBreak)
case op.BlockContinue:
    panic(BlockContinue)
```

### Handling in Application Code

The calling code catches these exceptions and handles them appropriately. For example, in [`Objects.ss`](stdlib/Data%20Types/Objects.ss:449):

```suneido
catch (e, "block:")
    if (e is "block:break")
        break
    // else block:continue ... so continue
```

This pattern is used throughout the standard library to allow `break` and `continue` to work with iterator methods like `Each`.

### Example

```suneido
// Break from within an Each block
Items.Each()
    {
    if (@value is target)
        break  // exits the Each loop
    }
```

```suneido
// Continue to next iteration
Numbers.Each()
    {
    if (@value is odd)
        continue  // skip odd numbers
    sum = sum + @value
    }
```

The `break` or `continue` within the block throws `"block:break"` or `"block:continue"`, which is caught by the `Each` implementation and converted to the appropriate loop control flow.

## Concurrency Support

Closures can be used in concurrent contexts. The [`SetConcurrent`](core/suclosure.go:71) method prepares a closure for concurrent use:

```go
func (b *SuClosure) SetConcurrent() {
    if b.concurrent {
        return
    }
    b.concurrent = true
    // make a copy of the locals - read-only since it will be shared
    v := slc.Clone(b.locals)
    // make them concurrent
    for _, x := range v {
        if x != nil {
            x.SetConcurrent()
        }
    }
    b.locals = v
    if b.this != nil {
        b.this.SetConcurrent()
    }
}
```

This clones the captured locals and marks values as concurrent so they can be safely shared across threads.

## Example

Consider this Suneido code:

```suneido
function makeCounter()
    count = 0
    return {
        count = count + 1
    }
end

counter = makeCounter()
counter()  // returns 1
counter()  // returns 2
```

The flow:

1. `makeCounter()` is called, creating a frame with `count = 0` on the stack
2. The block `{ count = count + 1 }` is compiled with `op.Closure`
3. At runtime, `op.Closure` moves the frame's locals to heap and creates a `SuClosure` with `count` captured
4. The closure is returned
5. When `counter()` is called, a new frame is created using the captured `locals` (containing `count`)
6. The closure increments the captured `count`, persisting the change

## Summary

gSuneido implements closures through:

1. **Compilation**: Closures share the outer function's local variable table
2. **Runtime capture**: The `op.Closure` opcode copies stack locals to heap
3. **Persistent state**: The `SuClosure` struct holds captured variables
4. **Block returns**: The `blockParent` frame reference enables non-local returns
5. **Concurrency**: Cloning ensures thread-safe sharing of captured variables
