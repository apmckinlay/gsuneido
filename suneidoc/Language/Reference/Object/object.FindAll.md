<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.FindAll

``` suneido
(value) => list
```

Returns a list of all the members containing the given value.

For example:

``` suneido
Object(1, 2, a: 1, b: 2).FindAll(1)
    => #(0, a)
```


See also:
[object.Any?](<object.Any?.md>),
[object.Every?](<object.Every?.md>),
[object.Find](<object.Find.md>),
[object.FindAllIf](<object.FindAllIf.md>),
[object.FindIf](<object.FindIf.md>),
[object.FindLastIf](<object.FindLastIf.md>),
[object.FindOne](<object.FindOne.md>),
[object.Has?](<object.Has?.md>),
[object.HasIf?](<object.HasIf?.md>)
