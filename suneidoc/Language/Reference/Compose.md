### Compose

``` suneido
(func, func ...) => callable
```

Compose is a convenience function to avoid the deep nesting of parenthesis from nested calls.

For example:

``` suneido
f = Compose(a, b, c)
```

is equivalent to:

``` suneido
f = function (@args) { c(b(a(@args))) }
```

Notice that the functions are passed in the order they will be applied (i.e. inside out) in the same way you would write a *nix pipeline a | b | c

Note: The returned value is actually a block.


See also:
[Curry](<Curry.md>),
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
