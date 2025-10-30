#### object.HasNonEmptyMember?

``` suneido
(members) => true or false
```

Returns true if the object has one of the given members and it is non-empty.

For example:

``` suneido
x = #(a: 1, b: 2, c: 3, b: 4)
x.HasNonEmptyMember?(#(c, d))
=> true
```