#### object.AddMany!

``` suneido
(value, n) => this
```

Adds **value** to the object **n** times.

For example:

``` suneido
Object().AddMany!(9, 3)
    => #(9, 9, 9)
```

**Note:** AddMany! modifies the object it is applied to, it does not create a new object.

See also:
[object.Add](<object.Delete.md>)