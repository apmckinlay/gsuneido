### By

``` suneido
(@fields)
```

Returns a block that can be used as the comparison function for [object.Sort!](<Object/object.Sort!.md>) to sort a list of objects or records by one or more fields. For example:

``` suneido
list = [[name: 'Fred', age: 23], [name: 'Andy', age: 45]]
list.Sort!(By(#age, #name))

=> [[name: "Fred", age: 23], [name: "Andy", age: 45]]
```


See also:
[object.Sort!](<Object/object.Sort!.md>),
[object.Sorted?](<Object/object.Sorted?.md>),
[object.SortWith!](<Object/object.SortWith!.md>)
