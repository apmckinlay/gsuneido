#### string.SplitOnLast

``` suneido
(separator = " ") => object
```

Splits the string at the last occurrence of the specified separator string.

Note: the separator string cannot be empty ("").

For example:

``` suneido
"metric tons 123".SplitOnLast() => Object("metric tons", "123")
```

See also:
[string.Split](<string.Split>)[string.SplitOnFirst](<string.SplitOnFirst>)[object.Join](<../Object/object.Join.md>)