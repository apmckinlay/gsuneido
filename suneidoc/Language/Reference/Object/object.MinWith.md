#### object.MinWith

``` suneido
(block) => value
```

Returns the value in the object that the block returns the smallest value for.

For example:

``` suneido
ob = #("one", "two", "three", "four")

ob.MinWith({|x| x.Size() }) => "one" // shortest string

ob.Min() => "four" // first alphabetically
```

**Note:** There must be at least one member.


See also:
[Cmp](<../Cmp.md>),
[Gt](<../Gt.md>),
[Min](<../Min.md>),
[Max](<../Max.md>),
[object.Min](<object.Min.md>),
[object.Max](<object.Max.md>),
[object.MaxWith](<object.MaxWith.md>),
[Suneido.StrictCompare](<../Suneido/Suneido.StrictCompare.md>)
