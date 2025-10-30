#### object.Without

``` suneido
(@values) => object
```

Returns a new object without any occurrences of the values. It does not modify the original object.

For example:

``` suneido
Object(12, 34, a: 34, b: 78).Without(34)
    => #(12, b: 78)
```


See also:
[object.Delete](<object.Delete.md>),
[object.DeleteIf](<object.DeleteIf.md>),
[object.Erase](<object.Erase.md>),
[object.Remove](<object.Remove.md>),
[object.RemoveIf](<object.RemoveIf.md>),
[object.Trim!](<object.Trim!.md>),
[object.WithoutFields](<object.WithoutFields.md>)
