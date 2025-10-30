#### object.MergeNew

``` suneido
(object) => this
```

Adds object members from the object passed in that do not exist in this object.

**Note:** this modifies the object it is applied to, it does not create a new object.

For example:

``` suneido
ob = Object(12, 24, a: 1, b: 2, c: 1)
Object(1, 2, 3, c: 3, d: 4).MergeNew(ob)
    => #(1, 2, 3, c: 3, b: 2, a: 1, d: 4)
```

See also:
[object.Merge](<object.Merge.md>),
[object.MergeUnion](<object.MergeUnion.md>)