### Filter

``` suneido
(iterable, block) => sequence
```

Returns a [sequence](<../Basic Data Types/Sequence.md>) containing the values for which block returns true. Filter is normally accessed as object/sequence.Filter

block is called with the value as its single argument

For example:

``` suneido
#(1, 2, 3, 4, 5).Filter({ it > 3 })
    => #(4, 5)
#(1, 2, 3, 4, 5).Filter(#Even?)
    => #(2, 4)
```

**Note:** Member names are ignored and will not be carried over to the result.


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
[Nof](<Nof.md>),
[Grep](<Grep.md>),
[Map](<Map.md>),
[Map2](<Map2.md>),
[Seq](<Seq.md>),
[Sequence](<Sequence.md>),
[string.Lines](<String/string.Lines.md>),
[Take](<Take.md>)
