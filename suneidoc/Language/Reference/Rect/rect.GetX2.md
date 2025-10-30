### rect.GetX2

``` suneido
( ) => number
```

Returns the x coordinate of the rectangle's "right side".

For example:

``` suneido
rect = Rect(50, 100, 10, 20)
rect.GetX2() => 60
```

**Note:** If the width of the rectangle is negative then `rect.GetX2() < rect.GetX()`.

See also: [rect.GetX](<rect.GetX.md>)