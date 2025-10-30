#### report.SelectFont

``` suneido
( font ) => old_font
```

Where font is an object containing one or more of:
name
: the font name e.g. "Courier"

size
: the font size in points, e.g. 14,
or a relative size as a string, e.g. "-1" or "+2"

weight
: 0 means don't care, 1 for lightest, 1000 for boldest

italic
: true or false

underline
: true or false

The default font is #(name: "Arial", size: 10, weight: 400)