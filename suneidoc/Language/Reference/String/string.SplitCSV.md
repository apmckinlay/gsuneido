#### string.SplitCSV

``` suneido
(fields = false) => object
```

SplitCSV extracts the values from a comma separated string.

-	double quotes around strings are removed
-	pairs of double quotes within strings are replaced with single double quotes
-	numeric values will be converted to numbers

If **fields** is specified then the values will be placed into named members.

For example:

``` suneido
'123,"35"" of string"'.SplitCSV()
    => #(123, '35" of string')
```

Or with named members:

``` suneido
'"Andrew McKinlay",41'.SplitCSV(#("name", "age"))
    => #(name: "Andrew McKinlay", age: 41)
```

See also: [object.JoinCSV](<../Object/object.JoinCSV.md>)


See also:
[string.AfterFirst](<string.AfterFirst.md>),
[string.AfterLast](<string.AfterLast.md>),
[string.BeforeFirst](<string.BeforeFirst.md>),
[string.BeforeLast](<string.BeforeLast.md>),
[string.CutFirst](<string.CutFirst.md>),
[string.CutLast](<string.CutLast.md>),
[string.Split](<string.Split.md>),
[string.SplitChars](<string.SplitChars.md>),
[string.SplitFixedLength](<string.SplitFixedLength.md>),
[string.SplitOnFirst](<string.SplitOnFirst.md>),
[string.SplitOnLast](<string.SplitOnLast.md>)
