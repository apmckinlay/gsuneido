#### object.AddUnique

``` suneido
(value) => this
```

Add a value to an object only if the value is not already in the object.  For example:

``` suneido
Object(12, 34).AddUnique(34) => Object(12, 34)
```

Note: AddUnique modifies the object it is applied to, it does not create a new object.