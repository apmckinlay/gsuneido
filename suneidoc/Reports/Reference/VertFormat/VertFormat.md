#### VertFormat

``` suneido
( format, ... )
```

Takes a list of formats and lays them out in a vertical column, one below the other.

The height of the Vert is the total of the heights of it contents.  
Its width is the maximum of its contents.

If displayed with a larger height, 
the excess space is distributed between the formats based on their Ystretch values.