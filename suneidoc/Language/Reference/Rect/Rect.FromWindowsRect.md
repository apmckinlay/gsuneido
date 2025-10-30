### Rect.FromWindowsRect

``` suneido
( object ) => rect
```

Constructs and returns a rectangle whose size and position are equivalent to the given Windows RECT structure. This method is "static" in the sense that it is called on the class Rect, not an instance of the class.

For example:

``` suneido
winrc = #(left: 50, top: 100, right: 60, bottom: 120)
Rect.FromWindowsRect(winrc)
    => Rect(50, 100, 10, 20)
```

See also: [rect.ToWindowsRect](<rect.ToWindowsRect.md>), [rect.IntoWindowsRect](<rect.IntoWindowsRect.md>), [Point.FromWindowsPoint](<../Point/Point.FromWindowsPoint.md>)