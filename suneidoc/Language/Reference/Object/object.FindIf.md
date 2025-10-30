#### object.FindIf

``` suneido
(callable) => member or false
```

Searches for the first value for which callable returns true.
Returns false if the value is not found.

callable is called with a value as its single argument.

callable can be anything that can be called, e.g. block, function, class, or instance.

For example:

``` suneido
#(38, 48, 32, 67, 89).FindIf({|x| x > 50}) => 3

#(12, 34, a: 56, b: 78).FindIf(function (x) { x > 50 }) => "a"

#((1, 2), (3, 4) (5, 6) (7, 8)).FindIf({|x| x[0] + x[1] > 8}) => 2
```

**Note:** Named members are not searched in any particular order.


See also:
[object.Any?](<object.Any?.md>),
[object.Every?](<object.Every?.md>),
[object.Find](<object.Find.md>),
[object.FindAll](<object.FindAll.md>),
[object.FindAllIf](<object.FindAllIf.md>),
[object.FindLastIf](<object.FindLastIf.md>),
[object.FindOne](<object.FindOne.md>),
[object.Has?](<object.Has?.md>),
[object.HasIf?](<object.HasIf?.md>)
