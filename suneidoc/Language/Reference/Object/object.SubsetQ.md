#### object.Subset?

``` suneido
(object) => true or false
```

Returns true if this is a subset of the specified object, false if not.

For example:

``` suneido
#(a, c).Subset?(#(a, b, c)) => true
```


See also:
[object.Difference](<object.Difference.md>),
[object.Disjoint?](<object.Disjoint?.md>),
[object.Intersect](<object.Intersect.md>),
[object.Intersects?](<object.Intersects?.md>),
[object.Union](<object.Union.md>)
