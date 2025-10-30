#### CanvasText

``` suneido
(text, x1, y1, x2, y2)
```

Uses the coordinates passed in to draw the text using DrawText.

For example: 

``` suneido
c = CanvasControl().Ctrl
text = CanvasText('Hello World', 50, 50, 100, 100)
text.SetColor(c.GetColor())
c.AddItem(text)
```

Will draw: 

![](<../../../res/canvastext.png>)