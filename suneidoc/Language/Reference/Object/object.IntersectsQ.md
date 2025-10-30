#### object.Intersects?

``` suneido
(object) => true or false
```

Returns true if the two objects have any common values. For example:

``` suneido
#(1, 2, 3).Intersects?(#(4, 5, 6))
    => false

#(1, 2, 3).Intersects?(#(2, 4))
    => true
```

Intersects? is the opposite of [object.Disjoint?](<object.Disjoint?.md>)


See also:
[object.Difference](<object.Difference.md>),
[object.Disjoint?](<object.Disjoint?.md>),
[object.Intersect](<object.Intersect.md>),
[object.Subset?](<object.Subset?.md>),
[object.Union](<object.Union.md>)
