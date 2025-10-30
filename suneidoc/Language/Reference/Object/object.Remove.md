#### object.Remove

``` suneido
(@values) => this
```

Removes all occurrences of the values from the object.

For example:

``` suneido
Object(12, 34, a: 34, b: 78).Remove(34)
    => #(12, b: 78)
```

**Note:** Remove modifies the object it is applied to, it does not create a new object.


See also:
[object.Delete](<object.Delete.md>),
[object.DeleteIf](<object.DeleteIf.md>),
[object.Erase](<object.Erase.md>),
[object.RemoveIf](<object.RemoveIf.md>),
[object.Trim!](<object.Trim!.md>),
[object.Without](<object.Without.md>),
[object.WithoutFields](<object.WithoutFields.md>)
