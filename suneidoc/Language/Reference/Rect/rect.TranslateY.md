### rect.TranslateY

``` suneido
(dy)
```

Shifts the y coordinate of the rectangle by **dy**. Equivalent to calling `rect.SetY(rect.GetY() + dy)`.

For example:

``` suneido
rect = Rect(50, 100, 10, 20)
rect.TranslateY(50) => Rect(50, 150, 10, 20)
```

See also: [rect.Translate](<rect.Translate.md>), [point.Translate](<../Point/point.Translate.md>), [rect.TranslateX](<rect.TranslateX.md>), [rect.SetY](<rect.SetY.md>)