#### object.Zip

``` suneido
() => list
```

Takes a list of lists and returns a new list of lists by grouping the i'th element from each of the original lists.  If the lists are not the same length, the "extra" elements of the longer lists are ignored (dropped).

Note: If the sub-lists are the same length, then Zip is reversible.

For example:

``` suneido
a = #(1, 2, 3, 4)
b = #(a, b, c, d)
z = Object(a, b).Zip()
    => #(#(1, "a"), #(2, "b"), #(3, "c"), #(4, "d"))
z.Zip()
    => #(#(1, 2, 3, 4), #("a", "b", "c", "d"))
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
[object.Map!](<object.Map!.md>),
[object.Map2](<object.Map2.md>),
[object.Nth](<object.Nth.md>),
object.PrevSeq,
[object.Reduce](<object.Reduce.md>),
[object.Take](<object.Take.md>)
