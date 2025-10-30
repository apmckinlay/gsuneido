#### object.Disjoint?

``` suneido
(object) => true or false
```

Returns true if the two objects do not have any common values. For example:

``` suneido
#(1, 2, 3).Disjoint?(#(4, 5, 6))
    => true

#(1, 2, 3).Disjoint?(#(2, 4))
    => false
```

Disjoint? is the opposite of [object.Intersects?](<object.Intersects?.md>)


See also:
[object.Difference](<object.Difference.md>),
[object.Intersect](<object.Intersect.md>),
[object.Intersects?](<object.Intersects?.md>),
[object.Subset?](<object.Subset?.md>),
[object.Union](<object.Union.md>)
