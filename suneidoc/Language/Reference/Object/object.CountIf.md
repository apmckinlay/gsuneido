#### object.CountIf

``` suneido
(block) => count
```

Returns the number of values for which block returns true. block can be anything callable. It is passed a value. For example:

``` suneido
#(1, 4, 3, 6, 3, 7).CountIf({ it > 3 })
    => 3
```

See also:
[object.Count](<object.Count.md>)