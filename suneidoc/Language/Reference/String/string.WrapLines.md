#### string.WrapLines

``` suneido
(width, stripBlank = false) => object
```

Split the string into lines no wider than width. Breaks on spaces or tabs where possible.

For example:

``` suneido
"hello there world".WrapLines(15)
    => #("hello there", "world")
```


See also:
[string.FirstLine](<string.FirstLine.md>),
[string.ChangeEol](<string.ChangeEol.md>),
[string.Lines](<string.Lines.md>),
[string.LineAtPosition](<string.LineAtPosition.md>),
[string.LineCount](<string.LineCount.md>),
[string.LineFromPosition](<string.LineFromPosition.md>),
[string.NthLine](<string.NthLine.md>),
[string.RemoveBlankLines](<string.RemoveBlankLines.md>)
