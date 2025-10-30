<div style="float:right"><span class="builtin">Builtin</span></div>

### Seq

``` suneido
(from = 0, to = int_max, by = 1) => sequence
```

Returns a *sequence* of values starting with **from** up to, but not including **to**, in increments of **by**.

It is comparable to:

``` suneido
i = from
while (i < to)
    i += by
```

For example:

``` suneido
Seq(3, 15, 3) => #(3, 6, 9, 12)
```

If only a single value is passed, it is taken as **to**, and **from** is taken as 0. For example:

``` suneido
Seq(4) // same as Seq(0, 4)
    => #(0, 1, 2, 3)
```

If no arguments are supplied, the result is an infinite sequence starting at zero. Attempts to instantiate will throw an error and Display will not attempt to show the contents.

The returned sequence is initially "virtual", i.e. it is simply an iterator over the original data. If you only iterate over the sequence (e.g. using a for-in loop) then an object containing all the values is never instantiated. However, if you access the sequence in most other ways (e.g. call a built-in object method on it) then a list will be created with the values from the sequence and that list is what will be used for any future operations. See also: [Basic Data Types > Sequence](<../Basic Data Types/Sequence.md>)

Note: When testing on the WorkSpace, the variable display will trigger instantiation of any sequences in local variables.

For example, the following does **not** generate an object with 10,000 values in it:

``` suneido
Seq(10000).Iter().Next()
    => 0
```


See also:
[Drop](<Drop.md>),
[FileLines](<FileLines.md>),
[Filter](<Filter.md>),
[Nof](<Nof.md>),
[Grep](<Grep.md>),
[Map](<Map.md>),
[Map2](<Map2.md>),
[Sequence](<Sequence.md>),
[string.Lines](<String/string.Lines.md>),
[Take](<Take.md>)
