### range.Minus

``` suneido
( range ) => range
```

where **range** is a Range type.

Subtracts the current range from **range** passing in.

For example:

``` suneido
range = Range(2, 7)
range.Minus(Range(1, 4)) => Range(4, 7)
```