### rect.GetY2

``` suneido
( ) => number
```

Returns the y coordinate of the rectangle's "bottom side".

For example:

``` suneido
rect = Rect(50, 100, 10, 20)
rect.GetY2() => 120
```

**Note:** If the height of the rectangle is negative then `rect.GetY2() < rect.GetY()`.

See also: [rect.GetY](<rect.GetY.md>)