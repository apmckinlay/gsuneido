#### object.Append

``` suneido
(object) => this
```

Adds the contents of the other object to this object. [object.Add](<object.Add.md>) is used for un-named values. Named values are simply set. For example:

``` suneido
Object(1, 2, a: 3, b: 4).Append(#(9, b: 5, c: 6))
    => Object(1, 2, 9, a: 3, b: 5, c: 6)
```

Note: Append modifies the object it is applied to, it does not create a new object.

For sequences, use [Concat](<../Concat.md>)

See also: 
[object.Merge](<object.Merge.md>)