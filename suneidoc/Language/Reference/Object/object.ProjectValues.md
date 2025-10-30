#### object.ProjectValues

``` suneido
(members) => object
```

Returns a new object listing the values of the specified members.

The members may be passed as an object or as individual arguments.

For example:

``` suneido
#(a: 12, b: 34, c: 56, d: 78).ProjectValues(#(b, c)) => #(34, 56)

#(a: 12, b: 34, c: 56, d: 78).ProjectValues(#b, #c) => #(34, 56)

#(12, 34, 56, 78, 90).ProjectValues(#(0, 2, 4)) => #(12, 56, 90)
```

**Note:** ProjectValues will throw an error if the object does not contain all the specified members.


See also:
[NameArgs](<../NameArgs.md>),
[object.Extract](<object.Extract.md>),
[object.ListToMembers](<object.ListToMembers.md>),
[object.ListToNamed](<object.ListToNamed.md>),
[object.Project](<object.Project.md>),
object.Slice
