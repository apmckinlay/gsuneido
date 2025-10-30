### string.LineCount

``` suneido
() => number
```

Returns the number of lines in the string.

For example:

``` suneido
"".LineCount()
    => 0

"hello world".LineCount()
    => 1

"hello world\n".LineCount()
    => 1

"hello\nworld".LineCount()
    => 2

"hello\nworld\n".LineCount()
    => 2
```


See also:
[string.FirstLine](<string.FirstLine.md>),
[string.ChangeEol](<string.ChangeEol.md>),
[string.Lines](<string.Lines.md>),
[string.LineAtPosition](<string.LineAtPosition.md>),
[string.LineFromPosition](<string.LineFromPosition.md>),
[string.NthLine](<string.NthLine.md>),
[string.RemoveBlankLines](<string.RemoveBlankLines.md>),
[string.WrapLines](<string.WrapLines.md>)



See also:
[string.Count](<string.Count.md>)
