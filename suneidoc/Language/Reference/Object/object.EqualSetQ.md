#### object.EqualSet?

``` suneido
(object) => true / false
```

Returns true if the object contain the same set of values, regardless of order. Member names are ignored.

For example:

``` suneido
#(4, 6, 3, 8, 2).EqualSet?(#(3, 2, 6, 8, 4))
    => true
```