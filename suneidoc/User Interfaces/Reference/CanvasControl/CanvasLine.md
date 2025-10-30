#### CanvasLine

``` suneido
(x1, y1, x2, y2)
```

Uses the coordinates passed in to draw a line using MoveTo and LineTo.

For example:

``` suneido
c = CanvasControl().Ctrl
c.AddItem(CanvasLine(30, 30, 60, 100))
```

Will draw:

![](<../../../res/canvasline.png>)