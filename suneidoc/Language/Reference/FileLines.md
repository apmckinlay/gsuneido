### FileLines

``` suneido
(file, block)
```

FileLines opens the file and passes a sequence of lines to the block.

For example:

``` suneido
FileLines(file)
	{ it.Filter(...).Map(...) }
```

This is somewhat equivalent to GetFile(file).Lines() except FileLines does not read the entire file into memory, it only reads a line at a time.


See also:
[Drop](<Drop.md>),
[Filter](<Filter.md>),
[Nof](<Nof.md>),
[Grep](<Grep.md>),
[Map](<Map.md>),
[Map2](<Map2.md>),
[Seq](<Seq.md>),
[Sequence](<Sequence.md>),
[string.Lines](<String/string.Lines.md>),
[Take](<Take.md>)
