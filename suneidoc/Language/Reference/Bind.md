<div style="float:right"><span class="builtin">Builtin</span></div>

### Bind

``` suneido
(func, arg ...) => callable
```

**Note**: In most cases it is simpler and more efficient to use a block.

Returns a new function with the initial arguments "bound" to the supplied ones.

Technically, this is not "currying" it is binding some of the arguments, which is something currying allows.

For example:

``` suneido
cmp0 = Bind(Cmp, 0)
cmp0(123)
=> 1

```

is equivalent to:

``` suneido
cmp0 = function (x) { return Cmp(0, x) }
```

Built-in since 2026-02-24


See also:
[Compose](<Compose.md>),
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
