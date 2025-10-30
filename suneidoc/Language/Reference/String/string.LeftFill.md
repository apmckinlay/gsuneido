#### string.LeftFill

``` suneido
(minSize, char = " ") => string
```

Returns a string padded on the left with blanks's (or the specified character) to be not less than minSize characters. If the string is already minSize or greater, no characters are added.

For example:

``` suneido
x = "hello"
x.LeftFill(8, '*')
    => "***hello"
```

See also:
[string.Center](<string.Center.md>),
[string.RightFill](<string.RightFill.md>),
[number.Pad](<../Number/number.Pad.md>)