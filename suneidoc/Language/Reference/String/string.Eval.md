<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Eval

``` suneido
() => value
```

Evaluates the string as Suneido code.

For example:

``` suneido
"123 + 456".Eval() => 579
```

**Note:** If the code does not return anything (e.g. a return with no value) then "" will be returned. This is not the same behavior as ServerEval, string.ServerEval, or [object.Eval](<../Object/object.Eval.md>)

**WARNING:** Eval is a **potential security risk**, especially when the string is coming from the user or some other external source. For example, the string could be 'System("del *.*")'. Wherever possible, use Global, string.Compile, or [string.SafeEval](<string.SafeEval.md>) instead

See also:
[string.Eval2](<string.Eval2.md>),
[ServerEval](<../ServerEval.md>),
[string.ServerEval](<string.ServerEval.md>),