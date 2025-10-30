<div style="float:right"><span class="builtin">Builtin</span></div>

### Sequence

``` suneido
(iterator) => sequence
```

Sequence wraps an "iterator", an instance of a Suneido class with Next() and Dup() methods so that it can be used with [for-in](<../Statements/for.md>) loops and so it can be treated interchangeably as an iterator or as a list of values. (Similar to how [Seq](<Seq.md>), [object.Members](<Object/object.Members.md>), [object.Values](<Object/object.Values.md>), and [object.Assocs](<Object/object.Assocs.md>) work.)

Sequence is necessary because Iter is a built-in method. Even if you define an Iter method on a class it will not be used (e.g. by for-in) because the built-in method takes precedence. (In Suneido built-in methods cannot be overridden by user defined classes.)

The Next() method of the supplied iterator should return the next value or the iterator itself (i.e. `this`) to signal there are no more values.

The Dup() method should return a new iterator for the sequence. If the sequence wraps another iterator, it should call Dup() on that iterator. (i.e. it should be a "deep" copy)

The Infinite?() method should return true or false. If an iterator returns true from Infinite?() then attempts to instantiate it will throw an error and [Display](<Display.md>) will not attempt to show the contents. This prevents getting into infinite loops. (e.g. see Primes) If an iterator wraps another iterator (e.g. Map, Filter) then Infinite?() should take into account its iterator's Infinite?()

For an example of a sequence source see Lines, for examples that read from a sequence and produce a new sequence see Map, Filter, Take, and Drop. Consuming a sequence is simply a matter of iterating it either with for-in or an explicit Iter.

Sequences methods can be called on Sequence's without forcing instantiation. And Sequences inherits from Objects, so this includes Objects methods. Although, some Objects methods may themselves force instantiation.

The returned sequence is initially "virtual", i.e. it is simply an iterator over the original data. If you only iterate over the sequence (e.g. using a for-in loop) then an object containing all the values is never instantiated. However, if you access the sequence in most other ways (e.g. call a built-in object method on it) then a list will be created with the values from the sequence and that list is what will be used for any future operations. See also: [Basic Data Types > Sequence](<../Basic Data Types/Sequence.md>)

Note: When testing on the WorkSpace, the variable display will trigger instantiation of any sequences in local variables.


See also:
[Drop](<Drop.md>),
[FileLines](<FileLines.md>),
[Filter](<Filter.md>),
[Nof](<Nof.md>),
[Grep](<Grep.md>),
[Map](<Map.md>),
[Map2](<Map2.md>),
[Seq](<Seq.md>),
[string.Lines](<String/string.Lines.md>),
[Take](<Take.md>)
