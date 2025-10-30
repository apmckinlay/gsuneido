#### string.SafeEval

``` suneido
() => value
```

Similar to [string.Eval](<string.Eval.md>), but only handles constants, i.e. the result of [Display](<../Display.md>).

**Note:** It is preferable to use SaveEval rather than Eval to avoid security issues.

If the string is not a constant, SafeEval will throw "invalid SafeEval"

For example:

``` suneido
"true".SafeEval() => true
"#(1, 2, 3)".SafeEval() => #(1, 2, 3)
```

Implemented using [string.Compile](<string.Compile.md>)

See also:
[string.ServerEval](<string.ServerEval.md>)