#### string.SplitChars

``` suneido
() => object
```

Returns a list of the characters in the string.

For example:

``` suneido
"hello".SplitChars()
```

This can be useful to apply further operations. For example:

``` suneido
"hello".SplitChars().Reverse!().Join() => "olleh"
"hello".SplitChars().Shuffle!().Join() => "lheol"
```

See also: [string.MapN](<string.MapN.md>)


See also:
[string.AfterFirst](<string.AfterFirst.md>),
[string.AfterLast](<string.AfterLast.md>),
[string.BeforeFirst](<string.BeforeFirst.md>),
[string.BeforeLast](<string.BeforeLast.md>),
[string.CutFirst](<string.CutFirst.md>),
[string.CutLast](<string.CutLast.md>),
[string.Split](<string.Split.md>),
[string.SplitCSV](<string.SplitCSV.md>),
[string.SplitFixedLength](<string.SplitFixedLength.md>),
[string.SplitOnFirst](<string.SplitOnFirst.md>),
[string.SplitOnLast](<string.SplitOnLast.md>)
