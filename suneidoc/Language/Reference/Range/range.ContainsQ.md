### range.Contains?

``` suneido
( value ) => True or False
```

Returns True if **value** is between the *low* and *high* ends of range, inclusive; otherwise, returns False.

For example:

``` suneido
range = Range(2, 7)
range.Contains?(1) => False
range.Contains?(2) => True
range.Contains?(5) => True
range.Contains?(7) => True
range.Contains?(8) => False
```