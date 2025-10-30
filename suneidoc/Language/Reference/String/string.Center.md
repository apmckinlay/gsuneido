#### string.Center

``` suneido
(minSize, char = " ") => string
```

Returns a string with text in the centre, and left and right sides are padded with blanks (or a specified character). If the string is already minSize or greater, no characters are added.

For example:

``` suneido
'hello world'.Center(20, '*')
    => "****hello world*****"
```

See also:
[string.RightFill](<string.RightFill.md>),
[string.LeftFill](<string.LeftFill.md>)