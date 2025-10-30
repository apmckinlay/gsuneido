#### object.Count

``` suneido
(value) => count
```

Returns the number of occurrences of the value. For example:

``` suneido
#(1, 4, 3, 6, 3, 7).Count(3)
    => 2
```

.Count() is equivalent to .Size() but for sequences it will avoid instantiation.

See also:
[object.CountIf](<object.CountIf.md>)