<div style="float:right"><span class="builtin">Builtin</span></div>

#### TypeChecker.RegisterSignatures

``` suneido
(signatures :object) => number
```

Supplies the checker with signatures for code written in Suneido, e.g. the Object and String methods in stdlib. Returns the number of signatures registered.

Each element of **signatures** is an object with:

`kind`
: "method" (the default), "free", or "static"

`receiver`
: For methods, the type the method is called on. One of "string", "number", "object", "record" (same as object), "date", "class", "function", or "sequence". 

`class`
: For statics, the name of the global class. Required for statics, not allowed otherwise.

`name`
: The name of the method or function.

`sig`
: The declaration, in the same format as builtin declarations, e.g. "(n :number = 0) :string". A missing return type is treated as `unknown`.

For example:

``` suneido
TypeChecker.RegisterSignatures(#(
	(receiver: object, name: Map, sig: "(block) :object"),
	(receiver: string, name: Lines, sig: "() :object"),
	(kind: free, name: Datetime, sig: "(string) :date"),
	(kind: static, class: Database, name: Check, sig: "() :string")))
=> 4
```

Each call replaces everything registered by the previous call, it does not add to it.
The builtin signatures however can not be replaced and/or modified.


See also: [TypeChecker.Annotations](<TypeChecker.Annotations.md>)