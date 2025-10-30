#### number.Pad

``` suneido
(minSize, char = '0') => string
```

Returns a string padded on the left with zero's (or the specified character) to be not less than minSize characters. If the number converted to a string is already minSize or greater, no padding is added.

For example:

``` suneido
x = 6
x.Pad(2)
    => "06"
```

**Note:** The absolute value of the number is used, so the sign is ignored.

See also:
[string.LeftFill](<../String/string.LeftFill.md>)