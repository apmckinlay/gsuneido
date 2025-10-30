### Concat

``` suneido
(@iterables) => sequence
```

Returns a [sequence](<../Basic Data Types/Sequence.md>) containing the elements from the iterables concatenated together.

Concat can also be accessed as [object/sequence.Concat](<Object/object.Concat.md>)

For example:

``` suneido
p = Primes().Take(4)
f = Fibonaccis().Take(4)
Concat(p, f) // or p.Concat(f)
	=> #(2, 3, 5, 7, 0, 1, 1, 2)
```

See also: [object.Append](<Object/object.Append.md>)


See also:
[Drop](<Drop.md>),
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
