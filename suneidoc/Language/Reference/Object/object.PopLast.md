<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.PopLast

``` suneido
() => value
```

Deletes the last list member and returns its value.

If the list is empty, it returns the object itself.

For example:

``` suneido
while not Same?(ob, last = ob.PopLast())
    {
    ...
    }
```

PopLast is built-in to make it atomic. Checking the size before popping is not thread-safe because another thread could make changes in between.

See also: 
[object.Last](<object.Last.md>),
[object.PopFirst](<object.PopFirst.md>)