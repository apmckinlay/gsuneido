#### object.ListToMembers

``` suneido
(@names) => object
```

Converts a list of values to an object with those values as the member "names".

This can be useful because object.Member?() is faster than object.Has?()

For example:

``` suneido
["one", "two", "three"].ListToMembers()
    => [one:, two:, three:]
```


See also:
[NameArgs](<../NameArgs.md>),
[object.Extract](<object.Extract.md>),
[object.ListToNamed](<object.ListToNamed.md>),
[object.Project](<object.Project.md>),
[object.ProjectValues](<object.ProjectValues.md>),
object.Slice
