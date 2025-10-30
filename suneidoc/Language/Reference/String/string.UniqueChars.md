#### string.UniqueChars

``` suneido
() => string
```

Returns a string containing the unique characters of this string, preserving the order of the characters.

For example:

``` suneido
"abracadabra".UniqueChars()
    => "abrcd"
```

Follow with [string.Divide](<string.Divide.md>) to get the characters as a list in an object, for example:

``` suneido
"abracadabra".UniqueChars().Divide()
    => #("a", "b", "r", "c", "d")
```

See also:
[object.UniqueValues](<../Object/object.UniqueValues.md>)