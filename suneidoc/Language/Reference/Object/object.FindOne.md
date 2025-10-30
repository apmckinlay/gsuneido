#### object.FindOne

``` suneido
(callable) => value or false
```

Searches the object for a value for which callable returns true.
Returns false if one is not found.

callable is called with a value as its single argument.

callable can be anything that can be called, e.g. block, function, class, or instance.

For example:

``` suneido
#((a: 1, b: 2), (a: 6, b: 7)).FindOne({ it.b.Odd?() })
    => #(a: 6, b: 7)
```

**Note:** Named members are not searched in any particular order.


See also:
[object.Any?](<object.Any?.md>),
[object.Every?](<object.Every?.md>),
[object.Find](<object.Find.md>),
[object.FindAll](<object.FindAll.md>),
[object.FindAllIf](<object.FindAllIf.md>),
[object.FindIf](<object.FindIf.md>),
[object.FindLastIf](<object.FindLastIf.md>),
[object.Has?](<object.Has?.md>),
[object.HasIf?](<object.HasIf?.md>)
