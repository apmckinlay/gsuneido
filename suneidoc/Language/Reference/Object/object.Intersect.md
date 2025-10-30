#### object.Intersect

``` suneido
(object) => object
```

Returns a new list of the values in this object that are also contained in the passed in object.

Preserves the order of the values.

For example:

``` suneido
#(23, 76, 34, 98, 56).Intersect(#(34, 56, 78))
    => #(34, 56)
```

**Note:** Member names are ignored and will not be carried over to the result.


See also:
[object.Difference](<object.Difference.md>),
[object.Disjoint?](<object.Disjoint?.md>),
[object.Intersects?](<object.Intersects?.md>),
[object.Subset?](<object.Subset?.md>),
[object.Union](<object.Union.md>)
