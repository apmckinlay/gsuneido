#### object.FindLastIf

``` suneido
(callable) => member or false
```

Searches for the last unnamed list value for which callable returns true.
Returns false if the value is not found.

callable is called with a value as its single argument.

callable can be anything that can be called, e.g. block, function, class, or instance.

For example:

``` suneido
#(38, 48, 32, 67, 89).FindLastIf({|x| x < 50}) => 2
```

**Note:** Named members are not searched.


See also:
[object.Any?](<object.Any?.md>),
[object.Every?](<object.Every?.md>),
[object.Find](<object.Find.md>),
[object.FindAll](<object.FindAll.md>),
[object.FindAllIf](<object.FindAllIf.md>),
[object.FindIf](<object.FindIf.md>),
[object.FindOne](<object.FindOne.md>),
[object.Has?](<object.Has?.md>),
[object.HasIf?](<object.HasIf?.md>)
