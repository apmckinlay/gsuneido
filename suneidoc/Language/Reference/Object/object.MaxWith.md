#### object.MaxWith

``` suneido
(block) => value
```

Returns the value in the object that the block returns the largest value for.

For example:

``` suneido
ob = #("one", "two", "three", "four")

ob.MaxWith({|x| x.Size() }) => "three" // longest string

ob.Max() => "two" // last alphabetically
```

**Note:** There must be at least one member.


See also:
[Cmp](<../Cmp.md>),
[Gt](<../Gt.md>),
[Min](<../Min.md>),
[Max](<../Max.md>),
[object.Min](<object.Min.md>),
[object.MinWith](<object.MinWith.md>),
[object.Max](<object.Max.md>),
[Suneido.StrictCompare](<../Suneido/Suneido.StrictCompare.md>)
