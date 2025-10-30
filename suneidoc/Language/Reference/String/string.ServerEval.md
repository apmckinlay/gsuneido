<div style="float:right"><span class="deprecated">Deprecated</span><span class="builtin">Builtin</span></div>

#### string.ServerEval

``` suneido
() => value
```

Evaluates the string as Suneido code *on the server*.

For example:

``` suneido
"123 + 456".ServerEval() => 579
```

**Warning:** Any newlines in the string are silently replaced by blanks.

If you're running standalone, string.ServerEval is identical to
[string.Eval](<string.Eval.md>).

**Note**: string.ServerEval is deprecated and may be removed in the future. Use [ServerEval](<../ServerEval.md>) instead