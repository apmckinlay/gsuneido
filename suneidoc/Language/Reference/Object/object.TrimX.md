#### object.Trim!

``` suneido
(@values) => this
```

Removes all leading and trailing occurrences of the values. This only modifies the unnamed list values of the object.

For example:

``` suneido
Object("", 0, 1, 0, "", 3, "").Trim!(0, "")
    => #(1, 0, "", 3)
```

**Note:** Trim! modifies the object it is applied to, it does not create a new object.


See also:
[object.Delete](<object.Delete.md>),
[object.DeleteIf](<object.DeleteIf.md>),
[object.Erase](<object.Erase.md>),
[object.Remove](<object.Remove.md>),
[object.RemoveIf](<object.RemoveIf.md>),
[object.Without](<object.Without.md>),
[object.WithoutFields](<object.WithoutFields.md>)
