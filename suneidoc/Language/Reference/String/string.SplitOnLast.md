#### string.SplitOnLast

``` suneido
(delimiter = " ") => object
```

**DEPRECATED** - prefer [string.CutLast](<string.CutLast.md>)

Splits the string at the last occurrence of the specified delimiter string.

For example:

``` suneido
"metric tons 123".SplitOnLast() => Object("metric tons", "123")
```


See also:
[string.AfterFirst](<string.AfterFirst.md>),
[string.AfterLast](<string.AfterLast.md>),
[string.BeforeFirst](<string.BeforeFirst.md>),
[string.BeforeLast](<string.BeforeLast.md>),
[string.CutFirst](<string.CutFirst.md>),
[string.CutLast](<string.CutLast.md>),
[string.Split](<string.Split.md>),
[string.SplitCSV](<string.SplitCSV.md>),
[string.SplitChars](<string.SplitChars.md>),
[string.SplitFixedLength](<string.SplitFixedLength.md>),
[string.SplitOnFirst](<string.SplitOnFirst.md>)
