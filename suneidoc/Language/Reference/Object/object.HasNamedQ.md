#### object.HasNamed?

``` suneido
() => true or false
```

Returns true if the object has named members, false if not.

For example:

``` suneido
#().HasNamed?() => false
#(12, 34, 56).HasNamed?() => false
#(12, 34, name: "Fred", age: 23).HasNamed?() => true
#(name: "Fred", age: 23).HasNamed?() => true
```