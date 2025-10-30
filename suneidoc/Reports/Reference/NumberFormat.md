### NumberFormat

``` suneido
(value = "", mask = false, width = false, w = false, font = false,
        color = false)
```

Similar to [OptionalNumberFormat](<OptionalNumberFormat.md>) except that a value of "" prints zero, instead of nothing.

width is in characters, based on average character width, and therefore dependent on the font.  w is in twips.  If there is a mask, its width takes precedence over the width and w arguments.

See [number.Format](<../../Language/Reference/Number/number.Format.md>) for an explanation of the mask.

Derived from [TextFormat](<TextFormat.md>) via OptionalNumberFormat