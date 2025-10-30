#### object.Map2

``` suneido
(block) => object
```

See [Map2](<../Map2.md>)

Note: If there are named members, object.Map2 will include and sort them.

For example:

``` suneido
#(99, c: 12, a: 34, b: 56).Map2({|m,v| m $ ' is ' $ v })
    => #("0 is 99", "a is 34", "b is 56", "c is 12")
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
[object.Nth](<object.Nth.md>),
object.PrevSeq,
[object.Reduce](<object.Reduce.md>),
[object.Take](<object.Take.md>),
[object.Zip](<object.Zip.md>)



See also:
[Drop](<../Drop.md>),
[FileLines](<../FileLines.md>),
[Filter](<../Filter.md>),
[Nof](<../Nof.md>),
[Grep](<../Grep.md>),
[Map](<../Map.md>),
[Map2](<../Map2.md>),
[Seq](<../Seq.md>),
[Sequence](<../Sequence.md>),
[string.Lines](<../String/string.Lines.md>),
[Take](<../Take.md>)
