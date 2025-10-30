### Nof

``` suneido
(n, block) => sequence
```

Nof produces a sequence of values from calling block n times. Normally accessed with `number.Of(block)`

For example:

``` suneido
5.Of(Random)
=> #(9430793933001491, 5509763035817541, 7347549908506805, 492112129290015, 3229830691105010)

8.Of({ "abcdef".RandChar() }).Join()
=> "feefcecf"
```

`Nof({ ... })` is equivalent to `Seq(n).Map(|unused| ... })`


See also:
[Drop](<Drop.md>),
[FileLines](<FileLines.md>),
[Filter](<Filter.md>),
[Grep](<Grep.md>),
[Map](<Map.md>),
[Map2](<Map2.md>),
[Seq](<Seq.md>),
[Sequence](<Sequence.md>),
[string.Lines](<String/string.Lines.md>),
[Take](<Take.md>)
