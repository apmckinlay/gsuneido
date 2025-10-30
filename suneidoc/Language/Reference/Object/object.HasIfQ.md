#### object.HasIf?

``` suneido
(block) => true or false
```

Returns true if the object contains a *value* for which the block returns true, false if not.

For example:

``` suneido
#(123, "abc", name: "Fred", age: 23).HasIf?({ String?(it) and it.Capitalized?() }) 
    => true
```

**Note**: object.HasIf? is identical to [object.Any?](<object.Any?.md>)


See also:
[object.Any?](<object.Any?.md>),
[object.Every?](<object.Every?.md>),
[object.Find](<object.Find.md>),
[object.FindAll](<object.FindAll.md>),
[object.FindAllIf](<object.FindAllIf.md>),
[object.FindIf](<object.FindIf.md>),
[object.FindLastIf](<object.FindLastIf.md>),
[object.FindOne](<object.FindOne.md>),
[object.Has?](<object.Has?.md>)
