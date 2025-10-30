#### object.Replace

``` suneido
(oldvalue, newvalue) => this
```

Replaces all occurrences of *oldvalue* with *newvalue*.

For example:

``` suneido
Object(1, 2, a: 1, b: 2).Replace(2, 22) => #(1, 22, a: 1, b: 22)
```

**Note:** Replace modifies the object it is applied to, it does not create a new object.