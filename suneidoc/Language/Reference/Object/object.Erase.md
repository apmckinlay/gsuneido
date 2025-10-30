<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Erase

``` suneido
(member ...) => object
```

Erases a member of an object.

``` suneido
Object(12, 34, a: 56, b: 78).Erase('a') => #(12, 34, b: 78)
```

Unlike Delete, if the member is within the un-named members, the following members are <u>not</u> moved down.

``` suneido
Object(12, 34, 56, 78).Erase(1) => #(12, 2: 56, 3: 78)
```

Note: Erase modifies the object it is applied to, it does not create a new object.


See also:
[object.Delete](<object.Delete.md>),
[object.DeleteIf](<object.DeleteIf.md>),
[object.Remove](<object.Remove.md>),
[object.RemoveIf](<object.RemoveIf.md>),
[object.Trim!](<object.Trim!.md>),
[object.Without](<object.Without.md>),
[object.WithoutFields](<object.WithoutFields.md>)
