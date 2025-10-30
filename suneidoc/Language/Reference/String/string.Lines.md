#### string.Lines

``` suneido
( ) => object
```

Returns a [Sequence](<../Sequence.md>) of newline separated lines from the string. Return characters are not skipped. A newline on the last line is optional. An empty string gives no lines.

For example:

``` suneido
s = "first line\nsecond line\r\nthird line\n"
s.Lines()
    => #("first line", "second line", "third line")
```

**Note**: To process the lines of a file, consider [FileLines](<../FileLines.md>) or [file.Readline](<../File/file.Readline.md>) rather than reading the whole file into memory with GetFile(...).Lines()


See also:
[string.FirstLine](<string.FirstLine.md>),
[string.ChangeEol](<string.ChangeEol.md>),
[string.LineAtPosition](<string.LineAtPosition.md>),
[string.LineCount](<string.LineCount.md>),
[string.LineFromPosition](<string.LineFromPosition.md>),
[string.NthLine](<string.NthLine.md>),
[string.RemoveBlankLines](<string.RemoveBlankLines.md>),
[string.WrapLines](<string.WrapLines.md>)
