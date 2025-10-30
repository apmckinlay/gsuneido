#### CanvasArc

``` suneido
(left, top, right, bottom, xStartArc, yStartArc, xEndArc, yEndArc)
```

Uses the coordinates passed in to draw an arc using Arc.

For example:

``` suneido
c = CanvasControl().Ctrl
c.AddItem(CanvasArc(50, 50, 200, 200, 30, 30, 700, 60))
```

Will draw:

![](<../../../res/canvasarc.png>)