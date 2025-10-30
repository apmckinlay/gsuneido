### range.Overlaps?

``` suneido
( range ) => True or False
```

where **range** is a Range.

Returns True if **range** overlaps with the current range.

For example:

``` suneido
range = Range(2, 7)
range.Overlaps?(Range(3, 7)) => True
range.Overlaps?(Range(8, 15)) => False
```