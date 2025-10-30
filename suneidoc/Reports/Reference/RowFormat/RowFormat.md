#### RowFormat

``` suneido
( format, ... data = False, font = False )
```

`Row` is the default output for `QueryFormat`.

Row constructs each format.
It then determines the correct font size (in the range 4 to 14 points) to fill the page width.

A Row will be variable height if any of its formats are variable.

Formats can have a "Span" member to allow them to use the space of following fields.
A Span of 1 will combine the width of the current column and one following.
Note: if the spanned columns had their own formats, they will be ignored.