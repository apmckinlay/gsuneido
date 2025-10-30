#### string.LineAtPosition

``` suneido
(pos) => line
```

Returns the line containing the specified position.

For example:

``` suneido
"one\ntwo\r\nthree".LineAtPosition(5)
    => "two"
```


See also:
[string.FirstLine](<string.FirstLine.md>),
[string.ChangeEol](<string.ChangeEol.md>),
[string.Lines](<string.Lines.md>),
[string.LineCount](<string.LineCount.md>),
[string.LineFromPosition](<string.LineFromPosition.md>),
[string.NthLine](<string.NthLine.md>),
[string.RemoveBlankLines](<string.RemoveBlankLines.md>),
[string.WrapLines](<string.WrapLines.md>)
