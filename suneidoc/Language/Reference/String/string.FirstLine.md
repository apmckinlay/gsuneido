#### string.FirstLine

``` suneido
(pos) => line
```

Returns the first line. i.e. up to the first '\n' with any trailing '\r' removed.

For example:

``` suneido
"one\r\ntwo\nthree".FirstLine()
    => "two"
```


See also:
[string.ChangeEol](<string.ChangeEol.md>),
[string.Lines](<string.Lines.md>),
[string.LineAtPosition](<string.LineAtPosition.md>),
[string.LineCount](<string.LineCount.md>),
[string.LineFromPosition](<string.LineFromPosition.md>),
[string.NthLine](<string.NthLine.md>),
[string.RemoveBlankLines](<string.RemoveBlankLines.md>),
[string.WrapLines](<string.WrapLines.md>)
