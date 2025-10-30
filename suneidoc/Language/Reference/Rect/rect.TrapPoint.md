### rect.TrapPoint

``` suneido
( p ) => q 
```

where **p** and **q** are instances of Point.

If the current rectangle contains the point **p**, returns **p**. Otherwise, returns a new point **q** which is the closest point to **p** that is contained within the current rectangle.

For example:

``` suneido
rect = Rect(50, 100, 10, 20)
p = Point(55, 115)
rect.TrapPoint(p)
    => p
p.SetX(45)
rect.TrapPoint(p)
    => Point(50, 115)
```

See also: [rect.ContainsPoint?](<rect.ContainsPoint?.md>)