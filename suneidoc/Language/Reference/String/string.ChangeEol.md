#### string.ChangeEol

``` suneido
(eol) => line
```

Replaces line endings ('\n' or '\r\n') with the specified string.

For example:

``` suneido
"one\ntwo\r\nthree".ChangeEol('\r\n').Escape()
    => "one\r\ntwo\r\nthree"
```


See also:
[string.FirstLine](<string.FirstLine.md>),
[string.Lines](<string.Lines.md>),
[string.LineAtPosition](<string.LineAtPosition.md>),
[string.LineCount](<string.LineCount.md>),
[string.LineFromPosition](<string.LineFromPosition.md>),
[string.NthLine](<string.NthLine.md>),
[string.RemoveBlankLines](<string.RemoveBlankLines.md>),
[string.WrapLines](<string.WrapLines.md>)
