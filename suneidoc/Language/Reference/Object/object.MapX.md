#### object.Map!

``` suneido
(callable) => this
```

For each member, Map! does:

``` suneido
[member] = callable([member])
```

callable can be anything that can be called e.g. block, function, class, or instance.

**Note:** Map! modifies the object it is applied to, it does not create a new object. It cannot be applied to a read-only object. Whereas [object.Map](<object.Map.md>) creates a new object and can be applied to a read-only object.

For example:

``` suneido
Object(1, 2, 3).Map!({|x| x * 2 }) => #(2, 4, 6)
```


See also:
[Compose](<../Compose.md>),
[Curry](<../Curry.md>),
[Memoize](<../Memoize.md>),
[MemoizeSingle](<../MemoizeSingle.md>),
[object.Any?](<object.Any?.md>),
[object.Drop](<object.Drop.md>),
[object.Every?](<object.Every?.md>),
[object.Filter](<object.Filter.md>),
[object.FlatMap](<object.FlatMap.md>),
[object.Fold](<object.Fold.md>),
[object.Map](<object.Map.md>),
[object.Map2](<object.Map2.md>),
[object.Nth](<object.Nth.md>),
object.PrevSeq,
[object.Reduce](<object.Reduce.md>),
[object.Take](<object.Take.md>),
[object.Zip](<object.Zip.md>)
