#### CanvasPolygon

``` suneido
(#((x, y), (x, y), (x, y) ...))
```

Uses the points passed in to draw a polygon using Polygon.

For example:

``` suneido
c = CanvasControl().Ctrl
c.AddItem(CanvasPolygon(#((x: 50, y: 50), (x: 150, y: 50), (x: 100, y: 150))))
```

Will draw:

![](<../../../res/canvaspolygon.png>)