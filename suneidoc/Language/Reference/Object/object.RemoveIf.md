#### object.RemoveIf

``` suneido
(block) => this
```

Removes values from the object for which the block returns true. The block can be anything callable. It is passed the value.

For example:

``` suneido
Object(1, 2, a: 3, b: 4).RemoveIf({ it is 1 or it is 3 })
    => #(2, b: 4)
```

Uses object.Delete so within the un-named members any following members are shifted over.

**Note:** RemoveIf modifies the object it is applied to, it does not create a new object.


See also:
[object.Delete](<object.Delete.md>),
[object.DeleteIf](<object.DeleteIf.md>),
[object.Erase](<object.Erase.md>),
[object.Remove](<object.Remove.md>),
[object.Trim!](<object.Trim!.md>),
[object.Without](<object.Without.md>),
[object.WithoutFields](<object.WithoutFields.md>)
