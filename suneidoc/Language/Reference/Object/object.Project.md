#### object.Project

``` suneido
(members) => object
```

Returns a new object containing the specified members.

The members may be passed as an object or as individual arguments.

For example:

``` suneido
#(a: 12, b: 34, c: 56, d: 78).Project(#(b, c)) => #(b: 34, c: 56)

#(a: 12, b: 34, c: 56, d: 78).Project(#b, #c) => #(b: 34, c: 56)
```

**Note:** Project will cause an error if the object does not contain all the specified members. (Unless the object has a default value.)


See also:
[NameArgs](<../NameArgs.md>),
[object.Extract](<object.Extract.md>),
[object.ListToMembers](<object.ListToMembers.md>),
[object.ListToNamed](<object.ListToNamed.md>),
[object.ProjectValues](<object.ProjectValues.md>),
object.Slice
