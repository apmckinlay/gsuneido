### Point.FromWindowsPoint

``` suneido
( object ) => point
```

Constructs and returns a point whose position is equivalent to the given Windows POINT structure. This method is "static" in the sense that it is called on the class Point, not an instance of the class.

For example:

``` suneido
winpt = #(x: 1, y: 1)
Point.FromWindowsPoint(winpt)
    => Point(1, 1)
```

See also: [point.ToWindowsPoint](<point.ToWindowsPoint.md>), [rect.FromWindowsRect](<../Rect/Rect.FromWindowsRect.md>)