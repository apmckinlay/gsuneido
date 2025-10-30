#### object.WithoutFields

``` suneido
(member ...) => object
```

Returns a new object without any of the specified members. It does not modify the original object.

For example:

``` suneido
Object(a: 12, b: 34, c: 56).Without('c', 'a')
    => #(b: 34)
```


See also:
[object.Delete](<object.Delete.md>),
[object.DeleteIf](<object.DeleteIf.md>),
[object.Erase](<object.Erase.md>),
[object.Remove](<object.Remove.md>),
[object.RemoveIf](<object.RemoveIf.md>),
[object.Trim!](<object.Trim!.md>),
[object.Without](<object.Without.md>)
