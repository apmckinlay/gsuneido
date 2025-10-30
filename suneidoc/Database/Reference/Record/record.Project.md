#### record.Project

``` suneido
(members) => record
```

Returns a new record containing the specified members.

The members may be passed as an object or as individual arguments.

For example:

``` suneido
[a: 12, b: 34, c: 56, d: 78].Project(#(b, c)) => [b: 34, c: 56]

[a: 12, b: 34, c: 56, d: 78].Project(#b, #c) => [b: 34, c: 56]
```

record.Project is very similar to [object.Project](<../../../Language/Reference/Object/object.Project.md>) except it returns a record.