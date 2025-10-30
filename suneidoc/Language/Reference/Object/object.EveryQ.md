#### object.Every?

``` suneido
(callable) => true / false
```

Returns true if callable returns true for <u>every</u> value in the object, false otherwise.

callable must return true or false. It can be anything that can be called e.g. block, function, class, or instance.

For example:

``` suneido
#(2, 4, 6, 8).Every?({ it % 2 is 0 })
    => true
#(2, 4, 5, 8).Every?({ it % 2 is 0 })
    => false
```


See also:
[object.Any?](<object.Any?.md>),
[object.Find](<object.Find.md>),
[object.FindAll](<object.FindAll.md>),
[object.FindAllIf](<object.FindAllIf.md>),
[object.FindIf](<object.FindIf.md>),
[object.FindLastIf](<object.FindLastIf.md>),
[object.FindOne](<object.FindOne.md>),
[object.Has?](<object.Has?.md>),
[object.HasIf?](<object.HasIf?.md>)



See also:
[Compose](<../Compose.md>),
[Curry](<../Curry.md>),
[Memoize](<../Memoize.md>),
[MemoizeSingle](<../MemoizeSingle.md>),
[object.Any?](<object.Any?.md>),
[object.Drop](<object.Drop.md>),
[object.Filter](<object.Filter.md>),
[object.FlatMap](<object.FlatMap.md>),
[object.Fold](<object.Fold.md>),
[object.Map](<object.Map.md>),
[object.Map!](<object.Map!.md>),
[object.Map2](<object.Map2.md>),
[object.Nth](<object.Nth.md>),
object.PrevSeq,
[object.Reduce](<object.Reduce.md>),
[object.Take](<object.Take.md>),
[object.Zip](<object.Zip.md>)
