<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Set_readonly

``` suneido
() => this
```

Set the object to read-only i.e. so it can no longer be modified.

If an object (or record) will be shared between concurrent threads it is a good idea to make it read-only.

For example:

``` suneido
ob = Object()
ob.num = 123
ob.Set_readonly()
ob.num = 456
    ERROR: can't modify readonly objects
```

**Note**: Set_readonly is recursively applied to any nested objects or records.

**Note**: When a record is read-only rules will still work, but their results can not be saved so they will be evaluated every time they are referenced.