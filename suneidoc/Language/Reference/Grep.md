### Grep

``` suneido
(iterable, regex, block = false) => sequence
```

Returns a [sequence](<../Basic Data Types/Sequence.md>) of the values that match the regular expression. The input sequence must be strings. Grep is normally accessed as object/sequence.Grep

For example:

``` suneido
#(one, two, three, four, five).Grep('e')
=> #("one", "three", "five")
```

If a block is supplied, it is called with the index in the input sequence and the input value. The block's return values will be the output sequence.

For example:

``` suneido
#(one, two, three, four, five).Grep('e', {|i, x| i $ ':' $ x })
=> #("0:one", "2:three", "4:five")
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
[Map](<Map.md>),
[Map2](<Map2.md>),
[Seq](<Seq.md>),
[Sequence](<Sequence.md>),
[string.Lines](<String/string.Lines.md>),
[Take](<Take.md>)
