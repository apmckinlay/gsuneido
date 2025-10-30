#### string.RemoveBlankLines

``` suneido
() => string
```

Returns the string with blank lines (zero or more spaces or tabs) removed.

For example:

``` suneido
"one\n \r\n \nthree".RemoveBlankLines().Escape()
    => "one\nthree"
```


See also:
[string.FirstLine](<string.FirstLine.md>),
[string.ChangeEol](<string.ChangeEol.md>),
[string.Lines](<string.Lines.md>),
[string.LineAtPosition](<string.LineAtPosition.md>),
[string.LineCount](<string.LineCount.md>),
[string.LineFromPosition](<string.LineFromPosition.md>),
[string.NthLine](<string.NthLine.md>),
[string.WrapLines](<string.WrapLines.md>)
