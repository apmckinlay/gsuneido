#### object.FlatMap

``` suneido
(block) => object
```

Applies the block to each value in this. The callable must return a list of results. These results are combined into a single list similar to [object.Flatten](<object.Flatten.md>).

For example:

``` suneido
twice = function (x) { return [x, x] }
twice(1)
    => [1, 1]
[1, 2].FlatMap(twice)
    => [1, 1, 2, 2]
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
[object.Fold](<object.Fold.md>),
[object.Map](<object.Map.md>),
[object.Map!](<object.Map!.md>),
[object.Map2](<object.Map2.md>),
[object.Nth](<object.Nth.md>),
object.PrevSeq,
[object.Reduce](<object.Reduce.md>),
[object.Take](<object.Take.md>),
[object.Zip](<object.Zip.md>)
