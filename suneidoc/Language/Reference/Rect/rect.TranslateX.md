### rect.TranslateX

``` suneido
(dx)
```

Shifts the x coordinate of the rectangle by **dx** Equivalent to calling `rect.SetX(rect.GetX() + dx)`.

For example:

``` suneido
rect = Rect(50, 100, 10, 20)
rect.TranslateX(-50) => Rect(0, 100, 10, 20)
```

See also: [rect.Translate](<rect.Translate.md>), [point.Translate](<../Point/point.Translate.md>), [rect.TranslateY](<rect.TranslateY.md>), [rect.SetX](<rect.SetX.md>)