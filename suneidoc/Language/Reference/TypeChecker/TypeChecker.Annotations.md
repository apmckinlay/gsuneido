<div style="float:right"><span class="builtin">Builtin</span></div>

#### TypeChecker.Annotations

``` suneido
() => object
```

Returns the signatures declared by gsuneido builtins - the base layer of the table the checker infers from. It does not include the signatures added by
[TypeChecker.RegisterSignatures](<TypeChecker.RegisterSignatures.md>).

Each element is an object with:

`kind`
: "method" (called on a value), "free" (a global function), or "static" (a method on a builtin class)

`prefix`
: The type the method is called on, e.g. "string", or the class for a static, e.g. "db". Empty for free functions.

`name`
: The name of the method or function.

`sig`
: The declaration, e.g. "(n) :string"

For example:

``` suneido
TypeChecker.Annotations().Filter({ it.name is "Split" })
=> #(#(kind: "method", prefix: "string", name: "Split",
    sig: "(separator :string = '') :object"))
```

The result is read-only, and the same object each time.

See also: [TypeChecker.RegisterSignatures](<TypeChecker.RegisterSignatures.md>)