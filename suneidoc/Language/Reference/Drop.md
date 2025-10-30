### Drop

``` suneido
(iterable, n) => sequence
```

Returns a [sequence](<../Basic Data Types/Sequence.md>) containing the elements from the iterable, skipping the first **n** elements. Drop is normally accessed as object/sequence.Drop

For example:

``` suneido
Seq(10).Drop(5)
    => #(5, 6, 7, 8, 9)
```


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



See also:
[FileLines](<FileLines.md>),
[Filter](<Filter.md>),
[Nof](<Nof.md>),
[Grep](<Grep.md>),
[Map](<Map.md>),
[Map2](<Map2.md>),
[Seq](<Seq.md>),
[Sequence](<Sequence.md>),
[string.Lines](<String/string.Lines.md>),
[Take](<Take.md>)
