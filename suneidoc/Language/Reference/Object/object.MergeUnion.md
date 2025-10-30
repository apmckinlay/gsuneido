#### object.MergeUnion

``` suneido
(object) => this
```

Adds un-named values from the object passed in that do not exist in this.

**Note:** this modifies the object it is applied to, it does not create a new object.

For example:

``` suneido
ob = Object(3, 4)
Object(1, 2, 3).MergeUnion(ob)
    => #(1, 2, 3, 4)
```

See also:
[object.Merge](<object.Merge.md>),
[object.MergeNew](<object.MergeNew.md>)