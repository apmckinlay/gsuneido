### NameArgs

``` suneido
(args, names, defaults = #())
```

Does argument handling similar to what Suneido normally does, i.e. allowing named and unnamed arguments as well as defaults.

If there are less defaults than names, they apply to the names at the *end*.

This is useful when you need to accept a variable number of arguments with @args, but still want to get specific arguments which might have been passed named or unnamed or not passed but given a default value. For example:

``` suneido
test = function (@args)
    { NameArgs(args, #(a, b, c), #(2, 3)) }

test()
    => missing argument: a

test(1)
    => [a: 1, b: 2, c: 3]

test(11, 22, 33)
    => [a: 11, b: 22, c: 33]

test(11, c: 33)
    => [a: 11, b: 2, c: 33]
```


See also:
[object.Extract](<Object/object.Extract.md>),
[object.ListToMembers](<Object/object.ListToMembers.md>),
[object.ListToNamed](<Object/object.ListToNamed.md>),
[object.Project](<Object/object.Project.md>),
[object.ProjectValues](<Object/object.ProjectValues.md>),
object.Slice
