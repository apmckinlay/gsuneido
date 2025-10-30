### Join

``` suneido
(separator, value ...) => string
```

Concatenates the values with the separator between them.

Empty values ("") are ignored and do not cause extra separators.

This is equivalent to:

``` suneido
Object(value ...).Remove("").Join(separator)
```

For example:

``` suneido
a = "one"
b = ""
c = "two"
Join(", ", a, b, c)
    => "one, two"
```
See also:
[Opt](<Opt.md>),
[object.Join](<Object/object.Join.md>)