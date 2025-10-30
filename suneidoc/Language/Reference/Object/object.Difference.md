#### object.Difference

``` suneido
(object) => object
```

Returns a new list of the values in this object that are <u>not</u> contained in the passed in object.

Preserves the order of the values.

For example:

``` suneido
#(23, 76, 34, 98, 56).Difference(#(34, 56, 78))
    => #(23, 76, 98)
```

**Note:** Member names are ignored and will not be carried over to the result.


See also:
[object.Disjoint?](<object.Disjoint?.md>),
[object.Intersect](<object.Intersect.md>),
[object.Intersects?](<object.Intersects?.md>),
[object.Subset?](<object.Subset?.md>),
[object.Union](<object.Union.md>)
