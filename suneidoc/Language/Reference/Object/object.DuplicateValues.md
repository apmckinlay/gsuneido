#### object.DuplicateValues

``` suneido
() => object
```

Returns an object containing all values that are duplicated in the original object. This includes values that are in named members.

For example:

``` suneido
ob = #(12, 34, 34, 36, a: 56, b: 12)
ob.DuplicateValues()
    => #(34, 12)
```