#### string.SplitOnFirst

``` suneido
(delimiter = " ") => object
```

**DEPRECATED** - prefer [string.CutFirst](<string.CutFirst.md>)

Splits the string at the first occurrence of the specified delimiter string.

For example:

``` suneido
"123 metric tons".SplitOnFirst() => Object("123", "metric tons")
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
[string.SplitOnLast](<string.SplitOnLast.md>)
