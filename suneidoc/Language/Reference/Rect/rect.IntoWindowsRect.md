### rect.IntoWindowsRect

``` suneido
( object ) => object
```

Copies the position and size of the current rectangle into a pre-existing Windows RECT structure, **object**, and returns a reference to **object**.

For example:

``` suneido
rect = Rect(50, 100, 10, 20)
winrc = rect.ToWindowsRect()
    => #(left: 50, top: 100, right: 60, bottom: 120)
rect.Set(x: 0, y: -5, h: 10)
rect.IntoWindowsRect(winrc)
    => #(left: 0, top: -5, right: 10, bottom: 5)
```

See also: [rect.ToWindowsRect](<rect.ToWindowsRect.md>), [Rect.FromWindowsRect](<Rect.FromWindowsRect.md>)