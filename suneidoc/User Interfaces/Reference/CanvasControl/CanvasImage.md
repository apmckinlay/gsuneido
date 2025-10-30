#### CanvasImage

``` suneido
(image, left, top, right, bottom)
```

Uses the coordinates passed in to draw an image using 
[Image](<../../../Language/Reference/Image.md>)

For example:

``` suneido
c = CanvasControl().Ctrl
c.AddItem(CanvasImage('c:/Windows/Gone Fishing.bmp', 20, 20, 100, 100))
```

Will draw:

![](<../../../res/canvasimage.png>)