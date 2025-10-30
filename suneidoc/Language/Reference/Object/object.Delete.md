<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Delete

``` suneido
(member ...) => this
(all:) => this
```

Deletes specified or all members of an object.

``` suneido
Object(12, 34, a: 56, b: 78).Delete('a') => #(12, 34, b: 78)

Object(12, 34, a: 56, b: 78).Delete(all:) => #()
```

If the members are within the un-named members, any following members are moved down.

``` suneido
Object(12, 34, 56, 78).Delete(2, 1) => #(12, 78)
```

**Warning:** If you are deleting multiple un-named members, they should be listed in reverse order (as in the above example) so earlier deletes do not invalidate the indexes of later ones.

**Note:** Delete modifies the object it is applied to, it does not create a new object.


See also:
[object.DeleteIf](<object.DeleteIf.md>),
[object.Erase](<object.Erase.md>),
[object.Remove](<object.Remove.md>),
[object.RemoveIf](<object.RemoveIf.md>),
[object.Trim!](<object.Trim!.md>),
[object.Without](<object.Without.md>),
[object.WithoutFields](<object.WithoutFields.md>)
