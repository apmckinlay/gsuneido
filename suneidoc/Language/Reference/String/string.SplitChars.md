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

See also:
[string.Split](<string.Split>)[string.MapN](<string.MapN.md>),
[object.Join](<../Object/object.Join.md>)