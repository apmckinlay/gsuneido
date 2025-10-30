#### object.Any?

``` suneido
(callable) => true / false
```

Returns true if callable returns true for <u>any</u> value in the object, false otherwise.

callable must return true or false. It can be anything that can be called e.g. block, function, class, or instance.

For example:

``` suneido
#(a, b, C, d).Any?({ it.Upper?() })
    => true
#(a, b, c, d).Any?({ it.Upper?() })
    => false
```

**Note**: object.Any? is identical to [object.HasIf?](<object.HasIf?.md>)


See also:
[object.Every?](<object.Every?.md>),
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
[object.Take](<object.Take.md>),
[object.Zip](<object.Zip.md>)
