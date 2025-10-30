<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.PopFirst

``` suneido
() => value
```

Deletes member 0 and returns its value.

If the list is empty, it returns the object itself.

For example:

``` suneido
while not Same?(ob, first = ob.PopFirst())
    {
    ...
    }
```

PopFirst is built-in to make it atomic. Checking the size before popping is not thread-safe because another thread could make changes in between.

See also: 
[object.First](<object.First.md>),
[object.PopLast](<object.PopLast.md>)