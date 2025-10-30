<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.DeleteIf

``` suneido
(block) => this
```

Deletes members of the object for which the block returns true. block can be anything callable. It is passed a member name.

``` suneido
Object(a: 1, B: 2, c: 3).DeleteIf({ it.Upper?() })
    => #(a: 1, c: 3)
```

Uses object.Delete so within the un-named members any following members are shifted over.

``` suneido
Object(12, 34, 56, 78).DeleteIf({ it is 1 or it is 2 })
    => #(12, 78)
```

**Note:** DeleteIf modifies the object it is applied to, it does not create a new object.


See also:
[object.Delete](<object.Delete.md>),
[object.Erase](<object.Erase.md>),
[object.Remove](<object.Remove.md>),
[object.RemoveIf](<object.RemoveIf.md>),
[object.Trim!](<object.Trim!.md>),
[object.Without](<object.Without.md>),
[object.WithoutFields](<object.WithoutFields.md>)
