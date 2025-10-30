### Curry

``` suneido
(func, arg ...) => callable
```

**Note**: In almost all cases it is simpler and more efficient to use a block.

Returns a new function with the initial arguments "bound" to the supplied ones.

Technically, this is not "currying" it is binding some of the arguments, which is something currying allows.

For example:

``` suneido
add1 = Curry(Add, 1)
```

is equivalent to:

``` suneido
add1 = function (x) { return Add(1, x) }
```

Note: The returned value is actually a class instance with a Call method.


See also:
[Compose](<Compose.md>),
[Memoize](<Memoize.md>),
[MemoizeSingle](<MemoizeSingle.md>),
[object.Any?](<Object/object.Any?.md>),
[object.Drop](<Object/object.Drop.md>),
[object.Every?](<Object/object.Every?.md>),
[object.Filter](<Object/object.Filter.md>),
[object.FlatMap](<Object/object.FlatMap.md>),
[object.Fold](<Object/object.Fold.md>),
[object.Map](<Object/object.Map.md>),
[object.Map!](<Object/object.Map!.md>),
[object.Map2](<Object/object.Map2.md>),
[object.Nth](<Object/object.Nth.md>),
object.PrevSeq,
[object.Reduce](<Object/object.Reduce.md>),
[object.Take](<Object/object.Take.md>),
[object.Zip](<Object/object.Zip.md>)
