# Closure Implementation Changes in gSuneido

This document details the recent architectural changes to how closures are implemented in gSuneido.

## The Problem: "Detached" Locals

In the previous implementation (documented in [`docs/Closures.md`](/docs/Closures.md)), closures shared the outer function's local variable table. When a closure was created at runtime (via the `op.Closure` opcode), the entire local variable stack of the enclosing frame was copied to the heap. 

The critical flaw emerged when closures became concurrent. The `SetConcurrent` method would clone the captured locals to allow independent mutation. This meant that if multiple concurrent closures were created from the same outer function, or if the outer function continued to execute concurrently with the closure, they each received their own independent copy of the locals. 

Changes made by one concurrent call were not visible to others. This led to confusing behavior because the closure *appeared* to share state with its siblings and parent, but in reality, the state was "detached."

Another issue was that closure parameters were stored in the locals (at a known offset) which made recursion problematic.

## The Solution: Explicit Shared Storage

The new architecture solves this by explicitly separating local variables (which belong to a specific function/closure invocation) from shared variables (which are captured and shared between a function invocation and its closures).

NOTE: Previously the same variable name in sibling blocks was assumed to be sharing. This is no longer the case. This is a potential breaking change, but the fix is to simple reference the variable name(s) in an outer scope.

### 1. Variable Index Space Encoding

Instead of introducing new opcodes for shared variables (like `LoadShared`, `StoreShared`), the new design cleverly encodes the storage location directly into the variable index (slot number):

*   **Indexes 0-191**: Local variables (stack-allocated).
*   **Indexes 192-255**: Shared variables (heap-allocated). The actual index into the shared storage is `stored_index - 192`.

This approach keeps the bytecode compact and works naturally with existing multi-variable opcodes (like `ForIn2`).

### 2. The `Shared` Struct

A new `Shared` struct was introduced in [`frame.go`](/core/frame.go):

```go
type Shared struct {
	values []Value
	MayLock
}
```

This struct holds the actual values of the captured variables. Crucially, it includes a `MayLock` (mutex) to support safe concurrent access.

### 3. Frame Layout Updates

The execution `Frame` ([`frame.go`](/core/frame.go)) was updated to include a reference to this shared storage:

```go
type Frame struct {
    // ...
	locals []Value   // Local variables (on the thread stack)
	shared *Shared   // Shared variable storage for closures
    // ...
}
```

Helper methods `getSlot(idx)` and `setSlot(idx, val)` were added to `Frame`. These methods inspect the index: if it's `< 192` (`SharedSlotStart`), they access `locals`; if it's `>= 192`, they access `shared.values` (applying locks if the shared block is marked concurrent).

### 4. Interpreter Routing

The main interpreter loop in [`core/interp.go`](/core/interp.go) was updated to use `getSlot` and `setSlot` instead of directly indexing into `fr.locals`. This means all variable access opcodes (`Load`, `Store`, `LoadStore`, `Dyload`, etc.) automatically route to the correct storage location based on the index space encoding.

### 5. Function Entry and Closure Creation

*   **Function Entry**: When a function is called ([`Thread.invoke`](/core/interp.go)), it allocates stack space for its locals. If the function has shared variables (determined by `len(fn.Names) > fn.Nstack`), it allocates a new `Shared` block on the heap. It then calls [`Frame.moveLocalsToShared()`](/core/frame.go) to copy any parameters that are captured into the shared block.
    *   Shared parameters are still initially present in locals on entry (due to arguments being pushed and passed on the stack), but the shared slot becomes the authoritative location for dynamic lookup and closure sharing.
*   **Closure Creation**: When a closure is created by `op.Closure` it no longer copies the entire local stack to the heap. Instead, it simply captures a reference to the parent frame's `shared` block.
*   **Closure Entry**: When a closure is called ([`SuClosure.Call`](/core/suclosure.go)), it sets up its frame similarly to a normal function. It gets its own stack space for its local variables, but it uses the `shared` reference captured during creation.

### 6. Compiler Slot Assignment (`compile/ast/blocks.go`)

The compiler was updated to pre-assign these slot indexes during the AST processing phase. The [`Blocks`](/compile/ast/blocks.go) function now performs a three-phase analysis:

1.  **Collect**: It traverses the AST, creating a `scope` for the outer function and each nested block. It assigns slot indexes starting at `0` for parameters and records all other variables as unassigned (`-1`).
2.  **Assign Shared**: It walks up the parent chain for each block. If a block uses a variable that exists in a parent scope, that variable is assigned a shared slot index (starting at `192` / `SharedSlotStart`). This ensures that all closures within the same outer function agree on the slot index for a given shared variable.
3.  **Assign Local**: Finally, any remaining unassigned variables in each scope are given local stack slot indexes (starting after the parameters).

Because shared slots are allocated across the full outer-scope sharing set, [`ParamSpec.Names`](/core/paramspec.go) can contain apparent "gaps" (where the name is "") from the perspective of an individual scope: names inherited for shared-slot alignment may be present even when that specific scope does not directly use them.

This pre-assignment allows the code generator to simply emit `Load` or `Store` opcodes with the correct index, knowing the runtime will route it appropriately.

The new compiler slot assignment has to be handled by [`Thread.locals()`](/core/thread.go)

## Why This is Better

1.  **True Sharing**: All closures created within the same outer function invocation now share the *exact same* `Shared` struct instance. If one closure modifies a shared variable, the change is immediately visible to the outer function and any sibling closures.
2.  **Correct Concurrency**: When a closure is marked concurrent (`SuClosure.SetConcurrent](/core/suclosure.go)), it sets the `concurrent` flag on the `Shared` struct. The `getSlot` and `setSlot` methods then use the mutex to ensure thread-safe access to the shared variables. The state is no longer detached; it is safely shared.
3.  **Efficiency**: Non-captured local variables remain on the stack, avoiding unnecessary heap allocations. The overhead of checking the index (`< 192`) is minimal and highly predictable for modern CPUs.

In summary, the changes move gSuneido from a "copy-on-capture" model (which broke down under concurrency) to a "reference-shared-state" model, ensuring correct semantics for closures in all contexts.
