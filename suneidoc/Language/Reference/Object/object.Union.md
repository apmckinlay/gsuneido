#### object.Union

``` suneido
(object) => object
```

Returns a new list of the distinct values from this object and the passed in object.

Preserves the order of the values.

For example:

``` suneido
#(23, 76, 34, 98, 56).Union(#(34, 56, 78))
    => #(23, 76, 34, 98, 56, 78)
```

**Note:** Member names are ignored and will not be carried over to the result.


See also:
[object.Difference](<object.Difference.md>),
[object.Disjoint?](<object.Disjoint?.md>),
[object.Intersect](<object.Intersect.md>),
[object.Intersects?](<object.Intersects?.md>),
[object.Subset?](<object.Subset?.md>)
