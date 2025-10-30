#### object.ListToNamed

``` suneido
(@names) => object
```

Converts a list of values to an object containing named values. If there are more list values than names they are ignored.

For example:

``` suneido
[1, 2, 3, 4].ListToNamed('a', 'b', 'c')
    => [a: 1, b: 2, c: 3]
```


See also:
[NameArgs](<../NameArgs.md>),
[object.Extract](<object.Extract.md>),
[object.ListToMembers](<object.ListToMembers.md>),
[object.Project](<object.Project.md>),
[object.ProjectValues](<object.ProjectValues.md>),
object.Slice
