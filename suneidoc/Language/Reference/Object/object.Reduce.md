#### object.Reduce

``` suneido
(block)
```

Apply a function of two arguments cumulatively to the values in a list, from left to right, reducing the sequence to a single value. Similar to [object.Fold](<object.Fold.md>) with the first element as the initial value.

**Note**: The list must have at least one value. If a list has a single value, that value will be the result.

For example: (although in practice you would use [object.Sum](<object.Sum.md>))

``` suneido
#(1,2,3,4,5).Reduce({|x,y| x + y }) 
    => 15
```

i.e. ((((1 + 2) + 3) + 4) + 5)


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
[object.Map!](<object.Map!.md>),
[object.Map2](<object.Map2.md>),
[object.Nth](<object.Nth.md>),
object.PrevSeq,
[object.Take](<object.Take.md>),
[object.Zip](<object.Zip.md>)
