### rect.ToWindowsRect

``` suneido
( ) => object
```

Returns an object compatible with the Windows RECT structure.

For example:

``` suneido
rect = Rect(50, 100, 10, 20)
rect.ToWindowsRect()
    => #(left: 50, top: 100, right: 60, bottom: 120)
```

**Note:** This method was formerly called "`GetWindowsRect`" and has been renamed.

See also: [Rect.FromWindowsRect](<Rect.FromWindowsRect.md>), [rect.IntoWindowsRect](<rect.IntoWindowsRect.md>), [point.ToWindowsPoint](<../Point/point.ToWindowsPoint.md>)