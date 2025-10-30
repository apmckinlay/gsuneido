### Map

``` suneido
(iterable, block) => sequence
```

Returns a [sequence](<../Basic Data Types/Sequence.md>) that applies the block to each element. Map does **not** modify the original object. Map is normally accessed as object/sequence.Map

For example:

``` suneido
Seq(4).Map({ it * it })
    => #(0, 1, 4, 9)
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
[Drop](<Drop.md>),
[FileLines](<FileLines.md>),
[Filter](<Filter.md>),
[Nof](<Nof.md>),
[Grep](<Grep.md>),
[Map2](<Map2.md>),
[Seq](<Seq.md>),
[Sequence](<Sequence.md>),
[string.Lines](<String/string.Lines.md>),
[Take](<Take.md>)
