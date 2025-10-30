<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Find

``` suneido
(value) => member or false
```

Searches the members of an object for a value. Returns false if the value is not found.

For example:

``` suneido
Object(38, 48, 32, 67, 89).Find(32) => 2
Object(12, 34, a: 56, b: 78).Find(56) => "a"
Object(12, 34, a: 56, b: 78).Find(90) => false
```

If the value is in multiple un-named members, the first one will be returned. For example:

``` suneido
#(1, 2, 2, 3).Find(2) => 1
```

**Note:** If the value is in multiple *named* members it is undefined which one will be returned.


See also:
[object.Any?](<object.Any?.md>),
[object.Every?](<object.Every?.md>),
[object.FindAll](<object.FindAll.md>),
[object.FindAllIf](<object.FindAllIf.md>),
[object.FindIf](<object.FindIf.md>),
[object.FindLastIf](<object.FindLastIf.md>),
[object.FindOne](<object.FindOne.md>),
[object.Has?](<object.Has?.md>),
[object.HasIf?](<object.HasIf?.md>)
