<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Eval2

``` suneido
() => object
```

Evaluates the string as Suneido code and returns an object containing the result. If there is no result (e.g. a return with no value) then the result object will be empty.

For example:

``` suneido
"123 + 456".Eval2() => #(579)

"return".Eval2() => #()
```

See also:
[string.Eval](<string.Eval.md>),
[string.ServerEval](<string.ServerEval.md>)