#### object.Merge

``` suneido
(object) => this
```

Adds to the current object each member from the object passed in. Both named and un-named members in the current object will be overwritten by corresponding members in the passed in argument.

**Note:** this modifies the object it is applied to, it does not create a new object.

For example:

``` suneido
ob = Object(12, 24, a: 1, b: 2)
Object(1, 2, 3, c: 3, d: 4).Merge(ob)
    => #(12, 24, 3, a: 1, b: 2, c: 3, d: 4)
```

See also:
[object.Append](<object.Append.md>)