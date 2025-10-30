### range.Includes?

``` suneido
( range ) => boolean
```

where **range** is a Range type.

Returns True if the current range includes the **range** passed in; otherwise, returns False.

For example:

``` suneido
range = Range(2, 7)
range.Includes?(Range(3, 7)) => True
range.Includes?(Range(1, 7)) => False
```