<div style="float:right"><span class="builtin">Builtin</span></div>

#### TypeChecker.Infer

``` suneido
(arguments :object, references = #(), config = #()) => object
```

Infers the types in each of the **arguments** sources
`arguments`
: arguments are object in order of source starting with the base class and ending with the current class

`references`
: an optional list, in the same form as arguments. References are not checked themselves and produce no diagnostics instead only used to infer return types

`config`
: An optional object of options - see below.

For example:

``` suneido
TypeChecker.Infer(#('class { Total(n :number) { return n * 2 } }'))
=>
#(  method: "TypeInfer",
    result: #(#(methods: #(Total: #("$return": "number", n: "number")), members: #())),
    diagnostics: #(errors: #(), warnings: #()),
    version: "Aug 18 2026 09:36 (go1.27.0 windows/amd64)")
```

A class that calls something else needs that something else in **references**, otherwise its type is unknown:

``` suneido
invoice = 'class
	{
	Total(amount :number)
		{
		return amount * TaxRate()
		}
	}'
taxrate = [name: "TaxRate", src: 'function () { return .07 }']

TypeChecker.Infer([invoice], [taxrate]).diagnostics
=> #(errors: #(), warnings: #())

TypeChecker.Infer([invoice]).diagnostics.warnings[0].msg
=> '[0.25] operator "*" operand has type unknown cannot prove it is number'
```

##### Config

`strictStringConcat`
: "off", "warn", or "error" (the default). How to report a `$` operand that isn't known to be a string.

`strictCrossTypeCompares`
: "off", "warn" (the default), or "error". How to report `<` `<=` `>` `>=` between operands of different types, which falls back to Suneido's cross-type ordering.

`confidence`
: Only keep diagnostics whose confidence matches, e.g. ">=0.70". A bare number means ">=". `>`, `<=`, `<`, and `==` also work, as does "all" (the default).

Unrecognized config members are ignored. A bad value throws.

``` suneido
src = 'class { Label(n :number) { return "n=" $ n } }'

TypeChecker.Infer([src]).diagnostics.errors[0]
=> #(class: "Class0", method: "Label", pos: 41, line: 1, col: 42,
    msg: '[0.20] operator "$" expects string, got number', flag: "strictStringConcat")

TypeChecker.Infer([src], #(), [strictStringConcat: "off"]).diagnostics
=> #(errors: #(), warnings: #())
```

See also: [TypeChecker.Annotate](<TypeChecker.Annotate.md>)