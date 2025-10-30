#### string.SplitOnFirst

``` suneido
(separator = " ") => object
```

Splits the string at the first occurrence of the specified separator string.

Note: the separator string cannot be empty ("").

For example:

``` suneido
"123 metric tons".SplitOnFirst() => Object("123", "metric tons")
```

See also:
[string.Split](<string.Split>)[string.SplitOnLast](<string.SplitOnLast>)[object.Join](<../Object/object.Join.md>)