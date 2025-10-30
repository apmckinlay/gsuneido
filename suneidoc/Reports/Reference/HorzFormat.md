### HorzFormat

``` suneido
( format, ... )
```

Takes a list of formats and lays them out in a horizontal row, side by side.

The widths of the Horz is the total of the widths of it contents.
Its height is the maximum of its contents.

If displayed with a larger width, 
the excess space is distributed between the formats based on their Xstretch values.